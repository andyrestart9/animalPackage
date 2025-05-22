
package main

import (
    "fmt"
    "os"
)

func main() {
    // 创建或打开文件
    file, err := os.Create("example.txt")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close() // 确保在函数结束时关闭文件

    // 写入内容到文件
    content := []byte("Hello, Go file writing!")
    _, err = file.Write(content)
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }

    fmt.Println("File written successfully!")
}
