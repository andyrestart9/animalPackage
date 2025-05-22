package main

import (
	"context" // 用於建立及控制 Context
	"fmt"     // 用於格式化輸出
)

func main() {
	// 1) 建立一個根 Context：永不取消、不帶 deadline 也不攜帶值
	ctx := context.Background()

	// 此時 ctx.String() 會印出 "context.Background"
	fmt.Println("context\t", ctx)
	// Err() 回傳 nil：Background 永遠不會被取消或逾時
	fmt.Println("context err:\t", ctx.Err())
	// 動態型別（Dynamic Type）為 backgroundCtx
	fmt.Printf("context type:\t%T\n", ctx)
	fmt.Println("-----------")

	// 2) 從根 Context 衍生出一個可取消的 Context，並取得 cancel 函式
	ctx, cancel := context.WithCancel(ctx)

	// 新的 ctx 仍是介面，靜態型別是 context.Context
	// 動態型別變成 *cancelCtx，代表它可被取消
	fmt.Println("context:\t", ctx)
	fmt.Println("context err:\t", ctx.Err()) // 仍然 nil，尚未呼叫 cancel
	fmt.Printf("context type:\t%T\n", ctx)   // *context.cancelCtx

	// cancel 是一個函式，可用來觸發取消
	fmt.Println("cancel:\t\t", cancel)
	fmt.Printf("cancel type:\t%T\n", cancel) // context.CancelFunc
	fmt.Println("----------")

	// 3) 呼叫 cancel()，主動取消該 Context
	cancel()

	// 一旦取消，ctx.Err() 會回傳 context.Canceled
	fmt.Println("context:\t", ctx)
	fmt.Println("context err:\t", ctx.Err()) // context.Canceled
	fmt.Printf("context type:\t%T\n", ctx)   // 仍然 *context.cancelCtx

	// cancel 函式本身沒有變，仍可重複呼叫（後續不會再改變 Err）
	fmt.Println("cancel:\t\t", cancel)
	fmt.Printf("cancel type:\t%T\n", cancel)
	fmt.Println("----------")
}

/*
context.Background()
回傳一個「空」的根 Context，不會取消也不帶值或逾時。

context.WithCancel(parent)
生成一個可取消的子 Context（動態型別為 *cancelCtx），並回傳一個 CancelFunc 給你用來觸發取消。

cancel()
呼叫後會立即將該 Context 以及所有以它為父的子 Context 標記為取消，Err() 開始回傳 context.Canceled，Done() channel 也會關閉。

靜態型別 vs 動態型別
靜態型別：ctx 宣告時就是 context.Context 介面，編譯時就固定。
動態型別：執行期隨賦值不同而變——先是 backgroundCtx，再是 *cancelCtx。
*/
