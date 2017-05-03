package main

import (
	"errors"

	"gopkg.in/yaml.v2"
)

type (
	// Flower struct represents a container which is exposing commands.
	Flower struct {
		Container  *Container
		Path       string
		FlowerData *FlowerData
	}

	// FlowerData struct represents a YAML file defining commands.
	FlowerData struct {
		Commands []*FlowerCommandData `yaml:"Commands"`
	}

	// FlowerCommandData struct represents a section in the YAML file defining a command.
	FlowerCommandData struct {
		Name              string                  `yaml:"Name"`
		Usage             string                  `yaml:"Usage,omitempty"`
		Help              string                  `yaml:"Help,omitempty"`
		Bin               string                  `yaml:"Bin"`
		Workdir           string                  `yaml:"Workdir,omitempty"`
		DockerExecOptions *DockerExecOptions      `yaml:"DockerExecOptions,omitempty"`
		Sub               []*FlowerCommandSubData `yaml:"Sub,omitempty"`
	}

	// FlowerCommandSubData struct represents a section in the YAML file defining
	// option/value/sub-command of a command or another option/value/sub-command.
	FlowerCommandSubData struct {
		Name string                  `yaml:"Name"`
		Sub  []*FlowerCommandSubData `yaml:"Sub,omitempty"`
	}
)

// NewFlower function instantiates a Flower.
func NewFlower(container *Container, path string) *Flower {
	return &Flower{
		Container: container,
		Path:      path,
	}
}

// Parse function retrieves data contained in a YAML file which path has been defined in the FLOWER_PATH
// container's environment variable.
// It also populates the FlowerData variable of the Flower instance.
// This method mights throw the errors encountered while running Container Exec or yaml Unmarshal methods.
func (flower *Flower) Parse() (*FlowerData, error) {
	command := []string{"cat", flower.Path}
	captured, err := flower.Container.Exec(command, &DockerExecOptions{}, true)
	if err != nil {
		return nil, err
	}

	flowerData := &FlowerData{}
	err = yaml.Unmarshal([]byte(captured), flowerData)
	if err != nil {
		return nil, err
	}

	flower.FlowerData = flowerData

	return flowerData, nil
}

// Exec function simply runs an available command from the Flower.
// If capture parameter is set to true, it sends the output of the command into a string.
// This method might throws the errors encountered while running GetFlowerCommandData or Container Exec methods.
func (flower *Flower) Exec(commandName string, capture bool, args []string) (string, error) {
	flowerCommandData, err := flower.GetFlowerCommandData(commandName)
	if err != nil {
		return "", err
	}

	var command []string
	if flowerCommandData.Workdir != "" {
		command = append([]string{"cd", flowerCommandData.Workdir, "&&"})
	}

	command = append([]string{flowerCommandData.Bin})
	for _, arg := range args {
		command = append([]string{arg})
	}

	var dockerExecOptions *DockerExecOptions
	switch flowerCommandData.DockerExecOptions {
	case nil:
		dockerExecOptions = flowerCommandData.DockerExecOptions
	default:
		dockerExecOptions = &DockerExecOptions{}
	}

	return flower.Container.Exec(command, dockerExecOptions, capture)
}

// GetFlowerCommandData returns a FlowerCommandData.
// If the Flower instance has not been parsed or the command does not exist, throws an error.
func (flower *Flower) GetFlowerCommandData(commandName string) (*FlowerCommandData, error) {
	if flower.FlowerData == nil {
		return nil, errors.New("Flower has not been parsed.")
	}

	for _, flowerCommandData := range flower.FlowerData.Commands {
		if flowerCommandData.Name == commandName {
			return flowerCommandData, nil
		}
	}

	return nil, errors.New("Command not found.")
}
