package main

import (
	"fmt"
	"sync"
)

func main() {
	// 建立三條無緩衝 channel：
	// even  用來接收偶數
	// odd   用來接收奇數
	// fanin 用來匯流 even/odd 的結果
	even := make(chan int)
	odd := make(chan int)
	fanin := make(chan int)

	// 啟動生產者 goroutine：
	// 將 0–9 拆分為偶數/奇數，分別送到 even 和 odd
	go send(even, odd)

	// 啟動匯流者 goroutine：
	// 從 even/odd 讀取後寫到 fanin，最後關閉 fanin
	go receive(even, odd, fanin)

	// 在 main goroutine 中同步讀取 fanin
	// range 會持續接收直到 fanin 被 close()
	for v := range fanin {
		fmt.Println(v) // 印出從 fanin 匯流過來的值
	}

	// fanin 關閉後跳出 range，印出結束訊息
	fmt.Println("abort to exit")
}

// send 只允許發送到 even 或 odd 通道
func send(even, odd chan<- int) {
	// 產生 0–9，按奇偶寫入對應通道
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			even <- i // 偶數送到 even，若沒人接收則會阻塞
		} else {
			odd <- i  // 奇數送到 odd，若沒人接收則會阻塞
		}
	}
	// 全部送完後，關閉 even/odd 通道，通知接收者 no more data
	close(even)
	close(odd)
}

// receive 只允許從 even/odd 接收，並發送到 fanin
func receive(even, odd <-chan int, fanin chan<- int) {
	var wg sync.WaitGroup
	wg.Add(2) // 兩條來源通道分別由兩個 goroutine 處理

	// 匯流偶數部分
	go func() {
		// 讀取 even channel，直到它被 close()
		for v := range even {
			fanin <- v // 發送到 fanin，若 main 沒接收則會阻塞
		}
		wg.Done() // 標記此 goroutine 處理完成
	}()

	// 匯流奇數部分
	go func() {
		for v := range odd {
			fanin <- v
		}
		wg.Done()
	}()

	// 等待上述兩個 goroutine 都 Done()
	wg.Wait()
	// 全部值都匯流完畢後，關閉 fanin，結束主程式的 range 迴圈
	close(fanin)
}
