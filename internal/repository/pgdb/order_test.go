package pgdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"testing"
	postgres "wb-internship-l0/internal/lib/pg"
	"wb-internship-l0/pkg/logger"
)

func TestOrderRepository_GetOrder_OK(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer mock.Close()

	log := logger.NewZap("dev")
	mockPg := &postgres.Postgres{
		Log: log,
		DB:  mock,
	}
	r := NewOrderRepository(mockPg)

	var id = "b563feb7b2b84b6test"
	var raw json.RawMessage = []byte(`{
		"order_uid": "b563feb7b2b84b6test",
		"track_number": "WBILMTESTTRACK",
		"entry": "WBIL",
		"delivery": {
			"name": "Test Testov",
			"phone": "+9720000000",
			"zip": "2639809",
			"city": "Kiryat Mozkin",
			"address": "Ploshad Mira 15",
			"region": "Kraiot",
			"email": "test@gmail.com"
		},
		"payment": {
			"transaction": "b563feb7b2b84b6test",
			"request_id": "",
			"currency": "USD",
			"provider": "wbpay",
			"amount": 1817,
			"payment_dt": 1637907727,
			"bank": "alpha",
			"delivery_cost": 1500,
			"goods_total": 317,
			"custom_fee": 0
		},
		"items": [
			{
				"chrt_id": 9934930,
				"track_number": "WBILMTESTTRACK",
				"price": 453,
				"rid": "ab4219087a764ae0btest",
				"name": "Mascaras",
				"sale": 30,
				"size": "0",
				"total_price": 317,
				"nm_id": 2389212,
				"brand": "Vivienne Sabo",
				"status": 202
			}
		],
		"locale": "en",
		"internal_signature": "",
		"customer_id": "test",
		"delivery_service": "meest",
		"shardkey": "9",
		"sm_id": 99,
		"date_created": "2021-11-26T06:22:19Z",
		"oof_shard": "1"
	}`)

	mock.ExpectExec("INSERT INTO orders_schema.order").
		WithArgs(id, raw).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err = r.AddOrder(context.Background(), id, raw); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOrderRepository_GetOrder_ErrAlreadyExists(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer mock.Close()

	log := logger.NewZap("dev")
	mockPg := &postgres.Postgres{
		Log: log,
		DB:  mock,
	}
	r := NewOrderRepository(mockPg)

	var id = "b563feb7b2b84b6test"
	var raw json.RawMessage = []byte(`{
		"order_uid": "b563feb7b2b84b6test",
		"track_number": "WBILMTESTTRACK",
		"entry": "WBIL",
		"delivery": {
			"name": "Test Testov",
			"phone": "+9720000000",
			"zip": "2639809",
			"city": "Kiryat Mozkin",
			"address": "Ploshad Mira 15",
			"region": "Kraiot",
			"email": "test@gmail.com"
		},
		"payment": {
			"transaction": "b563feb7b2b84b6test",
			"request_id": "",
			"currency": "USD",
			"provider": "wbpay",
			"amount": 1817,
			"payment_dt": 1637907727,
			"bank": "alpha",
			"delivery_cost": 1500,
			"goods_total": 317,
			"custom_fee": 0
		},
		"items": [
			{
				"chrt_id": 9934930,
				"track_number": "WBILMTESTTRACK",
				"price": 453,
				"rid": "ab4219087a764ae0btest",
				"name": "Mascaras",
				"sale": 30,
				"size": "0",
				"total_price": 317,
				"nm_id": 2389212,
				"brand": "Vivienne Sabo",
				"status": 202
			}
		],
		"locale": "en",
		"internal_signature": "",
		"customer_id": "test",
		"delivery_service": "meest",
		"shardkey": "9",
		"sm_id": 99,
		"date_created": "2021-11-26T06:22:19Z",
		"oof_shard": "1"
	}`)
	var expectedError = fmt.Errorf("%s: %w", "repository.order.AddOrder", ErrOrderAlreadyExists)

	mock.ExpectExec("INSERT INTO orders_schema.order").
		WithArgs(id, raw).
		WillReturnError(&pgconn.PgError{Code: "23505"})

	if err = r.AddOrder(context.Background(), id, raw); err == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	if errors.Is(err, expectedError) {
		t.Fatalf("expected error: %s, actual: %s", err, expectedError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOrderRepository_GetOrder_ErrUnexpected(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer mock.Close()

	log := logger.NewZap("dev")
	mockPg := &postgres.Postgres{
		Log: log,
		DB:  mock,
	}
	r := NewOrderRepository(mockPg)

	var id = ""
	var raw json.RawMessage = []byte("")

	mock.ExpectExec("INSERT INTO orders_schema.order").
		WithArgs(id, raw).
		WillReturnError(fmt.Errorf("unexpected error"))

	if err = r.AddOrder(context.Background(), id, raw); err == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
