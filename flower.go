package main

type Flower struct {
	containerId 	string
	flowerPath 	string
}

func NewFlower(containerId string, flowerPath string) *Flower {
	flower := &Flower {
		containerId: 	containerId,
		flowerPath:	flowerPath,
	}

	return flower
}