-- +goose Up
CREATE TABLE equity_universe (
    conid INT PRIMARY KEY,
    ticker VARCHAR(12) NOT NULL,
    exchange VARCHAR(12) NOT NULL,
    pe_ratio DOUBLE PRECISION,
    price_to_sales DOUBLE PRECISION,
    return_on_equity DOUBLE PRECISION,
    business_summary TEXT,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_equity_universe_ticker ON equity_universe(ticker);
-- +goose Down
DROP TABLE equity_universe;