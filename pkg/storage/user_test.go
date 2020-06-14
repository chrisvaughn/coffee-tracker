package storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	assert := assert.New(t)

	s, err := NewStorage()
	assert.NotNil(s)
	assert.NoError(err)

	ctx := context.Background()

	auth0ID := "testing|" + uuid.New().String()

	// try to get user with ID, it should not exist
	user, err := s.GetUserFromAuth0ID(ctx, auth0ID)
	assert.Nil(user)
	assert.NoError(err)

	// create the user
	user = &User{
		Auth0ID: auth0ID,
	}
	err = s.CreateUser(ctx, user)
	assert.NoError(err)
	assert.NotNil(user.Key)
	assert.NotNil(user.UpdatedDT)
	assert.NotZero(user.ID)

	// try to get user with ID, it should not exist
	fetchedUser, err := s.GetUserFromAuth0ID(ctx, auth0ID)
	assert.NotNil(fetchedUser)
	assert.NoError(err)
	assert.Equal(user.Auth0ID, fetchedUser.Auth0ID)

	// delete user
	err = s.DeleteUser(ctx, user)
	assert.NoError(err)

	// try to get user with ID, it should not exist
	user, err = s.GetUserFromAuth0ID(ctx, auth0ID)
	assert.Nil(user)
	assert.NoError(err)
}

func TestUserGetOrCreate(t *testing.T) {
	assert := assert.New(t)

	s, err := NewStorage()
	assert.NotNil(s)
	assert.NoError(err)

	ctx := context.Background()

	auth0ID := "testing|" + uuid.New().String()

	// try to get user with ID, it will be created
	user, err := s.GetOrCreateUser(ctx, auth0ID)
	assert.NotNil(user)
	assert.NoError(err)

	// try to get user with ID, it should not exist
	fetchedUser, err := s.GetUserFromAuth0ID(ctx, auth0ID)
	assert.NotNil(fetchedUser)
	assert.NoError(err)
	assert.Equal(user.Auth0ID, fetchedUser.Auth0ID)

	// try to get user with ID, it will be returned
	fetchedUser, err = s.GetOrCreateUser(ctx, auth0ID)
	assert.NotNil(user)
	assert.NoError(err)
	assert.Equal(user.Auth0ID, fetchedUser.Auth0ID)

	// delete user
	err = s.DeleteUser(ctx, user)
	assert.NoError(err)

	// try to get user with ID, it should not exist
	user, err = s.GetUserFromAuth0ID(ctx, auth0ID)
	assert.Nil(user)
	assert.NoError(err)
}
