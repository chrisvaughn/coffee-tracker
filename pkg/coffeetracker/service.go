package coffeetracker

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Service struct {
	Router *chi.Mux
}

func NewService() (*Service, error) {
	svc := Service{
		Router: chi.NewRouter(),
	}
	return &svc, nil
}

func (s *Service) SetupRoutes() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "views"))
	FileServer(s.Router, "/", filesDir)

	s.Router.Route("/api", func(r chi.Router) {
		r.Get("/coffees", s.dummyHandler)
		r.Post("/coffees", s.dummyHandler)

		r.Get("/coffees/{coffeeID}", s.dummyHandler)
		r.Put("/coffees/{coffeeID}", s.dummyHandler)
		r.Delete("/coffees/{coffeeID}", s.dummyHandler)
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
