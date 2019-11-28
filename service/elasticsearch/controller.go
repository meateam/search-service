package elasticsearch

import (
	"context"

	es "github.com/elastic/go-elasticsearch/v7"
)

// Controller is the search service business logic implementation using elasticsearch store.
type Controller struct {
	store *Store
}

// NewController returns a new controller.
func NewController(cfg es.Config) (*Controller, error) {
	store, err := newStore(cfg)
	if err != nil {
		return nil, err
	}

	return &Controller{store: store}, nil
}

// HealthCheck runs store's healthcheck and returns true if healthy, otherwise returns false
// and any error if occured.
func (c Controller) HealthCheck(ctx context.Context) (bool, error) {
	return c.store.HealthCheck(ctx)
}
