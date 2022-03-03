package main

import (
	"fmt"
	"sync"
	"time"
)

//GO里面MAP如何实现key不存在
//get操作等待 直到key存在或者超时 保证并发安全，且需要实现以下接口

type sp interface {
	//存入key /val，如果该key读取的goroutine挂起，则唤醒。此方法不会阻塞，时刻都可以立即执行并返回
	Out(key string, val interface{})
	//读取一个key，如果key不存在阻塞，等待key存在或者超时
	Rd(key string, timeout time.Duration) interface{}
}

type mmap struct {
	m         sync.Map
	notifyChs map[string]chan struct{}
	mux       sync.Mutex
}

func NewMmap() *mmap {
	return &mmap{
		m:         sync.Map{},
		notifyChs: make(map[string]chan struct{}),
		mux:       sync.Mutex{},
	}
}

func (m *mmap) Out(key string, val interface{}) {
	m.m.Store(key, val)
	m.mux.Lock()
	defer m.mux.Unlock()
	ch, ok := m.notifyChs[key]
	if !ok {
		return
	}
	close(ch)
	delete(m.notifyChs, key)
}

func (m *mmap) Rd(key string, timeout time.Duration) interface{} {
	v, ok := m.m.Load(key)
	if ok {
		return v
	}

	to := time.NewTimer(timeout)
	defer to.Stop()

	m.mux.Lock()
	ch, ok := m.notifyChs[key]
	if !ok {
		ch = make(chan struct{})
		m.notifyChs[key] = ch
	}
	m.mux.Unlock()

	select {
	case <-ch:
		v, _ := m.m.Load(key)
		return v
	case <-to.C:
		return nil
	}
}

func main() {
	mm := NewMmap()
	for i := 0 ; i < 5 ; i++{
		go func() {
			v := mm.Rd("1", 5*time.Second)
			fmt.Printf("%d get %+v\n", 1, v)
		}()
	}

	time.Sleep(1*time.Second)
	mm.Out("1","ll")


	select {}
}
