package main

import (
	"fmt"
	"sync"
)

func SetPro() {
	sm := new(sync.Map)
	for i := 0; i < 10; i++ {
		sm.Store(i, i+'a')
	}
	count := 0
	sm.Range(func(k, v interface{}) bool {
		fmt.Println(k, v)
		count++
		return true
	})
	fmt.Println(count)
}

func main() {
	SetPro()
}
