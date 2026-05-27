package protocol

import (
	"net/http"
	"organization_structure/internal/adapters/protocol/handlers"
	"organization_structure/internal/usecase"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-openapi/runtime/server-middleware/docui"
)

func NewRouter(s *usecase.UseCase) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/healpth"))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PATCH"},
	}))

	r.Route("/departments", func(r chi.Router) {
		r.Post("/", handlers.HCreateDepartment(s))
		r.Post("/{id}/employees", handlers.HCreateEmployee(s))
		r.Get("/{id}", handlers.HGetDepartment(s))
		r.Patch("/{id}", handlers.HUpdateDepartment(s))
		r.Delete("/{id}", handlers.HDeleteDepartment(s))
	})

	swaggerUIHandler := docui.SwaggerUI(r, docui.WithSpecURL("/openapi.yaml"))
	r.Handle("/docs", swaggerUIHandler)
	r.Handle("/openapi.yaml", http.FileServer(http.Dir("././openapi/")))
	r.Handle("/common.yaml", http.FileServer(http.Dir("././openapi/")))

	return r
}
