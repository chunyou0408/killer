package main

import (
	"fmt"
	"math/rand"
	"time"
)

// emergencyMsg 緊急通知
type emergencyMsg struct {
	name  string
	money int
}

// normalContract 一般訊息
type normalContract struct {
	mission_number int
	mission_real   bool
}

// Boss 幕後老闆
type Boss struct {
	name                string
	todoNormal          chan int
	sosChan             chan string
	emergencyNotifyChan chan *emergencyMsg
	killtopCh           chan *killtop
}

func NewBoss(name string, sos chan string, todoNormal chan int, emergencyNotifyCh chan *emergencyMsg, killtopCh chan *killtop) {
	b := &Boss{
		name:                name,
		sosChan:             sos,
		todoNormal:          todoNormal,
		emergencyNotifyChan: emergencyNotifyCh,
		killtopCh:           killtopCh,
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
		case kill := <-a.killtopCh:
			if kill.dieName == "Boss" {
				fmt.Println(a.name, ":我死了,被", kill.name, "殺死了")
				fmt.Println(a.name, "發送支援給仲介")
				a.sosChan <- kill.name
				return
			} else {
				a.killtopCh <- kill
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
