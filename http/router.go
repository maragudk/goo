package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"

	"maragu.dev/goo/html"
)

type Router struct {
	Mux chi.Router
}

func (r *Router) Get(path string, cb func(props html.PageProps) (Node, error)) {
	r.Mux.Get(path, ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return cb(getProps(r))
	}))
}

func (r *Router) Post(path string, cb func(props html.PageProps) (Node, error)) {
	r.Mux.Post(path, ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return cb(getProps(r))
	}))
}

func getProps(r *http.Request) html.PageProps {
	return html.PageProps{
		User: GetUserFromContext(r.Context()),
		Ctx:  r.Context(),
		Req:  r,
	}
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
