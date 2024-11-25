package appcontext

import "github.com/gregmulvaney/d9s/pkg/data"

type Context struct {
	Docker data.Docker

	ScreenHeight     int
	ScreenWidth      int
	MaxContentHeight int
}
