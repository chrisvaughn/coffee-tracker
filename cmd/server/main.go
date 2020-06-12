package main

import (
	"log"
	"net/http"

	"github.com/pseidemann/finish"

	"github.com/chrisvaughn/coffeetracker/pkg/coffeetracker"
	"github.com/chrisvaughn/coffeetracker/pkg/configuration"
)

func main() {

	cfg := configuration.GetConfiguration()
	service, err := coffeetracker.NewService()
	if err != nil {
		log.Fatal(err)
	}
	service.SetupRoutes()

	log.Printf("Accepting requests on port: %s", cfg.Port)
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
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
