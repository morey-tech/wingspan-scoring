package goals

import (
	"sort"
)

// PlayerScore represents a player's score for a goal
type PlayerScore struct {
	PlayerName string `json:"playerName"`
	Count      int    `json:"count"`
	Points     int    `json:"points"`
	Rank       int    `json:"rank"` // 1 = 1st place, 2 = 2nd place, etc.
}

// Green side scoring rules based on round number
var greenScoringRules = map[int]map[int]int{
	1: {1: 4, 2: 1, 3: 0}, // Round 1: 1st=4, 2nd=1, 3rd=0
	2: {1: 5, 2: 2, 3: 0}, // Round 2: 1st=5, 2nd=2, 3rd=0
	3: {1: 6, 2: 3, 3: 2}, // Round 3: 1st=6, 2nd=3, 3rd=2
	4: {1: 7, 2: 4, 3: 2}, // Round 4: 1st=7, 2nd=4, 3rd=2
}

// CalculateGreenScores calculates scores using the Green (competitive) scoring method
// Players are ranked by their counts, with ties resolved by averaging points
func CalculateGreenScores(playerCounts map[string]int, round int) []PlayerScore {
	if round < 1 || round > 4 {
		round = 1 // Default to round 1 if invalid
	}

	// Convert to slice for sorting
	scores := make([]PlayerScore, 0, len(playerCounts))
	for name, count := range playerCounts {
		scores = append(scores, PlayerScore{
			PlayerName: name,
			Count:      count,
		})
	}

	// Sort by count (descending) - players with same count get same rank
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].Count == scores[j].Count {
			return scores[i].PlayerName < scores[j].PlayerName // Alphabetical for stability
		}
		return scores[i].Count > scores[j].Count
	})

	// Assign ranks and calculate points
	currentRank := 1
	for i := 0; i < len(scores); i++ {
		// Find all players tied at this count
		tiedPlayers := []int{i}
		for j := i + 1; j < len(scores) && scores[j].Count == scores[i].Count; j++ {
			tiedPlayers = append(tiedPlayers, j)
		}

		// Calculate points for this rank group
		if len(tiedPlayers) > 1 {
			// Tie: sum points for all tied ranks and divide
			totalPoints := 0
			for r := currentRank; r < currentRank+len(tiedPlayers) && r <= 3; r++ {
				if pts, ok := greenScoringRules[round][r]; ok {
					totalPoints += pts
				}
			}
			avgPoints := totalPoints / len(tiedPlayers) // Integer division (rounds down)

			for _, idx := range tiedPlayers {
				scores[idx].Rank = currentRank
				// Players with 0 count always get 0 points, regardless of rank
				if scores[idx].Count == 0 {
					scores[idx].Points = 0
				} else {
					scores[idx].Points = avgPoints
				}
			}
		} else {
			// No tie
			scores[i].Rank = currentRank
			// Players with 0 count always get 0 points, regardless of rank
			if scores[i].Count == 0 {
				scores[i].Points = 0
			} else if pts, ok := greenScoringRules[round][currentRank]; ok && currentRank <= 3 {
				scores[i].Points = pts
			} else {
				scores[i].Points = 0 // 4th place and below get 0
			}
		}

		currentRank += len(tiedPlayers)
		i += len(tiedPlayers) - 1
	}

	return scores
}

// CalculateBlueScores calculates scores using the Blue (linear) scoring method
// Each player gets 1 point per item, maximum 5 points
func CalculateBlueScores(playerCounts map[string]int) []PlayerScore {
	scores := make([]PlayerScore, 0, len(playerCounts))

	for name, count := range playerCounts {
		points := count
		if points > 5 {
			points = 5 // Maximum 5 points
		}
		if points < 0 {
			points = 0
		}

		scores = append(scores, PlayerScore{
			PlayerName: name,
			Count:      count,
			Points:     points,
		})
	}

	// Sort by points descending for display consistency
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].Points == scores[j].Points {
			return scores[i].PlayerName < scores[j].PlayerName
		}
		return scores[i].Points > scores[j].Points
	})

	return scores
}
