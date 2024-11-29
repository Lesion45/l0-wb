package pgdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"wb-internship-l0/internal/entity"
	postgres "wb-internship-l0/internal/lib/pg"
	"wb-internship-l0/internal/repository"
)

// OrderRepository is a repository for managing orders in the database.
type OrderRepository struct {
	*postgres.Postgres
}

// NewOrderRepository creates a new instance of OrderRepository.
func NewOrderRepository(pg *postgres.Postgres) *OrderRepository {
	return &OrderRepository{pg}
}

// AddOrder adds a new order to the database.
// Returns an error if the insertion fails
func (r *OrderRepository) AddOrder(ctx context.Context, id string, data json.RawMessage) error {
	const op = "repository.order.AddOrder"

	query := `INSERT INTO orders_schema.order(OrderID, Data) VALUES(@id, @data)`
	args := pgx.NamedArgs{
		"id":   id,
		"data": data,
	}

	_, err := r.DB.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, repository.ErrOrderAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetOrder retrieves an order by its uid from the database.
// Returns entity.Order if the order is found, or an error if errors are occurred.
func (r *OrderRepository) GetOrder(ctx context.Context, id string) (entity.Order, error) {
	const op = "repository.order.GetOrder"

	var data json.RawMessage

	query := `SELECT OrderID, Data FROM orders_schema.order WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	err := r.DB.QueryRow(ctx, query, args).Scan(&data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{}, fmt.Errorf("%s: %w", op, repository.ErrOrderNotFound)
		}

		return entity.Order{}, fmt.Errorf("%s: %w", op, err)
	}

	order := entity.Order{
		UID:  id,
		Data: data,
	}

	return order, nil
}
