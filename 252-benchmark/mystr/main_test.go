package mystr

import (
	"fmt"     // 用于格式化输出，Example 函数中调用 fmt.Println 输出结果
	"strings" // 用于拆分字符串，测试和基准测试里使用 strings.Split
	"testing" // Go 内置测试框架，用于编写单元测试和基准测试
)

// TestCat 验证 Cat 函数能够将字符串切片按空格拼接成原始句子
func TestCat(t *testing.T) {
	s := "Shaken not stirred"      // 原始字符串
	xs := strings.Split(s, " ")    // 将字符串按空格拆分为切片
	s = Cat(xs)                    // 调用 Cat 函数拼接切片
	if s != "Shaken not stirred" { // 比较结果是否与原始字符串一致
		t.Error("got", s, "want", "Shaken not stirred")
	}
}

// TestJoin 验证 Join 函数能够将字符串切片按空格拼接成原始句子
func TestJoin(t *testing.T) {
	s := "Shaken not stirred"      // 原始字符串
	xs := strings.Split(s, " ")    // 将字符串按空格拆分为切片
	s = Join(xs)                   // 调用 Join 函数拼接切片（可能是对 strings.Join 的封装）
	if s != "Shaken not stirred" { // 比较结果是否与原始字符串一致
		t.Error("got", s, "want", "Shaken not stirred")
	}
}

// ExampleCat 用于文档测试，示例 Cat 函数的典型用法，并验证输出
func ExampleCat() {
	s := "Shaken not stirred"   // 原始字符串
	xs := strings.Split(s, " ") // 拆分为切片
	fmt.Println(Cat(xs))        // 打印拼接结果
	// Output:
	// Shaken not stirred
}

// ExampleJoin 用于文档测试，示例 Join 函数的典型用法，并验证输出
func ExampleJoin() {
	s := "Shaken not stirred"   // 原始字符串
	xs := strings.Split(s, " ") // 拆分为切片
	fmt.Println(Join(xs))       // 打印拼接结果
	// Output:
	// Shaken not stirred
}

// 用于基准测试的长文本常量，模拟真实场景下更大的输入规模
const s = "We ask ourselves, Who am I to be brilliant, gorgeous, talented, fabulous? Actually, who are you not to be? Your playing small does not serve the world. There is nothing enlightened about shrinking so that other people won't feel insecure around you. We are all meant to shine, as children do. We were born to make manifest the glory that is within us. It's not just in some of us; it's in everyone. And as we let our own light shine, we unconsciously give other people permission to do the same. As we are liberated from our own fear, our presence automatically liberates others. - Marianne Williamson"

var xs []string // 全局变量，避免在基准测试循环中重复拆分和分配内存

// BenchmarkCat 衡量 Cat 函数在大输入下的性能表现
func BenchmarkCat(b *testing.B) {
	// 1. 做好不需要计时的准备工作，例如生成输入数据
	xs = strings.Split(s, " ") // 在基准计时前拆分一次，减少额外开销
	/*
		在 Go 的基准测试（Benchmark）中，b.ResetTimer() 是 *testing.B 提供的一个方法，用来重置并重新启动内部计时器。它的主要作用是：

		清除之前执行的“非测量”阶段耗时
		在调用 BenchmarkXxx 函数时，通常会先做一些准备工作（比如分配输入数据、打开文件、建立连接等）。如果不重置计时器，那么这些准备工作所消耗的时间也会被算入最终的基准结果中。

		准确测量目标代码的执行时间
		在做完所有不想纳入测量的初始化之后，调用 b.ResetTimer()，它会将累积的耗时清零，然后从这一刻开始才计时。这样下面的循环里 b.N 次调用才是真正被衡量的部分。
	*/
	// 2. 重置计时器，忽略到目前为止的所有耗时
	b.ResetTimer() // 重置计时器：清除之前的耗时统计，确保下面循环中调用 Cat 的时间才被测量
	// 3. 基准循环：这里的每次迭代都会被计时
	for i := 0; i < b.N; i++ {
		Cat(xs) // 重复调用 Cat 函数
	}
}

// BenchmarkJoin 衡量 Join 函数在大输入下的性能表现
func BenchmarkJoin(b *testing.B) {
	xs = strings.Split(s, " ") // 在基准计时前拆分一次，减少额外开销
	b.ResetTimer()             // 重置计时器：清除之前的耗时统计，确保下面循环中调用 Join 的时间才被测量
	for i := 0; i < b.N; i++ {
		Join(xs) // 重复调用 Join 函数
	}
}
/*
基准测试后的报告
goos: darwin // 操作系统
goarch: arm64 // 架构
pkg: github.com/andyrestart9/animalPackage/252-benchmark/mystr // 包的导入路径
cpu: Apple M3 // CPU 型号
BenchmarkCat-8            142482              7692 ns/op
BenchmarkJoin-8          1964202               586.6 ns/op
PASS
ok      github.com/andyrestart9/animalPackage/252-benchmark/mystr       3.260s

BenchmarkCat-8 142482 7692 ns/op
BenchmarkCat：你在测试文件中定义的基准函数名。
-8：表示这次基准是在 GOMAXPROCS=8（也就是用了 8 个逻辑线程）下跑的。
142482：在规定的基准时间（默认 1 秒）里，总共执行了 142 482 次 BenchmarkCat 的循环。
7692 ns/op：平均每次调用耗时 7 692 纳秒（约 7.7 微秒）。
*/

