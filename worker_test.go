package main

import (
	"fmt"
	"testing"
	"time"
)

func TestWorker_Normal(t *testing.T) {
	worker := New(WorkerConfig{
		ConcurrencyNum: 10,
	})
	for i := 0; i < 100; i++ {
		j := i
		worker.Add(func() {
			fmt.Println("done ", j)
		})
	}
	worker.IsDone()
}
func TestWorker_ConcurrencyNum(t *testing.T) {
	worker := New(WorkerConfig{
		ConcurrencyNum: 10,
	})
	start := time.Now().Unix()
	for i := 0; i < 20; i++ {
		j := i
		worker.Add(func() {
			// make the producer is more faster than consumer, test the concurrencyNum works fine
			time.Sleep(1 * time.Second)
			fmt.Println("done ", j)
		})
	}
	worker.IsDone()
	end := time.Now().Unix()
	if end-start > 2 {
		fmt.Println(end - start)
		t.Error("concurrency works not correctly")
	}
	time.Sleep(1 * time.Second)

}

func TestWorker_Timer(t *testing.T) {
	worker := New(WorkerConfig{
		ConcurrencyNum: 10,
		TimeOut:        2000,
	})
	start := time.Now().Unix()
	for i := 0; i < 30; i++ {
		j := i
		worker.Add(func() {
			time.Sleep(1 * time.Second)
			fmt.Println("done ", j)
		})
		fmt.Println("add ", i)
	}
	worker.IsDone()
	end := time.Now().Unix()
	if end-start > 2 {
		fmt.Println(end - start)
		t.Error("timer works not correctly")
	}
}
func TestWorker_Exit(t *testing.T) {
	worker := New(WorkerConfig{
		ConcurrencyNum: 10,
	})
	start := time.Now().Unix()

	for i := 0; i < 20; i++ {
		j := i
		worker.Add(func() {
			time.Sleep(1 * time.Second)
			fmt.Println("done ", j)
		})
		fmt.Println("add ", i)
	}
	worker.Exit()
	//worker.IsDone() // IsDone is not necessary when Exit() is called
	end := time.Now().Unix()
	if end-start > 1 {
		fmt.Println(end - start)
		t.Error("exit works not correctly")
	}

}
