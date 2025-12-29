// Game History Page JavaScript

const ITEMS_PER_PAGE = 20;
let currentPage = 1;
let totalGames = 0;
let gameToDelete = null;

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    loadGames();
    initializeEventListeners();
    initializePlayerStats();
    initializeImportForm();
    initializeExportButton();
});

function initializeEventListeners() {
    const prevBtn = document.getElementById('prevPage');
    const nextBtn = document.getElementById('nextPage');
    const closeModalBtn = document.getElementById('closeModal');
    const closeConfirmBtn = document.getElementById('closeConfirmModal');
    const confirmDeleteBtn = document.getElementById('confirmDeleteBtn');
    const cancelDeleteBtn = document.getElementById('cancelDeleteBtn');

    if (prevBtn) {
        prevBtn.addEventListener('click', () => {
            if (currentPage > 1) {
                currentPage--;
                loadGames();
            }
        });
    }

    if (nextBtn) {
        nextBtn.addEventListener('click', () => {
            const maxPages = Math.ceil(totalGames / ITEMS_PER_PAGE);
            if (currentPage < maxPages) {
                currentPage++;
                loadGames();
            }
        });
    }

    if (closeModalBtn) {
        closeModalBtn.addEventListener('click', closeGameDetails);
    }

    if (closeConfirmBtn) {
        closeConfirmBtn.addEventListener('click', closeDeleteConfirmation);
    }

    if (confirmDeleteBtn) {
        confirmDeleteBtn.addEventListener('click', confirmDelete);
    }

    if (cancelDeleteBtn) {
        cancelDeleteBtn.addEventListener('click', closeDeleteConfirmation);
    }

    // Close modals when clicking outside
    const detailsModal = document.getElementById('gameDetailsModal');
    if (detailsModal) {
        detailsModal.addEventListener('click', (e) => {
            if (e.target === detailsModal) {
                closeGameDetails();
            }
        });
    }

    const confirmModal = document.getElementById('confirmDeleteModal');
    if (confirmModal) {
        confirmModal.addEventListener('click', (e) => {
            if (e.target === confirmModal) {
                closeDeleteConfirmation();
            }
        });
    }
}

async function loadGames() {
    const gamesList = document.getElementById('gamesList');
    gamesList.innerHTML = '<div class="loading">Loading game history...</div>';

    try {
        const offset = (currentPage - 1) * ITEMS_PER_PAGE;
        const response = await fetch(`/api/games?limit=${ITEMS_PER_PAGE}&offset=${offset}`);

        if (!response.ok) {
            throw new Error('Failed to load games');
        }

        const data = await response.json();
        totalGames = data.totalCount;

        // Update stats
        document.getElementById('totalGames').textContent = totalGames;

        // Display games
        displayGames(data.games);

        // Update pagination
        updatePagination();

    } catch (error) {
        console.error('Error loading games:', error);
        gamesList.innerHTML = '<div class="error">Failed to load game history. Please try again.</div>';
    }
}

function displayGames(games) {
    const gamesList = document.getElementById('gamesList');

    if (!games || games.length === 0) {
        gamesList.innerHTML = '<div class="no-games">No games recorded yet. Play a game and calculate game end scores to see it here!</div>';
        return;
    }

    gamesList.innerHTML = '';

    games.forEach(game => {
        const gameCard = createGameCard(game);
        gamesList.appendChild(gameCard);
    });
}

