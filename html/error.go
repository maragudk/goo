package html

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func ErrorPage(page PageFunc) g.Node {
	return page(PageProps{Title: "Something went wrong"},
		H1(g.Text("Something went wrong")),
	)
}

func NotFoundPage(page PageFunc) g.Node {
	return page(PageProps{Title: "Not found"},
		H1(g.Text("Not found")),
	)
}
