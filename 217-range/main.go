package main

import "fmt"

func main() {
	c := make(chan int)

	// send
	go func() {
		for i := 0; i < 5; i++ {
			c <- i
		}
		// 沒有用 close(c) 關閉 channel 的話在執行上會產生問題，因為 for range 循環會一直等待 channel 被關閉。
		// 會造成 fatal error: all goroutines are asleep - deadlock!
		// for range 迴圈的行為：
		// 在主 goroutine 中，你用 for v := range c 來接收通道內的資料。for range 會一直讀取通道，直到該通道被關閉。如果通道沒有被關閉，當資料讀取完畢後，for range 會一直等待新的資料，導致程式無法結束（也就是說會一直阻塞，形成死鎖的情形）。
		close(c)
	}()

	// receive
	for v := range c {
		fmt.Println(v)
	}

	fmt.Println("about to exit")
}
/*
在這個程式中，值是「一個接一個」地讀取，而不是等到所有五個值都發送完才一起讀取。

由於你使用的是無緩衝 channel（make(chan int)），所以發送操作會阻塞直到有接收者讀取這個值。也就是說，當匿名 goroutine 嘗試發送一個值到 channel 時，它必須等待主 goroutine 從 channel 讀取這個值後，才能繼續下一次發送。

主 goroutine 使用 for v := range c 循環，這個循環會在每次從 channel 讀取到一個值後就處理（印出）該值，而不是等到所有值都發送完畢。

因此，程式的運作流程大致如下：

匿名 goroutine 發送第 0 個值（0），這個操作被阻塞直到主 goroutine 從 channel 讀取值。

主 goroutine 讀取值並印出（0）。

匿名 goroutine 發送第 1 個值（1），同樣等待主 goroutine 讀取。

主 goroutine 讀取值並印出（1）。

如此反覆直到發送完 5 個值。

最後，發送完所有值後關閉 channel（如果有呼叫 close(c)），使得 for range 循環退出。

所以，程式不會等到發送完所有 5 個值才開始讀取，而是一邊發送、一邊讀取。
*/