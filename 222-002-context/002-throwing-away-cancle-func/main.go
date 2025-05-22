package main

import (
    "context" // 引入 context 套件，用於傳遞取消訊號、截止時間或攜帶值
    "fmt"     // 格式化輸出
)

func main() {
    // 建立一個根 Context：永不取消、不帶 deadline，也不攜帶任何值
    ctx := context.Background()

    // 印出 ctx 的 String() 表示（透過 Stringer 介面），通常是 "context.Background"
    fmt.Println("context\t", ctx)
    // 檢查 ctx 是否已經被取消或超時；Background 永遠回傳 nil
    fmt.Println("context err:\t", ctx.Err())
    // 用 %T 印出 ctx 的動態型別；Background 回傳 backgroundCtx（值型別）
    fmt.Printf("context type:\t%T\n", ctx)
    fmt.Println("-----------")

    // 對根 Context 建立一個可取消的衍生 Context
    // ctx.WithCancel 會回傳新的 Context 與對應的 cancel 函式
    // 這裡忽略了 cancel 函式 (用 _ 接收)，所以後續不會手動取消
    ctx, _ = context.WithCancel(ctx)

    // 現在 ctx 是新的可取消 Context
    // 它內部包了：
    //   parent = backgroundCtx
    //   cancelFunc = 一個能觸發取消的函式
    fmt.Println("context:\t", ctx)
    // Err() 仍然是 nil，因為我們還沒呼叫 cancel() 也沒逾時
    fmt.Println("context err:\t", ctx.Err())
    // 動態型別變成了 *cancelCtx，表示這是一個可取消的 Context
    fmt.Printf("context type:\t%T\n", ctx)
    fmt.Println("----------")
}
/*
context.Background()
返回一個空的根 Context，永不取消、也不帶任何值或 deadline。
動態型別通常是 backgroundCtx（值型別）。

context.WithCancel(parent)
創建一個新的 Context（動態型別為 *cancelCtx），它會監聽父 Context。
同時返回一個 cancel() 函式，呼叫後會使這個新 Context 收到取消訊號，並讓所有 <-ctx.Done() 解除阻塞、ctx.Err() 返回 context.Canceled。

靜態 vs 動態型別
靜態型別（ctx 宣告時）永遠是 context.Context 介面。
動態型別（執行期）第一階段是 backgroundCtx，第二階段是 *cancelCtx。
*/