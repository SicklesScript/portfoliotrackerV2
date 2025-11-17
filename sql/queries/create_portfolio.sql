-- name: CreatePortfolio :exec

INSERT INTO portfolios (id, created_at, updated_at, name, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
);