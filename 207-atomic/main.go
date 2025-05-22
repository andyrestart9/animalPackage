package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println("Cpus:", runtime.NumCPU())
	fmt.Println("Goroutines:", runtime.NumGoroutine())

	var counter int64

	const gs = 100
	var wg sync.WaitGroup
	wg.Add(gs)


	for i := 0; i < gs; i++ {
		go func() {
			atomic.AddInt64(&counter,1)
			runtime.Gosched()
			fmt.Println("Count\t",atomic.LoadInt64(&counter))
			wg.Done()
		}()
		fmt.Println("Goroutines:", runtime.NumGoroutine())
	}
	wg.Wait()
	fmt.Println("Goroutines:", runtime.NumGoroutine())
	fmt.Println("count:", counter)
	// 用 go run -race main.go 指令，這個指令會告訴你妳的程式有沒有 race condition，最後會出現類似 Found 2 data race(s) 的資訊
	/*
		怎麼找到 -race 這個 flag?
		`go help`
		Use "go help <command>" for more information about a command.
		`go help run`
		For more about build flags, see 'go help build'.
		`go help build`
		        -race
                enable data race detection.
                Supported only on linux/amd64, freebsd/amd64, darwin/amd64, darwin/arm64, windows/amd64,
                linux/ppc64le and linux/arm64 (only for 48-bit VMA).
	*/
}
