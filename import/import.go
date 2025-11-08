package importgames

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"wingspan-scoring/db"
	"wingspan-scoring/scoring"
)

// CSVRecord represents a single row in the import CSV
type CSVRecord struct {
	GameID           string
	Date             string
	IncludeOceania   string
	PlayerName       string
	BirdPoints       string
	BonusCards       string
	RoundGoals       string
	Eggs             string
	CachedFood       string
	TuckedCards      string
	NectarForest     string
	NectarGrassland  string
	NectarWetland    string
	UnusedFood       string
	Total            string
	Rank             string
}

// ImportError represents an error that occurred during import
type ImportError struct {
	Line    int
	GameID  string
	Message string
}

func (e ImportError) Error() string {
	return fmt.Sprintf("line %d (game %s): %s", e.Line, e.GameID, e.Message)
}

// ImportResult contains the results of an import operation
type ImportResult struct {
	GamesImported int
	Errors        []ImportError
}

// ParseCSV reads a CSV file and returns grouped game data
func ParseCSV(reader io.Reader) (map[string][]*CSVRecord, []ImportError) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, []ImportError{{Line: 1, Message: fmt.Sprintf("failed to read header: %v", err)}}
	}

	// Validate header
	expectedHeaders := []string{"GameID", "Date", "IncludeOceania", "PlayerName", "BirdPoints", "BonusCards", "RoundGoals", "Eggs", "CachedFood", "TuckedCards", "NectarForest", "NectarGrassland", "NectarWetland", "UnusedFood", "Total", "Rank"}
	if len(header) != len(expectedHeaders) {
		return nil, []ImportError{{Line: 1, Message: fmt.Sprintf("invalid header: expected %d columns, got %d", len(expectedHeaders), len(header))}}
	}

	// Group records by GameID
	gameRecords := make(map[string][]*CSVRecord)
	var errors []ImportError
	lineNum := 1 // Start at 1 for header

	for {
		lineNum++
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, ImportError{Line: lineNum, Message: fmt.Sprintf("failed to read row: %v", err)})
			continue
		}

		if len(row) != len(expectedHeaders) {
			errors = append(errors, ImportError{Line: lineNum, Message: fmt.Sprintf("invalid column count: expected %d, got %d", len(expectedHeaders), len(row))})
			continue
		}

		record := &CSVRecord{
			GameID:          strings.TrimSpace(row[0]),
			Date:            strings.TrimSpace(row[1]),
			IncludeOceania:  strings.TrimSpace(row[2]),
			PlayerName:      strings.TrimSpace(row[3]),
			BirdPoints:      strings.TrimSpace(row[4]),
			BonusCards:      strings.TrimSpace(row[5]),
			RoundGoals:      strings.TrimSpace(row[6]),
			Eggs:            strings.TrimSpace(row[7]),
			CachedFood:      strings.TrimSpace(row[8]),
			TuckedCards:     strings.TrimSpace(row[9]),
			NectarForest:    strings.TrimSpace(row[10]),
			NectarGrassland: strings.TrimSpace(row[11]),
			NectarWetland:   strings.TrimSpace(row[12]),
			UnusedFood:      strings.TrimSpace(row[13]),
			Total:           strings.TrimSpace(row[14]),
			Rank:            strings.TrimSpace(row[15]),
		}

		if record.GameID == "" {
			errors = append(errors, ImportError{Line: lineNum, GameID: record.GameID, Message: "GameID cannot be empty"})
			continue
		}

		gameRecords[record.GameID] = append(gameRecords[record.GameID], record)
	}

	return gameRecords, errors
}

