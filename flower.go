package main

type Flower struct {
	containerId string
	path        string
}

func NewFlower(containerId string, path string) *Flower {
	flower := &Flower {
		containerId: containerId,
		path:        path,
	}

	return flower
}