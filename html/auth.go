package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	"maragu.dev/goo/model"
)

func SignupPage(page PageFunc, newUser model.User) Node {
	props := PageProps{Title: "Sign up"}

	return page(props,
		authPageCard(
			Form(Action("/signup"), Method("post"), Class("space-y-6"),
				Div(Class("text-center"),
					h1(Text(`Sign up`)),
					a(Href("/login"), Text("or log in")),
				),

				If(newUser.Email.String() != "",
					alertBox(Raw(`There’s already a user with this email address. `), a(Href("/login"), Text("Log in instead?"))),
				),

				Div(
					label("name", "Name"),
					input(Type("text"), ID("name"), Name("name"), Value(newUser.Name), AutoComplete("name"),
						Placeholder("Me"), Required(), If(newUser.Name == "", AutoFocus())),
				),

				Div(
					label("email", "Email"),
					input(Type("email"), ID("email"), Name("email"), Value(newUser.Email.String()), AutoComplete("email"),
						Placeholder("me@example.com"), Required()),
				),

				Div(Class("flex items-center space-x-2"),
					Input(ID("accept"), Name("accept"), Type("checkbox"), Value("true"), Required(),
						Class("h-4 w-4 rounded border-gray-300 text-sky-600 focus:ring-sky-500")),
					Label(For("accept"), Class("text-gray-900"),
						Text(`I accept `),
						a(Href("https://www.maragu.dev/p/terms-of-service"), Target("_blank"), Text(`Terms of Service`)),
						Text(` and `),
						a(Href("https://www.maragu.dev/p/privacy-policy"), Target("_blank"), Text(`Privacy Policy`)),
						Text(`.`),
					),
				),

				ButtonPrimary(Type("submit"), Text(`Sign up`)),
			),
		),
	)
}

func SignupThanksPage(page PageFunc) Node {
	props := PageProps{Title: "Thanks!"}

	return page(props,
		authPageCard(
			Div(Class("text-center"),
				h1(Text(`Thanks!`)),
				Text(`Now check your inbox.`),
			),
		),
	)
}

func LoginPage(page PageFunc, email model.Email) Node {
	props := PageProps{Title: "Log in"}

	return page(props,
		authPageCard(
			Form(Action("/login/email"), Method("post"), Class("space-y-6"),
				Div(Class("text-center"),
					h1(Text(`Log in`)),
					a(Href("/signup"), Text("or sign up")),
				),

				If(email.String() != "",
					alertBox(Raw(`It doesn’t look like anyone’s signed up with that email address. `), a(Href("/signup"), Text("Sign up instead?"))),
				),

				Div(
					label("email", "Email"),
					input(Type("email"), ID("email"), Name("email"), AutoComplete("email"), Placeholder("me@example.com"), Required(), AutoFocus(), Value(email.String())),
				),

				ButtonPrimary(Type("submit"), Text(`Log in`)),
			),
		),
	)
}

func LoginCheckEmailPage(page PageFunc) Node {
	props := PageProps{Title: "Check your email"}

	return page(props,
		authPageCard(
			Div(Class("text-center"),
				h1(Text(`Check your email`)),
			),

			p("mt-8", Raw(`Now check your email and click the link in it.`)),
		),
	)
}

func LoginSubmitTokenPage(page PageFunc, token string) Node {
	props := PageProps{Title: "Log in"}

	return page(props,
		authPageCard(
			Form(Action("/login/token"), Method("post"),
				Input(Type("hidden"), Name("token"), Value(token)),

				ButtonPrimary(Type("submit"), Text(`Log in`)),
			),
		),
	)
}

func LoginUserInactivePage(page PageFunc) Node {
	props := PageProps{Title: "Your user is inactive"}

	return page(props,
		authPageCard(
			Div(Class("text-center"),
				h1(Text(`Account inactive`)),
			),

			p("mt-8", Raw(`Your account is inactive. If you think this is an error, `),
				a(Href("mailto:support@maragu.dk"), Text("reach out to support")), Text(".")),
		),
	)
}

func LoginTokenExpiredPage(page PageFunc) Node {
	props := PageProps{Title: "Your link has expired"}

	return page(props,
		authPageCard(
			Div(Class("text-center"),
				h1(Text(`Link expired`)),
			),

			p("mt-8", Raw(`Your link has expired. `), a(Href("/login"), Text("Log in again")), Text(".")),
		),
	)
}

func authPageCard(children ...Node) Node {
	return Div(Class("sm:mx-auto sm:w-full sm:max-w-md"),
		card(Group(children)),
	)
}
