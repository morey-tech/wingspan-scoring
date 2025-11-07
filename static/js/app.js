// Track current mode
let currentMode = 'blue'; // 'blue' or 'green'

// Player colors
const PLAYER_COLORS = ['blue', 'purple', 'green', 'red', 'yellow'];
const PLAYER_COLOR_NAMES = ['Blue', 'Purple', 'Green', 'Red', 'Yellow'];
// Default player names
const DEFAULT_PLAYER_NAMES = ['Dani', 'Nick', 'Player 3', 'Player 4', 'Player 5'];
// Default color assignments: Player 1 = blue, Player 2 = yellow, then others in order
const DEFAULT_PLAYER_COLORS = ['blue', 'yellow', 'green', 'purple', 'red'];

// Goal ID to Sprite ID Mapping
const goalIdToSpriteId = {
    'base-birds-forest': 'g-birds-in-forest',
    'base-birds-grassland': 'g-birds-in-grassland',
    'base-birds-wetland': 'g-birds-in-wetland',
    'base-birds-bowl-egg': 'g-bowl-with-egg',
    'base-birds-cavity-egg': 'g-cavity-with-egg',
    'base-birds-ground-egg': 'g-ground-with-egg',
    'base-birds-platform-egg': 'g-platform-with-egg',
    'base-eggs-forest': 'g-eggs-in-forest',
    'base-eggs-grassland': 'g-eggs-in-grassland',
    'base-eggs-wetland': 'g-eggs-in-wetland',
    'base-eggs-bowl': 'g-eggs-in-bowl',
    'base-eggs-cavity': 'g-eggs-in-cavity',
    'base-eggs-ground': 'g-eggs-in-ground',
    'base-eggs-platform': 'g-eggs-in-platform',
    'base-egg-sets': 'g-egg-habitat-sets',
    'base-total-birds': 'g-total-birds',
    'eu-birds-tucked': 'g-birds-with-tucked-cards',
    'eu-food-cost': 'g-birds-food-cost',
    'eu-birds-one-row': 'g-birds-in-one-row',
    'eu-filled-columns': 'g-filled-columns',
    'eu-brown-powers': 'g-brown-powers',
    'eu-white-no-powers': 'g-white-no-powers',
    'eu-birds-high-value': 'g-birds-over-4pt',
    'eu-birds-no-eggs': 'g-eggless-birds',
    'eu-food-supply': 'g-food-owned',
    'eu-cards-hand': 'g-birds-in-hand',
    'oc-beak-left': 'g-beak-lt',
    'oc-beak-right': 'g-beak-rt',
    'oc-invertebrate-cost': 'g-inv-in-cost',
    'oc-fruit-seed-cost': 'g-fruit-seed-in-cost',
    'oc-no-goal': 'g-no-goal',
    'oc-rat-fish-cost': 'g-rodent-fish-in-cost',
    'oc-cubes-play-bird': 'g-cubes-on-play-bird',
    'oc-birds-low-value': 'g-birds-3pt-or-under'
};

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

// Load SVG sprite sheet and inline it
async function loadSpriteSheet() {
    try {
        const response = await fetch('/static/images/svg/wingspan-sprites.svg');
        if (!response.ok) {
            throw new Error(`Failed to load sprite sheet: ${response.status}`);
        }
        const svgText = await response.text();

        // Create hidden container for sprite definitions
        const container = document.createElement('div');
        container.style.display = 'none';
        container.innerHTML = svgText;
        document.body.insertBefore(container, document.body.firstChild);

        console.log('SVG sprite sheet loaded successfully');
    } catch (error) {
        console.error('Error loading sprite sheet:', error);
    }
}

