package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Assassin 殺手
type Assassin struct {
	name                  string               // 名稱
	emergencyNotifyChan   chan *emergencyMsg   // 緊急任務頻道
	normalContractChan    chan *normalContract // 一般任務頻道
	normalmissioncomplete int                  // 任務完成數目
	killtopChan           chan *killtop        // 殺手殺掉上層的發送頻道
}
type killtop struct {
	name    string // 殺手名稱
	dieName string // 要殺的上層名稱
}

// KillHighLevelPerson 幹掉上層
func (a *Assassin) KillHighLevelPerson() {
	kill := &killtop{
		name: a.name,
	}
	// 指定要殺掉誰

	t := time.Now().UnixNano()
	r1 := rand.New(rand.NewSource(t))
	number := r1.Int() % 2
	if number == 1 {
		kill.dieName = "Boss"
	} else {
		kill.dieName = "Intermediary"
	}

	fmt.Println(a.name, "殺掉", kill.dieName)
	a.killtopChan <- kill
}

// NewAssassin
func NewAssassin(name string, normalContract chan *normalContract, emergencyNotifyCh chan *emergencyMsg, killtopCh chan *killtop) {

	a := &Assassin{
		name:                  name,
		normalmissioncomplete: 0,
		normalContractChan:    normalContract,
		emergencyNotifyChan:   emergencyNotifyCh,
		killtopChan:           killtopCh,
	}

	go a.running()
}

func (a *Assassin) running() {

	for {
		select {
		case no := <-a.normalContractChan:
			fmt.Println(a.name, ":收到任務,編號", no.mission_number, "判斷真偽中...")
			time.Sleep(time.Millisecond * 1)
			if no.mission_real {
				a.normalmissioncomplete++
				fmt.Println(a.name, ":任務", no.mission_number, "為真,執行任務 目前任務次數:", a.normalmissioncomplete)
				if a.normalmissioncomplete >= 20 {
					fmt.Println(a.name, "任務都做完了,殺掉上層")
					a.KillHighLevelPerson()
					return
				}
				time.Sleep(time.Millisecond * 10)
			}

		case em := <-a.emergencyNotifyChan:
			fmt.Println(a.name, ":收到緊急任務,殺掉", em.name)
			fmt.Println(a.name, ":擊殺成功,取得獎金", em.money, "美金")
			exit <- "GG"
			return
		}
	}
}
