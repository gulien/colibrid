package main

import (
	"strings"

	"github.com/fsouza/go-dockerclient"
)

type Source struct {
	client 	*docker.Client
	flowers map[string]*Flower
}

func NewSource() *Source {
	// initializes Docker client
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	// creates Source by registering Docker client
	source := &Source {
		client:	client,
	}

	return source
}

func (source *Source) findFlowers() {
	// initializes Flowers mapped by container id
	source.flowers = make(map[string]*Flower)

	// lists all running containers
	opts := docker.ListContainersOptions{All: true}
	containersInfo, err := source.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}

	// finds Flowers among running containers
	for _, containerInfo := range containersInfo {
		flower := source.createFlowerIfExists(containerInfo.ID)
		if flower != nil {
			source.flowers[containerInfo.ID] = flower
		}
	}
}

func (source *Source) createFlowerIfExists(containerId string) *Flower {
	container, err := source.client.InspectContainer(containerId)
	if err != nil {
		panic(err)
	}

	// parses env variables in order to find FLOWER_PATH value
	for _, envVariable := range container.Config.Env {
		envVariableParts := strings.Split(envVariable, "=")
		if envVariableParts[0] == "FLOWER_PATH" {
			// yaa! FLOWER_PATH found, let's create a Flower!
			flower := NewFlower(containerId, container.Name, envVariableParts[1])
			return flower
		}
	}

	return nil
}

func (source *Source) getFlower(containerIdOrName string) *Flower {
	flower := source.flowers[containerIdOrName]
	// alright, no flower found by container id, let's find it by container name
	if flower == nil {
		for _, flower := range source.flowers {
			if flower.containerName == containerIdOrName {
				return flower
			}
		}
	}

	return flower
}