// ValidateAndConvertGame validates a group of CSV records and converts them to a GameResult
func ValidateAndConvertGame(gameID string, records []*CSVRecord) (*db.GameResult, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("no players for game")
	}

	// All records should have the same date and oceania setting
	firstRecord := records[0]

	// Parse date
	createdAt, err := parseDate(firstRecord.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	// Parse IncludeOceania
	includeOceania, err := parseBool(firstRecord.IncludeOceania)
	if err != nil {
		return nil, fmt.Errorf("invalid IncludeOceania value: %v", err)
	}

	// Validate player count
	if len(records) < 2 || len(records) > 5 {
		return nil, fmt.Errorf("invalid player count: %d (must be 2-5)", len(records))
	}

	// Convert players
	players := make([]scoring.PlayerGameEnd, 0, len(records))
	seenRanks := make(map[int]bool)
	hasWinner := false

	for _, record := range records {
		// Validate consistency
		if record.Date != firstRecord.Date {
			return nil, fmt.Errorf("inconsistent dates within game")
		}
		if record.IncludeOceania != firstRecord.IncludeOceania {
			return nil, fmt.Errorf("inconsistent IncludeOceania within game")
		}

		player, err := convertPlayer(record, includeOceania)
		if err != nil {
			return nil, fmt.Errorf("player %s: %v", record.PlayerName, err)
		}

		// Validate rank uniqueness
		if seenRanks[player.Rank] {
			return nil, fmt.Errorf("duplicate rank %d", player.Rank)
		}
		seenRanks[player.Rank] = true

		if player.Rank == 1 {
			hasWinner = true
		}

		players = append(players, *player)
	}

	// Validate ranks are sequential
	for i := 1; i <= len(players); i++ {
		if !seenRanks[i] {
			return nil, fmt.Errorf("ranks must be sequential: missing rank %d", i)
		}
	}

	if !hasWinner {
		return nil, fmt.Errorf("no player with rank 1 (winner)")
	}

	// Calculate nectar scoring if Oceania expansion is included
	var nectarScoring *scoring.NectarScoring
	if includeOceania {
		nectarScoring = calculateNectarScoring(players)
	}

	// Determine winner
	var winnerName string
	var winnerScore int
	for _, player := range players {
		if player.Rank == 1 {
			winnerName = player.PlayerName
			winnerScore = player.Total
			break
		}
	}

	return &db.GameResult{
		CreatedAt:      createdAt,
		NumPlayers:     len(players),
		IncludeOceania: includeOceania,
		WinnerName:     winnerName,
		WinnerScore:    winnerScore,
		Players:        players,
		NectarScoring:  nectarScoring,
	}, nil
}

// convertPlayer converts a CSV record to a PlayerGameEnd
func convertPlayer(record *CSVRecord, includeOceania bool) (*scoring.PlayerGameEnd, error) {
	if record.PlayerName == "" {
		return nil, fmt.Errorf("player name cannot be empty")
	}

	birdPoints, err := parseInt(record.BirdPoints, "BirdPoints")
	if err != nil {
		return nil, err
	}

	bonusCards, err := parseInt(record.BonusCards, "BonusCards")
	if err != nil {
		return nil, err
	}

	roundGoals, err := parseInt(record.RoundGoals, "RoundGoals")
	if err != nil {
		return nil, err
	}

	eggs, err := parseInt(record.Eggs, "Eggs")
	if err != nil {
		return nil, err
	}

	cachedFood, err := parseInt(record.CachedFood, "CachedFood")
	if err != nil {
		return nil, err
	}

	tuckedCards, err := parseInt(record.TuckedCards, "TuckedCards")
	if err != nil {
		return nil, err
	}

	rank, err := parseInt(record.Rank, "Rank")
	if err != nil {
		return nil, err
	}
	if rank < 1 {
		return nil, fmt.Errorf("rank must be >= 1")
	}

	// Parse optional fields
	unusedFood := 0
	if record.UnusedFood != "" {
		unusedFood, err = parseInt(record.UnusedFood, "UnusedFood")
		if err != nil {
			return nil, err
		}
	}

	// Calculate total if not provided
	total := 0
	if record.Total != "" {
		total, err = parseInt(record.Total, "Total")
		if err != nil {
			return nil, err
		}
	} else {
		total = birdPoints + bonusCards + roundGoals + eggs + cachedFood + tuckedCards
	}

	player := &scoring.PlayerGameEnd{
		PlayerName:  record.PlayerName,
		BirdPoints:  birdPoints,
		BonusCards:  bonusCards,
		RoundGoals:  roundGoals,
		Eggs:        eggs,
		CachedFood:  cachedFood,
		TuckedCards: tuckedCards,
		UnusedFood:  unusedFood,
		Total:       total,
		Rank:        rank,
	}

	// Parse nectar fields if Oceania is included
	if includeOceania {
		nectarForest, err := parseInt(record.NectarForest, "NectarForest")
		if err != nil {
			return nil, err
		}
		nectarGrassland, err := parseInt(record.NectarGrassland, "NectarGrassland")
		if err != nil {
			return nil, err
		}
		nectarWetland, err := parseInt(record.NectarWetland, "NectarWetland")
		if err != nil {
			return nil, err
		}

		player.NectarForest = nectarForest
		player.NectarGrassland = nectarGrassland
		player.NectarWetland = nectarWetland
	}

	return player, nil
}

