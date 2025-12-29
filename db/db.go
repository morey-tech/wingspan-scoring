package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// Initialize opens the database connection and creates tables if they don't exist
func Initialize() error {
	// Get database path from environment variable or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/wingspan.db"
	}

	// Resolve absolute path
	if !filepath.IsAbs(dbPath) {
		absPath, err := filepath.Abs(dbPath)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path: %w", err)
		}
		dbPath = absPath
	}
	fmt.Println("Resolved database path:", dbPath)

	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		fmt.Printf("Failed to ping database: %v\n", err)
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// createTables creates the necessary database tables
func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS game_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		num_players INTEGER NOT NULL,
		include_oceania BOOLEAN NOT NULL,
		winner_name TEXT NOT NULL,
		winner_score INTEGER NOT NULL,
		players_json TEXT NOT NULL,
		nectar_json TEXT,
		round_breakdown_json TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_created_at ON game_results(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_winner_name ON game_results(winner_name);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		return err
	}

	// Migration for existing databases (safe to run multiple times)
	// ALTER TABLE will fail silently if column already exists
	_, _ = DB.Exec(`ALTER TABLE game_results ADD COLUMN round_breakdown_json TEXT`)

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
