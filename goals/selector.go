package goals

import (
	"crypto/rand"
	"math/big"
)

// RoundGoals represents the 4 goals selected for a game
type RoundGoals struct {
	Round1 Goal `json:"round1"`
	Round2 Goal `json:"round2"`
	Round3 Goal `json:"round3"`
	Round4 Goal `json:"round4"`
}

// SelectRandomGoals randomly selects 4 unique goals from the available pool
func SelectRandomGoals(availableGoals []Goal) (RoundGoals, error) {
	if len(availableGoals) < 4 {
		// If we don't have enough goals, just use what we have
		// This shouldn't happen in normal use, but handle it gracefully
		result := RoundGoals{}
		for i := 0; i < len(availableGoals) && i < 4; i++ {
			switch i {
			case 0:
				result.Round1 = availableGoals[i]
			case 1:
				result.Round2 = availableGoals[i]
			case 2:
				result.Round3 = availableGoals[i]
			case 3:
				result.Round4 = availableGoals[i]
			}
		}
		return result, nil
	}

	// Create a copy to avoid modifying the original slice
	goalsCopy := make([]Goal, len(availableGoals))
	copy(goalsCopy, availableGoals)

	// Fisher-Yates shuffle
	for i := len(goalsCopy) - 1; i > 0; i-- {
		// Generate cryptographically secure random index
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return RoundGoals{}, err
		}
		j := nBig.Int64()

		// Swap
		goalsCopy[i], goalsCopy[j] = goalsCopy[j], goalsCopy[i]
	}

	// Return the first 4 shuffled goals
	return RoundGoals{
		Round1: goalsCopy[0],
		Round2: goalsCopy[1],
		Round3: goalsCopy[2],
		Round4: goalsCopy[3],
	}, nil
}
