package repository

import (
	"context"
	"encoding/json"
	"wb-internship-l0/internal/entity"
	postgres "wb-internship-l0/internal/lib/pg"
	"wb-internship-l0/internal/repository/pgdb"
)

// Order defines an interface for order-related operations.
type Order interface {
	AddOrder(ctx context.Context, id string, data json.RawMessage) error
	GetOrder(ctx context.Context, id string) (entity.Order, error)
	GetAllOrders(ctx context.Context) ([]entity.Order, error)
}

// Repositories is a struct that aggregates various repositories.
type Repositories struct {
	Order
}

// NewRepositories returns a new instance of Repository.
func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Order: pgdb.NewOrderRepository(pg),
	}
}
