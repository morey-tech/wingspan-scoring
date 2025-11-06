package goals

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSelectRandomGoals_NormalCase tests selection with more than 4 goals
func TestSelectRandomGoals_NormalCase(t *testing.T) {
	availableGoals := GetAllGoals(true, false, false) // 16 base goals

	result, err := SelectRandomGoals(availableGoals)
	require.NoError(t, err)

	// Verify we got 4 goals
	goals := []Goal{result.Round1, result.Round2, result.Round3, result.Round4}
	assert.Len(t, goals, 4)

	// Verify all goals have non-empty IDs
	for i, goal := range goals {
		assert.NotEmpty(t, goal.ID, "Round %d goal should have an ID", i+1)
	}
}

// TestSelectRandomGoals_Uniqueness tests that selected goals are unique
func TestSelectRandomGoals_Uniqueness(t *testing.T) {
	availableGoals := GetAllGoals(true, false, false) // 16 base goals

	result, err := SelectRandomGoals(availableGoals)
	require.NoError(t, err)

	// Collect all IDs
	ids := map[string]bool{
		result.Round1.ID: true,
		result.Round2.ID: true,
		result.Round3.ID: true,
		result.Round4.ID: true,
	}

	// Verify all 4 IDs are unique
	assert.Len(t, ids, 4, "All 4 selected goals should be unique")
}

// TestSelectRandomGoals_ExactlyFourGoals tests selection when exactly 4 goals available
func TestSelectRandomGoals_ExactlyFourGoals(t *testing.T) {
	availableGoals := []Goal{
		{ID: "goal1", Name: "Goal 1", Expansion: "base"},
		{ID: "goal2", Name: "Goal 2", Expansion: "base"},
		{ID: "goal3", Name: "Goal 3", Expansion: "base"},
		{ID: "goal4", Name: "Goal 4", Expansion: "base"},
	}

	result, err := SelectRandomGoals(availableGoals)
	require.NoError(t, err)

	// Collect selected IDs
	selectedIDs := map[string]bool{
		result.Round1.ID: true,
		result.Round2.ID: true,
		result.Round3.ID: true,
		result.Round4.ID: true,
	}

	// All 4 should be selected (though order may vary)
	assert.Len(t, selectedIDs, 4)
	assert.True(t, selectedIDs["goal1"])
	assert.True(t, selectedIDs["goal2"])
	assert.True(t, selectedIDs["goal3"])
	assert.True(t, selectedIDs["goal4"])
}

// TestSelectRandomGoals_ThreeGoals tests selection with only 3 goals available
func TestSelectRandomGoals_ThreeGoals(t *testing.T) {
	availableGoals := []Goal{
		{ID: "goal1", Name: "Goal 1", Expansion: "base"},
		{ID: "goal2", Name: "Goal 2", Expansion: "base"},
		{ID: "goal3", Name: "Goal 3", Expansion: "base"},
	}

	result, err := SelectRandomGoals(availableGoals)
	require.NoError(t, err)

	// Should return the 3 goals in first 3 rounds
	assert.NotEmpty(t, result.Round1.ID)
	assert.NotEmpty(t, result.Round2.ID)
	assert.NotEmpty(t, result.Round3.ID)

	// Round 4 should be empty (zero value)
	assert.Empty(t, result.Round4.ID)
}

// TestSelectRandomGoals_TwoGoals tests selection with only 2 goals available
func TestSelectRandomGoals_TwoGoals(t *testing.T) {
	availableGoals := []Goal{
		{ID: "goal1", Name: "Goal 1", Expansion: "base"},
		{ID: "goal2", Name: "Goal 2", Expansion: "base"},
	}

	result, err := SelectRandomGoals(availableGoals)
	require.NoError(t, err)

	// Should return the 2 goals in first 2 rounds
	assert.NotEmpty(t, result.Round1.ID)
	assert.NotEmpty(t, result.Round2.ID)

	// Rounds 3 and 4 should be empty
	assert.Empty(t, result.Round3.ID)
	assert.Empty(t, result.Round4.ID)
}

