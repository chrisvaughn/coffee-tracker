package storage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCoffee(t *testing.T) {
	assert := assert.New(t)

	s, err := NewStorage()
	assert.NotNil(s)
	assert.NoError(err)

	ctx := context.Background()

	// create users
	auth0ID := "testing|" + uuid.New().String()
	user1, err := s.GetOrCreateUser(ctx, auth0ID)
	assert.NotNil(user1)
	assert.NoError(err)

	auth0ID = "testing|" + uuid.New().String()
	user2, err := s.GetOrCreateUser(ctx, auth0ID)
	assert.NotNil(user2)
	assert.NoError(err)

	// user1 should not have any coffees
	coffees, err := s.GetCoffeesByUser(ctx, user1)
	assert.Nil(coffees)
	assert.NoError(err)

	// create multiple coffees
	c1 := &Coffee{
		Name:  "Test Coffee1",
		Added: time.Now(),
	}
	err = s.CreateCoffee(ctx, c1, user1)
	assert.NoError(err)

	c2 := &Coffee{
		Name:  "Test Coffee2",
		Added: time.Now(),
	}
	s.CreateCoffee(ctx, c2, user1)
	assert.NoError(err)

	coffees, err = s.GetCoffeesByUser(ctx, user1)
	assert.Len(coffees, 2)
	assert.NoError(err)

	// user2 should not have any coffees
	coffees2, err := s.GetCoffeesByUser(ctx, user2)
	assert.Nil(coffees2)
	assert.NoError(err)
}
