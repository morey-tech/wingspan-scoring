package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"wingspan-goals/scoring"
)

// GameResult represents a saved game result
type GameResult struct {
	ID             int64                        `json:"id"`
	CreatedAt      time.Time                    `json:"createdAt"`
	NumPlayers     int                          `json:"numPlayers"`
	IncludeOceania bool                         `json:"includeOceania"`
	WinnerName     string                       `json:"winnerName"`
	WinnerScore    int                          `json:"winnerScore"`
	Players        []scoring.PlayerFinalScore   `json:"players"`
	NectarScoring  *scoring.NectarScoring       `json:"nectarScoring,omitempty"`
}

// SaveGameResult saves a game result to the database
func SaveGameResult(players []scoring.PlayerFinalScore, nectarScoring scoring.NectarScoring, includeOceania bool) (int64, error) {
	if len(players) == 0 {
		return 0, fmt.Errorf("no players provided")
	}

	// Find the winner (player with rank 1)
	var winnerName string
	var winnerScore int
	for _, p := range players {
		if p.Rank == 1 {
			winnerName = p.PlayerName
			winnerScore = p.Total
			break
		}
	}

	if winnerName == "" {
		return 0, fmt.Errorf("no winner found in player data")
	}

	// Marshal players to JSON
	playersJSON, err := json.Marshal(players)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal players: %w", err)
	}

	// Marshal nectar scoring to JSON (nullable)
	var nectarJSON *string
	if includeOceania {
		nectarBytes, err := json.Marshal(nectarScoring)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal nectar scoring: %w", err)
		}
		nectarStr := string(nectarBytes)
		nectarJSON = &nectarStr
	}

	// Insert into database
	query := `
		INSERT INTO game_results (num_players, include_oceania, winner_name, winner_score, players_json, nectar_json)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := DB.Exec(query, len(players), includeOceania, winnerName, winnerScore, string(playersJSON), nectarJSON)
	if err != nil {
		return 0, fmt.Errorf("failed to insert game result: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

// GetGameResult retrieves a single game result by ID
func GetGameResult(id int64) (*GameResult, error) {
	query := `
		SELECT id, created_at, num_players, include_oceania, winner_name, winner_score, players_json, nectar_json
		FROM game_results
		WHERE id = ?
	`

	row := DB.QueryRow(query, id)

	var result GameResult
	var playersJSON string
	var nectarJSON sql.NullString

	err := row.Scan(
		&result.ID,
		&result.CreatedAt,
		&result.NumPlayers,
		&result.IncludeOceania,
		&result.WinnerName,
		&result.WinnerScore,
		&playersJSON,
		&nectarJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game result not found")
		}
		return nil, fmt.Errorf("failed to query game result: %w", err)
	}

	// Unmarshal players
	if err := json.Unmarshal([]byte(playersJSON), &result.Players); err != nil {
		return nil, fmt.Errorf("failed to unmarshal players: %w", err)
	}

	// Unmarshal nectar scoring if present
	if nectarJSON.Valid {
		var nectar scoring.NectarScoring
		if err := json.Unmarshal([]byte(nectarJSON.String), &nectar); err != nil {
			return nil, fmt.Errorf("failed to unmarshal nectar scoring: %w", err)
		}
		result.NectarScoring = &nectar
	}

	return &result, nil
}

// GetAllGameResults retrieves all game results with pagination
func GetAllGameResults(limit, offset int) ([]GameResult, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}

	query := `
		SELECT id, created_at, num_players, include_oceania, winner_name, winner_score, players_json, nectar_json
		FROM game_results
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query game results: %w", err)
	}
	defer rows.Close()

	var results []GameResult
	for rows.Next() {
		var result GameResult
		var playersJSON string
		var nectarJSON sql.NullString

		err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.NumPlayers,
			&result.IncludeOceania,
			&result.WinnerName,
			&result.WinnerScore,
			&playersJSON,
			&nectarJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Unmarshal players
		if err := json.Unmarshal([]byte(playersJSON), &result.Players); err != nil {
			return nil, fmt.Errorf("failed to unmarshal players: %w", err)
		}

		// Unmarshal nectar scoring if present
		if nectarJSON.Valid {
			var nectar scoring.NectarScoring
			if err := json.Unmarshal([]byte(nectarJSON.String), &nectar); err != nil {
				return nil, fmt.Errorf("failed to unmarshal nectar scoring: %w", err)
			}
			result.NectarScoring = &nectar
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

// CountGameResults returns the total number of game results
func CountGameResults() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM game_results").Scan(&count)
	return count, err
}

// GetPlayerStats returns statistics for a specific player
func GetPlayerStats(playerName string) (map[string]interface{}, error) {
	// Get total games played
	var gamesPlayed int
	err := DB.QueryRow(`
		SELECT COUNT(*)
		FROM game_results
		WHERE players_json LIKE ?
	`, "%\""+playerName+"\"%").Scan(&gamesPlayed)
	if err != nil {
		return nil, fmt.Errorf("failed to count games played: %w", err)
	}

	// Get total wins
	var wins int
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM game_results
		WHERE winner_name = ?
	`, playerName).Scan(&wins)
	if err != nil {
		return nil, fmt.Errorf("failed to count wins: %w", err)
	}

	// Get average score - this requires parsing JSON, so we'll get all games and calculate
	query := `
		SELECT players_json
		FROM game_results
		WHERE players_json LIKE ?
	`

	rows, err := DB.Query(query, "%\""+playerName+"\"%")
	if err != nil {
		return nil, fmt.Errorf("failed to query player games: %w", err)
	}
	defer rows.Close()

	totalScore := 0
	gameCount := 0
	for rows.Next() {
		var playersJSON string
		if err := rows.Scan(&playersJSON); err != nil {
			continue
		}

		var players []scoring.PlayerFinalScore
		if err := json.Unmarshal([]byte(playersJSON), &players); err != nil {
			continue
		}

		for _, p := range players {
			if p.PlayerName == playerName {
				totalScore += p.Total
				gameCount++
				break
			}
		}
	}

	avgScore := 0.0
	if gameCount > 0 {
		avgScore = float64(totalScore) / float64(gameCount)
	}

	stats := map[string]interface{}{
		"playerName":   playerName,
		"gamesPlayed":  gamesPlayed,
		"wins":         wins,
		"averageScore": avgScore,
		"winRate":      0.0,
	}

	if gamesPlayed > 0 {
		stats["winRate"] = float64(wins) / float64(gamesPlayed) * 100
	}

	return stats, nil
}

// DeleteGameResult deletes a game result by ID
func DeleteGameResult(id int64) error {
	query := `DELETE FROM game_results WHERE id = ?`

	result, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete game result: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("game result not found")
	}

	return nil
}
