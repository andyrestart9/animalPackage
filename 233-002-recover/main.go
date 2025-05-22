package main

import "fmt"

// https://go.dev/blog/defer-panic-and-recover
func main() {
	// 呼叫 f()，若 f 中發生 panic 並被 recover，程式會繼續執行
	f()
	// f() 正常返回後印出訊息
	fmt.Println("Returned normally from f.")
}

func f() {
	// 註冊一個 defer 匿名函式，利用 recover() 捕捉 panic
	defer func() {
		// recover() 如果抓到 panic 的值，就不為 nil
		if r := recover(); r != nil {
			// 印出捕獲到的 panic 值
			fmt.Println("Recovered in f", r)
		}
	}()
	
	// 印出進到 g 函式的提示
	fmt.Println("Calling g.")
	// 呼叫遞迴函式 g，從 0 開始
	g(0)
	// 若 g 未 panic 或 panic 已被上層的 recover 捕獲，才會執行到這裡
	fmt.Println("Returned normally from g.")
}

func g(i int) {
	// 當遞迴層數 i 超過 3 時，觸發 panic
	if i > 3 {
		fmt.Println("Panicking!")
		// panic 時會中斷後續執行，並開始向上層展開堆疊
		panic(fmt.Sprintf("%v", i))
	}
	// 註冊 defer，當 g 返回或 panic 回溯時依 LIFO 執行
	// 會列印出當前的 i
	defer fmt.Println("Defer in g", i)
	
	// 印出當前的 i
	fmt.Println("Printing in g", i)
	// 遞迴呼叫 g，i 加 1
	g(i + 1)
}
