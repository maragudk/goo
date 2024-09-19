package http_test

import (
	"io"
	http2 "net/http"
	"testing"
	"time"

	"maragu.dev/is"

	"maragu.dev/goo/http"
	"maragu.dev/goo/sqltest"
)

func TestServer_Start(t *testing.T) {
	t.Run("can start and stop server", func(t *testing.T) {
		sqlHelper := sqltest.NewHelper(t)

		httpRouterInjector := func(r *http.Router) {
			r.Mux.Get("/", func(w http2.ResponseWriter, r *http2.Request) {
				_, _ = w.Write([]byte("OK"))
			})
		}

		s := http.NewServer(http.NewServerOptions{
			Address:            ":58232",
			AdminPassword:      "correct horse battery staple",
			BaseURL:            "http://localhost:8080",
			HTTPRouterInjector: httpRouterInjector,
			SecureCookie:       false,
			SQLHelper:          sqlHelper,
		})

		go func() {
			is.NotError(t, s.Start())
		}()
		defer func() {
			is.NotError(t, s.Stop())
		}()

		// I know we could check that the server is running here, but it's easier to just wait a bit
		time.Sleep(10 * time.Millisecond)

		res, err := http2.Get("http://localhost:58232/")
		is.NotError(t, err)
		is.Equal(t, http2.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		is.NotError(t, err)
		is.Equal(t, "OK", string(body))
	})
}
