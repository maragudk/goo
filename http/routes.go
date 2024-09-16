package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"maragu.dev/httph"
)

func (s *Server) setupRoutes() {
	s.mux.Group(func(r chi.Router) {
		r.Use(middleware.Compress(5))
		r.Use(middleware.RealIP)
		r.NotFound(NotFound(s.htmlPage))

		r.Group(func(r chi.Router) {
			r.Use(httph.VersionedAssets)

			Static(r)
		})

		r.Group(func(r chi.Router) {
			r.Use(s.sm.LoadAndSave)
			r.Use(httph.NoClickjacking)
			r.Use(Authenticate(false, s.sm, s.sqlHelper, s.log))

			// TODO should CSP middleware be loaded here already? Does ordering matter in chi?

			Signup(r, s.htmlPage, s.log, s.sqlHelper)
			Login(r, s.htmlPage, s.log, s.sqlHelper, s.sm)
			Logout(r, s.htmlPage, s.log, s.sm)

			if s.httpRouterInjector != nil {
				s.httpRouterInjector(r)
			}
		})
	})
}
