package storage

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

type Coffee struct {
	Name  string
	Added time.Time
}

type CoffeeWithID struct {
	*Coffee
	ID int64
}

func (s *Storage) GetCoffee(context context.Context, coffeeID int64, userKey *datastore.Key) (*Coffee, error) {
	coffee := &Coffee{}
	key := datastore.IDKey("Coffee", coffeeID, userKey)
	err := s.client.Get(context, key, coffee)
	return coffee, err
}

func (s *Storage) GetCoffeesByUser(context context.Context, userKey *datastore.Key) ([]*CoffeeWithID, error) {
	var coffees []*Coffee
	qry := datastore.NewQuery("Coffee").Ancestor(userKey)
	keys, err := s.client.GetAll(context, qry, &coffees)

	var coffeesWithID []*CoffeeWithID
	for i, key := range keys {
		coffeesWithID = append(coffeesWithID, &CoffeeWithID{coffees[i], key.ID})
	}
	return coffeesWithID, err
}

func (s *Storage) CreateCoffee(context context.Context, c *Coffee, userKey *datastore.Key) error {
	newKey := datastore.IncompleteKey("Coffee", userKey)
	_, err := s.client.Put(context, newKey, c)
	return err
}
