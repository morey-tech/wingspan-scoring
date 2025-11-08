package importgames

import (
	"os"
	"strings"
	"testing"

	"wingspan-scoring/db"
	"wingspan-scoring/scoring"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCSV_ValidData(t *testing.T) {
	csvData := `GameID,Date,IncludeOceania,PlayerName,BirdPoints,BonusCards,RoundGoals,Eggs,CachedFood,TuckedCards,NectarForest,NectarGrassland,NectarWetland,UnusedFood,Total,Rank
1,2024-01-15,false,Alice,45,12,18,9,5,3,0,0,0,2,92,1
1,2024-01-15,false,Bob,40,10,15,8,4,2,0,0,0,1,79,2`

	reader := strings.NewReader(csvData)
	gameRecords, errors := ParseCSV(reader)

	assert.Empty(t, errors)
	assert.Len(t, gameRecords, 1)
	assert.Len(t, gameRecords["1"], 2)
	assert.Equal(t, "Alice", gameRecords["1"][0].PlayerName)
	assert.Equal(t, "Bob", gameRecords["1"][1].PlayerName)
}

func TestParseCSV_MultipleGames(t *testing.T) {
	csvData := `GameID,Date,IncludeOceania,PlayerName,BirdPoints,BonusCards,RoundGoals,Eggs,CachedFood,TuckedCards,NectarForest,NectarGrassland,NectarWetland,UnusedFood,Total,Rank
1,2024-01-15,false,Alice,45,12,18,9,5,3,0,0,0,2,92,1
1,2024-01-15,false,Bob,40,10,15,8,4,2,0,0,0,1,79,2
2,2024-01-16,true,Carol,50,15,20,10,6,4,3,2,1,3,111,1
2,2024-01-16,true,Dave,48,12,18,9,5,3,2,3,2,2,102,2`

	reader := strings.NewReader(csvData)
	gameRecords, errors := ParseCSV(reader)

	assert.Empty(t, errors)
	assert.Len(t, gameRecords, 2)
	assert.Len(t, gameRecords["1"], 2)
	assert.Len(t, gameRecords["2"], 2)
}

func TestParseCSV_InvalidHeader(t *testing.T) {
	csvData := `GameID,Date,PlayerName
1,2024-01-15,Alice`

	reader := strings.NewReader(csvData)
	_, errors := ParseCSV(reader)

	assert.NotEmpty(t, errors)
	assert.Contains(t, errors[0].Error(), "invalid header")
}

func TestParseCSV_InvalidColumnCount(t *testing.T) {
	csvData := `GameID,Date,IncludeOceania,PlayerName,BirdPoints,BonusCards,RoundGoals,Eggs,CachedFood,TuckedCards,NectarForest,NectarGrassland,NectarWetland,UnusedFood,Total,Rank
1,2024-01-15,false,Alice`

	reader := strings.NewReader(csvData)
	_, errors := ParseCSV(reader)

	assert.NotEmpty(t, errors)
	// The CSV reader returns a different error message, so just check that there's an error
	assert.Contains(t, errors[0].Error(), "failed to read row")
}

func TestParseCSV_EmptyGameID(t *testing.T) {
	csvData := `GameID,Date,IncludeOceania,PlayerName,BirdPoints,BonusCards,RoundGoals,Eggs,CachedFood,TuckedCards,NectarForest,NectarGrassland,NectarWetland,UnusedFood,Total,Rank
,2024-01-15,false,Alice,45,12,18,9,5,3,0,0,0,2,92,1`

	reader := strings.NewReader(csvData)
	_, errors := ParseCSV(reader)

	assert.NotEmpty(t, errors)
	assert.Contains(t, errors[0].Error(), "GameID cannot be empty")
}

func TestValidateAndConvertGame_ValidGame(t *testing.T) {
	records := []*CSVRecord{
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "false",
			PlayerName: "Alice", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
			Eggs: "9", CachedFood: "5", TuckedCards: "3",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "2", Total: "92", Rank: "1",
		},
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "false",
			PlayerName: "Bob", BirdPoints: "40", BonusCards: "10", RoundGoals: "15",
			Eggs: "8", CachedFood: "4", TuckedCards: "2",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "1", Total: "79", Rank: "2",
		},
	}

	game, err := ValidateAndConvertGame("1", records)

	require.NoError(t, err)
	assert.Equal(t, 2, game.NumPlayers)
	assert.False(t, game.IncludeOceania)
	assert.Equal(t, "Alice", game.WinnerName)
	assert.Equal(t, 92, game.WinnerScore)
	assert.Len(t, game.Players, 2)
	assert.Equal(t, "Alice", game.Players[0].PlayerName)
	assert.Equal(t, 1, game.Players[0].Rank)
	assert.Equal(t, "Bob", game.Players[1].PlayerName)
	assert.Equal(t, 2, game.Players[1].Rank)
}

