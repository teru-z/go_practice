package main

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func main() {
	p := Person{
		FirstName: "John",
		LastName:  "Lennon",
		Age:       40,
	}
	loopSize := 10_000_000
	people := []Person{} // capacity未指定
	// people := make([]Person, 0, loopSize) // capacity=loopSize
	for i := 0; i < loopSize; i++ {
		people = append(people, p)
	}
}
