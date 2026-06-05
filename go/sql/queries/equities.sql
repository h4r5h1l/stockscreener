-- name: InsertEquity :exec
INSERT INTO equities (conid, ticker, name, exchange)
VALUES ($1, $2, $3, $4);
-- name: GetEquityByConid :one
SELECT *
FROM equities
WHERE conid = $1;
-- name: ListEquities :many
SELECT *
FROM equities
ORDER BY ticker;