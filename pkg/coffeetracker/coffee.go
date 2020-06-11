package coffeetracker

import (
	"encoding/json"
	"fmt"
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
		httputils.ErrorResponse(w, "invalid coffee_id", http.StatusBadRequest)
		return
	}

	auth0UserID := ctx.Value(AuthContextUserID).(string)
	if auth0UserID == "" {
		httputils.ErrorResponse(w, "did not get user id", 500)
	}
	fmt.Printf("%s\n", auth0UserID)
	user, err := s.storage.GetOrCreateUser(ctx, auth0UserID)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}

	coffee, err := s.storage.GetCoffee(ctx, coffeeID, user)
	if err != nil {
		httputils.ErrorResponse(w, err.Error(), 500)
	}
	b, _ := json.Marshal(coffee)
	_, _ = w.Write(b)
}
