package export

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"
	"wingspan-scoring/db"
	"wingspan-scoring/scoring"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportGamesToCSV_EmptyList(t *testing.T) {
	csvData, err := ExportGamesToCSV([]db.GameResult{})
	require.NoError(t, err)

	// Should only have header row
	reader := csv.NewReader(strings.NewReader(string(csvData)))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	assert.Len(t, records, 1, "Should have only header row")
	assert.Equal(t, csvHeader, records[0])
}

func TestExportGamesToCSV_SingleGameMultiplePlayers(t *testing.T) {
	games := []db.GameResult{
		{
			ID:             1,
			CreatedAt:      time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			NumPlayers:     2,
			IncludeOceania: false,
			WinnerName:     "Alice",
			WinnerScore:    92,
			Players: []scoring.PlayerGameEnd{
				{
					PlayerName:  "Alice",
					BirdPoints:  45,
					BonusCards:  12,
					RoundGoals:  18,
					Eggs:        9,
					CachedFood:  5,
					TuckedCards: 3,
					UnusedFood:  2,
					Total:       92,
					Rank:        1,
				},
				{
					PlayerName:  "Bob",
					BirdPoints:  40,
					BonusCards:  10,
					RoundGoals:  15,
					Eggs:        8,
					CachedFood:  4,
					TuckedCards: 2,
					UnusedFood:  1,
					Total:       79,
					Rank:        2,
				},
			},
		},
	}

	csvData, err := ExportGamesToCSV(games)
	require.NoError(t, err)

	reader := csv.NewReader(strings.NewReader(string(csvData)))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Header + 2 player rows
	assert.Len(t, records, 3)

	// Check header
	assert.Equal(t, csvHeader, records[0])

	// Check Alice's row
	assert.Equal(t, "1", records[1][0])          // GameID
	assert.Equal(t, "2024-01-15", records[1][1]) // Date
	assert.Equal(t, "false", records[1][2])      // IncludeOceania
	assert.Equal(t, "Alice", records[1][3])      // PlayerName
	assert.Equal(t, "45", records[1][4])         // BirdPoints
	assert.Equal(t, "12", records[1][5])         // BonusCards
	assert.Equal(t, "18", records[1][6])         // RoundGoals
	assert.Equal(t, "9", records[1][7])          // Eggs
	assert.Equal(t, "5", records[1][8])          // CachedFood
	assert.Equal(t, "3", records[1][9])          // TuckedCards
	assert.Equal(t, "0", records[1][10])         // NectarForest
	assert.Equal(t, "0", records[1][11])         // NectarGrassland
	assert.Equal(t, "0", records[1][12])         // NectarWetland
	assert.Equal(t, "2", records[1][13])         // UnusedFood
	assert.Equal(t, "92", records[1][14])        // Total
	assert.Equal(t, "1", records[1][15])         // Rank

	// Check Bob's row
	assert.Equal(t, "1", records[2][0])   // Same GameID
	assert.Equal(t, "Bob", records[2][3]) // PlayerName
	assert.Equal(t, "79", records[2][14]) // Total
	assert.Equal(t, "2", records[2][15])  // Rank
}

func TestExportGamesToCSV_MultipleGames(t *testing.T) {
	games := []db.GameResult{
		{
			ID:             1,
			CreatedAt:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			NumPlayers:     2,
			IncludeOceania: false,
			WinnerName:     "Alice",
			WinnerScore:    92,
			Players: []scoring.PlayerGameEnd{
				{PlayerName: "Alice", Total: 92, Rank: 1},
				{PlayerName: "Bob", Total: 79, Rank: 2},
			},
		},
		{
			ID:             2,
			CreatedAt:      time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			NumPlayers:     2,
			IncludeOceania: false,
			WinnerName:     "Carol",
			WinnerScore:    85,
			Players: []scoring.PlayerGameEnd{
				{PlayerName: "Carol", Total: 85, Rank: 1},
				{PlayerName: "Dave", Total: 80, Rank: 2},
			},
		},
	}

	csvData, err := ExportGamesToCSV(games)
	require.NoError(t, err)

	reader := csv.NewReader(strings.NewReader(string(csvData)))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Header + 2 games * 2 players = 5 rows
	assert.Len(t, records, 5)

	// Verify game IDs
	assert.Equal(t, "1", records[1][0])
	assert.Equal(t, "1", records[2][0])
	assert.Equal(t, "2", records[3][0])
	assert.Equal(t, "2", records[4][0])

	// Verify dates
	assert.Equal(t, "2024-01-15", records[1][1])
	assert.Equal(t, "2024-01-16", records[3][1])
}

