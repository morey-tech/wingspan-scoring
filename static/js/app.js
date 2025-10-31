// Track current mode
let currentMode = 'blue'; // 'blue' or 'green'

// Initialize
document.addEventListener('DOMContentLoaded', function() {
    const newGameBtn = document.getElementById('newGame');
    const toggleModeBtn = document.getElementById('toggleMode');

    newGameBtn.addEventListener('click', generateNewGame);
    toggleModeBtn.addEventListener('click', toggleScoringMode);
});

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
    } catch (error) {
        console.error('Error generating new game:', error);
        alert('Failed to generate new game. Please try again.');
    }
}

// Update the goal display with new goals
function updateGoalDisplay(goals) {
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

// Toggle between blue and green scoring modes
function toggleScoringMode() {
    const goalCard = document.getElementById('goalCard');
    const toggleBtn = document.getElementById('toggleMode');
    const blueTracks = goalCard.querySelectorAll('.blue-track');
    const greenTracks = goalCard.querySelectorAll('.green-track');

    if (currentMode === 'blue') {
        // Switch to green
        currentMode = 'green';
        goalCard.classList.remove('blue-side');
        goalCard.classList.add('green-side');
        toggleBtn.textContent = 'Switch to Blue Side';

        blueTracks.forEach(track => track.style.display = 'none');
        greenTracks.forEach(track => track.style.display = 'block');
    } else {
        // Switch to blue
        currentMode = 'blue';
        goalCard.classList.remove('green-side');
        goalCard.classList.add('blue-side');
        toggleBtn.textContent = 'Switch to Green Side';

        blueTracks.forEach(track => track.style.display = 'block');
        greenTracks.forEach(track => track.style.display = 'none');
    }
}

// Optional: Calculate scores (for future enhancement)
async function calculateScores(round, playerCounts) {
    try {
        const response = await fetch('/api/calculate-scores', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                mode: currentMode,
                round: round,
                playerCounts: playerCounts
            })
        });

        if (!response.ok) {
            throw new Error('Failed to calculate scores');
        }

        const scores = await response.json();
        return scores;
    } catch (error) {
        console.error('Error calculating scores:', error);
        return null;
    }
}
