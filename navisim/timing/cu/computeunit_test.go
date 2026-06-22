package cu

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/navisim/kernels"
	"gitlab.com/akita/navisim/protocol"
	"gitlab.com/akita/navisim/rdnainsts"
	"gitlab.com/akita/navisim/timing/wavefront"
)

type mockScheduler struct {
}

func (m *mockScheduler) Run(now akita.VTimeInSec) bool {
	return true
}

func (m *mockScheduler) Pause() {
}

func (m *mockScheduler) Resume() {
}

func (m *mockScheduler) Flush() {
}

type mockDecoder struct {
	Inst *rdnainsts.Inst
}

func (d *mockDecoder) Decode(buf []byte) (*rdnainsts.Inst, error) {
	return d.Inst, nil
}

func exampleGrid() *kernels.Grid {
	grid := kernels.NewGrid()

	grid.CodeObject = rdnainsts.NewHsaCo()
	grid.CodeObject.HsaCoHeader = new(rdnainsts.HsaCoHeader)

	packet := new(kernels.HsaKernelDispatchPacket)
	grid.Packet = packet

	wg := kernels.NewWorkGroup()
	wg.Packet = packet
	wg.CodeObject = grid.CodeObject
	grid.WorkGroups = append(grid.WorkGroups, wg)

	wf := kernels.NewWavefront(32)
	wf.WG = wg
	wg.Wavefronts = append(wg.Wavefronts, wf)

	return grid
}