// Update goal tiles with correct sprite references
function updateGoalTiles() {
    const goalInfoElements = document.querySelectorAll('.goal-info[data-goal-id]');

    goalInfoElements.forEach(element => {
        const goalId = element.getAttribute('data-goal-id');
        const spriteId = goalIdToSpriteId[goalId];

        if (spriteId) {
            const useElement = element.querySelector('.goal-tile-sprite');
            if (useElement) {
                useElement.setAttribute('href', `#${spriteId}`);
            }
        } else {
            console.warn(`No sprite mapping found for goal ID: ${goalId}`);
        }
    });
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

// Update visual state of score boxes based on "No Goal" status
function updateScoreBoxesVisualState() {
    // Update all rounds
    for (let round = 1; round <= 4; round++) {
        const scoreBoxes = document.querySelectorAll(`.score-box[data-round="${round}"]`);
        const hasNoGoal = isNoGoal(round);

        scoreBoxes.forEach(box => {
            const score = parseInt(box.dataset.score);
            if (hasNoGoal && score !== 0) {
                box.classList.add('disabled-no-goal');
                box.style.opacity = '0.3';
                box.style.cursor = 'not-allowed';
            } else {
                box.classList.remove('disabled-no-goal');
                box.style.opacity = '';
                box.style.cursor = '';
            }
        });
    }
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

        // Update the data-goal-id attribute
        goalInfo.setAttribute('data-goal-id', selectedGoal.id);
    }

    // Update the goal tile sprite
    updateGoalTiles();

    // Update visual state of score boxes
    updateScoreBoxesVisualState();

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

    // Initialize game end section
    initializeGameEndSection();

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
            checkbox.addEventListener('change', async (e) => {
                await fetchAllGoals();

                // Special handling for Oceania checkbox - also controls nectar visibility
                if (id === 'oceania') {
                    gameEndState.includeOceania = e.target.checked;
                    const nectarHeaders = document.querySelectorAll('.nectar-header');
                    const nectarCells = document.querySelectorAll('.nectar-cell');

                    nectarHeaders.forEach(header => {
                        header.style.display = gameEndState.includeOceania ? '' : 'none';
                    });

                    nectarCells.forEach(cell => {
                        cell.style.display = gameEndState.includeOceania ? '' : 'none';
                    });

                    // Hide/show nectar rules
                    const nectarRules = document.querySelector('.nectar-rules');
                    if (nectarRules) {
                        nectarRules.style.display = gameEndState.includeOceania ? '' : 'none';
                    }

                    // Hide/show nectar breakdown in results
                    const nectarBreakdown = document.getElementById('nectarBreakdown');
                    if (nectarBreakdown) {
                        nectarBreakdown.style.display = gameEndState.includeOceania ? '' : 'none';
                    }

                    saveGameEndState();
                }
            });
        }
    });

    // Initialize goals and sprites
    (async function() {
        // Load sprite sheet first
        await loadSpriteSheet();

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

        // Update goal tiles with sprite references
        updateGoalTiles();

        // Update visual state of score boxes for "No Goal" rounds
        updateScoreBoxesVisualState();

        // Initialize round highlighting on page load
        if (!hasSavedState) {
            updateCurrentRoundHighlight();
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
    // Update game end section with new player count
    if (typeof generateGameEndPlayerRows === 'function') {
        generateGameEndPlayerRows();
    }
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
            // Update game end section with new player name
            if (typeof generateGameEndPlayerRows === 'function') {
                generateGameEndPlayerRows();
            }
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

// Get the current round (first round where not all players have scored)
function getCurrentRound() {
    // If no players, return null
    if (gameState.players.length === 0) {
        return null;
    }

    // Check each round in order
    for (let round = 1; round <= 4; round++) {
        // Check if any player hasn't scored this round
        for (const player of gameState.players) {
            let hasScored = false;

            // Check all cube placements for this player in this round
            for (const [key, colors] of Object.entries(gameState.cubePlacements)) {
                const parts = key.split('-');
                const r = parseInt(parts[0]);

                if (r === round && colors.includes(player.color)) {
                    hasScored = true;
                    break;
                }
            }

            // If this player hasn't scored this round, it's the current round
            if (!hasScored) {
                return round;
            }
        }
    }

    // All rounds complete
    return null;
}

// Update the visual highlighting of the current round
function updateCurrentRoundHighlight() {
    // Get all round rows
    const roundRows = document.querySelectorAll('.round-row');

    // Clear all existing highlights
    roundRows.forEach(row => {
        row.classList.remove('current-round', 'completed-round');
    });

    // Determine current round
    const currentRound = getCurrentRound();

    if (currentRound === null) {
        // All rounds complete - mark all as completed
        roundRows.forEach(row => {
            row.classList.add('completed-round');
        });
        return;
    }

    // Apply highlighting based on round status
    roundRows.forEach((row, index) => {
        const round = index + 1;
        if (round < currentRound) {
            row.classList.add('completed-round');
        } else if (round === currentRound) {
            row.classList.add('current-round');
        }
        // Future rounds (round > currentRound) get no special class
    });
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

// Check if a round has "No Goal" selected
function isNoGoal(round) {
    const roundKey = `round${round}`;
    return gameState.goals &&
           gameState.goals[roundKey] &&
           gameState.goals[roundKey].id === 'oc-no-goal';
}

// Show warning when trying to place cube on non-zero box during "No Goal" round
function showNoGoalWarning(box) {
    // Remove any existing warning
    const existingWarning = document.querySelector('.no-goal-warning');
    if (existingWarning) {
        existingWarning.remove();
    }

    const warning = document.createElement('div');
    warning.className = 'no-goal-warning';
    warning.textContent = 'No Goal: Only 0 points allowed';
    warning.style.position = 'fixed';
    warning.style.backgroundColor = '#ff6b6b';
    warning.style.color = 'white';
    warning.style.padding = '8px 12px';
    warning.style.borderRadius = '4px';
    warning.style.fontSize = '14px';
    warning.style.zIndex = '10000';
    warning.style.boxShadow = '0 2px 8px rgba(0,0,0,0.2)';

    const rect = box.getBoundingClientRect();
    warning.style.left = `${rect.left}px`;
    warning.style.top = `${rect.top - 40}px`;

    document.body.appendChild(warning);

    // Auto-remove after 2 seconds
    setTimeout(() => {
        warning.remove();
    }, 2000);
}

// Handle score box click
function handleScoreBoxClick(e) {
    const box = e.currentTarget;
    const round = parseInt(box.dataset.round);
    const score = parseInt(box.dataset.score);
    const position = box.dataset.position;

    // Check if this is a "No Goal" round and restrict to zero points only
    if (isNoGoal(round) && score !== 0) {
        // Show a brief message and prevent interaction
        showNoGoalWarning(box);
        return;
    }

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

    // Check if player is already in this box
    const currentPlacements = gameState.cubePlacements[key] || [];
    const isAlreadyPlaced = currentPlacements.includes(playerColor);

    if (isAlreadyPlaced) {
        // Deselect: Remove player from this box
        const index = currentPlacements.indexOf(playerColor);
        if (index > -1) {
            currentPlacements.splice(index, 1);
            if (currentPlacements.length === 0) {
                delete gameState.cubePlacements[key];
            }
        }
    } else {
        // Select/Move: Remove from other boxes in this round, then add to this box
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
        gameState.cubePlacements[key].push(playerColor);
    }

    // Re-render all boxes in this round
    const roundBoxes = document.querySelectorAll(`.score-box[data-round="${round}"]`);
    roundBoxes.forEach(box => {
        const boxScore = parseInt(box.dataset.score);
        const boxPosition = box.dataset.position;
        renderCubesInBox(box, round, boxScore, boxPosition);
    });

    // Update round highlighting
    updateCurrentRoundHighlight();
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
    updateCurrentRoundHighlight();
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
        clearAllGameEndScores(true); // Clear game-end scores without confirmation

        // Re-fetch goals based on new expansion selection
        await fetchAllGoals();

        // Update goal tiles after new game
        updateGoalTiles();

        // Update visual state of score boxes for "No Goal" rounds
        updateScoreBoxesVisualState();
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
            const goalInfoElement = row.querySelector('.goal-info');

            if (nameElement) nameElement.textContent = goal.name;
            if (descElement) descElement.textContent = goal.description;
            if (goalInfoElement && goal.id) {
                goalInfoElement.setAttribute('data-goal-id', goal.id);
            }
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
    let isConfirmed = confirm("Are you sure you want to switch sides? This will remove all cube placements!")
    if(!isConfirmed){
        return
    }
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

        // Auto-sync round goals to game end section
        if (typeof updateGameEndRoundGoals === 'function') {
            updateGameEndRoundGoals();
        }
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

                // Update round highlighting based on loaded state
                updateCurrentRoundHighlight();

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

// ============================================================================
// GAME END SCORING FUNCTIONALITY
// ============================================================================

const GAME_END_STORAGE_KEY = 'wingspan-game-end-scores';

// Game end state
let gameEndState = {
    includeOceania: true,
    players: []
};

// Initialize game end section
function initializeGameEndSection() {
    const calculateBtn = document.getElementById('calculateBtn');
    const clearGameEndBtn = document.getElementById('clearGameEndBtn');

    if (calculateBtn) {
        calculateBtn.addEventListener('click', calculateGameEndScores);
    }
    if (clearGameEndBtn) {
        clearGameEndBtn.addEventListener('click', clearAllGameEndScores);
    }

    // Load saved game end state
    loadGameEndState();

    // Generate player rows for game end table
    generateGameEndPlayerRows();
}

// Generate player rows for game end table
function generateGameEndPlayerRows() {
    const tbody = document.getElementById('gameEndTableBody');
    if (!tbody) return;

    tbody.innerHTML = '';
    const numPlayers = gameState.players.length;

    // Helper function to calculate column-major tabindex
    // This makes tab navigation go down columns (between players) instead of across rows
    const getTabIndex = (playerIndex, columnIndex) => {
        return (columnIndex * numPlayers) + playerIndex + 1;
    };

    for (let i = 0; i < numPlayers; i++) {
        const player = gameState.players[i];
        const playerNum = i + 1;
        const playerName = player.name || `Player ${playerNum}`;
        const playerColor = player.color || PLAYER_COLORS[i];

        // Calculate round goal score for this player
        const roundGoalScore = calculatePlayerRoundGoalScore(player.color);

        const row = document.createElement('tr');
        row.className = 'player-row';
        row.dataset.playerIndex = playerNum;

        row.innerHTML = `
            <td class="player-name-cell">
                <div class="player-name-wrapper">
                    <span class="player-color-indicator ${playerColor}"></span>
                    <span class="player-name-display">${playerName}</span>
                </div>
            </td>
            <td class="score-cell">
                <input type="number"
                       inputmode="numeric"
                       class="score-input game-end-input"
                       min="0"
                       placeholder="0"
                       tabindex="${getTabIndex(i, 0)}"
                       data-player="${playerNum}"
                       data-field="birdPoints">
            </td>
            <td class="score-cell">
                <input type="number"
                       inputmode="numeric"
                       class="score-input game-end-input"
                       min="0"
                       placeholder="0"
                       tabindex="${getTabIndex(i, 1)}"
                       data-player="${playerNum}"
                       data-field="bonusCards">
            </td>
            <td class="score-cell">
                <input type="number"
                       inputmode="numeric"
                       class="score-input game-end-input"
                       min="0"
                       value="${roundGoalScore}"
                       data-player="${playerNum}"
                       data-field="roundGoals"
                       readonly
                       tabindex="-1"
                       style="background-color: #f0f0f0;">
            </td>
            <td class="score-cell">
                <input type="number"
                       inputmode="numeric"
                       class="score-input game-end-input"
                       min="0"
                       placeholder="0"
                       tabindex="${getTabIndex(i, 2)}"
                       data-player="${playerNum}"
                       data-field="eggs">
            </td>
            <td class="score-cell">
                <input type="number"
                       inputmode="numeric"
                       class="score-input game-end-input"
                       min="0"
                       placeholder="0"
                       tabindex="${getTabIndex(i, 3)}"
                       data-player="${playerNum}"
                       data-field="cachedFood">
            </td>
            <td class="score-cell">
                <input type="number"
                       inputmode="numeric"
                       class="score-input game-end-input"
                       min="0"
                       placeholder="0"
                       tabindex="${getTabIndex(i, 4)}"
                       data-player="${playerNum}"
                       data-field="tuckedCards">
            </td>
            <td class="score-cell nectar-cell">
                <input type="number"
                       inputmode="numeric"
                       class="score-input nectar-input game-end-input"
                       min="0"
                       placeholder="0"
                       tabindex="${getTabIndex(i, 5)}"
                       data-player="${playerNum}"
                       data-field="nectarForest">
            </td>
            <td class="score-cell nectar-cell">
                <input type="number"
                       inputmode="numeric"
                       class="score-input nectar-input game-end-input"
                       min="0"
                       placeholder="0"
                       tabindex="${getTabIndex(i, 6)}"
                       data-player="${playerNum}"
                       data-field="nectarGrassland">
            </td>
            <td class="score-cell nectar-cell">
                <input type="number"
                       inputmode="numeric"
                       class="score-input nectar-input game-end-input"
                       min="0"
                       placeholder="0"
                       tabindex="${getTabIndex(i, 7)}"
                       data-player="${playerNum}"
                       data-field="nectarWetland">
            </td>
            <td class="score-cell tiebreaker-cell">
                <input type="number"
                       inputmode="numeric"
                       class="score-input tiebreaker-input game-end-input"
                       min="0"
                       placeholder="0"
                       tabindex="${getTabIndex(i, 8)}"
                       data-player="${playerNum}"
                       data-field="unusedFood">
            </td>
            <td class="total-cell">
                <span class="total-display">0</span>
            </td>
        `;

        tbody.appendChild(row);
    }

    // Add event listeners to all game end inputs
    document.querySelectorAll('.game-end-input').forEach(input => {
        input.addEventListener('change', saveGameEndState);
        input.addEventListener('focus', (e) => {
            if (e.target.value) {
                e.target.select();
            }
        });
        // Add Enter key handling to navigate like Tab
        input.addEventListener('keydown', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                const currentTabIndex = parseInt(input.getAttribute('tabindex'));
                // Find the next input with the next tabindex
                const nextInput = document.querySelector(`.game-end-input[tabindex="${currentTabIndex + 1}"]`);
                if (nextInput) {
                    nextInput.focus();
                }
            }
        });
    });

    // Apply saved state
    applyGameEndSavedState();

    // Update nectar visibility
    const oceaniaToggle = document.getElementById('oceaniaToggle');
    if (oceaniaToggle) {
        handleOceaniaToggle({ target: oceaniaToggle });
    }
}

// Calculate round goal score for a specific player color
function calculatePlayerRoundGoalScore(playerColor) {
    let totalScore = 0;

    // For each round (1-4), find the score for this player's color
    for (let round = 1; round <= 4; round++) {
        // Check all cube placement keys for this round
        for (const [key, colors] of Object.entries(gameState.cubePlacements)) {
            const parts = key.split('-');
            const r = parseInt(parts[0]);
            const score = parseInt(parts[1]);

            // If this key is for the current round and contains this player's color
            if (r === round && colors && colors.includes(playerColor)) {
                totalScore += score;
                break; // Found score for this round, move to next round
            }
        }
    }

    return totalScore;
}

// Update game end section when round scores change (auto-sync)
function updateGameEndRoundGoals() {
    const numPlayers = gameState.players.length;

    for (let i = 0; i < numPlayers; i++) {
        const player = gameState.players[i];
        const playerNum = i + 1;
        const roundGoalScore = calculatePlayerRoundGoalScore(player.color);

        const roundGoalsInput = document.querySelector(
            `.game-end-input[data-player="${playerNum}"][data-field="roundGoals"]`
        );

        if (roundGoalsInput) {
            roundGoalsInput.value = roundGoalScore;
        }
    }
}

// Calculate game end scores
async function calculateGameEndScores() {
    const players = [];

    // Gather player data
    for (let i = 0; i < gameState.players.length; i++) {
        const player = gameState.players[i];
        const playerNum = i + 1;
        const playerName = player.name || `Player ${playerNum}`;

        const playerData = {
            playerName: playerName,
            birdPoints: parseInt(document.querySelector(`.game-end-input[data-player="${playerNum}"][data-field="birdPoints"]`).value) || 0,
            bonusCards: parseInt(document.querySelector(`.game-end-input[data-player="${playerNum}"][data-field="bonusCards"]`).value) || 0,
            roundGoals: parseInt(document.querySelector(`.game-end-input[data-player="${playerNum}"][data-field="roundGoals"]`).value) || 0,
            eggs: parseInt(document.querySelector(`.game-end-input[data-player="${playerNum}"][data-field="eggs"]`).value) || 0,
            cachedFood: parseInt(document.querySelector(`.game-end-input[data-player="${playerNum}"][data-field="cachedFood"]`).value) || 0,
            tuckedCards: parseInt(document.querySelector(`.game-end-input[data-player="${playerNum}"][data-field="tuckedCards"]`).value) || 0,
            nectarForest: parseInt(document.querySelector(`.game-end-input[data-player="${playerNum}"][data-field="nectarForest"]`).value) || 0,
            nectarGrassland: parseInt(document.querySelector(`.game-end-input[data-player="${playerNum}"][data-field="nectarGrassland"]`).value) || 0,
            nectarWetland: parseInt(document.querySelector(`.game-end-input[data-player="${playerNum}"][data-field="nectarWetland"]`).value) || 0,
            unusedFood: parseInt(document.querySelector(`.game-end-input[data-player="${playerNum}"][data-field="unusedFood"]`).value) || 0
        };

        players.push(playerData);
    }

    try {
        const response = await fetch('/api/calculate-game-end', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                players: players,
                includeOceania: gameEndState.includeOceania
            })
        });

        if (!response.ok) {
            throw new Error('Failed to calculate scores');
        }

        const result = await response.json();
        displayGameEndResults(result);
    } catch (error) {
        console.error('Error calculating scores:', error);
        alert('Error calculating game end scores. Please try again.');
    }
}

