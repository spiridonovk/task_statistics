package engine

import (
	"sync"
	"time"
)

const timeDelay = 60

var allClients *clients

type clients struct {
	mu    *sync.Mutex
	chans []chan struct{}
}

type Statistics struct {
	Mu         *sync.RWMutex
	Statistics map[string]int
}

var Stat, tempStat Statistics

func Init() {
	tempStat = newStat()
	allClients = newClients()
	Stat = newStat()
}
func newStat() Statistics {
	return Statistics{
		Mu:         &sync.RWMutex{},
		Statistics: make(map[string]int),
	}
}
func newClients() *clients {
	return &clients{
		mu:    &sync.Mutex{},
		chans: make([]chan struct{}, 0),
	}
}
func AddToStat(cmd string) {
	Stat.Mu.Lock()
	tempStat.Statistics[cmd] += 1
	Stat.Mu.Unlock()
}

func AddNewClient(ch chan struct{}) {
	allClients.mu.Lock()
	allClients.chans = append(allClients.chans, ch)
	allClients.mu.Unlock()
}
func sendToClients() {
	allClients.mu.Lock()
	for _, ch := range allClients.chans {
		ch <- struct{}{}
	}
	allClients.mu.Unlock()
}

func Timer() {
	ticker := time.NewTicker(timeDelay * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				Stat = tempStat
				tempStat = newStat()
				sendToClients()
			}
		}
	}()
}
