package scoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCalculateGameEndScores_WithoutOceania tests basic game-end scoring without Oceania
func TestCalculateGameEndScores_WithoutOceania(t *testing.T) {
	players := []PlayerGameEnd{
		{
			PlayerName:  "Alice",
			BirdPoints:  50,
			BonusCards:  10,
			RoundGoals:  15,
			Eggs:        8,
			CachedFood:  3,
			TuckedCards: 5,
			UnusedFood:  2,
		},
		{
			PlayerName:  "Bob",
			BirdPoints:  45,
			BonusCards:  12,
			RoundGoals:  18,
			Eggs:        7,
			CachedFood:  4,
			TuckedCards: 3,
			UnusedFood:  3,
		},
	}

	result, nectarScoring := CalculateGameEndScores(players, false)

	// Alice total: 50+10+15+8+3+5 = 91
	assert.Equal(t, 91, result[0].Total)
	assert.Equal(t, "Alice", result[0].PlayerName)
	assert.Equal(t, 1, result[0].Rank)

	// Bob total: 45+12+18+7+4+3 = 89
	assert.Equal(t, 89, result[1].Total)
	assert.Equal(t, "Bob", result[1].PlayerName)
	assert.Equal(t, 2, result[1].Rank)

	// Nectar scoring should be empty
	assert.Empty(t, nectarScoring.Forest)
	assert.Empty(t, nectarScoring.Grassland)
	assert.Empty(t, nectarScoring.Wetland)
}

// TestCalculateGameEndScores_WithOceania tests game-end scoring with Oceania nectar
func TestCalculateGameEndScores_WithOceania(t *testing.T) {
	players := []PlayerGameEnd{
		{
			PlayerName:      "Alice",
			BirdPoints:      50,
			BonusCards:      10,
			RoundGoals:      15,
			Eggs:            8,
			CachedFood:      3,
			TuckedCards:     5,
			NectarForest:    3, // 1st place
			NectarGrassland: 2, // 2nd place
			NectarWetland:   1, // 3rd place
			UnusedFood:      2,
		},
		{
			PlayerName:      "Bob",
			BirdPoints:      45,
			BonusCards:      12,
			RoundGoals:      18,
			Eggs:            7,
			CachedFood:      4,
			TuckedCards:     3,
			NectarForest:    2, // 2nd place
			NectarGrassland: 3, // 1st place
			NectarWetland:   2, // 1st place
			UnusedFood:      3,
		},
	}

	result, nectarScoring := CalculateGameEndScores(players, true)

	// Check nectar scoring
	// Forest: Alice=3 (1st=5pts), Bob=2 (2nd=2pts)
	assert.Equal(t, 5, nectarScoring.Forest["Alice"])
	assert.Equal(t, 2, nectarScoring.Forest["Bob"])

	// Grassland: Bob=3 (1st=5pts), Alice=2 (2nd=2pts)
	assert.Equal(t, 2, nectarScoring.Grassland["Alice"])
	assert.Equal(t, 5, nectarScoring.Grassland["Bob"])

	// Wetland: Bob=2 (1st=5pts), Alice=1 (2nd=2pts)
	assert.Equal(t, 2, nectarScoring.Wetland["Alice"])
	assert.Equal(t, 5, nectarScoring.Wetland["Bob"])

	// Bob total: 45+12+18+7+4+3 + (2+5+5) = 89 + 12 = 101
	assert.Equal(t, 101, result[0].Total)
	assert.Equal(t, "Bob", result[0].PlayerName)
	assert.Equal(t, 1, result[0].Rank) // Bob (highest total - winner)

	// Alice total: 50+10+15+8+3+5 + (5+2+2) = 91 + 9 = 100
	assert.Equal(t, 100, result[1].Total)
	assert.Equal(t, "Alice", result[1].PlayerName)
	assert.Equal(t, 2, result[1].Rank) // Alice (second place)
}

// TestScoreHabitat_SingleWinner tests nectar scoring with clear 1st and 2nd place
func TestScoreHabitat_SingleWinner(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", NectarForest: 5},
		{PlayerName: "Bob", NectarForest: 3},
		{PlayerName: "Carol", NectarForest: 1},
	}

	points := scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarForest })

	assert.Equal(t, 5, points["Alice"]) // 1st place = 5 points
	assert.Equal(t, 2, points["Bob"])   // 2nd place = 2 points
	assert.Equal(t, 0, points["Carol"]) // 3rd place = 0 points
}

