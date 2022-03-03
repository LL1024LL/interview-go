package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
写代码实现两个 goroutine，
其中一个产生随机数并写入到 go channel 中，
另外一个从 channel 中读取数字并打印到标准输出。
最终输出五个随机数。
 */
func main() {
	ch := make(chan int, 1)
	rand.Seed(time.Now().Unix())
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0 ; i < 5 ; i++{
			time.Sleep(1*time.Second)
			r := rand.Intn(10)
			ch <-  r
		}
		close(ch)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			v, ok := <- ch
			if !ok {
				break
			}
			fmt.Println(v)
		}
	}()

	wg.Wait()
	select{}
}
