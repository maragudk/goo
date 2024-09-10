package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"maragu.dev/snorkel"

	"maragu.dev/goo/sql"
)

type Server struct {
	adminPassword   string
	baseURL         string
	db              *sql.Database
	log             *snorkel.Logger
	mux             *chi.Mux
	server          *http.Server
	setupUserRoutes func(chi.Router)
	sm              *scs.SessionManager
}

type NewServerOptions struct {
	AdminPassword string
	BaseURL       string
	DB            *sql.Database
	Log           *snorkel.Logger
	Routes        func(chi.Router)
	SecureCookie  bool
}

func NewServer(opts NewServerOptions) *Server {
	if opts.Log == nil {
		opts.Log = snorkel.New(snorkel.Options{W: io.Discard})
	}

	mux := chi.NewMux()

	sm := scs.New()
	sm.Store = sqlite3store.New(opts.DB.DB.DB)
	sm.Lifetime = 365 * 24 * time.Hour
	sm.Cookie.Secure = opts.SecureCookie

	return &Server{
		adminPassword: opts.AdminPassword,
		baseURL:       strings.TrimSuffix(opts.BaseURL, "/"),
		db:            opts.DB,
		log:           opts.Log,
		mux:           mux,
		server: &http.Server{
			Addr:              ":8080",
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
		setupUserRoutes: opts.Routes,
		sm:              sm,
	}
}

func (s *Server) Start() error {
	s.log.Event("Starting http server", 1)

	s.setupRoutes()

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	s.log.Event("Stopping http server", 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.log.Event("Stopped http server", 1)
	return nil
}
