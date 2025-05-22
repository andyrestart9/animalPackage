package main

import "fmt"

func main() {
	// 呼叫 a()，示範 defer 的立即參數評估與延遲執行
	a()
	// 分隔線
	fmt.Println("---")
	
	// 呼叫 b()，示範 defer 在迴圈中的註冊與 LIFO 執行順序
	b()
	// 分隔線
	fmt.Println("---")
	
	// 呼叫 c() 並接收回傳值，示範命名回傳值與 defer 修改返回參數
	i := c()
	// 印出 c() 的最終回傳值
	fmt.Println(i)
}

func a() {
	// 宣告區域變數 i，初始值為 0
	i := 0
	// 註冊 defer fmt.Println(i)
	// 立即評估參數 i（此時為 0），並延遲執行列印
	defer fmt.Println(i)
	// 將 i 加 1，變成 1
	i++
	// return 前會執行已註冊的 defer，列印出當時評估的參數值 0
	return
}

func b() {
	// 迴圈從 0 到 3，每次迭代都註冊 defer fmt.Println(i)
	// 立即評估當前 i 值並延遲列印
	for i := 0; i < 4; i++ {
		defer fmt.Println(i)
	}
	// b() 返回前，所有 defer 依後進先出順序執行，依序列印：3, 2, 1, 0
}

func c() (i int) {
	// 使用命名回傳值 i，初始值為 0
	// 註冊匿名 defer，在 return 之前執行 i++，作用於命名回傳值
	defer func() { i++ }()
	// 執行 return 1：
	// 1) 將命名回傳值 i 設為 1
	// 2) 執行 defer 內部 i++，使 i 變為 2
	return 1
}
