package goals

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCalculateGreenScores_Round1_NoTies tests basic Round 1 scoring without ties
func TestCalculateGreenScores_Round1_NoTies(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 5,
		"Bob":   3,
		"Carol": 1,
	}

	scores := CalculateGreenScores(playerCounts, 1)

	// Round 1: 1st=4, 2nd=1, 3rd=0
	assert.Len(t, scores, 3)
	assert.Equal(t, "Alice", scores[0].PlayerName)
	assert.Equal(t, 5, scores[0].Count)
	assert.Equal(t, 4, scores[0].Points)
	assert.Equal(t, 1, scores[0].Rank)

	assert.Equal(t, "Bob", scores[1].PlayerName)
	assert.Equal(t, 3, scores[1].Count)
	assert.Equal(t, 1, scores[1].Points)
	assert.Equal(t, 2, scores[1].Rank)

	assert.Equal(t, "Carol", scores[2].PlayerName)
	assert.Equal(t, 1, scores[2].Count)
	assert.Equal(t, 0, scores[2].Points)
	assert.Equal(t, 3, scores[2].Rank)
}

// TestCalculateGreenScores_Round2_NoTies tests Round 2 scoring values
func TestCalculateGreenScores_Round2_NoTies(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 10,
		"Bob":   7,
		"Carol": 3,
	}

	scores := CalculateGreenScores(playerCounts, 2)

	// Round 2: 1st=5, 2nd=2, 3rd=0
	assert.Equal(t, 5, scores[0].Points)
	assert.Equal(t, 2, scores[1].Points)
	assert.Equal(t, 0, scores[2].Points)
}

// TestCalculateGreenScores_Round3_NoTies tests Round 3 scoring values
func TestCalculateGreenScores_Round3_NoTies(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 10,
		"Bob":   7,
		"Carol": 3,
	}

	scores := CalculateGreenScores(playerCounts, 3)

	// Round 3: 1st=6, 2nd=3, 3rd=2
	assert.Equal(t, 6, scores[0].Points)
	assert.Equal(t, 3, scores[1].Points)
	assert.Equal(t, 2, scores[2].Points)
}

// TestCalculateGreenScores_Round4_NoTies tests Round 4 scoring values
func TestCalculateGreenScores_Round4_NoTies(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 10,
		"Bob":   7,
		"Carol": 3,
	}

	scores := CalculateGreenScores(playerCounts, 4)

	// Round 4: 1st=7, 2nd=4, 3rd=2
	assert.Equal(t, 7, scores[0].Points)
	assert.Equal(t, 4, scores[1].Points)
	assert.Equal(t, 2, scores[2].Points)
}

// TestCalculateGreenScores_TwoWayTieForFirst tests 2-player tie for 1st place
func TestCalculateGreenScores_TwoWayTieForFirst(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 5,
		"Bob":   5,
		"Carol": 2,
	}

	scores := CalculateGreenScores(playerCounts, 1)

	// Round 1: 1st=4, 2nd=1
	// Two tied for 1st: (4+1)/2 = 2 points each (integer division rounds down)
	assert.Equal(t, "Alice", scores[0].PlayerName) // Alphabetical order
	assert.Equal(t, 2, scores[0].Points)
	assert.Equal(t, 1, scores[0].Rank)

	assert.Equal(t, "Bob", scores[1].PlayerName)
	assert.Equal(t, 2, scores[1].Points)
	assert.Equal(t, 1, scores[1].Rank)

	assert.Equal(t, "Carol", scores[2].PlayerName)
	assert.Equal(t, 0, scores[2].Points)
	assert.Equal(t, 3, scores[2].Rank)
}

// TestCalculateGreenScores_TwoWayTieForSecond tests 2-player tie for 2nd place
func TestCalculateGreenScores_TwoWayTieForSecond(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 7,
		"Bob":   4,
		"Carol": 4,
		"Dave":  1,
	}

	scores := CalculateGreenScores(playerCounts, 3)

	// Round 3: 1st=6, 2nd=3, 3rd=2
	// Alice: 1st place = 6 points
	// Bob & Carol tied for 2nd: (3+2)/2 = 2 points each
	// Dave: 4th place = 0 points
	assert.Equal(t, "Alice", scores[0].PlayerName)
	assert.Equal(t, 6, scores[0].Points)
	assert.Equal(t, 1, scores[0].Rank)

	assert.Equal(t, "Bob", scores[1].PlayerName)
	assert.Equal(t, 2, scores[1].Points)
	assert.Equal(t, 2, scores[1].Rank)

	assert.Equal(t, "Carol", scores[2].PlayerName)
	assert.Equal(t, 2, scores[2].Points)
	assert.Equal(t, 2, scores[2].Rank)

	assert.Equal(t, "Dave", scores[3].PlayerName)
	assert.Equal(t, 0, scores[3].Points)
	assert.Equal(t, 4, scores[3].Rank)
}

