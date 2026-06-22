# A Tutorial of Go and GCN3Sim Handout

## Environment

Connecting to the environment

```bash
ssh gcn3sim_guest@lowfreq.ece.neu.edu
```

What is in the docker image:

1. The most recent stable release of Go
1. GOPATH is configured at `/go`
1. Vim is installed
1. Vim-go, a vim plugin for Go development, is configured
1. GCN3Sim is installed at `/go/src/gitlab.com/yaotsu/gcn3`
1. Some code that is going to be used in the tutorial can be found in `/go/src/github.com/syifan/gcn3sim_tutorial`

## Demo 1: Go HelloWorld

1. Go to the folder `/go/src/github.com/syifan/gcn3sim_tutorial/`.
1. Create a folder called `helloworld` and `cd` into it.
1. Create a file called `main.go` by typing `vim main.go`. You should see the code there as vim-go already created it for you. The code should look like this.

    ```go
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello World.")
    }
    ```

1. Run `go build`. You should see a `helloworld` executable generated.
1. Run `./helloworld` and check the result.

## Demo 2: Installing GCN3Sim

1. GCN3Sim is installed in the virtual machine and we need to uninstall it first. We can run `rm -rf /go/src/gitlab.com/yaotsu` to remove everything related to GCN3Sim.
1. Run `go get gitlab.com/yaotsu/gcn3` from any directory. You should see the folder `/go/src/gitlab.com/yaotsu/gcn3` is cloned from the git repository.
1. Change directory to `/go/src/gitlab.com/yaotsu/gcn3/`
1. Run `go get -t ./...` and all the dependencies are installed automatically.
1. Change directory to `/go/src/gitlab.com/yaotsu/gcn3/samples/firsim`.
1. Run `go build` and execute `./firsim`. If there is no error, you have run the first GCN3Sim simulation.

## Demo 3: Implementing Simplified GCN3Sim Event System

1. Create a folder `/go/src/github.com/syifan/gcn3sim_tutorial/events/`. In the newly created folder, create two subfolders as `events` and `eventmain`.
1. Add file `events.go` to the `events` folder.

    ```go
    package events

    type TimeInSec float64

    type Event interface {
      Time() TimeInSec
    }

    type EventBase struct {
      eventTime TimeInSec
    }

    func (e *EventBase) Time() TimeInSec {
      return e.eventTime
    }

    type SomeEvent struct {
      *EventBase
      // Some other fields
    }

    func NewSomeEvent(t TimeInSec) *SomeEvent {
      eventBase := new(EventBase)
      eventBase.eventTime = t
      event := SomeEvent{eventBase}
      return &event
    }
    ```

1. Add `engine.go` into the `events` folder.

    ```go
    package events

    import "fmt"

    type Engine struct {
      events []Event
    }

    func (e *Engine) Schedule(event Event) {
      e.events = append(e.events, event)
    }

    func (e *Engine) Run() {
      for len(e.events) > 0 {
        minTime := TimeInSec(0)
        minIndex := 0
        for i := 0; i < len(e.events); i++ {
          if minTime < e.events[i].Time() {
            minTime = e.events[i].Time()
            minIndex = i
          }
        }

        fmt.Printf("Event time %f\n", e.events[minIndex].Time())

        e.events = append(e.events[:minIndex], e.events[minIndex+1:]...)
      }
    }
    ```

1. Add `main.go` into the `eventmain` folder.

    ```go
    package main

    import (
      "math/rand"

      "github.com/syifan/gcn3sim_tutorial_examples/events/events"
    )

    func main() {
      engine := new(events.Engine)
      for i := 0; i < 100; i++ {
        e := events.NewSomeEvent(events.TimeInSec(rand.Float64()))
        engine.Schedule(e)
      }
      engine.Run()
    }
    ```

1. In the `eventmain` folder, run `go build` and `./eventmain`. Check the program output.

## Demo 4: A Ping-pong Simulation: How to Write Component for GCN3Sim

