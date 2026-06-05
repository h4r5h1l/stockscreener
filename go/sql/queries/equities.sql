-- name: UpsertEquityBase :exec
INSERT INTO equity_universe (conid, ticker, exchange)
VALUES ($1, $2, $3) ON CONFLICT (conid) DO
UPDATE
SET ticker = EXCLUDED.ticker,
    exchange = EXCLUDED.exchange;
-- name: UpdateEquityMetrics :exec
UPDATE equity_universe
SET pe_ratio = $1,
    price_to_sales = $2,
    return_on_equity = $3,
    business_summary = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE conid = $5;
-- name: ListEquitiesForSync :many
SELECT conid,
    ticker
FROM equity_universe;