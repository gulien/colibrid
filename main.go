package main

import "fmt"

func main() {
	source := NewSource()
	source.findFlowers()

	for _, flower := range source.flowersById {
		fmt.Printf(flower.path)
	}
}