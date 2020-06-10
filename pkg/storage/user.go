package storage

import (
	"context"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type User struct {
	Auth0ID string
}

func (s *Storage) GetOrCreateUser(context context.Context, auth0ID string) (*User, error) {
	user, err := s.GetUserFromAuth0ID(context, auth0ID)
	if err != nil && err != iterator.Done {
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
	user := &User{}
	query := datastore.NewQuery("User").
		Filter("Auth0ID =", auth0ID)
	query.Limit(1)
	it := s.client.Run(context, query)
	_, err := it.Next(user)
	return user, err
}

func (s *Storage) CreateUser(context context.Context, u *User) error {
	newKey := datastore.IncompleteKey("User", nil)
	_, err := s.client.Put(context, newKey, u)
	return err
}
