/*
在程式語言中，反射（Reflection）指的是程式在執行時（而非編譯時）能夠：

檢查變數的型別結構
取得某個值是什麼型別（int、struct、slice、interface…）
查詢該型別有哪些欄位、方法、標籤（Tag）等

動態讀取或修改值
不管這個值原本宣告為何種型別，都可以透過反射機制去讀出它的內容，甚至在允許的情況下修改它

動態呼叫方法
根據字串名稱呼叫某個物件的方法，而不需要在編譯期就寫死

Go 裡的反射：reflect 套件

reflect.TypeOf(x) → 回傳一個 reflect.Type，可以用來問
Kind()：它是哪一種基本分類（struct、ptr、slice…）
NumField()、Field(i)：如果是 struct，取出欄位資訊
NumMethod()、Method(i)：列出方法簽名

reflect.ValueOf(x) → 回傳一個 reflect.Value，代表那個具體的值，並可：
.Interface() → 轉回空介面 interface{}
.Int()/.String()/… → 如果是對應的基本型別，就能拿到真實值
.Set(...) → 如果這個 Value 可寫（addressable & settable），就能修改它的內容
.MethodByName("Foo") → 動態取得並呼叫名為 Foo 的方法

典型應用場景

序列化／反序列化
JSON、XML、ORM 等套件需要把任意結構體轉成字串或從字串還原成物件。

框架與工具
Web 框架（自動綁定 JSON 請求到結構體）、DI（依賴注入）、測試框架（自動跑測試、生成文件）等，都依賴反射來發現程式的結構和標籤。

通用函式庫
比如一個通用的 Validate(obj interface{})，它可以在執行時讀取結構體的 validate:"required" 標籤，依據不同型別做檢查。

總結
反射是一把「在執行期打開型別和記憶體結構」的利器，能讓程式更靈活、更通用，但也要謹慎使用，留意其性能與可維護性的代價。
*/
/*
「在執行期打開型別和記憶體結構」的意思是：透過反射，你可以在程式執行時，不必在編譯時就知道資料的具體型別，就能：

檢視任意變數的結構
例如傳進來的是一個 struct，你可以列出它有哪些欄位、欄位名稱、型別、標籤（tag）。
甚至可以算出它在記憶體裡的排布（例如欄位偏移、大小、對齊），雖然 Go 的 reflect 並不直接暴露所有偏移細節，但你可以看到 reflect.StructField.Offset。

讀寫底層的值
用 reflect.ValueOf(x) 拿到一個可讀寫的 Value，可以 .Interface() 取回原始值，或者（在可寫入時）用 .Set(...) 修改它。
例如，你可以寫一段程式，對任何結構體欄位加上前綴字串、或者讓所有非匿名字串欄位都清空。

動態呼叫方法
reflect.ValueOf(x).MethodByName("Foo") 可以在不知道 x 是哪種型別的情況下，用字串去呼叫它的某個方法。
*/
package main

import (
	"fmt"     // 用於格式化輸出
	"reflect" // 反射，用於在執行期檢查型別與欄位
)

// Person 是一個結構示例，用於示範反射
type Person struct {
	Name string `myTag:"name" secondTag:"姓名"`             // 這是在 Go struct 欄位後面的「結構標籤」（struct tag）。它其實就是一個純文字的 metadata，語法是放在反引號 `…` 裡面
	Age  int    `myTag:"age" secondTag:"number of years"` // json:"age" 也是結構標籤
}

// inspect 使用反射來檢查傳入物件的型別與欄位值
func inspect(i interface{}) {
	t := reflect.TypeOf(i)  // 取得變數 i 的動態型別 (reflect.Type)
	v := reflect.ValueOf(i) // 取得變數 i 的反射值 (reflect.Value)

	// 在 reflect 裡，Elem() 代表「把外層包裝打開，拿到裡面的元素 (element)」
	// var n int = 7	直接存值	0x1000 → 7
	// var p *int = &n	存位址	0x2000 → 0x1000 (再去 0x1000 拿到 7)
	// n 直接把 7 放在自己的位址 0x1000。
	// p 存在 0x2000，裡面放的是「另一個位址 0x1000」。
	// 換句話說：p 只是一張「地圖」，真正的 7 還躺在別處。
	// 因此 p 的值不是 7，而是「通往 7 的路標」——這就叫「包了一層」。
	// Elem() 每呼叫一次就會拆最外面那一層

	// 如果傳入的是指標 (pointer)，則取其所指向的元素型別和值
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // 取得指標指向的元素型別
		v = v.Elem() // 取得指標指向的元素值
	}

	// 印出結構體的名稱，例如 "Person"
	fmt.Println("Type:", t.Name())

	// 迴圈列出這個結構體的所有欄位
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)               // StructField，包含欄位名稱、型別、標籤等資訊
		val := v.Field(i).Interface() // 取得欄位的實際值
		tag := f.Tag

		// %q 會把字串用雙引號括起並做跳脫
		fmt.Printf(" Field %q (%s) = %v, tag=%q,\t myTag=%q, secondTag=%q\t\n",
			f.Name,                 // 欄位名稱
			f.Type,                 // 欄位型別
			val,                    // 欄位值
			tag,                    // 讀取 json 標籤
			f.Tag.Get("myTag"),     // 讀取 myTag 標籤
			f.Tag.Get("secondTag"), // 讀取 myTag 標籤
		)
	}
}

