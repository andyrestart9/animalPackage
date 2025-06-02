package main

import "fmt"

// 定義一個基底型別 Animal，帶一個方法 Speak
type Animal struct {
    Name string
}

func (a *Animal) Speak() {
    fmt.Println(a.Name, "makes a sound")
}

// 範例 1：Dog 匿名嵌入 Animal
// ------------------------------------
// 這裡沒有給嵌入的 Animal 欄位取名字，稱為「匿名欄位」。
// 這會把 Animal 上的方法全部「提升」(promote) 到 Dog。
type Dog struct {
    Animal  // 匿名嵌入
    Breed   string
}

// 範例 2：Cat 命名欄位嵌入 Animal
// ------------------------------------
// 這裡把 Animal 欄位命名為 Ani，方法不會被提升到 Cat。
// 想用 Animal 的方法就必須寫 c.Ani.Speak()
type Cat struct {
    Ani   Animal  // 命名欄位
    Color string
}

func main() {
    // 建立 Dog 實例
    d := Dog{
        Animal: Animal{Name: "Buddy"},
        Breed:  "Beagle", // unused write to field Breed 是靜態分析工具（像是 go vet -unused、staticcheck）在提醒你：你在結構體 literal 裡給 Breed、Color 賦了值，卻從來沒有在程式的其他地方讀（使用）過這兩個欄位。
    }
    // 雖然 Speak 定義在 Animal 上，因為匿名嵌入，Dog 直接「提升」了它，
    // 所以可以直接呼叫 d.Speak()（底層會執行 d.Animal.Speak()）
    d.Animal.Speak() // 輸出：Buddy makes a sound
    d.Speak() // 輸出：Buddy makes a sound

    // 建立 Cat 實例
    c := Cat{
        Ani:   Animal{Name: "Kitty"},
        Color: "White",
    }
    // 下面這行編譯會失敗，因為 Cat 沒有被提升 Speak 方法
    // c.Speak()

    // 正確的呼叫方式，要透過命名欄位 Ani：
    c.Ani.Speak() // 輸出：Kitty makes a sound
}
