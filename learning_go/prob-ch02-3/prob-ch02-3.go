package main

import "fmt"

func main() {
	var b byte = 1<<8 - 1
	var smallI int32 = 1<<31 - 1
	var bigI uint64 = 1<<64 - 1

	fmt.Println("b:", b)
	fmt.Println("smallI:", smallI)
	fmt.Println("bigI:", bigI)

	fmt.Println("b+1:", b+1)
	fmt.Println("smallI+1:", smallI+1)
	fmt.Println("bigI+1:", bigI+1)
}
