package data

import (
	"github.com/docker/docker/client"
)

type Docker struct {
	*client.Client
}

func NewDockerClient() (d Docker) {
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	d.Client = c

	return d
}
