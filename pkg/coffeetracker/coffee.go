package coffeetracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"github.com/chrisvaughn/coffeetracker/pkg/httputils"
	"github.com/chrisvaughn/coffeetracker/pkg/storage"
)

func (s *Service) getCoffee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	coffeeStr := chi.URLParam(r, "coffeeID")
	coffeeID, err := strconv.ParseInt(coffeeStr, 10, 64)
	if err != nil {
		httputils.ErrorResponse(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	coffee, err := s.storage.GetCoffee(ctx, coffeeID)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}
	b, _ := json.Marshal(coffee)
	w.Write(b)
}

func (s *Service) putCoffee(w http.ResponseWriter, r *http.Request) {
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