// TestScoreHabitat_TwoWayTieForFirst tests 2-player tie for 1st place
func TestScoreHabitat_TwoWayTieForFirst(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", NectarForest: 4},
		{PlayerName: "Bob", NectarForest: 4},
		{PlayerName: "Carol", NectarForest: 2},
	}

	points := scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarForest })

	// Alice and Bob tied for 1st: (5+2)/2 = 3 points each (integer division rounds down)
	assert.Equal(t, 3, points["Alice"])
	assert.Equal(t, 3, points["Bob"])
	assert.Equal(t, 0, points["Carol"]) // 3rd place = 0 points
}

// TestScoreHabitat_TwoWayTieForSecond tests 2-player tie for 2nd place
func TestScoreHabitat_TwoWayTieForSecond(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", NectarGrassland: 5},
		{PlayerName: "Bob", NectarGrassland: 3},
		{PlayerName: "Carol", NectarGrassland: 3},
	}

	points := scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarGrassland })

	assert.Equal(t, 5, points["Alice"]) // 1st place = 5 points
	// Bob and Carol tied for 2nd: (2+0)/2 = 1 point each
	assert.Equal(t, 1, points["Bob"])
	assert.Equal(t, 1, points["Carol"])
}

// TestScoreHabitat_ThreeWayTieForFirst tests 3-player tie for 1st place
func TestScoreHabitat_ThreeWayTieForFirst(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", NectarWetland: 6},
		{PlayerName: "Bob", NectarWetland: 6},
		{PlayerName: "Carol", NectarWetland: 6},
		{PlayerName: "Dave", NectarWetland: 2},
	}

	points := scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarWetland })

	// Three tied for 1st: (5+2+0)/3 = 7/3 = 2 points each (integer division)
	assert.Equal(t, 2, points["Alice"])
	assert.Equal(t, 2, points["Bob"])
	assert.Equal(t, 2, points["Carol"])
	assert.Equal(t, 0, points["Dave"]) // 4th place = 0 points
}

// TestScoreHabitat_AllPlayersTied tests when all players have same nectar count
func TestScoreHabitat_AllPlayersTied(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", NectarForest: 3},
		{PlayerName: "Bob", NectarForest: 3},
		{PlayerName: "Carol", NectarForest: 3},
	}

	points := scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarForest })

	// All tied for 1st: (5+2+0)/3 = 7/3 = 2 points each
	assert.Equal(t, 2, points["Alice"])
	assert.Equal(t, 2, points["Bob"])
	assert.Equal(t, 2, points["Carol"])
}

// TestScoreHabitat_ZeroNectarExcluded tests that players with 0 nectar don't get points
func TestScoreHabitat_ZeroNectarExcluded(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", NectarForest: 5},
		{PlayerName: "Bob", NectarForest: 3},
		{PlayerName: "Carol", NectarForest: 0}, // Should be excluded
		{PlayerName: "Dave", NectarForest: 0},  // Should be excluded
	}

	points := scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarForest })

	assert.Equal(t, 5, points["Alice"]) // 1st place
	assert.Equal(t, 2, points["Bob"])   // 2nd place

	// Players with 0 nectar should not be in the scoring map
	_, carolExists := points["Carol"]
	_, daveExists := points["Dave"]
	assert.False(t, carolExists)
	assert.False(t, daveExists)
}

// TestScoreHabitat_NoNectar tests habitat where no one has nectar
func TestScoreHabitat_NoNectar(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", NectarForest: 0},
		{PlayerName: "Bob", NectarForest: 0},
	}

	points := scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarForest })

	assert.Empty(t, points) // No points awarded
}

// TestScoreHabitat_SinglePlayer tests habitat with only one player having nectar
func TestScoreHabitat_SinglePlayer(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", NectarForest: 7},
		{PlayerName: "Bob", NectarForest: 0},
	}

	points := scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarForest })

	assert.Equal(t, 5, points["Alice"]) // 1st place = 5 points
	_, bobExists := points["Bob"]
	assert.False(t, bobExists) // Bob excluded
}

// TestCalculateNectarPoints_MultipleHabitats tests nectar scoring across all 3 habitats
func TestCalculateNectarPoints_MultipleHabitats(t *testing.T) {
	players := []PlayerGameEnd{
		{
			PlayerName:      "Alice",
			NectarForest:    5,
			NectarGrassland: 2,
			NectarWetland:   3,
		},
		{
			PlayerName:      "Bob",
			NectarForest:    3,
			NectarGrassland: 4,
			NectarWetland:   3,
		},
	}

	scoring := calculateNectarPoints(players)

	// Forest: Alice=5 (1st=5pts), Bob=3 (2nd=2pts)
	assert.Equal(t, 5, scoring.Forest["Alice"])
	assert.Equal(t, 2, scoring.Forest["Bob"])

	// Grassland: Bob=4 (1st=5pts), Alice=2 (2nd=2pts)
	assert.Equal(t, 2, scoring.Grassland["Alice"])
	assert.Equal(t, 5, scoring.Grassland["Bob"])

	// Wetland: Alice=3, Bob=3 (tied for 1st: (5+2)/2=3 each)
	assert.Equal(t, 3, scoring.Wetland["Alice"])
	assert.Equal(t, 3, scoring.Wetland["Bob"])
}

