package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// 啟動兩條 boring 生產者，並將它們的輸出合併到同一條 channel
	c := fanIn(boring("Joe"), boring("Ann"))

	// 從合併後的 channel 中取出 10 條訊息並印出
	for i := 0; i < 10; i++ {
		fmt.Println(<-c) // 從 c 接收一個字串並印出，阻塞直到有值
	}

	// 讀完 10 筆後結束程式
	fmt.Println("You're both boring: I'm leaving.")
}

// boring 返回一條只可接收 (<-chan) 的 string channel
// 背景 goroutine 產生訊息後隨機暫停部分微秒
func boring(msg string) <-chan string {
	c := make(chan string) // 建立無緩衝 channel
	go func() {
		// 產生 0–9 訊息
		for i := 0; i < 10; i++ {
			// 將格式化字串送到 channel，若沒人接收會阻塞
			c <- fmt.Sprintf("%s %d", msg, i)
			// rand.Intn(1e6) 的 1e6 意味著 1 * 10^6 = 1000000
			// 這裡隨機停 0～999 微秒，以製造兩個生產者交錯效果
			time.Sleep(time.Duration(rand.Intn(1e6)) * time.Microsecond)
		}
		// ※ 注意：這裡沒有 close(c)，goroutine 結束後 channel 仍開啟但不再有新值
	}()
	return c // 回傳只可接收的 channel
}

// fanIn 將兩條只讀 string channel 的輸出匯流到一條新的 channel
func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string) // 建立無緩衝 channel 作為匯流出口

	// 第一條 background goroutine：不斷從 input1 讀取並寫入 c
	go func() {
		for {
			// <-input1：從 input1 接收，阻塞直到有值
			// c <- …    ：將接收到的值發送到 c，阻塞直到有人接收
			c <- <-input1
		}
	}()

	// 從 input2 不斷讀取並寫入 c
	go func() {
		for {
			c <- <-input2
		}
	}()

	return c // 回傳匯流後的 channel
}

/*
c <- <-input2
這句話其實是把從 input2 通道接收到的值，立刻再發送到 c 通道──等同於：
```
// 兩步驟寫法
v := <-input2  // 先從 input2 接收一個值
c <- v         // 再把這個值發送到 c
```
而 c <- <-input2 則是把這兩步合成一行：
<-input2 —— 從 input2 通道接收（receive）一個值。
c <- (那個值) —— 把剛接收到的值，發送（send）到 c 通道。
注意這裡的兩個 <- 並不是位元運算，而是分別表示「從通道接收」和「往通道發送」。
*/
