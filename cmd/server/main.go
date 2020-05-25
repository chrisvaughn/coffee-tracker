package main

import (
	"log"
	"net/http"
	"os"

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
	if err := http.ListenAndServe(":"+port, service.Router); err != nil {
		log.Fatal(err)
	}
}
