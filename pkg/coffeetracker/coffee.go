package coffeetracker

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/chrisvaughn/coffeetracker/pkg/httputils"
)

func (s *Service) getCoffee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	coffeeStr := chi.URLParam(r, "coffeeID")
	coffeeID, err := strconv.ParseInt(coffeeStr, 10, 64)
	if err != nil {
		httputils.ErrorResponse(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	coffee, err := s.storage.GetCoffee(ctx, coffeeID, 1)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}
	b, _ := json.Marshal(coffee)
	_, _ = w.Write(b)
}
