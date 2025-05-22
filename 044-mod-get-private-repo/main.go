package main

import (
	"fmt"
	privateRepo "github.com/andyrestart9/private-repo"
)

func main() {
	// 要使用 private 的 git repository 要先執行 export GOPRIVATE=github.com/andyrestart9/private-repo 命令
	// go get github.com/andyrestart9/private-repo
	// go mod tidy
	// 多個私有倉庫可以在同一個 GOPRIVATE 裡逗號分隔時，中間不要加空格。例如：export GOPRIVATE=git.company.com/*,github.com/your-org/*,gitlab.com/anotherteam/repo

	fmt.Println(privateRepo.Desc())
}
