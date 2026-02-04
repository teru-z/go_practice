package main

import (
	"fmt"
	"math/rand"
)

func main() {
	randomInts := []int{}
	for range 100 {
		randomInts = append(randomInts, rand.Intn(100))
	}
	for _, i := range randomInts {
		var message string
		switch {
		case i%2 == 0 && i%3 == 0:
			message = "Six!"
		case i%2 == 0:
			message = "Two!"
		case i%3 == 0:
			message = "Three!"
		default:
			message = "Never mind"
		}
		fmt.Println(message)
	}
}