func TestValidateAndConvertGame_WithOceania(t *testing.T) {
	records := []*CSVRecord{
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "true",
			PlayerName: "Alice", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
			Eggs: "9", CachedFood: "5", TuckedCards: "3",
			NectarForest: "5", NectarGrassland: "3", NectarWetland: "2",
			UnusedFood: "2", Total: "92", Rank: "1",
		},
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "true",
			PlayerName: "Bob", BirdPoints: "40", BonusCards: "10", RoundGoals: "15",
			Eggs: "8", CachedFood: "4", TuckedCards: "2",
			NectarForest: "3", NectarGrassland: "5", NectarWetland: "4",
			UnusedFood: "1", Total: "79", Rank: "2",
		},
	}

	game, err := ValidateAndConvertGame("1", records)

	require.NoError(t, err)
	assert.True(t, game.IncludeOceania)
	assert.NotNil(t, game.NectarScoring)
	assert.Equal(t, 5, game.Players[0].NectarForest)
	assert.Equal(t, 3, game.Players[0].NectarGrassland)
	assert.Equal(t, 2, game.Players[0].NectarWetland)
}

func TestValidateAndConvertGame_InvalidPlayerCount(t *testing.T) {
	// Only 1 player
	records := []*CSVRecord{
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "false",
			PlayerName: "Alice", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
			Eggs: "9", CachedFood: "5", TuckedCards: "3",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "2", Total: "92", Rank: "1",
		},
	}

	_, err := ValidateAndConvertGame("1", records)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid player count")
}

func TestValidateAndConvertGame_NoWinner(t *testing.T) {
	records := []*CSVRecord{
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "false",
			PlayerName: "Alice", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
			Eggs: "9", CachedFood: "5", TuckedCards: "3",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "2", Total: "92", Rank: "2",
		},
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "false",
			PlayerName: "Bob", BirdPoints: "40", BonusCards: "10", RoundGoals: "15",
			Eggs: "8", CachedFood: "4", TuckedCards: "2",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "1", Total: "79", Rank: "3",
		},
	}

	_, err := ValidateAndConvertGame("1", records)
	assert.Error(t, err)
	// Sequential rank check happens before winner check, so we get "missing rank 1" error
	assert.Contains(t, err.Error(), "missing rank 1")
}

func TestValidateAndConvertGame_DuplicateRank(t *testing.T) {
	records := []*CSVRecord{
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "false",
			PlayerName: "Alice", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
			Eggs: "9", CachedFood: "5", TuckedCards: "3",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "2", Total: "92", Rank: "1",
		},
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "false",
			PlayerName: "Bob", BirdPoints: "40", BonusCards: "10", RoundGoals: "15",
			Eggs: "8", CachedFood: "4", TuckedCards: "2",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "1", Total: "79", Rank: "1",
		},
	}

	_, err := ValidateAndConvertGame("1", records)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate rank")
}

func TestValidateAndConvertGame_NonSequentialRanks(t *testing.T) {
	records := []*CSVRecord{
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "false",
			PlayerName: "Alice", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
			Eggs: "9", CachedFood: "5", TuckedCards: "3",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "2", Total: "92", Rank: "1",
		},
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "false",
			PlayerName: "Bob", BirdPoints: "40", BonusCards: "10", RoundGoals: "15",
			Eggs: "8", CachedFood: "4", TuckedCards: "2",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "1", Total: "79", Rank: "3",
		},
	}

	_, err := ValidateAndConvertGame("1", records)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ranks must be sequential")
}

func TestValidateAndConvertGame_InconsistentDates(t *testing.T) {
	records := []*CSVRecord{
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "false",
			PlayerName: "Alice", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
			Eggs: "9", CachedFood: "5", TuckedCards: "3",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "2", Total: "92", Rank: "1",
		},
		{
			GameID: "1", Date: "2024-01-16", IncludeOceania: "false",
			PlayerName: "Bob", BirdPoints: "40", BonusCards: "10", RoundGoals: "15",
			Eggs: "8", CachedFood: "4", TuckedCards: "2",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "1", Total: "79", Rank: "2",
		},
	}

	_, err := ValidateAndConvertGame("1", records)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "inconsistent dates")
}

