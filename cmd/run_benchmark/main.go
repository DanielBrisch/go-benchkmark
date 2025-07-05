package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-benchmark/internal/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/marcboeker/go-duckdb/v2"
)

func main() {
	http.HandleFunc("/benchmark", handleBenchmark)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleBenchmark(w http.ResponseWriter, r *http.Request) {
	var memStart runtime.MemStats
	runtime.ReadMemStats(&memStart)

	db, err := sql.Open("duckdb", "duck.db")
	if err != nil {
		http.Error(w, "Failed to open DuckDB connection: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	csvPath := "internal/data/large_dataset.csv"
	absPath, err := filepath.Abs(csvPath)
	if err != nil {
		log.Fatalf("Failed to resolve CSV path: %v", err)
	}

	loadStart := time.Now()
	_, err = db.Exec(fmt.Sprintf(
		`CREATE OR REPLACE TABLE vendas AS SELECT * FROM read_csv_auto('%s');`,
		absPath,
	))
	if err != nil {
		http.Error(w, "Failed to load CSV into DuckDB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	loadDuration := time.Since(loadStart)

	sqlPath, err := filepath.Abs("internal/data/pipeline.sql")
	absSQLPath, err := filepath.Abs(sqlPath)
	if err != nil {
		http.Error(w, "Failed to resolve pipeline.sql path: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sqlBytes, err := os.ReadFile(absSQLPath)
	if err != nil {
		http.Error(w, "Failed to read pipeline.sql: "+err.Error(), http.StatusInternalServerError)
		return
	}

	queryStart := time.Now()
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		http.Error(w, "Failed to execute SQL pipeline: "+err.Error(), http.StatusInternalServerError)
		return
	}
	queryDuration := time.Since(queryStart)

	var memEnd runtime.MemStats
	runtime.ReadMemStats(&memEnd)
	peakMemoryMB := (memEnd.Sys - memStart.Sys) / 1024 / 1024

	result := models.BenchmarkResult{
		Language:     "go",
		LoadTimeSec:  loadDuration.Seconds(),
		QueryTimeSec: queryDuration.Seconds(),
		PeakMemoryMB: peakMemoryMB,
		DatasetSize:  "1_000_000 rows",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}
