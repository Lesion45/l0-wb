-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS orders_schema;

CREATE TABLE IF NOT EXISTS orders_schema.order(
   OrderID VARCHAR(255) PRIMARY KEY,
   Data JSONB
);
-- +goose StatementEnd