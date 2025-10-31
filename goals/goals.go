package goals

// Goal represents a single end-of-round goal
type Goal struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Expansion   string `json:"expansion"` // "base", "european", "oceania"
}

// All base game goals (8 tiles, 16 goals total)
var BaseGameGoals = []Goal{
	// Tile 1
	{
		ID:          "base-birds-forest",
		Name:        "Birds in Forest",
		Description: "Count the total number of birds you have played in your forest habitat",
		Expansion:   "base",
	},
	{
		ID:          "base-birds-grassland",
		Name:        "Birds in Grassland",
		Description: "Count the total number of birds you have played in your grassland habitat",
		Expansion:   "base",
	},
	// Tile 2
	{
		ID:          "base-birds-wetland",
		Name:        "Birds in Wetland",
		Description: "Count the total number of birds you have played in your wetland habitat",
		Expansion:   "base",
	},
	{
		ID:          "base-birds-bowl-egg",
		Name:        "Birds with Bowl Nests + Egg",
		Description: "Count birds with a bowl nest that have at least 1 egg (star nests count)",
		Expansion:   "base",
	},
	// Tile 3
	{
		ID:          "base-birds-cavity-egg",
		Name:        "Birds with Cavity Nests + Egg",
		Description: "Count birds with a cavity nest that have at least 1 egg (star nests count)",
		Expansion:   "base",
	},
	{
		ID:          "base-birds-ground-egg",
		Name:        "Birds with Ground Nests + Egg",
		Description: "Count birds with a ground nest that have at least 1 egg",
		Expansion:   "base",
	},
	// Tile 4
	{
		ID:          "base-birds-platform-egg",
		Name:        "Birds with Platform Nests + Egg",
		Description: "Count birds with a platform nest that have at least 1 egg (star nests count)",
		Expansion:   "base",
	},
	{
		ID:          "base-eggs-forest",
		Name:        "Eggs in Forest",
		Description: "Count the total number of eggs in your forest habitat (multiple eggs on one bird each count)",
		Expansion:   "base",
	},
	// Tile 5
	{
		ID:          "base-eggs-grassland",
		Name:        "Eggs in Grassland",
		Description: "Count the total number of eggs in your grassland habitat (multiple eggs on one bird each count)",
		Expansion:   "base",
	},
	{
		ID:          "base-eggs-wetland",
		Name:        "Eggs in Wetland",
		Description: "Count the total number of eggs in your wetland habitat (multiple eggs on one bird each count)",
		Expansion:   "base",
	},
	// Tile 6
	{
		ID:          "base-eggs-bowl",
		Name:        "Eggs on Bowl Nests",
		Description: "Count the total number of eggs on birds with a bowl nest (star nests count)",
		Expansion:   "base",
	},
	{
		ID:          "base-eggs-cavity",
		Name:        "Eggs on Cavity Nests",
		Description: "Count the total number of eggs on birds with a cavity nest (star nests count)",
		Expansion:   "base",
	},
	// Tile 7
	{
		ID:          "base-eggs-ground",
		Name:        "Eggs on Ground Nests",
		Description: "Count the total number of eggs on birds with a ground nest",
		Expansion:   "base",
	},
	{
		ID:          "base-eggs-platform",
		Name:        "Eggs on Platform Nests",
		Description: "Count the total number of eggs on birds with a platform nest (star nests count)",
		Expansion:   "base",
	},
	// Tile 8
	{
		ID:          "base-egg-sets",
		Name:        "Sets of Eggs in Each Habitat",
		Description: "Count sets of eggs (1 set = 1 egg in wetland + 1 egg in grassland + 1 egg in forest)",
		Expansion:   "base",
	},
	{
		ID:          "base-total-birds",
		Name:        "Total Birds Played",
		Description: "Count the total number of birds you have played",
		Expansion:   "base",
	},
}

