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

// TestCalculateGameEndScores_PlayerCountMatrix runs a complete end-game calculation for
// 2, 3, and 4 players and checks the total, the ranking order, and the nectar award in
// each habitat. Nectar places are worth 5 (1st) and 2 (2nd); ties split the sum of the
// places they occupy, rounded down.
func TestCalculateGameEndScores_PlayerCountMatrix(t *testing.T) {
	tests := []struct {
		name        string
		players     []PlayerGameEnd
		wantTotals  map[string]int
		wantRanks   map[string]int
		wantForest  map[string]int
		wantWetland map[string]int
	}{
		{
			name: "2 players",
			players: []PlayerGameEnd{
				// 30+10+15+12+3+5 = 75, +5 forest +5 wetland = 85
				{PlayerName: "Alice", BirdPoints: 30, BonusCards: 10, RoundGoals: 15, Eggs: 12, CachedFood: 3, TuckedCards: 5, NectarForest: 6, NectarWetland: 4, UnusedFood: 2},
				// 28+8+11+14+2+4 = 67, +2 forest +2 wetland = 71
				{PlayerName: "Bob", BirdPoints: 28, BonusCards: 8, RoundGoals: 11, Eggs: 14, CachedFood: 2, TuckedCards: 4, NectarForest: 3, NectarWetland: 1, UnusedFood: 5},
			},
			wantTotals:  map[string]int{"Alice": 85, "Bob": 71},
			wantRanks:   map[string]int{"Alice": 1, "Bob": 2},
			wantForest:  map[string]int{"Alice": 5, "Bob": 2},
			wantWetland: map[string]int{"Alice": 5, "Bob": 2},
		},
		{
			name: "3 players",
			players: []PlayerGameEnd{
				// 40+12+18+10+0+6 = 86, +5 forest +0 wetland = 91
				{PlayerName: "Alice", BirdPoints: 40, BonusCards: 12, RoundGoals: 18, Eggs: 10, CachedFood: 0, TuckedCards: 6, NectarForest: 9, NectarWetland: 0, UnusedFood: 1},
				// 35+9+14+16+4+2 = 80, +2 forest +5 wetland = 87
				{PlayerName: "Bob", BirdPoints: 35, BonusCards: 9, RoundGoals: 14, Eggs: 16, CachedFood: 4, TuckedCards: 2, NectarForest: 5, NectarWetland: 7, UnusedFood: 3},
				// 31+7+9+11+1+3 = 62, +0 forest +2 wetland = 64
				{PlayerName: "Carol", BirdPoints: 31, BonusCards: 7, RoundGoals: 9, Eggs: 11, CachedFood: 1, TuckedCards: 3, NectarForest: 2, NectarWetland: 4, UnusedFood: 0},
			},
			wantTotals:  map[string]int{"Alice": 91, "Bob": 87, "Carol": 64},
			wantRanks:   map[string]int{"Alice": 1, "Bob": 2, "Carol": 3},
			wantForest:  map[string]int{"Alice": 5, "Bob": 2, "Carol": 0},
			wantWetland: map[string]int{"Bob": 5, "Carol": 2},
		},
		{
			name: "4 players with a three-way nectar tie for second",
			players: []PlayerGameEnd{
				// 44+14+20+13+2+7 = 100, +5 forest +5 wetland = 110
				{PlayerName: "Alice", BirdPoints: 44, BonusCards: 14, RoundGoals: 20, Eggs: 13, CachedFood: 2, TuckedCards: 7, NectarForest: 10, NectarWetland: 6, UnusedFood: 4},
				// 38+11+16+15+3+5 = 88, +0 forest +0 wetland = 88
				{PlayerName: "Bob", BirdPoints: 38, BonusCards: 11, RoundGoals: 16, Eggs: 15, CachedFood: 3, TuckedCards: 5, NectarForest: 4, NectarWetland: 2, UnusedFood: 2},
				// 33+10+12+12+1+4 = 72, +0 forest +0 wetland = 72
				{PlayerName: "Carol", BirdPoints: 33, BonusCards: 10, RoundGoals: 12, Eggs: 12, CachedFood: 1, TuckedCards: 4, NectarForest: 4, NectarWetland: 2, UnusedFood: 6},
				// 29+6+8+9+0+2 = 54, +0 forest +0 wetland = 54
				{PlayerName: "Dave", BirdPoints: 29, BonusCards: 6, RoundGoals: 8, Eggs: 9, CachedFood: 0, TuckedCards: 2, NectarForest: 4, NectarWetland: 2, UnusedFood: 1},
			},
			wantTotals: map[string]int{"Alice": 110, "Bob": 88, "Carol": 72, "Dave": 54},
			wantRanks:  map[string]int{"Alice": 1, "Bob": 2, "Carol": 3, "Dave": 4},
			// Bob, Carol and Dave all tie for 2nd: (2+0+0)/3 = 0 each.
			wantForest:  map[string]int{"Alice": 5, "Bob": 0, "Carol": 0, "Dave": 0},
			wantWetland: map[string]int{"Alice": 5, "Bob": 0, "Carol": 0, "Dave": 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players, nectar := CalculateGameEndScores(tt.players, true)

			assert.Len(t, players, len(tt.wantTotals))

			gotTotals := make(map[string]int, len(players))
			gotRanks := make(map[string]int, len(players))
			for _, p := range players {
				gotTotals[p.PlayerName] = p.Total
				gotRanks[p.PlayerName] = p.Rank
			}
			assert.Equal(t, tt.wantTotals, gotTotals)
			assert.Equal(t, tt.wantRanks, gotRanks)

			for name, want := range tt.wantForest {
				assert.Equal(t, want, nectar.Forest[name], "forest nectar for %s", name)
			}
			for name, want := range tt.wantWetland {
				assert.Equal(t, want, nectar.Wetland[name], "wetland nectar for %s", name)
			}

			// CalculateGameEndScores sorts in place, so the returned slice must
			// already be in ranking order.
			for i := 1; i < len(players); i++ {
				assert.LessOrEqual(t, players[i-1].Rank, players[i].Rank,
					"players are not returned in ranking order")
			}
		})
	}
}

