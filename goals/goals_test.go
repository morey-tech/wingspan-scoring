package goals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBaseGameGoals_Count tests that we have the correct number of base game goals
func TestBaseGameGoals_Count(t *testing.T) {
	assert.Len(t, BaseGameGoals, 16, "Base game should have 16 goals (8 tiles * 2 sides)")
}

// TestEuropeanGoals_Count tests that we have the correct number of European expansion goals
func TestEuropeanGoals_Count(t *testing.T) {
	assert.Len(t, EuropeanGoals, 10, "European expansion should have 10 goals (5 tiles * 2 sides)")
}

// TestOceaniaGoals_Count tests that we have the correct number of Oceania expansion goals
func TestOceaniaGoals_Count(t *testing.T) {
	assert.Len(t, OceaniaGoals, 8, "Oceania expansion should have 8 goals")
}

// TestBaseGameGoals_AllFieldsPopulated tests that all base game goals have required fields
func TestBaseGameGoals_AllFieldsPopulated(t *testing.T) {
	for _, goal := range BaseGameGoals {
		assert.NotEmpty(t, goal.ID, "Goal ID should not be empty")
		assert.NotEmpty(t, goal.Name, "Goal Name should not be empty")
		assert.NotEmpty(t, goal.Description, "Goal Description should not be empty")
		assert.Equal(t, "base", goal.Expansion, "Base game goals should have expansion='base'")
	}
}

// TestEuropeanGoals_AllFieldsPopulated tests that all European goals have required fields
func TestEuropeanGoals_AllFieldsPopulated(t *testing.T) {
	for _, goal := range EuropeanGoals {
		assert.NotEmpty(t, goal.ID, "Goal ID should not be empty")
		assert.NotEmpty(t, goal.Name, "Goal Name should not be empty")
		assert.NotEmpty(t, goal.Description, "Goal Description should not be empty")
		assert.Equal(t, "european", goal.Expansion, "European goals should have expansion='european'")
	}
}

// TestOceaniaGoals_AllFieldsPopulated tests that all Oceania goals have required fields
func TestOceaniaGoals_AllFieldsPopulated(t *testing.T) {
	for _, goal := range OceaniaGoals {
		assert.NotEmpty(t, goal.ID, "Goal ID should not be empty")
		assert.NotEmpty(t, goal.Name, "Goal Name should not be empty")
		assert.NotEmpty(t, goal.Description, "Goal Description should not be empty")
		assert.Equal(t, "oceania", goal.Expansion, "Oceania goals should have expansion='oceania'")
	}
}

// TestAllGoals_UniqueIDs tests that all goal IDs are unique across all expansions
func TestAllGoals_UniqueIDs(t *testing.T) {
	allGoals := GetAllGoals(true, true, true)

	idMap := make(map[string]bool)
	for _, goal := range allGoals {
		assert.False(t, idMap[goal.ID], "Duplicate goal ID found: %s", goal.ID)
		idMap[goal.ID] = true
	}
}

// TestGetAllGoals_OnlyBase tests getting only base game goals
func TestGetAllGoals_OnlyBase(t *testing.T) {
	goals := GetAllGoals(true, false, false)

	assert.Len(t, goals, 16)
	for _, goal := range goals {
		assert.Equal(t, "base", goal.Expansion)
	}
}

// TestGetAllGoals_OnlyEuropean tests getting only European expansion goals
func TestGetAllGoals_OnlyEuropean(t *testing.T) {
	goals := GetAllGoals(false, true, false)

	assert.Len(t, goals, 10)
	for _, goal := range goals {
		assert.Equal(t, "european", goal.Expansion)
	}
}

// TestGetAllGoals_OnlyOceania tests getting only Oceania expansion goals
func TestGetAllGoals_OnlyOceania(t *testing.T) {
	goals := GetAllGoals(false, false, true)

	assert.Len(t, goals, 8)
	for _, goal := range goals {
		assert.Equal(t, "oceania", goal.Expansion)
	}
}

// TestGetAllGoals_BaseAndEuropean tests getting base + European goals
func TestGetAllGoals_BaseAndEuropean(t *testing.T) {
	goals := GetAllGoals(true, true, false)

	assert.Len(t, goals, 26) // 16 + 10

	baseCount := 0
	europeanCount := 0
	for _, goal := range goals {
		if goal.Expansion == "base" {
			baseCount++
		} else if goal.Expansion == "european" {
			europeanCount++
		}
	}

	assert.Equal(t, 16, baseCount)
	assert.Equal(t, 10, europeanCount)
}

// TestGetAllGoals_BaseAndOceania tests getting base + Oceania goals
func TestGetAllGoals_BaseAndOceania(t *testing.T) {
	goals := GetAllGoals(true, false, true)

	assert.Len(t, goals, 24) // 16 + 8

	baseCount := 0
	oceaniaCount := 0
	for _, goal := range goals {
		if goal.Expansion == "base" {
			baseCount++
		} else if goal.Expansion == "oceania" {
			oceaniaCount++
		}
	}

	assert.Equal(t, 16, baseCount)
	assert.Equal(t, 8, oceaniaCount)
}

// TestGetAllGoals_EuropeanAndOceania tests getting European + Oceania goals
func TestGetAllGoals_EuropeanAndOceania(t *testing.T) {
	goals := GetAllGoals(false, true, true)

	assert.Len(t, goals, 18) // 10 + 8

	europeanCount := 0
	oceaniaCount := 0
	for _, goal := range goals {
		if goal.Expansion == "european" {
			europeanCount++
		} else if goal.Expansion == "oceania" {
			oceaniaCount++
		}
	}

	assert.Equal(t, 10, europeanCount)
	assert.Equal(t, 8, oceaniaCount)
}

