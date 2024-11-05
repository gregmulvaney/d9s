package context

import (
	docker "github.com/docker/docker/client"
)

type ProgramContext struct {
	ScreenWidth       int
	ScreenHeight      int
	MainContentHeight int
	ShowCommandView   bool
	DockerClient      *docker.Client
}
