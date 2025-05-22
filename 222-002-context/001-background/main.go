package main

import (
	"context" // 引入 context 套件，用於跨多段呼叫傳遞取消、逾時與值
	"fmt"     // 引入 fmt 套件，用於格式化輸出
	"reflect"
)

func main() {
	// 建立一個根 Context：永不取消、不帶 deadline，也不攜帶任何值
	// 常用於 main、初始化或測試中作為所有衍生 Context 的起點
	ctx := context.Background()

	// 透過 Stringer 接口，也就是透過 context.go 的 func (backgroundCtx) String() string ，印出 ctx 的字串表示
	fmt.Println("context:\t", ctx)

	// 呼叫 Err() 檢查 Context 是否已取消或逾時
	// Background() 永遠不會被取消也沒有 deadline，所以 Err() 回傳 nil
	fmt.Println("context err:\t", ctx.Err())

	// 用 %T 印出 ctx 的動態型別，Context.Background() 回傳的是值型別 context.backgroundCtx
	// %T 在底層就是呼叫了 reflect.TypeOf 來取得值的動態型別，再把它轉成字串輸出。
	// 等同於：
	// t := reflect.TypeOf(x)
	// fmt.Println(t.String())
	// reflect.TypeOf(x) 回傳的是一個 reflect.Type 物件，程式裡可以對它做進一步的反射操作（像 t.Kind()、t.NumField() 等）。
	// %T 只是格式化輸出指令，直接印出型別名稱，不能在程式裡拿到那個 Type 物件來做其它事情。
	fmt.Printf("context type:\t%T\n", ctx)
	fmt.Println("context type:\t", reflect.TypeOf(ctx))

	// 印出分隔線，方便後面輸出區分
	fmt.Println("----------")

	fmt.Println("什麼是 reflect ？")
	// reflect.TypeOf(x) 回傳的是一個 reflect.Type 物件，程式裡可以對它做進一步的反射操作（像 t.Kind()、t.NumField() 等）。
	// 取得動態型別
	t := reflect.TypeOf(ctx)

	// 你可以問它是哪種「Kind」
	fmt.Println("Kind():", t.Kind()) // interface

	// 你可以問它的「名稱」與「包路徑」
	fmt.Println("Name():", t.Name())       // backgroundCtx
	fmt.Println("PkgPath():", t.PkgPath()) // context

	// 你可以看它的方法數，以及一個一個列出方法簽名
	fmt.Println("NumMethod():", t.NumMethod())
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Printf(" Method %d: %s%s\n", i, m.Name, m.Type)
	}

	// 如果它是指標或介面，你還可以取它的底下元素
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		elem := t.Elem()
		fmt.Println("Underlying element type:", elem)
	}

	/*
		在 Go 裡，每個變數都有一個靜態型別（static type），也就是你在程式中宣告時寫的那個型別；但如果這個變數是介面（interface）類型，那它在執行時會“攜帶”一個動態型別（dynamic type），代表了介面值裡實際存放的那個具體型別。

		靜態型別 vs 動態型別
		```
		var x interface{}       // x 的靜態型別是 interface{}
		x = "hello"             // 把一個 string 值賦給 x
		```
		靜態型別：x 在編譯期就被視為 interface{}，只有介面上定義的方法可用。
		動態型別：執行到賦值時，x 真正存放的是一個 string，所以它的動態型別就是 string。

		為什麼動態型別重要？
		反射
		reflect.TypeOf(x) 或 fmt.Printf("%T", x) 就是讀出 x 的動態型別。

		介面滿意度
		一個值即使是 interface，底下的具體型別決定了它能呼叫哪些方法。

		多型（Polymorphism）
		在函式接收 io.Reader、error、context.Context 這些介面時，動態型別告訴你實際使用的是哪一種實作。

		以 context 為例
		```
		ctx := context.Background()   // 滑鼠移到 ctx 上面會顯示 var ctx context.Context ，表示靜態型別是 context.Context 介面
		fmt.Println(ctx) // 印出 context.Background ，因為 Println 會檢查 ctx 是否實作了 fmt.Stringer 介面（也就是是否有 String() string 方法），有的話就印出 backgroundCtx 的 String() string 的返回值
		fmt.Printf("%T\n", ctx)      // 印出動態型別，例如 context.backgroundCtx
		```
		你宣告的靜態型別是 Context（介面），但 Background() 回傳的底層實作是 backgroundCtx，這就是它的動態型別。

		小結
		靜態型別：編譯時就確定（你在程式裡定義的那個型別）。
		動態型別：執行時，介面變數實際持有的那個具體型別。
		在非介面變數上，靜態型別就等於動態型別；在介面變數上，動態型別可以是多種具體型別之一。
	*/
	/*
	   静态类型（Static Type）
	   定义阶段就确定：当你写
	   ```
	   var ctx context.Context = context.Background()
	   ```
	   在编译时，编译器只知道 ctx 是 context.Context 这个「接口类型」。

	   可用的方法受限：因为静态类型是接口，编译器只允许你对 ctx 调用接口里声明的那四个方法：
	   ```
	   ctx.Deadline()
	   ctx.Done()
	   ctx.Err()
	   ctx.Value(key)
	   ```
	   即便在运行时它的具体实现还有别的方法，你也不能直接在代码里写 ctx.SomeOtherMethod()，否则编译不通过。

	   动态类型（Dynamic Type）
	   运行时才确定：在运行时，接口值 ctx 里实际装的是某个具体类型（concrete type）的实例。对于 context.Background() 来说，那个具体类型就是 backgroundCtx。
	   真正的方法实现来自动态类型：当你调用 ctx.Err() 时，它到底执行哪个 Err()，是看动态类型 backgroundCtx 的实现，而不是接口本身。

	   接口值的内部结构
	   一个接口值在内存里大致包含两部分：
	   类型信息（Type）：告诉运行时这个接口值底下是哪种具体类型，比如 backgroundCtx
	   数据指针/值（Value）：指向或存放那个具体类型的实际数据

	   因此：
	   编译期，你只能看到「Type 是 interface」，只准调用接口方法。
	   运行期，当接口被赋予一个具体值后，Type 字段就记录了它的 dynamic type，Value 字段存了那份数据；调用接口方法时，就是跳到这份动态类型上真正的方法实现去执行。
	   举例
	   ```
	   ctx := context.Background()      // static type: context.Context
	   //          ↑
	   //      编译器知道你只能调接口这四个方法

	   fmt.Printf("%T\n", ctx)
	   // 输出：context.backgroundCtx
	   //   ↑ 动态类型，运行时才知道它到底是哪种具体类型
	   ```
	   静态类型决定了「我在代码里能写哪些调用」。
	   动态类型决定了「这些调用在运行时会跑到哪份具体实现上」。
	   这就是接口在 Go 里多态（polymorphism）的核心：编译时看接口，运行时看具体类型。
	*/

	/*
		在 Go 中，「concrete type」（具體型別）指的就是實際的資料型別——不像介面（interface）那樣只定義方法簽名，具體型別會有真實的結構（struct）、基礎型別（int、string）、或是指向這些結構的指標（*MyStruct）等。

		跟介面的對比
		介面（interface）：只定義行為（方法集），是一種「抽象」／「契約」。
		```
		type Reader interface {
		    Read(p []byte) (n int, err error)
		}
		```

		具體型別（concrete type）：真正實作這些行為的方法，有實際的欄位或底層表示。
		```
		type File struct {
		    fd uintptr
		    // …
		}
		func (f *File) Read(p []byte) (n int, err error) { … }
		```

		當你有一個介面變數：
		```
		var r io.Reader = os.Stdin
		```
		靜態型別是 io.Reader（抽象介面）。
		動態型別則是被賦值的那個具體型別，在這裡就是 *os.File（具體實作了 Read 的型別）。

		常見的 Concrete Type 範例
		基本型別：int、string、bool
		結構型別：struct{ X, Y int }、type Person struct { Name string }
		指標型別：*Person、*MyStruct
		陣列、切片、映射：[]string、map[string]int
		這些都是真正能夠被「建立實例」(new、make、或字面量) 的型別，有具體的記憶體佈局和方法實作。

		為什麼重要？
		多態（Polymorphism）：介面只知道「我有這些方法簽名」，真正執行時要呼叫的就是那個 concrete type 上的方法實作。
		反射（Reflection）：reflect.TypeOf(x) 看到的是 concrete type，才能知道要做哪些底層操作。
		效能：理解 concrete type 才能知道記憶體分配、方法呼叫是值接收還是指標接收等細節。

		簡單來說，concrete type = interface 的實際實作，它們才是編譯器和 runtime 真正操作、分配記憶體、執行程式邏輯的「實體」。
	*/
	/*
		You said:
		type Reader interface 是介面
		concrete type 是
		func (f *File) Read(p []byte) (n int, err error) { … } 的 *File
		func (x *Xxx) Read(p []byte) (n int, err error) { … } 的 *Xxx
		的意思嗎？

		ChatGPT said:
		没错。你举的例子里：
		```
		type Reader interface {
		    Read(p []byte) (n int, err error)
		}
		```
		这是接口（interface），只定义了要有一个 Read 方法。
		而
		```
		func (f *File) Read(p []byte) (n int, err error) { … }
		```
		和
		```
		func (x *Xxx) Read(p []byte) (n int, err error) { … }
		```
		分别表明 *File 类型和 *Xxx 类型都“具体实现”了这个方法，因此它们就是 Reader 的 concrete type（具体类型）。换句话说：
		接口（Reader）：定义行为契约——需要一个 Read 方法。
		具体类型（*File、*Xxx）：真正写出 Read 方法的那几个类型，也就是运行时存放在接口变量里、真正执行方法体的那部分。
	*/
}

/*
context 就是 Go 提供的一個跨 goroutine 協調機制，它讓你能：

取消：從最上層發出取消訊號，所有 <-ctx.Done() 的 goroutine 都會立刻醒過來、結束工作。

逾時／截止：在同一個 Context 上設定 deadline 或 timeout，讓所有衍生的工作自動遵守最終時間限制。

攜帶值：在請求生命週期內傳遞共用資料（例如 Trace ID、使用者資訊），不必一層層改函式簽名。

換句話說，當你在多個 goroutine 間需要「一起停」、「一起過期」或「共享一些參數」時，就可以透過 context 來統一管理和協調。
*/