// TestCalculateGreenScores_ThreeWayTieForFirst tests 3-player tie for 1st place
func TestCalculateGreenScores_ThreeWayTieForFirst(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 8,
		"Bob":   8,
		"Carol": 8,
		"Dave":  2,
	}

	scores := CalculateGreenScores(playerCounts, 3)

	// Round 3: 1st=6, 2nd=3, 3rd=2
	// Three tied for 1st: (6+3+2)/3 = 11/3 = 3 points each (integer division)
	for i := 0; i < 3; i++ {
		assert.Equal(t, 8, scores[i].Count)
		assert.Equal(t, 3, scores[i].Points)
		assert.Equal(t, 1, scores[i].Rank)
	}

	// Verify alphabetical ordering
	assert.Equal(t, "Alice", scores[0].PlayerName)
	assert.Equal(t, "Bob", scores[1].PlayerName)
	assert.Equal(t, "Carol", scores[2].PlayerName)

	assert.Equal(t, "Dave", scores[3].PlayerName)
	assert.Equal(t, 0, scores[3].Points)
	assert.Equal(t, 4, scores[3].Rank)
}

// TestCalculateGreenScores_AllPlayersTied tests when all players have the same count
func TestCalculateGreenScores_AllPlayersTied(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 5,
		"Bob":   5,
		"Carol": 5,
	}

	scores := CalculateGreenScores(playerCounts, 1)

	// Round 1: 1st=4, 2nd=1, 3rd=0
	// All tied for 1st: (4+1+0)/3 = 5/3 = 1 point each
	for i := 0; i < 3; i++ {
		assert.Equal(t, 5, scores[i].Count)
		assert.Equal(t, 1, scores[i].Points)
		assert.Equal(t, 1, scores[i].Rank)
	}
}

// TestCalculateGreenScores_FourthPlaceAndBeyond tests that 4th+ place gets 0 points
func TestCalculateGreenScores_FourthPlaceAndBeyond(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 10,
		"Bob":   7,
		"Carol": 5,
		"Dave":  3,
		"Eve":   1,
	}

	scores := CalculateGreenScores(playerCounts, 1)

	// Round 1: 1st=4, 2nd=1, 3rd=0
	assert.Equal(t, 4, scores[0].Points)
	assert.Equal(t, 1, scores[1].Points)
	assert.Equal(t, 0, scores[2].Points)
	assert.Equal(t, 0, scores[3].Points) // 4th place = 0
	assert.Equal(t, 0, scores[4].Points) // 5th place = 0
}

// TestCalculateGreenScores_InvalidRound tests that invalid rounds default to round 1
func TestCalculateGreenScores_InvalidRound(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 5,
		"Bob":   3,
	}

	// Test round 0
	scores := CalculateGreenScores(playerCounts, 0)
	assert.Equal(t, 4, scores[0].Points) // Should use Round 1 scoring

	// Test round 5
	scores = CalculateGreenScores(playerCounts, 5)
	assert.Equal(t, 4, scores[0].Points) // Should use Round 1 scoring

	// Test negative round
	scores = CalculateGreenScores(playerCounts, -1)
	assert.Equal(t, 4, scores[0].Points) // Should use Round 1 scoring
}

// TestCalculateGreenScores_ZeroCounts tests players with zero counts
func TestCalculateGreenScores_ZeroCounts(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 5,
		"Bob":   0,
		"Carol": 0,
	}

	scores := CalculateGreenScores(playerCounts, 1)

	assert.Equal(t, "Alice", scores[0].PlayerName)
	assert.Equal(t, 4, scores[0].Points)

	// Bob and Carol tied for 2nd with 0: (1+0)/2 = 0 points each
	assert.Equal(t, "Bob", scores[1].PlayerName)
	assert.Equal(t, 0, scores[1].Count)
	assert.Equal(t, 0, scores[1].Points)
	assert.Equal(t, 2, scores[1].Rank)

	assert.Equal(t, "Carol", scores[2].PlayerName)
	assert.Equal(t, 0, scores[2].Count)
	assert.Equal(t, 0, scores[2].Points)
	assert.Equal(t, 2, scores[2].Rank)
}

