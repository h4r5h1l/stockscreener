package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	// Update imports
	"github.com/h4r5h1l/stockscreener/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	// 1. Connect using pgxpool
	// Use the project's migration DB string (no embedded credentials, disable TLS for local dev)
	connStr := "postgres://localhost:5432/stockscreener?sslmode=disable"
	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Println("Database Connection Error:", err)
		return
	}
	defer dbpool.Close()

	// sqlc.New(dbpool) works perfectly because pgxpool implements DBTX
	queries := db.New(dbpool)

	fmt.Println("Running Python from Go...!")

	// 2. Execute Python
	cmd := exec.Command("uv", "run", "--directory", "../ibpython", "main.py", "fetch_universe")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Runtime Execution Error: %v\nOutput: %s\n", err, string(out))
	}

	var equities []struct {
		Ticker   string `json:"ticker"`
		Conid    int32  `json:"conid"`
		Exchange string `json:"exchange"`
	}

	if err := json.Unmarshal(out, &equities); err != nil {
		fmt.Println("JSON Unmarshal Error:", err)
		return
	}

	// 3. Populate Database
	for _, eq := range equities {
		err := queries.UpsertEquityBase(ctx, db.UpsertEquityBaseParams{
			Conid:    eq.Conid,
			Ticker:   eq.Ticker,
			Exchange: eq.Exchange,
		})
		if err != nil {
			fmt.Printf("Database Upsert Error for %s: %v\n", eq.Ticker, err)
		}
	}

	fmt.Println("Table populated successfully!")
}
