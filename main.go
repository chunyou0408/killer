package main

import "fmt"

var exit chan string
var killtopCh chan *killtop
var sos chan string
var normalContractCh chan *normalContract
var emergencyNotifyCh chan *emergencyMsg
var todoNormal chan int

func main() {
	initSetting()
	fmt.Println("結束", <-exit)
}

func initSetting() {
	// 全部結束的頻道
	exit = make(chan string)
	// 用來呼叫另一方啟動警急任務
	sos = make(chan string)
	// 老闆傳給仲介任務資料
	todoNormal = make(chan int)
	//
	killtopCh = make(chan *killtop)
	// 一般任務頻道
	normalContractCh = make(chan *normalContract)
	// 警急任務頻道
	emergencyNotifyCh = make(chan *emergencyMsg)

	NewBoss("Boss", sos, todoNormal, emergencyNotifyCh, killtopCh)
	NewIntermediary("Intermediary", sos, normalContractCh, todoNormal, emergencyNotifyCh, killtopCh)
	NewAssassin("JohnWick", normalContractCh, emergencyNotifyCh, killtopCh)
	NewAssassin("JasonBourne", normalContractCh, emergencyNotifyCh, killtopCh)
	NewAssassin("EthanHunt", normalContractCh, emergencyNotifyCh, killtopCh)
}