// TestDetermineRankings_NoTies tests ranking with clear winners
func TestDetermineRankings_NoTies(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Carol", Total: 85, UnusedFood: 2},
		{PlayerName: "Alice", Total: 100, UnusedFood: 3},
		{PlayerName: "Bob", Total: 90, UnusedFood: 5},
	}

	determineRankings(players)

	// Should be sorted by total score descending
	assert.Equal(t, "Alice", players[0].PlayerName)
	assert.Equal(t, 1, players[0].Rank)

	assert.Equal(t, "Bob", players[1].PlayerName)
	assert.Equal(t, 2, players[1].Rank)

	assert.Equal(t, "Carol", players[2].PlayerName)
	assert.Equal(t, 3, players[2].Rank)
}

// TestDetermineRankings_TieBrokenByUnusedFood tests tiebreaker mechanism
func TestDetermineRankings_TieBrokenByUnusedFood(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", Total: 100, UnusedFood: 3},
		{PlayerName: "Bob", Total: 100, UnusedFood: 5}, // Same total, more unused food
		{PlayerName: "Carol", Total: 90, UnusedFood: 2},
	}

	determineRankings(players)

	// Bob should win due to tiebreaker (more unused food)
	assert.Equal(t, "Bob", players[0].PlayerName)
	assert.Equal(t, 1, players[0].Rank)

	assert.Equal(t, "Alice", players[1].PlayerName)
	assert.Equal(t, 2, players[1].Rank)

	assert.Equal(t, "Carol", players[2].PlayerName)
	assert.Equal(t, 3, players[2].Rank)
}

// TestDetermineRankings_CompleteTie tests when players have identical totals and unused food
func TestDetermineRankings_CompleteTie(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", Total: 100, UnusedFood: 5},
		{PlayerName: "Bob", Total: 100, UnusedFood: 5}, // Complete tie
		{PlayerName: "Carol", Total: 90, UnusedFood: 3},
	}

	determineRankings(players)

	// Alice and Bob should share rank 1
	assert.Equal(t, 1, players[0].Rank)
	assert.Equal(t, 1, players[1].Rank)

	// Carol should be rank 3 (not 2, because two players tied for 1st)
	assert.Equal(t, "Carol", players[2].PlayerName)
	assert.Equal(t, 3, players[2].Rank)
}

// TestDetermineRankings_MultipleGroups tests multiple tied groups
func TestDetermineRankings_MultipleGroups(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", Total: 100, UnusedFood: 5},
		{PlayerName: "Bob", Total: 100, UnusedFood: 5}, // Tied with Alice for 1st
		{PlayerName: "Carol", Total: 90, UnusedFood: 3},
		{PlayerName: "Dave", Total: 90, UnusedFood: 3}, // Tied with Carol for 3rd
		{PlayerName: "Eve", Total: 80, UnusedFood: 2},
	}

	determineRankings(players)

	// Alice and Bob share rank 1
	assert.Equal(t, 1, players[0].Rank)
	assert.Equal(t, 1, players[1].Rank)

	// Carol and Dave share rank 3
	assert.Equal(t, 3, players[2].Rank)
	assert.Equal(t, 3, players[3].Rank)

	// Eve is rank 5
	assert.Equal(t, "Eve", players[4].PlayerName)
	assert.Equal(t, 5, players[4].Rank)
}

// TestDetermineRankings_SinglePlayer tests edge case with one player
func TestDetermineRankings_SinglePlayer(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", Total: 100, UnusedFood: 5},
	}

	determineRankings(players)

	assert.Equal(t, 1, players[0].Rank)
}

