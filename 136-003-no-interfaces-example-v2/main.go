package main

import "fmt"

type Human struct{ Name string }
func (h Human) Speak() string { return "Hi, I'm " + h.Name }

type Robot struct{ ID int }
func (r Robot) Speak() string { return fmt.Sprintf("Beep! I am robot #%d", r.ID) }

func Announce(v any) {           // 接受空介面
	switch s := v.(type) {       // 逐一列舉型別
	case Human:
		fmt.Println(s.Speak())
	case Robot:
		fmt.Println(s.Speak())
	default:
		panic("unsupported type")
	}
}

func main() {
	Announce(Human{"Andy"})
	Announce(Robot{42})
}