// TestCalculateGreenScores_SinglePlayer tests edge case with only one player
func TestCalculateGreenScores_SinglePlayer(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 5,
	}

	scores := CalculateGreenScores(playerCounts, 1)

	assert.Len(t, scores, 1)
	assert.Equal(t, "Alice", scores[0].PlayerName)
	assert.Equal(t, 4, scores[0].Points)
	assert.Equal(t, 1, scores[0].Rank)
}

// TestCalculateGreenScores_EmptyInput tests edge case with no players
func TestCalculateGreenScores_EmptyInput(t *testing.T) {
	playerCounts := map[string]int{}

	scores := CalculateGreenScores(playerCounts, 1)

	assert.Len(t, scores, 0)
}

// TestCalculateBlueScores_NormalCounts tests basic Blue scoring
func TestCalculateBlueScores_NormalCounts(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 3,
		"Bob":   1,
		"Carol": 5,
	}

	scores := CalculateBlueScores(playerCounts)

	assert.Len(t, scores, 3)

	// Should be sorted by points descending (Carol=5, Alice=3, Bob=1)
	assert.Equal(t, "Carol", scores[0].PlayerName)
	assert.Equal(t, 5, scores[0].Count)
	assert.Equal(t, 5, scores[0].Points)

	assert.Equal(t, "Alice", scores[1].PlayerName)
	assert.Equal(t, 3, scores[1].Count)
	assert.Equal(t, 3, scores[1].Points)

	assert.Equal(t, "Bob", scores[2].PlayerName)
	assert.Equal(t, 1, scores[2].Count)
	assert.Equal(t, 1, scores[2].Points)
}

// TestCalculateBlueScores_MaximumCap tests that counts > 5 are capped at 5 points
func TestCalculateBlueScores_MaximumCap(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 10,
		"Bob":   7,
		"Carol": 100,
	}

	scores := CalculateBlueScores(playerCounts)

	// All should have 5 points (capped)
	for _, score := range scores {
		assert.Equal(t, 5, score.Points)
	}

	// When all have same points, sorted alphabetically
	assert.Equal(t, "Alice", scores[0].PlayerName)
	assert.Equal(t, 10, scores[0].Count) // Count preserved

	assert.Equal(t, "Bob", scores[1].PlayerName)
	assert.Equal(t, 7, scores[1].Count) // Count preserved

	assert.Equal(t, "Carol", scores[2].PlayerName)
	assert.Equal(t, 100, scores[2].Count) // Count preserved
}

// TestCalculateBlueScores_NegativeCounts tests that negative counts result in 0 points
func TestCalculateBlueScores_NegativeCounts(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 3,
		"Bob":   -5,
		"Carol": -1,
	}

	scores := CalculateBlueScores(playerCounts)

	assert.Equal(t, "Alice", scores[0].PlayerName)
	assert.Equal(t, 3, scores[0].Points)

	// Negative counts should result in 0 points
	for _, score := range scores {
		if score.PlayerName == "Bob" || score.PlayerName == "Carol" {
			assert.Equal(t, 0, score.Points)
		}
	}
}

// TestCalculateBlueScores_ZeroCounts tests players with zero counts
func TestCalculateBlueScores_ZeroCounts(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 0,
		"Bob":   0,
		"Carol": 2,
	}

	scores := CalculateBlueScores(playerCounts)

	assert.Equal(t, "Carol", scores[0].PlayerName)
	assert.Equal(t, 2, scores[0].Points)

	// Alice and Bob should have 0 points, alphabetically sorted
	assert.Equal(t, "Alice", scores[1].PlayerName)
	assert.Equal(t, 0, scores[1].Points)

	assert.Equal(t, "Bob", scores[2].PlayerName)
	assert.Equal(t, 0, scores[2].Points)
}

