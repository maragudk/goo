package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	g "github.com/maragudk/gomponents"
	ghttp "github.com/maragudk/gomponents/http"
	"maragu.dev/errors"
	"maragu.dev/httph"
	"maragu.dev/snorkel"

	"maragu.dev/goo/html"
	"maragu.dev/goo/model"
)

type contextKey string

const contextUserKey = contextKey("user")
const sessionUserIDKey = "userID"

// getUserFromContext, which may be nil if the user is not authenticated.
func getUserFromContext(ctx context.Context) *model.User {
	user := ctx.Value(contextUserKey)
	if user == nil {
		return nil
	}
	return user.(*model.User)
}

type signupper interface {
	Signup(ctx context.Context, u model.User) (model.User, error)
}

type signupRequest struct {
	Name   string
	Email  model.Email
	Accept bool
}

func (s signupRequest) Validate() error {
	if s.Name == "" {
		return errors.New("name cannot be empty")
	}

	if !s.Email.IsValid() {
		return errors.New("email is invalid")
	}

	if !s.Accept {
		return errors.New("not accepted")
	}

	return nil
}

func Signup(mux chi.Router, page html.PageFunc, log *snorkel.Logger, db signupper) {
	mux.Get("/signup", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		user := getUserFromContext(r.Context())
		if user != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return nil, nil
		}

		return html.SignupPage(page, html.PageProps{}, model.User{}), nil
	}))

	mux.Post("/signup", httph.FormHandler(func(w http.ResponseWriter, r *http.Request, req signupRequest) {
		h := ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
			// TODO this should be middleware
			user := getUserFromContext(r.Context())
			if user != nil {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return nil, nil
			}

			u := model.User{Name: req.Name, Email: req.Email}

			if u, err := db.Signup(r.Context(), u); err != nil {
				if errors.Is(err, model.ErrorEmailConflict) {
					return html.SignupPage(page, html.PageProps{}, u), nil
				}
				log.Event("Error signing up", 1, "error", err)
				return html.ErrorPage(page), err
			}

			http.Redirect(w, r, "/signup/thanks", http.StatusFound)
			return nil, nil
		})

		h(w, r)
	}))
}
