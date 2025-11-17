-- name: CreateTransaction :exec

INSERT INTO transactions (id, date, stock_name, ticker, shares, price_per_share, portfolio_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5
);