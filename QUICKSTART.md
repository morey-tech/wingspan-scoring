# Quick Start Guide

## Running the Application

### Option 1: Run Locally with Go

```bash
go run main.go
```

Open your browser to: http://localhost:8080

### Option 2: Build and Run Binary

```bash
go build -o wingspan-goals .
./wingspan-goals
```

Open your browser to: http://localhost:8080

### Option 3: Run with Docker

```bash
# Build the image
docker build -t wingspan-goals .

# Run the container
docker run -p 8080:8080 wingspan-goals
```

Open your browser to: http://localhost:8080

### Option 4: Run with Docker Compose

```bash
docker-compose up -d
```

Open your browser to: http://localhost:8080

To stop: `docker-compose down`

## Using the App

1. **Select Expansions**
   - Check/uncheck boxes for Base Game, European, and Oceania expansions
   - At least one expansion must be selected

2. **Generate New Game**
   - Click "New Game" button to randomly select 4 goals
   - Goals will be displayed for rounds 1-4

3. **Toggle Scoring Mode**
   - Click "Switch to Green/Blue Side" to change scoring systems
   - **Blue Side**: 1 point per item, max 5 points (beginner-friendly)
   - **Green Side**: Competitive ranking (1st/2nd/3rd place)

4. **View Goals**
   - Each round shows the goal name and description
   - Scoring tracks match the physical game card

## Scoring Reference

### Blue Side (Linear)
- Each player scores 1 point per item counted
- Maximum 5 points per round
- All players can score

### Green Side (Competitive)
- Players ranked by count (most to least)
- Points awarded to top 3 places:
  - **Round 1**: 4 / 1 / 0
  - **Round 2**: 5 / 2 / 0
  - **Round 3**: 6 / 3 / 2
  - **Round 4**: 7 / 4 / 2
- **Ties**: Points averaged and rounded down
  - Example: Two players tied for 1st in Round 3
  - Combined points: 6 + 3 = 9
  - Each gets: 9 ÷ 2 = 4.5 → 4 points

## Printing

The app is print-friendly! Just use your browser's print function (Ctrl+P / Cmd+P) to print the goal card.

## API Usage

### Get New Random Goals
```bash
curl -X POST "http://localhost:8080/api/new-game" \
  -d "base=true&european=true&oceania=false"
```

### Calculate Scores
```bash
curl -X POST "http://localhost:8080/api/calculate-scores" \
  -H "Content-Type: application/json" \
  -d '{
    "mode": "green",
    "round": 3,
    "playerCounts": {
      "Alice": 5,
      "Bob": 3,
      "Charlie": 5
    }
  }'
```

## Troubleshooting

**Port 8080 already in use?**
```bash
# Find process using port 8080
lsof -i :8080

# Kill it or use a different port
PORT=8081 ./wingspan-goals
```

**Docker build fails?**
- Make sure Docker is installed and running
- Check you have enough disk space
- Try: `docker system prune` to clean up

**Goals not updating?**
- Hard refresh your browser (Ctrl+Shift+R / Cmd+Shift+R)
- Clear browser cache
- Check browser console for errors (F12)

## Next Steps

See [README.md](README.md) for full documentation, API details, and contribution guidelines.