// TestCalculateBlueScores_AlphabeticalSorting tests that ties are sorted alphabetically
func TestCalculateBlueScores_AlphabeticalSorting(t *testing.T) {
	playerCounts := map[string]int{
		"Zara":  3,
		"Alice": 3,
		"Mike":  3,
	}

	scores := CalculateBlueScores(playerCounts)

	// All have same points, should be alphabetically sorted
	assert.Equal(t, "Alice", scores[0].PlayerName)
	assert.Equal(t, "Mike", scores[1].PlayerName)
	assert.Equal(t, "Zara", scores[2].PlayerName)
}

// TestCalculateBlueScores_SinglePlayer tests edge case with only one player
func TestCalculateBlueScores_SinglePlayer(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 4,
	}

	scores := CalculateBlueScores(playerCounts)

	assert.Len(t, scores, 1)
	assert.Equal(t, "Alice", scores[0].PlayerName)
	assert.Equal(t, 4, scores[0].Points)
}

// TestCalculateBlueScores_EmptyInput tests edge case with no players
func TestCalculateBlueScores_EmptyInput(t *testing.T) {
	playerCounts := map[string]int{}

	scores := CalculateBlueScores(playerCounts)

	assert.Len(t, scores, 0)
}

// TestCalculateBlueScores_BoundaryValues tests exact boundary at 5 points
func TestCalculateBlueScores_BoundaryValues(t *testing.T) {
	playerCounts := map[string]int{
		"Four": 4,
		"Five": 5,
		"Six":  6,
	}

	scores := CalculateBlueScores(playerCounts)

	// Find each player in results
	for _, score := range scores {
		switch score.PlayerName {
		case "Four":
			assert.Equal(t, 4, score.Points)
		case "Five":
			assert.Equal(t, 5, score.Points) // Exactly at boundary
		case "Six":
			assert.Equal(t, 5, score.Points) // Capped
			assert.Equal(t, 6, score.Count)  // Count preserved
		}
	}
}

// TestCalculateGreenScores_ZeroCountsRound3 tests that players with zero count get zero points
// This is a regression test for: https://github.com/morey-tech/wingspan-scoring/issues/80
// Bug: In rounds 3 and 4, players with 0 count were getting placement-based points instead of 0
func TestCalculateGreenScores_ZeroCountsRound3(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 5,
		"Bob":   0, // Should get 0 points, not placement points
		"Carol": 0, // Should get 0 points, not placement points
	}

	scores := CalculateGreenScores(playerCounts, 3)

	// Round 3: 1st=6, 2nd=3, 3rd=2
	// Alice should get 1st place points
	assert.Equal(t, "Alice", scores[0].PlayerName)
	assert.Equal(t, 6, scores[0].Points)
	assert.Equal(t, 1, scores[0].Rank)

	// Bob and Carol have 0 count, so they should ALWAYS get 0 points
	// regardless of their placement ranking
	for i := 1; i < len(scores); i++ {
		if scores[i].Count == 0 {
			assert.Equal(t, 0, scores[i].Points,
				"Player %s with 0 count should get 0 points, not placement-based points",
				scores[i].PlayerName)
		}
	}
}

// TestCalculateGreenScores_ZeroCountsRound4 tests zero count behavior in round 4
func TestCalculateGreenScores_ZeroCountsRound4(t *testing.T) {
	playerCounts := map[string]int{
		"Alice": 8,
		"Bob":   3,
		"Carol": 0, // Should get 0 points
	}

	scores := CalculateGreenScores(playerCounts, 4)

	// Round 4: 1st=7, 2nd=4, 3rd=2
	assert.Equal(t, "Alice", scores[0].PlayerName)
	assert.Equal(t, 7, scores[0].Points)

	assert.Equal(t, "Bob", scores[1].PlayerName)
	assert.Equal(t, 4, scores[1].Points)

	// Carol has 0 count, should get 0 points (not 2 for 3rd place)
	assert.Equal(t, "Carol", scores[2].PlayerName)
	assert.Equal(t, 0, scores[2].Count)
	assert.Equal(t, 0, scores[2].Points,
		"Player with 0 count should get 0 points, not 3rd place points")
}