/*
json:"name"
這是在 Go struct 欄位後面的「結構標籤」（struct tag）。它其實就是一個純文字的 metadata，語法是放在反引號 `…` 裡面，格式大致長這樣：
```
type User struct {
    Name string `json:"name"`     // 這裡是 struct tag
}
```
當你用 encoding/json 這類依賴反射的套件去序列化/反序列化時，套件會讀取這段 tag：
json:"name" 告訴它「把這個欄位在 JSON 裡對應的鍵名叫做 name」
預設沒有 tag 時，鍵名會和欄位名稱一樣（例如 "Name"）
你也可以加上選項，像 json:"name,omitempty"（空字串/0值就不輸出）

%q
這是一個 fmt 的格式化動詞，用來「把字串用雙引號括起來並做必要的跳脫」：
```
s := "Hello\nWorld"
fmt.Printf("%s\n", s)  // 直接輸出：Hello
                        //           World
fmt.Printf("%q\n", s)  // 輸出： "Hello\nWorld"
```
%s：原樣輸出字串內容，遇到 \n 就換行
%q：輸出 Go 語法風格的字面值（literal），在前後加上雙引號，並將換行、制表符等做 \n、\t 的跳脫。
*/

/*
這段程式在執行期「打開」了 Person 這個型別，看它有哪些欄位、型別、標籤，並讀出每個欄位的值。

完全不需要在編譯期就寫死 Person，同理也能對任何其他結構體做相同處理。

總結
打開型別：用反射取得一個值的 Type，看它是 struct、slice、map、interface…，並可以列出欄位或方法。
打開記憶體結構：用反射的欄位資訊（StructField.Offset、Field(i).CanSet()）和 Value API，讀取或修改底層資料。
這就是所謂「在執行期打開型別和記憶體結構」的涵義：反射讓程式能夠動態探索和操作它本來看不見的「內部」！
*/

func main() {
	p := Person{"Alice", 30}
	fmt.Println(p.Name)
	fmt.Println(p.Age)
	inspect(&p)
}

/*
You said:
為什麼要用反射？
	fmt.Println(p.Name)
	fmt.Println(p.Age)
不是也可以拿到姓名年齡嗎？

ChatGPT said:
直接用 p.Name、p.Age 當然可以，而且在你已經知道 Person 型別、要硬寫字段存取時，直接取用永遠是最簡單、最安全、也最快的做法——不需要反射，也不用動態解析。

那為什麼還要用反射？ 反射的價值在於「不事先知道具體型別、還能在執行期動態處理各種不同結構」：

通用化／自動化程式庫
比如 JSON 序列化，你不會為每一個 struct 寫一大串 MarshalJSON，而是用一套反射機制自動把所有字段和 tag 拿出來。
ORM、表單綁定、依賴注入這類框架，也都得在執行期「掃描」任何傳進來的 struct，根據 tag 自動把資料塞進去。

動態場景
你可能從外部取得一個 interface{}，或從 plugin、配置載入某個型別，直到執行時才知道它是什麼。這時就沒辦法 p.Name ——因為編譯期你根本不知道它是不是有 Name 字段。
用反射就能在運行中：「這個值是 struct？讓我列出它所有字段，再動態取值或改值」。

Tag 驅動行為
結構標籤（json:"name,omitempty"、db:"user_id"、validate:"required"……）不會影響編譯，也不會變成字段。唯有反射能讀出這些 tag，根據 tag 決定序列化鍵名、資料驗證規則或資料庫映射。

減少樣板碼
假設你有 50 種 struct，都要印出所有字段、檢查 tag 或做某種轉換，如果都要手寫，維護成本爆表。反射讓你只寫一次泛用函式，就能自動擴展到所有型別。

總結：
反射：當你需要寫「對任何型別都有效」或「根據執行期資訊（如 tag、未知結構）」來動態操作時，才用得上。它帶來靈活性和通用性，但也犧牲了可讀性、型別安全與效能。


You said:
不事先知道具體型別、還能在執行期動態處理各種不同結構是什麼意思？


ChatGPT said:
舉個最常見的例子：你想寫一個「把任意結構體序列化成 JSON」的函式，你根本不可能一開始就把所有可能的 struct 型別都寫進去像這樣：
```
func Marshal(v interface{}) ([]byte, error) {
    if u, ok := v.(User); ok {
        // 針對 User 結構寫一大段序列化程式
    } else if o, ok := v.(Order); ok {
        // 針對 Order 結構寫另一段序列化程式
    } else if p, ok := v.(Product); ok {
        // …
    }
    // 無窮無盡……
}
```
你並不知道最終會傳進來的是哪一個 struct，也不可能事先寫好所有分支。這時候，就必須在執行期看一看：

這個 v 究竟是個 struct？
它有多少個欄位？欄位叫什麼、型別是什麼？
每個欄位又帶了哪些 tag（像 json:"name"、db:"user_id"）？

透過 reflect.TypeOf(v) 和 reflect.ValueOf(v)，你就能動態：

找出它是哪一個 struct（或 slice、map、基本型別），
用迴圈把所有欄位一個一個讀出來（t.NumField()、v.Field(i)），
根據 tag 決定鍵名、或者做某些轉換，
最後把結果拼成 JSON 字串。

這整個過程，你事先完全不用知道傳進來的是 User 還是 Order，只要符合「有結構體欄位」這件事，就能同一套程式處理所有型別，這就是「不事先知道具體型別、還能在執行期動態處理各種不同結構」的意義。

— 如果不用反射，你就必須為每一種可能的 struct 手動寫一段序列化邏輯；但有了反射，你只要寫一次「泛型」迴圈，就可以通吃所有 struct，大大減少重複程式碼，也更容易維護。
*/
