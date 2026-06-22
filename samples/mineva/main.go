package main

import (
	"flag"

	"gitlab.com/akita/navisim/benchmarks/dnn/mineva"
	"gitlab.com/akita/navisim/samples/runner"
)

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := mineva.NewBenchmark(runner.GPUDriver)

	runner.AddBenchmark(benchmark)

	runner.Run()
}
