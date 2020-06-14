package storage

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/datastore"
)

type Coffee struct {
	Name      string         `json:"name"`
	AddedDT   time.Time      `json:"added_dt"`
	UpdatedDT time.Time      `json:"updated_dt"`
	Key       *datastore.Key `datastore:"__key__" json:"-"`
	ID        int64          `datastore:"-" json:"id"`
}

func (x *Coffee) LoadKey(k *datastore.Key) error {
	x.Key = k
	x.ID = x.Key.ID
	return nil
}

func (x *Coffee) Load(ps []datastore.Property) error {
	if err := datastore.LoadStruct(x, ps); err != nil {
		return err
	}
	return nil
}

func (x *Coffee) Save() ([]datastore.Property, error) {
	x.UpdatedDT = time.Now()
	return datastore.SaveStruct(x)
}

func (s *Storage) GetCoffeeByID(context context.Context, coffeeID int64, user *User) (*Coffee, error) {
	coffee := &Coffee{}
	key := datastore.IDKey("Coffee", coffeeID, user.Key)
	err := s.client.Get(context, key, coffee)
	if err != nil && errors.Is(err, datastore.ErrNoSuchEntity) {
		return nil, nil
	}
	return coffee, err
}

func (s *Storage) GetAllCoffeesForUser(context context.Context, user *User) ([]*Coffee, error) {
	var coffees []*Coffee
	qry := datastore.NewQuery("Coffee").Ancestor(user.Key)
	_, err := s.client.GetAll(context, qry, &coffees)
	return coffees, err
}

func (s *Storage) CreateCoffee(context context.Context, c *Coffee, user *User) error {
	newKey := datastore.IncompleteKey("Coffee", user.Key)
	key, err := s.client.Put(context, newKey, c)
	if err != nil {
		return err
	}
	c.Key = key
	c.ID = key.ID
	return nil
}

func (s *Storage) UpdateCoffee(context context.Context, c *Coffee) error {
	_, err := s.client.Put(context, c.Key, c)
	return err
}

func (s *Storage) DeleteCoffee(context context.Context, c *Coffee) error {
	return s.client.Delete(context, c.Key)
}
