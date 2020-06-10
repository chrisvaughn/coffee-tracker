package coffeetracker

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	storage, err := storage.NewStorage()
	if err != nil {
		return nil, err
	}
	svc := Service{
		Router:  chi.NewRouter(),
		storage: storage,
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

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "views"))
	FileServer(s.Router, "/", filesDir)

	s.Router.Route("/api", func(r chi.Router) {
		r.Use(AuthMiddleware().Handler)
		r.Use(GetUserID)
		r.Get("/coffees", s.getCoffees)
		r.Post("/coffees", s.postCoffees)

		r.Get("/coffees/{coffeeID:[0-9]+}", s.getCoffee)
		r.Put("/coffees/{coffeeID:[0-9]+}", s.dummyHandler)
		r.Delete("/coffees/{coffeeID:[0-9]+}", s.dummyHandler)
	})
}

func (s *Service) dummyHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Made it")
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
