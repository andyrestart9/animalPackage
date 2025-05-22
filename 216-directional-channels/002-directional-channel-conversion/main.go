package main

import "fmt"

func main() {
	// 建立一個雙向 channel，能同時發送與接收資料，並設定緩衝區大小為 1
	ch := make(chan int, 1)

	// 將雙向 channel 透過顯式轉換成 send-only channel
	// 寫法可以是：(chan<- int)(ch)
	var sendCh chan<- int = (chan<- int)(ch)
	// 使用 send-only channel 發送資料 (合法)
	sendCh <- 100

	// 將雙向 channel 透過顯式轉換成 receive-only channel
	// 寫法可以是：(<-chan int)(ch)
	var recvCh <-chan int = (<-chan int)(ch)
	// 使用 receive-only channel 接收資料 (合法)
	value := <-recvCh
	fmt.Println("接收到的數值:", value)

	// ---------------------------
	// 以下示範嘗試將 directional channels 轉回雙向 channel 的情況
	// 這類轉換是不允許的，取消註解後會出現編譯錯誤

	// 嘗試將 send-only channel 轉換回雙向 channel (非法)
	//  (chan int)(sendCh) // 編譯錯誤：cannot convert sendCh (type chan<- int) to type chan int

	// 嘗試將 receive-only channel 轉換回雙向 channel (非法)
	// (chan int)(recvCh) // 編譯錯誤：cannot convert recvCh (type <-chan int) to type chan int

	// 結論：從雙向 channel 轉換為 directional channels 可以用顯式轉換
	//         但單向 channel（send-only 或 receive-only）不能被轉換回雙向 channel，
	//         這是 Go 語言型別安全性的一部分。
}
