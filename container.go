package main

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

// Container struct represents a Docker's container.
type Container struct {
	client  *docker.Client
	ID      string
	ShortID string
	Name    string
	Env     []string
}

// NewContainer function instantiates a Container.
func NewContainer(client *docker.Client, id string) *Container {
	container := &Container{
		client:  client,
		ID:      id,
		ShortID: id[:12],
	}

	inspected, err := client.InspectContainer(container.ID)
	if err != nil {
		panic(err)
	}

	container.Name = strings.TrimPrefix(inspected.Name, "/")
	container.Env = inspected.Config.Env

	return container
}

// Exec function runs a command from current Container instance.
// If capture parameter is set to true, it sends the output of the command into a string.
func (container *Container) Exec(command []string, capture bool) (string, error) {
	captured := ""

	var createExecOptions docker.CreateExecOptions
	if capture {
		createExecOptions = docker.CreateExecOptions{
			AttachStdin:  false,
			AttachStdout: true,
			AttachStderr: false,
			Tty:          false,
			Cmd:          command,
			Container:    container.ID,
		}
	} else {
		createExecOptions = docker.CreateExecOptions{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
			Cmd:          command,
			Container:    container.ID,
		}
	}

	exec, err := container.client.CreateExec(createExecOptions)
	if err != nil {
		return captured, err
	}

	if capture {
		reader, writer, _ := os.Pipe()

		startExecOptions := docker.StartExecOptions{
			OutputStream: writer,
			RawTerminal:  false,
		}

		err = container.client.StartExec(exec.ID, startExecOptions)
		if err != nil {
			return captured, err
		}

		commandOutput := make(chan string)
		go func() {
			var buffer bytes.Buffer
			io.Copy(&buffer, reader)
			commandOutput <- buffer.String()
		}()

		writer.Close()
		captured = <-commandOutput
	} else {
		startExecOptions := docker.StartExecOptions{
			InputStream:  os.Stdin,
			OutputStream: os.Stdout,
			ErrorStream:  os.Stderr,
			RawTerminal:  true,
		}

		err = container.client.StartExec(exec.ID, startExecOptions)
		if err != nil {
			return captured, err
		}
	}

	return captured, nil
}

// GetEnvValue function retrieves the value of an environment variable of
// the current Container instance. If no environment variable found,
// returns an empty string.
func (container *Container) GetEnvValue(keyName string) string {
	for _, envStr := range container.Env {
		envStrParts := strings.SplitN(envStr, "=", 2)
		if envStrParts[0] == keyName {
			return envStrParts[1]
		}
	}

	return ""
}
