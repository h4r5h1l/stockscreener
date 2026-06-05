-- +goose Up
CREATE TABLE equities (
    conid INTEGER PRIMARY KEY,
    ticker VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    exchange VARCHAR(50) NOT NULL
);
-- +goose Down
DROP TABLE equities;