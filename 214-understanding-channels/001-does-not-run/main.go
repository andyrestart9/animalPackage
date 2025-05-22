package main

import "fmt"

func main() {
	c := make(chan int) // 建立一個無緩衝的 channel

	c <- 42 // 嘗試送值進 channel（會阻塞直到有 goroutine 讀取）

	fmt.Println(<-c) // 嘗試從 channel 讀值（永遠不會執行到這行）
}

/*
在 Go 中，無緩衝（同步）channel 的特性是：在發送操作（c <- 42）進行時，必須要有另一個 goroutine 同時在等待接收，才能使發送成功。如果沒有同時進行的接收操作，發送就會阻塞。而在這段程式碼中，發送和接收都在同一個 goroutine（main 函數）中依序進行。

當程式執行到 c <- 42 時，因為 channel 沒有緩衝區，所以發送操作會等待接收者出現。然而，由於 main goroutine 正在等待發送操作完成，它無法繼續執行到 fmt.Println(<-c) 進行接收。這就造成了一種情況：發送操作無法完成，而接收操作又無法啟動，最終導致死鎖（deadlock）。
*/
