package main

import "fmt"

/*
andylin@Andys-MacBook-Pro 002-wrong-demonstration-2 % go run main.go
close(c) 尚未被執行, channel c 未被關閉, 但是 <-c 會持續收到 0 值, v = 0, ok = true
close(c) 尚未被執行, channel c 未被關閉, 但是 <-c 會持續收到 0 值, v = 1, ok = true
close(c) 尚未被執行, channel c 未被關閉, 但是 <-c 會持續收到 0 值, v = 2, ok = true
close(c) 尚未被執行, channel c 未被關閉, 但是 <-c 會持續收到 0 值, v = 3, ok = true
close(c) 尚未被執行, channel c 未被關閉, 但是 <-c 會持續收到 0 值, v = 4, ok = true
close(c) 尚未被執行, channel c 未被關閉, 但是 <-c 會持續收到 0 值, v = 5, ok = true
close(c) 尚未被執行, channel c 未被關閉, 但是 <-c 會持續收到 0 值, v = 6, ok = true
close(c) 尚未被執行, channel c 未被關閉, 但是 <-c 會持續收到 0 值, v = 7, ok = true
close(c) 尚未被執行, channel c 未被關閉, 但是 <-c 會持續收到 0 值, v = 8, ok = true
close(c) 尚未被執行, channel c 未被關閉, 但是 <-c 會持續收到 0 值, v = 9, ok = true
close(c) 已執行, channel c 已被關閉, 但是 <-c 會持續收到 0 值, v = 0, ok = false
close(c) 已執行, channel c 已被關閉, 但是 <-c 會持續收到 0 值, v = 0, ok = false
close(c) 已執行, channel c 已被關閉, 但是 <-c 會持續收到 0 值, v = 0, ok = false
about to exit

為什麼會是這樣的執行結果？

当你在 select 里这样写：
```
select {
case v := <-c:
    fmt.Println(v)
case <-q:
    return
}
```
在 gen goroutine 里发生的顺序是：
送出 0…9
close(c)
尝试 q <- 1，此时因为 receive 还在不停地抢 c 分支，又还没到去抢 q，所以 gen 就在这里 阻塞 了。

此时回到 receive：
对于 case v := <-c：
只要 channel c 关闭了，再对它做 <-c 永远都不会阻塞，而是“就绪”地返回零值（这里是 0）。

对于 case <-q：
当 gen 卡在 q<-1 上时，这个分支也是“就绪”的（因为有一个挂起的发送正等着配对）。

当 select 看到 两个都“就绪”时，会在它们之间随机选一个。这样，你就会看到：
有时候选到 c 分支，就打印一个 0；
有时候选到 q 分支，就退出循环。
观察到连续三个 0 很可能是因为连续三次 select 都恰好选中了 c 分支，直到第四次才选 q 分支并 return。
*/

func main() {
	q := make(chan int)
	c := gen(q)

	receive(c, q)

	fmt.Println("about to exit")
}

func receive(c, q <-chan int) {
	for {
		select {
		case v, ok := <-c:
			if ok {
				fmt.Printf("close(c) 尚未被執行, channel c 未被關閉, 但是 <-c 會持續收到 0 值, v = %v, ok = %v\n", v, ok)
			}
			if !ok {
				fmt.Printf("close(c) 已執行, channel c 已被關閉, 但是 <-c 會持續收到 0 值, v = %v, ok = %v\n", v, ok)
			}
		case <-q:
			return
		}

	}
}

func gen(q chan<- int) <-chan int {
	c := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			c <- i
		}
		close(c)
		q <- 1
	}()

	return c
}
