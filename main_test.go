package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"wingspan-scoring/db"
	"wingspan-scoring/goals"
	"wingspan-scoring/scoring"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a temporary database for testing
func setupTestDB(t *testing.T) func() {
	// Save original DB
	originalDB := db.DB

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wingspan-test-*")
	require.NoError(t, err)

	// Set custom DB_PATH
	customPath := filepath.Join(tmpDir, "test.db")
	originalDBPath := os.Getenv("DB_PATH")
	os.Setenv("DB_PATH", customPath)

	// Initialize database
	err = db.Initialize()
	require.NoError(t, err)
	require.NotNil(t, db.DB)

	// Return cleanup function
	return func() {
		db.Close()
		db.DB = originalDB
		os.Setenv("DB_PATH", originalDBPath)
		os.RemoveAll(tmpDir)
	}
}

// TestHandleNewGame_ValidRequest tests POST /api/new-game with valid expansion selections
func TestHandleNewGame_ValidRequest(t *testing.T) {
	testCases := []struct {
		name           string
		base           string
		european       string
		oceania        string
		expectedMinLen int
	}{
		{"Base only", "true", "false", "false", 4},
		{"European only", "false", "true", "false", 4},
		{"Oceania only", "false", "false", "true", 4},
		{"Base + European", "true", "true", "false", 4},
		{"All expansions", "true", "true", "true", 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("base", tc.base)
			form.Add("european", tc.european)
			form.Add("oceania", tc.oceania)

			req := httptest.NewRequest(http.MethodPost, "/api/new-game", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handleNewGame(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var result goals.RoundGoals
			err := json.NewDecoder(w.Body).Decode(&result)
			require.NoError(t, err)

			// Should have 4 unique goals
			assert.NotEmpty(t, result.Round1.ID)
			assert.NotEmpty(t, result.Round2.ID)
			assert.NotEmpty(t, result.Round3.ID)
			assert.NotEmpty(t, result.Round4.ID)
		})
	}
}

// TestHandleNewGame_DefaultToBase tests that no selection defaults to base game
func TestHandleNewGame_DefaultToBase(t *testing.T) {
	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/api/new-game", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handleNewGame(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result goals.RoundGoals
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)

	// Should have selected goals (from base game as default)
	assert.NotEmpty(t, result.Round1.ID)
}

// TestHandleNewGame_InvalidMethod tests that non-POST methods are rejected
func TestHandleNewGame_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/new-game", nil)
	w := httptest.NewRecorder()

	handleNewGame(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// TestHandleGetGoals_DifferentExpansions tests GET /api/goals with various expansions
func TestHandleGetGoals_DifferentExpansions(t *testing.T) {
	testCases := []struct {
		name         string
		base         string
		european     string
		oceania      string
		expectedLen  int
	}{
		{"Base only", "true", "false", "false", 16},
		{"European only", "false", "true", "false", 10},
		{"Oceania only", "false", "false", "true", 8},
		{"All expansions", "true", "true", "true", 34},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/api/goals?base=" + tc.base + "&european=" + tc.european + "&oceania=" + tc.oceania
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			handleGetGoals(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var result []goals.Goal
			err := json.NewDecoder(w.Body).Decode(&result)
			require.NoError(t, err)
			assert.Len(t, result, tc.expectedLen)
		})
	}
}

// TestHandleGetGoals_DefaultToAll tests that no parameters defaults to all expansions
func TestHandleGetGoals_DefaultToAll(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/goals", nil)
	w := httptest.NewRecorder()

	handleGetGoals(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []goals.Goal
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)
	assert.Len(t, result, 34) // All expansions
}

