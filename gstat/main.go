package main

import (
	"fmt"
	"github.com/buetow/gstat/diskstats"
	"github.com/buetow/gstat/process"
	"time"
)

type twoP struct {
	first  process.Process
	second process.Process
}
type processMap map[string]twoP

func timedGather(timerChan <-chan bool, dRxChan chan<- diskstats.Diskstats, pRxChan chan<- process.Process) {
	for {
		switch <-timerChan {
		case false:
			{
				break
			}
		case true:
			{
				go diskstats.Gather(dRxChan)
				go process.Gather(pRxChan)
			}
		}
	}
	close(dRxChan)
	close(pRxChan)
}

func receiveD(dRxChan <-chan diskstats.Diskstats) {
	for d := range dRxChan {
		//diskstats.Print()
		// Implemented later
		_ = d
	}
}

func compareP(lastP processMap) {
	for _, val := range lastP {
		fmt.Printf("%d <-> %d\n", val.first.Pid, val.second.Pid)
	}
}

func receiveP(pRxChan <-chan process.Process) {
	lastP := make(processMap)
	flag := false

	for p := range pRxChan {
		if p.Last {
			if flag {
				compareP(lastP)
			}
			flag = !flag
		} else {
			if val, ok := lastP[p.Id]; ok {
				if flag {
					lastP[p.Id] = twoP{first: val.first, second: p}
				} else {
					lastP[p.Id] = twoP{first: p, second: val.second}
				}
			} else {
				lastP[p.Id] = twoP{first: p}
			}
		}
	}
}

func main() {
	timerChan := make(chan bool)
	dRxChan := make(chan diskstats.Diskstats)
	pRxChan := make(chan process.Process)

	go timedGather(timerChan, dRxChan, pRxChan)
	go receiveD(dRxChan)
	go receiveP(pRxChan)

	for counter := 0; counter < 5; counter++ {
		timerChan <- true
		time.Sleep(time.Second * 1)
	}
	timerChan <- false

	fmt.Println("Good bye")
}
