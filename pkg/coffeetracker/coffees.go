package coffeetracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chrisvaughn/coffeetracker/pkg/httputils"
	"github.com/chrisvaughn/coffeetracker/pkg/storage"
)

func (s *Service) getCoffees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth0UserID := ctx.Value(AuthContextUserID).(string)
	if auth0UserID == "" {
		httputils.ErrorResponse(w, "did not get user id", 500)
	}
	fmt.Printf("%s\n", auth0UserID)
	key, _, err := s.storage.GetOrCreateUser(ctx, auth0UserID)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}
	coffees, err := s.storage.GetCoffeesByUser(ctx, key)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}
	b, _ := json.Marshal(coffees)
	_, _ = w.Write(b)
}

func (s *Service) postCoffees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	auth0UserID := ctx.Value(AuthContextUserID).(string)
	if auth0UserID == "" {
		httputils.ErrorResponse(w, "did not get user id", 500)
	}
	fmt.Printf("%s\n", auth0UserID)
	key, _, err := s.storage.GetOrCreateUser(ctx, auth0UserID)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}

	decoder := json.NewDecoder(r.Body)

	var coffee storage.Coffee
	err = decoder.Decode(&coffee)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 400)
	}

	coffee.Added = time.Now()
	fmt.Println(coffee.Name)

	err = s.storage.CreateCoffee(ctx, &coffee, key)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}
}
