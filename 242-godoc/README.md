# How to use godoc

## RUN Go Documentation Server

<https://pkg.go.dev/golang.org/x/tools/cmd/godoc>

如果沒有安裝 godoc 執行下面指令安裝

從 Go 1.12 開始，godoc 已被移出標準發行版，變為一個在 golang.org/x/tools 庫中的額外工具。你需要手動安裝它：

```sh
go install golang.org/x/tools/cmd/godoc@latest
```

安裝完之後，確認你已經把 $GOPATH/bin（或 $GOBIN）加到環境變數 PATH：

```sh
# 假設你的 GOPATH 是預設的 $HOME/go
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
source ~/.zshrc
```

如果能顯示路徑，就代表安裝並設定成功了。

```sh
which godoc
```

godoc 本身並沒有提供像 godoc help 這樣的子命令，但你可以用標準的 flag 方式來顯示使用說明：

```sh
godoc --help
```

run Go Documentation Server on local

```sh
godoc -http=localhost:6060
```

go to link and see Go Documentation Server

<http://localhost:6060/>

看我們在 242-godoc/mymath 寫的註解

<http://localhost:6060/pkg/mymodule/242-godoc/mymath/>

## 如何在 godoc.org 穰別人看到我的文件

go to <https://godoc.org> 搜尋 pkg 在 github 的 URL ，之後就能在 godoc.org 看到自己的 pkg 的文件

例如：

搜尋 `https://github.com/andyrestart9/puppy`

搜尋 `https://github.com/andyrestart9/animalPackage/tree/main/242-godoc/mymath`

不會馬上看到，收錄需要時間

## 為什麼 go doc 和 godoc 不會展示 package main

在 Go 里，go doc 和 godoc 默认只展示“导出”的符号（即首字母大写的类型、函数、变量等），而且 package main 是一个命令包，不能被别的包 import，所以从文档工具的角度它“没有任何导出标识符”，只会显示包级注释，不会列出函数或变量。

godoc 的 HTTP 服务同样遵循“只展示导出符号”的原则，且它专为 library（可被 import 的包）设计。要在 Web 界面下强行查看 main 包
