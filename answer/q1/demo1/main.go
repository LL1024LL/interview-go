package main

import (
	"fmt"
	"time"
)

func main()  {
	ch1 := make(chan struct{},1)
	ch2 := make(chan struct{},1)
	ch1 <- struct{}{}
	go func() {
		for {
			select{
			case <-ch1:
				PrintWithSleep("1")
				ch2 <- struct{}{}
			}
		}
	}()
	go func() {
		for {
			select{
			case <- ch2:
				PrintWithSleep("a")
				ch1 <- struct{}{}
			}
		}
	}()
	select{}
}
func PrintWithSleep(str string){
	fmt.Println(str)
	time.Sleep(1*time.Second)
}
