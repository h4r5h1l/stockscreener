package main

import (
	"bufio"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/h4r5h1l/stockscreener/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	connStr := "postgres://localhost:5432/stockscreener?sslmode=disable"
	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal("DB Error:", err)
	}
	defer dbpool.Close()
	queries := db.New(dbpool)

	// PHASE 1: Fetch Universe
	fmt.Println("Syncing Universe...")
	cmd := exec.Command("uv", "run", "--directory", "../ibpython", "main.py", "fetch_universe")
	out, _ := cmd.CombinedOutput()

	var equities []struct {
		Ticker string `json:"ticker"`
		Conid  int32  `json:"conid"`
	}
	json.Unmarshal(out, &equities)

	// PHASE 2: Stream Fundamentals
	fmt.Println("Streaming Fundamentals...")
	conids := []string{}
	for _, e := range equities {
		conids = append(conids, fmt.Sprintf("%d", e.Conid))
	}

	args := append([]string{"run", "--directory", "../ibpython", "main.py", "stream_fundamentals"}, conids...)
	cmd = exec.Command("uv", args...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	var currentXML strings.Builder
	var currentConID int32
	parsing := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "START_CONID:") {
			parsing = true
			fmt.Sscanf(line, "START_CONID:%d", &currentConID)
			currentXML.Reset()
		} else if line == "END_XML_BLOCK" {
			parsing = false
			processXML(ctx, queries, currentConID, currentXML.String())
		} else if parsing {
			currentXML.WriteString(line)
		}
	}
	cmd.Wait()
}

func processXML(ctx context.Context, q *db.Queries, conid int32, data string) {
	if data == "NOT_FOUND" {
		return
	}

	// Define your XML mapping structure here
	type Report struct {
		Ratios struct {
			Ratio []struct {
				ID    string  `xml:"ID,attr"`
				Value float64 `xml:",chardata"`
			} `xml:"Ratio"`
		} `xml:"Ratios"`
	}

	var r Report
	xml.Unmarshal([]byte(data), &r)

	// Extract specific data (example: P/E Ratio)
	for _, ratio := range r.Ratios.Ratio {
		if ratio.ID == "PERatio" {
			peValue := pgtype.Float8{
				Float64: ratio.Value,
				Valid:   true,
			}

			q.UpdateEquityMetrics(ctx, db.UpdateEquityMetricsParams{
				Conid:   conid,
				PeRatio: peValue, // Ensure you have this field in your DB
			})
			fmt.Printf("Updated ConID %d: PE Ratio %f\n", conid, ratio.Value)
		}
	}
}
