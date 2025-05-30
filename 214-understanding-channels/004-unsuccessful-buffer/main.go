package main

import "fmt"

func main() {
	c := make(chan int, 1)

	c <- 42
	c <- 43

	fmt.Println(<-c)
}
/*
這段程式碼會產生死鎖。原因如下：

你建立了一個帶有緩衝區大小為 1 的 channel：c := make(chan int, 1)。

當你執行 c <- 42 時，42 成功被存入 channel，因為緩衝區是空的。

接著執行 c <- 43 時，由於緩衝區已滿（只有 1 個位置），這個 send 操作會阻塞，等待有空位。

由於目前仍在同一個 goroutine（main）中執行，且沒有同時執行的接收操作，程式就卡在這裡，導致死鎖。
*/

/*
盡量使用 unbuffered channels 而非 buffered channels，原因有以下幾點：

嚴格同步的保證：
使用 unbuffered channels 意味著每一次的送出操作必須等待對應的接收操作，這樣能夠確保程式中資料的傳遞是同步進行的。換句話說，資料從送出到被接收之間形成了一個「交織」的連結，這樣可以降低出錯的風險，並使邏輯更容易推理。

明確的控制流：
當你使用 unbuffered channels 時，程式碼必須明確地安排好誰在送資料、誰在收資料，這種「鎖定」機制迫使你在設計時思考並確保所有操作都有配對，這有助於發現潛在的同步問題。

避免潛在隱患：
雖然 buffered channels 可以讓生產者和消費者解耦，從而在某些情況下提升效能，但它們也容易讓你忽略或掩蓋程式中不恰當的同步邏輯。用 unbuffered channels 則可以讓你及早發現設計上的漏洞，因為任何未能及時進行的接收操作都會立即導致阻塞，從而迫使你修正問題。

總結來說：
盡量避免依賴 buffered channels，除非確有必要。優先使用 unbuffered channels，讓你的程式在送出與接收之間始終存在一個明確且同步的交接點，這樣可以使程式的同步行為更直觀，降低錯誤發生的可能。
*/