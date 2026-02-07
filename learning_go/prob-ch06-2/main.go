package main

import "fmt"

func UpdateSlice(ss []string, s string) {
	ss[len(ss)-1] = s
	fmt.Println(ss)
}

func GrowSlice(ss []string, s string) {
	ss = append(ss, s)
	fmt.Println(ss)
}

func main() {
	ss1 := []string{"1", "2", "3"}
	ss2 := []string{"1", "2", "3"}

	fmt.Println("Before function calling:")
	fmt.Println(ss1)
	fmt.Println(ss2)

	fmt.Println("\nInside function calling:")
	UpdateSlice(ss1, "4")
	GrowSlice(ss2, "4")

	fmt.Println("\nAfter function calling:")
	fmt.Println(ss1)
	fmt.Println(ss2)
}
