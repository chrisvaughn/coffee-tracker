package storage

import (
	"context"

	"cloud.google.com/go/datastore"
)

type User struct {
	Auth0ID string
	Key     *datastore.Key `datastore:"__key__"`
}

func (s *Storage) GetOrCreateUser(context context.Context, auth0ID string) (*User, error) {
	user, err := s.GetUserFromAuth0ID(context, auth0ID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	user = &User{
		Auth0ID: auth0ID,
	}
	err = s.CreateUser(context, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Storage) GetUserFromAuth0ID(context context.Context, auth0ID string) (*User, error) {
	var users []*User
	query := datastore.NewQuery("User").Filter("Auth0ID =", auth0ID).Limit(1)
	_, err := s.client.GetAll(context, query, &users)
	if len(users) == 1 {
		return users[0], err
	}
	return nil, err
}

func (s *Storage) CreateUser(context context.Context, u *User) error {
	newKey := datastore.IncompleteKey("User", nil)
	_, err := s.client.Put(context, newKey, u)
	return err
}
