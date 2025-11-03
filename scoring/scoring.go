package scoring

import (
	"sort"
)

// PlayerGameEnd represents a player's complete end-game score
type PlayerGameEnd struct {
	PlayerName       string `json:"playerName"`
	BirdPoints       int    `json:"birdPoints"`
	BonusCards       int    `json:"bonusCards"`
	RoundGoals       int    `json:"roundGoals"`
	Eggs             int    `json:"eggs"`
	CachedFood       int    `json:"cachedFood"`
	TuckedCards      int    `json:"tuckedCards"`
	NectarForest     int    `json:"nectarForest"`     // Oceania expansion
	NectarGrassland  int    `json:"nectarGrassland"`  // Oceania expansion
	NectarWetland    int    `json:"nectarWetland"`    // Oceania expansion
	UnusedFood       int    `json:"unusedFood"`       // For tiebreaker
	Total            int    `json:"total"`
	Rank             int    `json:"rank"`             // 1 = winner, 2 = second, etc.
}

// NectarScoring represents nectar points awarded per habitat
type NectarScoring struct {
	Forest     map[string]int `json:"forest"`     // playerName -> points
	Grassland  map[string]int `json:"grassland"`  // playerName -> points
	Wetland    map[string]int `json:"wetland"`    // playerName -> points
}

// CalculateGameEndScores calculates all player scores including nectar competitive scoring
func CalculateGameEndScores(players []PlayerGameEnd, includeOceania bool) ([]PlayerGameEnd, NectarScoring) {
	nectarScoring := NectarScoring{
		Forest:    make(map[string]int),
		Grassland: make(map[string]int),
		Wetland:   make(map[string]int),
	}

	// Calculate nectar points if Oceania is included
	if includeOceania {
		nectarScoring = calculateNectarPoints(players)
	}

	// Calculate totals for each player
	for i := range players {
		players[i].Total = players[i].BirdPoints +
			players[i].BonusCards +
			players[i].RoundGoals +
			players[i].Eggs +
			players[i].CachedFood +
			players[i].TuckedCards

		// Add nectar points
		if includeOceania {
			players[i].Total += nectarScoring.Forest[players[i].PlayerName]
			players[i].Total += nectarScoring.Grassland[players[i].PlayerName]
			players[i].Total += nectarScoring.Wetland[players[i].PlayerName]
		}
	}

	// Determine rankings with tiebreaker
	determineRankings(players)

	return players, nectarScoring
}

// calculateNectarPoints calculates competitive nectar scoring (1st=5pts, 2nd=2pts per habitat)
func calculateNectarPoints(players []PlayerGameEnd) NectarScoring {
	scoring := NectarScoring{
		Forest:    make(map[string]int),
		Grassland: make(map[string]int),
		Wetland:   make(map[string]int),
	}

	// Score each habitat separately
	scoring.Forest = scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarForest })
	scoring.Grassland = scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarGrassland })
	scoring.Wetland = scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarWetland })

	return scoring
}

// scoreHabitat scores a single habitat using competitive nectar rules with official tie-breaking
// Tie Rule: Players who tie split and average their placement points (rounded down)
// Example: 2 players tie for 1st place split (5+2)/2 = 3 points each
func scoreHabitat(players []PlayerGameEnd, getNectar func(PlayerGameEnd) int) map[string]int {
	points := make(map[string]int)

	// Group players by nectar count
	type playerCount struct {
		name  string
		count int
	}
	var counts []playerCount
	for _, p := range players {
		count := getNectar(p)
		if count > 0 { // Only consider players with nectar
			counts = append(counts, playerCount{name: p.PlayerName, count: count})
		}
	}

	if len(counts) == 0 {
		return points // No one has nectar in this habitat
	}

	// Sort by count descending
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].count > counts[j].count
	})

	// Available point values for places: [5, 2, 0, 0, ...] (1st=5, 2nd=2, rest=0)
	placeValues := []int{5, 2}

	position := 0
	for position < len(counts) {
		currentCount := counts[position].count

		// Find all players tied at this position
		tiedPlayers := []string{}
		for i := position; i < len(counts) && counts[i].count == currentCount; i++ {
			tiedPlayers = append(tiedPlayers, counts[i].name)
		}

		// Calculate the average of the place values for tied positions
		numTied := len(tiedPlayers)
		totalPoints := 0
		for i := 0; i < numTied; i++ {
			if position+i < len(placeValues) {
				totalPoints += placeValues[position+i]
			}
			// After 2nd place (index 1), all values are 0
		}

		// Split and average (rounded down)
		avgPoints := totalPoints / numTied

		// Award points to all tied players
		for _, playerName := range tiedPlayers {
			points[playerName] = avgPoints
		}

		// Move to next position
		position += numTied
	}

	return points
}

// determineRankings assigns ranks to players based on total score and tiebreakers
func determineRankings(players []PlayerGameEnd) {
	// Sort players by total (descending), then by unused food (descending) for tiebreaker
	sort.Slice(players, func(i, j int) bool {
		if players[i].Total == players[j].Total {
			return players[i].UnusedFood > players[j].UnusedFood
		}
		return players[i].Total > players[j].Total
	})

	// Assign ranks (handle ties - players with same total and unused food share rank)
	currentRank := 1
	for i := range players {
		if i > 0 {
			// Check if this player ties with previous
			if players[i].Total != players[i-1].Total || players[i].UnusedFood != players[i-1].UnusedFood {
				currentRank = i + 1
			}
		}
		players[i].Rank = currentRank
	}
}
