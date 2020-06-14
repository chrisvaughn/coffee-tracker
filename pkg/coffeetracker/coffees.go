package coffeetracker

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chrisvaughn/coffeetracker/pkg/httputils"
	"github.com/chrisvaughn/coffeetracker/pkg/storage"
)

func (s *Service) getCoffees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value(AuthContextUser).(*storage.User)
	if user == nil {
		httputils.ErrorResponse(w, "did not get user", 500)
	}

	coffees, err := s.storage.GetAllCoffeesForUser(ctx, user)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}
	b, _ := json.Marshal(coffees)
	_, _ = w.Write(b)
}

func (s *Service) postCoffees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := ctx.Value(AuthContextUser).(*storage.User)
	if user == nil {
		httputils.ErrorResponse(w, "did not get user", 500)
	}

	decoder := json.NewDecoder(r.Body)

	var coffee storage.Coffee
	err := decoder.Decode(&coffee)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 400)
	}

	coffee.AddedDT = time.Now()
	err = s.storage.CreateCoffee(ctx, &coffee, user)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}
	w.WriteHeader(201)
}
