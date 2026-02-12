package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := range 10 {
			ch <- i*2 + 1
		}
	}()

	go func() {
		defer wg.Done()
		for i := range 10 {
			ch <- i * 2
		}
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := range ch {
		fmt.Printf("%d が書き込まれました\n", i)
	}
}
