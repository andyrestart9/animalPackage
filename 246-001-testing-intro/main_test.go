// 測試檔案命名必須以 _test.go 結尾
// 測試程式與被測程式在同一個 package
// 測試函式必須是 Test 為前綴

// 在 Go 裡，一個套件（package）對應一個目錄，裡面所有 .go 檔如果都宣告了同樣的 package 名稱，就屬於同一個套件。
// go test 會對每個套件分別編譯、執行同目錄下所有 _test.go 檔案裡的測試。
package main

import (
    "testing"
)

// TestMySum 用來測試 mySum 函式是否能正確地將兩個整數相加
// 函式名稱必須以 Test 為前綴（後面跟要測試的函式或描述），
// 這樣 go test 才會自動辨識並執行它
func TestMySum(t *testing.T) {
    // 這裡的參數 t *testing.T：
    // - t 是參數名稱，可以隨意取，但習慣用 t
    // - *testing.T 是指向 testing.T 結構體的指針
    //   1. testing.T 提供了 Error、Fail、Log…等方法
    //   2. 用指針可以讓這些方法修改同一份測試狀態（標記失敗、緩存日誌等）
    //   3. 如果用值類型（testing.T），方法改動只會作用在拷貝上，框架感知不到

    // 呼叫要測試的函式 mySum，傳入 2 和 3，並將結果存到 x
    x := mySum(2, 3)
    
    // 檢查回傳值是否等於預期的 5
    if x != 5 {
        // 如果不相等，就報告錯誤：顯示「Expected 5 Got x」
        // t.Error 會記錄錯誤並標記測試失敗，但不會停止後續測試
        t.Error("Expected", 5, "Got", x)
    }
}
/*
為什麼 t.Error 用 “go to definition” 找到的是 func (c *common) Error(args ...any) ？

因为：

在 src/testing/testing.go 中，T 类型是这样声明的——它匿名嵌入了一个名为 common 的结构体：
```
// T is a type passed to Test functions to manage test state…
type T struct {
    common         // 匿名字段
    isParallel bool
    isEnvSet   bool
    context    *testContext
}
```
通过这种匿名嵌入（embedding），common 上定义的所有方法都会“提升”（promote）到 T 上，使得 *T 实例可以直接调用 Error、Fail、Log 等方法 
Reddit
。
*/