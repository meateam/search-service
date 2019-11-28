package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

// Service is a structure used for handling Permission Service grpc requests.
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