/*
b.N 是什麼？

b 是传入的 *testing.B 对象，它负责控制和记录基准测试的细节。
b.N 是 testing.B 上的一个整型字段，表示这次基准测试里应该跑多少次循环。
在运行 go test -bench=. 时，测试框架会自动：
先猜一个初始的 b.N（比如 1、2、5、10……）。
运行循环并测时，看耗时是否足够。
如果不够，就增大 b.N 重测，直到耗时达到某个阈值（通常是几十到上百毫秒），保证测量结果有统计意义。
这样你在循环里写 for i := 0; i < b.N; i++，就不用关心具体要跑多少次，框架会帮你调整到一个合适的迭代次数，得出稳定的 “每次调用多少纳秒” 的结果。

看 Go 标准库中 testing 包的实现就能发现这段逻辑。以下摘自 src/testing/benchmark.go（Go 1.21+ 版本）中负责“定时模式”跑基准的核心部分
```
// b.loop.n == 0：表示还未做任何迭代，需要先进行“探测”以确定合适的循环次数。
// 如果不是固定迭代模式，就进入下面的 for 循环，不断增长 n，
// 直到累计耗时超过 benchTime（默认为 1s）为止。
d := b.benchTime.d
for n := int64(1); !b.failed && b.duration < d && n < 1e9; {
    last := n
    goalns := d.Nanoseconds()
    prevIters := int64(b.N)           // 上一次实际运行的迭代次数
    // predictN 根据目标时间、上次迭代次数和上次测量耗时，预测下一次迭代次数 n
    n = int64(predictN(goalns, prevIters, b.duration.Nanoseconds(), last))
    b.runN(int(n))                     // 设置 b.N = n 并跑 n 次基准函数
}

```
第一次迭代时，n 从 1 开始；
每跑完一次 b.runN(int(n))，内部会更新 b.N = n 并把这次循环的耗时累加到 b.duration；
只要 b.duration < benchTime（默认 1 秒），就继续预测更大的 n，再次跑基准；
最终跳出循环后，b.N 就是达成目标时间所需的迭代次数，用来计算 “ns/op”。
以上逻辑就体现了“Go 基准测试会自动调整 b.N ，直到测得的耗时达到至少 -benchtime 为止”这一行为。

---

n := int64(1); !b.failed && b.duration < d && n < 1e9; 是什麼？

这一行是 Go 中 for 循环的头部，分三部分：
```
for n := int64(1); !b.failed && b.duration < d && n < 1e9; {
    // ...
}
```
**初始化（n := int64(1)）** 将循环变量 n定义为int64类型并初始化为1`，表示第一次要跑 1 次基准。
条件（!b.failed && b.duration < d && n < 1e9）

循环会一直执行，直到这个布尔表达式为 false，具体包含三个子条件：
1. !b.failed：基准测试还未标记为失败 (b.failed == false)；
2. b.duration < d：当前累计测得的总耗时 b.duration（一个 time.Duration）仍然低于目标时长 d（通常是 benchTime，如 1 秒）；
3. n < 1e9：为了防止 n 过度膨胀（无限制增长到极大值），限定最多增长到 10^9 次。
只有当 三者都为真 时，循环才继续，否则退出。

后置语句
这里省略了第三段（通常写在 for 头部的 post 语句），意味着循环内部会自己更新 n（通过 predictN 和 b.runN），而不是用 for i++ 这种写法。

整体逻辑就是：
1. 从 n=1 次开始做一次基准测量；
2. 如果测得的总时间还不到目标，并且没有出错，且 n 未超过上限，就预测一个新的更大 n，再跑一次；
3. 重复上述过程，直到测时满足要求为止。如此便能自动找到一个合适的循环次数 n，让基准测试在足够长的时间窗口内进行，以获得稳定、可靠的 “ns/op” 指标。

---

d := b.benchTime.d 是什麼？

在这段基准测试的源码里：
```
// 在 B 结构体中，benchTime 存储了由 -test.benchtime 标志解析得来
// 的 durationOrCountFlag 类型，其中的 d 字段就是基准要跑的目标时长。
const (
    // 默认情况下，benchTime.d = 1 * time.Second
    benchTime = durationOrCountFlag{d: 1 * time.Second}
)

type B struct {
    // ...
    benchTime durationOrCountFlag // 来自 -test.benchtime 标志
    // ...
}

// 在自动调整循环次数的 for 里，第一步就是把这个时长取出来：
d := b.benchTime.d
for n := int64(1); !b.failed && b.duration < d && n < 1e9; {
    // ...
}
```
b.benchTime：是 B 结构体中的一个字段，类型为 durationOrCountFlag，它封装了用户通过 -test.benchtime 指定的运行时长（或固定次数）。
.d：这个 durationOrCountFlag 结构体的一个成员，类型是 time.Duration，保存了“要测”的总时长（默认 1 秒）。
d := b.benchTime.d：就是把这个目标 time.Duration 拷贝给局部变量 d，后续的循环条件里就用 b.duration < d 来判断「当前已耗时」是否还不足「目标时长」，从而决定是否继续增加迭代次数。

所以，d := b.benchTime.d 的作用，就是把用户期望的基准运行时长（如默认的 1 s 或者 -benchtime=500ms）取出来，用于控制基准测试自动增量循环的停止条件。

---

var (
	benchTime = durationOrCountFlag{d: 1 * time.Second} // changed during test of testing package
)
是什麼？

这行代码出自 testing 包内部，用来控制基准测试的默认运行时长：
```
var (
    benchTime = durationOrCountFlag{d: 1 * time.Second} // 改动仅在测试 testing 包自身时生效
)
```
逐部分解释：
1. var benchTime = …
声明了一个包级别变量 benchTime，所有基准测试都会读取它来知道「要跑多长时间」。
2. durationOrCountFlag{d: 1 * time.Second}
durationOrCountFlag 是一个既能表示“时长”（Duration）也能表示“固定迭代次数”（Count）的标志类型。这里初始化时只设置了它的 .d 字段为 1 * time.Second，也就是“默认要测 1 秒”。
3. // changed during test of testing package
这条注释说明：当 Go 在自测（self-test）testing 包本身的时候，会临时修改这个变量（比如改成更短或更长），以便测试基准框架的各个分支逻辑，而不影响普通用户通过 go test -bench 得到的默认行为。

所以，总结一下：
平时：benchTime.d == 1s，意味着每个 BenchmarkXxx 会自动扩大 b.N，直到累计跑够 1 秒为止，再报告结果。
测试 testing 包时：这一值会被修改，以覆盖各种边界场景（比如更短的时长、固定次数模式等），从而验证框架本身的正确性。
*/

