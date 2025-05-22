package main

import "fmt"

// Go 類型系統採用 三級分類：
// 基礎值類型：int/float/struct/array 等
// 引用語義類型：slice/map/channel
// 特殊值類型：interface/function pointer

// 指標類型 (Pointer Types)
// 存儲變量的內存地址，指向其他數據的引用，但指標本身是值類型（複製指標時複製地址值）
func intDeltaPointer(n *int) {
	*n = 43
}

// 值類型 (Value Types)
// 直接存儲數據本身，賦值或傳參時會複製完整數據副本：
// int, float32/64, bool, string, rune, byte 等。
func intDeltaValue(n int) {
	n++
	fmt.Println("n in intDeltaValue is:", n)
}

// **reference type（引用類型）**指數據結構本身不直接存儲值，而是通過內部指針間接引用底層數據。這類型的變量在賦值或傳遞時，共享同一份底層數據，修改會反映到所有引用該數據的變量。
// Slice, Map, Channel 等。
func sliceDelta(ii []int) {
	ii[0] = 99
}

func main() {

	a := 42
	fmt.Println(a)
	intDeltaPointer(&a)
	fmt.Println(a)
	fmt.Println("-----------------------")

	b := 55
	fmt.Println(b)
	intDeltaValue(b)
	fmt.Println(b)
	fmt.Println("-----------------------")

	xi := []int{1, 2, 3, 4}
	fmt.Println(xi)
	sliceDelta(xi)
	fmt.Println(xi)
}
