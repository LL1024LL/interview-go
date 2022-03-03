package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main(){
	testData := MakeTestData(1000000, 10)
	res, err := CurrentFind(2,testData ,11, 1*time.Second )
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func MakeTestData(n int, upper int) []int {
	rand.Seed(time.Now().Unix())
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(upper)
	}
	return arr[:]
}

var (
	ErrNotFind = errors.New("Not Found")
	ErrTimeOut = errors.New("Timeout, Not Found")
)

func CurrentFind(curCnt int, input []int, target int, timeout time.Duration) (res bool, err error) {

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	stop := make(chan struct{})
	var stopped int32
	// 超时控制
	go func() {
		select {
		case <-timer.C:
			if atomic.CompareAndSwapInt32(&stopped, 0, 1) {
				close(stop)
				err = ErrTimeOut
			}
		}
	}()
	if curCnt > len(input) {
		curCnt = len(input)
	}
	size := len(input) / curCnt
	wg := sync.WaitGroup{}
	for _, arr := range splitArr(input, size) {
		arr := arr
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok := findNum(arr, target, stop)
			if ok && atomic.CompareAndSwapInt32(&stopped, 0, 1) {
				// 找到了
				res = true
				close(stop)
			}
		}()
	}
	wg.Wait()
	// 未超时，未找到
	if !res && err == nil {
		err = ErrNotFind
	}
	return
}

func splitArr(input []int, size int) [][]int {
	var arrs [][]int
	for i := 0; i < len(input); i += size{
		e := i + size
		if e > len(input) {
			e = len(input)
		}
		arrs = append(arrs, input[i:e])
	}
	return arrs
}

func findNum(input []int, target int, stop chan struct{}) bool {
	for _, v := range input {
		select {
		case <-stop:
			return false
		default:
		}
		if v == target {
			return true
		}
	}
	return false
}
