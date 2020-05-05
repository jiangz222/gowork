package main

import (
	"sync"
)

type WaitGroup struct {
	waitGroup  sync.WaitGroup
	addNumbers int
	mu         sync.Mutex
}

func (wg *WaitGroup) Add(delta int) {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	wg.addNumbers += delta
	wg.waitGroup.Add(delta)
}

// Done 重复Done不会panic，和原生的不同
func (wg *WaitGroup) Done() {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	if wg.addNumbers == 0 {
		return
	}
	wg.addNumbers -= 1
	wg.waitGroup.Done()
}

func (wg *WaitGroup) Wait() {
	wg.waitGroup.Wait()
}

func (wg *WaitGroup) Len() int {
	return wg.addNumbers
}

// DoneAll 一次把剩下的所有等待都结束
func (wg *WaitGroup) DoneAll() {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	for i := 0; i < wg.addNumbers; i++ {
		wg.waitGroup.Done()
	}
	wg.addNumbers = 0

}