func TestExportGamesToCSV_WithOceaniaExpansion(t *testing.T) {
	games := []db.GameResult{
		{
			ID:             1,
			CreatedAt:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			NumPlayers:     2,
			IncludeOceania: true,
			WinnerName:     "Alice",
			WinnerScore:    111,
			Players: []scoring.PlayerGameEnd{
				{
					PlayerName:      "Alice",
					BirdPoints:      50,
					BonusCards:      15,
					RoundGoals:      20,
					Eggs:            10,
					CachedFood:      6,
					TuckedCards:     4,
					NectarForest:    3,
					NectarGrassland: 2,
					NectarWetland:   1,
					UnusedFood:      3,
					Total:           111,
					Rank:            1,
				},
				{
					PlayerName:      "Bob",
					BirdPoints:      48,
					BonusCards:      12,
					RoundGoals:      18,
					Eggs:            9,
					CachedFood:      5,
					TuckedCards:     3,
					NectarForest:    2,
					NectarGrassland: 3,
					NectarWetland:   2,
					UnusedFood:      2,
					Total:           102,
					Rank:            2,
				},
			},
			NectarScoring: &scoring.NectarScoring{
				Forest:    map[string]int{"Alice": 5, "Bob": 2},
				Grassland: map[string]int{"Alice": 2, "Bob": 5},
				Wetland:   map[string]int{"Alice": 5, "Bob": 2},
			},
		},
	}

	csvData, err := ExportGamesToCSV(games)
	require.NoError(t, err)

	reader := csv.NewReader(strings.NewReader(string(csvData)))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	assert.Len(t, records, 3) // Header + 2 players

	// Check Oceania flag
	assert.Equal(t, "true", records[1][2])
	assert.Equal(t, "true", records[2][2])

	// Check Alice's nectar values
	assert.Equal(t, "3", records[1][10]) // NectarForest
	assert.Equal(t, "2", records[1][11]) // NectarGrassland
	assert.Equal(t, "1", records[1][12]) // NectarWetland

	// Check Bob's nectar values
	assert.Equal(t, "2", records[2][10]) // NectarForest
	assert.Equal(t, "3", records[2][11]) // NectarGrassland
	assert.Equal(t, "2", records[2][12]) // NectarWetland
}

func TestExportGamesToCSV_HeaderMatchesExpectedFormat(t *testing.T) {
	expectedHeader := []string{
		"GameID",
		"Date",
		"IncludeOceania",
		"PlayerName",
		"BirdPoints",
		"BonusCards",
		"RoundGoals",
		"Eggs",
		"CachedFood",
		"TuckedCards",
		"NectarForest",
		"NectarGrassland",
		"NectarWetland",
		"UnusedFood",
		"Total",
		"Rank",
	}

	assert.Equal(t, expectedHeader, csvHeader)
	assert.Len(t, csvHeader, 16, "CSV should have 16 columns")
}

func TestExportGamesToCSV_SpecialCharactersInPlayerName(t *testing.T) {
	games := []db.GameResult{
		{
			ID:             1,
			CreatedAt:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			NumPlayers:     1,
			IncludeOceania: false,
			WinnerName:     "Alice, Jr.",
			WinnerScore:    92,
			Players: []scoring.PlayerGameEnd{
				{
					PlayerName: "Alice, Jr.",
					Total:      92,
					Rank:       1,
				},
			},
		},
	}

	csvData, err := ExportGamesToCSV(games)
	require.NoError(t, err)

	// CSV should properly escape the comma in the name
	reader := csv.NewReader(strings.NewReader(string(csvData)))
	records, err := reader.ReadAll()
	require.NoError(t, err)

	assert.Equal(t, "Alice, Jr.", records[1][3])
}
