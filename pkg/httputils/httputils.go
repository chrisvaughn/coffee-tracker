package httputils

import (
	"encoding/json"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"error": error,
	}
	b, _ := json.Marshal(resp)
	_, _ = w.Write(b)
}
