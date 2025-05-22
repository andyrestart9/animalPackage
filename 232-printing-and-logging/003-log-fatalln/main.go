package main

import (
	"fmt" // 提供格式化輸出功能
	"log" // 提供日誌與致命錯誤輸出
	"os"  // 提供作業系統功能（如檔案操作）
)

func main() {
	// 註冊延遲執行 foo()，只有在 main 正常返回時才會調用
	// 若程式透過 os.Exit 或 log.Fatalln (內部呼叫 os.Exit) 結束，defer 不會執行
	defer foo()

	// 嘗試以唯讀模式開啟 "xx.txt"，只關心可能的錯誤
	_, err := os.Open("xx.txt")
	if err != nil {
		// 若開檔失敗，log.Fatalln 會先輸出錯誤訊息，然後呼叫 os.Exit(1) 直接終止程式
		log.Fatalln("err happened:", err)
	}

	// 如果成功開檔且 main 正常結束，defer foo() 才會在此後被執行
}

func foo() {
	// 提示：當呼叫 os.Exit 時，所有 defer 註冊的函式都不會執行
	fmt.Println("When os.Exit() is called, deferred functions don't run")
}

/*
os.Exit
Exit causes the current program to exit with the given status code.
Conventionally, code zero indicates success, non-zero an error.
The program terminates immediately; deferred functions are not run.
*/
