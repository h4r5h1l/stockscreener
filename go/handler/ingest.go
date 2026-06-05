package handler

import (
	"context"       //required for database operations
	"encoding/json" // decode JSON responses from Python API
	"log"           //print errors
	"net/http"      ///call Python API

	"github.com/h4r5h1l/stockscreener/internal/db" // SQLC-generated database package
	"github.com/jackc/pgx/v5/pgxpool"              // PostgreSQL connection pool
)

type Equity struct {
	ConID    int    `json:"conid"`
	Ticker   string `json:"ticker"`
	Exchange string `json:"exchange"`
}

// Connects to your PostgreSQL database
func IngestEquities() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://localhost:5432/stockscreener?sslmode=disable")
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	q := db.New(pool)

	// Calls your Python API (/all-equities)
	resp, err := http.Get("http://localhost:5010")
	if err != nil {
		log.Fatal("Error calling Python API:", err)
	}
	defer resp.Body.Close()

	// Parses the JSON list of equities
	var equities []Equity
	if err := json.NewDecoder(resp.Body).Decode(&equities); err != nil {
		log.Fatal("Error decoding JSON response:", err)
	}
	// Inserts each equity into your equities table using SQLC
	for _, eq := range equities {
		err := q.InsertEquity(ctx, db.InsertEquityParams{
			Conid:    int32(eq.ConID),
			Ticker:   eq.Ticker,
			Exchange: eq.Exchange,
		})
		if err != nil {
			log.Printf("Error inserting equity %s: %v", eq.Ticker, err)
		}
	}
	log.Printf("Successfully ingested %d equities", len(equities))
}
