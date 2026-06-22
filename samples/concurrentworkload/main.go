package main

import (
	"flag"

	_ "net/http/pprof"

	"gitlab.com/akita/navisim/benchmarks/amdappsdk/bitonicsort"
	"gitlab.com/akita/navisim/benchmarks/heteromark/fir"
	"gitlab.com/akita/navisim/samples/runner"
)

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	firBenchmark := fir.NewBenchmark(runner.GPUDriver)
	firBenchmark.Length = 10240
	firBenchmark.SelectGPU([]int{1, 2})

	bsBenchmark := bitonicsort.NewBenchmark(runner.GPUDriver)
	bsBenchmark.Length = 64
	bsBenchmark.SelectGPU([]int{3})

	runner.AddBenchmarkWithoutSettingGPUsToUse(firBenchmark)
	runner.AddBenchmarkWithoutSettingGPUsToUse(bsBenchmark)

	runner.Run()
}
