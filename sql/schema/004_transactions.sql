-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE transactions (
    id UUID PRIMARY KEY NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    stock_name TEXT NOT NULL,
    ticker TEXT NOT NULL,
    shares DECIMAL(18, 6) NOT NULL,
    price_per_share DECIMAL(18, 6) NOT NULL,
    portfolio_id UUID NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE transactions;