// Final Score Calculator for Wingspan

const STORAGE_KEY = 'wingspan-final-scores';
const GOALS_PAGE_STORAGE_KEY = 'wingspanGameState';

// Player colors (matching goals page)
const PLAYER_COLORS = ['blue', 'purple', 'green', 'red', 'yellow'];

// State
let gameState = {
    numPlayers: 4,
    includeOceania: true,
    players: [], // Will store player names and colors from goals page
    playerNames: [],
    playerColors: [],
    roundGoalScores: [] // Total round goal scores from goals page
};

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    loadGameState();
    initializeEventListeners();
    generatePlayerRows();
});

// Event Listeners
function initializeEventListeners() {
    document.getElementById('numPlayers').addEventListener('change', handlePlayerCountChange);
    document.getElementById('oceaniaToggle').addEventListener('change', handleOceaniaToggle);
    document.getElementById('calculateBtn').addEventListener('click', calculateFinalScores);
    document.getElementById('clearBtn').addEventListener('click', clearAllScores);
}

function handlePlayerCountChange(e) {
    gameState.numPlayers = parseInt(e.target.value);
    generatePlayerRows();
    saveGameState();
}

function handleOceaniaToggle(e) {
    gameState.includeOceania = e.target.checked;
    const nectarHeaders = document.querySelectorAll('.nectar-header');
    const nectarCells = document.querySelectorAll('.nectar-cell');

    nectarHeaders.forEach(header => {
        header.style.display = gameState.includeOceania ? '' : 'none';
    });

    nectarCells.forEach(cell => {
        cell.style.display = gameState.includeOceania ? '' : 'none';
    });

    // Hide/show nectar rules
    const nectarRules = document.querySelector('.nectar-rules');
    if (nectarRules) {
        nectarRules.style.display = gameState.includeOceania ? '' : 'none';
    }

    // Hide/show nectar breakdown in results
    const nectarBreakdown = document.getElementById('nectarBreakdown');
    if (nectarBreakdown) {
        nectarBreakdown.style.display = gameState.includeOceania ? '' : 'none';
    }

    saveGameState();
}

// Generate player input rows
function generatePlayerRows() {
    const tbody = document.getElementById('scoreTableBody');
    tbody.innerHTML = '';

    for (let i = 1; i <= gameState.numPlayers; i++) {
        const row = document.createElement('tr');
        row.className = 'player-row';
        row.dataset.playerIndex = i;

        // Get player name and color from goals page data if available
        const playerName = (gameState.playerNames && gameState.playerNames[i-1]) || `Player ${i}`;
        const playerColor = (gameState.playerColors && gameState.playerColors[i-1]) || PLAYER_COLORS[(i-1) % PLAYER_COLORS.length];
        const roundGoalScore = (gameState.roundGoalScores && gameState.roundGoalScores[i-1]) || 0;

        row.innerHTML = `
            <td class="player-name-cell">
                <div class="player-name-wrapper">
                    <span class="player-color-indicator ${playerColor}"></span>
                    <input type="text"
                           class="player-name-input"
                           placeholder="Player ${i}"
                           value="${playerName}"
                           data-player="${i}"
                           data-field="name">
                </div>
            </td>
            <td class="score-cell">
                <input type="number"
                       class="score-input"
                       min="0"
                       value="0"
                       data-player="${i}"
                       data-field="birdPoints">
            </td>
            <td class="score-cell">
                <input type="number"
                       class="score-input"
                       min="0"
                       value="0"
                       data-player="${i}"
                       data-field="bonusCards">
            </td>
            <td class="score-cell">
                <input type="number"
                       class="score-input"
                       min="0"
                       value="${roundGoalScore}"
                       data-player="${i}"
                       data-field="roundGoals">
            </td>
            <td class="score-cell">
                <input type="number"
                       class="score-input"
                       min="0"
                       value="0"
                       data-player="${i}"
                       data-field="eggs">
            </td>
            <td class="score-cell">
                <input type="number"
                       class="score-input"
                       min="0"
                       value="0"
                       data-player="${i}"
                       data-field="cachedFood">
            </td>
            <td class="score-cell">
                <input type="number"
                       class="score-input"
                       min="0"
                       value="0"
                       data-player="${i}"
                       data-field="tuckedCards">
            </td>
            <td class="score-cell nectar-cell">
                <input type="number"
                       class="score-input nectar-input"
                       min="0"
                       value="0"
                       data-player="${i}"
                       data-field="nectarForest">
            </td>
            <td class="score-cell nectar-cell">
                <input type="number"
                       class="score-input nectar-input"
                       min="0"
                       value="0"
                       data-player="${i}"
                       data-field="nectarGrassland">
            </td>
            <td class="score-cell nectar-cell">
                <input type="number"
                       class="score-input nectar-input"
                       min="0"
                       value="0"
                       data-player="${i}"
                       data-field="nectarWetland">
            </td>
            <td class="score-cell tiebreaker-cell">
                <input type="number"
                       class="score-input tiebreaker-input"
                       min="0"
                       value="0"
                       data-player="${i}"
                       data-field="unusedFood">
            </td>
            <td class="total-cell">
                <span class="total-display">0</span>
            </td>
        `;

        tbody.appendChild(row);
    }

    // Add event listeners to all inputs
    document.querySelectorAll('.score-input, .player-name-input').forEach(input => {
        input.addEventListener('change', saveGameState);
    });

    // Apply saved state
    applySavedState();

    // Update nectar visibility
    handleOceaniaToggle({ target: document.getElementById('oceaniaToggle') });
}

