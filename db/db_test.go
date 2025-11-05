package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitialize_DefaultPath tests database initialization with default path
func TestInitialize_DefaultPath(t *testing.T) {
	// Save original DB and defer restore
	originalDB := DB
	defer func() { DB = originalDB }()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wingspan-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Clear DB_PATH environment variable
	originalDBPath := os.Getenv("DB_PATH")
	defer os.Setenv("DB_PATH", originalDBPath)
	os.Unsetenv("DB_PATH")

	// Initialize database
	err = Initialize()
	require.NoError(t, err)
	require.NotNil(t, DB)

	// Verify database file was created in default location
	_, err = os.Stat("./data/wingspan.db")
	assert.NoError(t, err)

	// Verify we can ping the database
	err = DB.Ping()
	assert.NoError(t, err)

	// Clean up
	Close()
}

// TestInitialize_CustomPath tests database initialization with custom path from environment
func TestInitialize_CustomPath(t *testing.T) {
	// Save original DB and defer restore
	originalDB := DB
	defer func() { DB = originalDB }()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wingspan-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set custom DB_PATH
	customPath := filepath.Join(tmpDir, "custom", "test.db")
	originalDBPath := os.Getenv("DB_PATH")
	defer os.Setenv("DB_PATH", originalDBPath)
	os.Setenv("DB_PATH", customPath)

	// Initialize database
	err = Initialize()
	require.NoError(t, err)
	require.NotNil(t, DB)

	// Verify database file was created at custom location
	_, err = os.Stat(customPath)
	assert.NoError(t, err)

	// Verify we can ping the database
	err = DB.Ping()
	assert.NoError(t, err)

	// Clean up
	Close()
}

// TestInitialize_CreatesDirectory tests that Initialize creates parent directories
func TestInitialize_CreatesDirectory(t *testing.T) {
	// Save original DB and defer restore
	originalDB := DB
	defer func() { DB = originalDB }()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wingspan-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set DB_PATH to a path with non-existent parent directories
	deepPath := filepath.Join(tmpDir, "level1", "level2", "level3", "test.db")
	originalDBPath := os.Getenv("DB_PATH")
	defer os.Setenv("DB_PATH", originalDBPath)
	os.Setenv("DB_PATH", deepPath)

	// Initialize database
	err = Initialize()
	require.NoError(t, err)
	require.NotNil(t, DB)

	// Verify all parent directories were created
	_, err = os.Stat(filepath.Dir(deepPath))
	assert.NoError(t, err)

	// Verify database file was created
	_, err = os.Stat(deepPath)
	assert.NoError(t, err)

	// Clean up
	Close()
}

// TestInitialize_CreatesTablesAndIndexes tests that Initialize creates proper schema
func TestInitialize_CreatesTablesAndIndexes(t *testing.T) {
	// Save original DB and defer restore
	originalDB := DB
	defer func() { DB = originalDB }()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wingspan-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set custom DB_PATH
	customPath := filepath.Join(tmpDir, "test.db")
	originalDBPath := os.Getenv("DB_PATH")
	defer os.Setenv("DB_PATH", originalDBPath)
	os.Setenv("DB_PATH", customPath)

	// Initialize database
	err = Initialize()
	require.NoError(t, err)
	require.NotNil(t, DB)

	// Verify game_results table exists
	var tableName string
	err = DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='game_results'").Scan(&tableName)
	assert.NoError(t, err)
	assert.Equal(t, "game_results", tableName)

	// Verify indexes exist
	rows, err := DB.Query("SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='game_results'")
	require.NoError(t, err)
	defer rows.Close()

	indexes := []string{}
	for rows.Next() {
		var indexName string
		err := rows.Scan(&indexName)
		require.NoError(t, err)
		indexes = append(indexes, indexName)
	}

	// Should have at least the two indexes we created (plus possibly an auto-created one for PRIMARY KEY)
	assert.Contains(t, indexes, "idx_created_at")
	assert.Contains(t, indexes, "idx_winner_name")

	// Clean up
	Close()
}