// TestDetermineRankings_ThreeAndFourPlayerTies covers rank assignment when 3 or 4 players
// are involved in a tie, including the unused-food tiebreaker.
func TestDetermineRankings_ThreeAndFourPlayerTies(t *testing.T) {
	t.Run("3 players, tie for first broken by unused food", func(t *testing.T) {
		players := []PlayerGameEnd{
			{PlayerName: "Alice", Total: 90, UnusedFood: 2},
			{PlayerName: "Bob", Total: 90, UnusedFood: 5},
			{PlayerName: "Carol", Total: 80, UnusedFood: 9},
		}

		determineRankings(players)

		assert.Equal(t, "Bob", players[0].PlayerName, "more unused food wins the tiebreak")
		assert.Equal(t, 1, players[0].Rank)
		assert.Equal(t, "Alice", players[1].PlayerName)
		assert.Equal(t, 2, players[1].Rank)
		assert.Equal(t, "Carol", players[2].PlayerName)
		assert.Equal(t, 3, players[2].Rank)
	})

	t.Run("3 players, unbreakable tie for first shares rank 1", func(t *testing.T) {
		players := []PlayerGameEnd{
			{PlayerName: "Alice", Total: 90, UnusedFood: 3},
			{PlayerName: "Bob", Total: 90, UnusedFood: 3},
			{PlayerName: "Carol", Total: 70, UnusedFood: 1},
		}

		determineRankings(players)

		assert.Equal(t, 1, players[0].Rank)
		assert.Equal(t, 1, players[1].Rank)
		// Two players occupy 1st, so the next player is 3rd, not 2nd.
		assert.Equal(t, 3, players[2].Rank)
	})

	t.Run("4 players, tie for second shares rank 2 and pushes last to rank 4", func(t *testing.T) {
		players := []PlayerGameEnd{
			{PlayerName: "Alice", Total: 100, UnusedFood: 1},
			{PlayerName: "Bob", Total: 85, UnusedFood: 4},
			{PlayerName: "Carol", Total: 85, UnusedFood: 4},
			{PlayerName: "Dave", Total: 60, UnusedFood: 0},
		}

		determineRankings(players)

		ranks := make(map[string]int, len(players))
		for _, p := range players {
			ranks[p.PlayerName] = p.Rank
		}
		assert.Equal(t, map[string]int{"Alice": 1, "Bob": 2, "Carol": 2, "Dave": 4}, ranks)
	})

	t.Run("4 players, all tied share rank 1", func(t *testing.T) {
		players := []PlayerGameEnd{
			{PlayerName: "Alice", Total: 75, UnusedFood: 2},
			{PlayerName: "Bob", Total: 75, UnusedFood: 2},
			{PlayerName: "Carol", Total: 75, UnusedFood: 2},
			{PlayerName: "Dave", Total: 75, UnusedFood: 2},
		}

		determineRankings(players)

		for _, p := range players {
			assert.Equal(t, 1, p.Rank, "%s should share rank 1", p.PlayerName)
		}
	})
}

