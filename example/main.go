package main

import (
	"bytes"
	"fmt"
	"github.com/pyroscope-io/godeltaprof"
	_ "github.com/pyroscope-io/godeltaprof/http/pprof"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

//go:noinline
func work(n int) {
	// revive:disable:empty-block this is fine because this is a example app, not real production code
	for i := 0; i < n; i++ {
	}
	fmt.Printf("work\n")
	// revive:enable:empty-block
}

var m sync.Mutex

func fastFunction(wg *sync.WaitGroup) {
	m.Lock()
	defer m.Unlock()

	work(200000000)

	wg.Done()
}

func slowFunction(wg *sync.WaitGroup) {
	m.Lock()
	defer m.Unlock()

	work(800000000)
	wg.Done()
}

func main() {
	go func() {
		err := http.ListenAndServe("localhost:6060", http.DefaultServeMux)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		var deltaHeapProfiler = godeltaprof.NewHeapProfiler()
		var deltaBlockProfiler = godeltaprof.NewBlockProfiler()
		var deltaMutexProfiler = godeltaprof.NewMutexProfiler()
		for {
			time.Sleep(10 * time.Second)
			deltaHeapProfiler.Profile(bytes.NewBuffer(nil))
			deltaBlockProfiler.Profile(bytes.NewBuffer(nil))
			deltaMutexProfiler.Profile(bytes.NewBuffer(nil))
		}
	}()
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	for {
		wg := sync.WaitGroup{}
		wg.Add(2)
		go fastFunction(&wg)
		go slowFunction(&wg)
		wg.Wait()
	}
}
