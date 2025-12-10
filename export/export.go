package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"wingspan-scoring/db"
)

// CSV header matching the import format
var csvHeader = []string{
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

// ExportGamesToCSV converts game results to CSV format matching the import format.
// Each player in each game becomes one row in the CSV.
func ExportGamesToCSV(games []db.GameResult) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	if err := writer.Write(csvHeader); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, game := range games {
		gameID := strconv.FormatInt(game.ID, 10)
		date := game.CreatedAt.Format("2006-01-02")
		includeOceania := strconv.FormatBool(game.IncludeOceania)

		for _, player := range game.Players {
			row := []string{
				gameID,
				date,
				includeOceania,
				player.PlayerName,
				strconv.Itoa(player.BirdPoints),
				strconv.Itoa(player.BonusCards),
				strconv.Itoa(player.RoundGoals),
				strconv.Itoa(player.Eggs),
				strconv.Itoa(player.CachedFood),
				strconv.Itoa(player.TuckedCards),
				strconv.Itoa(player.NectarForest),
				strconv.Itoa(player.NectarGrassland),
				strconv.Itoa(player.NectarWetland),
				strconv.Itoa(player.UnusedFood),
				strconv.Itoa(player.Total),
				strconv.Itoa(player.Rank),
			}

			if err := writer.Write(row); err != nil {
				return nil, fmt.Errorf("failed to write CSV row for game %d, player %s: %w",
					game.ID, player.PlayerName, err)
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}
