package service

import (
	"context"
	pb "github.com/meateam/search-service/proto"
)

// Controller is an interface for the business logic of the search.Service which uses a Store.
type Controller interface {
	CreateFile(ctx context.Context, req *pb.File) (*pb.CreateFileResponse, error)
	Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error)
	Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error)
	Update(ctx context.Context, req *pb.File) (*pb.UpdateResponse, error)
	HealthCheck(ctx context.Context) (bool, error)
}
