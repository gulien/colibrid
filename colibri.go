package main

import (
	"github.com/fsouza/go-dockerclient"
)

// Colibri struct helps for discovering
// and caching Flowers.
type Colibri struct {
	client *docker.Client
	cache  map[string]*Flower
}

// NewColibri function instantiates a Colibri.
func NewColibri() *Colibri {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	return &Colibri{
		client: client,
	}
}

// Discover function finds running containers which are exposing commands
// and populates its cache.
func (colibri *Colibri) Discover() {
	// lists all running containers
	opts := docker.ListContainersOptions{All: true}
	containersInfo, err := colibri.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}

	// finds Flowers among running containers
	tmpCache := make(map[string]*Flower)
	for _, containerInfo := range containersInfo {
		switch colibri.cache[containerInfo.ID] {
		case nil:
			container := NewContainer(colibri.client, containerInfo.ID)
			path := container.GetEnvValue("FLOWER_PATH")
			if path != "" {
				tmpCache[containerInfo.ID] = NewFlower(container, path)
			}
		default:
			tmpCache[containerInfo.ID] = colibri.cache[containerInfo.ID]
		}
	}

	// refreshes the cache
	colibri.cache = tmpCache
}

// GetFlower function returns a Flower by its short id or name.
// If there is no corresponding Flower in its cache, returns nil.
func (colibri *Colibri) GetFlower(identifier string) *Flower {
	for _, flower := range colibri.cache {
		if flower.Container.ShortID == identifier || flower.Container.Name == identifier {
			return flower
		}
	}

	return nil
}

// ListIdentifiers function returns the list of containers' short ids
// and names from its cache.
func (colibri *Colibri) ListIdentifiers() []string {
	return append(colibri.listNames(), colibri.listShortIDs()...)
}

// listShortIDs function returns the list of containers' short ids
// from its cache.
func (colibri *Colibri) listShortIDs() []string {
	shortIDs := make([]string, len(colibri.cache))

	i := 0
	for _, flower := range colibri.cache {
		shortIDs[i] = flower.Container.ShortID
		i++
	}

	return shortIDs
}

// listNames function returns the list of containers' names
// from its cache.
func (colibri *Colibri) listNames() []string {
	names := make([]string, len(colibri.cache))

	i := 0
	for _, flower := range colibri.cache {
		names[i] = flower.Container.Name
		i++
	}

	return names
}
