# Wingspan Scoring

A comprehensive scoring platform for the Wingspan board game. Track round goals, calculate end-game scores, save game results, and view player statistics. Built with Go and designed for containerized deployment on Kubernetes.

## Overview

Wingspan Scoring provides a complete digital scoring solution for Wingspan gameplay sessions:

- Random and manual round goal selection with expansion support
- Interactive score tracking with visual player cubes
- Complete end-game calculator with Oceania nectar scoring
- Persistent game history with SQLite database
- Player statistics tracking (wins, averages, win rates)
- Mobile-optimized responsive design

This is an unofficial fan-made tool. Wingspan is designed by Elizabeth Hargrave and published by Stonemaier Games.

## Features

### Round Goal Management

**Goal Selection:**
- Random selection of 4 unique goals from a pool of 36
- Manual goal selection with clickable interface
- Expansion filtering: Base Game (16 goals), European Expansion (10 goals), Oceania Expansion (10 goals)
- Side selection persistence across page navigation

**Scoring Modes:**
- **Blue Side** (Linear): 1 point per item, maximum 5 points per round
- **Green Side** (Competitive): Ranking-based scoring with 1st/2nd/3rd place points
  - Round 1: 4/1/0, Round 2: 5/2/0, Round 3: 6/3/2, Round 4: 7/4/2
- Automatic tie resolution (split and average points, rounded down)

### Interactive Score Tracking

- 2-5 player support with custom player names
- Colored player cubes: Blue, Purple, Green, Red, Yellow
- Click-to-place scoring system in score boxes
- Automatic score calculation with running totals across all 4 rounds
- Winner highlighting and tie handling
- Clear all cubes functionality
- Session persistence via browser localStorage

### Game End Calculator

Complete end-game scoring with all Wingspan categories:

- **Bird Points**: Total points from played bird cards
- **Bonus Cards**: Points from completed end-of-game goals
- **Round Goals**: Points earned across all 4 rounds
- **Eggs**: Points from eggs on bird cards
- **Cached Food**: Points from food tokens on bird cards
- **Tucked Cards**: Points from cards tucked under birds

**Oceania Expansion Support:**
- Nectar competitive scoring per habitat (Forest, Grassland, Wetland)
- 1st place: 5 points, 2nd place: 2 points
- Tie resolution with split/average points (rounded down)

**Tiebreaker:**
- Unused food tokens used for final ranking ties

**Features:**
- Mobile-optimized number input with numpad
- Automatic ranking calculation
- Player name and count preservation from Round Goals page

### Game History & Persistence

- Save game results to SQLite database
- View complete game history with pagination
- Delete individual game records
- Per-game details:
  - All player scores and rankings
  - Winner identification
  - Timestamp of game completion
  - Score breakdown by category

### Player Statistics (API Only)

The `/api/stats/{player}` endpoint provides performance tracking across all saved games:

- **Games Played**: Total number of games for each player
- **Total Wins**: Number of 1st place finishes
- **Average Score**: Mean score across all games
- **Win Rate**: Percentage of games won

