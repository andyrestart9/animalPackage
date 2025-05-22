package main

import (
	"fmt" // 提供格式化輸出功能
	"log" // 提供日誌功能，用於列印錯誤
	"math" // 提供數學函式庫，包含 Sqrt
)

// 定義自訂錯誤類型，包含錯誤名稱與原始錯誤
type wrongMathError struct {
	name string // 錯誤名稱
	err  error  // 原始錯誤
}

// 為 wrongMathError 實作 Error() 方法，以符合 error 介面
func (w wrongMathError) Error() string {
	// 回傳結合 name 與原始 err 的格式化錯誤訊息
	return fmt.Sprintf("a wrong math error occurred: %v %v", w.name, w.err)
}

func main() {
	// 呼叫 sqrt 嘗試計算 -10 的平方根，只關心錯誤
	_, err := sqrt(-10)
	if err != nil {
		// 若 err 不為 nil，將錯誤寫入日誌
		log.Println(err)
	}
}

// sqrt 計算輸入 f 的平方根，如果 f < 0，回傳自訂錯誤
func sqrt(f float64) (float64, error) {
	if f < 0 {
		// 建立一個描述性原始錯誤
		wme := fmt.Errorf("wrong math redux: square root of negative number: %v", f)
		// 回傳 0 及 wrongMathError 實例
		return 0, wrongMathError{"sqrt err", wme}
	}
	// 正常情況下，使用 math.Sqrt 計算並回傳結果，err 為 nil
	return math.Sqrt(f), nil
}