function createGameCard(game) {
    const card = document.createElement('div');
    card.className = 'game-card';
    card.dataset.gameId = game.id;

    const date = new Date(game.createdAt);
    const formattedDate = formatDate(date);

    // Get all player names for display
    const playerNames = game.players.map(p => p.playerName).join(', ');

    card.innerHTML = `
        <div class="game-card-header">
            <div class="game-date">${formattedDate}</div>
            <div class="game-badges">
                <div class="game-badge">${game.numPlayers} Players</div>
                ${game.includeOceania ? '<div class="game-badge oceania">Oceania</div>' : ''}
            </div>
        </div>
        <div class="game-card-body">
            <div class="game-winner">
                <span class="winner-icon">🏆</span>
                <span class="winner-name">${game.winnerName}</span>
                <span class="winner-score">${game.winnerScore} points</span>
            </div>
            <div class="game-players">
                <strong>Players:</strong> ${playerNames}
            </div>
        </div>
        <div class="game-card-footer">
            <button class="btn-view-details" onclick="viewGameDetails(${game.id})">View Details</button>
            <button class="btn-delete" onclick="showDeleteConfirmation(${game.id})">Delete</button>
        </div>
    `;

    return card;
}

function formatDate(date) {
    const options = {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    };
    return date.toLocaleDateString('en-US', options);
}

function updatePagination() {
    const pagination = document.getElementById('pagination');
    const prevBtn = document.getElementById('prevPage');
    const nextBtn = document.getElementById('nextPage');
    const pageInfo = document.getElementById('pageInfo');

    const maxPages = Math.ceil(totalGames / ITEMS_PER_PAGE);

    if (maxPages <= 1) {
        pagination.style.display = 'none';
        return;
    }

    pagination.style.display = 'flex';
    pageInfo.textContent = `Page ${currentPage} of ${maxPages}`;

    prevBtn.disabled = currentPage === 1;
    nextBtn.disabled = currentPage === maxPages;
}

async function viewGameDetails(gameId) {
    const modal = document.getElementById('gameDetailsModal');
    const modalBody = document.getElementById('modalBody');

    modalBody.innerHTML = '<div class="loading">Loading game details...</div>';
    modal.style.display = 'flex';

    try {
        const response = await fetch(`/api/games/${gameId}`);

        if (!response.ok) {
            throw new Error('Failed to load game details');
        }

        const game = await response.json();
        displayGameDetails(game);

    } catch (error) {
        console.error('Error loading game details:', error);
        modalBody.innerHTML = '<div class="error">Failed to load game details.</div>';
    }
}

