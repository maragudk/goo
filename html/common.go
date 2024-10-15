package html

import (
	"context"
	"net/http"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents-heroicons/v2/mini"
	. "github.com/maragudk/gomponents/html"

	"maragu.dev/goo/model"
)

type PageProps struct {
	Title       string
	Description string
	User        *model.User
	Ctx         context.Context
	Req         *http.Request
}

type PageFunc = func(props PageProps, children ...g.Node) g.Node

func FavIcons(themeColor string) g.Node {
	return g.Group([]g.Node{
		Link(Rel("apple-touch-icon"), g.Attr("sizes", "180x180"), Href("/apple-touch-icon.png")),
		Link(Rel("icon"), Type("image/png"), g.Attr("sizes", "32x32"), Href("/favicon-32x32.png")),
		Link(Rel("icon"), Type("image/png"), g.Attr("sizes", "16x16"), Href("/favicon-16x16.png")),
		Link(Rel("manifest"), Href("/manifest.json")),
		Meta(Name("msapplication-TileColor"), Content(themeColor)),
		Meta(Name("theme-color"), Content(themeColor)),
	})
}

func card(children ...g.Node) g.Node {
	return Div(Class("bg-white py-8 px-4 shadow rounded-lg sm:px-10"), g.Group(children))
}

func label(id, text string) g.Node {
	return Label(For(id), Class("block text-sm text-gray-700 mb-1"), g.Text(text))
}

func LabelAndInput(name string, children ...g.Node) g.Node {
	return Div(
		label(name, name),
		input(ID(name), g.Group(children)),
	)
}

func input(children ...g.Node) g.Node {
	return Input(Class("block w-full rounded-md border border-gray-300 focus:border-primary-500 px-3 py-2 placeholder-gray-400 shadow-sm sm:text-sm text-gray-900 focus:ring-primary-500"), g.Group(children))
}

func button(children ...g.Node) g.Node {
	return Button(Class("block w-full rounded-md bg-primary-600 hover:bg-primary-700 px-4 py-2 font-medium text-white focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 font-sans"), g.Group(children))
}

func h1(children ...g.Node) g.Node {
	return H1(Class("font-medium text-gray-900 text-xl"), g.Group(children))
}

func p(class string, children ...g.Node) g.Node {
	return P(Class("text-gray-900 "+class), g.Group(children))
}

func a(children ...g.Node) g.Node {
	return A(Class("text-primary-600 hover:text-primary-500"), g.Group(children))
}

func alertBox(children ...g.Node) g.Node {
	return Div(Class("rounded-md bg-yellow-50 p-4"),
		Div(Class("flex items-center space-x-2"),
			Div(Class("flex-shrink-0"),
				mini.ExclamationTriangle(Class("h-5 w-5 text-yellow-400")),
			),
			Div(Class("text-yellow-700"),
				g.Group(children),
			),
		),
	)
}
