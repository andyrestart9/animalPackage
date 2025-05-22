package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.Open("xx.txt")
	if err != nil {
		log.Println("err happened:", err) // 2025/05/21 19:38:36 err happened: open xx.txt: no such file or directory
		return
	}
	defer f.Close()
}
