package coffeetracker

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/chrisvaughn/coffeetracker/pkg/httputils"
	"github.com/chrisvaughn/coffeetracker/pkg/storage"
)

func (s *Service) getCoffee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	coffeeStr := chi.URLParam(r, "coffeeID")
	coffeeID, err := strconv.ParseInt(coffeeStr, 10, 64)
	if err != nil {
		httputils.ErrorResponse(w, "invalid coffee_id", http.StatusBadRequest)
		return
	}

	user := ctx.Value(AuthContextUser).(*storage.User)
	if user == nil {
		httputils.ErrorResponse(w, "did not get user", 500)
	}

	coffee, err := s.storage.GetCoffeeByID(ctx, coffeeID, user)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}
	b, _ := json.Marshal(coffee)
	_, _ = w.Write(b)
}
