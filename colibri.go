package main

import (
	"github.com/fsouza/go-dockerclient"
)

type Colibri struct {
	Client 	*docker.Client
	Cache 	map[string]*Flower
}

func NewColibri() *Colibri {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	return &Colibri{
		Client:	client,
	}
}

func (colibri *Colibri) Refresh() {
	// lists all running containers
	opts := docker.ListContainersOptions{All: true}
	containersInfo, err := colibri.Client.ListContainers(opts)
	if err != nil {
		panic(err)
	}

	// finds Flowers among running containers
	tmpCache := make(map[string]*Flower)
	for _, containerInfo := range containersInfo {
		switch colibri.Cache[containerInfo.ID] {
		case nil:
			container := NewContainer(colibri.Client, containerInfo.ID)
			path := container.GetEnvValue("FLOWER_PATH")
			if path != "" {
				tmpCache[containerInfo.ID] = NewFlower(container, path)
			}
			break
		default:
			tmpCache[containerInfo.ID] = colibri.Cache[containerInfo.ID]
		}
	}

	// refreshes the cache
	colibri.Cache = tmpCache
}

func (colibri *Colibri) GetFlower(identifier string) *Flower {
	for _, flower := range colibri.Cache {
		if flower.Container.ShortID == identifier || flower.Container.Name == identifier {
			return flower
		}
	}

	return nil
}

func (colibri *Colibri) ListNames() []string {
	names := make([]string, len(colibri.Cache))

	i := 0
	for _, flower := range colibri.Cache {
		names[i] = flower.Container.Name
		i++
	}

	return names
}

func (colibri *Colibri) ListShortIDs() []string {
	shortIDs := make([]string, len(colibri.Cache))

	i := 0
	for _, flower := range colibri.Cache {
		shortIDs[i] = flower.Container.ShortID
		i++
	}

	return shortIDs
}