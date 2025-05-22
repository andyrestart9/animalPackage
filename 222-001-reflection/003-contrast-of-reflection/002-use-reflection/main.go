package main

import (
	"fmt"     // 用於格式化輸出
	"reflect" // 用於反射，檢查運行時型別和值
)

// 定義兩種不同的 struct，演示反射的泛用性
type Person struct {
	Name string // 人名
	Age  int    // 年齡
}

type Book struct {
	Title  string // 書名
	Author string // 作者
}

/*
為什麼反射好？
泛用性：一支 printStruct 就能處理 Person、Book、甚至 Order、User……所有結構體。

維護成本低：新增新型別時，不必再寫新的列印函式。

可讀結構標籤：更可擴充為讀 json、db 標籤，自動決定鍵名、核對欄位。

只要你需要「對不確定的多種型別做同樣的欄位操作」，反射就能極大減少重複程式碼。
*/

// printStruct 是一支通用函式，能在執行期動態讀取並列印任何 struct 的所有欄位
func printStruct(x interface{}) {
	// 取得傳入值的動態型別與反射值
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	// 如果 x 是指標 (pointer)，取得它所指向的元素
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // 取得指標指向的型別
		v = v.Elem() // 取得指標指向的值
	}

	// 只處理 struct 類型，否則警告並返回
	if t.Kind() != reflect.Struct {
		fmt.Println("只能列印 struct 或 *struct")
		return
	}

	// 印出 struct 的名稱
	fmt.Println("Type:", t.Name())

	// 迴圈遍歷所有欄位
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)             // reflect.StructField，包含欄位名稱、型別、標籤等
		value := v.Field(i).Interface() // reflect.Value 轉回 interface{}，取得實際值

		// 列印欄位名稱與欄位值
		fmt.Printf("  %s = %v\n", field.Name, value)
	}
}

func main() {
    // 建立 Person 與 Book 實例
    p := Person{"Alice", 30}
    b := Book{"1984", "Orwell"}

    // 同一支函式即可處理不同結構，展示反射的泛用性
    printStruct(p)
    fmt.Println("-----")
    printStruct(&b) // 傳入指標同樣支援
}