func TestValidateAndConvertGame_InconsistentOceania(t *testing.T) {
	records := []*CSVRecord{
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "false",
			PlayerName: "Alice", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
			Eggs: "9", CachedFood: "5", TuckedCards: "3",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "2", Total: "92", Rank: "1",
		},
		{
			GameID: "1", Date: "2024-01-15", IncludeOceania: "true",
			PlayerName: "Bob", BirdPoints: "40", BonusCards: "10", RoundGoals: "15",
			Eggs: "8", CachedFood: "4", TuckedCards: "2",
			NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
			UnusedFood: "1", Total: "79", Rank: "2",
		},
	}

	_, err := ValidateAndConvertGame("1", records)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "inconsistent IncludeOceania")
}

func TestConvertPlayer_ValidPlayer(t *testing.T) {
	record := &CSVRecord{
		PlayerName: "Alice", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
		Eggs: "9", CachedFood: "5", TuckedCards: "3",
		NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
		UnusedFood: "2", Total: "92", Rank: "1",
	}

	player, err := convertPlayer(record, false)

	require.NoError(t, err)
	assert.Equal(t, "Alice", player.PlayerName)
	assert.Equal(t, 45, player.BirdPoints)
	assert.Equal(t, 12, player.BonusCards)
	assert.Equal(t, 18, player.RoundGoals)
	assert.Equal(t, 9, player.Eggs)
	assert.Equal(t, 5, player.CachedFood)
	assert.Equal(t, 3, player.TuckedCards)
	assert.Equal(t, 2, player.UnusedFood)
	assert.Equal(t, 92, player.Total)
	assert.Equal(t, 1, player.Rank)
}

func TestConvertPlayer_AutoCalculateTotal(t *testing.T) {
	record := &CSVRecord{
		PlayerName: "Alice", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
		Eggs: "9", CachedFood: "5", TuckedCards: "3",
		NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
		UnusedFood: "2", Total: "", Rank: "1",
	}

	player, err := convertPlayer(record, false)

	require.NoError(t, err)
	assert.Equal(t, 92, player.Total) // 45 + 12 + 18 + 9 + 5 + 3 = 92
}

func TestConvertPlayer_EmptyPlayerName(t *testing.T) {
	record := &CSVRecord{
		PlayerName: "", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
		Eggs: "9", CachedFood: "5", TuckedCards: "3",
		NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
		UnusedFood: "2", Total: "92", Rank: "1",
	}

	_, err := convertPlayer(record, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "player name cannot be empty")
}

