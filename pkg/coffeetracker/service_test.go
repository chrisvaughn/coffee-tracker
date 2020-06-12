package coffeetracker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	assert := assert.New(t)

	s, err := NewService()
	assert.NoError(err)
	assert.NotNil(s)
	s.SetupRoutes()
}
