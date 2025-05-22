package main // 定義主套件

import "fmt" // 匯入 fmt 套件以使用輸出功能

func main() {
	even := make(chan int)  // 建立一條無緩衝整數通道，用來傳送偶數
	odd := make(chan int)   // 建立一條無緩衝整數通道，用來傳送奇數
	quit := make(chan bool) // 建立一條無緩衝布林通道，用來傳送結束訊號

	go send(even, odd, quit) // 啟動一個背景 goroutine 執行 send，開始並行傳送資料

	receive(even, odd, quit) // 在 main goroutine 中同步呼叫 receive，等待並處理通道資料

	fmt.Println("aboutt to exit") // receive 返回後印出結束訊息，程式結束
}

// send 只允許發送資料到各通道
func send(even, odd chan<- int, quit chan<- bool) {
	for i := 0; i < 10; i++ { // 迴圈從 0 數到 9
		if i%2 == 0 { // 判斷是否為偶數
			even <- i // 將偶數 i 發送到 even 通道，無接收時會阻塞
		} else {
			odd <- i // 將奇數 i 發送到 odd 通道，無接收時會阻塞
		}
	}
	close(quit) // 關閉 quit 通道，通知接收方不會再有新訊號
}

// receive 只允許從各通道接收資料
func receive(even, odd <-chan int, quit <-chan bool) {
	for { // 無窮迴圈，直到在 quit case 裡 return
		select {
		case v := <-even: // 如果 even 通道有值就讀取
			fmt.Println("the value reveice from the even channel:", v) // 印出從 even 接收到的值
		case v := <-odd: // 如果 odd 通道有值就讀取
			fmt.Println("the value reveice from the odd channel:", v) // 印出從 odd 接收到的值
		case v, ok := <-quit: // 嘗試從 quit 通道接收，ok 表示通道是否仍開啟
			if !ok { // 如果 ok == false，代表通道已被關閉
				fmt.Println("from commoa ok", v, ok) // 印出零值及 ok 狀態(false)
				return                               // 結束 receive，跳出迴圈
			} else {
				fmt.Println("from comma ok:", v) // 如果 ok == true，印出接收到的訊號
			}
		}
	}
}
/*
在 Go 裡，channel（特別是無緩衝 channel）的送／收其實是一種「點對點的同步 rendezvous」，不是說 channel 本身被鎖住再打開，而是送與收的操作要互相配對才能各自完成。下面分步解釋流程：

一、無緩衝 (unbuffered) channel 的流程
假設有：
```
c := make(chan int)  // 無緩衝 channel
```
1. 送（send）先於收（receive）
Goroutine A 執行 c <- 42，發現沒有任何 goroutine 正在 <-c，
它就「阻塞」在這一行，並在 channel 內部排入「送方隊列」。
直到某個 Goroutine B 跑到 v := <-c：
B 看到有送方排隊，立即「配對」──將 42 傳給它，
解除 A 的阻塞（send 完成）並解除 B 的阻塞（receive 完成），
兩邊才同時繼續往下跑。

2. 收（receive）先於送（send）
Goroutine B 執行 v := <-c，發現沒有送方，
它就「阻塞」在這一行，並在 channel 內部排入「收方隊列」。
當 Goroutine A 執行 c <- 42：
A 看到有收方排隊，立即配對──把 42 傳給 B，
解除 A 的阻塞（send 完成）並解除 B 的阻塞（receive 完成）。
重點：不管誰先到，第一個到的人都會阻塞、等待對方；真正的送／收只有在「雙方都到站」時才同時完成。

二、帶緩衝 (buffered) channel 的流程
```
c := make(chan int, N)  // 緩衝大小 N
```
送：
如果 buffer 未滿（目前少於 N 個元素），c<-v 會立即成功，不阻塞；
如果 buffer 已滿（已有 N 個元素），c<-v 就會阻塞，等待有 <-c 釋出空位。

收：
如果 buffer 非空（至少有 1 個元素），v := <-c 會立即成功，不阻塞；
如果 buffer 為空，v := <-c 就會阻塞，等待有 c<-… 放入元素。
當一側阻塞了（因為 buffer 滿或空），另一側一旦到達且操作可完成，就會解除阻塞，兩者才繼續。

三、在 select 裡的「就緒」(ready) vs 「阻塞」(blocked)
```
select {
case c <- v:   // 只有在 send 不會阻塞（unbuffered 有收方在等，或 buffered 未滿）時才是就緒
    …
case v := <-c: // 只有在 receive 不會阻塞（unbuffered 有送方在等，或 buffered 非空）時才是就緒
    …
default:       // default 永遠就緒
    …
}
```
就緒：該 case 的 channel 操作可以「立即」完成，不用等對方。
阻塞：該 case 執行時會停下來等待對方配對，不會被視為就緒候選。
select 只會從所有就緒的 case 中隨機選一個執行；如果都阻塞且有 default，就跑 default；如果都阻塞且沒 default，整個 select 本身就阻塞。

小結
無緩衝 channel
送或收先到的一方會阻塞；第二方到達後，兩者配對，解除阻塞。
帶緩衝 channel
緩衝未滿／非空時一側不阻塞；滿／空時才阻塞。
select 的就緒判定
只有那些「操作不會阻塞」的 case 才是就緒，才會被選中。

這就是 Go channel 在送／收時的阻塞與就緒機制。理解了這個 rendezvous（會合）模型，就能更精準地設計並發程式、避免 deadlock 或意外的非同步行為。
*/

