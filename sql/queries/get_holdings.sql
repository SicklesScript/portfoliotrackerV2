-- name: GetHoldings :many
SELECT 
    p.name as portfolio_name,
    t.ticker,
    t.stock_name,
    SUM(t.shares)::TEXT as total_shares,
    AVG(t.price_per_share)::TEXT as avg_cost_per_share,
    SUM(t.shares * t.price_per_share)::TEXT as total_invested
FROM transactions t
JOIN portfolios p ON t.portfolio_id = p.id
WHERE p.user_id = $1
GROUP BY p.name, t.ticker, t.stock_name
ORDER BY p.name, t.ticker;