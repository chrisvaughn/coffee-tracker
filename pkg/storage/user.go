package storage

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

type User struct {
	Auth0ID   string         `json:"-"`
	AddedDT   time.Time      `json:"added_dt"`
	UpdatedDT time.Time      `json:"updated_dt"`
	Key       *datastore.Key `datastore:"__key__" json:"-"`
	ID        int64          `datastore:"-" json:"id"`
}

func (x *User) LoadKey(k *datastore.Key) error {
	x.Key = k
	x.ID = x.Key.ID
	return nil
}

func (x *User) Load(ps []datastore.Property) error {
	if err := datastore.LoadStruct(x, ps); err != nil {
		return err
	}
	return nil
}

func (x *User) Save() ([]datastore.Property, error) {
	x.UpdatedDT = time.Now()
	return datastore.SaveStruct(x)
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
	key, err := s.client.Put(context, newKey, u)
	if err != nil {
		return err
	}
	u.Key = key
	u.ID = key.ID
	return nil
}

func (s *Storage) DeleteUser(context context.Context, u *User) error {
	return s.client.Delete(context, u.Key)
}
