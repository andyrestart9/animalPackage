package main

import (
	"fmt" // 提供格式化輸出功能，這裡用來打印錯誤訊息
	"os"  // 提供作業系統功能，用於開啟與關閉檔案
)

func main() {
	// 嘗試以唯讀模式開啟檔案 "xx.txt"
	// os.Open 會回傳 *os.File（檔案物件）與 error（錯誤資訊）
	f, err := os.Open("xx.txt")
	if err != nil {
		// 如果 err 不為 nil，代表開檔失敗（例如：檔案不存在或權限不足）
		// 將錯誤訊息印出
		fmt.Println("err happened:",err) // err happened:  open xx.txt: no such file or directory
		// 遇到錯誤後提前結束 main 函式，程式終止
		return
	}
	// defer 可延遲執行 f.Close()，直到 main 函式返回前才關閉檔案
	// 確保資源一定會被釋放，避免檔案遺漏未關閉
	defer f.Close()

}
