package elasticsearch

import (
	"context"

	es "github.com/elastic/go-elasticsearch/v7"
)

// Store holds the elasticsearch and implements Store interface.
type Store struct {
	client *es.Client
}

func newStore(cfg es.Config) (*Store, error) {
	client, err := es.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Store{client: client}, nil
}

// HealthCheck checks the health of the service, returns true if healthy, or false otherwise.
func (s Store) HealthCheck(ctx context.Context) (bool, error) {
	if _, err := s.client.Info(s.client.API.Info.WithContext(ctx)); err != nil {
		return false, err
	}

	return true, nil
}
