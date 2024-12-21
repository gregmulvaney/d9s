package appcontext

import (
	docker "github.com/docker/docker/client"
)

type Context struct {
	Docker *docker.Client
}