// TestInitialize_IdempotentTableCreation tests that Initialize can be called multiple times
func TestInitialize_IdempotentTableCreation(t *testing.T) {
	// Save original DB and defer restore
	originalDB := DB
	defer func() { DB = originalDB }()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wingspan-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set custom DB_PATH
	customPath := filepath.Join(tmpDir, "test.db")
	originalDBPath := os.Getenv("DB_PATH")
	defer os.Setenv("DB_PATH", originalDBPath)
	os.Setenv("DB_PATH", customPath)

	// Initialize database first time
	err = Initialize()
	require.NoError(t, err)

	// Close the connection
	Close()

	// Initialize database second time (should not error)
	err = Initialize()
	assert.NoError(t, err)
	assert.NotNil(t, DB)

	// Clean up
	Close()
}

// TestClose_ValidConnection tests closing a valid database connection
func TestClose_ValidConnection(t *testing.T) {
	// Save original DB and defer restore
	originalDB := DB
	defer func() { DB = originalDB }()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wingspan-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set custom DB_PATH
	customPath := filepath.Join(tmpDir, "test.db")
	originalDBPath := os.Getenv("DB_PATH")
	defer os.Setenv("DB_PATH", originalDBPath)
	os.Setenv("DB_PATH", customPath)

	// Initialize database
	err = Initialize()
	require.NoError(t, err)

	// Close should succeed
	err = Close()
	assert.NoError(t, err)

	// Verify database is actually closed (Ping should fail)
	err = DB.Ping()
	assert.Error(t, err)
}

// TestClose_NilConnection tests closing when DB is nil
func TestClose_NilConnection(t *testing.T) {
	// Save original DB and defer restore
	originalDB := DB
	defer func() { DB = originalDB }()

	// Set DB to nil
	DB = nil

	// Close should not error when DB is nil
	err := Close()
	assert.NoError(t, err)
}

// TestCreateTables_VerifySchema tests that the schema is created correctly
func TestCreateTables_VerifySchema(t *testing.T) {
	// Save original DB and defer restore
	originalDB := DB
	defer func() { DB = originalDB }()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wingspan-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test database
	customPath := filepath.Join(tmpDir, "test.db")
	DB, err = sql.Open("sqlite", customPath)
	require.NoError(t, err)

	// Call createTables
	err = createTables()
	require.NoError(t, err)

	// Verify schema by checking column information
	rows, err := DB.Query("PRAGMA table_info(game_results)")
	require.NoError(t, err)
	defer rows.Close()

	columns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var dfltValue sql.NullString

		err := rows.Scan(&cid, &name, &dataType, &notNull, &dfltValue, &pk)
		require.NoError(t, err)
		columns[name] = true
	}

	// Verify all expected columns exist
	expectedColumns := []string{
		"id", "created_at", "num_players", "include_oceania",
		"winner_name", "winner_score", "players_json", "nectar_json",
	}

	for _, col := range expectedColumns {
		assert.True(t, columns[col], "Expected column %s to exist", col)
	}

	// Clean up
	Close()
}

// TestInitialize_AbsolutePath tests initialization with an absolute path
func TestInitialize_AbsolutePath(t *testing.T) {
	// Save original DB and defer restore
	originalDB := DB
	defer func() { DB = originalDB }()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wingspan-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Use an absolute path directly
	absolutePath := filepath.Join(tmpDir, "absolute.db")
	originalDBPath := os.Getenv("DB_PATH")
	defer os.Setenv("DB_PATH", originalDBPath)
	os.Setenv("DB_PATH", absolutePath)

	// Initialize database
	err = Initialize()
	require.NoError(t, err)
	require.NotNil(t, DB)

	// Verify database file was created at exact location
	_, err = os.Stat(absolutePath)
	assert.NoError(t, err)

	// Clean up
	Close()
}
