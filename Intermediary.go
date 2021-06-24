package main

import (
	"fmt"
	"time"
)

// Intermediary 仲介
type Intermediary struct {
	name                string               //名稱
	todoNormal          chan int             //正確的任務編號
	sosChan             chan string          // 當被殺時,最後一刻傳送求救訊號給另一個人
	emergencyNotifyChan chan *emergencyMsg   // 緊急任務頻道
	normalContractChan  chan *normalContract // 一般任務頻道
	killtopChan         chan *killtop        // 殺手殺掉上層的發送頻道
}

func NewIntermediary(name string, sos chan string, normalContract chan *normalContract, todoNormal chan int, emergencyNotifyCh chan *emergencyMsg, killtopCh chan *killtop) {
	i := &Intermediary{
		name:                name,
		sosChan:             sos,
		todoNormal:          todoNormal,
		normalContractChan:  normalContract,
		emergencyNotifyChan: emergencyNotifyCh,
		killtopChan:         killtopCh,
	}

	go i.running()

}

func (a *Intermediary) running() {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for {
		select {
		case number := <-a.todoNormal:

			for i := 1; i <= 30; i++ {
				no := &normalContract{
					mission_number: i,
					mission_real:   false,
				}
				if number == i {
					no.mission_real = true
				}
				a.normalContractChan <- no
			}

		case kill := <-a.killtopChan:
			if kill.dieName == "Intermediary" {
				fmt.Println(a.name, ":我死了,被", kill.name, "殺死了")
				fmt.Println(a.name, "發送支援給老闆")
				a.sosChan <- kill.name
				return
			} else {
				a.killtopChan <- kill
			}

		case name := <-a.sosChan:
			fmt.Println(a.name, ":老闆死了,被", name, "殺死了")
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
