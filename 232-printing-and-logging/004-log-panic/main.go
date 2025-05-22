package main

import (
	"fmt" // 提供格式化輸出功能
	"log" // 提供日誌功能，包含 Panicln
	"os"  // 提供作業系統功能（如檔案操作）
)

func main() {
	// 註冊延遲呼叫 foo()，不論 main 正常返回或因 panic 而展開堆疊，都會執行
	defer foo()

	// 嘗試以唯讀模式開啟 "xx.txt"，只取回錯誤值
	_, err := os.Open("xx.txt")
	if err != nil {
		// 若開檔失敗，log.Panicln 會先輸出錯誤訊息，接著呼叫 panic()
		// panic 會觸發堆疊展開，並依序執行所有註冊的 defer
		log.Panicln("err happened:", err)
	}

	// 若沒有發生 panic，main 結束前也會執行 defer foo()
}

func foo() {
	// 印出提示：只有在調用 os.Exit 時，deferred 函式才不會執行
	fmt.Println("When os.Exit() is called, deferred functions don't run")
}
