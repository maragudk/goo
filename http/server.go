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

	"maragu.dev/goo/llm"

	"maragu.dev/goo/html"
	"maragu.dev/goo/sql"
)

type Server struct {
	adminPassword      string
	baseURL            string
	htmlPage           html.PageFunc
	httpRouterInjector func(*Router)
	llmClient          *llm.OpenAIClient
	log                *snorkel.Logger
	r                  *Router
	server             *http.Server
	sm                 *scs.SessionManager
	sqlHelper          *sql.Helper
}

type NewServerOptions struct {
	Address            string
	AdminPassword      string
	BaseURL            string
	HTMLPage           html.PageFunc
	HTTPRouterInjector func(*Router)
	LLMClient          *llm.OpenAIClient
	Log                *snorkel.Logger
	SecureCookie       bool
	SQLHelper          *sql.Helper
}

func NewServer(opts NewServerOptions) *Server {
	if opts.Log == nil {
		opts.Log = snorkel.New(snorkel.Options{W: io.Discard})
	}

	if opts.Address == "" {
		opts.Address = ":8080"
	}

	mux := chi.NewMux()

	sm := scs.New()
	sm.Store = sqlite3store.New(opts.SQLHelper.DB.DB)
	sm.Lifetime = 365 * 24 * time.Hour
	sm.Cookie.Secure = opts.SecureCookie

	return &Server{
		adminPassword:      opts.AdminPassword,
		baseURL:            strings.TrimSuffix(opts.BaseURL, "/"),
		httpRouterInjector: opts.HTTPRouterInjector,
		htmlPage:           opts.HTMLPage,
		llmClient:          opts.LLMClient,
		log:                opts.Log,
		r:                  &Router{Mux: mux},
		server: &http.Server{
			Addr:              opts.Address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
		sm:        sm,
		sqlHelper: opts.SQLHelper,
	}
}

func (s *Server) Start() error {
	s.log.Event("Starting http server", 1, "url", s.baseURL)

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
