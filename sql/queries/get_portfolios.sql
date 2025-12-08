-- name: GetPortfoliosForUser :many
SELECT id, name FROM portfolios WHERE user_id = $1;