// TestHandleCalculateScores_GreenMode tests green (competitive) scoring
func TestHandleCalculateScores_GreenMode(t *testing.T) {
	requestBody := map[string]interface{}{
		"mode":  "green",
		"round": 1,
		"playerCounts": map[string]int{
			"Alice": 5,
			"Bob":   3,
		},
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/calculate-scores", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleCalculateScores(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result []goals.PlayerScore
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify scores (Round 1: 1st=4, 2nd=1)
	assert.Equal(t, 4, result[0].Points) // Alice
	assert.Equal(t, 1, result[1].Points) // Bob
}

// TestHandleCalculateScores_BlueMode tests blue (linear) scoring
func TestHandleCalculateScores_BlueMode(t *testing.T) {
	requestBody := map[string]interface{}{
		"mode": "blue",
		"playerCounts": map[string]int{
			"Alice": 3,
			"Bob":   7, // Should be capped at 5
		},
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/calculate-scores", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleCalculateScores(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []goals.PlayerScore
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Find each player's score
	for _, score := range result {
		if score.PlayerName == "Alice" {
			assert.Equal(t, 3, score.Points)
		} else if score.PlayerName == "Bob" {
			assert.Equal(t, 5, score.Points) // Capped at 5
		}
	}
}

// TestHandleCalculateScores_InvalidJSON tests error handling for invalid JSON
func TestHandleCalculateScores_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/calculate-scores", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleCalculateScores(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestHandleCalculateScores_InvalidMethod tests non-POST methods are rejected
func TestHandleCalculateScores_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/calculate-scores", nil)
	w := httptest.NewRecorder()

	handleCalculateScores(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// TestHandleCalculateGameEnd_ValidRequest tests game end calculation and saving
func TestHandleCalculateGameEnd_ValidRequest(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	requestBody := map[string]interface{}{
		"includeOceania": false,
		"players": []map[string]interface{}{
			{
				"playerName":  "Alice",
				"birdPoints":  50,
				"bonusCards":  10,
				"roundGoals":  15,
				"eggs":        8,
				"cachedFood":  3,
				"tuckedCards": 5,
				"unusedFood":  2,
			},
			{
				"playerName":  "Bob",
				"birdPoints":  45,
				"bonusCards":  12,
				"roundGoals":  18,
				"eggs":        7,
				"cachedFood":  4,
				"tuckedCards": 3,
				"unusedFood":  3,
			},
		},
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/calculate-game-end", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleCalculateGameEnd(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result struct {
		Players       []scoring.PlayerGameEnd `json:"players"`
		NectarScoring scoring.NectarScoring   `json:"nectarScoring"`
		GameID        int64                   `json:"gameId"`
	}
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)

	assert.Len(t, result.Players, 2)
	assert.Greater(t, result.GameID, int64(0))

	// Verify totals were calculated
	assert.Greater(t, result.Players[0].Total, 0)
	assert.Greater(t, result.Players[1].Total, 0)

	// Verify rankings were assigned
	assert.GreaterOrEqual(t, result.Players[0].Rank, 1)
	assert.GreaterOrEqual(t, result.Players[1].Rank, 1)
}

// TestHandleCalculateGameEnd_InvalidJSON tests error handling
func TestHandleCalculateGameEnd_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/calculate-game-end", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleCalculateGameEnd(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestHandleCalculateGameEnd_InvalidMethod tests non-POST methods are rejected
func TestHandleCalculateGameEnd_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/calculate-game-end", nil)
	w := httptest.NewRecorder()

	handleCalculateGameEnd(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// TestHandleGetGames_ValidRequest tests GET /api/games with pagination
func TestHandleGetGames_ValidRequest(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save some games first
	for i := 0; i < 3; i++ {
		players := []scoring.PlayerGameEnd{
			{PlayerName: "Alice", Total: 100, Rank: 1},
		}
		nectarScoring := scoring.NectarScoring{}
		_, err := db.SaveGameResult(players, nectarScoring, false)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/games?limit=10&offset=0", nil)
	w := httptest.NewRecorder()

	handleGetGames(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result struct {
		Games      []db.GameResult `json:"games"`
		TotalCount int             `json:"totalCount"`
		Limit      int             `json:"limit"`
		Offset     int             `json:"offset"`
	}
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)

	assert.Len(t, result.Games, 3)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)
}

// TestHandleGetGames_InvalidMethod tests non-GET methods are rejected
func TestHandleGetGames_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/games", nil)
	w := httptest.NewRecorder()

	handleGetGames(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// TestHandleGetGame_ValidID tests GET /api/games/{id}
func TestHandleGetGame_ValidID(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save a game
	players := []scoring.PlayerGameEnd{
		{PlayerName: "Alice", Total: 100, Rank: 1},
	}
	nectarScoring := scoring.NectarScoring{}
	id, err := db.SaveGameResult(players, nectarScoring, false)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/games/"+strconv.FormatInt(id, 10), nil)
	w := httptest.NewRecorder()

	handleGetGame(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result db.GameResult
	err = json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, id, result.ID)
	assert.Equal(t, "Alice", result.WinnerName)
}

// TestHandleGetGame_InvalidID tests error for invalid game ID
func TestHandleGetGame_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/games/invalid", nil)
	w := httptest.NewRecorder()

	handleGetGame(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestHandleGetGame_NotFound tests error for non-existent game
func TestHandleGetGame_NotFound(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/games/999", nil)
	w := httptest.NewRecorder()

	handleGetGame(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestHandleDeleteGame_ValidID tests DELETE /api/games/{id}
func TestHandleDeleteGame_ValidID(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save a game
	players := []scoring.PlayerGameEnd{
		{PlayerName: "Alice", Total: 100, Rank: 1},
	}
	nectarScoring := scoring.NectarScoring{}
	id, err := db.SaveGameResult(players, nectarScoring, false)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/api/games/"+strconv.FormatInt(id, 10), nil)
	w := httptest.NewRecorder()

	handleDeleteGame(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result["success"].(bool))

	// Verify game was actually deleted
	_, err = db.GetGameResult(id)
	assert.Error(t, err)
}

// TestHandleDeleteGame_InvalidID tests error for invalid game ID
func TestHandleDeleteGame_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/games/invalid", nil)
	w := httptest.NewRecorder()

	handleDeleteGame(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestHandleGetPlayerStats_ValidPlayer tests GET /api/stats/{playerName}
func TestHandleGetPlayerStats_ValidPlayer(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save some games
	games := [][]scoring.PlayerGameEnd{
		{
			{PlayerName: "Alice", Total: 100, Rank: 1},
			{PlayerName: "Bob", Total: 90, Rank: 2},
		},
		{
			{PlayerName: "Alice", Total: 95, Rank: 2},
			{PlayerName: "Carol", Total: 105, Rank: 1},
		},
	}

	for _, players := range games {
		nectarScoring := scoring.NectarScoring{}
		_, err := db.SaveGameResult(players, nectarScoring, false)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/stats/Alice", nil)
	w := httptest.NewRecorder()

	handleGetPlayerStats(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "Alice", result["playerName"])
	assert.Equal(t, float64(2), result["gamesPlayed"])
	assert.Equal(t, float64(1), result["wins"])
}

// TestHandleGetPlayerStats_EmptyName tests error for empty player name
func TestHandleGetPlayerStats_EmptyName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/stats/", nil)
	w := httptest.NewRecorder()

	handleGetPlayerStats(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestHandleGetPlayerStats_InvalidMethod tests non-GET methods are rejected
func TestHandleGetPlayerStats_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/stats/Alice", nil)
	w := httptest.NewRecorder()

	handleGetPlayerStats(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// TestParseIntDefault tests the parseIntDefault helper function
func TestParseIntDefault(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		defVal   int
		expected int
	}{
		{"Valid integer", "42", 10, 42},
		{"Invalid string", "abc", 10, 10},
		{"Empty string", "", 5, 5},
		{"Zero", "0", 10, 0},
		{"Negative", "-5", 10, -5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := parseIntDefault(tc.input, tc.defVal)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestHandleGameRoute_MethodRouting tests that handleGameRoute routes to correct handler
func TestHandleGameRoute_MethodRouting(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Test GET routing
	req := httptest.NewRequest(http.MethodGet, "/api/games/999", nil)
	w := httptest.NewRecorder()
	handleGameRoute(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code) // Will be NotFound since game doesn't exist

	// Test DELETE routing
	req = httptest.NewRequest(http.MethodDelete, "/api/games/999", nil)
	w = httptest.NewRecorder()
	handleGameRoute(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code) // Will fail since game doesn't exist

	// Test invalid method
	req = httptest.NewRequest(http.MethodPost, "/api/games/1", nil)
	w = httptest.NewRecorder()
	handleGameRoute(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
