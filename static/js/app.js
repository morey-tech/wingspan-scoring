// Track current mode
let currentMode = 'blue'; // 'blue' or 'green'

// Player colors
const PLAYER_COLORS = ['blue', 'purple', 'green', 'red', 'yellow'];
const PLAYER_COLOR_NAMES = ['Blue', 'Purple', 'Green', 'Red', 'Yellow'];
// Default player names
const DEFAULT_PLAYER_NAMES = ['Dani', 'Nick', 'Player 3', 'Player 4', 'Player 5'];
// Default color assignments: Player 1 = blue, Player 2 = yellow, then others in order
const DEFAULT_PLAYER_COLORS = ['blue', 'yellow', 'green', 'purple', 'red'];

// Game state
let gameState = {
    players: [],
    cubePlacements: {}, // { "round-score": ["playerColor1", "playerColor2"] }
    goals: null // Current round goals (will be captured from HTML or loaded from server)
};

// All available goals (fetched from API)
let allGoals = [];

// Fetch all available goals from API
async function fetchAllGoals() {
    try {
        const base = document.getElementById('base').checked;
        const european = document.getElementById('european').checked;
        const oceania = document.getElementById('oceania').checked;

        const response = await fetch(`/api/goals?base=${base}&european=${european}&oceania=${oceania}`);
        if (!response.ok) {
            throw new Error('Failed to fetch goals');
        }

        allGoals = await response.json();
        return allGoals;
    } catch (error) {
        console.error('Error fetching goals:', error);
        return [];
    }
}

// Show goal selection menu
function showGoalMenu(element, round) {
    // Remove any existing goal menu
    const existingMenu = document.querySelector('.goal-menu');
    if (existingMenu) {
        existingMenu.remove();
    }

    const menu = document.createElement('div');
    menu.className = 'goal-menu';

    // Get currently selected goals to filter them out
    const rounds = ['round1', 'round2', 'round3', 'round4'];
    const selectedGoalIds = [];
    if (gameState.goals) {
        rounds.forEach(r => {
            if (gameState.goals[r] && gameState.goals[r].id) {
                selectedGoalIds.push(gameState.goals[r].id);
            }
        });
    }

    // Get current goal for this round
    const roundKey = `round${round}`;
    const currentGoalId = gameState.goals && gameState.goals[roundKey] ? gameState.goals[roundKey].id : null;

    // Get expansion label
    const getExpansionLabel = (expansion) => {
        if (expansion === 'base') return 'Base';
        if (expansion === 'european') return 'European';
        if (expansion === 'oceania') return 'Oceania';
        return expansion;
    };

    // Build menu HTML
    const goalItems = allGoals
        .map(goal => {
            const isSelected = goal.id === currentGoalId;
            const isDisabled = selectedGoalIds.includes(goal.id) && !isSelected;
            const classes = ['goal-menu-item'];
            if (isSelected) classes.push('selected');
            if (isDisabled) classes.push('disabled');

            return `
                <div class="${classes.join(' ')}" data-goal-id="${goal.id}">
                    <div class="goal-content">
                        <strong class="menu-goal-name">${goal.name}</strong>
                        <span class="menu-goal-description">${goal.description}</span>
                        <span class="expansion-badge ${goal.expansion}">${getExpansionLabel(goal.expansion)}</span>
                    </div>
                    ${isSelected ? '<span class="checkmark">✓</span>' : ''}
                </div>
            `;
        }).join('');

    menu.innerHTML = `
        <div class="goal-menu-header">Select Round ${round} Goal</div>
        <div class="goal-menu-items">
            ${goalItems}
        </div>
        <div class="goal-menu-actions">
            <button class="menu-btn menu-btn-close">Close</button>
        </div>
    `;

    // Position menu near the clicked element
    const rect = element.getBoundingClientRect();
    menu.style.position = 'fixed';
    menu.style.left = `${rect.left}px`;
    menu.style.top = `${rect.bottom + 5}px`;

    document.body.appendChild(menu);

    // Add event listeners to goal items (excluding disabled ones)
    menu.querySelectorAll('.goal-menu-item:not(.disabled)').forEach(item => {
        item.addEventListener('click', () => {
            const goalId = item.dataset.goalId;
            handleGoalSelection(round, goalId);
            menu.remove();
        });
    });

    // Close button handler
    menu.querySelector('.menu-btn-close').addEventListener('click', () => {
        menu.remove();
    });

    // Click outside to close
    setTimeout(() => {
        document.addEventListener('click', function closeMenu(e) {
            if (!menu.contains(e.target) && e.target !== element) {
                menu.remove();
                document.removeEventListener('click', closeMenu);
            }
        });
    }, 100);
}