function displayGameDetails(game) {
    const modalBody = document.getElementById('modalBody');

    const date = new Date(game.createdAt);
    const formattedDate = formatDate(date);

    let html = `
        <div class="game-details">
            <div class="details-header">
                <div class="detail-item">
                    <strong>Date:</strong> ${formattedDate}
                </div>
                <div class="detail-item">
                    <strong>Players:</strong> ${game.numPlayers}
                </div>
                <div class="detail-item">
                    <strong>Expansion:</strong> ${game.includeOceania ? 'With Oceania' : 'Base Game'}
                </div>
            </div>

            <div class="details-scores">
                <h4>Game End Scores</h4>
                <table class="details-table">
                    <thead>
                        <tr>
                            <th>Rank</th>
                            <th>Player</th>
                            <th>Birds</th>
                            <th>Bonus</th>
                            <th>Goals</th>
                            <th>Eggs</th>
                            <th>Food</th>
                            <th>Cards</th>
                            ${game.includeOceania ? '<th>Nectar</th>' : ''}
                            <th class="total-col">Total</th>
                        </tr>
                    </thead>
                    <tbody>
    `;

    game.players.forEach(player => {
        const isWinner = player.rank === 1;
        const nectarTotal = game.includeOceania ?
            (game.nectarScoring.forest[player.playerName] || 0) +
            (game.nectarScoring.grassland[player.playerName] || 0) +
            (game.nectarScoring.wetland[player.playerName] || 0) : 0;

        html += `
            <tr class="${isWinner ? 'winner-row' : ''}">
                <td>${getRankDisplay(player.rank)}</td>
                <td><strong>${player.playerName}</strong></td>
                <td>${player.birdPoints}</td>
                <td>${player.bonusCards}</td>
                <td>${player.roundGoals}</td>
                <td>${player.eggs}</td>
                <td>${player.cachedFood}</td>
                <td>${player.tuckedCards}</td>
                ${game.includeOceania ? `<td>${nectarTotal}</td>` : ''}
                <td class="total-col"><strong>${player.total}</strong></td>
            </tr>
        `;
    });

    html += `
                    </tbody>
                </table>
            </div>
    `;

    // Add round goals breakdown if available
    if (game.roundBreakdown && Object.keys(game.roundBreakdown).length > 0) {
        html += `
            <div class="details-round-breakdown">
                <h4>Round Goals Breakdown</h4>
                <table class="details-table">
                    <thead>
                        <tr>
                            <th>Player</th>
                            <th>Round 1</th>
                            <th>Round 2</th>
                            <th>Round 3</th>
                            <th>Round 4</th>
                            <th class="total-col">Total</th>
                        </tr>
                    </thead>
                    <tbody>
        `;

        game.players.forEach(player => {
            const breakdown = game.roundBreakdown[player.playerName];
            if (breakdown) {
                const total = breakdown.round1 + breakdown.round2 + breakdown.round3 + breakdown.round4;
                html += `
                    <tr>
                        <td><strong>${player.playerName}</strong></td>
                        <td>${breakdown.round1}</td>
                        <td>${breakdown.round2}</td>
                        <td>${breakdown.round3}</td>
                        <td>${breakdown.round4}</td>
                        <td class="total-col"><strong>${total}</strong></td>
                    </tr>
                `;
            }
        });

        html += `
                    </tbody>
                </table>
            </div>
        `;
    }

    // Add nectar breakdown if Oceania is included
    if (game.includeOceania && game.nectarScoring) {
        html += `
            <div class="details-nectar">
                <h4>Nectar Scoring Breakdown</h4>
                <div class="nectar-breakdown-grid">
                    ${formatNectarHabitat('Forest 🌲', game.nectarScoring.forest, game.players)}
                    ${formatNectarHabitat('Grassland 🌾', game.nectarScoring.grassland, game.players)}
                    ${formatNectarHabitat('Wetland 💧', game.nectarScoring.wetland, game.players)}
                </div>
            </div>
        `;
    }

    html += '</div>';

    modalBody.innerHTML = html;
}

function formatNectarHabitat(habitatName, scoring, players) {
    let html = `
        <div class="nectar-habitat-detail">
            <h5>${habitatName}</h5>
            <div class="nectar-players">
    `;

    // Sort players by points awarded (descending)
    const sorted = Object.entries(scoring)
        .sort((a, b) => b[1] - a[1])
        .filter(([_, points]) => points > 0);

    if (sorted.length === 0) {
        html += '<div class="no-nectar">No nectar scored</div>';
    } else {
        sorted.forEach(([playerName, points]) => {
            // Find the player to get their nectar count
            const player = players.find(p => p.playerName === playerName);
            const nectarCount = getNectarCountForHabitat(player, habitatName);

            html += `
                <div class="nectar-player-item">
                    <span class="nectar-player-name">${playerName}</span>
                    <span class="nectar-player-count">(${nectarCount} nectar)</span>
                    <span class="nectar-player-points">${points} pts</span>
                </div>
            `;
        });
    }

    html += `
            </div>
        </div>
    `;

    return html;
}

function getNectarCountForHabitat(player, habitatName) {
    if (!player) return 0;

    if (habitatName.includes('Forest')) return player.nectarForest || 0;
    if (habitatName.includes('Grassland')) return player.nectarGrassland || 0;
    if (habitatName.includes('Wetland')) return player.nectarWetland || 0;

    return 0;
}

function getRankDisplay(rank) {
    if (rank === 1) return '🥇 1st';
    if (rank === 2) return '🥈 2nd';
    if (rank === 3) return '🥉 3rd';
    return `${rank}th`;
}

function closeGameDetails() {
    const modal = document.getElementById('gameDetailsModal');
    modal.style.display = 'none';
}

function showDeleteConfirmation(gameId) {
    gameToDelete = gameId;
    const modal = document.getElementById('confirmDeleteModal');
    modal.style.display = 'flex';
}

