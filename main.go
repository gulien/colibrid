package main

import "fmt"

func main() {
	source := NewSource()
	source.findFlowers()

	for _, flower := range source.flowers {
		fmt.Printf(flower.flowerPath)
	}
}