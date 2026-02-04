package main

import "fmt"

func main() {
	var total int
	for i := range 10 {
		total = total + i
		fmt.Printf("i=%v total=%v\n", i, total)
	}
	fmt.Println("total=", total)
}