var _ = Describe("ComputeUnit", func() {
	var (
		mockCtrl          *gomock.Controller
		cu                *ComputeUnit
		engine            *MockEngine
		wfDispatcher      *MockWfDispatcher
		decoder           *mockDecoder
		toInstMem         *MockPort
		toScalarMem       *MockPort
		toVectorMem       *MockPort
		toACE             *MockPort
		toCP              *MockPort
		ldsDecoder        *MockSubComponent
		ldsUnit           *MockSubComponent
		instMem           *MockPort
		branchUnits       []*MockSubComponent
		vectorMemDecoders []*MockSubComponent
		vectorMemUnits    []*MockSubComponent
		scalarDecoders    []*MockSubComponent
		vectorDecoders    []*MockSubComponent
		scalarUnits       []*MockSubComponent
		simdUnits         []*MockSubComponent
		schedulers        []*mockScheduler
		grid              *kernels.Grid
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		engine = NewMockEngine(mockCtrl)
		wfDispatcher = NewMockWfDispatcher(mockCtrl)
		decoder = new(mockDecoder)

		ldsDecoder = NewMockSubComponent(mockCtrl)
		ldsUnit = NewMockSubComponent(mockCtrl)

		cu = NewComputeUnit("cu", engine)
		cu.WfDispatcher = wfDispatcher
		cu.Decoder = decoder
		cu.Freq = 1
		for i := 0; i < 4; i++ {
			cu.SRegFiles = append(cu.SRegFiles, NewSimpleRegisterFile(1024, 0))
			cu.VRegFiles = append(cu.VRegFiles, NewSimpleRegisterFile(4096, 64))
			schedulers = append(schedulers, new(mockScheduler))
			branchUnits = append(branchUnits, NewMockSubComponent(mockCtrl))
			vectorMemDecoders = append(vectorMemDecoders, NewMockSubComponent(mockCtrl))
			vectorMemUnits = append(vectorMemUnits, NewMockSubComponent(mockCtrl))
			scalarDecoders = append(scalarDecoders, NewMockSubComponent(mockCtrl))
			vectorDecoders = append(vectorDecoders, NewMockSubComponent(mockCtrl))
			scalarUnits = append(scalarUnits, NewMockSubComponent(mockCtrl))
			simdUnits = append(simdUnits, NewMockSubComponent(mockCtrl))
		}

		cu.LDSDecoder = ldsDecoder
		cu.LDSUnit = ldsUnit
		for i := 0; i < 4; i++ {
			cu.BranchUnits = append(cu.BranchUnits, branchUnits[i])
			cu.VectorMemDecoders = append(cu.VectorMemDecoders, vectorMemDecoders[i])
			cu.VectorMemUnits = append(cu.VectorMemUnits, vectorMemUnits[i])
			cu.ScalarDecoders = append(cu.ScalarDecoders, scalarDecoders[i])
			cu.VectorDecoders = append(cu.VectorDecoders, vectorDecoders[i])
			cu.ScalarUnits = append(cu.ScalarUnits, scalarUnits[i])
			cu.SIMDUnits = append(cu.SIMDUnits, simdUnits[i])
			cu.Schedulers = append(cu.Schedulers, schedulers[i])
		}

		for i := 0; i < 4; i++ {
			cu.WfPools = append(cu.WfPools, NewWavefrontPool(10))
		}

		toInstMem = NewMockPort(mockCtrl)
		toACE = NewMockPort(mockCtrl)
		toScalarMem = NewMockPort(mockCtrl)
		toVectorMem = NewMockPort(mockCtrl)
		cu.ToInstMem = toInstMem
		cu.ToACE = toACE
		cu.ToScalarMem = toScalarMem
		cu.ToVectorMem = toVectorMem

		instMem = NewMockPort(mockCtrl)
		cu.InstMem = instMem

		toCP = NewMockPort(mockCtrl)

		cu.ToCP = toCP

		grid = exampleGrid()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when processing MapWGReq", func() {
		var (
			req *protocol.MapWGReq
		)

		BeforeEach(func() {
			wg := grid.WorkGroups[0]
			wg.Wavefronts = make([]*kernels.Wavefront, 2)
			wg.Wavefronts[0] = kernels.NewWavefront(32)
			wg.Wavefronts[1] = kernels.NewWavefront(32)
			location1 := protocol.WfDispatchLocation{
				Wavefront:  wg.Wavefronts[0],
				SIMDID:     1,
				VGPROffset: 100,
				SGPROffset: 10,
				LDSOffset:  100,
			}
			location2 := protocol.WfDispatchLocation{
				Wavefront:  wg.Wavefronts[1],
				SIMDID:     2,
				VGPROffset: 200,
				SGPROffset: 200,
				LDSOffset:  200,
			}

			builder := protocol.MapWGReqBuilder{}.
				WithSendTime(6).
				WithSrc(nil).
				WithDst(cu.ToACE).
				WithWG(wg).
				AddWf(location1).
				AddWf(location2)
			req = builder.Build()
			req.RecvTime = 10

			toACE.EXPECT().Retrieve(gomock.Any()).Return(req)
		})

		It("should dispatch wavefront", func() {
			wfDispatcher.EXPECT().
				DispatchWf(akita.VTimeInSec(11),
					gomock.Any(), req.Wavefronts[0])
			wfDispatcher.EXPECT().
				DispatchWf(akita.VTimeInSec(11),
					gomock.Any(), req.Wavefronts[1])
			engine.EXPECT().Schedule(gomock.Any())

			cu.processInputFromACE(11)

			Expect(cu.WfPools[1].wfs).To(HaveLen(1))
			Expect(cu.WfPools[2].wfs).To(HaveLen(1))
		})
	})

	Context("when handling DataReady from ToInstMem Port", func() {
		var (
			wf        *wavefront.Wavefront
			dataReady *mem.DataReadyRsp
		)
		BeforeEach(func() {
			wf = new(wavefront.Wavefront)
			inst := wavefront.NewInst(nil)
			wf.SetDynamicInst(inst)
			wf.PC = 0x1000

			req := mem.ReadReqBuilder{}.
				WithSendTime(8).
				WithSrc(cu.ToInstMem).
				WithDst(instMem).
				WithAddress(0x100).
				WithByteSize(64).
				Build()

			dataReady = mem.DataReadyRspBuilder{}.
				WithSendTime(10).
				WithSrc(instMem).
				WithDst(cu.ToInstMem).
				WithRspTo(req.ID).
				WithData([]byte{
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
				}).
				Build()

			dataReady.RecvTime = 10
			dataReady.EventTime = 10
			toInstMem.EXPECT().Retrieve(gomock.Any()).Return(dataReady)

			info := new(InstFetchReqInfo)
			info.Wavefront = wf
			info.Req = req
			cu.InFlightInstFetch = append(cu.InFlightInstFetch, info)
		})

		It("should handle fetch return", func() {
			madeProgress := cu.processInputFromInstMem(10)

			//Expect(wf.State).To(Equal(WfFetched))
			Expect(wf.LastFetchTime).To(BeNumerically("~", 10))
			Expect(wf.PC).To(Equal(uint64(0x1000)))
			Expect(cu.InFlightInstFetch).To(HaveLen(0))
			Expect(wf.InstBuffer).To(HaveLen(64))
			Expect(madeProgress).To(BeTrue())
		})
	})

	Context("should handle DataReady from ToScalarMem port", func() {
		var (
			wf *wavefront.Wavefront
		)

		BeforeEach(func() {
			rawWf := grid.WorkGroups[0].Wavefronts[0]
			wf = wavefront.NewWavefront(rawWf)
			wf.SRegOffset = 0
			wf.OutstandingScalarMemAccess = 1
		})

		It("should handle scalar data load return", func() {
			read := mem.ReadReqBuilder{}.
				WithSendTime(8).
				WithSrc(cu.ToScalarMem).
				WithAddress(0x100).
				WithByteSize(64).
				Build()

			info := new(ScalarMemAccessInfo)
			info.Inst = wavefront.NewInst(rdnainsts.NewInst())
			info.Wavefront = wf
			info.DstSGPR = rdnainsts.SReg(0)
			info.Req = read
			cu.InFlightScalarMemAccess = append(cu.InFlightScalarMemAccess, info)

			rsp := mem.DataReadyRspBuilder{}.
				WithSendTime(10).
				WithRspTo(read.ID).
				WithData(rdnainsts.Uint32ToBytes(32)).
				Build()
			rsp.RecvTime = 10
			toScalarMem.EXPECT().Retrieve(gomock.Any()).Return(rsp)

			cu.processInputFromScalarMem(10)

			access := RegisterAccess{
				Reg:        rdnainsts.SReg(0),
				RegCount:   1,
				WaveOffset: 0,
				Data:       make([]byte, 4),
			}
			cu.SRegFiles[wf.SIMDID].Read(access)
			Expect(rdnainsts.BytesToUint32(access.Data)).To(Equal(uint32(32)))
			Expect(wf.OutstandingScalarMemAccess).To(Equal(0))
			Expect(cu.InFlightScalarMemAccess).To(HaveLen(0))
		})
	})

	Context("should handle DataReady from ToVectorMem", func() {
		var (
			rawWf *kernels.Wavefront
			wf    *wavefront.Wavefront
			inst  *wavefront.Inst
			read  *mem.ReadReq
			info  VectorMemAccessInfo
		)

		BeforeEach(func() {
			rawWf = grid.WorkGroups[0].Wavefronts[0]
			inst = wavefront.NewInst(rdnainsts.NewInst())
			inst.FormatType = rdnainsts.FLAT
			wf = wavefront.NewWavefront(rawWf)
			wf.SIMDID = 0
			wf.SetDynamicInst(inst)
			wf.VRegOffset = 0
			wf.OutstandingVectorMemAccess = 1
			wf.OutstandingScalarMemAccess = 1

			read = mem.ReadReqBuilder{}.
				WithSendTime(8).
				WithAddress(0x100).
				WithByteSize(16).
				CanWaitForCoalesce().
				Build()

			info = VectorMemAccessInfo{}
			info.Read = read
			info.Wavefront = wf
			info.Inst = inst
			info.laneInfo = []vectorMemAccessLaneInfo{
				{0, rdnainsts.VReg(0), 1, 0},
				{1, rdnainsts.VReg(0), 1, 4},
				{2, rdnainsts.VReg(0), 1, 8},
				{3, rdnainsts.VReg(0), 1, 12},
			}
			cu.InFlightVectorMemAccess = append(
				cu.InFlightVectorMemAccess, info)

			dataReady := mem.DataReadyRspBuilder{}.
				WithSendTime(10).
				WithRspTo(read.ID).
				WithData(make([]byte, 16)).
				Build()
			for i := 0; i < 4; i++ {
				copy(dataReady.Data[i*4:i*4+4], rdnainsts.Uint32ToBytes(uint32(i)))
			}
			toVectorMem.EXPECT().Retrieve(gomock.Any()).Return(dataReady)
		})

		It("should handle vector data load return, and the return is not the last one for an instruction", func() {
			cu.processInputFromVectorMem(10)

			for i := 0; i < 4; i++ {
				access := RegisterAccess{}
				access.RegCount = 1
				access.WaveOffset = 0
				access.LaneID = i
				access.Reg = rdnainsts.VReg(0)
				access.Data = make([]byte, access.RegCount*4)
				cu.VRegFiles[0].Read(access)
				Expect(rdnainsts.BytesToUint32(access.Data)).To(Equal(uint32(i)))
			}

			Expect(wf.OutstandingVectorMemAccess).To(Equal(1))
			Expect(wf.OutstandingScalarMemAccess).To(Equal(1))
			Expect(cu.InFlightVectorMemAccess).To(HaveLen(0))
		})

		It("should handle vector data load return, and the return is the last one for an instruction", func() {
			read.CanWaitForCoalesce = false

			cu.processInputFromVectorMem(10)

			Expect(wf.OutstandingVectorMemAccess).To(Equal(0))
			Expect(wf.OutstandingScalarMemAccess).To(Equal(0))
			for i := 0; i < 4; i++ {
				access := RegisterAccess{}
				access.RegCount = 1
				access.WaveOffset = 0
				access.LaneID = i
				access.Reg = rdnainsts.VReg(0)
				access.Data = make([]byte, access.RegCount*4)
				cu.VRegFiles[0].Read(access)
				Expect(rdnainsts.BytesToUint32(access.Data)).To(Equal(uint32(i)))
			}
		})
	})

	Context("handle write done respond from ToVectorMem port", func() {
		var (
			rawWf    *kernels.Wavefront
			inst     *wavefront.Inst
			wf       *wavefront.Wavefront
			info     VectorMemAccessInfo
			writeReq *mem.WriteReq
			doneRsp  *mem.WriteDoneRsp
		)

		BeforeEach(func() {
			rawWf = grid.WorkGroups[0].Wavefronts[0]
			inst = wavefront.NewInst(rdnainsts.NewInst())
			inst.FormatType = rdnainsts.FLAT
			wf = wavefront.NewWavefront(rawWf)
			wf.SIMDID = 0
			wf.SetDynamicInst(inst)
			wf.VRegOffset = 0
			wf.OutstandingVectorMemAccess = 1
			wf.OutstandingScalarMemAccess = 1
			wf.OutstandingVectorStoreAccess = 1

			writeReq = mem.WriteReqBuilder{}.
				WithSendTime(8).
				WithAddress(0x100).
				CanWaitForCoalesce().
				Build()

			info = VectorMemAccessInfo{}
			info.Wavefront = wf
			info.Inst = inst
			info.Write = writeReq
			cu.InFlightVectorMemAccess = append(cu.InFlightVectorMemAccess, info)

			doneRsp = mem.WriteDoneRspBuilder{}.
				WithSendTime(10).
				WithRspTo(writeReq.ID).
				Build()
			toVectorMem.EXPECT().Retrieve(gomock.Any()).Return(doneRsp)
		})

		It("should handle vector data store return and the return is not the last one from an instruction", func() {
			madeProgress := cu.processInputFromVectorMem(10)

			Expect(cu.InFlightVectorMemAccess).To(HaveLen(0))
			Expect(madeProgress).To(BeTrue())
		})

		It("should handle vector data store return and the return is the last one from an instruction", func() {
			writeReq.CanWaitForCoalesce = false

			cu.processInputFromVectorMem(10)

			Expect(wf.OutstandingVectorStoreAccess).To(Equal(0))
			Expect(wf.OutstandingScalarMemAccess).To(Equal(0))
			Expect(cu.InFlightVectorMemAccess).To(HaveLen(0))
		})
	})

	Context("should handle flush request", func() {
		It("should handle a pipeline flush request from CU", func() {
			req := protocol.CUPipelineFlushReqBuilder{}.
				WithSrc(nil).
				WithDst(cu.ToCP).
				WithSendTime(10).
				Build()

			toCP.EXPECT().Retrieve(akita.VTimeInSec(11)).Return(req)

			cu.processInputFromCP(11)

			Expect(cu.inCPRequestProcessingStage).To(BeIdenticalTo(req))
			Expect(cu.isFlushing).To(BeTrue())
			Expect(cu.currentFlushReq).To(BeIdenticalTo(req))
		})

		It("should flush internal CU buffers", func() {
			info := new(InstFetchReqInfo)
			cu.InFlightInstFetch = append(cu.InFlightInstFetch, info)

			scalarMemInfo := new(ScalarMemAccessInfo)
			cu.InFlightScalarMemAccess = append(cu.InFlightScalarMemAccess, scalarMemInfo)

			vectorMemInfo := VectorMemAccessInfo{}
			cu.InFlightVectorMemAccess = append(cu.InFlightVectorMemAccess, vectorMemInfo)

			cu.flushCUBuffers()

			Expect(cu.InFlightInstFetch).To(BeNil())
			Expect(cu.InFlightVectorMemAccess).To(BeNil())
			Expect(cu.InFlightScalarMemAccess).To(BeNil())
		})

		It("should handle a restart request", func() {
			cu.isPaused = true

			req := protocol.CUPipelineRestartReqBuilder{}.
				WithSendTime(10).
				WithSrc(nil).
				WithDst(cu.ToCP).
				Build()

			toCP.EXPECT().Retrieve(gomock.Any()).Return(req)
			toCP.EXPECT().Send(gomock.Any())

			cu.processInputFromCP(11)
			Expect(cu.isPaused).To(BeTrue())
			Expect(cu.isSendingOutShadowBufferReqs).To(BeTrue())
		})

		It("should flush the full CU", func() {
			req := protocol.CUPipelineFlushReqBuilder{}.
				WithSrc(nil).
				WithDst(cu.ToCP).
				WithSendTime(10).
				Build()

			cu.currentFlushReq = req

			info := new(InstFetchReqInfo)
			cu.InFlightInstFetch = append(cu.InFlightInstFetch, info)

			scalarMemInfo := new(ScalarMemAccessInfo)
			cu.InFlightScalarMemAccess = append(cu.InFlightScalarMemAccess, scalarMemInfo)

			vectorMemInfo := VectorMemAccessInfo{}
			cu.InFlightVectorMemAccess = append(cu.InFlightVectorMemAccess, vectorMemInfo)

			for i := 0; i < 4; i++ {
				branchUnits[i].EXPECT().Flush()
				scalarUnits[i].EXPECT().Flush()
				scalarDecoders[i].EXPECT().Flush()
				simdUnits[i].EXPECT().Flush()
				vectorDecoders[i].EXPECT().Flush()
				vectorMemDecoders[i].EXPECT().Flush()
				vectorMemUnits[i].EXPECT().Flush()
			}

			ldsUnit.EXPECT().Flush()
			ldsDecoder.EXPECT().Flush()

			cu.flushPipeline(10)

			Expect(cu.InFlightInstFetch).To(BeNil())
			Expect(cu.InFlightVectorMemAccess).To(BeNil())
			Expect(cu.InFlightScalarMemAccess).To(BeNil())

			Expect(cu.shadowInFlightInstFetch).To(Not(BeNil()))
			Expect(cu.shadowInFlightVectorMemAccess).To(Not(BeNil()))
			Expect(cu.shadowInFlightScalarMemAccess).To(Not(BeNil()))

			Expect(cu.toSendToCP).NotTo(BeNil())
			Expect(cu.isFlushing).To(BeFalse())
			Expect(cu.isPaused).To(BeTrue())
		})

		It("should not restart a CU where there are shadow buffer reqs pending", func() {
			info := new(InstFetchReqInfo)
			req := mem.ReadReqBuilder{}.
				WithSendTime(8).
				WithSrc(cu.ToInstMem).
				WithDst(instMem).
				WithAddress(0x100).
				WithByteSize(64).
				Build()
			info.Req = req

			cu.shadowInFlightInstFetch = append(cu.InFlightInstFetch, info)

			scalarMemInfo := new(ScalarMemAccessInfo)
			scalarMemInfo.Req = req
			cu.shadowInFlightScalarMemAccess = append(cu.InFlightScalarMemAccess, scalarMemInfo)

			vectorMemInfo := VectorMemAccessInfo{}
			vectorMemInfo.Read = req
			cu.shadowInFlightVectorMemAccess = append(cu.InFlightVectorMemAccess, vectorMemInfo)

			toInstMem.EXPECT().Send(gomock.Any())
			toVectorMem.EXPECT().Send(gomock.Any())
			toScalarMem.EXPECT().Send(gomock.Any())

			cu.checkShadowBuffers(11)
		})

		It("should restart a CU where there are  no shadow buffer reqs pending", func() {
			cu.shadowInFlightInstFetch = nil
			cu.shadowInFlightScalarMemAccess = nil
			cu.shadowInFlightVectorMemAccess = nil

			cu.checkShadowBuffers(11)

			Expect(cu.isPaused).To(BeFalse())
		})
	})
})
