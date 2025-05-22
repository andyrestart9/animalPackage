package main

import "fmt"

func main() {
	eve := make(chan int)
	odd := make(chan int)
	quit := make(chan int)

	// 送出端放到背景 goroutine
	// 發送前必須有對應的接收，接收前必須有對應的發送，兩者必須配對才能解除阻塞。
	// go send：背景開始送資料，但第一筆送出後會因無接收而阻塞，不會結束程式
	go send(eve, odd, quit)

	// 接收端在 main 中同步呼叫，直到收到 quit 才結束
	// receive：在 main 中同步呼叫，select 會一邊收一邊把背景 goroutine 的阻塞解除，直到讀到 quit 才跳出
	receive(eve, odd, quit)

	// 送（producer） 與 收（consumer） 必須在不同 goroutine ⇒ 先啟動 go send，再同步呼叫 receive，才能互相配對、避免死鎖。
	// main 必須有阻塞點 才能等到真正做完所有工作後再離開。
	fmt.Println("about to exit")
}

// 只允許從這三條 channel 接收
func receive(e, o, q <-chan int) {
	for {
		// receive 裡的 select 會同時監聽三條 channel，依照收到的訊息路徑分類處理，最後遇到 quit 就 return，結束迴圈
		select {
		case v := <-e:
			fmt.Println("from the eve channel:", v)
		case v := <-o:
			fmt.Println("from the odd channel:", v)
		case v := <-q:
			fmt.Println("from the quit channel:", v)
			// return 會立刻結束整個函式的執行，並把控制權交還給呼叫端，它不只是跳出迴圈。如果你在迴圈裡面寫了 return，那個函式後面的任何程式碼（包括迴圈之外的）都不會再執行。
			return
		}
	}

}

// 只允許往這三條 channel 發送
func send(e, o, q chan<- int) {
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			e <- i
		} else {
			o <- i
		}
	}
	q <- 0
}
/*
為什麼要把 send 放進 goroutine，但 receive 卻不放？
無緩衝 channel 的特性：
發送前必須有對應的接收，接收前必須有對應的發送，兩者必須配對才能解除阻塞。

如果 send 同步執行：
send(eve, odd, quit)
receive(eve, odd, quit)
send 在 main 內同步呼叫時，一開始就嘗試推送 0 到 eve，因為 receive 還沒開始而阻塞 → 死鎖。

把 receive 也丟 goroutine：
go send(...)
go receive(...)
fmt.Println("about to exit")
main goroutine 沒有阻塞點，立刻印出 "about to exit" 然後程式結束，背景的 send/receive goroutine 會被強制殺掉 → 資料不一定完整被處理。

正確做法：
go send(...)
receive(...)      // main 阻塞在這裡，等到收到 quit 才回來
fmt.Println("about to exit")
go send：背景開始送資料，但第一筆送出後會因無接收而阻塞，不會結束程式

receive：在 main 中同步呼叫，select 會一邊收一邊把背景 goroutine 的阻塞解除，直到讀到 quit 才跳出

最後印出退出訊息並結束程式

---

送（producer） 與 收（consumer） 必須在不同 goroutine ⇒ 先啟動 go send，再同步呼叫 receive，才能互相配對、避免死鎖。

main 必須有阻塞點 才能等到真正做完所有工作後再離開。

directional channels（chan<- / <-chan）＋ select，能讓發送端與接收端在型別上各司其職、程式結構更清晰。

千萬不要依賴 -race 的副作用或看似能「偶然成功」的非決定性行為，務必使用同步原語（像這裡的同步呼叫或 sync.WaitGroup）確保程式正確、穩定地執行完畢。
*/
/*
用
// send
go send(eve, odd, quit)

// receive
go receive(eve, odd, quit)
看起來似乎沒有問題

它「看起來」沒問題，多半是因為你是用 -race 來跑，race 檢測器會大幅放慢排程、延遲退出，讓那兩隻 goroutine 有機會把資料全印出來。但這並不是正確的做法，也不是保證可用的行為——它完全仰賴排程時機。

為什麼 -race 下看起來沒問題？
調度延遲：race 檢測器會在各 goroutine 間插入更多同步、更多檢查，整體跑得更慢、更「公平」，因而 often send/receive 有時間跑完才 exit。

輸出緩衝：fmt 的輸出也可能因排程時機被完整 flush。

但這都只是副作用，當你移除 -race，或是在快排程環境下執行，主程式隨時都可能先回去而把子 goroutine 一起殺掉——就看你踩到哪個排程切換點，行為是非決定性的。
*/

