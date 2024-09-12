package http_test

import (
	"io"
	http2 "net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"maragu.dev/is"

	"maragu.dev/goo/http"
	"maragu.dev/goo/sqltest"
)

func TestServer_Start(t *testing.T) {
	t.Run("can start and stop server", func(t *testing.T) {
		sqlHelper := sqltest.NewHelper(t)

		httpRouterInjector := func(r chi.Router) {
			r.Get("/", func(w http2.ResponseWriter, r *http2.Request) {
				_, _ = w.Write([]byte("OK"))
			})
		}

		s := http.NewServer(http.NewServerOptions{
			AdminPassword:      "correct horse battery staple",
			BaseURL:            "http://localhost:8080",
			HTTPRouterInjector: httpRouterInjector,
			SecureCookie:       false,
			SQLHelper:          sqlHelper,
		})

		go s.Start()
		defer s.Stop()

		res, err := http2.Get("http://localhost:8080/")
		is.NotError(t, err)
		is.Equal(t, http2.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		is.NotError(t, err)
		is.Equal(t, "OK", string(body))
	})
}