// Display game end results
function displayGameEndResults(result) {
    const resultsSection = document.getElementById('resultsSection');
    resultsSection.style.display = 'block';

    // Clear all winner highlighting first
    document.querySelectorAll('#gameEndTableBody .player-row').forEach(row => {
        row.classList.remove('winner-row');
    });

    // Update total cells in the table - match by player name
    result.players.forEach(player => {
        // Find the row with this player's name
        const nameDisplays = document.querySelectorAll('.player-name-display');
        let playerRowIndex = null;

        nameDisplays.forEach((display, idx) => {
            if (display.textContent === player.playerName) {
                playerRowIndex = idx + 1;
            }
        });

        if (playerRowIndex) {
            const totalCell = document.querySelector(`#gameEndTableBody tr[data-player-index="${playerRowIndex}"] .total-display`);
            if (totalCell) {
                totalCell.textContent = player.total;
            }

            // Highlight winner(s)
            const row = document.querySelector(`#gameEndTableBody tr[data-player-index="${playerRowIndex}"]`);
            if (row && player.rank === 1) {
                row.classList.add('winner-row');
            }
        }
    });

    // Display nectar breakdown
    if (gameEndState.includeOceania) {
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

// Clear all game end scores
function clearAllGameEndScores(skipConfirm = false) {
    if (!skipConfirm && !confirm('Clear all game end scores? This cannot be undone.')) {
        return;
    }

    document.querySelectorAll('.game-end-input:not([readonly])').forEach(input => {
        input.value = '';
    });

    document.querySelectorAll('#gameEndTableBody .total-display').forEach(cell => {
        cell.textContent = '0';
    });

    document.querySelectorAll('#gameEndTableBody .player-row').forEach(row => {
        row.classList.remove('winner-row');
    });

    const resultsSection = document.getElementById('resultsSection');
    resultsSection.style.display = 'none';

    saveGameEndState();
}

// Save game end state to localStorage
function saveGameEndState() {
    const state = {
        includeOceania: gameEndState.includeOceania,
        players: []
    };

    // Save all player data
    for (let i = 1; i <= gameState.players.length; i++) {
        const playerData = {};
        document.querySelectorAll(`.game-end-input[data-player="${i}"]`).forEach(input => {
            playerData[input.dataset.field] = input.value;
        });
        state.players.push(playerData);
    }

    localStorage.setItem(GAME_END_STORAGE_KEY, JSON.stringify(state));
}

// Load game end state from localStorage
function loadGameEndState() {
    const saved = localStorage.getItem(GAME_END_STORAGE_KEY);
    if (saved) {
        try {
            const state = JSON.parse(saved);
            gameEndState.includeOceania = state.includeOceania !== undefined ? state.includeOceania : true;
            gameEndState.players = state.players || [];

            // Update UI - sync with Oceania checkbox in Round Goals section
            const oceaniaCheckbox = document.getElementById('oceania');
            if (oceaniaCheckbox) {
                oceaniaCheckbox.checked = gameEndState.includeOceania;

                // Apply nectar visibility based on loaded state
                const nectarHeaders = document.querySelectorAll('.nectar-header');
                const nectarCells = document.querySelectorAll('.nectar-cell');

                nectarHeaders.forEach(header => {
                    header.style.display = gameEndState.includeOceania ? '' : 'none';
                });

                nectarCells.forEach(cell => {
                    cell.style.display = gameEndState.includeOceania ? '' : 'none';
                });

                // Hide/show nectar rules
                const nectarRules = document.querySelector('.nectar-rules');
                if (nectarRules) {
                    nectarRules.style.display = gameEndState.includeOceania ? '' : 'none';
                }

                // Hide/show nectar breakdown in results
                const nectarBreakdown = document.getElementById('nectarBreakdown');
                if (nectarBreakdown) {
                    nectarBreakdown.style.display = gameEndState.includeOceania ? '' : 'none';
                }
            }
        } catch (error) {
            console.error('Error loading game end state:', error);
        }
    } else {
        // No saved state - sync with current Oceania checkbox state
        const oceaniaCheckbox = document.getElementById('oceania');
        if (oceaniaCheckbox) {
            gameEndState.includeOceania = oceaniaCheckbox.checked;
        }
    }
}

// Apply saved game end state to inputs
function applyGameEndSavedState() {
    gameEndState.players.forEach((playerData, index) => {
        const playerNum = index + 1;
        if (playerNum > gameState.players.length) return;

        Object.keys(playerData).forEach(field => {
            // Skip roundGoals as it's auto-synced from round goals
            if (field === 'roundGoals') {
                return;
            }

            const input = document.querySelector(`.game-end-input[data-player="${playerNum}"][data-field="${field}"]`);
            if (input && playerData[field]) {
                input.value = playerData[field];
            }
        });
    });
}

// Navigation Toggle
document.addEventListener('DOMContentLoaded', () => {
    const toggleOptions = document.querySelectorAll('.toggle-option');
    const slider = document.querySelector('.toggle-slider');

    toggleOptions.forEach(option => {
        option.addEventListener('click', function() {
            const targetPage = this.dataset.page;

            // Determine target URL
            let targetUrl = '/';
            if (targetPage === 'history') {
                targetUrl = '/history';
            }

            // Animate slider to clicked position
            if (targetPage === 'home') {
                slider.classList.remove('toggle-right');
                slider.classList.add('toggle-left');
            } else {
                slider.classList.remove('toggle-left');
                slider.classList.add('toggle-right');
            }

            // Navigate after animation completes
            setTimeout(() => {
                window.location.href = targetUrl;
            }, 300);
        });
    });
});
