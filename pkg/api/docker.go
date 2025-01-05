package api

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
)

type Api struct {
	Client *docker.Client
}

func Init() (a Api) {
	client, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	a.Client = client

	return a
}

func (a Api) FetchContainers() ([]types.Container, error) {
	containers, err := a.Client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (a Api) StartContainer(id string) error {
	err := a.Client.ContainerStart(context.Background(), id, container.StartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (a Api) StopContainer(id string) error {
	err := a.Client.ContainerStop(context.Background(), id, container.StopOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (a Api) RestartContainer(id string) error {
	err := a.Client.ContainerRestart(context.Background(), id, container.StopOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (a Api) RemoveContainer(id string) error {
	err := a.Client.ContainerRemove(context.Background(), id, container.RemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (a Api) ImagesFetch() ([]image.Summary, error) {
	images, err := a.Client.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (a Api) ImagesPrune() (image.PruneReport, error) {
	report, err := a.Client.ImagesPrune(context.Background(), filters.Args{})
	if err != nil {
		return image.PruneReport{}, err
	}
	return report, nil
}

func (a Api) NetworksFetch() ([]network.Summary, error) {
	networks, err := a.Client.NetworkList(context.Background(), network.ListOptions{})
	if err != nil {
		return nil, err
	}
	return networks, nil
}
