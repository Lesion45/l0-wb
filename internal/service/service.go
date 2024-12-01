package service

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"

	"wb-internship-l0/internal/repository"
	"wb-internship-l0/pkg/cache"
)

// Order defines the interface for managing orders.
type Order interface {
	GetOrder(ctx context.Context, id string) (json.RawMessage, error)
	SaveOrder(ctx context.Context, id string, data json.RawMessage) error
}

// Services aggregates all application service interfaces.
type Services struct {
	Order
}

// ServicesDependencies holds dependencies required to create services.
type ServicesDependencies struct {
	Log   *zap.Logger
	Cache cache.Cache
	Repos *repository.Repositories
}

// NewServices initializes and returns a Services struct with all dependencies resolved.
func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Order: NewOrderService(deps.Log, deps.Cache, deps.Repos.Order),
	}
}
