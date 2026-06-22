package cu

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/navisim/rdnaemu"
	"gitlab.com/akita/navisim/rdnainsts"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/pipelining"
	"gitlab.com/akita/util/tracing"
)

// A Builder can construct a fully functional Compute Unit.
type Builder struct {
	engine                 akita.Engine
	freq                   akita.Freq
	name                   string
	simdCount              int
	schedulerCount         int
	vgprCount              []int
	sgprCount              []int
	log2CachelineSize      uint64
	numSinglePrecisionUnit int

	scratchpadPreparer ScratchpadPreparer
	alu                rdnaemu.ALU

	visTracer        tracing.Tracer
	enableVisTracing bool
}

// MakeBuilder returns a default builder object
func MakeBuilder() Builder {
	var b Builder
	b.freq = 1000 * akita.MHz
	b.simdCount = 4
	b.schedulerCount = 4
	b.sgprCount = []int{2560, 2560, 2560, 2560}
	b.vgprCount = []int{32768, 32768, 32768, 32768}
	b.log2CachelineSize = 7
	b.numSinglePrecisionUnit = 32

	return b
}

// WithEngine sets the engine to use.
func (b Builder) WithEngine(engine akita.Engine) Builder {
	b.engine = engine
	return b
}

// WithFreq sets the frequency.
func (b Builder) WithFreq(f akita.Freq) Builder {
	b.freq = f
	return b
}

// WithSIMDCount sets the number of SIMD unit in the ComputeUnit.
func (b Builder) WithSIMDCount(n int) Builder {
	b.simdCount = n
	b.schedulerCount = n
	return b
}

// WithVGPRCount sets the number of VGPRs associated with each SIMD Unit.
func (b Builder) WithVGPRCount(counts []int) Builder {
	if len(counts) != b.simdCount {
		panic("counts must have a length that equals to the SIMD count")
	}

	b.vgprCount = counts
	return b
}

// WithSGPRCount equals the number of SGPRs in the Compute Unit.
func (b Builder) WithSGPRCount(counts []int) Builder {
	if len(counts) != b.simdCount {
		panic("counts must have a length that equals to the SIMD count")
	}
	b.sgprCount = counts
	return b
}

// WithLog2CachelineSize sets the cacheline size as a power of 2.
func (b Builder) WithLog2CachelineSize(n uint64) Builder {
	b.log2CachelineSize = n
	return b
}

// WithNumSinglePrecisionUnit sets the number of lanes per SIMD.
func (b Builder) WithNumSinglePrecisionUnit(n int) Builder {
	b.numSinglePrecisionUnit = n
	return b
}

// WithVisTracer adds a tracer to the builder.
func (b Builder) WithVisTracer(t tracing.Tracer) Builder {
	b.enableVisTracing = true
	b.visTracer = t
	return b
}

// Build returns a newly constructed compute unit according to the
// configuration.
func (b *Builder) Build(name string) *ComputeUnit {
	b.name = name
	cu := NewComputeUnit(name, b.engine)
	cu.Freq = b.freq
	cu.Decoder = rdnainsts.NewDisassembler()
	cu.WfDispatcher = NewWfDispatcher(cu)
	cu.InFlightVectorMemAccessLimit = 512
	cu.log2CacheLineSize = b.log2CachelineSize
	cu.numSinglePrecisionUnit = b.numSinglePrecisionUnit

	b.alu = rdnaemu.NewALU(nil)
	b.scratchpadPreparer = NewScratchpadPreparerImpl(cu)

	for i := 0; i < 4; i++ {
		cu.WfPools = append(cu.WfPools, NewWavefrontPool(20))
	}

	b.equipSchedulers(cu)
	b.equipScalarUnits(cu)
	b.equipSIMDUnits(cu)
	b.equipLDSUnit(cu)
	b.equipVectorMemoryUnit(cu)
	b.equipRegisterFiles(cu)

	return cu
}

func (b *Builder) equipSchedulers(cu *ComputeUnit) {
	for i := 0; i < b.schedulerCount; i++ {
		name := fmt.Sprintf(b.name+".Scheduler%d", i)
		fetchArbitor := NewFetchArbiter(cu.WfPools[i])
		fetchArbitor.InstBufByteSize = 256
		issueArbitor := NewIssueArbiter(cu.WfPools[i])
		scheduler := NewScheduler(name, cu, fetchArbitor, issueArbitor, cu.WfPools[i])
		cu.Schedulers = append(cu.Schedulers, scheduler)
	}
}