/*
go send(eve, odd, quit)
receive(eve, odd, quit)
是先執行 send 還是先執行 receive？

當你寫：
go send(eve, odd, quit)
receive(eve, odd, quit)
go send(...) 只是「啟動一個新的 goroutine 去執行 send」然後馬上返回。

緊接著，main goroutine 就會立刻呼叫 receive(...)。

也就是說，在語句層面，receive 是緊接在 go send 之後同步執行的。至於真正 send 函式裡面的程式碼（第一個 e <- 0）和 receive 裡面的程式碼（第一個 <-e）哪個先跑，就完全交給 Go 的排程器了──它們是並行的，沒有硬性先後順序。

如果排程器先跑到 main，receive 會立刻進入 select、在 <-e 那裡阻塞，等待 send 發值。

如果排程器先跑到 send，send 在 e <- 0 那裡阻塞，直到 receive 開始讀。

總結：

語句執行順序：go send(...) → receive(...)

函式內部真正的送／收執行時機：取決於排程器，send 和 receive 是並發跑的，沒有絕對誰先誰後。


---

send 和 receive 是並發跑的，那還沒有 send 第一個 e <- 0 就先 receive 裡面的程式碼第一個 <-e 不會有問題嗎？


不會有問題。對無緩衝 channel 而言，<-e 和 e<- 都會阻塞直到對方到位，這正是它的「同步 rendezvous（會合）」機制：
receive 先執行 <-e
主 goroutine 在 receive 裡的 case v := <-e 會先到，就在這裡阻塞、等著有人送資料。
因為它是同步阻塞，所以並不會跳過，也不會 panic，只是暫時停在那一行。
send goroutine 執行 e<-0
背景的 send goroutine 一開始跑到 if i%2==0，執行 e<-i。
它也會因為 e<-0 找不到接收者而阻塞，直到另一端（也就是剛才那個在 <-e 的 receive）還在等。
兩邊同時到阻塞點，就會解開阻塞
一旦 send 嘗試寫，receive 正好在等，那麼兩者就配對成功：
send 的阻塞被解除──資料寫入完成
receive 的阻塞被解除──讀到剛寫入的值
這樣的雙向同步保證了「無論誰先到，都會等到對方」，所以不會出現「先讀不到就 crash」或「先寫沒人讀就死鎖」的情況（只會單方面阻塞，直到配對）。
*/

/*
沒有先後順序那寫成
	receive(eve, odd, quit)
	go send(eve, odd, quit)
可以嗎？


不行，這樣會馬上死鎖。原因在於：

// 1. 直接同步呼叫 receive —> 這裡就阻塞在 select 等消息
receive(eve, odd, quit)

// 2. send 還沒被啟動，go send(...) 完全沒機會執行
go send(eve, odd, quit)
receive 會一直在 select 的那幾個 <-chan 上阻塞，直到有第一筆 eve、odd 或 quit 的發送進來；

但你的 send goroutine 是在 receive 完成後才啟動──可惜 receive 永遠都不會自己 return（quit 值也是由 send 最後才送出的），所以主程式卡在那裡、永遠等不到，send 也不會被呼叫，整組程式就死鎖了。

要讓它正確執行，send 一定要先（或同時）啟動，才能讓後續的 receive 有資料可讀。常見寫法是：
go send(eve, odd, quit)  // 先啟動 producer
receive(eve, odd, quit)  // main 在這裡阻塞，等到 quit 再 return
如此 send 才能「先到先阻塞」，receive 才能「先阻塞再被喚醒」，兩邊才能正常配對、程式也才能最終跳出並印出 "about to exit"。

---
那為什麼你說他們是併發，沒有先後順序？

這裡有兩層「先後順序」要分清楚：

1. 啟動 goroutine 的先後順序
go send(...) 先啟動 send goroutine，再呼叫 receive(...)，才是真正把「送」跟「收」放到兩個獨立的執行緒（goroutine）裡面，這時候它們才真正並發運作。

如果你把 receive(...) 放在前面，主 goroutine 就在那裡同步阻塞，連 go send 都還沒執行到，根本不會有第二支 goroutine；此時並沒有並發，只有一支 goroutine 在 receive 裡面無限等待。

所以，要達到真 的併發，你必須先用 go send(...) 把 send 放到背景，讓 main goroutine 能繼續跑到 receive(...)；這兩支 goroutine 才能真正「同時」在 channel 上做送／收。

2. 真正送與收在 channel 上的同步「會合」
一旦兩支 goroutine 都在 channel 操作時，兩邊對同一條無緩衝 channel 的 e<- 和 <-e 真正執行，就沒有固定的先後──

如果送的 goroutine 比較早跑到 e<-0，它就先阻塞，等接收那端到位。

如果接收的 goroutine 比較早跑到 <-e，它就先阻塞，等送那端到位。
當兩端「都到站」了，操作才同時完成，這才是 channel 無緩衝同步的本質。

總結
啟動順序（你要先 go send，才呼叫 receive）決定了是不是有兩支真正獨立的 goroutine。

操作會合（e<- vs. <-e）一旦在兩支 goroutine 裡，就是真正並行、沒有固定誰先誰後，只要有一端到站就會各自阻塞，直到另一端也到站。
*/