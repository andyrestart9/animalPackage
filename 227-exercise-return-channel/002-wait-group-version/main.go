package main

import (
    "fmt"
    "sync"
)
/*
Exam Questions: launches 10 goroutines, each goroutine adds 10 numbers to a channel
pull the numbers off the channel and print them
*/
func main() {
    x, y := 10, 10           // x = 工作者数量, y = 每个工作者发送的值数量
    c := gen(x, y)           // gen 会启动工作者并返回一个 channel

    // 从 channel c 中接收所有值直到其关闭
    // 我们知道 gen 会在发送完 x*y 个值后关闭 c
    for v := range c {
        fmt.Println(v)       // 打印来自每个工作者的值（0..y-1）
    }

    fmt.Println("所有值处理完毕，程序退出")
}

// gen 启动 x 个工作者 goroutine，每个发送 y 个值到同一个 channel。
// 使用 WaitGroup 确保当所有工作者完成后，由协调者关闭 channel。
func gen(x, y int) <-chan int {
    c := make(chan int)      // 无缓冲的整数 channel

    // 协调者 goroutine：启动工作者，等待它们完成后关闭 c
    go func() {
        var wg sync.WaitGroup
        wg.Add(x)            // 添加 x 个工作者到 WaitGroup

        // 启动 x 个工作者
        for i := 0; i < x; i++ {
            go func(id int) {
                defer wg.Done()  // 函数返回时表示该工作者完成

                // 每个工作者发送 y 个值：0, 1, ..., y-1
                for j := 0; j < y; j++ {
                    c <- j      // 发送时会阻塞，直到有接收者准备好
                }
            }(i)
        }

        wg.Wait()               // 等待所有工作者完成发送

        // 只有协调者负责关闭 channel，
        // 确保 close(c) 只执行一次且在所有发送完成后
        close(c)
    }()

    // 立即返回 channel，工作者在后台并行发送值
    return c
}
