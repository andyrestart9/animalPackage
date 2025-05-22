package main

import "fmt"

/*
執行 go run main.go 會不斷打印 0
因為；
从一个已关闭的 channel 接收有两条语义：
接收操作 不会阻塞，而是立刻就绪。
它会返回该类型的 零值，并且如果你用“comma-ok”格式还能拿到一个 ok == false

c 一旦被 close(c)，所有后续对 c 的 <-c 都会 立刻返回 0，并且不会阻塞。
在 select 中，任何就绪的分支都有机会被选中；因为 case v := <-c 在关闭后永远就绪，它就会一直被跑。

close(c) 之後 c 還是會一直傳送值嗎？
不會，再往已關閉的 channel 發送值會直接 panic，但從已關閉的 channel 接收則永遠不會阻塞，會一直拿到該類型的「零值」。

具體來說：
發送端：
```
close(c)
c <- 1      // panic: send on closed channel
```
你一旦呼叫 close(c)，後面任何 c <- … 都會在執行期立刻 panic，因為 channel 已經關閉，不能再放入

為什麼 <-c 在 c 關閉後永遠「就緒」？
正常情況
只有當 channel 裡有值，<-c 才能立刻拿到值並且不阻塞；否則就要等（阻塞）到有東西送進來。

已關閉的 channel
Go 保證：一旦 close(c)，所有對它的接收 v := <-c
不會阻塞，而是立即回傳一個零值（如 0）；
如果你用雙值寫法 v, ok := <-c，還能拿到 ok==false 表示「已關閉」。
由於「不會阻塞」，在 select 判斷時這個 case <-c: 就一直被視為「就緒」，永遠有資格被選中。
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
		case v := <-c:
			fmt.Println(v)
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
		// q <- 1
	}()

	return c
}
