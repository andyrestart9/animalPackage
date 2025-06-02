package mystr

import "strings"

func Cat(xs []string) string {
	s := xs[0]
	for _, v := range xs[1:] {
		s += " "
		s += v
	}
	return s
}

func Join(xs []string) string {
	return strings.Join(xs, " ")
}
/*
Cat VS strings.Join

从性能和内存开销来看，Join（基于 Builder 预分配）明显优于 Cat（用 += 循环拼接），原因有兩点：

时间复杂度

1. Cat 每次 s += " "+v 都是在背后重新分配、拷贝整个中间字符串，所以总体是 O(n²)（n 是切片长度），随着元素增多代价急剧上升。
每一次执行 s += v，Go 都会：
申请一块新内存：大小是 len(s) + len(v)，因为要放下之前的 s 和这次要加的 v。
把旧的 s 全部拷贝过去，再把 v 拷进去。

Join 先计算出总长度 n，一次性通过 b.Grow(n) 预分配缓冲区，然后用 WriteString 按序写入，整体是 O(n)，只遍历、拷贝一次。


2. 内存分配
Cat 在循环里每次增长都会触发新的分配和旧数据拷贝，分配次数 ≈ 元素个数。
Join 只调用一次 Grow，只分配一次大块内存，大大减少 GC 压力和分配开销。
*/