// TestCalculateGameEndScores_ComplexScenario tests a realistic multi-player game
func TestCalculateGameEndScores_ComplexScenario(t *testing.T) {
	players := []PlayerGameEnd{
		{
			PlayerName:      "Alice",
			BirdPoints:      48,
			BonusCards:      12,
			RoundGoals:      14,
			Eggs:            9,
			CachedFood:      5,
			TuckedCards:     6,
			NectarForest:    4,
			NectarGrassland: 3,
			NectarWetland:   2,
			UnusedFood:      3,
		},
		{
			PlayerName:      "Bob",
			BirdPoints:      52,
			BonusCards:      8,
			RoundGoals:      16,
			Eggs:            7,
			CachedFood:      4,
			TuckedCards:     4,
			NectarForest:    3,
			NectarGrassland: 4,
			NectarWetland:   5,
			UnusedFood:      5,
		},
		{
			PlayerName:      "Carol",
			BirdPoints:      45,
			BonusCards:      14,
			RoundGoals:      12,
			Eggs:            10,
			CachedFood:      3,
			TuckedCards:     7,
			NectarForest:    5,
			NectarGrassland: 2,
			NectarWetland:   3,
			UnusedFood:      2,
		},
	}

	result, nectarScoring := CalculateGameEndScores(players, true)

	// Verify nectar scoring
	// Forest: Carol=5 (1st=5), Alice=4 (2nd=2), Bob=3 (3rd=0)
	assert.Equal(t, 5, nectarScoring.Forest["Carol"])
	assert.Equal(t, 2, nectarScoring.Forest["Alice"])
	assert.Equal(t, 0, nectarScoring.Forest["Bob"])

	// Grassland: Bob=4 (1st=5), Alice=3 (2nd=2), Carol=2 (3rd=0)
	assert.Equal(t, 5, nectarScoring.Grassland["Bob"])
	assert.Equal(t, 2, nectarScoring.Grassland["Alice"])
	assert.Equal(t, 0, nectarScoring.Grassland["Carol"])

	// Wetland: Bob=5 (1st=5), Carol=3 (2nd=2), Alice=2 (3rd=0)
	assert.Equal(t, 5, nectarScoring.Wetland["Bob"])
	assert.Equal(t, 2, nectarScoring.Wetland["Carol"])
	assert.Equal(t, 0, nectarScoring.Wetland["Alice"])

	// Verify totals
	// Alice: 48+12+14+9+5+6 + (2+2+0) = 94 + 4 = 98
	// Bob: 52+8+16+7+4+4 + (0+5+5) = 91 + 10 = 101
	// Carol: 45+14+12+10+3+7 + (5+0+2) = 91 + 7 = 98

	// Find each player's result
	var aliceResult, bobResult, carolResult *PlayerGameEnd
	for i := range result {
		switch result[i].PlayerName {
		case "Alice":
			aliceResult = &result[i]
		case "Bob":
			bobResult = &result[i]
		case "Carol":
			carolResult = &result[i]
		}
	}

	assert.NotNil(t, aliceResult)
	assert.NotNil(t, bobResult)
	assert.NotNil(t, carolResult)

	assert.Equal(t, 98, aliceResult.Total)
	assert.Equal(t, 101, bobResult.Total)
	assert.Equal(t, 98, carolResult.Total)

	// Bob should be 1st with 101
	assert.Equal(t, 1, bobResult.Rank)

	// Alice and Carol both have 98, but Alice has more unused food (3 vs 2)
	assert.Equal(t, 2, aliceResult.Rank)
	assert.Equal(t, 3, carolResult.Rank)
}

// TestCalculateGameEndScores_EmptyPlayers tests edge case with no players
func TestCalculateGameEndScores_EmptyPlayers(t *testing.T) {
	players := []PlayerGameEnd{}

	result, nectarScoring := CalculateGameEndScores(players, true)

	assert.Len(t, result, 0)
	assert.Empty(t, nectarScoring.Forest)
	assert.Empty(t, nectarScoring.Grassland)
	assert.Empty(t, nectarScoring.Wetland)
}

// TestScoreHabitat_FourPlusPlayers tests scoring with more than 3 players
func TestScoreHabitat_FourPlusPlayers(t *testing.T) {
	players := []PlayerGameEnd{
		{PlayerName: "Alice", NectarForest: 8},
		{PlayerName: "Bob", NectarForest: 6},
		{PlayerName: "Carol", NectarForest: 4},
		{PlayerName: "Dave", NectarForest: 2},
		{PlayerName: "Eve", NectarForest: 1},
	}

	points := scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarForest })

	assert.Equal(t, 5, points["Alice"]) // 1st = 5
	assert.Equal(t, 2, points["Bob"])   // 2nd = 2
	assert.Equal(t, 0, points["Carol"]) // 3rd = 0
	assert.Equal(t, 0, points["Dave"])  // 4th = 0
	assert.Equal(t, 0, points["Eve"])   // 5th = 0
}
