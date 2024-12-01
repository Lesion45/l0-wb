package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"wb-internship-l0/internal/repository"
	"wb-internship-l0/internal/repository/pgdb"
	"wb-internship-l0/pkg/cache"
)

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotFound      = errors.New("order not found")
	ErrCache              = errors.New("failed to save order in cache")
)

// OrderService provides methods to manage orders.
type OrderService struct {
	Log   *zap.Logger
	Cache cache.Cache
	Repo  repository.Order
}

// NewOrderService initializes and returns a new OrderService.
func NewOrderService(log *zap.Logger, cache cache.Cache, repo repository.Order) *OrderService {
	return &OrderService{
		Log:   log,
		Cache: cache,
		Repo:  repo,
	}
}

// SaveOrder saves a new order to the database and cache.
func (s *OrderService) SaveOrder(ctx context.Context, id string, data json.RawMessage) error {
	const op = "service.OrderService.SaveOrder"
	s.Log.With(
		zap.String("op", op),
	)

	s.Log.Info("Attempting to save order")

	err := s.Repo.AddOrder(ctx, id, data)
	if err != nil {
		if errors.Is(err, pgdb.ErrOrderAlreadyExists) {
			s.Log.Warn("Order already exists",
				zap.String("orderID", id),
			)

			return fmt.Errorf("%s: %w", op, ErrOrderAlreadyExists)
		}

		s.Log.Error("Failed to save order to database",
			zap.String("orderID", id),
			zap.Error(err),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	s.Log.Info("Order successfully saved to database")

	err = s.Cache.Set(id, data)
	if err != nil {
		s.Log.Warn("Failed to save order to cache",
			zap.String("orderID", id),
			zap.Error(err),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	s.Log.Info("Order successfully saved to cache")

	return nil
}

// GetOrder retrieves an order by its ID, checking the cache first.
func (s *OrderService) GetOrder(ctx context.Context, id string) (json.RawMessage, error) {
	const op = "service.OrderService.GetOrder"
	s.Log.With(
		zap.String("op", op),
	)

	s.Log.Info("Attempting to found order")

	dataFromCache, found := s.Cache.Get(id)
	if found {
		return dataFromCache.(json.RawMessage), nil
	}

	order, err := s.Repo.GetOrder(ctx, id)
	if err != nil {
		if errors.Is(err, pgdb.ErrOrderNotFound) {
			s.Log.Warn("Order not found",
				zap.String("orderID", id),
				zap.Error(err),
			)

			return nil, fmt.Errorf("%s: %w", op, ErrOrderNotFound)
		}

		s.Log.Error("Failed to get order",
			zap.Error(err),
		)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.Log.Info("Order successfully found")

	return order.Data, nil
}