// TestGetAllGoals_AllExpansions tests getting all goals from all expansions
func TestGetAllGoals_AllExpansions(t *testing.T) {
	goals := GetAllGoals(true, true, true)

	assert.Len(t, goals, 34) // 16 + 10 + 8

	baseCount := 0
	europeanCount := 0
	oceaniaCount := 0
	for _, goal := range goals {
		switch goal.Expansion {
		case "base":
			baseCount++
		case "european":
			europeanCount++
		case "oceania":
			oceaniaCount++
		}
	}

	assert.Equal(t, 16, baseCount)
	assert.Equal(t, 10, europeanCount)
	assert.Equal(t, 8, oceaniaCount)
}

// TestGetAllGoals_NoExpansions tests getting goals when no expansions are selected
func TestGetAllGoals_NoExpansions(t *testing.T) {
	goals := GetAllGoals(false, false, false)

	assert.Len(t, goals, 0)
}

// TestBaseGameGoals_IDFormat tests that base game IDs follow the expected format
func TestBaseGameGoals_IDFormat(t *testing.T) {
	for _, goal := range BaseGameGoals {
		assert.Contains(t, goal.ID, "base-", "Base game goal IDs should start with 'base-'")
	}
}

// TestEuropeanGoals_IDFormat tests that European goal IDs follow the expected format
func TestEuropeanGoals_IDFormat(t *testing.T) {
	for _, goal := range EuropeanGoals {
		assert.Contains(t, goal.ID, "eu-", "European goal IDs should start with 'eu-'")
	}
}

// TestOceaniaGoals_IDFormat tests that Oceania goal IDs follow the expected format
func TestOceaniaGoals_IDFormat(t *testing.T) {
	for _, goal := range OceaniaGoals {
		assert.Contains(t, goal.ID, "oc-", "Oceania goal IDs should start with 'oc-'")
	}
}

// TestBaseGameGoals_SpecificGoals tests that specific expected base goals exist
func TestBaseGameGoals_SpecificGoals(t *testing.T) {
	expectedIDs := []string{
		"base-birds-forest",
		"base-birds-grassland",
		"base-birds-wetland",
		"base-total-birds",
		"base-egg-sets",
	}

	idMap := make(map[string]bool)
	for _, goal := range BaseGameGoals {
		idMap[goal.ID] = true
	}

	for _, expectedID := range expectedIDs {
		assert.True(t, idMap[expectedID], "Expected base goal %s not found", expectedID)
	}
}

// TestEuropeanGoals_SpecificGoals tests that specific expected European goals exist
func TestEuropeanGoals_SpecificGoals(t *testing.T) {
	expectedIDs := []string{
		"eu-birds-tucked",
		"eu-food-cost",
		"eu-filled-columns",
		"eu-cards-hand",
	}

	idMap := make(map[string]bool)
	for _, goal := range EuropeanGoals {
		idMap[goal.ID] = true
	}

	for _, expectedID := range expectedIDs {
		assert.True(t, idMap[expectedID], "Expected European goal %s not found", expectedID)
	}
}

// TestOceaniaGoals_SpecificGoals tests that specific expected Oceania goals exist
func TestOceaniaGoals_SpecificGoals(t *testing.T) {
	expectedIDs := []string{
		"oc-beak-left",
		"oc-beak-right",
		"oc-no-goal",
		"oc-birds-low-value",
	}

	idMap := make(map[string]bool)
	for _, goal := range OceaniaGoals {
		idMap[goal.ID] = true
	}

	for _, expectedID := range expectedIDs {
		assert.True(t, idMap[expectedID], "Expected Oceania goal %s not found", expectedID)
	}
}

// TestOceaniaGoals_NoGoalExists tests that the special "No Goal" exists in Oceania
func TestOceaniaGoals_NoGoalExists(t *testing.T) {
	var noGoal *Goal
	for i := range OceaniaGoals {
		if OceaniaGoals[i].ID == "oc-no-goal" {
			noGoal = &OceaniaGoals[i]
			break
		}
	}

	assert.NotNil(t, noGoal, "No Goal should exist in Oceania expansion")
	assert.Equal(t, "No Goal", noGoal.Name)
	assert.Contains(t, noGoal.Description, "No goal is scored")
}

// TestGoalDescriptions_NotEmpty tests that all goals have meaningful descriptions
func TestGoalDescriptions_NotEmpty(t *testing.T) {
	allGoals := GetAllGoals(true, true, true)

	for _, goal := range allGoals {
		assert.Greater(t, len(goal.Description), 10,
			"Goal %s should have a meaningful description", goal.ID)
	}
}

// TestGetAllGoals_OrderPreserved tests that goals are returned in the expected order
func TestGetAllGoals_OrderPreserved(t *testing.T) {
	goals := GetAllGoals(true, true, true)

	// First 16 should be base, next 10 European, last 8 Oceania
	for i := 0; i < 16; i++ {
		assert.Equal(t, "base", goals[i].Expansion,
			"First 16 goals should be base game")
	}

	for i := 16; i < 26; i++ {
		assert.Equal(t, "european", goals[i].Expansion,
			"Goals 16-25 should be European")
	}

	for i := 26; i < 34; i++ {
		assert.Equal(t, "oceania", goals[i].Expansion,
			"Goals 26-33 should be Oceania")
	}
}
