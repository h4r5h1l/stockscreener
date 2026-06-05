-- +goose Up
CREATE TABLE equity_universe (
    conid INT PRIMARY KEY,
    ticker VARCHAR(12) NOT NULL,
    exchange VARCHAR(12) NOT NULL,
    pe_ratio NUMERIC(12, 4),
    price_to_sales NUMERIC(12, 4),
    return_on_equity NUMERIC(12, 4),
    business_summary TEXT,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_equity_universe_ticker ON equity_universe(ticker);
-- +goose Down
DROP TABLE equity_universe;