/*
在 Go 的執行模型裡，**阻塞（blocked）和就緒（ready 或 runnable）**是描述 goroutine（或在 select 裡面各個 case）在等待與執行上的兩種不同狀態：

1. Goroutine 的狀態
就緒（Ready 或 Runnable）
指這個 goroutine 已經準備好執行，並排在調度器（scheduler）的就緒隊列裡，等到系統分配到 CPU 時就會執行。
例如，一個 goroutine 從 channel 成功接收到資料、或一個 timer 到時、或剛從函式呼叫回來，都會被標記成就緒。

阻塞（Blocked 或 Waiting）
指這個 goroutine 因為某些條件而被暫停執行，不在就緒隊列裡，也不佔用 CPU。只有等到它等待的事件發生了（像 channel 有對應的對端、mutex 解鎖、timer 到期）才會被重新標記為就緒。

典型的阻塞情況：
ch <- v（發送到無緩衝 channel 或緩衝已滿的 channel）
v := <-ch（接收無緩衝 channel 或空的 channel）
sync.Mutex.Lock()（但 mutex 還沒被 Unlock）
time.Sleep(...) 期間也算阻塞直到時間到

2. select 裡的 case 就緒 vs 阻塞
在一個 select 語句中，每個 case 都代表一次 channel 操作（送或收）—調度器要先判斷哪些 case “不會阻塞”，才能選擇執行其中一個：
```
select {
case v := <-even:  // 只有當 even channel 裡**已有值**或已關閉時，這個 case 才是「就緒」
    // …
case odd <- x:    // 只有當 odd channel 緩衝未滿或已關閉（send on closed 會 panic 除外）時，這個 case 才是「就緒」
    // …
default:          // default 永遠「就緒」，因為它絕不會阻塞
}
```

就緒（ready）
該 case 的 channel 操作可以立即完成而不會停下來等待。例如：
讀取有資料的 buffered channel
向還沒滿的 buffered channel 發送
從已關閉的 channel 接收（永遠就緒，回傳零值＋ok=false）
select 只會從所有就緒的 case 中隨機挑一個執行。

阻塞（blocked）
該 case 的 channel 操作如果執行，會因為沒有對端或 channel 狀態不符合而必須「停下來等待」。例如：
從空的 unbuffered channel 或 empty buffered channel 接收
向滿的 buffered channel 發送
在判斷時，這些阻塞的 case 根本不列入候選，select 會忽略它們。

如果沒有任何 case 就緒且沒有 default
select 本身就會像一次 channel 操作那樣進入阻塞狀態，直到至少有一個 case 變成就緒。

就緒：能馬上執行、不用等對端；
阻塞：必須等對端配對才能繼續，否則就停在那裡等。
*/