Note: Currently accessible via API only. No web UI for viewing statistics yet (see issue #22).

## Technology Stack

**Backend:**
- Go 1.24
- SQLite database via `modernc.org/sqlite` (pure Go, no CGO required)
- Standard library HTTP server
- Embedded file system (`embed.FS`) for templates and static assets
- JSON API endpoints

**Frontend:**
- Vanilla JavaScript (no frameworks)
- Responsive CSS with mobile optimization
- Print-friendly styling
- Browser localStorage for session persistence

**Database:**
- SQLite 3 (pure Go implementation)
- Automatic database initialization
- Tables: `game_results` with indexes on `created_at` and `winner_name`
- JSON columns for complex data (players array, nectar scoring)
- Configurable database path via environment variable

**Container:**
- Multi-stage Podman/Docker build
- Build stage: `golang:1.24-alpine`
- Runtime stage: Red Hat UBI 10 Minimal
- Non-root user (UID 1000)
- Image size: ~30MB
- Published to GitHub Container Registry

**CI/CD:**
- GitHub Actions workflow
- Automated container builds on push to main
- Container image tagging: `latest` + git SHA
- Published to `ghcr.io`

## Quick Start

### Local Development

**Prerequisites:**
- Go 1.24 or higher

**Run:**
```bash
go run main.go
```

Access the application at http://localhost:8080

### Podman/Docker

**Build:**
```bash
podman build -t wingspan-scoring .
```

**Run:**
```bash
podman run -p 8080:8080 -v ./data:/app/data:Z wingspan-scoring
```

**Access:**
Open http://localhost:8080 in your browser

**Note:** The `-v` flag mounts a local directory for database persistence. The `:Z` flag is required for SELinux contexts.

### Pre-built Container Image

```bash
podman run -p 8080:8080 -v ./data:/app/data:Z ghcr.io/morey-tech/wingspan-scoring:latest
```

### Red Hat OpenShift Dev Spaces

The project includes a `devfile.yaml` for cloud-based development in Red Hat OpenShift Dev Spaces.

**Pre-installed Tools:**
- Go 1.25+
- Git
- GitHub CLI (gh) - automatically installed on workspace startup
- Podman/Buildah container tools

**Automatic Setup:**
- Go dependencies are automatically downloaded
- GitHub CLI is checked and installed if needed
- Setup runs via `.devcontainer/postCreate.sh` script

**Quick Start:**
1. Open the project in Dev Spaces (or VS Code with Dev Containers)
2. Wait for automatic setup to complete (~10 seconds first time, ~2 seconds after)
3. Use `gh` CLI for GitHub operations (PRs, issues, etc.)
4. Run predefined commands: `build`, `run`, `test`, `test-coverage`

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `DB_PATH` | SQLite database file path | `./data/wingspan.db` |

### Database

The application automatically creates the SQLite database on first run. The database file is stored at `./data/wingspan.db` by default.

**Schema:**
- `game_results` table with columns:
  - `id` (INTEGER PRIMARY KEY)
  - `created_at` (TIMESTAMP, indexed)
  - `winner_name` (TEXT, indexed)
  - `players` (JSON array)
  - `nectar_scoring` (JSON object)

## Kubernetes Deployment

The application is designed for deployment on Kubernetes/OpenShift. Reference deployment configuration:

**Repository:** [morey-tech/homelab](https://github.com/morey-tech/homelab/tree/184c1d1331a19a1f501a318397a655d7981e6922/kubernetes/ocp-mgmt/applications/wingspan-scoring)

**Key Resources:**
- **Deployment**: Single replica, resource limits (250Mi memory, 10m CPU request)
- **Service**: Exposes port 8080
- **Route/Ingress**: External access configuration
- **PersistentVolumeClaim**: 8GB storage for SQLite database (NVMe-backed LVM)

**Probes:**
- Startup, readiness, and liveness probes on `/` endpoint
- 10-second intervals, 1-second timeout, 3 failure threshold

**Container Image:**
```yaml
image: ghcr.io/morey-tech/wingspan-scoring:latest
imagePullPolicy: Always
```

**Volume Mount:**
```yaml
volumeMounts:
  - name: db-storage
    mountPath: /app/data
```

## Project Structure

```
wingspan-scoring/
├── main.go                     # HTTP server, routes, API handlers
├── goals/
│   ├── goals.go               # Goal definitions (36 total: Base, European, Oceania)
│   ├── selector.go            # Random selection algorithm (Fisher-Yates shuffle)
│   └── scorer.go              # Round goal scoring logic and tie resolution
├── db/
│   ├── db.go                  # Database initialization and connection
│   └── game_results.go        # CRUD operations for game results and stats
├── scoring/
│   └── scoring.go             # End-game scoring calculations (nectar, ranking)
├── templates/
│   ├── index.html             # Main page (Round Goals + Game End Calculator)
│   └── history.html           # Game history viewer
├── static/
│   ├── css/
│   │   ├── styles.css         # Round goals card styling
│   │   ├── game-end.css       # Game end calculator styling
│   │   └── history.css        # Game history page styling
│   └── js/
│       ├── app.js             # Main page interactions
│       └── history.js         # Game history page logic
├── data/                       # SQLite database directory (gitignored)
├── .github/
│   └── workflows/
│       └── container-build.yml # CI/CD pipeline
├── .devcontainer/
│   ├── devcontainer.json      # VS Code dev container config
│   └── ContainerFile          # Development container image
├── Containerfile              # Production multi-stage build
├── go.mod                     # Go module (Go 1.24)
├── CLAUDE.md                  # Development guidelines
└── README.md
```

## API Reference

### Web Pages

| Endpoint | Description |
|----------|-------------|
| `GET /` | Main application page (round goals + game end calculator) |
| `GET /history` | Game history viewer with pagination |

### API Endpoints

| Method | Endpoint | Description | Request Body / Params |
|--------|----------|-------------|----------------------|
| `POST` | `/api/new-game` | Generate new random goal set | Form: `base`, `european`, `oceania` (booleans) |
| `GET` | `/api/goals` | List all available goals | Query: `base`, `european`, `oceania` (booleans) |
| `POST` | `/api/calculate-scores` | Calculate round goal rankings | JSON: `{mode, round, playerCounts}` |
| `POST` | `/api/calculate-game-end` | Calculate end-game scores | JSON: player scores and nectar data |
| `GET` | `/api/games` | Retrieve game history | Query: `limit`, `offset` (pagination) |
| `GET` | `/api/games/{id}` | Get specific game result | Path: game ID |
| `DELETE` | `/api/games/{id}` | Delete game result | Path: game ID |
| `GET` | `/api/stats/{player}` | Get player statistics | Path: player name |

### Example API Usage

**Generate New Game:**
```bash
curl -X POST http://localhost:8080/api/new-game \
  -d "base=true&european=true&oceania=false"
```

**Get Player Stats:**
```bash
curl http://localhost:8080/api/stats/Alice
```

**Get Game History:**
```bash
curl http://localhost:8080/api/games?limit=10&offset=0
```

## Development

### Build

```bash
go build -o wingspan-scoring
./wingspan-scoring
```

### Format Code

```bash
go fmt ./...
```

### Build Container

```bash
podman build -t wingspan-scoring .
```

### Development Workflow

This project follows a Git workflow with feature branches. See [CLAUDE.md](CLAUDE.md) for full guidelines:

- Always create feature branches before making changes
- Never commit directly to `main`
- Use conventional commit format: `type: description`
- Run `go fmt` and `go build` before committing

**Branch Naming:**
- `feature/` - New features
- `fix/` - Bug fixes
- `refactor/` - Code refactoring
- `docs/` - Documentation updates

### Dev Container

VS Code dev container configuration is available in `.devcontainer/`. The container includes:
- Go 1.24
- Claude Code extension
- GitHub PR extension

## Goal Database

### Base Game (16 goals)

- Birds in specific habitats (Forest, Grassland, Wetland)
- Birds with specific nest types + eggs
- Eggs in habitats or on nest types
- Sets of eggs across all habitats
- Total birds played

### European Expansion (10 goals)

- Birds with tucked cards
- Food cost symbols
- Birds concentrated in one row
- Filled columns
- Birds with specific power types (brown/when played/end of round)
- Birds by point value ranges
- Birds with no eggs
- Food tokens and cards in hand

### Oceania Expansion (10 goals)

- Beak direction (facing left/right)
- Food symbols in costs
- Nectar tokens collected
- Action cube placement patterns
- Birds by point value ranges
- Total tucked cards across all birds

## Contributing

Contributions are welcome! Please follow the development guidelines in [CLAUDE.md](CLAUDE.md).

**Ideas for Enhancement:**
- Asia expansion goal support
- Game state export/import
- Multi-language support
- Dark mode theme
- Advanced statistics and charts
- PDF export of game results

## License

This is an unofficial fan-made tool. Wingspan is a trademark of Stonemaier Games. This project is not affiliated with or endorsed by Stonemaier Games.

## Credits

- **Game Design**: Elizabeth Hargrave
- **Publisher**: Stonemaier Games
- **Web Application**: Community-developed open source tool

## Support

For issues, bug reports, or feature requests, please open an issue on the GitHub repository.
