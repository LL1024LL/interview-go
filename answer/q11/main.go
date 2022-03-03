package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
场景：在一个高并发的web服务器中，要限制IP的频繁访问。

现模拟100个IP同时并发访问服务器，每个IP要重复访问1000次。

每个IP三分钟之内只能访问一次。修改以下代码完成该过程，要求能成功输出 success:100
*/

type Ban struct {
	lock sync.Mutex
	visitIPs map[string]time.Time
}

func NewBan() *Ban {
	return &Ban{visitIPs: map[string]time.Time{}}
}
func (o *Ban) visit(ip string) bool {
	o.lock.Lock()
	defer o.lock.Unlock()

	v, ok := o.visitIPs[ip]
	if ok && !v.Before(time.Now().Add(-2*time.Second)){
		return true
	}
	o.visitIPs[ip] = time.Now()
	return false
}
func main() {
	var success int32
	ban := NewBan()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		for j := 0; j < 100; j++ {
			j := j
			wg.Add(1)
			go func() {
				defer wg.Done()
				ip := fmt.Sprintf("192.168.1.%d", j)
				if !ban.visit(ip) {
					atomic.AddInt32(&success, 1)
				}
			}()
		}
		time.Sleep(1*time.Second)
	}
	wg.Wait()

	fmt.Println("success:", success)
}
