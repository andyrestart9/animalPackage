package main

import (
    "context"   // 用於建立可取消的 Context
    "fmt"       // 格式化輸出
    "runtime"   // 查詢目前 goroutine 數量
    "time"      // 時間相關函式
)

func main() {
    // 1) 建立一個可取消的 Context，並取得對應的 cancel 函式
    ctx, cancel := context.WithCancel(context.Background())

    // 初始檢查：Err() 回傳 nil；目前只有 main goroutine + runtime 的底層
    fmt.Println("error check 1:", ctx.Err())
    fmt.Println("num goroutines 1:", runtime.NumGoroutine())

    // 2) 啟動背景工作 goroutine，它會定期印出「working n」
    go func() {
        n := 0
        for {
            select {
            case <-ctx.Done():
                // 當 ctx 被取消後，從 Done() channel 讀到值，結束 goroutine
                return
            default:
                // 尚未取消，繼續工作
                n++
                time.Sleep(200 * time.Millisecond)
                fmt.Println("working", n)
            }
        }
    }()

    // 3) 讓主 goroutine 阻塞 2 秒，讓工作 goroutine 有時間執行
    time.Sleep(2 * time.Second)
    // 取消前再次檢查：Err() 仍為 nil；goroutine 數量應為多了那支工作 goroutine
    fmt.Println("error check 2:", ctx.Err())
    fmt.Println("num goroutines 2:", runtime.NumGoroutine())

    // 4) 主動取消 context，觸發 ctx.Done()
    fmt.Println("about to cancel context")
    cancel()  // 這會關閉 Done() channel，讓背景 goroutine 停止
    fmt.Println("cancelled context")

    // 5) 再讓 main 阻塞 2 秒，觀察工作 goroutine 是否已退出
    time.Sleep(2 * time.Second)
    // 取消後檢查：Err() 回傳 context.Canceled；goroutine 數量回到原本（不含已退出的那支）
    fmt.Println("error check 3:", ctx.Err())
    fmt.Println("num goroutines 3:", runtime.NumGoroutine())
}
