# Wingspan Round End Goals Web App

A web application that randomly selects round end goals for the Wingspan board game and displays them in a format similar to the physical round end goal card. Built with Go and deployable as a Docker container.

## Features

- **Random Goal Selection**: Randomly selects 4 unique round end goals for a game
- **Multiple Expansions**: Support for Base Game, European Expansion, and Oceania Expansion goals
- **Dual Scoring Modes**:
  - **Blue Side** (Beginner-friendly): Linear scoring, 1 point per item, maximum 5 points
  - **Green Side** (Competitive): Ranking-based scoring with 1st, 2nd, and 3rd place points
- **Visual Fidelity**: Card design that closely matches the physical game card
- **Responsive Design**: Works on desktop, tablet, and mobile devices
- **Print-Friendly**: Optimized CSS for printing the goal card

## Game Information

Wingspan is a competitive, card-driven, engine-building board game designed by Elizabeth Hargrave and published by Stonemaier Games. This is an unofficial fan-made tool to enhance gameplay.

### Scoring Modes

**Blue Side (Linear Scoring)**
- Each player scores 1 point per item counted
- Maximum of 5 points per round
- All players can score points

**Green Side (Competitive Scoring)**
- Players compete for 1st, 2nd, and 3rd place
- Points vary by round:
  - Round 1: 4 / 1 / 0
  - Round 2: 5 / 2 / 0
  - Round 3: 6 / 3 / 2
  - Round 4: 7 / 4 / 2
- Ties are resolved by averaging points and rounding down

## Quick Start

### Running Locally

**Prerequisites:**
- Go 1.22 or higher

**Steps:**
```bash
# Run the application
go run main.go

# Open browser to http://localhost:8080
```

### Running with Docker

**Build the Docker image:**
```bash
docker build -t wingspan-goals .
```

**Run the container:**
```bash
docker run -p 8080:8080 wingspan-goals
```

**Access the application:**
Open your browser to http://localhost:8080

## Usage

1. **Select Expansions**: Check the boxes for which expansions to include (Base Game, European, Oceania)
2. **Generate New Game**: Click "New Game" to randomly select 4 goals
3. **Toggle Scoring Mode**: Click "Switch to Green/Blue Side" to change between scoring modes
4. **View Goals**: The four round end goals are displayed in order from Round 1 to Round 4

## Project Structure

```
wingspan-goals/
├── main.go                 # HTTP server and handlers
├── goals/
│   ├── goals.go           # Goal definitions (Base, European, Oceania)
│   ├── selector.go        # Random selection algorithm
│   └── scorer.go          # Scoring calculation logic
├── templates/
│   └── index.html         # Main UI template
├── static/
│   ├── css/
│   │   └── styles.css     # Card styling (blue/green themes)
│   └── js/
│       └── app.js         # Frontend interactions
├── Dockerfile             # Multi-stage container build
├── .dockerignore
├── go.mod
└── README.md
```

## API Endpoints

- `GET /` - Main page with randomly selected goals
- `POST /api/new-game` - Generate new goal set (JSON response)
  - Form params: `base`, `european`, `oceania` (boolean)
- `GET /api/goals` - List all available goals
  - Query params: `base`, `european`, `oceania` (boolean)
- `POST /api/calculate-scores` - Calculate player rankings
  - JSON body: `{mode: "green"|"blue", round: 1-4, playerCounts: {...}}`

## Goal Database

### Base Game (16 goals)
- Birds in habitats (Forest, Grassland, Wetland)
- Birds with specific nest types + eggs
- Eggs in habitats
- Eggs on specific nest types
- Sets of eggs across all habitats
- Total birds played

### European Expansion (10 goals)
- Birds with tucked cards
- Food cost symbols
- Birds in one row
- Filled columns
- Birds with specific power types
- Birds by point value
- Birds with no eggs
- Food/cards in hand

### Oceania Expansion (10 goals)
- Beak direction (left/right)
- Food symbols in costs
- Nectar tokens
- Action cube placement
- Birds by point value
- Total tucked cards

## Docker Details

**Image Size:** ~20MB (Alpine-based)

**Environment Variables:**
- `PORT` - Server port (default: 8080)

**Health Check:**
- Endpoint: `http://localhost:8080/`
- Interval: 30 seconds

**Example Docker Compose:**
```yaml
version: '3.8'
services:
  wingspan-goals:
    build: .
    ports:
      - "8080:8080"
    restart: unless-stopped
```

## Development

**Build locally:**
```bash
go build -o wingspan-goals .
./wingspan-goals
```

**Run tests:**
```bash
go test ./...
```

## Credits

- **Game Design**: Elizabeth Hargrave
- **Publisher**: Stonemaier Games
- **Web App**: Unofficial fan-made tool

## License

This is an unofficial fan-made tool. Wingspan is a trademark of Stonemaier Games. This project is not affiliated with or endorsed by Stonemaier Games.

## Contributing

Contributions are welcome! Possible enhancements:
- Add Asia expansion goals
- Interactive scoring calculator
- Session persistence (save current game)
- Export/share game configuration
- Multiple language support
- Theme customization

## Support

For issues or questions, please open an issue on the project repository.