/*
how to trace the variable "b.N" ?

cmd 點擊 b.N
type B struct {
    // ...
	N                int
	// ...
}

找到賦值給 b.N 的地方
func (b *B) runN(n int) {
	// ...
	b.N = n
	b.parallelism = 1
	b.ResetTimer()
	b.StartTimer()
	b.benchFunc(b)
	b.StopTimer()
	// ...
}
找到賦值給 b.N 的地方發現賦值給 b.N 後紧接着它用当前的 *B（也就是带着这个新 b.N）去调用 b.benchFunc(b)

b.benchFunc(b) 意味着 benchFunc 是 B 结构体的一个字段，类型是 func(*B)，它在这里被调用。要弄清这到底调用的是哪个函数实现，就必须找到“给这个字段赋值”的地方。
func (b *B) Run(name string, f func(b *B)) bool {
    // ..
	sub := &B{
		common: common{
			signal:  make(chan bool),
			name:    benchName,
			parent:  &b.common,
			level:   b.level + 1,
			creator: pc[:n],
			w:       b.w,
			chatty:  b.chatty,
			bench:   true,
		},
		importPath: b.importPath,
		benchFunc:  f, // benchFunc 在這裡被賦值
		benchTime:  b.benchTime,
		context:    b.context,
	}
    // ...
	if sub.run1() {
		sub.run() // &B 賦值給 sub ，後在這裡發現 sub.run 被調用
	}
	// ...
}
&B 賦值給 sub ，後在這裡發現 sub.run 被調用

追踪 sub.run
func (b *B) run() {
	// ...
		b.doBench()
	// ...
}

追踪 b.doBench
func (b *B) doBench() BenchmarkResult {
	go b.launch()
	// ...
}

追踪 b.launch
func (b *B) launch() {
	// ...
	b.runN(int(n))
	// ...
}

串通整个调用链
1. b.Run(name, f)
在 Run 方法里把 f（你的 BenchmarkXxx）赋给 sub.benchFunc，然后调用 sub.run1() / sub.run().
2. sub.run()
根据上下文分支到 doBench()（无子基准）或 processBench()（有子基准），但最终都要调用 doBench().
3. doBench()
Fork 一个 goroutine 去执行 launch()，并在当前协程 ← b.signal 阻塞。
4. launch()
如果用户用 -benchtime=100x，直接 runN(100).
否则按时间扩容：循环预测新的 n，不断调用 runN(n)，直到累积时长 ≥ 目标（默认为 1s）。
5. runN(n)
把 n 写入 b.N，然后调用 b.benchFunc(b) —— 也就是你的 BenchmarkXxx(b)，内部的 for i := 0; i < b.N; i++ { … } 就在这里执行。

这样，你就完整地从 benchFunc: f → sub.run1()/sub.run() → doBench() → launch() → runN(n) → b.benchFunc(b) 把 b.N 赋值和消费的全链路理清了。
*/
