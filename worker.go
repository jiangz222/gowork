package main

import (
	"sync"
	"time"
)

type WorkerConfig struct {
	ConcurrencyNum int // max number of concurrency tasks
	TimeOut        int // mills seconds
}

func (wc WorkerConfig) Valid() bool {
	if wc.ConcurrencyNum == 0 {
		return false
	}
	return true
}

type Worker struct {
	config  WorkerConfig
	tasks   chan func()
	wg      WaitGroup
	mu      sync.Mutex
	running bool
	t       *time.Timer
	quit    chan struct{}
}

func New(config WorkerConfig) (instance *Worker) {
	if !config.Valid() {
		return nil
	}
	instance = &Worker{config: config, tasks: make(chan func(), 4*config.ConcurrencyNum), quit: make(chan struct{})}

	go instance.Run()
	if config.TimeOut > 0 {
		instance.t = time.NewTimer(time.Duration(instance.config.TimeOut) * time.Millisecond)
		go instance.timeout()
	}
	instance.mu.Lock()
	defer instance.mu.Unlock()
	instance.running = true
	return
}
func (w *Worker) timeout() {
	select {
	case <-w.t.C:
		w.t = nil
		w.exit()
	}
}
func (w *Worker) Add(f func()) {
	w.mu.Lock()
	if w.running == false {
		return
	}
	w.mu.Unlock()
	go func() {
		for {
			select {
			case w.tasks <- f:
				return
			case <-w.quit:
				return
			}
		}
	}()
}

func (w *Worker) Run() {
	for i := 0; i < w.config.ConcurrencyNum; i++ {
		go func(j int) {
			for {
				select {
				case task := <-w.tasks:
					if !w.running {
						return
					}
					w.wg.Add(1)
					task()
					w.wg.Done()

				case <-w.quit:
					return
				}
			}
		}(i)
	}
}

func (w *Worker) exit() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.running == true {
		w.running = false
		if w.t != nil {
			w.t.Stop()
		}
		close(w.quit)
		w.wg.DoneAll()
	}
}

func (w *Worker) Exit() {
	w.exit()
}

func (w *Worker) IsDone() bool {
	w.wg.Wait()
	return true
}
