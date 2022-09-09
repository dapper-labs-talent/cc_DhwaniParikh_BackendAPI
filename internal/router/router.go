package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"net/http"
	"server/internal/config"
	"server/internal/models"
)

func New(config *config.Config, logger *zerolog.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.GetHead)

	return r
}

func AddControllersToGroup(path string, g chi.Router, controllers []models.Controller) {

	var routes []models.Route
	for _, controller := range controllers {
		routes = append(routes, controller.Routes()...)
	}
	addRoutesToGroup(path, g, routes)
}

func addRoutesToGroup(path string, api chi.Router, routes []models.Route) {
	api.Route(path, func(r chi.Router) {
		for _, route := range routes {
			switch route.Method {
			case http.MethodPost:
				{
					r.Post(route.Path, route.HandlerWithMiddleware())
					break
				}
			case http.MethodGet:
				{
					r.Get(route.Path, route.HandlerWithMiddleware())
					break
				}
			case http.MethodDelete:
				{
					r.Delete(route.Path, route.HandlerWithMiddleware())
					break
				}
			case http.MethodPut:
				{
					r.Put(route.Path, route.HandlerWithMiddleware())
					break
				}
			case http.MethodPatch:
				{
					r.Patch(route.Path, route.HandlerWithMiddleware())
					break
				}
			}
		}
	})

}
