package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"wingspan-scoring/db"
	"wingspan-scoring/goals"
	importgames "wingspan-scoring/import"
	"wingspan-scoring/scoring"
)

//go:embed templates static
var content embed.FS

var tmpl *template.Template

// version is set at build time via -ldflags
var version = "dev"

func init() {
	var err error
	tmpl, err = template.ParseFS(content, "templates/*.html")
	if err != nil {
		log.Fatal("Error parsing templates:", err)
	}
}

type PageData struct {
	Goals        goals.RoundGoals
	HasGoals     bool
	BaseGame     bool
	European     bool
	Oceania      bool
	NumPlayers   int
	PageTitle    string
	PageSubtitle string
	CurrentPage  string
	Version      string
}

func main() {
	// Initialize database
	if err := db.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	log.Println("Database initialized successfully")

	// Serve static files
	fs := http.FileServer(http.FS(content))
	http.Handle("/static/", fs)

	// Routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/history", handleHistory)
	http.HandleFunc("/api/version", handleVersion)
	http.HandleFunc("/api/new-game", handleNewGame)
	http.HandleFunc("/api/goals", handleGetGoals)
	http.HandleFunc("/api/calculate-scores", handleCalculateScores)
	http.HandleFunc("/api/calculate-game-end", handleCalculateGameEnd)
	http.HandleFunc("/api/games", handleGetGames)
	http.HandleFunc("/api/games/", handleGameRoute)
	http.HandleFunc("/api/stats/", handleGetPlayerStats)
	http.HandleFunc("/api/leaderboard", handleGetLeaderboard)
	http.HandleFunc("/api/import", handleImportGames)

	log.Println("Starting Wingspan Scoring server on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggingMiddleware(http.DefaultServeMux)))
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
		Goals:        selectedGoals,
		HasGoals:     true,
		BaseGame:     true,
		European:     true,
		Oceania:      true,
		NumPlayers:   4, // Default to 4 players
		PageTitle:    "Round Goals",
		PageSubtitle: "Round End Goals",
		CurrentPage:  "home",
		Version:      version,
	}

	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"version": version,
	})
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

func handleCalculateGameEnd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Players        []scoring.PlayerGameEnd `json:"players"`
		IncludeOceania bool                    `json:"includeOceania"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate game end scores
	players, nectarScoring := scoring.CalculateGameEndScores(request.Players, request.IncludeOceania)

	// Save game result to database
	gameID, err := db.SaveGameResult(players, nectarScoring, request.IncludeOceania)
	if err != nil {
		log.Printf("Failed to save game result: %v", err)
		// Don't fail the request - just log the error
	} else {
		log.Printf("Saved game result with ID: %d", gameID)
	}

	response := struct {
		Players       []scoring.PlayerGameEnd `json:"players"`
		NectarScoring scoring.NectarScoring    `json:"nectarScoring"`
		GameID        int64                    `json:"gameId"`
	}{
		Players:       players,
		NectarScoring: nectarScoring,
		GameID:        gameID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleHistory(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		BaseGame:     true,
		European:     true,
		Oceania:      true,
		NumPlayers:   4,
		PageTitle:    "Game History",
		PageSubtitle: "Game History",
		CurrentPage:  "history",
		Version:      version,
	}

	err := tmpl.ExecuteTemplate(w, "history.html", data)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func handleGetGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse pagination parameters
	limit := parseIntDefault(r.URL.Query().Get("limit"), 50)
	offset := parseIntDefault(r.URL.Query().Get("offset"), 0)

	// Get games from database
	games, err := db.GetAllGameResults(limit, offset)
	if err != nil {
		log.Printf("Failed to get games: %v", err)
		http.Error(w, "Failed to retrieve games", http.StatusInternalServerError)
		return
	}

	// Get total count
	totalCount, err := db.CountGameResults()
	if err != nil {
		log.Printf("Failed to count games: %v", err)
		totalCount = 0
	}

	response := struct {
		Games      []db.GameResult `json:"games"`
		TotalCount int             `json:"totalCount"`
		Limit      int             `json:"limit"`
		Offset     int             `json:"offset"`
	}{
		Games:      games,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGameRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetGame(w, r)
	case http.MethodDelete:
		handleDeleteGame(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetGame(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	path := r.URL.Path
	idStr := path[len("/api/games/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	// Get game from database
	game, err := db.GetGameResult(id)
	if err != nil {
		log.Printf("Failed to get game %d: %v", id, err)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func handleDeleteGame(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	path := r.URL.Path
	idStr := path[len("/api/games/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	// Delete game from database
	err = db.DeleteGameResult(id)
	if err != nil {
		log.Printf("Failed to delete game %d: %v", id, err)
		http.Error(w, "Failed to delete game", http.StatusInternalServerError)
		return
	}

	log.Printf("Deleted game with ID: %d", id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Game deleted successfully",
	})
}

func handleGetPlayerStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract player name from path
	path := r.URL.Path
	playerName := path[len("/api/stats/"):]
	if playerName == "" {
		http.Error(w, "Player name required", http.StatusBadRequest)
		return
	}

	// Get player stats from database
	stats, err := db.GetPlayerStats(playerName)
	if err != nil {
		log.Printf("Failed to get stats for player %s: %v", playerName, err)
		http.Error(w, "Failed to retrieve player stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func handleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get leaderboard stats from database
	leaderboard, err := db.GetLeaderboardStats()
	if err != nil {
		log.Printf("Failed to get leaderboard stats: %v", err)
		http.Error(w, "Failed to retrieve leaderboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}

// Helper to parse int with default
func parseIntDefault(s string, defaultVal int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultVal
}

func handleImportGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get the file from form data
	file, header, err := r.FormFile("csvFile")
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Processing import file: %s (%d bytes)", header.Filename, header.Size)

	// Import games
	result, err := importgames.ImportGames(file)
	if err != nil {
		log.Printf("Import failed: %v", err)

		// Return error details
		response := struct {
			Success       bool                        `json:"success"`
			Message       string                      `json:"message"`
			GamesImported int                         `json:"gamesImported"`
			Errors        []importgames.ImportError   `json:"errors"`
		}{
			Success:       false,
			Message:       err.Error(),
			GamesImported: result.GamesImported,
			Errors:        result.Errors,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("Successfully imported %d games", result.GamesImported)

	// Return success response
	response := struct {
		Success       bool                        `json:"success"`
		Message       string                      `json:"message"`
		GamesImported int                         `json:"gamesImported"`
		Errors        []importgames.ImportError   `json:"errors"`
	}{
		Success:       true,
		Message:       "Import completed successfully",
		GamesImported: result.GamesImported,
		Errors:        result.Errors,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
