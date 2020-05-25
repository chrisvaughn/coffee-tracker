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
	coffees, err := s.storage.GetCoffees(ctx)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}
	b, _ := json.Marshal(coffees)
	w.Write(b)
}

func (s *Service) postCoffees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)

	var coffee storage.Coffee
	err := decoder.Decode(&coffee)

	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 400)
	}
	coffee.Added = time.Now()
	fmt.Println(coffee.Name)

	err = s.storage.CreateCoffee(ctx, &coffee)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}
}