1. Go to the folder `/go/src/github.com/syifan/gcn3sim_tutorial/pingpong`. You should see several files are already provided there.
1. Append the following code to `agent.go` file.

    ```go
    type PingAgent struct {
      *core.ComponentBase
      engine core.Engine
    }

    func NewPingAgent(name string, engine core.Engine) *PingAgent {
      base := core.NewComponentBase(name)
      agent := new(PingAgent)
      agent.ComponentBase = base
      agent.engine = engine

      agent.AddPort("Out")

      return agent
    }

    func (a *PingAgent) Recv(r core.Req) *core.Error {
      switch req := r.(type) {
      case *PingReq:
        pongEvent := NewSendPongEvent(
          r.RecvTime()+core.VTimeInSec(rand.Float64()), a)
        pongEvent.pingTime = r.SendTime()
        pongEvent.dst = r.Src()
        a.engine.Schedule(pongEvent)
      case *PongReq:
        fmt.Printf("send: %f, recv: %f, latency: %f\n",
          req.pingTime, req.RecvTime(), req.RecvTime()-req.pingTime)
      }
      return nil
    }

    func (a *PingAgent) Handle(evt core.Event) error {
      switch event := evt.(type) {
      case *PingEvent:
        pingReq := NewPingReq()
        pingReq.SetSrc(a)
        pingReq.SetDst(event.dst)
        pingReq.SetSendTime(event.Time())
        a.GetConnection("Out").Send(pingReq)

      case *SendPongEvent:
        pongReq := NewPongReq()
        pongReq.pingTime = event.pingTime
        pongReq.SetSrc(a)
        pongReq.SetDst(event.dst)
        pongReq.SetSendTime(event.Time())
        a.GetConnection("Out").Send(pongReq)
      }
      return nil
    }
    ```

1. Run `go build` and `./pingpong` and check the result.

## Demo 5: Prepare `VectorAdd` Benchmark for GCN3Sim

1. Go to the folder `/go/src/github.com/syifan/gcn3sim_tutorial/vectoradd` directory. You should see a `vector_add.cl` file there. It is the Vector Add OpenCL kernel. Feel free to open the file and check the content.
1. Let's compile the kernel with `clang-ocl` by running `clang-ocl vector_add.cl -mcpu=gfx803 -o vector_add.hsaco`. In case you want to check the assembly of the compiled kernel, you can run `clang-ocl vector_add.cl -mcpu=gfx803 -S -o vector_add.isa`.
1. Write the following code in `vectoradd.go`

    ```go
    package main

    import (
      "fmt"

        "gitlab.com/yaotsu/gcn3/driver"
        "gitlab.com/yaotsu/gcn3/kernels"
        "gitlab.com/yaotsu/gcn3/platform"
    )

    type VecAddArgs struct {
        a, b, c                   driver.GPUPtr
        n, padding                int32
        hiddenX, hiddenY, hiddenZ uint64
    }

    func main() {
        _, gpu, driver, mem := platform.BuildEmuPlatform()
        kernel := kernels.LoadProgram("vector_add.hsaco", "vec_add")
        storage := mem.Storage

        dataSize := 1024
        dA := driver.AllocateMemory(storage, uint64(dataSize*4))
        dB := driver.AllocateMemory(storage, uint64(dataSize*4))
        dC := driver.AllocateMemory(storage, uint64(dataSize*4))

        hA := make([]float32, dataSize)
        hB := make([]float32, dataSize)
        hC := make([]float32, dataSize)

        for i := 0; i < dataSize; i++ {
            hA[i] = float32(i)
            hB[i] = float32(i)
        }

        driver.MemoryCopyHostToDevice(dA, hA, storage)
        driver.MemoryCopyHostToDevice(dB, hB, storage)

        args := VecAddArgs{dA, dB, dC, int32(dataSize), 0, 0, 0, 0}
        driver.LaunchKernel(kernel, gpu, storage,
            [3]uint32{uint32(dataSize), 1, 1},
            [3]uint16{64, 1, 1},
            &args)

        driver.MemoryCopyDeviceToHost(hC, dC, storage)
        for i := 0; i < dataSize; i++ {
            fmt.Printf("[%d]: %f\n", i, hC[i])
        }
    }
    ```
1. Run `go build` and `./vectoradd`.

## Demo 6: Use the GCN3Sim Visualization Tool

1. Go to the folder `/go/src/gitlab.com/yaotsu/gcn3/samples/firsim/` and run `go build`.
1. Run command `./firsim -timing -trace-inst -dataSize 65536`. You will see a large number of trace files generated in the current folder. They are the trace files for each CU.
1. Go to folder `/go/src/gitlab.com/yaotsu/gcn3/trace/vis` and run `go build`
1. Run the visualization server by using

    ```bash
      ./vis -http=:6060 ../../samples/firsim/inst_GPU.CU00.trace
    ```

1. If you run on it on your local environment, a browser window will popup. However, since you are running it in our tutorial container environment, you need to manually open your browser  and type in the following URL: `http://lowfreq.ece.neu.edu:[port]`. Replace the port number with the port number you are prompted at the time when you log in the system. If you do not remember the number, run `source ~/.bashrc` to print it again.
1. Feel free to drag in the minimap at the bottom of the interface and hover your mouse pointer on the pipeline diagram.
