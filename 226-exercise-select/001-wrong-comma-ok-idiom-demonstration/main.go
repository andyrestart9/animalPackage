package main

import "fmt"

/*
% go run main.go
panic: send on closed channel

You said:
為什麼是 panic: send on closed channel 不是 panic: receive on closed channel？我發送之後才關閉 channel


ChatGPT said:
Go 运行到那行 close(c) 的时机，几乎总是在新启动的 goroutine 真的有机会执行 c <- 42 之前：

你在主 goroutine 里立刻写了
```
go func() { c <- 42 }()
close(c)
```
go func… 只是告诉运行时“请在后台执行这个函数”，它并不会等这个匿名函数跑到 c<-42 才继续往下。

紧接着主 goroutine 就执行了 close(c)，把 channel 关掉了。

背景 goroutine 终于开始执行 c <- 42 的时候，channel 已经关了，往已关闭的 channel 发送，这就会在 Go 里触发
```
panic: send on closed channel
```

所以你看到的是 “send on closed channel”——因为错误发生在那一刻后台 goroutine 在做 send，而不是在接收那端。
*/

func main() {
	c := make(chan int)

	go func() {
		c <- 42
	}()

	close(c)
	v, ok := <-c
	fmt.Println(v, ok)
}
