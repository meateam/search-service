package service

import (
	"context"
	// pb "github.com/meateam/search-service/proto"
)

// Controller is an interface for the business logic of the search.Service which uses a Store.
type Controller interface {
	HealthCheck(ctx context.Context) (bool, error)
}
