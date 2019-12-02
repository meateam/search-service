package elasticsearch

import (
	"context"
	"fmt"

	pb "github.com/meateam/search-service/proto"
	es "github.com/olivere/elastic/v7"
)

// Controller is the search service business logic implementation using elasticsearch store.
type Controller struct {
	store *Store
}

// NewController returns a new controller.
func NewController(cfg []es.ClientOptionFunc, index string) (*Controller, error) {
	store, err := newStore(cfg, index)
	if err != nil {
		return nil, err
	}

	return &Controller{store: store}, nil
}

// HealthCheck runs store's healthcheck and returns true if healthy, otherwise returns false
// and any error if occurred.
func (c Controller) HealthCheck(ctx context.Context) (bool, error) {
	return c.store.HealthCheck(ctx)
}

// CreateFile creates a file in store and returns its unique ID.
func (c Controller) CreateFile(ctx context.Context, req *pb.File) (*pb.CreateFileResponse, error) {
	id, err := c.store.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.CreateFileResponse{Id: id}, nil
}

// Search retrieves a list of the file ids that match the search term, and any error if occurred.
func (c Controller) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	query := es.NewMultiMatchQuery(req.GetTerm())
	ids, err := c.store.GetAll(ctx, query)
	if err != nil {
		return nil, err
	}

	return &pb.SearchResponse{Ids: ids}, nil
}

// Delete retrieves a file id and id the match file by fild id from store, and any error if occurred.
func (c Controller) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	id := req.GetId()
	if id == "" {
		return nil, fmt.Errorf("file id is required")
	}

	res, err := c.store.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{Id: res}, nil
}

// Update retrieves a file and update the match file id, and any error if occurred.
func (c Controller) Update(ctx context.Context, req *pb.File) (*pb.UpdateResponse, error) {
	id := req.GetId()
	if id == "" {
		return nil, fmt.Errorf("file id is required")
	}

	res, err := c.store.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateResponse{Id: res}, nil
}
