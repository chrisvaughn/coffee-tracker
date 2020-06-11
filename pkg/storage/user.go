package storage

import (
	"context"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type User struct {
	Auth0ID string
}

func (s *Storage) GetOrCreateUser(context context.Context, auth0ID string) (*datastore.Key, *User, error) {
	key, user, err := s.GetUserFromAuth0ID(context, auth0ID)
	if err != nil && err != iterator.Done {
		return nil, nil, err
	}
	if key != nil && user != nil {
		return key, user, nil
	}
	user = &User{
		Auth0ID: auth0ID,
	}
	key, err = s.CreateUser(context, user)
	if err != nil {
		return nil, nil, err
	}
	return key, user, nil
}

func (s *Storage) GetUserFromAuth0ID(context context.Context, auth0ID string) (*datastore.Key, *User, error) {
	user := &User{}
	query := datastore.NewQuery("User").Filter("Auth0ID =", auth0ID).Limit(1)
	it := s.client.Run(context, query)
	key, err := it.Next(user)
	return key, user, err
}

func (s *Storage) CreateUser(context context.Context, u *User) (*datastore.Key, error) {
	newKey := datastore.IncompleteKey("User", nil)
	return s.client.Put(context, newKey, u)
}
