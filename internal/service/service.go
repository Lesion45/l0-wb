package service

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"

	"wb-internship-l0/internal/repository"
	"wb-internship-l0/pkg/cache"
)

type Order interface {
	GetOrder(ctx context.Context, id string) (json.RawMessage, error)
	SaveOrder(ctx context.Context, id string, data json.RawMessage) error
}

type Services struct {
	Order
}

type ServicesDependencies struct {
	Log   *zap.Logger
	Cache cache.Cache
	Repos *repository.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Order: NewOrderService(deps.Log, deps.Cache, deps.Repos.Order),
	}
}