// Handle goal selection
function handleGoalSelection(round, goalId) {
    // Find the selected goal
    const selectedGoal = allGoals.find(g => g.id === goalId);
    if (!selectedGoal) return;

    // Update gameState
    const roundKey = `round${round}`;
    if (!gameState.goals) {
        gameState.goals = {};
    }
    gameState.goals[roundKey] = selectedGoal;

    // Update the display
    const goalInfo = document.querySelector(`.goal-info[data-round="${round}"]`);
    if (goalInfo) {
        const goalName = goalInfo.querySelector('.goal-name');
        const goalDescription = goalInfo.querySelector('.goal-description');

        if (goalName) goalName.textContent = selectedGoal.name;
        if (goalDescription) goalDescription.textContent = selectedGoal.description;
    }

    // Save state
    saveGameState();
}

// Initialize
document.addEventListener('DOMContentLoaded', function() {
    const newGameBtn = document.getElementById('newGame');
    const toggleModeBtn = document.getElementById('toggleMode');
    const numPlayersSelect = document.getElementById('numPlayers');
    const clearScoresBtn = document.getElementById('clearScores');

    newGameBtn.addEventListener('click', generateNewGame);
    toggleModeBtn.addEventListener('click', toggleScoringMode);
    numPlayersSelect.addEventListener('change', handlePlayerCountChange);
    clearScoresBtn.addEventListener('click', clearAllCubes);

    // Set initial button color (we start in blue mode, button says "Switch to Green", so make it green)
    toggleModeBtn.classList.add('green-mode');

    // Initialize players
    initializePlayers(parseInt(numPlayersSelect.value));

    // Add click handlers to all score boxes
    initializeScoreBoxHandlers();

    // Add click event listeners to goal-info areas
    document.querySelectorAll('.goal-info').forEach(goalInfo => {
        goalInfo.addEventListener('click', function(e) {
            const round = parseInt(this.dataset.round);
            showGoalMenu(this, round);
        });
    });

    // Add event listeners for expansion checkboxes to re-fetch goals
    const expansionCheckboxes = ['base', 'european', 'oceania'];
    expansionCheckboxes.forEach(id => {
        const checkbox = document.getElementById(id);
        if (checkbox) {
            checkbox.addEventListener('change', async () => {
                await fetchAllGoals();
            });
        }
    });

    // Initialize goals
    (async function() {
        // Fetch all available goals
        await fetchAllGoals();

        // Load saved state if available
        const hasSavedState = loadGameState();

        // If no saved state with goals, capture goals from server-rendered HTML
        if (!hasSavedState || !gameState.goals) {
            // Capture the goals that were rendered by the server
            gameState.goals = captureGoalsFromDisplay();
            saveGameState();
        }
    })();
});

// Initialize players based on count
function initializePlayers(count) {
    gameState.players = [];
    for (let i = 0; i < count; i++) {
        gameState.players.push({
            id: i,
            name: DEFAULT_PLAYER_NAMES[i] || `Player ${i + 1}`,
            color: DEFAULT_PLAYER_COLORS[i],
            scores: [null, null, null, null] // Round 1-4 scores
        });
    }
    renderPlayerList();
    renderScoreTable();
}

// Render player list in the setup area
function renderPlayerList() {
    const playerList = document.getElementById('playerList');
    playerList.innerHTML = '';

    gameState.players.forEach((player, index) => {
        const playerDiv = document.createElement('div');
        playerDiv.className = 'player-item';

        playerDiv.innerHTML = `
            <div class="player-cube ${player.color} clickable-cube" data-player-id="${player.id}" title="Click to change color"></div>
            <input
                type="text"
                value="${player.name}"
                class="player-name-input"
                data-player-id="${player.id}"
                placeholder="Player ${index + 1}"
            />
        `;
        playerList.appendChild(playerDiv);

        // Add event listener for name changes
        const input = playerDiv.querySelector('.player-name-input');
        input.addEventListener('change', (e) => {
            gameState.players[player.id].name = e.target.value || `Player ${index + 1}`;
            renderScoreTable();
            saveGameState();
        });

        // Add event listener for cube click (color picker)
        const cube = playerDiv.querySelector('.player-cube');
        cube.addEventListener('click', (e) => {
            showColorPicker(e.currentTarget, player.id);
        });
    });
}

