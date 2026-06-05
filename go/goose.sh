#!/usr/bin/env bash
set -a
GOOSE_DRIVER=postgres
GOOSE_DBSTRING="postgres://localhost:5432/stockscreener?sslmode=disable"
GOOSE_MIGRATION_DIR="$(pwd)/sql/schema"
set +a

goose "$@"