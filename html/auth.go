package html

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents-heroicons/v2/mini"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"

	"maragu.dev/goo/model"
)

func SignupPage(page PageFunc, props PageProps, newUser model.User) g.Node {
	props.Title = "Sign up"

	return page(props,
		authPageCard(
			Form(Action("/signup"), Method("post"), Class("space-y-6"), hx.Boost("false"),
				Div(Class("text-center"),
					h1(g.Text(`Sign up`)),
					a(Href("/login"), g.Text("or log in")),
				),

				g.If(newUser.Email.String() != "",
					alertBox(g.Raw(`There’s already a user with this email address. `), a(Href("/login"), g.Text("Log in instead?"))),
				),

				Div(
					label("name", "Name"),
					input(Type("text"), ID("name"), Name("name"), Value(newUser.Name), AutoComplete("name"),
						Placeholder("Me"), Required(), g.If(newUser.Name == "", AutoFocus())),
				),

				Div(
					label("email", "Email"),
					input(Type("email"), ID("email"), Name("email"), Value(newUser.Email.String()), AutoComplete("email"),
						Placeholder("me@example.com"), Required()),
				),

				Div(Class("flex items-center space-x-2"),
					Input(ID("accept"), Name("accept"), Type("checkbox"), Value("true"), Required(),
						Class("h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500")),
					Label(For("accept"), Class("text-gray-900"),
						g.Text(`I accept `),
						a(Href("/legal/terms-of-service"), Target("_blank"), g.Text(`Terms of Service`)),
						g.Text(` and `),
						a(Href("/legal/privacy-policy"), Target("_blank"), g.Text(`Privacy Policy`)),
						g.Text(`.`),
					),
				),

				button(Type("submit"), g.Text(`Sign up`)),
			),
		),
	)
}

func SignupThanksPage(page PageFunc, props PageProps) g.Node {
	props.Title = "Thanks!"

	return page(props,
		authPageCard(
			Div(Class("text-center"),
				h1(g.Text(`Thanks!`)),
				g.Text(`Now check your inbox.`),
			),
		),
	)
}

func authPageCard(children ...g.Node) g.Node {
	return Div(Class("sm:mx-auto sm:w-full sm:max-w-md"),
		card(g.Group(children)),
	)
}

func card(children ...g.Node) g.Node {
	return Div(Class("bg-white py-8 px-4 shadow rounded-lg sm:px-10"), g.Group(children))
}

func label(id, text string) g.Node {
	return Label(For(id), Class("block text-sm text-gray-700 mb-1"), g.Text(text))
}

func input(children ...g.Node) g.Node {
	return Input(Class("block w-full rounded-md border border-gray-300 focus:border-primary-500 px-3 py-2 placeholder-gray-400 shadow-sm sm:text-sm text-gray-900 focus:ring-primary-500"), g.Group(children))
}

func button(children ...g.Node) g.Node {
	return Button(Class("block w-full rounded-md bg-primary-600 hover:bg-primary-700 px-4 py-2 font-medium text-white focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"), g.Group(children))
}

func h1(children ...g.Node) g.Node {
	return H1(Class("font-medium text-gray-900 text-xl"), g.Group(children))
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
