package main

import (
	"fmt"
	"time"
)

func main(){
	ch := make(chan struct{})
	go func(){
		ch <- struct{}{}
		fmt.Println("1 done")
	}()
	go func() {
		ch <- struct{}{}
		fmt.Println("2 done")
	}()
	go func() {
		<- ch
	}()
	time.Sleep(1*time.Second)
	close(ch)
	fmt.Println("end")
	select{}
}
