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