func TestConvertPlayer_InvalidBirdPoints(t *testing.T) {
	record := &CSVRecord{
		PlayerName: "Alice", BirdPoints: "invalid", BonusCards: "12", RoundGoals: "18",
		Eggs: "9", CachedFood: "5", TuckedCards: "3",
		NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
		UnusedFood: "2", Total: "92", Rank: "1",
	}

	_, err := convertPlayer(record, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid BirdPoints")
}

func TestConvertPlayer_NegativeValue(t *testing.T) {
	record := &CSVRecord{
		PlayerName: "Alice", BirdPoints: "-5", BonusCards: "12", RoundGoals: "18",
		Eggs: "9", CachedFood: "5", TuckedCards: "3",
		NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
		UnusedFood: "2", Total: "92", Rank: "1",
	}

	_, err := convertPlayer(record, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be non-negative")
}

func TestConvertPlayer_InvalidRank(t *testing.T) {
	record := &CSVRecord{
		PlayerName: "Alice", BirdPoints: "45", BonusCards: "12", RoundGoals: "18",
		Eggs: "9", CachedFood: "5", TuckedCards: "3",
		NectarForest: "0", NectarGrassland: "0", NectarWetland: "0",
		UnusedFood: "2", Total: "92", Rank: "0",
	}

	_, err := convertPlayer(record, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rank must be >= 1")
}

func TestParseBool_ValidValues(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"yes", true},
		{"Yes", true},
		{"1", true},
		{"y", true},
		{"false", false},
		{"False", false},
		{"FALSE", false},
		{"no", false},
		{"No", false},
		{"0", false},
		{"n", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseBool(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseBool_InvalidValue(t *testing.T) {
	_, err := parseBool("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid boolean value")
}

func TestParseDate_ValidFormats(t *testing.T) {
	tests := []string{
		"2024-01-15 14:30:05",
		"2024-01-15 14:30",
		"2024-01-15",
		"01/15/2024 14:30:05",
		"01/15/2024 14:30",
		"01/15/2024",
		"1/15/2024 14:30:05",
		"1/15/2024 14:30",
		"1/15/2024",
	}

	for _, dateStr := range tests {
		t.Run(dateStr, func(t *testing.T) {
			_, err := parseDate(dateStr)
			assert.NoError(t, err)
		})
	}
}

func TestParseDate_InvalidFormat(t *testing.T) {
	_, err := parseDate("invalid-date")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to parse date")
}

func TestCalculateNectarScoring(t *testing.T) {
	players := []scoring.PlayerGameEnd{
		{PlayerName: "Alice", NectarForest: 5, NectarGrassland: 3, NectarWetland: 2},
		{PlayerName: "Bob", NectarForest: 3, NectarGrassland: 5, NectarWetland: 4},
		{PlayerName: "Carol", NectarForest: 2, NectarGrassland: 2, NectarWetland: 5},
	}

	nectar := calculateNectarScoring(players)

	// Alice should win Forest (5 > 3 > 2)
	assert.Equal(t, 5, nectar.Forest["Alice"])
	assert.Equal(t, 2, nectar.Forest["Bob"])

	// Bob should win Grassland (5 > 3 > 2)
	assert.Equal(t, 5, nectar.Grassland["Bob"])
	assert.Equal(t, 2, nectar.Grassland["Alice"])

	// Carol should win Wetland (5 > 4 > 2)
	assert.Equal(t, 5, nectar.Wetland["Carol"])
	assert.Equal(t, 2, nectar.Wetland["Bob"])
}

func TestImportGames_Integration(t *testing.T) {
	csvData := `GameID,Date,IncludeOceania,PlayerName,BirdPoints,BonusCards,RoundGoals,Eggs,CachedFood,TuckedCards,NectarForest,NectarGrassland,NectarWetland,UnusedFood,Total,Rank
1,2024-01-15,false,Alice,45,12,18,9,5,3,0,0,0,2,92,1
1,2024-01-15,false,Bob,40,10,15,8,4,2,0,0,0,1,79,2
2,2024-01-16,true,Carol,50,15,20,10,6,4,3,2,1,3,111,1
2,2024-01-16,true,Dave,48,12,18,9,5,3,2,3,2,2,102,2`

	// Create temporary database
	tmpDB, err := os.CreateTemp("", "test-import-*.db")
	require.NoError(t, err)
	defer os.Remove(tmpDB.Name())

	// Initialize the database with the temporary path
	os.Setenv("DB_PATH", tmpDB.Name())
	defer os.Unsetenv("DB_PATH")

	err = db.Initialize()
	require.NoError(t, err)
	defer db.Close()

	reader := strings.NewReader(csvData)
	result, err := ImportGames(reader)

	require.NoError(t, err)
	assert.Equal(t, 2, result.GamesImported)
	assert.Empty(t, result.Errors)

	// Verify games were saved
	count, err := db.CountGameResults()
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestImportGames_WithErrors(t *testing.T) {
	csvData := `GameID,Date,IncludeOceania,PlayerName,BirdPoints,BonusCards,RoundGoals,Eggs,CachedFood,TuckedCards,NectarForest,NectarGrassland,NectarWetland,UnusedFood,Total,Rank
1,2024-01-15,false,Alice,45,12,18,9,5,3,0,0,0,2,92,1
2,invalid-date,false,Bob,40,10,15,8,4,2,0,0,0,1,79,1`

	tmpDB, err := os.CreateTemp("", "test-import-*.db")
	require.NoError(t, err)
	defer os.Remove(tmpDB.Name())

	// Initialize the database with the temporary path
	os.Setenv("DB_PATH", tmpDB.Name())
	defer os.Unsetenv("DB_PATH")

	err = db.Initialize()
	require.NoError(t, err)
	defer db.Close()

	reader := strings.NewReader(csvData)
	result, err := ImportGames(reader)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "import failed")
	assert.NotEmpty(t, result.Errors)
	assert.Equal(t, 0, result.GamesImported) // Should not import anything if there are errors
}
