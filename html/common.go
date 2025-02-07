package html

import (
	"context"
	"net/http"

	. "maragu.dev/gomponents"
	"maragu.dev/gomponents-heroicons/v3/mini"
	. "maragu.dev/gomponents/html"

	"maragu.dev/goo/model"
)

type PageProps struct {
	Title       string
	Description string
	User        *model.User
	Ctx         context.Context
	Req         *http.Request
}

type PageFunc = func(props PageProps, children ...Node) Node

func FavIcons(name, themeColor string) Node {
	return Group([]Node{
		// <link rel="icon" type="image/png" href="/favicon-96x96.png" sizes="96x96" />
		Link(Rel("icon"), Type("image/png"), Href("/favicon-96x96.png"), Attr("sizes", "96x96")),
		// <link rel="icon" type="image/svg+xml" href="/favicon.svg" />
		Link(Rel("icon"), Type("image/svg+xml"), Href("/favicon.svg")),

		// <link rel="shortcut icon" href="/favicon.ico" />
		Link(Rel("shortcut icon"), Href("/favicon.ico")),

		// <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png" />
		Link(Rel("apple-touch-icon"), Attr("sizes", "180x180"), Href("/apple-touch-icon.png")),

		// <meta name="apple-mobile-web-app-title" content="name" />
		Meta(Name("apple-mobile-web-app-title"), Content(name)),

		// <link rel="manifest" href="/site.webmanifest" />
		Link(Rel("manifest"), Href("/manifest.json")),
	})
}

func card(children ...Node) Node {
	return Div(Class("bg-white py-8 px-4 shadow rounded-lg sm:px-10"), Group(children))
}

func label(id, text string) Node {
	return Label(For(id), Class("block text-sm text-gray-700 mb-1"), Text(text))
}

func LabelAndInput(name string, children ...Node) Node {
	return Div(
		label(name, name),
		input(ID(name), Name(name), Group(children)),
	)
}

func input(children ...Node) Node {
	return Input(Class("block w-full rounded-md border border-gray-300 focus:border-primary-500 px-3 py-2 placeholder-gray-400 shadow-sm sm:text-sm text-gray-900 focus:ring-primary-500"), Group(children))
}

func ButtonPrimary(children ...Node) Node {
	return Button(Class("block w-full rounded-md bg-primary-600 hover:bg-primary-700 px-4 py-2 font-medium text-white focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 font-sans"), Group(children))
}

func h1(children ...Node) Node {
	return H1(Class("font-medium text-gray-900 text-xl"), Group(children))
}

func p(class string, children ...Node) Node {
	return P(Class("text-gray-900 "+class), Group(children))
}

func a(children ...Node) Node {
	return A(Class("text-primary-600 hover:text-primary-500"), Group(children))
}

func alertBox(children ...Node) Node {
	return Div(Class("rounded-md bg-yellow-50 p-4"),
		Div(Class("flex items-center space-x-2"),
			Div(Class("flex-shrink-0"),
				mini.ExclamationTriangle(Class("h-5 w-5 text-yellow-400")),
			),
			Div(Class("text-yellow-700"),
				Group(children),
			),
		),
	)
}
