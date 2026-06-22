package main

import (
	"flag"

	"gitlab.com/akita/navisim/benchmarks/polybench/atax"
	"gitlab.com/akita/navisim/samples/runner"
)

var xFlag = flag.Int("x", 4096, "The width of the matrix.")
var yFlag = flag.Int("y", 4096, "The height of the matrix.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := atax.NewBenchmark(runner.GPUDriver)
	benchmark.NX = *xFlag
	benchmark.NY = *yFlag

	runner.AddBenchmark(benchmark)

	runner.Run()
}