func (b *Builder) equipScalarUnits(cu *ComputeUnit) {
	for i := 0; i < b.simdCount; i++ {
		branchUnit := NewBranchUnit(cu, b.scratchpadPreparer, b.alu)
		scalarDecoder := NewDecodeUnit(cu)

		scalarUnit := NewScalarUnit(cu, b.scratchpadPreparer, b.alu)
		scalarUnit.log2CachelineSize = b.log2CachelineSize
		scalarDecoder.AddExecutionUnit(scalarUnit)
		cu.BranchUnits = append(cu.BranchUnits, branchUnit)
		cu.ScalarUnits = append(cu.ScalarUnits, scalarUnit)
		cu.ScalarDecoders = append(cu.ScalarDecoders, scalarDecoder)
	}
}

func (b *Builder) equipSIMDUnits(cu *ComputeUnit) {
	for i := 0; i < b.simdCount; i++ {
		vectorDecoder := NewDecodeUnit(cu)
		name := fmt.Sprintf(b.name+".SIMD%d", i)
		simdUnit := NewSIMDUnit(cu, name, b.scratchpadPreparer, b.alu)
		if b.enableVisTracing {
			tracing.CollectTrace(simdUnit, b.visTracer)
		}
		vectorDecoder.AddExecutionUnit(simdUnit)
		cu.SIMDUnits = append(cu.SIMDUnits, simdUnit)
		cu.VectorDecoders = append(cu.VectorDecoders, vectorDecoder)
	}
}

func (b *Builder) equipLDSUnit(cu *ComputeUnit) {
	ldsDecoder := NewDecodeUnit(cu)
	cu.LDSDecoder = ldsDecoder

	ldsUnit := NewLDSUnit(cu, b.scratchpadPreparer, b.alu)
	cu.LDSUnit = ldsUnit

	for i := 0; i < b.simdCount; i++ {
		ldsDecoder.AddExecutionUnit(ldsUnit)
	}
}

func (b *Builder) equipVectorMemoryUnit(cu *ComputeUnit) {
	for i := 0; i < b.simdCount; i++ {
		vectorMemDecoder := NewDecodeUnit(cu)
		coalescer := &defaultCoalescer{
			log2CacheLineSize: b.log2CachelineSize,
		}
		vectorMemoryUnit := NewVectorMemoryUnit(cu, b.scratchpadPreparer, coalescer)
		vectorMemoryUnit.postInstructionPipelineBuffer = util.NewBuffer(8)
		vectorMemoryUnit.instructionPipeline = pipelining.NewPipeline(
			cu.Name()+".VectorMemoryUnit.InstPipeline",
			6, 1,
			vectorMemoryUnit.postInstructionPipelineBuffer)

		vectorMemoryUnit.postTransactionPipelineBuffer = util.NewBuffer(8)
		vectorMemoryUnit.transactionPipeline = pipelining.NewPipeline(
			cu.Name()+".VectorMemoryUnit.TransactionPipeline",
			60, 1,
			vectorMemoryUnit.postTransactionPipelineBuffer)
		vectorMemDecoder.AddExecutionUnit(vectorMemoryUnit)
		cu.VectorMemUnits = append(cu.VectorMemUnits, vectorMemoryUnit)
		cu.VectorMemDecoders = append(cu.VectorMemDecoders, vectorMemDecoder)
	}
}

func (b *Builder) equipRegisterFiles(cu *ComputeUnit) {
	for i := 0; i < b.simdCount; i++ {
		sRegFile := NewSimpleRegisterFile(uint64(b.sgprCount[i]*4), 0)
		cu.SRegFiles = append(cu.SRegFiles, sRegFile)
	}

	for i := 0; i < b.simdCount; i++ {
		vRegFile := NewSimpleRegisterFile(uint64(b.vgprCount[i]*4), 1024)
		cu.VRegFiles = append(cu.VRegFiles, vRegFile)
	}
}
