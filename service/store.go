package service

import (
	"context"
)

// Store is an interface for handling the storing of permissions.
type Store interface {
	Create(ctx context.Context, permission Search) (Search, error)
	Get(ctx context.Context, filter interface{}) (Search, error)
	GetAll(ctx context.Context, filter interface{}) ([]Search, error)
	Delete(ctx context.Context, filter interface{}) (Search, error)
	HealthCheck(ctx context.Context) (bool, error)
}
