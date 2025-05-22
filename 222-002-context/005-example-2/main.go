package main

import (
    "context" // 用於建立可取消的 Context
    "fmt"     // 用於格式化輸出
)

func main() {
    // 1) 建立一個可取消的 Context，並拿到 cancel 函式
    ctx, cancel := context.WithCancel(context.Background())
    // 2) 確保 main return 前一定呼叫 cancel()，避免 gen 的 goroutine 外洩
    defer cancel()

    // 3) 從 gen(ctx) 取得的 channel 中讀值
    //    gen 會持續產生 1, 2, 3, ... 直到 ctx 被取消
    for v := range gen(ctx) {
        fmt.Println(v)
        if v == 5 {
            // 4) 讀到 5 就跳出 for-loop
            //    但 channel 並未關閉，所以必須 cancel() 才讓 gen 停下
            break
        }
    }
    // 5) 因為 defer cancel()，在此離開 main 時會呼叫 cancel()
    //    觸發 gen 中的 <-ctx.Done() 分支，終止那個 goroutine
}

// gen 接受一個 Context，回傳只讀的 int channel：<-chan int
func gen(ctx context.Context) <-chan int {
    dst := make(chan int) // 建立無緩衝 channel，供外部接收
    n := 1                // 從 1 開始計數

    go func() {
        // 這個匿名 goroutine 會持續工作，直到 ctx 被取消
        for {
            select {
            case <-ctx.Done():
                // 當 ctx 被 cancel() 或過期時走到這裡
                // 直接 return 停止 goroutine，避免外洩
                return
            case dst <- n:
                // 否則就把 n 傳到 dst，然後 n++
                n++
            }
        }
    }()

    // 回傳給 caller：請注意，這條 channel 不會自己關閉
    // range gen(ctx) 必須搭配 break + cancel() 才能停止迴圈
    return dst
}
