package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	// 建立兩條無緩衝 channel：
	// c1 作為輸入工作項管道，c2 作為輸出結果管道
	c1 := make(chan int)
	c2 := make(chan int)

	// 啟動 populate goroutine，負責將 0–9 寫入 c1，並在結束後關閉 c1
	go populate(c1)

	// 啟動 fanOutIn goroutine，讀取 c1 的值並派發給多個 worker 處理，
	// 最終結果寫入 c2，全部處理完畢後關閉 c2
	go fanOutIn(c1, c2)

	// 主程式從 c2 依序讀取處理結果，直到 c2 被關閉
	for v := range c2 {
		fmt.Println(v) // 印出每個處理後的值
	}

	// c2 關閉後跳出迴圈，印出結束訊息
	fmt.Println("about to exit")
}

// populate 將 0–9 寫入輸入 channel c，寫完後關閉 channel
func populate(c chan int) {
	for i := 0; i < 10; i++ {
		c <- i // 發送工作項 i 到 c，若無接收者則阻塞
	}
	close(c) // 關閉 c，通知接收者不再有新資料
}

// fanOutIn 從 c1 讀取每個工作項，對每個項目啟動獨立的 goroutine 處理，
// 處理結果寫入 c2，所有工作完成後關閉 c2
func fanOutIn(c1, c2 chan int) {
	var wg sync.WaitGroup

	// 從 c1 循環接收工作項，直到 c1 被關閉
	for v := range c1 {
		wg.Add(1) // 每啟動一個工作就增加 WaitGroup 計數，新增一筆待處理工作

		// 為每個值並行啟動一個 worker goroutine 做耗時運算並將結果寫入 c2
		go func(v2 int) {
			// 將計算結果寫入 c2，若 main 還沒讀取則會阻塞
			c2 <- timeConsumingWork(v2)
			wg.Done() // 完成一筆工作後呼叫 Done ，標記此筆工作完成
		}(v) // 傳入當前的 v，避免閉包共用同一個 v
	}

	wg.Wait() // 等待所有 worker/goroutine 完成
	close(c2) // 關閉 c2，通知主程式沒有更多資料
}

// timeConsumingWork 模擬一段耗時工作：
// 隨機等待 0–999 毫秒，然後回傳輸入值
func timeConsumingWork(n int) int {
	// 隨機睡眠，模擬計算或 I/O 延遲
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	return n // 回傳原值或計算後結果
}

/*
整個流程大致分成四個階段：

populate → c1
```
for i := 0; i < 10; i++ {
  c1 <- i
}
close(c1)
```
populate 依次把 0–9 發送到 c1，然後關閉 c1。

fanOutIn 讀 c1 並發送到各工作 goroutine
```
for v := range c1 {
  wg.Add(1)
  go func(v2 int) {
    c2 <- timeConsumingWork(v2)
    wg.Done()
  }(v)
}
wg.Wait()
close(c2)
```
for v := range c1 會同步地從 c1 一個一個取出值 v。
每取到一個，就馬上啟動一個新的 goroutine 執行 timeConsumingWork(v)，並把結果寫到 c2。
這裡的 go 讓所有工作可以並行執行，互不等待。
wg 保證等所有 goroutine 都呼叫了 Done()（也就是都把結果送完）後，才關閉 c2。

timeConsumingWork 模擬耗時
```
time.Sleep(...)
return n + rand.Intn(1000)
```
每個工作會先 Sleep 一段隨機時間，再計算並回傳結果。

main 讀 c2 並印出
```
for v := range c2 {
  fmt.Println(v)
}
```
這個迴圈會阻塞在 <-c2，每當有一個工作 goroutine 把結果寫入 c2，主程式就立刻讀出、印出，直到 c2 被關閉為止。

總結
populate：同步送 0–9 → c1
fanOutIn：同步取 c1 → 並行處理（每個值一個 goroutine）→ 寫結果到 c2
main：同步讀 c2 → 印結果 → 等所有結果出完才結束
這樣的架構能夠把耗時工作「fan-out」到多個 goroutine，同時再「fan-in」到一個結果管道，實現高效且簡潔的並行處理。
*/
