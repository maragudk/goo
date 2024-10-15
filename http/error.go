package http

import (
	"net/http"

	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"

	"maragu.dev/goo/html"
)

type httpError struct {
	Code int
}

func (n httpError) Error() string {
	return http.StatusText(n.Code)
}

func (n httpError) StatusCode() int {
	return n.Code
}

func NotFound(page html.PageFunc) http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return html.NotFoundPage(page), httpError{Code: http.StatusNotFound}
	})
}