// Calculate final scores
async function calculateFinalScores() {
    const players = [];

    // Gather player data
    for (let i = 1; i <= gameState.numPlayers; i++) {
        const nameInput = document.querySelector(`input[data-player="${i}"][data-field="name"]`);
        const playerName = nameInput.value || `Player ${i}`;

        const player = {
            playerName: playerName,
            birdPoints: parseInt(document.querySelector(`input[data-player="${i}"][data-field="birdPoints"]`).value) || 0,
            bonusCards: parseInt(document.querySelector(`input[data-player="${i}"][data-field="bonusCards"]`).value) || 0,
            roundGoals: parseInt(document.querySelector(`input[data-player="${i}"][data-field="roundGoals"]`).value) || 0,
            eggs: parseInt(document.querySelector(`input[data-player="${i}"][data-field="eggs"]`).value) || 0,
            cachedFood: parseInt(document.querySelector(`input[data-player="${i}"][data-field="cachedFood"]`).value) || 0,
            tuckedCards: parseInt(document.querySelector(`input[data-player="${i}"][data-field="tuckedCards"]`).value) || 0,
            nectarForest: parseInt(document.querySelector(`input[data-player="${i}"][data-field="nectarForest"]`).value) || 0,
            nectarGrassland: parseInt(document.querySelector(`input[data-player="${i}"][data-field="nectarGrassland"]`).value) || 0,
            nectarWetland: parseInt(document.querySelector(`input[data-player="${i}"][data-field="nectarWetland"]`).value) || 0,
            unusedFood: parseInt(document.querySelector(`input[data-player="${i}"][data-field="unusedFood"]`).value) || 0
        };

        players.push(player);
    }

    try {
        const response = await fetch('/api/calculate-final-score', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                players: players,
                includeOceania: gameState.includeOceania
            })
        });

        if (!response.ok) {
            throw new Error('Failed to calculate scores');
        }

        const result = await response.json();
        displayResults(result);
    } catch (error) {
        console.error('Error calculating scores:', error);
        alert('Error calculating final scores. Please try again.');
    }
}

