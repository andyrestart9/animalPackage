package main

import "fmt"

type person struct {
	first string
}

type secretAgent struct {
	person
	ltk bool
}

func (p person) speak() {
	fmt.Println("I am", p.first)
}

func (sa secretAgent) speak() {
	fmt.Println("I am a secret agent", sa.first)
}

// person 實作了 human interface 所以 person 類型的值也是 human 類型
// secretAgent 實作了 human interface 所以 secretAgent 類型的值也是 human 類型
type human interface {
	speak()
}

// 傳入 human interface 就可以調用 human interface 的方法
func saySomething(h human) {
	h.speak()
}

func main() {
	sa1 := secretAgent{
		person: person{
			first: "James",
		},
		ltk: true,
	}

	p2 := person{
		first: "Jenny",
	}

	// sa1.speak()
	// p2.speak()

	saySomething(sa1)
	saySomething(p2)
}

// func (r receiver) identifier(p parameter(s)) (return(s)) { code }