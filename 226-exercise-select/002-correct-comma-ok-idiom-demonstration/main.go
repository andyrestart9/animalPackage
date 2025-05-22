package main

import "fmt"

func main() {
	c := make(chan int)

	go func() {
		c <- 42
		close(c)
	}()

	v, ok := <-c
	fmt.Println(v, ok) // 42 true
	// 往已关闭的 channel 接收（<-c）则永远不会 panic，只会无阻塞地取出零值并返回 ok==false。
	v, ok = <-c
	fmt.Println(v, ok) // 0 false
	v, ok = <-c
	fmt.Println(v, ok) // 0 false
}
