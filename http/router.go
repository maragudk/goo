package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	g "github.com/maragudk/gomponents"
	ghttp "github.com/maragudk/gomponents/http"

	"maragu.dev/goo/html"
)

type Router struct {
	Mux chi.Router
}

func (r *Router) Get(path string, cb func(props html.PageProps) (g.Node, error)) {
	r.Mux.Get(path, ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		props := html.PageProps{
			User: getUserFromContext(r.Context()),
			Ctx:  r.Context(),
			Req:  r,
		}
		return cb(props)
	}))
}

func (r *Router) Group(cb func(r *Router)) {
	r.Mux.Group(func(mux chi.Router) {
		cb(&Router{Mux: mux})
	})
}

func (r *Router) Use(middleware func(http.Handler) http.Handler) {
	r.Mux.Use(middleware)
}

func (r *Router) NotFound(h http.HandlerFunc) {
	r.Mux.NotFound(h)
}
