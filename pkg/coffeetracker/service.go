package coffeetracker

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/chrisvaughn/coffeetracker/pkg/storage"
)

type Service struct {
	Router  *chi.Mux
	storage *storage.Storage
}

func NewService() (*Service, error) {
	strg, err := storage.NewStorage()
	if err != nil {
		return nil, err
	}
	svc := Service{
		Router:  chi.NewRouter(),
		storage: strg,
	}
	return &svc, nil
}

func (s *Service) SetupRoutes() {
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	s.Router.Use(cors.Handler)

	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.Router.Route("/api", func(r chi.Router) {
		r.Use(AuthMiddleware().Handler)
		r.Use(s.GetUserMiddleware)

		r.Get("/coffees", s.getCoffees)
		r.Post("/coffees", s.postCoffees)

		r.Get("/coffees/{coffeeID:[0-9]+}", s.getCoffee)
		r.Put("/coffees/{coffeeID:[0-9]+}", s.putCoffee)
		r.Delete("/coffees/{coffeeID:[0-9]+}", s.deleteCoffee)
	})
}
