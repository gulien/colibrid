package main

import (
	"errors"

	"github.com/fsouza/go-dockerclient"
)

// Colibri struct helps for discovering
// and caching Flowers.
type Colibri struct {
	client        *docker.Client
	cache         map[string]*Flower
	CurrentFlower *Flower
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
// and populates its cache. Returns the number of these containers.
func (colibri *Colibri) Discover() int {
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

	return len(colibri.cache)
}

// FlyTo function is a wrapper of GetFlower function and Flower's Parse
// function. It also populates the CurrentFlower variable of the Colibri
// instance.
func (colibri *Colibri) FlyTo(identifier string) (*FlowerData, error) {
	colibri.CurrentFlower = colibri.GetFlower(identifier)

	if colibri.CurrentFlower == nil {
		return nil, errors.New("Unknown container: is it a flower?")
	}

	flowerData, err := colibri.CurrentFlower.Parse()
	if err != nil {
		return flowerData, err
	}

	return flowerData, nil
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
	return append(colibri.ListNames(), colibri.ListShortIDs()...)
}

// ListShortIDs function returns the list of containers' short ids
// from its cache.
func (colibri *Colibri) ListShortIDs() []string {
	shortIDs := make([]string, len(colibri.cache))

	i := 0
	for _, flower := range colibri.cache {
		shortIDs[i] = flower.Container.ShortID
		i++
	}

	return shortIDs
}

// ListNames function returns the list of containers' names
// from its cache.
func (colibri *Colibri) ListNames() []string {
	names := make([]string, len(colibri.cache))

	i := 0
	for _, flower := range colibri.cache {
		names[i] = flower.Container.Name
		i++
	}

	return names
}
