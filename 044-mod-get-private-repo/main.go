package main

import (
	"fmt"
	privateRepo "github.com/andyrestart9/private-repo"
)

func main() {
	// 要使用 private 的 git repository 要先執行 export GOPRIVATE=github.com/andyrestart9/private-repo 命令
	fmt.Println(privateRepo.Desc())
}
