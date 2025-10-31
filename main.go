package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"wingspan-goals/goals"
)

//go:embed templates static
var content embed.FS

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFS(content, "templates/*.html")
	if err != nil {
		log.Fatal("Error parsing templates:", err)
	}
}

type PageData struct {
	Goals     goals.RoundGoals
	HasGoals  bool
	BaseGame  bool
	European  bool
	Oceania   bool
	NumPlayers int
}

func main() {
	// Serve static files
	fs := http.FileServer(http.FS(content))
	http.Handle("/static/", fs)

	// Routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/new-game", handleNewGame)
	http.HandleFunc("/api/goals", handleGetGoals)
	http.HandleFunc("/api/calculate-scores", handleCalculateScores)

	log.Println("Starting Wingspan Goals server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	// Default: generate random goals from all expansions
	availableGoals := goals.GetAllGoals(true, true, true)
	selectedGoals, err := goals.SelectRandomGoals(availableGoals)
	if err != nil {
		http.Error(w, "Error selecting goals", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Goals:     selectedGoals,
		HasGoals:  true,
		BaseGame:  true,
		European:  true,
		Oceania:   true,
		NumPlayers: 4, // Default to 4 players
	}

	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func handleNewGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse expansion preferences
	includeBase := r.FormValue("base") == "true"
	includeEuropean := r.FormValue("european") == "true"
	includeOceania := r.FormValue("oceania") == "true"

	// Default to base game if nothing selected
	if !includeBase && !includeEuropean && !includeOceania {
		includeBase = true
	}

	availableGoals := goals.GetAllGoals(includeBase, includeEuropean, includeOceania)
	selectedGoals, err := goals.SelectRandomGoals(availableGoals)
	if err != nil {
		http.Error(w, "Error selecting goals", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(selectedGoals)
}

func handleGetGoals(w http.ResponseWriter, r *http.Request) {
	includeBase := r.URL.Query().Get("base") == "true"
	includeEuropean := r.URL.Query().Get("european") == "true"
	includeOceania := r.URL.Query().Get("oceania") == "true"

	// Default to all if no parameters
	if !includeBase && !includeEuropean && !includeOceania {
		includeBase = true
		includeEuropean = true
		includeOceania = true
	}

	allGoals := goals.GetAllGoals(includeBase, includeEuropean, includeOceania)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allGoals)
}

func handleCalculateScores(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Mode         string         `json:"mode"` // "green" or "blue"
		Round        int            `json:"round"`
		PlayerCounts map[string]int `json:"playerCounts"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var scores []goals.PlayerScore
	if request.Mode == "green" {
		scores = goals.CalculateGreenScores(request.PlayerCounts, request.Round)
	} else {
		scores = goals.CalculateBlueScores(request.PlayerCounts)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}

// Helper to parse int with default
func parseIntDefault(s string, defaultVal int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultVal
}
