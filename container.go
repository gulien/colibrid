package main

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

type Container struct {
	Client 	*docker.Client
	ID	string
	ShortID	string
	Name 	string
	Env 	[]string
}

func NewContainer(client *docker.Client, id string) *Container {
	container := &Container{
		Client:		client,
		ID:		id,
		ShortID: 	id[:12],
	}

	inspected, err := client.InspectContainer(container.ID)
	if err != nil {
		panic(err)
	}

	container.Name = strings.TrimPrefix(inspected.Name, "/")
	container.Env = inspected.Config.Env

	return container
}

func (container *Container) Exec(command []string, capture bool) (string, error) {
	captured := ""
	var createExecOptions docker.CreateExecOptions

	if capture {
		createExecOptions = docker.CreateExecOptions{
			AttachStdin: 	false,
			AttachStdout:	true,
			AttachStderr: 	false,
			Tty:		false,
			Cmd:		command,
			Container:	container.ID,
		}
	} else {
		createExecOptions = docker.CreateExecOptions{
			AttachStdin: 	true,
			AttachStdout:	true,
			AttachStderr: 	true,
			Tty:		true,
			Cmd:		command,
			Container:	container.ID,
		}
	}

	exec, err := container.Client.CreateExec(createExecOptions)
	if err != nil {
		return captured, err
	}

	if capture {
		reader, writer, _ := os.Pipe()

		startExecOptions := docker.StartExecOptions{
			OutputStream: writer,
			RawTerminal:  false,
		}

		err = container.Client.StartExec(exec.ID, startExecOptions)
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
		captured = <- commandOutput
	} else {
		startExecOptions := docker.StartExecOptions{
			InputStream: 	os.Stdin,
			OutputStream: 	os.Stdout,
			ErrorStream: 	os.Stderr,
			RawTerminal: 	true,
		}

		err = container.Client.StartExec(exec.ID, startExecOptions)
		if err != nil {
			return captured, err
		}
	}

	return captured, nil
}

func (container *Container) GetEnvValue(keyName string) string {
	for _, envStr := range container.Env {
		envStrParts := strings.SplitN(envStr, "=", 2)
		if envStrParts[0] == keyName {
			return envStrParts[1]
		}
	}

	return ""
}