// TestSelectRandomGoals_OneGoal tests selection with only 1 goal available
func TestSelectRandomGoals_OneGoal(t *testing.T) {
	availableGoals := []Goal{
		{ID: "goal1", Name: "Goal 1", Expansion: "base"},
	}

	result, err := SelectRandomGoals(availableGoals)
	require.NoError(t, err)

	// Should return the 1 goal in round 1
	assert.Equal(t, "goal1", result.Round1.ID)

	// Rounds 2, 3, and 4 should be empty
	assert.Empty(t, result.Round2.ID)
	assert.Empty(t, result.Round3.ID)
	assert.Empty(t, result.Round4.ID)
}

// TestSelectRandomGoals_NoGoals tests selection with empty goal list
func TestSelectRandomGoals_NoGoals(t *testing.T) {
	availableGoals := []Goal{}

	result, err := SelectRandomGoals(availableGoals)
	require.NoError(t, err)

	// All rounds should be empty
	assert.Empty(t, result.Round1.ID)
	assert.Empty(t, result.Round2.ID)
	assert.Empty(t, result.Round3.ID)
	assert.Empty(t, result.Round4.ID)
}

// TestSelectRandomGoals_DoesNotModifyOriginal tests that original slice is not modified
func TestSelectRandomGoals_DoesNotModifyOriginal(t *testing.T) {
	availableGoals := []Goal{
		{ID: "goal1", Name: "Goal 1", Expansion: "base"},
		{ID: "goal2", Name: "Goal 2", Expansion: "base"},
		{ID: "goal3", Name: "Goal 3", Expansion: "base"},
		{ID: "goal4", Name: "Goal 4", Expansion: "base"},
		{ID: "goal5", Name: "Goal 5", Expansion: "base"},
	}

	// Make a copy to compare later
	originalCopy := make([]Goal, len(availableGoals))
	copy(originalCopy, availableGoals)

	// Call SelectRandomGoals
	_, err := SelectRandomGoals(availableGoals)
	require.NoError(t, err)

	// Verify original slice is unchanged
	for i := range availableGoals {
		assert.Equal(t, originalCopy[i].ID, availableGoals[i].ID)
		assert.Equal(t, originalCopy[i].Name, availableGoals[i].Name)
	}
}

// TestSelectRandomGoals_Randomness tests that selection is actually random
func TestSelectRandomGoals_Randomness(t *testing.T) {
	availableGoals := GetAllGoals(true, false, false) // 16 base goals

	// Run selection 10 times
	var firstRoundIDs []string
	for i := 0; i < 10; i++ {
		result, err := SelectRandomGoals(availableGoals)
		require.NoError(t, err)
		firstRoundIDs = append(firstRoundIDs, result.Round1.ID)
	}

	// Check that we got at least 2 different results (very high probability with 16 goals)
	// This is a statistical test - it could theoretically fail by chance, but extremely unlikely
	uniqueIDs := make(map[string]bool)
	for _, id := range firstRoundIDs {
		uniqueIDs[id] = true
	}

	assert.GreaterOrEqual(t, len(uniqueIDs), 2,
		"With 10 random selections from 16 goals, should get at least 2 different Round1 goals")
}

// TestSelectRandomGoals_AllGoalsSelectable tests that all goals can be selected
func TestSelectRandomGoals_AllGoalsSelectable(t *testing.T) {
	availableGoals := []Goal{
		{ID: "goal1", Name: "Goal 1", Expansion: "base"},
		{ID: "goal2", Name: "Goal 2", Expansion: "base"},
		{ID: "goal3", Name: "Goal 3", Expansion: "base"},
		{ID: "goal4", Name: "Goal 4", Expansion: "base"},
		{ID: "goal5", Name: "Goal 5", Expansion: "base"},
	}

	// Run many selections and track which goals appear
	selectedGoals := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		result, err := SelectRandomGoals(availableGoals)
		require.NoError(t, err)

		selectedGoals[result.Round1.ID] = true
		selectedGoals[result.Round2.ID] = true
		selectedGoals[result.Round3.ID] = true
		selectedGoals[result.Round4.ID] = true
	}

	// All 5 goals should have been selected at least once across 100 iterations
	assert.Len(t, selectedGoals, 5, "All available goals should be selectable")
	assert.True(t, selectedGoals["goal1"])
	assert.True(t, selectedGoals["goal2"])
	assert.True(t, selectedGoals["goal3"])
	assert.True(t, selectedGoals["goal4"])
	assert.True(t, selectedGoals["goal5"])
}