// Show color picker menu
function showColorPicker(cubeElement, playerId) {
    // Remove any existing color picker
    const existingPicker = document.querySelector('.color-picker-menu');
    if (existingPicker) {
        existingPicker.remove();
    }

    const player = gameState.players.find(p => p.id === playerId);
    if (!player) return;

    // Get list of already-used colors (excluding this player's current color)
    const usedColors = gameState.players
        .filter(p => p.id !== playerId)
        .map(p => p.color);

    const menu = document.createElement('div');
    menu.className = 'color-picker-menu';

    const colorOptions = PLAYER_COLORS.map((color, index) => {
        const isSelected = player.color === color;
        const isUsed = usedColors.includes(color);
        const colorName = PLAYER_COLOR_NAMES[index];

        return `
            <div class="color-picker-item ${isSelected ? 'selected' : ''} ${isUsed ? 'disabled' : ''}"
                 data-color="${color}">
                <div class="player-cube ${color}"></div>
                <span>${colorName}</span>
                ${isSelected ? '<span class="checkmark">✓</span>' : ''}
            </div>
        `;
    }).join('');

    menu.innerHTML = `
        <div class="color-picker-header">Select Color</div>
        <div class="color-picker-items">
            ${colorOptions}
        </div>
    `;

    // Position menu near the cube
    const rect = cubeElement.getBoundingClientRect();
    menu.style.position = 'fixed';
    menu.style.left = `${rect.left}px`;
    menu.style.top = `${rect.bottom + 5}px`;

    document.body.appendChild(menu);

    // Add event listeners to color options
    menu.querySelectorAll('.color-picker-item:not(.disabled)').forEach(item => {
        item.addEventListener('click', () => {
            const newColor = item.dataset.color;
            handlePlayerColorChange(playerId, newColor);
            menu.remove();
        });
    });

    // Close menu when clicking outside
    setTimeout(() => {
        document.addEventListener('click', function closePicker(e) {
            if (!menu.contains(e.target) && e.target !== cubeElement) {
                menu.remove();
                document.removeEventListener('click', closePicker);
            }
        });
    }, 100);
}

// Handle player color change
function handlePlayerColorChange(playerId, newColor) {
    const player = gameState.players.find(p => p.id === playerId);
    if (!player) return;

    const oldColor = player.color;

    // Update player color
    player.color = newColor;

    // Update all cube placements with this player's cubes to use new color
    for (const key in gameState.cubePlacements) {
        const colors = gameState.cubePlacements[key];
        const index = colors.indexOf(oldColor);
        if (index > -1) {
            // Check if this placement belongs to this player
            // (by checking if only this player has this color in this position)
            colors[index] = newColor;
        }
    }

    // Re-render everything
    renderPlayerList();
    renderScoreTable();
    renderAllCubes();
    saveGameState();
}

// Render score table
function renderScoreTable() {
    const tbody = document.getElementById('scoreTableBody');
    tbody.innerHTML = '';

    gameState.players.forEach(player => {
        const row = document.createElement('tr');
        const scores = calculatePlayerScores(player);
        const total = scores.reduce((sum, score) => sum + (score || 0), 0);

        // Check if this player is winning
        const isWinner = isPlayerWinning(player.id, total);
        if (isWinner && total > 0) {
            row.classList.add('winning-player');
        }

        row.innerHTML = `
            <td>
                <div class="player-cell">
                    <div class="player-cube ${player.color}"></div>
                    <span>${player.name}</span>
                </div>
            </td>
            <td>${scores[0] !== null ? scores[0] : '-'}</td>
            <td>${scores[1] !== null ? scores[1] : '-'}</td>
            <td>${scores[2] !== null ? scores[2] : '-'}</td>
            <td>${scores[3] !== null ? scores[3] : '-'}</td>
            <td class="total-score"><strong>${total}</strong></td>
        `;
        tbody.appendChild(row);
    });
}