// TestScoreHabitat_ThreeAndFourPlayerShapes covers nectar placements that only occur once
// a game has more than two players.
func TestScoreHabitat_ThreeAndFourPlayerShapes(t *testing.T) {
	tests := []struct {
		name   string
		nectar map[string]int
		want   map[string]int
	}{
		// Two tie for 1st, third takes what is left of the places: (5+2)/2 = 3.
		{"3P tie for first", map[string]int{"Alice": 7, "Bob": 7, "Carol": 3},
			map[string]int{"Alice": 3, "Bob": 3, "Carol": 0}},
		// Two tie for 2nd: (2+0)/2 = 1.
		{"3P tie for second", map[string]int{"Alice": 7, "Bob": 3, "Carol": 3},
			map[string]int{"Alice": 5, "Bob": 1, "Carol": 1}},
		// All three tie: (5+2+0)/3 = 2.
		{"3P all tied", map[string]int{"Alice": 4, "Bob": 4, "Carol": 4},
			map[string]int{"Alice": 2, "Bob": 2, "Carol": 2}},
		{"4P no ties", map[string]int{"Alice": 9, "Bob": 7, "Carol": 5, "Dave": 3},
			map[string]int{"Alice": 5, "Bob": 2, "Carol": 0, "Dave": 0}},
		// Two tie for 1st: (5+2)/2 = 3; the rest get nothing.
		{"4P tie for first", map[string]int{"Alice": 9, "Bob": 9, "Carol": 5, "Dave": 3},
			map[string]int{"Alice": 3, "Bob": 3, "Carol": 0, "Dave": 0}},
		// All four tie: (5+2+0+0)/4 = 1.
		{"4P all tied", map[string]int{"Alice": 4, "Bob": 4, "Carol": 4, "Dave": 4},
			map[string]int{"Alice": 1, "Bob": 1, "Carol": 1, "Dave": 1}},
		// A player with no nectar in a habitat cannot place there, so 2nd goes to
		// the next player who does have some.
		{"4P two players have no nectar", map[string]int{"Alice": 6, "Bob": 4, "Carol": 0, "Dave": 0},
			map[string]int{"Alice": 5, "Bob": 2, "Carol": 0, "Dave": 0}},
		{"3P only one player has nectar", map[string]int{"Alice": 6, "Bob": 0, "Carol": 0},
			map[string]int{"Alice": 5, "Bob": 0, "Carol": 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := make([]PlayerGameEnd, 0, len(tt.nectar))
			for _, name := range []string{"Alice", "Bob", "Carol", "Dave"} {
				if count, ok := tt.nectar[name]; ok {
					players = append(players, PlayerGameEnd{PlayerName: name, NectarForest: count})
				}
			}

			points := scoreHabitat(players, func(p PlayerGameEnd) int { return p.NectarForest })

			for name, want := range tt.want {
				assert.Equal(t, want, points[name], "nectar points for %s", name)
			}

			// A habitat is worth 7 points in total (5 + 2); ties may only ever
			// award less through rounding, never more.
			total := 0
			for _, p := range points {
				total += p
			}
			assert.LessOrEqual(t, total, 7, "habitat awarded %d points, max is 7", total)
		})
	}
}
