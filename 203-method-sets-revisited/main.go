package main

import (
	"fmt"
	"math"
)

type circle struct {
	radius float64
}

type shape interface {
	area() float64
}

func (c *circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func info(s shape) {
	fmt.Println("area:", s.area())
}

func main() {
	c := circle{5}
	// 只有指標型別 type *circle 實現 type shape interface，值型別 type circle 沒有實現 type shape interface
	// 所以下面這行或出現錯誤 cannot use c (variable of type circle) as shape value in argument to info: circle does not implement shape (method area has pointer receiver)
	// 詳細說明看下面的第二則註解
	// info(c)

	// 為什麼 func (c *circle) area() float64 是指標接收器，但是 c.area() 可以運作，而不用寫成 (&c).area() ?
	// 因為這裡的 c 是一個可尋址的值（addressable value），當你呼叫方法時，Go 編譯器會自動將它的地址傳給該方法。也就是說，c.area() 實際上會被轉換為 (&c).area()
	// 可尋址的（Addressable）的表達式：變數、結構體欄位、陣列或切片元素
	// 不可尋址的（Non-addressable）的表達式：常量與字面值、臨時值、函數返回值、map 的元素、直接從通道接收的值
	// 詳細可以看下面第一個註解
	fmt.Println(c.area())
}

/*
一、方法呼叫的自動取址

說明
如果某個方法是用指標接收器定義的（例如 func (c *circle) area() float64），那麼只有指標型別（*circle）的方法集中包含此方法。
但如果你擁有一個可尋址的值（addressable value），例如變數 c，即使它的型別是 circle（值型別），當你呼叫方法時，Go 編譯器會自動將它的地址傳給該方法。
也就是說，c.area() 實際上會被轉換為 (&c).area()。
範例
package main

import (
	"fmt"
	"math"
)

type circle struct {
	radius float64
}

// 方法定義在 *circle 上（指標接收器）
func (c *circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func main() {
	// c 為 circle 值，且是個變數（可尋址的）
	c := circle{5}

	// 雖然 c 的型別是 circle，但因為 c 是 addressable，
	// 編譯器自動將 c 的地址傳入，所以這裡相當於 (&c).area()
	fmt.Println("Area:", c.area())
}
在上面的例子中，由於 c 是一個變數（可尋址的），因此即使 area() 是定義在 *circle 上，我們仍可以直接使用 c.area()。

二、介面實作的注意點

說明
介面（interface）定義了一組方法的簽名，但不關心方法是如何實作的。
如果介面要求的方法是用指標接收器實作的，那麼只有該類型的指標型別才會擁有這個方法。
自動取址機制只適用於方法呼叫時。當你把一個值賦給介面變數時，不會自動取址，介面的判斷是根據變數的方法集來決定是否符合。
範例
package main

import (
	"fmt"
	"math"
)

type shape interface {
	area() float64
}

type circle struct {
	radius float64
}

// 使用指標接收器定義 area() 方法，只有 *circle 方法集有此方法
func (c *circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

// info 需要傳入實作了 shape 介面的值
func info(s shape) {
	fmt.Println("Area from interface:", s.area())
}

func main() {
	c := circle{5}

	// 呼叫 c.area() 時，由於 c 是 addressable，
	// 自動轉換為 (&c).area()，因此這裡是可以呼叫的。
	fmt.Println("Direct call:", c.area())

	// 但若嘗試將 c（值型別）賦值給 shape 介面，則會出現錯誤，
	// 因為 circle 的方法集中不包含指標接收器定義的 area() 方法。
	// info(c) // 這行會編譯失敗

	// 正確做法：傳入指標 &c，因為 *circle 方法集中包含 area() 方法
	info(&c)
}
在這個例子中：

直接呼叫 c.area() 能夠正常工作，因為自動取址機制在方法呼叫時生效。
但將 c 賦值給介面變數 s（例如執行 info(c)）會產生錯誤，因為 c（circle 值）的方法集中沒有 area() 方法；只有 *circle 的方法集中才有此方法，所以必須使用 &c。
三、哪些值是可尋址的？哪些不是？

可尋址的（Addressable）的表達式
變數：直接宣告的變數有明確的記憶體位置。
x := 10
fmt.Println(&x) // 可取址
結構體欄位：若結構體變數可尋址，其欄位也可尋址。
type Person struct {
    Name string
}
p := Person{Name: "Alice"}
fmt.Println(&p.Name) // 可取址
陣列或切片元素：如果陣列或切片本身是可尋址的，其元素也可取址。
nums := []int{1, 2, 3}
fmt.Println(&nums[0]) // 可取址
不可尋址的（Non-addressable）的表達式
常量與字面值：沒有固定記憶體位置。
// fmt.Println(&42)   // 編譯錯誤
// fmt.Println(&"hi") // 編譯錯誤
臨時值（例如運算結果）：
a := 10
// fmt.Println(&(a + 5)) // 編譯錯誤，a+5 是臨時值
函數返回值（直接取值）：
func getNum() int { return 100 }
// fmt.Println(&getNum()) // 編譯錯誤，getNum() 的返回值不是可尋址的
map 的元素：
m := map[string]int{"key": 1}
// fmt.Println(&m["key"]) // 編譯錯誤，map 元素不可取址
直接從通道接收的值：
ch := make(chan int, 1)
ch <- 50
// fmt.Println(&(<-ch)) // 編譯錯誤，通道接收表達式不可取址
總結

方法呼叫自動取址：
當你呼叫定義在指標接收器上的方法（例如 c.area()），只要 c 是一個可尋址的變數，編譯器會自動將其地址傳入，等同於 (&c).area()。
介面實作注意點：
自動取址機制只在方法呼叫時生效，介面的判斷是根據變數的方法集來決定。如果介面方法是用指標接收器定義的，只有指標型別才符合要求，因此必須傳入 &c 而非 c。
哪些值是可尋址的：變數、結構體欄位、陣列/切片元素、指標解引用的結果等是可尋址的；而常量、臨時值、函數返回值（未存變數）、map 的元素、通道接收表達式則不是可尋址的。
*/

/*
從頭理解什麼是介面方法，以及為什麼用指標接收器與值接收器實作相同介面方法，最終會影響到該類型是否滿足介面。

1. 介面（interface）和介面方法

在 Go 語言中，介面定義了一組方法的簽名（名稱、參數、返回值）。例如：

type shape interface {
    area() float64
}
這裡的 area() float64 就是一個介面方法，它只描述了「契約」：任何實現 shape 介面的類型都必須有一個 area 方法，該方法無參數並返回一個 float64。

注意： 在介面定義中並沒有說明該方法是用指標接收器還是值接收器來實現，這是由具體類型在實作時決定的。

2. 方法集（Method Set）的概念

在 Go 中，每個類型都有一個「方法集」，這個集合決定了該類型擁有哪些方法，也就是說，只有在方法集中出現的方法才被認為是這個類型所實現的。

值接收器：如果一個方法是用值接收器來定義，那麼該方法屬於該類型的值的方法集。同時，指標類型也會「繼承」這個方法。
例如：

func (c circle) area() float64 {
    return 3.14 * c.radius * c.radius
}
在這種情況下，circle 的值擁有 area 方法，而 *circle 也能調用此方法。
指標接收器：如果一個方法是用指標接收器定義的，那麼該方法只屬於指標類型的方法集，而不屬於值類型。
例如：

func (c *circle) area() float64 {
    return 3.14 * c.radius * c.radius
}
在這種情況下，只有 *circle 擁有 area 方法，而 circle 的值不具備該方法。
這就解釋了為什麼即使介面只定義了方法簽名，不指定接收器形式，但實際上我們在實現時必須注意使用何種接收器。

3. 完整範例說明

範例 1：使用指標接收器實作
在這個例子中，我們用指標接收器定義 area() 方法：

package main

import (
	"fmt"
	"math"
)

type shape interface {
	area() float64
}

type circle struct {
	radius float64
}

// 使用指標接收器定義方法
func (c *circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func main() {
	var s shape

	c := circle{radius: 5}

	// 直接傳入 circle 值會錯誤，因為 circle 的方法集中沒有 area() 方法
	// s = c // 編譯錯誤

	// 正確：傳入 circle 的指標，*circle 的方法集包含 area() 方法
	s = &c
	fmt.Println("Area:", s.area())
}
說明：

定義介面 shape 要求有一個 area() 方法。
在 circle 上使用指標接收器實作 area() 方法，因此只有 *circle 擁有這個方法。
在 main 中，c 是 circle 值，直接將 c 賦值給 shape 介面會出錯，因為 circle 的方法集中沒有 area() 方法。必須傳入 &c（即 *circle）才正確。
範例 2：使用值接收器實作
這次我們用值接收器定義 area() 方法：

package main

import (
	"fmt"
	"math"
)

type shape interface {
	area() float64
}

type circle struct {
	radius float64
}

// 使用值接收器定義方法
func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func main() {
	var s shape

	c := circle{radius: 5}

	// 直接使用 circle 的值就可以，因為 circle 的方法集中包含 area() 方法
	s = c
	fmt.Println("Area (value):", s.area())

	// 同時，指標類型 *circle 也能調用此方法
	s = &c
	fmt.Println("Area (pointer):", s.area())
}
說明：

同樣定義了介面 shape。
在 circle 上使用值接收器實作 area() 方法，這樣 circle 的值就擁有該方法，且 *circle 也能調用該方法。
在 main 中，既可以將 c（值）賦予介面，也可以將 &c（指標）賦予介面，兩者都能正確呼叫 area() 方法。
4. 總結

介面方法 只是介面定義中的方法簽名，不涉及指標或值的選擇。
方法實作 的方式（指標接收器或值接收器）會影響到具體類型的方法集，進而決定該類型（或該類型的指標）是否能實現介面：
如果方法是指標接收器，只有該類型的指標會有該方法，只有指標類型能滿足介面。
如果方法是值接收器，則該方法屬於值和指標兩個方法集，兩者都能滿足介面要求。

*/