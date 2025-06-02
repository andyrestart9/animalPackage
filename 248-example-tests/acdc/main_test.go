package acdc

import "fmt"

// ExampleSum 定义了一个针对 Sum 函数的示例。
// Go 测试框架会把所有名称以 Example 为前缀的函数当作“示例测试”来运行。
// 它既能出现在 godoc 生成的文档里，也能被 go test 用来校验输出是否符合预期。
func ExampleSum() {
	// 调用 acdc 包里的 Sum 函数，并把结果打印到标准输出
	fmt.Println(Sum(2, 3))
	// Output:
	// 5
	//
	// 上面的注释告诉 go test：运行 ExampleSum 时，
	// 捕获到的标准输出（println 的内容）必须**精确**匹配下面注释里的文本（这里是 “5”），
	// 否则该示例测试会被判定为失败。
}

// ExampleSum 定义了一个针对 Sum 函数的示例。
//
// 为什么要写 Example 函数？
// 1. 文档示例：godoc 会自动把这些示例摘出来，放在生成的包文档里，
//    让使用者看到函数的典型用法和输出。
// 2. 自动校验：像单元测试一样，go test 也会执行这些示例，
//    并对比它们打印的输出和注释里的 // Output: 部分，
//    保证示例代码随时可用且正确。
//
// 命名规则：
// - 函数必须以 Example 开头，可以紧跟被示例的函数名（这里是 Sum），
//   也可以额外加上后缀（例如 ExampleSum_withNegative）。
// - 签名无参数无返回值：func ExampleXxx() { … }。
//
// Output 注释格式：
// - 必须写在示例函数体的最后，用一行 // Output: 开头，
//   紧接着若干行 // 后跟实际期望的每一行输出。
// - 比对时包含换行和空格都要一致。
//
// 运行方式：
//   go test ./acdc
//   go test 会检测到 ExampleSum，执行它并验证输出；
//   也可以用 go test -v 来查看所有示例的执行情况。
//
// 这样既保证了示例代码的可用性，也让文档始终和实际行为保持同步。
