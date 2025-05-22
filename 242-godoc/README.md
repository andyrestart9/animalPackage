# How to use godoc

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