// TestCalculateGreenScores_PlayerCountMatrix exercises green-side scoring across the
// full 2/3/4-player x round x tie-shape matrix.
//
// Expected points follow the official Wingspan tie rule: players who tie add together
// the points for every place they collectively occupy and divide evenly, rounding down.
// Places beyond 3rd are worth 0 in every round.
//
//	Round 1: 1st=4 2nd=1 3rd=0    Round 3: 1st=6 2nd=3 3rd=2
//	Round 2: 1st=5 2nd=2 3rd=0    Round 4: 1st=7 2nd=4 3rd=2
func TestCalculateGreenScores_PlayerCountMatrix(t *testing.T) {
	tests := []struct {
		name   string
		round  int
		counts map[string]int
		// want maps player name -> expected points
		want map[string]int
	}{
		// --- 2 players (verified-good baseline) ---
		{"2P/R1/no ties", 1, map[string]int{"Alice": 5, "Bob": 3}, map[string]int{"Alice": 4, "Bob": 1}},
		{"2P/R2/no ties", 2, map[string]int{"Alice": 5, "Bob": 3}, map[string]int{"Alice": 5, "Bob": 2}},
		{"2P/R3/no ties", 3, map[string]int{"Alice": 5, "Bob": 3}, map[string]int{"Alice": 6, "Bob": 3}},
		{"2P/R4/no ties", 4, map[string]int{"Alice": 5, "Bob": 3}, map[string]int{"Alice": 7, "Bob": 4}},
		// (4+1)/2, (5+2)/2, (6+3)/2, (7+4)/2
		{"2P/R1/tie for 1st", 1, map[string]int{"Alice": 5, "Bob": 5}, map[string]int{"Alice": 2, "Bob": 2}},
		{"2P/R2/tie for 1st", 2, map[string]int{"Alice": 5, "Bob": 5}, map[string]int{"Alice": 3, "Bob": 3}},
		{"2P/R3/tie for 1st", 3, map[string]int{"Alice": 5, "Bob": 5}, map[string]int{"Alice": 4, "Bob": 4}},
		{"2P/R4/tie for 1st", 4, map[string]int{"Alice": 5, "Bob": 5}, map[string]int{"Alice": 5, "Bob": 5}},

		// --- 3 players ---
		{"3P/R1/no ties", 1, map[string]int{"Alice": 5, "Bob": 3, "Carol": 1}, map[string]int{"Alice": 4, "Bob": 1, "Carol": 0}},
		{"3P/R2/no ties", 2, map[string]int{"Alice": 5, "Bob": 3, "Carol": 1}, map[string]int{"Alice": 5, "Bob": 2, "Carol": 0}},
		{"3P/R3/no ties", 3, map[string]int{"Alice": 5, "Bob": 3, "Carol": 1}, map[string]int{"Alice": 6, "Bob": 3, "Carol": 2}},
		{"3P/R4/no ties", 4, map[string]int{"Alice": 5, "Bob": 3, "Carol": 1}, map[string]int{"Alice": 7, "Bob": 4, "Carol": 2}},

		// Two tie for 1st, third player takes 3rd place (2nd is consumed by the tie).
		{"3P/R1/tie for 1st", 1, map[string]int{"Alice": 5, "Bob": 5, "Carol": 1}, map[string]int{"Alice": 2, "Bob": 2, "Carol": 0}},
		{"3P/R2/tie for 1st", 2, map[string]int{"Alice": 5, "Bob": 5, "Carol": 1}, map[string]int{"Alice": 3, "Bob": 3, "Carol": 0}},
		{"3P/R3/tie for 1st", 3, map[string]int{"Alice": 5, "Bob": 5, "Carol": 1}, map[string]int{"Alice": 4, "Bob": 4, "Carol": 2}},
		{"3P/R4/tie for 1st", 4, map[string]int{"Alice": 5, "Bob": 5, "Carol": 1}, map[string]int{"Alice": 5, "Bob": 5, "Carol": 2}},

		// Two tie for 2nd: they split 2nd + 3rd.
		{"3P/R1/tie for 2nd", 1, map[string]int{"Alice": 5, "Bob": 3, "Carol": 3}, map[string]int{"Alice": 4, "Bob": 0, "Carol": 0}},
		{"3P/R2/tie for 2nd", 2, map[string]int{"Alice": 5, "Bob": 3, "Carol": 3}, map[string]int{"Alice": 5, "Bob": 1, "Carol": 1}},
		{"3P/R3/tie for 2nd", 3, map[string]int{"Alice": 5, "Bob": 3, "Carol": 3}, map[string]int{"Alice": 6, "Bob": 2, "Carol": 2}},
		{"3P/R4/tie for 2nd", 4, map[string]int{"Alice": 5, "Bob": 3, "Carol": 3}, map[string]int{"Alice": 7, "Bob": 3, "Carol": 3}},

		// All three tie: split 1st + 2nd + 3rd.
		{"3P/R1/all tied", 1, map[string]int{"Alice": 4, "Bob": 4, "Carol": 4}, map[string]int{"Alice": 1, "Bob": 1, "Carol": 1}},
		{"3P/R2/all tied", 2, map[string]int{"Alice": 4, "Bob": 4, "Carol": 4}, map[string]int{"Alice": 2, "Bob": 2, "Carol": 2}},
		{"3P/R3/all tied", 3, map[string]int{"Alice": 4, "Bob": 4, "Carol": 4}, map[string]int{"Alice": 3, "Bob": 3, "Carol": 3}},
		{"3P/R4/all tied", 4, map[string]int{"Alice": 4, "Bob": 4, "Carol": 4}, map[string]int{"Alice": 4, "Bob": 4, "Carol": 4}},

		// --- 4 players ---
		{"4P/R1/no ties", 1, map[string]int{"Alice": 7, "Bob": 5, "Carol": 3, "Dave": 1}, map[string]int{"Alice": 4, "Bob": 1, "Carol": 0, "Dave": 0}},
		{"4P/R2/no ties", 2, map[string]int{"Alice": 7, "Bob": 5, "Carol": 3, "Dave": 1}, map[string]int{"Alice": 5, "Bob": 2, "Carol": 0, "Dave": 0}},
		{"4P/R3/no ties", 3, map[string]int{"Alice": 7, "Bob": 5, "Carol": 3, "Dave": 1}, map[string]int{"Alice": 6, "Bob": 3, "Carol": 2, "Dave": 0}},
		{"4P/R4/no ties", 4, map[string]int{"Alice": 7, "Bob": 5, "Carol": 3, "Dave": 1}, map[string]int{"Alice": 7, "Bob": 4, "Carol": 2, "Dave": 0}},

		// Tie for 1st pushes the remaining players to 3rd and 4th.
		{"4P/R3/tie for 1st", 3, map[string]int{"Alice": 7, "Bob": 7, "Carol": 3, "Dave": 1}, map[string]int{"Alice": 4, "Bob": 4, "Carol": 2, "Dave": 0}},
		{"4P/R4/tie for 1st", 4, map[string]int{"Alice": 7, "Bob": 7, "Carol": 3, "Dave": 1}, map[string]int{"Alice": 5, "Bob": 5, "Carol": 2, "Dave": 0}},

		// Three tie for 2nd: split 2nd + 3rd + 4th.
		{"4P/R3/three tie for 2nd", 3, map[string]int{"Alice": 7, "Bob": 4, "Carol": 4, "Dave": 4}, map[string]int{"Alice": 6, "Bob": 1, "Carol": 1, "Dave": 1}},
		{"4P/R4/three tie for 2nd", 4, map[string]int{"Alice": 7, "Bob": 4, "Carol": 4, "Dave": 4}, map[string]int{"Alice": 7, "Bob": 2, "Carol": 2, "Dave": 2}},

		// Two tie for 3rd: split 3rd + 4th.
		{"4P/R3/tie for 3rd", 3, map[string]int{"Alice": 7, "Bob": 5, "Carol": 3, "Dave": 3}, map[string]int{"Alice": 6, "Bob": 3, "Carol": 1, "Dave": 1}},
		{"4P/R4/tie for 3rd", 4, map[string]int{"Alice": 7, "Bob": 5, "Carol": 3, "Dave": 3}, map[string]int{"Alice": 7, "Bob": 4, "Carol": 1, "Dave": 1}},

		// Two separate tie groups in one round.
		{"4P/R3/tie 1st and tie 3rd", 3, map[string]int{"Alice": 7, "Bob": 7, "Carol": 3, "Dave": 3}, map[string]int{"Alice": 4, "Bob": 4, "Carol": 1, "Dave": 1}},

		// All four tie: split 1st + 2nd + 3rd + 4th.
		{"4P/R1/all tied", 1, map[string]int{"Alice": 4, "Bob": 4, "Carol": 4, "Dave": 4}, map[string]int{"Alice": 1, "Bob": 1, "Carol": 1, "Dave": 1}},
		{"4P/R2/all tied", 2, map[string]int{"Alice": 4, "Bob": 4, "Carol": 4, "Dave": 4}, map[string]int{"Alice": 1, "Bob": 1, "Carol": 1, "Dave": 1}},
		{"4P/R3/all tied", 3, map[string]int{"Alice": 4, "Bob": 4, "Carol": 4, "Dave": 4}, map[string]int{"Alice": 2, "Bob": 2, "Carol": 2, "Dave": 2}},
		{"4P/R4/all tied", 4, map[string]int{"Alice": 4, "Bob": 4, "Carol": 4, "Dave": 4}, map[string]int{"Alice": 3, "Bob": 3, "Carol": 3, "Dave": 3}},

		// A player with none of the goal item does not qualify for a place, so the
		// places below them go unawarded rather than sliding up.
		{"3P/R3/one player has none", 3, map[string]int{"Alice": 5, "Bob": 3, "Carol": 0}, map[string]int{"Alice": 6, "Bob": 3, "Carol": 0}},
		{"4P/R4/two players have none", 4, map[string]int{"Alice": 5, "Bob": 3, "Carol": 0, "Dave": 0}, map[string]int{"Alice": 7, "Bob": 4, "Carol": 0, "Dave": 0}},
		{"4P/R3/only one player qualifies", 3, map[string]int{"Alice": 5, "Bob": 0, "Carol": 0, "Dave": 0}, map[string]int{"Alice": 6, "Bob": 0, "Carol": 0, "Dave": 0}},
		{"4P/R4/nobody qualifies", 4, map[string]int{"Alice": 0, "Bob": 0, "Carol": 0, "Dave": 0}, map[string]int{"Alice": 0, "Bob": 0, "Carol": 0, "Dave": 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scores := CalculateGreenScores(tt.counts, tt.round)

			assert.Len(t, scores, len(tt.counts))

			got := make(map[string]int, len(scores))
			for _, s := range scores {
				got[s.PlayerName] = s.Points
			}
			assert.Equal(t, tt.want, got)

			// Total points awarded may never exceed what the goal mat offers for
			// the round, no matter how the players tie.
			maxAvailable := 0
			for place := 1; place <= 3 && place <= len(tt.counts); place++ {
				maxAvailable += greenScoringRules[tt.round][place]
			}
			total := 0
			for _, p := range got {
				total += p
			}
			assert.LessOrEqual(t, total, maxAvailable,
				"awarded %d points but round %d only offers %d across %d players",
				total, tt.round, maxAvailable, len(tt.counts))
		})
	}
}