function closeDeleteConfirmation() {
    gameToDelete = null;
    const modal = document.getElementById('confirmDeleteModal');
    modal.style.display = 'none';
}

async function confirmDelete() {
    if (!gameToDelete) return;

    try {
        const response = await fetch(`/api/games/${gameToDelete}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            throw new Error('Failed to delete game');
        }

        // Close the confirmation modal
        closeDeleteConfirmation();

        // Reload the games list
        await loadGames();

    } catch (error) {
        console.error('Error deleting game:', error);
        alert('Failed to delete game. Please try again.');
    }
}

// ===== Leaderboard Functions =====

function initializePlayerStats() {
    // Load leaderboard on page load
    loadLeaderboard();
}

async function loadLeaderboard() {
    const statsContent = document.getElementById('statsContent');
    statsContent.innerHTML = '<div class="loading">Loading leaderboard...</div>';

    try {
        const response = await fetch('/api/leaderboard');

        if (!response.ok) {
            throw new Error('Failed to load leaderboard');
        }

        const leaderboard = await response.json();
        displayLeaderboard(leaderboard);

    } catch (error) {
        console.error('Error loading leaderboard:', error);
        statsContent.innerHTML = '<div class="error">Failed to load leaderboard. Please try again.</div>';
    }
}

function displayLeaderboard(leaderboard) {
    const statsContent = document.getElementById('statsContent');

    // Helper function to create a leaderboard card
    function createLeaderCard(icon, category, leader) {
        const playerName = leader.playerName || 'N/A';
        const score = leader.score || 0;

        return `
            <div class="stat-card">
                <div class="stat-icon">${icon}</div>
                <div class="stat-value">${score}</div>
                <div class="stat-label">${category}</div>
                <div class="stat-leader">${playerName}</div>
            </div>
        `;
    }

    statsContent.innerHTML = `
        <div class="stats-cards">
            ${createLeaderCard('🏆', 'Total Score', leaderboard.totalScore)}
            ${createLeaderCard('🐦', 'Bird Points', leaderboard.birdPoints)}
            ${createLeaderCard('🎯', 'Bonus Cards', leaderboard.bonusCards)}
            ${createLeaderCard('🎪', 'Round Goals', leaderboard.roundGoals)}
            ${createLeaderCard('🥚', 'Eggs', leaderboard.eggs)}
            ${createLeaderCard('🍎', 'Cached Food', leaderboard.cachedFood)}
            ${createLeaderCard('🃏', 'Tucked Cards', leaderboard.tuckedCards)}
            ${createLeaderCard('🌲', 'Nectar Forest', leaderboard.nectarForest)}
            ${createLeaderCard('🌾', 'Nectar Grassland', leaderboard.nectarGrassland)}
            ${createLeaderCard('💧', 'Nectar Wetland', leaderboard.nectarWetland)}
        </div>
    `;
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

// Import functionality
function initializeImportForm() {
    const importForm = document.getElementById('importForm');
    const csvFileInput = document.getElementById('csvFile');
    const fileNameSpan = document.getElementById('fileName');

    if (csvFileInput && fileNameSpan) {
        csvFileInput.addEventListener('change', (e) => {
            const fileName = e.target.files[0]?.name || 'No file chosen';
            fileNameSpan.textContent = fileName;
        });
    }

    if (importForm) {
        importForm.addEventListener('submit', handleImportSubmit);
    }
}

async function handleImportSubmit(e) {
    e.preventDefault();

    const csvFileInput = document.getElementById('csvFile');
    const importStatus = document.getElementById('importStatus');
    const importProgress = document.getElementById('importProgress');
    const importErrors = document.getElementById('importErrors');

    // Hide previous messages
    importStatus.style.display = 'none';
    importErrors.style.display = 'none';

    // Validate file selection
    if (!csvFileInput.files || csvFileInput.files.length === 0) {
        showImportError('Please select a CSV file to import.');
        return;
    }

    const file = csvFileInput.files[0];

    // Validate file type
    if (!file.name.endsWith('.csv')) {
        showImportError('Please select a valid CSV file.');
        return;
    }

    // Show progress
    importProgress.style.display = 'block';

    // Create form data
    const formData = new FormData();
    formData.append('csvFile', file);

    try {
        const response = await fetch('/api/import', {
            method: 'POST',
            body: formData
        });

        const result = await response.json();

        // Hide progress
        importProgress.style.display = 'none';

        if (result.success) {
            showImportSuccess(`Successfully imported ${result.gamesImported} game(s)!`);

            // Reset form
            csvFileInput.value = '';
            document.getElementById('fileName').textContent = 'No file chosen';

            // Reload games list and stats
            setTimeout(() => {
                loadGames();
                initializePlayerStats();
            }, 1000);
        } else {
            showImportError(`Import failed: ${result.message}`);

            // Display detailed errors if available
            if (result.errors && result.errors.length > 0) {
                displayImportErrors(result.errors);
            }
        }
    } catch (error) {
        importProgress.style.display = 'none';
        showImportError(`Import failed: ${error.message}`);
    }
}

function showImportSuccess(message) {
    const importStatus = document.getElementById('importStatus');
    importStatus.textContent = message;
    importStatus.className = 'import-status success';
    importStatus.style.display = 'block';

    // Auto-hide after 5 seconds
    setTimeout(() => {
        importStatus.style.display = 'none';
    }, 5000);
}

function showImportError(message) {
    const importStatus = document.getElementById('importStatus');
    importStatus.textContent = message;
    importStatus.className = 'import-status error';
    importStatus.style.display = 'block';
}

function displayImportErrors(errors) {
    const importErrors = document.getElementById('importErrors');

    let errorHTML = '<h4>Import Errors:</h4><ul>';
    errors.forEach(error => {
        const line = error.Line ? `Line ${error.Line}` : '';
        const gameID = error.GameID ? `Game ${error.GameID}` : '';
        const location = [line, gameID].filter(Boolean).join(' - ');
        errorHTML += `<li>${location ? `${location}: ` : ''}${error.Message}</li>`;
    });
    errorHTML += '</ul>';

    importErrors.innerHTML = errorHTML;
    importErrors.style.display = 'block';
}

// Export functionality
function initializeExportButton() {
    const exportBtn = document.getElementById('exportGamesBtn');
    if (exportBtn) {
        exportBtn.addEventListener('click', exportGames);
    }
}

async function exportGames() {
    const exportBtn = document.getElementById('exportGamesBtn');
    const importStatus = document.getElementById('importStatus');

    // Disable button while exporting
    exportBtn.disabled = true;
    exportBtn.textContent = 'Exporting...';

    // Hide previous status messages
    importStatus.style.display = 'none';

    try {
        const response = await fetch('/api/export');

        if (!response.ok) {
            throw new Error('Failed to export games');
        }

        // Get the CSV data as a blob
        const blob = await response.blob();

        // Create a download link and trigger it
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'wingspan-games-export.csv';
        document.body.appendChild(a);
        a.click();

        // Cleanup
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);

        showExportSuccess('Games exported successfully!');

    } catch (error) {
        console.error('Error exporting games:', error);
        showExportError(`Export failed: ${error.message}`);
    } finally {
        // Re-enable button
        exportBtn.disabled = false;
        exportBtn.textContent = 'Export All Games';
    }
}

function showExportSuccess(message) {
    const importStatus = document.getElementById('importStatus');
    importStatus.textContent = message;
    importStatus.className = 'import-status success';
    importStatus.style.display = 'block';

    // Auto-hide after 5 seconds
    setTimeout(() => {
        importStatus.style.display = 'none';
    }, 5000);
}

function showExportError(message) {
    const importStatus = document.getElementById('importStatus');
    importStatus.textContent = message;
    importStatus.className = 'import-status error';
    importStatus.style.display = 'block';
}
