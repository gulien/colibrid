package main

import (
	"gopkg.in/yaml.v2"
)

type (
	// Flower struct represents a container which is exposing commands.
	Flower struct {
		Container *Container
		Path      string
	}

	// FlowerData struct represents a YAML file defining commands.
	FlowerData struct {
		Commands []FlowerCommandData `yaml:"commands"`
	}

	// FlowerCommandData struct represents a section in the YAML file defining a command.
	FlowerCommandData struct {
		Name    string                 `yaml:"name"`
		Bin     string                 `yaml:"bin"`
		Context string                 `yaml:"context,omitempty"`
		User    string                 `yaml:"user,omitempty"`
		Usage   string                 `yaml:"usage,omitempty"`
		Help    string                 `yaml:"help,omitempty"`
		Sub     []FlowerCommandSubData `yaml:"sub,omitempty"`
	}

	// FlowerCommandSubData struct represents a section in the YAML file defining
	// option/value/sub-command of a command or another option/value/sub-command.
	FlowerCommandSubData struct {
		Name  string                 `yaml:"name"`
		Usage string                 `yaml:"usage,omitempty"`
		Help  string                 `yaml:"help,omitempty"`
		Sub   []FlowerCommandSubData `yaml:"sub,omitempty"`
	}
)

// NewFlower function instantiates a Flower.
func NewFlower(container *Container, path string) *Flower {
	return &Flower{
		Container: container,
		Path:      path,
	}
}

// Parse function retrieves data contained in a YAML file
// which path has been defined in the FLOWER_PATH
// container's environment variable.
func (flower *Flower) Parse() (*FlowerData, error) {
	command := []string{"cat", flower.Path}
	captured, err := flower.Container.Exec(command, true)
	if err != nil {
		return nil, err
	}

	flowerData := &FlowerData{}
	err = yaml.Unmarshal([]byte(captured), flowerData)
	if err != nil {
		return nil, err
	}

	return flowerData, nil
}
