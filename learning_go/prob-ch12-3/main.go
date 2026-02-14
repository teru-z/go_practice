package main

import (
	"fmt"
	"math"
	"sync"
)

func buildSquareRootMap() map[int]float64 {
	sqrtMap := map[int]float64{}
	for i := range 100_000 {
		sqrtMap[i] = math.Sqrt(float64(i))
	}
	return sqrtMap
}

func getSqrt(i int) float64 {
	sqrtMap := sync.OnceValue(func() map[int]float64 {
		return buildSquareRootMap()
	})
	return sqrtMap()[i]
}

func main() {
	for i := range 100 {
		n := i * 1000
		fmt.Printf("%d %g\n", n, getSqrt(n))
	}
}
