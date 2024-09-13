package html

import (
	g "github.com/maragudk/gomponents"
)

type PageProps struct {
	Title       string
	Description string
}

type PageFunc = func(props PageProps, children ...g.Node) g.Node
