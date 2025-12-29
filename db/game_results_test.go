package db

import (
	"os"
	"path/filepath"
	"testing"
	"wingspan-scoring/scoring"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a temporary database for testing
func setupTestDB(t *testing.T) func() {
	// Save original DB
	originalDB := DB

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wingspan-test-*")
	require.NoError(t, err)

	// Set custom DB_PATH
	customPath := filepath.Join(tmpDir, "test.db")
	originalDBPath := os.Getenv("DB_PATH")
	os.Setenv("DB_PATH", customPath)

	// Initialize database
	err = Initialize()
	require.NoError(t, err)
	require.NotNil(t, DB)

	// Return cleanup function
	return func() {
		Close()
		DB = originalDB
		os.Setenv("DB_PATH", originalDBPath)
		os.RemoveAll(tmpDir)
	}
}

// TestSaveGameResult_WithoutOceania tests saving a game without Oceania expansion
func TestSaveGameResult_WithoutOceania(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	players := []scoring.PlayerGameEnd{
		{PlayerName: "Alice", Total: 100, Rank: 1},
		{PlayerName: "Bob", Total: 90, Rank: 2},
	}
	nectarScoring := scoring.NectarScoring{}

	id, err := SaveGameResult(players, nectarScoring, false)
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// Verify the game was saved correctly
	result, err := GetGameResult(id)
	require.NoError(t, err)
	assert.Equal(t, 2, result.NumPlayers)
	assert.False(t, result.IncludeOceania)
	assert.Equal(t, "Alice", result.WinnerName)
	assert.Equal(t, 100, result.WinnerScore)
	assert.Len(t, result.Players, 2)
	assert.Nil(t, result.NectarScoring)
}

// TestSaveGameResult_WithOceania tests saving a game with Oceania expansion
func TestSaveGameResult_WithOceania(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	players := []scoring.PlayerGameEnd{
		{PlayerName: "Alice", Total: 105, Rank: 1},
		{PlayerName: "Bob", Total: 95, Rank: 2},
	}
	nectarScoring := scoring.NectarScoring{
		Forest:    map[string]int{"Alice": 5, "Bob": 2},
		Grassland: map[string]int{"Alice": 3, "Bob": 5},
		Wetland:   map[string]int{"Alice": 2, "Bob": 0},
	}

	id, err := SaveGameResult(players, nectarScoring, true)
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// Verify the game was saved correctly
	result, err := GetGameResult(id)
	require.NoError(t, err)
	assert.Equal(t, 2, result.NumPlayers)
	assert.True(t, result.IncludeOceania)
	assert.Equal(t, "Alice", result.WinnerName)
	assert.Equal(t, 105, result.WinnerScore)
	assert.NotNil(t, result.NectarScoring)
	assert.Equal(t, 5, result.NectarScoring.Forest["Alice"])
	assert.Equal(t, 2, result.NectarScoring.Forest["Bob"])
}

// TestSaveGameResult_EmptyPlayers tests error when no players provided
func TestSaveGameResult_EmptyPlayers(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	players := []scoring.PlayerGameEnd{}
	nectarScoring := scoring.NectarScoring{}

	id, err := SaveGameResult(players, nectarScoring, false)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Contains(t, err.Error(), "no players provided")
}

// TestSaveGameResult_NoWinner tests error when no winner (rank 1) found
func TestSaveGameResult_NoWinner(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	players := []scoring.PlayerGameEnd{
		{PlayerName: "Alice", Total: 100, Rank: 2},
		{PlayerName: "Bob", Total: 90, Rank: 3},
	}
	nectarScoring := scoring.NectarScoring{}

	id, err := SaveGameResult(players, nectarScoring, false)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Contains(t, err.Error(), "no winner found")
}

// TestSaveGameResult_MultiplePlayers tests saving with multiple players
func TestSaveGameResult_MultiplePlayers(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	players := []scoring.PlayerGameEnd{
		{PlayerName: "Alice", Total: 110, Rank: 1, BirdPoints: 50, BonusCards: 10},
		{PlayerName: "Bob", Total: 100, Rank: 2, BirdPoints: 45, BonusCards: 12},
		{PlayerName: "Carol", Total: 95, Rank: 3, BirdPoints: 40, BonusCards: 15},
		{PlayerName: "Dave", Total: 85, Rank: 4, BirdPoints: 38, BonusCards: 8},
	}
	nectarScoring := scoring.NectarScoring{}

	id, err := SaveGameResult(players, nectarScoring, false)
	require.NoError(t, err)

	// Verify all players were saved
	result, err := GetGameResult(id)
	require.NoError(t, err)
	assert.Equal(t, 4, result.NumPlayers)
	assert.Len(t, result.Players, 4)
	assert.Equal(t, "Alice", result.WinnerName)
	assert.Equal(t, 110, result.WinnerScore)
}

// TestGetGameResult_NotFound tests getting a non-existent game
func TestGetGameResult_NotFound(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	result, err := GetGameResult(999)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

// TestGetGameResult_ValidGame tests retrieving an existing game
func TestGetGameResult_ValidGame(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save a game first
	players := []scoring.PlayerGameEnd{
		{PlayerName: "Alice", Total: 100, Rank: 1, BirdPoints: 50},
		{PlayerName: "Bob", Total: 90, Rank: 2, BirdPoints: 45},
	}
	nectarScoring := scoring.NectarScoring{}

	id, err := SaveGameResult(players, nectarScoring, false)
	require.NoError(t, err)

	// Retrieve the game
	result, err := GetGameResult(id)
	require.NoError(t, err)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, "Alice", result.WinnerName)
	assert.Equal(t, 100, result.WinnerScore)
	assert.Len(t, result.Players, 2)
	assert.Equal(t, 50, result.Players[0].BirdPoints)
}

// TestGetAllGameResults_Empty tests getting results when database is empty
func TestGetAllGameResults_Empty(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	results, err := GetAllGameResults(10, 0)
	require.NoError(t, err)
	assert.Len(t, results, 0)
}

// TestGetAllGameResults_Pagination tests pagination functionality
func TestGetAllGameResults_Pagination(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save 5 games
	for i := 0; i < 5; i++ {
		players := []scoring.PlayerGameEnd{
			{PlayerName: "Alice", Total: 100 + i, Rank: 1},
		}
		nectarScoring := scoring.NectarScoring{}
		_, err := SaveGameResult(players, nectarScoring, false)
		require.NoError(t, err)
	}

	// Get first 2 results
	results, err := GetAllGameResults(2, 0)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Get next 2 results
	results2, err := GetAllGameResults(2, 2)
	require.NoError(t, err)
	assert.Len(t, results2, 2)

	// Verify they're different
	assert.NotEqual(t, results[0].ID, results2[0].ID)

	// Get last result
	results3, err := GetAllGameResults(2, 4)
	require.NoError(t, err)
	assert.Len(t, results3, 1)
}

// TestGetAllGameResults_DefaultLimit tests default limit when 0 is provided
func TestGetAllGameResults_DefaultLimit(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save 3 games
	for i := 0; i < 3; i++ {
		players := []scoring.PlayerGameEnd{
			{PlayerName: "Alice", Total: 100, Rank: 1},
		}
		nectarScoring := scoring.NectarScoring{}
		_, err := SaveGameResult(players, nectarScoring, false)
		require.NoError(t, err)
	}

	// Get results with limit 0 (should use default)
	results, err := GetAllGameResults(0, 0)
	require.NoError(t, err)
	assert.Len(t, results, 3) // Should get all 3 with default limit of 50
}

// TestGetAllGameResults_OrderByCreatedAt tests results are ordered by creation date
func TestGetAllGameResults_OrderByCreatedAt(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save 3 games with different winners
	winners := []string{"Alice", "Bob", "Carol"}
	var ids []int64
	for _, winner := range winners {
		players := []scoring.PlayerGameEnd{
			{PlayerName: winner, Total: 100, Rank: 1},
		}
		nectarScoring := scoring.NectarScoring{}
		id, err := SaveGameResult(players, nectarScoring, false)
		require.NoError(t, err)
		ids = append(ids, id)
	}

	// Get all results
	results, err := GetAllGameResults(10, 0)
	require.NoError(t, err)
	assert.Len(t, results, 3)

	// Verify all 3 games are returned (order may vary if created in same second)
	returnedIDs := make(map[int64]bool)
	for _, result := range results {
		returnedIDs[result.ID] = true
	}
	for _, expectedID := range ids {
		assert.True(t, returnedIDs[expectedID], "Expected game ID %d to be in results", expectedID)
	}
}

// TestCountGameResults_Empty tests count when no games exist
func TestCountGameResults_Empty(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	count, err := CountGameResults()
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

// TestCountGameResults_Multiple tests count with multiple games
func TestCountGameResults_Multiple(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save 7 games
	for i := 0; i < 7; i++ {
		players := []scoring.PlayerGameEnd{
			{PlayerName: "Alice", Total: 100, Rank: 1},
		}
		nectarScoring := scoring.NectarScoring{}
		_, err := SaveGameResult(players, nectarScoring, false)
		require.NoError(t, err)
	}

	count, err := CountGameResults()
	require.NoError(t, err)
	assert.Equal(t, 7, count)
}

// TestGetPlayerStats_NoGames tests stats for a player with no games
func TestGetPlayerStats_NoGames(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	stats, err := GetPlayerStats("NonExistent")
	require.NoError(t, err)
	assert.Equal(t, "NonExistent", stats["playerName"])
	assert.Equal(t, 0, stats["gamesPlayed"])
	assert.Equal(t, 0, stats["wins"])
	assert.Equal(t, 0.0, stats["averageScore"])
	assert.Equal(t, 0.0, stats["winRate"])
}

// TestGetPlayerStats_WithGamesAndWins tests stats for a player with wins
func TestGetPlayerStats_WithGamesAndWins(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save 3 games where Alice plays and wins 2
	games := [][]scoring.PlayerGameEnd{
		{
			{PlayerName: "Alice", Total: 100, Rank: 1},
			{PlayerName: "Bob", Total: 90, Rank: 2},
		},
		{
			{PlayerName: "Bob", Total: 105, Rank: 1},
			{PlayerName: "Alice", Total: 95, Rank: 2},
		},
		{
			{PlayerName: "Alice", Total: 110, Rank: 1},
			{PlayerName: "Carol", Total: 85, Rank: 2},
		},
	}

	for _, players := range games {
		nectarScoring := scoring.NectarScoring{}
		_, err := SaveGameResult(players, nectarScoring, false)
		require.NoError(t, err)
	}

	stats, err := GetPlayerStats("Alice")
	require.NoError(t, err)
	assert.Equal(t, "Alice", stats["playerName"])
	assert.Equal(t, 3, stats["gamesPlayed"])
	assert.Equal(t, 2, stats["wins"])

	// Average: (100 + 95 + 110) / 3 = 305 / 3 = 101.666...
	avgScore := stats["averageScore"].(float64)
	assert.InDelta(t, 101.67, avgScore, 0.01)

	// Win rate: 2/3 * 100 = 66.666...%
	winRate := stats["winRate"].(float64)
	assert.InDelta(t, 66.67, winRate, 0.01)
}

// TestGetPlayerStats_GamesButNoWins tests stats for a player who never wins
func TestGetPlayerStats_GamesButNoWins(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save 2 games where Bob plays but never wins
	games := [][]scoring.PlayerGameEnd{
		{
			{PlayerName: "Alice", Total: 100, Rank: 1},
			{PlayerName: "Bob", Total: 90, Rank: 2},
		},
		{
			{PlayerName: "Carol", Total: 105, Rank: 1},
			{PlayerName: "Bob", Total: 95, Rank: 2},
		},
	}

	for _, players := range games {
		nectarScoring := scoring.NectarScoring{}
		_, err := SaveGameResult(players, nectarScoring, false)
		require.NoError(t, err)
	}

	stats, err := GetPlayerStats("Bob")
	require.NoError(t, err)
	assert.Equal(t, "Bob", stats["playerName"])
	assert.Equal(t, 2, stats["gamesPlayed"])
	assert.Equal(t, 0, stats["wins"])

	// Average: (90 + 95) / 2 = 92.5
	avgScore := stats["averageScore"].(float64)
	assert.Equal(t, 92.5, avgScore)

	// Win rate: 0%
	assert.Equal(t, 0.0, stats["winRate"])
}

// TestGetPlayerStats_MultipleGames tests accurate stat calculation across many games
func TestGetPlayerStats_MultipleGames(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save 10 games where Alice wins 7
	for i := 0; i < 10; i++ {
		var players []scoring.PlayerGameEnd
		if i < 7 {
			players = []scoring.PlayerGameEnd{
				{PlayerName: "Alice", Total: 100 + i, Rank: 1},
				{PlayerName: "Bob", Total: 90 + i, Rank: 2},
			}
		} else {
			players = []scoring.PlayerGameEnd{
				{PlayerName: "Bob", Total: 105 + i, Rank: 1},
				{PlayerName: "Alice", Total: 95 + i, Rank: 2},
			}
		}
		nectarScoring := scoring.NectarScoring{}
		_, err := SaveGameResult(players, nectarScoring, false)
		require.NoError(t, err)
	}

	stats, err := GetPlayerStats("Alice")
	require.NoError(t, err)
	assert.Equal(t, 10, stats["gamesPlayed"])
	assert.Equal(t, 7, stats["wins"])

	// Win rate: 7/10 * 100 = 70%
	assert.Equal(t, 70.0, stats["winRate"])
}

// TestDeleteGameResult_Existing tests deleting an existing game
func TestDeleteGameResult_Existing(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save a game
	players := []scoring.PlayerGameEnd{
		{PlayerName: "Alice", Total: 100, Rank: 1},
	}
	nectarScoring := scoring.NectarScoring{}
	id, err := SaveGameResult(players, nectarScoring, false)
	require.NoError(t, err)

	// Verify it exists
	result, err := GetGameResult(id)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Delete it
	err = DeleteGameResult(id)
	assert.NoError(t, err)

	// Verify it no longer exists
	result, err = GetGameResult(id)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestDeleteGameResult_NotFound tests deleting a non-existent game
func TestDeleteGameResult_NotFound(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	err := DeleteGameResult(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestDeleteGameResult_UpdatesCount tests that deletion updates count
func TestDeleteGameResult_UpdatesCount(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save 3 games
	var ids []int64
	for i := 0; i < 3; i++ {
		players := []scoring.PlayerGameEnd{
			{PlayerName: "Alice", Total: 100, Rank: 1},
		}
		nectarScoring := scoring.NectarScoring{}
		id, err := SaveGameResult(players, nectarScoring, false)
		require.NoError(t, err)
		ids = append(ids, id)
	}

	// Verify count
	count, err := CountGameResults()
	require.NoError(t, err)
	assert.Equal(t, 3, count)

	// Delete one game
	err = DeleteGameResult(ids[1])
	require.NoError(t, err)

	// Verify count decreased
	count, err = CountGameResults()
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

// TestSaveAndRetrieve_ComplexGame tests saving and retrieving a complex game with all fields
func TestSaveAndRetrieve_ComplexGame(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	players := []scoring.PlayerGameEnd{
		{
			PlayerName:      "Alice",
			BirdPoints:      50,
			BonusCards:      12,
			RoundGoals:      14,
			Eggs:            9,
			CachedFood:      5,
			TuckedCards:     6,
			NectarForest:    4,
			NectarGrassland: 3,
			NectarWetland:   2,
			UnusedFood:      3,
			Total:           105,
			Rank:            1,
		},
		{
			PlayerName:      "Bob",
			BirdPoints:      45,
			BonusCards:      10,
			RoundGoals:      12,
			Eggs:            8,
			CachedFood:      4,
			TuckedCards:     5,
			NectarForest:    3,
			NectarGrassland: 4,
			NectarWetland:   3,
			UnusedFood:      2,
			Total:           95,
			Rank:            2,
		},
	}

	nectarScoring := scoring.NectarScoring{
		Forest:    map[string]int{"Alice": 5, "Bob": 2},
		Grassland: map[string]int{"Alice": 2, "Bob": 5},
		Wetland:   map[string]int{"Alice": 3, "Bob": 3},
	}

	// Save the game
	id, err := SaveGameResult(players, nectarScoring, true)
	require.NoError(t, err)

	// Retrieve the game
	result, err := GetGameResult(id)
	require.NoError(t, err)

	// Verify all fields
	assert.Equal(t, 2, result.NumPlayers)
	assert.True(t, result.IncludeOceania)
	assert.Equal(t, "Alice", result.WinnerName)
	assert.Equal(t, 105, result.WinnerScore)

	// Verify players
	require.Len(t, result.Players, 2)
	alice := result.Players[0]
	assert.Equal(t, "Alice", alice.PlayerName)
	assert.Equal(t, 50, alice.BirdPoints)
	assert.Equal(t, 12, alice.BonusCards)
	assert.Equal(t, 14, alice.RoundGoals)
	assert.Equal(t, 9, alice.Eggs)
	assert.Equal(t, 5, alice.CachedFood)
	assert.Equal(t, 6, alice.TuckedCards)
	assert.Equal(t, 4, alice.NectarForest)
	assert.Equal(t, 3, alice.NectarGrassland)
	assert.Equal(t, 2, alice.NectarWetland)
	assert.Equal(t, 3, alice.UnusedFood)
	assert.Equal(t, 105, alice.Total)
	assert.Equal(t, 1, alice.Rank)

	// Verify nectar scoring
	require.NotNil(t, result.NectarScoring)
	assert.Equal(t, 5, result.NectarScoring.Forest["Alice"])
	assert.Equal(t, 2, result.NectarScoring.Forest["Bob"])
	assert.Equal(t, 2, result.NectarScoring.Grassland["Alice"])
	assert.Equal(t, 5, result.NectarScoring.Grassland["Bob"])
	assert.Equal(t, 3, result.NectarScoring.Wetland["Alice"])
	assert.Equal(t, 3, result.NectarScoring.Wetland["Bob"])
}

// TestGetLeaderboardStats_Empty tests leaderboard with no games
func TestGetLeaderboardStats_Empty(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	leaderboard, err := GetLeaderboardStats()
	require.NoError(t, err)
	require.NotNil(t, leaderboard)

	// All categories should have empty values
	assert.Equal(t, "", leaderboard.TotalScore.PlayerName)
	assert.Equal(t, 0, leaderboard.TotalScore.Score)
	assert.Equal(t, "", leaderboard.BirdPoints.PlayerName)
	assert.Equal(t, 0, leaderboard.BirdPoints.Score)
}

// TestGetLeaderboardStats_SingleGame tests leaderboard with one game
func TestGetLeaderboardStats_SingleGame(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	players := []scoring.PlayerGameEnd{
		{
			PlayerName:      "Alice",
			BirdPoints:      50,
			BonusCards:      12,
			RoundGoals:      14,
			Eggs:            9,
			CachedFood:      5,
			TuckedCards:     6,
			NectarForest:    5,
			NectarGrassland: 2,
			NectarWetland:   0,
			Total:           103,
			Rank:            1,
		},
		{
			PlayerName:      "Bob",
			BirdPoints:      45,
			BonusCards:      10,
			RoundGoals:      12,
			Eggs:            8,
			CachedFood:      4,
			TuckedCards:     5,
			NectarForest:    2,
			NectarGrassland: 5,
			NectarWetland:   0,
			Total:           91,
			Rank:            2,
		},
	}
	nectarScoring := scoring.NectarScoring{}
	_, err := SaveGameResult(players, nectarScoring, true)
	require.NoError(t, err)

	leaderboard, err := GetLeaderboardStats()
	require.NoError(t, err)
	require.NotNil(t, leaderboard)

	// Verify Alice has highest total score
	assert.Equal(t, "Alice", leaderboard.TotalScore.PlayerName)
	assert.Equal(t, 103, leaderboard.TotalScore.Score)

	// Verify Alice has highest bird points
	assert.Equal(t, "Alice", leaderboard.BirdPoints.PlayerName)
	assert.Equal(t, 50, leaderboard.BirdPoints.Score)

	// Verify Alice has highest bonus cards
	assert.Equal(t, "Alice", leaderboard.BonusCards.PlayerName)
	assert.Equal(t, 12, leaderboard.BonusCards.Score)

	// Verify Bob has highest nectar grassland
	assert.Equal(t, "Bob", leaderboard.NectarGrassland.PlayerName)
	assert.Equal(t, 5, leaderboard.NectarGrassland.Score)
}

// TestGetLeaderboardStats_MultipleGames tests leaderboard with different leaders per category
func TestGetLeaderboardStats_MultipleGames(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Game 1: Alice wins with high bird points
	game1 := []scoring.PlayerGameEnd{
		{
			PlayerName:  "Alice",
			BirdPoints:  60,
			BonusCards:  10,
			RoundGoals:  12,
			Eggs:        8,
			CachedFood:  4,
			TuckedCards: 5,
			Total:       99,
			Rank:        1,
		},
		{
			PlayerName:  "Bob",
			BirdPoints:  45,
			BonusCards:  15,
			RoundGoals:  10,
			Eggs:        7,
			CachedFood:  3,
			TuckedCards: 4,
			Total:       84,
			Rank:        2,
		},
	}
	_, err := SaveGameResult(game1, scoring.NectarScoring{}, false)
	require.NoError(t, err)

	// Game 2: Bob wins with high total score and high eggs
	game2 := []scoring.PlayerGameEnd{
		{
			PlayerName:  "Bob",
			BirdPoints:  50,
			BonusCards:  12,
			RoundGoals:  14,
			Eggs:        15,
			CachedFood:  5,
			TuckedCards: 6,
			Total:       102,
			Rank:        1,
		},
		{
			PlayerName:  "Carol",
			BirdPoints:  48,
			BonusCards:  11,
			RoundGoals:  13,
			Eggs:        9,
			CachedFood:  10,
			TuckedCards: 4,
			Total:       95,
			Rank:        2,
		},
	}
	_, err = SaveGameResult(game2, scoring.NectarScoring{}, false)
	require.NoError(t, err)

	// Game 3: Carol wins with high round goals and tucked cards
	game3 := []scoring.PlayerGameEnd{
		{
			PlayerName:  "Carol",
			BirdPoints:  52,
			BonusCards:  13,
			RoundGoals:  18,
			Eggs:        10,
			CachedFood:  6,
			TuckedCards: 12,
			Total:       111,
			Rank:        1,
		},
		{
			PlayerName:  "Alice",
			BirdPoints:  49,
			BonusCards:  14,
			RoundGoals:  11,
			Eggs:        8,
			CachedFood:  5,
			TuckedCards: 7,
			Total:       94,
			Rank:        2,
		},
	}
	_, err = SaveGameResult(game3, scoring.NectarScoring{}, false)
	require.NoError(t, err)

	leaderboard, err := GetLeaderboardStats()
	require.NoError(t, err)
	require.NotNil(t, leaderboard)

	// Verify different leaders for different categories
	assert.Equal(t, "Carol", leaderboard.TotalScore.PlayerName)
	assert.Equal(t, 111, leaderboard.TotalScore.Score)

	assert.Equal(t, "Alice", leaderboard.BirdPoints.PlayerName)
	assert.Equal(t, 60, leaderboard.BirdPoints.Score)

	assert.Equal(t, "Bob", leaderboard.BonusCards.PlayerName)
	assert.Equal(t, 15, leaderboard.BonusCards.Score)

	assert.Equal(t, "Carol", leaderboard.RoundGoals.PlayerName)
	assert.Equal(t, 18, leaderboard.RoundGoals.Score)

	assert.Equal(t, "Bob", leaderboard.Eggs.PlayerName)
	assert.Equal(t, 15, leaderboard.Eggs.Score)

	assert.Equal(t, "Carol", leaderboard.CachedFood.PlayerName)
	assert.Equal(t, 10, leaderboard.CachedFood.Score)

	assert.Equal(t, "Carol", leaderboard.TuckedCards.PlayerName)
	assert.Equal(t, 12, leaderboard.TuckedCards.Score)
}

// TestGetLeaderboardStats_WithNectar tests leaderboard with nectar scoring
func TestGetLeaderboardStats_WithNectar(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Game with Oceania expansion (nectar scoring)
	players := []scoring.PlayerGameEnd{
		{
			PlayerName:      "Alice",
			BirdPoints:      50,
			BonusCards:      12,
			RoundGoals:      14,
			Eggs:            9,
			CachedFood:      5,
			TuckedCards:     6,
			NectarForest:    5,
			NectarGrassland: 2,
			NectarWetland:   5,
			Total:           108,
			Rank:            1,
		},
		{
			PlayerName:      "Bob",
			BirdPoints:      45,
			BonusCards:      10,
			RoundGoals:      12,
			Eggs:            8,
			CachedFood:      4,
			TuckedCards:     5,
			NectarForest:    2,
			NectarGrassland: 5,
			NectarWetland:   2,
			Total:           93,
			Rank:            2,
		},
	}
	nectarScoring := scoring.NectarScoring{
		Forest:    map[string]int{"Alice": 5, "Bob": 2},
		Grassland: map[string]int{"Alice": 2, "Bob": 5},
		Wetland:   map[string]int{"Alice": 5, "Bob": 2},
	}
	_, err := SaveGameResult(players, nectarScoring, true)
	require.NoError(t, err)

	leaderboard, err := GetLeaderboardStats()
	require.NoError(t, err)
	require.NotNil(t, leaderboard)

	// Verify nectar leaders
	assert.Equal(t, "Alice", leaderboard.NectarForest.PlayerName)
	assert.Equal(t, 5, leaderboard.NectarForest.Score)

	assert.Equal(t, "Bob", leaderboard.NectarGrassland.PlayerName)
	assert.Equal(t, 5, leaderboard.NectarGrassland.Score)

	assert.Equal(t, "Alice", leaderboard.NectarWetland.PlayerName)
	assert.Equal(t, 5, leaderboard.NectarWetland.Score)
}

// TestGetLeaderboardStats_TiedScores tests that first player alphabetically wins ties
func TestGetLeaderboardStats_TiedScores(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Game where both players have same total score
	players := []scoring.PlayerGameEnd{
		{
			PlayerName: "Zoe",
			Total:      100,
			BirdPoints: 50,
			Rank:       1,
		},
		{
			PlayerName: "Alice",
			Total:      100,
			BirdPoints: 50,
			Rank:       2,
		},
	}
	_, err := SaveGameResult(players, scoring.NectarScoring{}, false)
	require.NoError(t, err)

	leaderboard, err := GetLeaderboardStats()
	require.NoError(t, err)
	require.NotNil(t, leaderboard)

	// First encountered player (Zoe) should be the leader since we don't have special tie logic
	assert.Equal(t, "Zoe", leaderboard.TotalScore.PlayerName)
	assert.Equal(t, 100, leaderboard.TotalScore.Score)
}

// TestSaveGameResult_WithRoundBreakdown tests saving and retrieving round goal breakdown
func TestSaveGameResult_WithRoundBreakdown(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	players := []scoring.PlayerGameEnd{
		{
			PlayerName: "Alice",
			RoundGoals: 15,
			RoundGoalsBreakdown: &scoring.RoundGoalBreakdown{
				Round1: 3,
				Round2: 5,
				Round3: 4,
				Round4: 3,
			},
			Total: 100,
			Rank:  1,
		},
		{
			PlayerName: "Bob",
			RoundGoals: 12,
			RoundGoalsBreakdown: &scoring.RoundGoalBreakdown{
				Round1: 2,
				Round2: 3,
				Round3: 4,
				Round4: 3,
			},
			Total: 90,
			Rank:  2,
		},
	}
	nectarScoring := scoring.NectarScoring{}

	// Save game with round breakdown
	id, err := SaveGameResult(players, nectarScoring, false)
	require.NoError(t, err)

	// Retrieve and verify breakdown
	result, err := GetGameResult(id)
	require.NoError(t, err)
	assert.NotNil(t, result.RoundBreakdown)

	// Verify Alice's breakdown
	aliceBreakdown := result.RoundBreakdown["Alice"]
	require.NotNil(t, aliceBreakdown)
	assert.Equal(t, 3, aliceBreakdown.Round1)
	assert.Equal(t, 5, aliceBreakdown.Round2)
	assert.Equal(t, 4, aliceBreakdown.Round3)
	assert.Equal(t, 3, aliceBreakdown.Round4)

	// Verify Bob's breakdown
	bobBreakdown := result.RoundBreakdown["Bob"]
	require.NotNil(t, bobBreakdown)
	assert.Equal(t, 2, bobBreakdown.Round1)
	assert.Equal(t, 3, bobBreakdown.Round2)
	assert.Equal(t, 4, bobBreakdown.Round3)
	assert.Equal(t, 3, bobBreakdown.Round4)
}

// TestSaveGameResult_BackwardCompatibility_NoBreakdown tests that games without breakdown still work
func TestSaveGameResult_BackwardCompatibility_NoBreakdown(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Save game without round breakdown (backward compatibility)
	players := []scoring.PlayerGameEnd{
		{PlayerName: "Alice", RoundGoals: 15, Total: 100, Rank: 1},
		{PlayerName: "Bob", RoundGoals: 12, Total: 90, Rank: 2},
	}
	nectarScoring := scoring.NectarScoring{}

	id, err := SaveGameResult(players, nectarScoring, false)
	require.NoError(t, err)

	// Retrieve and verify breakdown is nil
	result, err := GetGameResult(id)
	require.NoError(t, err)
	assert.Nil(t, result.RoundBreakdown) // Should be nil for games without breakdown
}
