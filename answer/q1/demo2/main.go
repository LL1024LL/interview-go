package main

import (
	"fmt"
	"sync"
	"time"
)

func main()  {
	m1 := &sync.Mutex{}
	cond1 := sync.Cond{L:m1 }
	m2 := &sync.Mutex{}
	cond2 := sync.Cond{L: m2}
	go func() {
		for {
			PrintWithSleep("1")

			m2.Lock()
			cond2.Signal()
			m2.Unlock()

			m1.Lock()
			cond1.Wait()
			m1.Unlock()
		}
	}()
	go func() {
		for {
			m2.Lock()
			cond2.Wait()
			m2.Unlock()

			PrintWithSleep("a")

			m1.Lock()
			cond1.Signal()
			m1.Unlock()
		}
	}()
	time.Sleep(10*time.Second)
}

func PrintWithSleep(str string){
	fmt.Println(str)
	time.Sleep(1*time.Second)
}