// Calculate player scores from cube placements
function calculatePlayerScores(player) {
    const scores = [null, null, null, null];

    for (let round = 1; round <= 4; round++) {
        // Find which score box this player's cube is in for this round
        for (const [key, colors] of Object.entries(gameState.cubePlacements)) {
            const parts = key.split('-');
            const r = parseInt(parts[0]);
            const score = parseInt(parts[1]);
            // parts[2] is the position, which we don't need for scoring
            if (r === round && colors.includes(player.color)) {
                scores[round - 1] = score;
                break;
            }
        }
    }

    return scores;
}

// Check if player is winning
function isPlayerWinning(playerId, playerTotal) {
    const totals = gameState.players.map((p, idx) => {
        const scores = calculatePlayerScores(p);
        return scores.reduce((sum, score) => sum + (score || 0), 0);
    });
    const maxTotal = Math.max(...totals);
    return playerTotal === maxTotal && maxTotal > 0;
}

// Handle player count change
function handlePlayerCountChange(e) {
    const newCount = parseInt(e.target.value);
    if (confirm('Changing player count will clear all placed cubes. Continue?')) {
        initializePlayers(newCount);
        clearAllCubes(true); // Skip confirmation - already confirmed above
    } else {
        e.target.value = gameState.players.length;
    }
}

// Initialize click handlers for score boxes
function initializeScoreBoxHandlers() {
    const scoreBoxes = document.querySelectorAll('.score-box');
    scoreBoxes.forEach(box => {
        box.addEventListener('click', handleScoreBoxClick);
    });
}

// Handle score box click
function handleScoreBoxClick(e) {
    const box = e.currentTarget;
    const round = parseInt(box.dataset.round);
    const score = parseInt(box.dataset.score);
    const position = box.dataset.position;

    // Show player selection menu
    showPlayerMenu(box, round, score, position);
}

// Show player selection menu
function showPlayerMenu(box, round, score, position) {
    // Remove any existing menu
    const existingMenu = document.querySelector('.player-menu');
    if (existingMenu) {
        existingMenu.remove();
    }

    const menu = document.createElement('div');
    menu.className = 'player-menu';

    // Get current placements for this box
    const key = `${round}-${score}-${position}`;
    const currentPlacements = gameState.cubePlacements[key] || [];

    menu.innerHTML = `
        <div class="player-menu-header">Select Player</div>
        <div class="player-menu-items">
            ${gameState.players.map(player => {
                const isPlaced = currentPlacements.includes(player.color);
                return `
                    <div class="player-menu-item ${isPlaced ? 'selected' : ''}" data-player-color="${player.color}">
                        <div class="player-cube ${player.color}"></div>
                        <span>${player.name}</span>
                        ${isPlaced ? '<span class="checkmark">✓</span>' : ''}
                    </div>
                `;
            }).join('')}
        </div>
        <div class="player-menu-actions">
            <button class="menu-btn menu-btn-close">Close</button>
        </div>
    `;

    // Position menu near the box
    const rect = box.getBoundingClientRect();
    menu.style.position = 'fixed';
    menu.style.left = `${rect.left}px`;
    menu.style.top = `${rect.bottom + 5}px`;

    document.body.appendChild(menu);

    // Add event listeners
    menu.querySelectorAll('.player-menu-item').forEach(item => {
        item.addEventListener('click', () => {
            const playerColor = item.dataset.playerColor;
            toggleCubePlacement(round, score, position, playerColor);
            renderCubesInBox(box, round, score, position);
            renderScoreTable();
            saveGameState();

            // Update menu to show selection
            showPlayerMenu(box, round, score, position);
        });
    });

    menu.querySelector('.menu-btn-close').addEventListener('click', () => {
        menu.remove();
    });

    // Close menu when clicking outside
    setTimeout(() => {
        document.addEventListener('click', function closeMenu(e) {
            if (!menu.contains(e.target) && e.target !== box) {
                menu.remove();
                document.removeEventListener('click', closeMenu);
            }
        });
    }, 100);
}