// TestSelectRandomGoals_DifferentExpansions tests selection from mixed expansions
func TestSelectRandomGoals_DifferentExpansions(t *testing.T) {
	availableGoals := GetAllGoals(true, true, false) // Base + European = 26 goals

	result, err := SelectRandomGoals(availableGoals)
	require.NoError(t, err)

	// Verify we got 4 unique goals
	ids := map[string]bool{
		result.Round1.ID: true,
		result.Round2.ID: true,
		result.Round3.ID: true,
		result.Round4.ID: true,
	}
	assert.Len(t, ids, 4)

	// All selected goals should be from either base or european
	for _, goal := range []Goal{result.Round1, result.Round2, result.Round3, result.Round4} {
		assert.Contains(t, []string{"base", "european"}, goal.Expansion)
	}
}

// TestSelectRandomGoals_AllExpansions tests selection from all expansions
func TestSelectRandomGoals_AllExpansions(t *testing.T) {
	availableGoals := GetAllGoals(true, true, true) // All 36 goals

	result, err := SelectRandomGoals(availableGoals)
	require.NoError(t, err)

	// Verify we got 4 unique goals
	ids := map[string]bool{
		result.Round1.ID: true,
		result.Round2.ID: true,
		result.Round3.ID: true,
		result.Round4.ID: true,
	}
	assert.Len(t, ids, 4)

	// All selected goals should be from a valid expansion
	for _, goal := range []Goal{result.Round1, result.Round2, result.Round3, result.Round4} {
		assert.Contains(t, []string{"base", "european", "oceania"}, goal.Expansion)
	}
}

// TestSelectRandomGoals_StructureComplete tests that returned structure has all fields
func TestSelectRandomGoals_StructureComplete(t *testing.T) {
	availableGoals := GetAllGoals(true, false, false)

	result, err := SelectRandomGoals(availableGoals)
	require.NoError(t, err)

	// Verify all rounds have complete goal data
	for i, goal := range []Goal{result.Round1, result.Round2, result.Round3, result.Round4} {
		assert.NotEmpty(t, goal.ID, "Round %d should have ID", i+1)
		assert.NotEmpty(t, goal.Name, "Round %d should have Name", i+1)
		assert.NotEmpty(t, goal.Description, "Round %d should have Description", i+1)
		assert.NotEmpty(t, goal.Expansion, "Round %d should have Expansion", i+1)
	}
}

// TestRoundGoals_FieldTypes tests that RoundGoals structure has correct types
func TestRoundGoals_FieldTypes(t *testing.T) {
	var roundGoals RoundGoals

	// Initialize with test goals
	roundGoals.Round1 = Goal{ID: "test1", Name: "Test 1", Expansion: "base"}
	roundGoals.Round2 = Goal{ID: "test2", Name: "Test 2", Expansion: "base"}
	roundGoals.Round3 = Goal{ID: "test3", Name: "Test 3", Expansion: "base"}
	roundGoals.Round4 = Goal{ID: "test4", Name: "Test 4", Expansion: "base"}

	// Verify all fields are accessible
	assert.Equal(t, "test1", roundGoals.Round1.ID)
	assert.Equal(t, "test2", roundGoals.Round2.ID)
	assert.Equal(t, "test3", roundGoals.Round3.ID)
	assert.Equal(t, "test4", roundGoals.Round4.ID)
}

// TestSelectRandomGoals_ConsistentLength tests that result always has 4 rounds
func TestSelectRandomGoals_ConsistentLength(t *testing.T) {
	testCases := []struct {
		name     string
		numGoals int
	}{
		{"0 goals", 0},
		{"1 goal", 1},
		{"3 goals", 3},
		{"4 goals", 4},
		{"10 goals", 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			goals := make([]Goal, tc.numGoals)
			for i := 0; i < tc.numGoals; i++ {
				goals[i] = Goal{
					ID:        string(rune('A' + i)),
					Name:      "Goal",
					Expansion: "base",
				}
			}

			result, err := SelectRandomGoals(goals)
			require.NoError(t, err)

			// Structure should always have all 4 fields (even if some are empty)
			_ = result.Round1
			_ = result.Round2
			_ = result.Round3
			_ = result.Round4
		})
	}
}
