package main

func main() {
	// 建立一個 send-only channel（只能發送）
	var sendCh chan<- int = make(chan<- int, 1)
	// 正確：使用 send-only channel 發送數值
	sendCh <- 100

	// 這裡我們示範宣告一個 receive-only channel，
	var recvCh <-chan int = make(<-chan int, 1)
	_ = recvCh // 避免「declared and not used」的錯誤

	// 以下為錯誤示範，若取消註解將導致編譯錯誤：
	// 1. 從 send-only channel 嘗試接收數值（錯誤）：
	// invalid operation: cannot receive from send-only channel sendCh (variable of type chan<- int)
	//    x := <-sendCh
	//
	// 2. 向 receive-only channel 嘗試發送數值（錯誤）：
	// invalid operation: cannot send to receive-only channel recvCh (variable of type <-chan int)
	//    recvCh <- 300
}
