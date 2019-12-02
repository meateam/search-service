package service

import (
	"context"
	"time"

	pb "github.com/meateam/search-service/proto"
	"github.com/sirupsen/logrus"
)

// Service is a structure used for handling Search Service grpc requests.
type Service struct {
	logger     *logrus.Logger
	controller Controller
}

// HealthCheck checks the health of the service, returns true if healthy, or false otherwise.
func (s Service) HealthCheck() bool {
	timeoutCtx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()
	healthy, err := s.controller.HealthCheck(timeoutCtx)
	if err != nil {
		s.logger.Errorf("%v", err)
		return false
	}

	return healthy
}

// NewService creates a Service and returns it.
func NewService(controller Controller, logger *logrus.Logger) Service {
	return Service{controller: controller, logger: logger}
}

// CreateFile is the request handler for creating a file.
func (s Service) CreateFile(ctx context.Context, req *pb.File) (*pb.CreateFileResponse, error) {
	return s.controller.CreateFile(ctx, req)
}

// Search is the request handler for searching a file.
func (s Service) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	return s.controller.Search(ctx, req)
}

// Delete is the request handler for deleting a file.
func (s Service) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return s.controller.Delete(ctx, req)
}

// Update is the request handler for updating a file.
func (s Service) Update(ctx context.Context, req *pb.File) (*pb.UpdateResponse, error) {
	return s.controller.Update(ctx, req)
}
