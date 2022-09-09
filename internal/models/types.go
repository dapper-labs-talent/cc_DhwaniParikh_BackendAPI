package models

import (
	"net/http"
)

type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []MiddlewareFunc
}

type Controller interface {
	Routes() []Route
}

func (r *Route) HandlerWithMiddleware() http.HandlerFunc {
	handler := r.Handler
	for _, m := range r.Middleware {
		handler = m(handler)
	}
	return handler
}