// Toggle cube placement
function toggleCubePlacement(round, score, position, playerColor) {
    const key = `${round}-${score}-${position}`;

    // First, remove this player from any other box in this round
    for (const k in gameState.cubePlacements) {
        if (k.startsWith(`${round}-`)) {
            const index = gameState.cubePlacements[k].indexOf(playerColor);
            if (index > -1) {
                gameState.cubePlacements[k].splice(index, 1);
                if (gameState.cubePlacements[k].length === 0) {
                    delete gameState.cubePlacements[k];
                }
            }
        }
    }

    // Add to the new box
    if (!gameState.cubePlacements[key]) {
        gameState.cubePlacements[key] = [];
    }

    if (!gameState.cubePlacements[key].includes(playerColor)) {
        gameState.cubePlacements[key].push(playerColor);
    }

    // Re-render all boxes in this round
    const roundBoxes = document.querySelectorAll(`.score-box[data-round="${round}"]`);
    roundBoxes.forEach(box => {
        const boxScore = parseInt(box.dataset.score);
        const boxPosition = box.dataset.position;
        renderCubesInBox(box, round, boxScore, boxPosition);
    });
}

// Render cubes in a score box
function renderCubesInBox(box, round, score, position) {
    const key = `${round}-${score}-${position}`;
    const container = box.querySelector('.cube-container');
    if (!container) return;

    container.innerHTML = '';
    const placements = gameState.cubePlacements[key] || [];

    placements.forEach(color => {
        const cube = document.createElement('div');
        cube.className = `placed-cube ${color}`;
        container.appendChild(cube);
    });
}

// Clear all cubes
function clearAllCubes(skipConfirm = false) {
    if (!skipConfirm && !confirm('Clear all placed cubes?')) {
        return;
    }

    gameState.cubePlacements = {};

    // Clear all cube containers
    document.querySelectorAll('.cube-container').forEach(container => {
        container.innerHTML = '';
    });

    renderScoreTable();
    saveGameState();
}

// Generate a new game with selected expansions
async function generateNewGame() {
    const base = document.getElementById('base').checked;
    const european = document.getElementById('european').checked;
    const oceania = document.getElementById('oceania').checked;

    // Ensure at least one expansion is selected
    if (!base && !european && !oceania) {
        alert('Please select at least one expansion');
        return;
    }

    if (Object.keys(gameState.cubePlacements).length > 0) {
        if (!confirm('Generate a new game? This will clear all placed cubes.')) {
            return;
        }
    }

    try {
        const response = await fetch('/api/new-game', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `base=${base}&european=${european}&oceania=${oceania}`
        });

        if (!response.ok) {
            throw new Error('Failed to generate new game');
        }

        const goals = await response.json();
        updateGoalDisplay(goals);
        clearAllCubes();

        // Re-fetch goals based on new expansion selection
        await fetchAllGoals();
    } catch (error) {
        console.error('Error generating new game:', error);
        alert('Failed to generate new game. Please try again.');
    }
}

// Capture current goals from HTML display
function captureGoalsFromDisplay() {
    const rounds = ['round1', 'round2', 'round3', 'round4'];
    const goalCard = document.getElementById('goalCard');
    const roundRows = goalCard.querySelectorAll('.round-row');
    const goals = {};

    rounds.forEach((round, index) => {
        const row = roundRows[index];
        if (row) {
            const nameElement = row.querySelector('.goal-name');
            const descElement = row.querySelector('.goal-description');

            if (nameElement && descElement) {
                const goalName = nameElement.textContent;
                const goalDescription = descElement.textContent;

                // Try to match with full goal data from allGoals to get the ID
                const matchedGoal = allGoals.find(g => g.name === goalName);

                if (matchedGoal) {
                    // Use the full goal object with ID
                    goals[round] = matchedGoal;
                } else {
                    // Fallback to basic data if no match found
                    goals[round] = {
                        name: goalName,
                        description: goalDescription
                    };
                }
            }
        }
    });

    return goals;
}

// Update the goal display with new goals (internal function, doesn't save)
function setGoalDisplay(goals) {
    const rounds = ['round1', 'round2', 'round3', 'round4'];
    const goalCard = document.getElementById('goalCard');
    const roundRows = goalCard.querySelectorAll('.round-row');

    rounds.forEach((round, index) => {
        const goal = goals[round];
        const row = roundRows[index];

        if (goal && row) {
            const nameElement = row.querySelector('.goal-name');
            const descElement = row.querySelector('.goal-description');

            if (nameElement) nameElement.textContent = goal.name;
            if (descElement) descElement.textContent = goal.description;
        }
    });
}

// Update the goal display with new goals and save to state
function updateGoalDisplay(goals) {
    setGoalDisplay(goals);
    gameState.goals = goals;
    saveGameState();
}

