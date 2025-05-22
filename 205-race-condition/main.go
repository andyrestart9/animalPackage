package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	fmt.Println("Cpus:", runtime.NumCPU())
	fmt.Println("Goroutines:", runtime.NumGoroutine())

	counter := 0

	const gs = 100
	var wg sync.WaitGroup
	wg.Add(gs)

	for i := 0; i < gs; i++ {
		go func() {
			v := counter
			/*
			runtime.Gosched() 是 Go 語言中 runtime 套件提供的一個函數，用來讓當前 goroutine 主動讓出 CPU 控制權，從而讓其他 goroutine 有機會運行。

			詳細說明
			讓出控制權：
			當調用 runtime.Gosched() 時，當前 goroutine 會暫時放棄執行，但它不會被阻塞或終止，只是暫時讓出 CPU，等待調度器再次安排它運行。
			用於多 goroutine 調度：
			在多 goroutine 的環境中，如果某個 goroutine 長時間佔用 CPU，其他 goroutine 可能無法及時執行。調用 runtime.Gosched() 可以使得調度器有機會將 CPU 分配給其他 goroutine，從而提高程序的並發性和響應性。
			不會使 goroutine 停止：
			與 time.Sleep() 不同，runtime.Gosched() 只是短暫地讓出 CPU，而不會使當前 goroutine 暫停一段固定時間。當其他 goroutine 都執行完畢後，當前 goroutine 可能會很快地再次獲得執行權。
			*/
			// runtime.Gosched() 可以讓我們更方便觀察 race condition 因為當其他 goroutine 要用 CPU 時，當前 goroutine 會讓出 CPU 控制權
			runtime.Gosched()
			v++
			counter = v
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
