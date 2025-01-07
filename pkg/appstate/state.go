package appstate

import (
	docker "github.com/docker/docker/client"
)

type State struct {
	// Current height of the content element
	ContentHeight int

	DockerClient *docker.Client
}
