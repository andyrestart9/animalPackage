package main

import "fmt"

func main() {
	// ===== 1. 類型轉換 (Conversion) =====
	var i int = 42
	// 將 int 類型轉換成 float64
	var f float64 = float64(i)
	fmt.Printf("Conversion: i (int) = %d → f (float64) = %f\n", i, f)

	// ===== 2. 類型斷言 (Assertion) =====
	// 把任何值先放進一個 interface{}，模擬動態類型
	var x interface{} = "hello, world"

	// 帶 ok 的斷言 (安全用法，不會 panic)
	//   嘗試把 x 當成 string
	if s, ok := x.(string); ok {
		fmt.Println("Assertion with OK succeeded:", s)
	} else {
		fmt.Println("Assertion with OK failed")
	}

	// 若想示範 panic 並捕捉，可用下面函式 safeAssert，示範不匹配的類型斷言會 panic
	safeAssert(x)
}

// safeAssert 演示不帶 ok 的斷言失敗時會 panic，並以 recover 捕捉
func safeAssert(v interface{}) {
	// 註冊 defer，捕捉後續的 panic，避免程式崩潰
	defer func() {
		if r := recover(); r != nil {
			// recover() 會返回 panic 傳出的值；若不為 nil，代表確實發生了 panic
			fmt.Println("Recovered from panic in assertion:", r)
		}
	}()

	// 這裡斷言會失敗並 panic
	// 嘗試將 interface{} 內部的值斷言為 int
	// 由於 v 底層實際存放的是 string，不是 int
	// 執行到此處時，Go 會檢查動態類型：
	//   if v.dynamicType != int { panic(...) }
	// 因而觸發 panic，信息類似：
	//   interface conversion: interface {} is string, not int
	n := v.(int)

	// 這行永遠不會被執行，因為前一行已 panic
	fmt.Println("This will never print:", n)
}
