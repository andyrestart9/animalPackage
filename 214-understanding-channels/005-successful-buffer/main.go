package main

import "fmt"

func main() {
	c := make(chan int, 2)

	c <- 42
	c <- 43

	fmt.Println(<-c)
	fmt.Println(<-c)
}
/*
這段程式碼是正確的，不會產生死鎖。解釋如下：

使用 make(chan int, 2) 建立了一個緩衝區大小為 2 的 channel，這意味著最多可以暫存 2 個值而不會阻塞。

當執行 c <- 42 和 c <- 43 時，由於 channel 內部緩衝區可以容納這兩個值，所以這兩個 send 操作都不會阻塞。

最後，fmt.Println(<-c) 分別從 channel 中取出兩個值（會依照先進先出 FIFO 的順序，也就是先取出 42，再取出 43），程式順利完成。

因此，這段程式碼既不會阻塞也不會產生死鎖。
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