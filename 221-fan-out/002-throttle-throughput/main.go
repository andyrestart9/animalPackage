package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	// 建立兩條無緩衝 channel：
	// c1 用於輸入工作項，c2 用於輸出處理結果
	c1 := make(chan int)
	c2 := make(chan int)

	// 啟動 populate goroutine，生成 0–9 並寫入 c1，完成後關閉 c1
	go populate(c1)

	// 啟動 fanOutIn goroutine，以多個 worker 並行處理 c1 的資料，並將結果寫入 c2
	go fanOutIn(c1, c2)

	// 主程式同步讀取 c2，直到 c2 被關閉
	for v := range c2 {
		fmt.Println(v) // 印出每個處理後的結果
	}

	// c2 關閉後跳出迴圈，印出結束訊息
	fmt.Println("about to exit")
}

// populate 從 0 到 9 寫入 channel c，然後關閉 c
func populate(c chan int) {
	for i := 0; i < 10; i++ {
		c <- i    // 發送工作項 i，若無接收者則會阻塞
	}
	close(c)   // 關閉 c，通知所有接收者不再有新值
}

// fanOutIn 使用固定數量的 worker goroutine 並行消費 c1
// 並將處理結果寫入 c2，所有工作完成後再關閉 c2
func fanOutIn(c1, c2 chan int) {
	var wg sync.WaitGroup
	const goroutines = 3    // 定義並行 worker 的數量
	wg.Add(goroutines)      // 設定 WaitGroup 需要等待的 goroutine 數量

	// 啟動多個 worker goroutine
	for i := 0; i < goroutines; i++ {
		go func() {
			// 每個 worker 從 c1 連續讀取直到 c1 被關閉
			for v := range c1 {
				// 呼叫耗時函式處理 v，並將結果發送到 c2
				c2 <- timeConsumingWork(v)
			}
			wg.Done() // 表示此 worker 已完成所有任務
		}()
	}

	wg.Wait()   // 等待所有 worker 完成
	close(c2)   // 關閉 c2，通知主程式結束接收
}

// timeConsumingWork 模擬耗時工作：
// 隨機延遲 0–1999999 微秒後返回輸入值
func timeConsumingWork(n int) int {
	// rand.Intn(2000000) 產生 0–1999999 的隨機整數
	// 乘以 time.Microsecond 將其轉換為微秒
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(2000000)))
	return n
}
