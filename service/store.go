package service

import (
	"context"
	pb "github.com/meateam/search-service/proto"
)

// Store is an interface for handling the storing of files.
type Store interface {
	Create(ctx context.Context, file *pb.File) (string, error)
	GetAll(ctx context.Context, filter interface{}) ([]string, error)
	Delete(ctx context.Context, filter interface{}) (*pb.File, error)
	HealthCheck(ctx context.Context) (bool, error)
}
