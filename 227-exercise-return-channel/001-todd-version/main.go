package main

import (
	"fmt"
	"runtime"
)

/*
Exam Questions: launches 10 goroutines, each goroutine adds 10 numbers to a channel
pull the numbers off the channel and print them
*/
func main() {
	x := 10
	y := 10

	c := gen(x, y)

	for i := 0; i < x*y; i++ {
		fmt.Println(i, <-c)
	}
	fmt.Println("3. ROUTINES: ", runtime.NumGoroutine())
}

// 在執行 c := gen(x, y) 已經拿到 return c ，表示 gen(x, y int) <-chan int 裡面的 for loop 已經跑完了，為什麼可以再用 for i := 0; i < x*y; i++ 接收值？
func gen(x, y int) <-chan int {
	c := make(chan int)
	fmt.Println("1. ROUTINES: ", runtime.NumGoroutine())
	// ① 外層 for：這裡只是快速地啟動 x 條 goroutine
	for i := 0; i < x; i++ {
		go func() {
			// ② 這裡的內層 for，是在各自的 goroutine 內並行執行的
			for j := 0; j < y; j++ {
				c <- j
			}
		}()
		fmt.Println("2. ROUTINES: ", runtime.NumGoroutine())
	}
	// ③ 立即 return，根本沒等那些 goroutine 把值送完
	return c
}
/*
You said:
在執行 c := gen(x, y) 已經拿到 return c ，表示 gen(x, y int) <-chan int 裡面的 for loop 已經跑完了，為什麼可以再用 for i := 0; i < x*y; i++ 接收值？


ChatGPT said:
在 Go 裡面，go 語句「啟動」一條 goroutine 之後不會等它跑完才繼續下一行，整個函式也不會等那些 goroutine 結束才 return。

換到你的程式：
```
func gen(x, y int) <-chan int {
    c := make(chan int)
    // ① 外層 for：這裡只是快速地啟動 x 條 goroutine
    for i := 0; i < x; i++ {
        go func() {
            // ② 這裡的內層 for，是在各自的 goroutine 內並行執行的
            for j := 0; j < y; j++ {
                c <- j
            }
        }()
    }
    // ③ 立即 return，根本沒等那些 goroutine 把值送完
    return c
}
```
外層 for i := 0; i < x; i++
這只是在 gen() 這個函式裏面「啟動」x 條 goroutine，然後馬上就走到 return c，把 channel 回傳給呼叫者。

內層 for j := 0; j < y; j++
這才是真正把 0…y-1 寫到 channel c 裏的地方，但它們都在各自的 goroutine 中「並行」執行。gen() 回傳後，這些 goroutine 才會被 scheduler 安排去跑它們的內層迴圈。

為什麼接收端能收到 x*y 個值？
gen() 回傳後，背景 goroutine 開始跑各自的 c<-j。
主程式在 for i := 0; i < x*y; i++ 裏對 <-c 進行 x*y 次接收。
因為你正好接收 x*y 次，所以每次 background goroutine 的阻塞式 send（無緩衝 channel）都會“凑對”到一次主程式的 receive，最終收滿所有值。

關鍵區別
gen() 的外層 for：同步啟動 goroutine，速度極快，不牽涉 channel send/receive，就結束並 return。
各 goroutine 裏的 send：執行在返回之後，真正跟 main() 的 receive 交互。

如果你把 gen() 寫成：
```
for i := 0; i < x; i++ {
    for j := 0; j < y; j++ {
        c <- j      // 這裡會在 gen() 內阻塞，直到 main 來接
    }
}
return c
```
那 gen() 根本不會 return，因為 send 永遠在等待 receive，main 也無法走到那條 receive，因此 deadlock。

總結
go func(){…}() 只是啟動一條新 goroutine，gen() 立刻就 return，不等那條 goroutine 跑完。
返回的是「空的」channel 參考，之後所有 send 都發生在背景 goroutine，main 用 x*y 次 <-c 正好把它們一一接走。
*/