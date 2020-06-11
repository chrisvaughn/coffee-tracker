package storage

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

type Coffee struct {
	Name  string
	Added time.Time
	Key   *datastore.Key `datastore:"__key__"`
}

func (s *Storage) GetCoffee(context context.Context, coffeeID int64, user *User) (*Coffee, error) {
	coffee := &Coffee{}
	key := datastore.IDKey("Coffee", coffeeID, user.Key)
	err := s.client.Get(context, key, coffee)
	return coffee, err
}

func (s *Storage) GetCoffeesByUser(context context.Context, user *User) ([]*Coffee, error) {
	var coffees []*Coffee
	qry := datastore.NewQuery("Coffee").Ancestor(user.Key)
	_, err := s.client.GetAll(context, qry, &coffees)
	return coffees, err
}

func (s *Storage) CreateCoffee(context context.Context, c *Coffee, user *User) error {
	newKey := datastore.IncompleteKey("Coffee", user.Key)
	_, err := s.client.Put(context, newKey, c)
	return err
}