// TestCalculateGreenScores_BoardValuesAsCounts covers the translation the browser
// actually performs: it sends the printed value of the goal-mat box a player's cube sits
// on as that player's "count", relying on the backend to re-derive places from the
// ordering. Placing cubes on the correct boxes must reproduce the printed points.
func TestCalculateGreenScores_BoardValuesAsCounts(t *testing.T) {
	// Printed value of each place box on the green goal mat, by round.
	board := map[int][]int{
		1: {4, 1, 0, 0},
		2: {5, 2, 0, 0},
		3: {6, 3, 2, 0},
		4: {7, 4, 2, 0},
	}
	names := []string{"Alice", "Bob", "Carol", "Dave"}

	for numPlayers := 2; numPlayers <= 4; numPlayers++ {
		for round := 1; round <= 4; round++ {
			t.Run(fmt.Sprintf("%dP/R%d", numPlayers, round), func(t *testing.T) {
				counts := make(map[string]int, numPlayers)
				for i := 0; i < numPlayers; i++ {
					counts[names[i]] = board[round][i]
				}

				scores := CalculateGreenScores(counts, round)

				for i, s := range scores {
					want := board[round][i]
					assert.Equal(t, want, s.Points,
						"%d-player round %d: cube on the %d-point box scored %d",
						numPlayers, round, want, s.Points)

					// Rank is only recoverable when the printed value identifies
					// the box uniquely. In rounds 1 and 2 the 3rd-place and
					// 4th-place boxes are both printed 0, so a 4-player game
					// reports both players at rank 3. Harmless -- both boxes are
					// worth 0 points -- but the distinction is genuinely lost.
					if occurrences(board[round][:numPlayers], want) == 1 {
						assert.Equal(t, i+1, s.Rank)
					}
				}
			})
		}
	}
}

// occurrences counts how many times v appears in vals.
func occurrences(vals []int, v int) int {
	n := 0
	for _, x := range vals {
		if x == v {
			n++
		}
	}
	return n
}
