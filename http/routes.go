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

		r.Group(func(r chi.Router) {
			r.Use(httph.VersionedAssets)

			Static(r)
		})

		r.Group(func(r chi.Router) {
			r.Use(s.sm.LoadAndSave)
			r.Use(httph.NoClickjacking)

			if s.setupUserRoutes != nil {
				s.setupUserRoutes(r)
			}
		})
	})
}
