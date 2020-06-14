package storage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCoffees(t *testing.T) {
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
	coffees, err := s.GetAllCoffeesForUser(ctx, user1)
	assert.Nil(coffees)
	assert.NoError(err)

	// create multiple coffees
	c1 := &Coffee{
		Name:    "Test Coffee1",
		AddedDT: time.Now(),
	}
	err = s.CreateCoffee(ctx, c1, user1)
	assert.NoError(err)
	assert.NotNil(c1.UpdatedDT)
	assert.NotNil(c1.Key)
	assert.NotZero(c1.ID)

	c2 := &Coffee{
		Name:    "Test Coffee2",
		AddedDT: time.Now(),
	}
	err = s.CreateCoffee(ctx, c2, user1)
	assert.NoError(err)

	coffees, err = s.GetAllCoffeesForUser(ctx, user1)
	assert.Len(coffees, 2)
	assert.NoError(err)

	c, err := s.GetCoffeeByID(ctx, coffees[0].ID, user1)
	assert.NoError(err)
	assert.Equal(c1.Name, c.Name)

	// user2 should not have any coffees
	coffees2, err := s.GetAllCoffeesForUser(ctx, user2)
	assert.Nil(coffees2)
	assert.NoError(err)
}

func TestCoffeeCRUD(t *testing.T) {
	assert := assert.New(t)

	s, err := NewStorage()
	assert.NotNil(s)
	assert.NoError(err)

	ctx := context.Background()
	auth0ID := "testing|" + uuid.New().String()
	user, err := s.GetOrCreateUser(ctx, auth0ID)
	assert.NotNil(user)
	assert.NoError(err)

	// create coffee
	c1 := &Coffee{
		Name:    "Test Coffee CRUD",
		AddedDT: time.Now(),
	}
	err = s.CreateCoffee(ctx, c1, user)
	assert.NoError(err)

	// update coffee
	c1.Name = "Test Coffee Updated"
	err = s.UpdateCoffee(ctx, c1)
	assert.NoError(err)
	updatedDT := c1.UpdatedDT

	c1.Name = "Test Coffee Updated again"
	err = s.UpdateCoffee(ctx, c1)
	assert.NoError(err)
	assert.Greater(c1.UpdatedDT.UnixNano(), updatedDT.UnixNano())

	// delete coffee
	err = s.DeleteCoffee(ctx, c1)
	assert.NoError(err)

	c2, err := s.GetCoffeeByID(ctx, c1.ID, user)
	assert.NoError(err)
	assert.Nil(c2)
}