// European Expansion goals (5 tiles, 10 goals)
var EuropeanGoals = []Goal{
	{
		ID:          "eu-birds-tucked",
		Name:        "Birds with Tucked Cards",
		Description: "Count the total number of birds that have at least 1 tucked card",
		Expansion:   "european",
	},
	{
		ID:          "eu-food-cost",
		Name:        "Food Cost of Played Birds",
		Description: "Count the total number of food symbols in the food cost of your bird cards",
		Expansion:   "european",
	},
	{
		ID:          "eu-birds-one-row",
		Name:        "Birds in One Row",
		Description: "Count birds in the single habitat row where you have the most birds",
		Expansion:   "european",
	},
	{
		ID:          "eu-filled-columns",
		Name:        "Filled Columns",
		Description: "Count the number of columns with all 5 spaces filled",
		Expansion:   "european",
	},
	{
		ID:          "eu-brown-powers",
		Name:        "Birds with Brown Powers",
		Description: "Count the total number of birds with brown (when activated) powers",
		Expansion:   "european",
	},
	{
		ID:          "eu-white-no-powers",
		Name:        "Birds with White/No Powers",
		Description: "Count the total number of birds with white (when played) or no powers",
		Expansion:   "european",
	},
	{
		ID:          "eu-birds-high-value",
		Name:        "Birds Worth > 4 Points",
		Description: "Count the total number of birds worth more than 4 victory points",
		Expansion:   "european",
	},
	{
		ID:          "eu-birds-no-eggs",
		Name:        "Birds with No Eggs",
		Description: "Count the total number of birds that have no eggs on them",
		Expansion:   "european",
	},
	{
		ID:          "eu-food-supply",
		Name:        "Food in Personal Supply",
		Description: "Count the total number of food tokens in your personal supply",
		Expansion:   "european",
	},
	{
		ID:          "eu-cards-hand",
		Name:        "Bird Cards in Hand",
		Description: "Count the total number of bird cards in your hand",
		Expansion:   "european",
	},
}

// Oceania Expansion goals (5 tiles, 10 goals)
var OceaniaGoals = []Goal{
	{
		ID:          "oc-beak-left",
		Name:        "Beak Pointing Left",
		Description: "Count the total number of birds whose beak is pointing left",
		Expansion:   "oceania",
	},
	{
		ID:          "oc-beak-right",
		Name:        "Beak Pointing Right",
		Description: "Count the total number of birds whose beak is pointing right",
		Expansion:   "oceania",
	},
	{
		ID:          "oc-invertebrate-cost",
		Name:        "Invertebrate in Food Cost",
		Description: "Count the number of invertebrate symbols in the food cost of your bird cards",
		Expansion:   "oceania",
	},
	{
		ID:          "oc-fruit-seed-cost",
		Name:        "Fruit + Seed in Food Cost",
		Description: "Count the total number of fruit and seed symbols in the food cost of your bird cards",
		Expansion:   "oceania",
	},
	{
		ID:          "oc-no-goal",
		Name:        "No Goal",
		Description: "No goal is scored this round. Keep your action cube and gain 1 extra turn in all following rounds",
		Expansion:   "oceania",
	},
	{
		ID:          "oc-rat-fish-cost",
		Name:        "Rat + Fish in Food Cost",
		Description: "Count the total number of rat and fish symbols in the food cost of your bird cards",
		Expansion:   "oceania",
	},
	{
		ID:          "oc-cubes-play-bird",
		Name:        "Cubes on Play a Bird",
		Description: "Count the total number of action cubes on the \"Play a Bird\" action",
		Expansion:   "oceania",
	},
	{
		ID:          "oc-birds-low-value",
		Name:        "Birds Worth ≤ 3 Points",
		Description: "Count the total number of birds worth 3 or fewer victory points",
		Expansion:   "oceania",
	},
}

// GetAllGoals returns all goals from specified expansions
func GetAllGoals(includeBase, includeEuropean, includeOceania bool) []Goal {
	var allGoals []Goal

	if includeBase {
		allGoals = append(allGoals, BaseGameGoals...)
	}
	if includeEuropean {
		allGoals = append(allGoals, EuropeanGoals...)
	}
	if includeOceania {
		allGoals = append(allGoals, OceaniaGoals...)
	}

	return allGoals
}
