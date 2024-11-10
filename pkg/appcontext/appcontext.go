package appcontext

import (
	docker "github.com/docker/docker/client"
)

type Context struct {
	ScreenWidth       int
	ScreenHeight      int
	MainContentHeight int
	ShowCommandView   bool
	DockerClient      *docker.Client
	Keymap            []KeyItem
}

type KeyItem struct {
	Key   string
	Value string
}
