package main

import (
	"fmt" // 提供格式化輸出功能，用於在終端顯示錯誤
	"log" // 提供日誌功能，這裡用於將錯誤寫入檔案
	"os"  // 提供作業系統功能，用於檔案的建立、開啟與關閉
)

func main() {
	// 建立（或覆寫）log.txt 檔案，若檔案不存在則建立，存在則清空
	f, err := os.Create("log.txt")
	if err != nil {
		// 如果建立檔案失敗（例如：權限不足、磁碟滿），就在終端印出錯誤
		fmt.Println(err)
	}
	// 延遲關閉 log.txt，確保 main 結束前檔案一定會被關閉
	defer f.Close()

	// 將標準日誌輸出導向到 log.txt，之後呼叫 log.Println 都會寫入此檔
	log.SetOutput(f)

	// 嘗試以唯讀模式開啟 xx.txt，只取回錯誤值
	f2, err := os.Open("xx.txt")
	if err != nil {
		// 如果開檔失敗，將錯誤訊息寫入 log.txt
		log.Println("err happened:", err)
	}
	// 延遲關閉 xx.txt，確保若成功開啟也能正確關閉
	defer f2.Close()

	// 提示使用者檢查當前目錄下的 log.txt 內容
	fmt.Println("check the log.txt file in the directory")
}
