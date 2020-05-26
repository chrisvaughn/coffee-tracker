package storage

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"
)

type Storage struct {
	client *datastore.Client
}

func NewStorage() (*Storage, error) {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		return nil, err
	}

	s := &Storage{client}
	return s, nil
}