// Toggle between blue and green scoring modes
function toggleScoringMode() {
    const goalCard = document.getElementById('goalCard');
    const toggleBtn = document.getElementById('toggleMode');
    const blueTracks = goalCard.querySelectorAll('.blue-track');
    const greenTracks = goalCard.querySelectorAll('.green-track');

    // Clear all cube placements when switching modes
    clearAllCubes(true);

    if (currentMode === 'blue') {
        // Switch to green
        currentMode = 'green';
        goalCard.classList.remove('green-side');
        goalCard.classList.add('blue-side');
        toggleBtn.textContent = 'Switch to Blue Side';
        toggleBtn.classList.remove('green-mode');

        blueTracks.forEach(track => track.style.display = 'none');
        greenTracks.forEach(track => track.style.display = 'block');
    } else {
        // Switch to blue
        currentMode = 'blue';
        goalCard.classList.remove('blue-side');
        goalCard.classList.add('green-side');
        toggleBtn.textContent = 'Switch to Green Side';
        toggleBtn.classList.add('green-mode');

        blueTracks.forEach(track => track.style.display = 'block');
        greenTracks.forEach(track => track.style.display = 'none');
    }

    // Re-render cubes in visible boxes (will be empty after clear)
    renderAllCubes();

    // Save the mode change
    saveGameState();
}

// Render all cubes
function renderAllCubes() {
    for (const [key, colors] of Object.entries(gameState.cubePlacements)) {
        const parts = key.split('-');
        const round = parseInt(parts[0]);
        const score = parseInt(parts[1]);
        const position = parts[2];
        const boxes = document.querySelectorAll(`.score-box[data-round="${round}"][data-score="${score}"][data-position="${position}"]`);
        boxes.forEach(box => {
            if (box.closest('.scoring-track').style.display !== 'none') {
                renderCubesInBox(box, round, score, position);
            }
        });
    }
}

// Save game state to localStorage
function saveGameState() {
    try {
        const stateToSave = {
            ...gameState,
            currentMode: currentMode
        };
        localStorage.setItem('wingspanGameState', JSON.stringify(stateToSave));
    } catch (e) {
        console.error('Failed to save game state:', e);
    }
}

// Load game state from localStorage
function loadGameState() {
    try {
        const saved = localStorage.getItem('wingspanGameState');
        if (saved) {
            const loaded = JSON.parse(saved);
            // Only load if player count matches
            const currentCount = parseInt(document.getElementById('numPlayers').value);
            if (loaded.players && loaded.players.length === currentCount) {
                // Restore currentMode (default to 'blue' if not found)
                const savedMode = loaded.currentMode || 'blue';

                // Extract gameState properties (exclude currentMode)
                gameState = {
                    players: loaded.players,
                    cubePlacements: loaded.cubePlacements || {},
                    goals: loaded.goals || null
                };

                renderPlayerList();
                renderScoreTable();

                // Restore goals if they exist in saved state
                if (gameState.goals) {
                    setGoalDisplay(gameState.goals);
                }

                // Apply the visual state for the saved mode
                applyModeVisualState(savedMode);

                // Render cubes after visual state is applied
                renderAllCubes();

                return true; // Successfully loaded saved state
            }
        }
    } catch (e) {
        console.error('Failed to load game state:', e);
    }
    return false; // No saved state found
}

// Apply visual state for the current mode
function applyModeVisualState(mode) {
    currentMode = mode;
    const goalCard = document.getElementById('goalCard');
    const toggleBtn = document.getElementById('toggleMode');
    const blueTracks = goalCard.querySelectorAll('.blue-track');
    const greenTracks = goalCard.querySelectorAll('.green-track');

    if (currentMode === 'green') {
        // Green mode
        goalCard.classList.remove('green-side');
        goalCard.classList.add('blue-side');
        toggleBtn.textContent = 'Switch to Blue Side';
        toggleBtn.classList.remove('green-mode');

        blueTracks.forEach(track => track.style.display = 'none');
        greenTracks.forEach(track => track.style.display = 'block');
    } else {
        // Blue mode
        goalCard.classList.remove('blue-side');
        goalCard.classList.add('green-side');
        toggleBtn.textContent = 'Switch to Green Side';
        toggleBtn.classList.add('green-mode');

        blueTracks.forEach(track => track.style.display = 'block');
        greenTracks.forEach(track => track.style.display = 'none');
    }
}
