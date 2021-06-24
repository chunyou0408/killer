package main

import (
	"fmt"
	"math/rand"
	"time"
)

// emergencyMsg 緊急通知
type emergencyMsg struct {
	name  string // 通緝的殺手名稱
	money int    // 獎金
}

// normalContract 一般訊息
type normalContract struct {
	mission_number int  // 任務編號
	mission_real   bool // 任務真假
}

// Boss 幕後老闆
type Boss struct {
	name                string             // 名稱
	todoNormal          chan int           // 正確的任務編號
	sosChan             chan string        // 當被殺時,最後一刻傳送求救訊號給另一個人
	emergencyNotifyChan chan *emergencyMsg // 緊急任務頻道
	killtopChan         chan *killtop      // 殺手殺掉上層的發送頻道
}

func NewBoss(name string, sos chan string, todoNormal chan int, emergencyNotifyCh chan *emergencyMsg, killtopCh chan *killtop) {
	b := &Boss{
		name:                name,
		sosChan:             sos,
		todoNormal:          todoNormal,
		emergencyNotifyChan: emergencyNotifyCh,
		killtopChan:         killtopCh,
	}

	go b.running()

}

func (a *Boss) running() {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t := time.Now().UnixNano()
			r1 := rand.New(rand.NewSource(t))
			number := r1.Int()%30 + 1
			a.todoNormal <- number
		case kill := <-a.killtopChan:
			if kill.dieName == "Boss" {
				fmt.Println(a.name, ":我死了,被", kill.name, "殺死了")
				fmt.Println(a.name, "發送支援給仲介")
				a.sosChan <- kill.name
				return
			} else {
				a.killtopChan <- kill
			}

		case name := <-a.sosChan:
			fmt.Println(a.name, ":仲介死了,被", name, "殺死了")
			fmt.Println(a.name, "發布警急任務,殺掉", name)

			em := &emergencyMsg{
				name:  name,
				money: 10000000,
			}
			a.emergencyNotifyChan <- em

			return
		}
	}
}
