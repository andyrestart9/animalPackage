package main

import "fmt"

func main() {
	c := make(chan int) // 宣告無緩衝 channel

	// 在新的 goroutine 中送出 42
	go func() {
		c <- 42
	}()

	fmt.Println(<-c) // main goroutine 嘗試從 channel 接收資料
}

/*
為什麼這樣不會死鎖？
因為：

go func() { c <- 42 }() 這行會開一個新的 goroutine，它會去發送資料。

同時，main goroutine 在跑 fmt.Println(<-c)，它在等資料。

這兩個 goroutine 可以「同步」完成一次傳遞：一個送、一個收，配對成功
*/
