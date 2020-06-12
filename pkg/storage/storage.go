package storage

import (
	"context"

	"cloud.google.com/go/datastore"

	"github.com/chrisvaughn/coffeetracker/pkg/configuration"
)

type Storage struct {
	client *datastore.Client
}

func NewStorage() (*Storage, error) {
	ctx := context.Background()

	cfg := configuration.GetConfiguration()

	client, err := datastore.NewClient(ctx, cfg.GoogleCloudProject)
	if err != nil {
		return nil, err
	}

	s := &Storage{client}
	return s, nil
}
