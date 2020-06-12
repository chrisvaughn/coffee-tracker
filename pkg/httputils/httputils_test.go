package httputils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	ErrorResponse(w, "test", http.StatusBadRequest)
	assert.Equal(http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(err)
	assert.Contains(resp, "error")
	assert.Equal("test", resp["error"].(string))
}
