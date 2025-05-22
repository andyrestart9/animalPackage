package main

import "fmt"

// 兩種不同的 struct
type Person struct {
    Name string
    Age  int
}

type Book struct {
    Title  string
    Author string
}

func printPerson(p Person) {
    fmt.Println("Person.Name:", p.Name)
    fmt.Println("Person.Age: ", p.Age)
}

func printBook(b Book) {
    fmt.Println("Book.Title: ", b.Title)
    fmt.Println("Book.Author:", b.Author)
}

func main() {
    p := Person{"Alice", 30}
    b := Book{"1984", "Orwell"}

    // 必須分別呼叫不同的函式
    printPerson(p)
    fmt.Println("-----")
    printBook(b)
}
