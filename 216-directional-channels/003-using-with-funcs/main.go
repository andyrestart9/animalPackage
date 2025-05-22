package main

import "fmt"

func main() {
	c := make(chan int)

	// send
	go foo(c)

	// receive
	// 直接呼叫 bar(c) 時，這個函數在主 goroutine 中同步執行，因此會在 <-c 處阻塞等待值，直到 foo(c) 傳送數值過來。也就是說，主 goroutine 因為執行 bar(c) 而被迫停下，等待接收值，這樣就保證了通道中的數值能正確地被接收到和處理。
	bar(c)

	// receive
	// 如果你使用 go bar(c)，那麼 bar(c) 就會在一個新的 goroutine 中運行，而主 goroutine 在啟動這個 goroutine 後不會停下來等待它完成。若主 goroutine 後續沒有其他阻塞的操作，它可能會在 bar(c) 還未接收到數值前結束程式，導致程式提前退出，看不到輸出結果或導致其他同步問題。
	// go bar(c)

	fmt.Println("about to exit")
}

func foo(c chan<- int) {
	c <- 42
}

func bar(c <-chan int) {
	fmt.Println(<-c)
}
