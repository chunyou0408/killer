package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Assassin 殺手
type Assassin struct {
	name                  string
	emergencyNotifyChan   chan *emergencyMsg
	normalContractChan    chan *normalContract
	normalmissioncomplete int
	killtopCh             chan *killtop
}
type killtop struct {
	name    string
	dieName string
}

// KillHighLevelPerson 幹掉上層
func (a *Assassin) KillHighLevelPerson() {
	killtopCh = make(chan *killtop)
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
	a.killtopCh <- kill
}

// NewAssassin
func NewAssassin(name string, normalContract chan *normalContract, emergencyNotifyCh chan *emergencyMsg, killtopCh chan *killtop) {

	a := &Assassin{
		name:                  name,
		normalmissioncomplete: 0,
		normalContractChan:    normalContract,
		emergencyNotifyChan:   emergencyNotifyCh,
		killtopCh:             killtopCh,
	}

	go a.running()
}

func (a *Assassin) running() {

	for {
		select {
		case no := <-a.normalContractChan:
			// fmt.Println(a.name, ":收到任務,編號", no.mission_number, "判斷真偽中...")
			time.Sleep(time.Millisecond * 1)
			if no.mission_real {
				a.normalmissioncomplete++
				fmt.Println(a.name, ":任務",no.mission_number,"為真,執行任務 目前任務次數:", a.normalmissioncomplete)
				if a.normalmissioncomplete >= 20 {
					fmt.Println(a.name, "任務都做完了,殺掉上層")
					a.KillHighLevelPerson()
					return
				}
				time.Sleep(time.Millisecond * 10)
			}

		case em := <-a.emergencyNotifyChan:
			fmt.Println(a.name, ":收到緊急任務,殺掉", em.name)
			fmt.Println(a.name, ":擊殺成功,取得獎金", em.money)
			exit <- "GG"
		}
	}
}