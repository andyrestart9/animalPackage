package main

import (
	"fmt"
	"reflect"
)

/*
Compile-time（編譯期） vs Run-time（執行期）

編譯期：編譯器把你的源碼翻譯成可執行檔的那段過程。
發生在你執行 go build 或 go run 開始前半段。

編譯器會做：
語法檢查
靜態型別檢查（確保你只呼叫型別所允許的方法）
介面實作檢查（確保某個 concrete type 符合某介面）
最後產生機器碼或中間物件檔

執行期：可執行檔真正跑起來到跑完退出的那段過程。
包括：goroutine 調度、記憶體分配／釋放、反射调用、型別斷言、panic/recover、IO 操作等。
*/

// Shaper 定义了一个接口：静态类型在编译期就已确定
type Shaper interface {
	Area() float64
}

// Square 是一个具体类型，实现了 Shaper 接口
type Square struct{ side float64 }

// 值接收器方法：在运行期会根据动态类型调度到这里
func (s Square) Area() float64 { return s.side * s.side }

func main() {
	// ---------------- Compile-time（编译期） ----------------

	// 声明 s 为 Shaper 接口：静态类型是 Shaper ，
	// 滑鼠移到 s 上面會出現 var s Shaper 表示他的靜態型別
	var s Shaper // ← 编译期就知道 s 的静态类型是 Shaper（interface）

	// 编译器在这里检查 Square 是否实现了 Shaper（必须有 Area 方法）
	// 1) 编译期检查：
	//    编译器在这里会确认 Square 是否实现了 Shaper（也就是有没有签名为 Area() float64 的方法）。
	//    如果没有，编译就会报错：cannot use sq (type Square) as type Shaper in assignment
	sq := Square{4}
	s = sq // ← 编译器同意后，编译通过，把 Square 值赋给接口 s，编译期通过检查

	// 如果写成 s.Foo()，编译器会报错：undefined method Foo for Shaper

	// ---------------- Run-time（运行期） ----------------

	// 程序真正跑到 s = sq 这一行时，才会把那份 Square{4} 的值打包进接口变量 s：
	// 在运行期，接口值内部存了：
	//   dynamic type = Square
	//   dynamic value = Square{4}

	// 打印静态类型——接口 Shaper 本身
	staticType := reflect.TypeOf((*Shaper)(nil)).Elem()
	fmt.Println("Static type of s:", staticType) // 输出: Static type of s: main.Shaper

	// 打印动态类型——运行期才知道接口里实际是哪种 concrete type
	fmt.Printf("Dynamic type of s: %T\n", s) // 输出: Dynamic type of s: main.Square

	// 动态分发：根据动态类型调用对应的 Area() 实现
	fmt.Println("Area:", s.Area()) // 输出: Area: 16
}
