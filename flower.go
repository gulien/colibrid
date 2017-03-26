package main

type Flower struct {
	containerId 	string
	containerName 	string
	path        	string
}

func NewFlower(containerId string, containerName string, path string) *Flower {
	flower := &Flower {
		containerId: 	containerId,
		containerName: 	containerName,
		path:        	path,
	}

	return flower
}