// Display results
function displayResults(result) {
    const resultsSection = document.getElementById('resultsSection');
    resultsSection.style.display = 'block';

    // Clear all winner highlighting first
    document.querySelectorAll('.player-row').forEach(row => {
        row.classList.remove('winner-row');
    });

    // Update total cells in the table - match by player name, not index
    result.players.forEach(player => {
        // Find the row with this player's name
        const nameInputs = document.querySelectorAll('.player-name-input');
        let playerRowIndex = null;

        nameInputs.forEach((input, idx) => {
            if (input.value === player.playerName) {
                playerRowIndex = idx + 1;
            }
        });

        if (playerRowIndex) {
            const totalCell = document.querySelector(`tr[data-player-index="${playerRowIndex}"] .total-display`);
            if (totalCell) {
                totalCell.textContent = player.total;
            }

            // Highlight winner(s)
            const row = document.querySelector(`tr[data-player-index="${playerRowIndex}"]`);
            if (row && player.rank === 1) {
                row.classList.add('winner-row');
            }
        }
    });

    // Display nectar breakdown
    if (gameState.includeOceania) {
        displayNectarBreakdown(result.nectarScoring, result.players);
    }

    // Display final rankings
    displayRankings(result.players);

    // Scroll to results
    resultsSection.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
}

// Display nectar scoring breakdown
function displayNectarBreakdown(nectarScoring, players) {
    const habitats = ['forest', 'grassland', 'wetland'];
    const habitatNames = {
        forest: '🌲 Forest',
        grassland: '🌾 Grassland',
        wetland: '💧 Wetland'
    };

    habitats.forEach(habitat => {
        const resultDiv = document.getElementById(`nectar${habitat.charAt(0).toUpperCase() + habitat.slice(1)}Results`);
        if (!resultDiv) return;

        const scores = nectarScoring[habitat] || {};

        // Create sorted list of players by nectar count
        const playerNectar = players.map(p => ({
            name: p.playerName,
            count: p[`nectar${habitat.charAt(0).toUpperCase() + habitat.slice(1)}`],
            points: scores[p.playerName] || 0
        })).filter(p => p.count > 0);

        playerNectar.sort((a, b) => b.count - a.count);

        if (playerNectar.length === 0) {
            resultDiv.innerHTML = '<p class="no-nectar">No nectar in this habitat</p>';
        } else {
            resultDiv.innerHTML = playerNectar.map(p => `
                <div class="nectar-result-item ${p.points === 5 ? 'first-place' : p.points === 2 ? 'second-place' : ''}">
                    <span class="player-name">${p.name}</span>
                    <span class="nectar-count">${p.count} nectar</span>
                    <span class="nectar-points">${p.points} pts</span>
                </div>
            `).join('');
        }
    });
}

// Display final rankings
function displayRankings(players) {
    const rankingsList = document.getElementById('rankingsList');

    const rankingsHtml = players.map(player => {
        const rankLabel = player.rank === 1 ? '🏆 Winner' : `${getOrdinal(player.rank)} Place`;
        const tiebreakerInfo = players.filter(p => p.total === player.total).length > 1
            ? ` (${player.unusedFood} unused food)`
            : '';

        return `
            <div class="ranking-item ${player.rank === 1 ? 'winner' : ''}">
                <div class="rank-label">${rankLabel}</div>
                <div class="player-info">
                    <span class="player-name">${player.playerName}</span>
                    <span class="player-score">${player.total} points${tiebreakerInfo}</span>
                </div>
            </div>
        `;
    }).join('');

    rankingsList.innerHTML = rankingsHtml;
}

// Helper to get ordinal (1st, 2nd, 3rd, etc.)
function getOrdinal(n) {
    const s = ['th', 'st', 'nd', 'rd'];
    const v = n % 100;
    return n + (s[(v - 20) % 10] || s[v] || s[0]);
}

// Clear all scores
function clearAllScores() {
    if (!confirm('Clear all scores? This cannot be undone.')) {
        return;
    }

    document.querySelectorAll('.score-input').forEach(input => {
        input.value = '0';
    });

    document.querySelectorAll('.total-display').forEach(cell => {
        cell.textContent = '0';
    });

    document.querySelectorAll('.player-row').forEach(row => {
        row.classList.remove('winner-row');
    });

    const resultsSection = document.getElementById('resultsSection');
    resultsSection.style.display = 'none';

    saveGameState();
}

