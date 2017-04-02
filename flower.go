package main

import (
	"gopkg.in/yaml.v2"
)

type Flower struct {
	Container 	*Container
	Path		string
}

type FlowerData struct {
	Version 	string `yaml:"version,omitempty"`
	Commands 	[]struct {
		Name 	string `yaml:"name"`
		Bin	string `yaml:"bin"`
		Context	string `yaml:"context,omitempty"`
		User 	string `yaml:"user,omitempty"`
		Usage 	string `yaml:"usage,omitempty"`
		Help 	string `yaml:"help,omitempty"`
	} `yaml:"commands"`
}

func NewFlower(container *Container, path string) *Flower {
	return &Flower{
		Container: 	container,
		Path:		path,
	}
}

func (flower *Flower) Parse() (*FlowerData, error) {
	flowerData := &FlowerData{}

	command := []string{"cat", flower.Path}
	captured, err := flower.Container.Exec(command, true)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(captured), flowerData)
	if err != nil {
		return nil, err
	}

	return flowerData, nil
}