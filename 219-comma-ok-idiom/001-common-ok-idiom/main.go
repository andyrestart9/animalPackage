package main

import "fmt"

func main() {
	c := make(chan int)
	go func() {
		c <- 42
		
		// close(c) 關閉 channel 後，所有後續的接收 (<-c) 都不會阻塞，而是立即回傳零值與 ok=false。
		// 這就是第二次接收得到 (0, false) 的原因。
		close(c)
	}()

	// v 會是剛剛匿名 goroutine 發送的 42。
	// ok 會是 true，表示這次接收是從一條「開啟中且有值」的 channel 得到的。
	v, ok := <-c

	fmt.Println(v, ok)

	// v 會是整數的「零值」0。
	// ok 會是 false，表示這次接收來自一條「已關閉且空」的 channel。
	v, ok = <-c

	fmt.Println(v, ok)
}