// calculateNectarScoring calculates nectar scoring from player nectar counts
func calculateNectarScoring(players []scoring.PlayerGameEnd) *scoring.NectarScoring {
	type nectarCount struct {
		playerName string
		count      int
	}

	// Helper to find top 2 players
	findTopTwo := func(counts []nectarCount) (first string, second string) {
		if len(counts) == 0 {
			return "", ""
		}

		// Sort by count (descending)
		for i := 0; i < len(counts)-1; i++ {
			for j := i + 1; j < len(counts); j++ {
				if counts[j].count > counts[i].count {
					counts[i], counts[j] = counts[j], counts[i]
				}
			}
		}

		first = counts[0].playerName
		if len(counts) > 1 && counts[1].count > 0 {
			second = counts[1].playerName
		}

		return first, second
	}

	// Collect nectar counts
	forestCounts := make([]nectarCount, 0, len(players))
	grasslandCounts := make([]nectarCount, 0, len(players))
	wetlandCounts := make([]nectarCount, 0, len(players))

	for _, p := range players {
		forestCounts = append(forestCounts, nectarCount{p.PlayerName, p.NectarForest})
		grasslandCounts = append(grasslandCounts, nectarCount{p.PlayerName, p.NectarGrassland})
		wetlandCounts = append(wetlandCounts, nectarCount{p.PlayerName, p.NectarWetland})
	}

	// Find winners
	forestFirst, forestSecond := findTopTwo(forestCounts)
	grasslandFirst, grasslandSecond := findTopTwo(grasslandCounts)
	wetlandFirst, wetlandSecond := findTopTwo(wetlandCounts)

	// Build scoring maps
	nectarScoring := &scoring.NectarScoring{
		Forest:    make(map[string]int),
		Grassland: make(map[string]int),
		Wetland:   make(map[string]int),
	}

	if forestFirst != "" {
		nectarScoring.Forest[forestFirst] = 5
	}
	if forestSecond != "" {
		nectarScoring.Forest[forestSecond] = 2
	}

	if grasslandFirst != "" {
		nectarScoring.Grassland[grasslandFirst] = 5
	}
	if grasslandSecond != "" {
		nectarScoring.Grassland[grasslandSecond] = 2
	}

	if wetlandFirst != "" {
		nectarScoring.Wetland[wetlandFirst] = 5
	}
	if wetlandSecond != "" {
		nectarScoring.Wetland[wetlandSecond] = 2
	}

	return nectarScoring
}

// parseInt parses a string to an integer with validation
func parseInt(s, fieldName string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("%s cannot be empty", fieldName)
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %v", fieldName, err)
	}

	if val < 0 {
		return 0, fmt.Errorf("%s must be non-negative", fieldName)
	}

	return val, nil
}

// parseBool parses a string to a boolean
func parseBool(s string) (bool, error) {
	lower := strings.ToLower(s)
	switch lower {
	case "true", "yes", "1", "y":
		return true, nil
	case "false", "no", "0", "n":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}

// parseDate attempts to parse a date string in multiple formats
func parseDate(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"01/02/2006 15:04:05",
		"01/02/2006 15:04",
		"01/02/2006",
		"1/2/2006 15:04:05",
		"1/2/2006 15:04",
		"1/2/2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", s)
}

// ImportGames imports games from CSV data
func ImportGames(reader io.Reader) (*ImportResult, error) {
	gameRecords, parseErrors := ParseCSV(reader)

	result := &ImportResult{
		Errors: parseErrors,
	}

	// Convert and validate each game
	games := make([]*db.GameResult, 0, len(gameRecords))
	for gameID, records := range gameRecords {
		game, err := ValidateAndConvertGame(gameID, records)
		if err != nil {
			result.Errors = append(result.Errors, ImportError{
				GameID:  gameID,
				Message: err.Error(),
			})
			continue
		}
		games = append(games, game)
	}

	// If there were any validation errors, don't import anything
	if len(result.Errors) > 0 {
		return result, fmt.Errorf("import failed: %d errors found", len(result.Errors))
	}

	// Import all games
	for _, game := range games {
		nectarScoring := scoring.NectarScoring{}
		if game.NectarScoring != nil {
			nectarScoring = *game.NectarScoring
		}

		if _, err := db.SaveGameResult(game.Players, nectarScoring, game.IncludeOceania); err != nil {
			return result, fmt.Errorf("failed to save game: %v", err)
		}
		result.GamesImported++
	}

	return result, nil
}