// Save game state to localStorage
function saveGameState() {
    const state = {
        numPlayers: gameState.numPlayers,
        includeOceania: gameState.includeOceania,
        players: []
    };

    // Save all player data
    for (let i = 1; i <= gameState.numPlayers; i++) {
        const playerData = {};
        document.querySelectorAll(`input[data-player="${i}"]`).forEach(input => {
            playerData[input.dataset.field] = input.value;
        });
        state.players.push(playerData);
    }

    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
}

// Load game state from localStorage
function loadGameState() {
    // First, try to load player data from the goals page
    loadPlayersFromGoalsPage();

    // Then load final score specific data (scores, etc.)
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved) {
        try {
            const state = JSON.parse(saved);

            // Only override player count if we didn't get it from goals page
            if (!gameState.playerNames || gameState.playerNames.length === 0) {
                gameState.numPlayers = state.numPlayers || 4;
            }

            gameState.includeOceania = state.includeOceania !== undefined ? state.includeOceania : true;
            gameState.players = state.players || [];

            // Update UI
            document.getElementById('numPlayers').value = gameState.numPlayers;
            document.getElementById('oceaniaToggle').checked = gameState.includeOceania;
        } catch (error) {
            console.error('Error loading saved state:', error);
        }
    }
}

// Load player data from the goals page localStorage
function loadPlayersFromGoalsPage() {
    try {
        const goalsPageData = localStorage.getItem(GOALS_PAGE_STORAGE_KEY);
        if (goalsPageData) {
            const data = JSON.parse(goalsPageData);
            if (data.players && data.players.length > 0) {
                gameState.numPlayers = data.players.length;
                gameState.playerNames = data.players.map(p => p.name || `Player ${p.id + 1}`);
                gameState.playerColors = data.players.map(p => p.color || PLAYER_COLORS[p.id]);

                // Calculate round goal scores for each player
                gameState.roundGoalScores = calculateRoundGoalScores(data.players, data.cubePlacements || {});

                // Update UI
                document.getElementById('numPlayers').value = gameState.numPlayers;

                console.log(`Loaded ${gameState.numPlayers} players from goals page:`, gameState.playerNames);
                console.log('Round goal scores:', gameState.roundGoalScores);
            }
        }
    } catch (error) {
        console.error('Error loading players from goals page:', error);
    }
}

// Calculate total round goal scores for each player from cube placements
function calculateRoundGoalScores(players, cubePlacements) {
    const scores = players.map(player => {
        let totalScore = 0;

        // For each round (1-4), find the score for this player's color
        for (let round = 1; round <= 4; round++) {
            // Check all cube placement keys for this round
            for (const [key, colors] of Object.entries(cubePlacements)) {
                const parts = key.split('-');
                const r = parseInt(parts[0]);
                const score = parseInt(parts[1]);

                // If this key is for the current round and contains this player's color
                if (r === round && colors && colors.includes(player.color)) {
                    totalScore += score;
                    break; // Found score for this round, move to next round
                }
            }
        }

        return totalScore;
    });

    return scores;
}

// Apply saved state to inputs
function applySavedState() {
    gameState.players.forEach((playerData, index) => {
        const playerNum = index + 1;
        if (playerNum > gameState.numPlayers) return;

        Object.keys(playerData).forEach(field => {
            // Skip roundGoals if we have data from the goals page
            if (field === 'roundGoals' && gameState.roundGoalScores && gameState.roundGoalScores.length > 0) {
                return; // Don't overwrite round goals from goals page
            }

            const input = document.querySelector(`input[data-player="${playerNum}"][data-field="${field}"]`);
            if (input) {
                input.value = playerData[field];
            }
        });
    });
}
