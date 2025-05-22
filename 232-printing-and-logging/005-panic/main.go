package main

import (
	"fmt" // 提供格式化輸出功能
	"os"  // 提供作業系統功能（如檔案操作）
)

func main() {
	// 註冊延遲呼叫 foo()
	// 無論 main 正常返回或發生 panic 進行堆疊展開，都會執行 foo()
	defer foo()

	// 嘗試以唯讀模式開啟 "xx.txt"，只取回錯誤值 err
	_, err := os.Open("xx.txt")
	if err != nil {
		// 若開檔失敗，呼叫 panic(err) 觸發恐慌（panic）
		// panic 會開始堆疊展開，並執行所有已註冊的 defer
		panic(err)
	}

	// 若成功開檔且未發生 panic，main 函式返回後，defer foo() 也會被執行
}

func foo() {
	// 印出提示文字
	// 注意：此訊息說明 os.Exit 不會執行 defer，但此範例中使用 panic，
	// panic 時 defer 仍會執行，因此可以看到這行輸出
	fmt.Println("When os.Exit() is called, deferred functions don't run")
}
