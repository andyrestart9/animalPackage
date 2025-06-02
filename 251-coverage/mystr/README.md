# Coverage

## Command

```sh
go test -cover
go test -coverprofile c.out
go tool cover -html=c.out
```

## 如何找到 go tool cover -html=c.out

```sh
go help
#         tool        run specified go tool
# Use "go help <command>" for more information about a command.
#         testflag        testing flags
# Use "go help <topic>" for more information about that topic.
# 用 go help 找到 tool 和 testflag 的用法

go help testflag
# -coverprofile cover.out
    # Write a coverage profile to the file after all tests have passed.

go help tool
# With no arguments it prints the list of known tools.

go tool
# cover
# 發現有 cover 工具可以用

go tool cover -help
# Usage of 'go tool cover':
# Given a coverage profile produced by 'go test':
        # go test -coverprofile=c.out
# Open a web browser displaying annotated source code:
        # go tool cover -html=c.out
# 發現 cover 工具的用法

go tool cover -html=c.out
```
