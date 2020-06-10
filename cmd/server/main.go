package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pseidemann/finish"

	"github.com/chrisvaughn/coffeetracker/pkg/coffeetracker"
)

func main() {

	service, err := coffeetracker.NewService()
	if err != nil {
		log.Fatal(err)
	}
	service.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: service.Router,
	}

	fin := finish.New()
	fin.Add(srv)

	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	fin.Wait()
}
