package html

import (
	g "github.com/maragudk/gomponents"

	"maragu.dev/goo/model"
)

type PageProps struct {
	Title       string
	Description string
	User        *model.User
}

type PageFunc = func(props PageProps, children ...g.Node) g.Node
