package main

import "fmt"

type dog struct {
	first string
}

func (d dog) walk() {
	fmt.Println("My name is", d.first, "and I'm walking.")
}

// 兩種寫法等價但指針版本更簡潔
// func (d *dog) run() { d.first = "Rover" }      // 自動解引用
// func (d *dog) run() { (*d).first = "Rover" }   // 手動解引用

func (d *dog) run() {
	d.first = "Rover"
	fmt.Println("My name is", d.first, "and I'm running.")
}

func main() {
	d1 := dog{"Henry"}
	d1.walk()
	d1.run()


	d2 := &dog{"Padget"}
	d2.walk()
	d2.run()
}

// Go指針接收器的自動解引用機制
// 1. 語法糖設計原則
// Go語言對結構體指針設計了自動解引用語法糖，主要基於以下考量：

// 減少星號(*)操作符的視覺干擾
// 統一值/指針接收器的訪問方式
// 避免開發者混淆值傳遞與指針傳遞的邊界

// type dog struct { first string }

// // 兩種寫法等價但指針版本更簡潔
// func (d *dog) run() { d.first = "Rover" }      // 自動解引用
// func (d *dog) run() { (*d).first = "Rover" }   // 手動解引用

// 2. 底層操作對比
// 編譯器處理流程：
// d := &dog{"Henry"}
// d.run()

// // 實際展開步驟：
// 1. 取地址：_dptr := d
// 2. 自動插入解引用：(*_dptr).first = "Rover"
// 3. 生成機器碼：MOVQ $_dptr, AX; MOVQ "Rover", (AX)

// 3. 操作符優先級陷阱
// 錯誤示範的編譯錯誤分析：
// *d.first = "Rover"  // 編譯錯誤：invalid indirect of d.first (type string)
// 等價於：
// *(d.first) = "Rover"  // 對string類型執行解引用操作

// 正確解引用應寫為：
// (*d).first = "Rover"  // 但Go不需要手動寫

// 4. 指針接收器的特殊行為
// 三種等效訪問方式：
// func (d *dog) run() {
//     d.first = "A"      // 語法糖模式
//     (*d).first = "B"   // 顯式解引用
//     (&*d).first = "C"  // 先解引用再取地址
// }
// // 三種寫法最終操作同一塊記憶體