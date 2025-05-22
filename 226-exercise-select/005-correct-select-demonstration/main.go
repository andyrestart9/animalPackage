package main

import "fmt"

/*
正確的做法應該在 func gen(q chan<- int) <-chan int 裡面先執行 q <- 1 ，再執行 close(c)
因為：
把那個 q<-1 的動作放在了關閉 c 之前：
```
go func() {
    for i := 0; i < 10; i++ {
        c <- i
    }
    q <- 1     // ← 先送結束訊號（會阻塞到 receive 抢到它）
    close(c)   // ← 再關閉 c
}()
```
這麼一來，當 gen 發完 0…9 之後，它馬上卡在 q<-1，等 receive 裡的 case <-q 被選中、把 1 接走以後才會繼續執行 close(c)。

而在你的 receive 裡：
```
for {
  select {
  case v := <-c:
    fmt.Println(v)
  case <-q:
    return
  }
}
```
在收到 0…9 之後，c 此刻暫時還沒被關閉（因為要等你從 q 接走那個 1 以後才關它），不過又沒有任何新的 c<- 在跑，於是 <-c 會阻塞。

相反地，q<-1 已經在那頭掛著一個等待配對的發送，所以 case <-q: 就是就緒的分支──select 看到這裡能馬上跑，就直接跑了 return，結束 receive，跳出 for，印出 about to exit。

總結：
把 q<-1 放在 close(c) 之前，讓 receive 先搶到 q 的分支就 return，還來不及看到任何 <-c 在 c 被關後回傳零值（那才會打印 0）。
也就是說，你這次的「先 q<-1、後 close(c)」改寫，確保了 receive 在 c 真正關閉、開始返零值之前就已經跳出去了。
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
		q <- 1
		close(c)
	}()

	return c
}
