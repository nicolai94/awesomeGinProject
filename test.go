package main

import (
	"fmt"
	"sync"
)

func say(num int64, ch chan int64, wg *sync.WaitGroup) {
	defer wg.Done()
	counter := int64(1)
	for i := int64(0); i < num; i++ {
		counter++
		counter = counter * 2
	}
	ch <- counter
}

func main() {
	var wg sync.WaitGroup
	ch1 := make(chan int64)
	ch2 := make(chan int64)
	l1 := int64(30)
	l2 := int64(20)
	wg.Add(2)
	go say(l1, ch1, &wg)
	go say(l2, ch2, &wg)
	res1 := <-ch1
	res2 := <-ch2
	wg.Wait()

	res := res1 + res2
	fmt.Println(res)
}
