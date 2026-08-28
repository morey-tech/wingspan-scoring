/**
 * Shared actions for the browser tests, so specs read as game actions rather than CSS
 * selectors: set the player count, place a cube, rename a player, calculate the game end.
 *
 * Two things to know before adding to these:
 *
 *   - openApp() deals goals with Oceania disabled and then re-enables it. The Oceania
 *     "No Goal" tile makes its round reject any cube worth more than zero, which would
 *     make the scoring specs fail at random.
 *   - waitForTablesSettled() is the only safe way to wait for a table. See its comment.
 */
const { expect } = require('@playwright/test');

// Printed value of each green goal-mat place box, by round. Index 0 is 1st place.
const GREEN_BOARD = {
    1: [4, 1, 0, 0],
    2: [5, 2, 0, 0],
    3: [6, 3, 2, 0],
    4: [7, 4, 2, 0],
};

// data-position attribute for each green place box, in placing order.
const GREEN_POSITIONS = ['1st', '2nd', '3rd', 'other'];

// Default player names the app assigns, in seat order.
const DEFAULT_NAMES = ['Dani', 'Nick', 'Player 3', 'Player 4', 'Player 5'];

/**
 * Load the app with a clean localStorage so tests never inherit a previous game, and deal
 * a goal set that contains no Oceania "No Goal" tile.
 *
 * The home page deals random goals from every expansion, and the Oceania "No Goal" tile
 * makes its round refuse any cube worth more than zero -- which would make the scoring
 * specs fail at random. Dealing without Oceania guarantees four real goals; re-checking
 * the box afterwards restores the nectar columns without re-dealing.
 */
async function openApp(page) {
    await page.goto('/');
    await page.evaluate(() => window.localStorage.clear());
    await page.reload();
    await expect(page.locator('#playerList .player-item').first()).toBeVisible();

    await page.uncheck('#oceania');
    await page.click('#newGame');
    await expect(page.locator('.goal-info[data-goal-id="oc-no-goal"]')).toHaveCount(0);
    await page.check('#oceania');
}

/** Set the number of players, accepting the "this clears all cubes" confirmation. */
async function setPlayerCount(page, count) {
    page.once('dialog', (dialog) => dialog.accept());
    await page.selectOption('#numPlayers', String(count));
    await expect(page.locator('#playerList .player-item')).toHaveCount(count);
    await waitForGameEndTableInSync(page);
}

/** Switch the goal mat to the green (competitive) side, accepting the confirmation. */
async function switchToGreenSide(page) {
    page.once('dialog', (dialog) => dialog.accept());
    await page.click('#toggleMode');
    await expect(page.locator('.green-track').first()).toBeVisible();
}

/**
 * Place a player's cube on a box in the given round.
 * On the green side pass a position such as '1st'; on the blue side pass the number of
 * items as a string, e.g. '3'.
 */
async function placeCube(page, round, position, playerName) {
    const box = page.locator(`.score-box[data-round="${round}"][data-position="${position}"]`);
    await box.scrollIntoViewIfNeeded();
    await box.click();

    // The menu is position:fixed just below the box, so for boxes low on the page it can
    // render off-screen. Dispatch the click rather than steering a real mouse to it.
    const menu = page.locator('.player-menu');
    await expect(menu).toHaveCount(1);
    await menu.locator('.player-menu-item', { hasText: playerName }).dispatchEvent('click');

    // Selecting a player saves the game and then re-opens the menu with the player
    // ticked. Wait for that second menu before closing, or the handler will resurrect it
    // right after we dismiss it.
    await expect(
        page.locator('.player-menu .player-menu-item.selected', { hasText: playerName })
    ).toHaveCount(1);
    await page.locator('.player-menu .menu-btn-close').dispatchEvent('click');
    await expect(menu).toHaveCount(0);

    // Placing a cube kicks off a burst of /api/calculate-scores requests that redraw both
    // tables; wait for them to settle so the next read sees final values.
    await waitForTablesSettled(page);
}

/**
 * Wait until both score tables reflect the current game state.
 *
 * The app rebuilds each table asynchronously and swaps the finished render in atomically,
 * so while a rebuild is in flight the DOM still shows the *previous* table -- complete,
 * self-consistent, and out of date. Waiting for row counts, network idle, or for the
 * markup to stop changing would all accept that stale table. Instead, assert the two
 * invariants that only hold once the newest render has landed:
 *
 *   1. Every player has a row, and a round shows a score exactly when that player has a
 *      cube in that round.
 *   2. The Game End "Round Goals" cell equals the player's Round End Scoring total.
 */
async function waitForTablesSettled(page) {
    await expect
        .poll(() =>
            page.evaluate(() => {
                const rows = Array.from(document.querySelectorAll('#scoreTableBody tr'));
                if (rows.length !== gameState.players.length) {
                    return false;
                }

                const goalInputs = document.querySelectorAll(
                    '#gameEndTableBody .game-end-input[data-field="roundGoals"]'
                );
                if (goalInputs.length !== gameState.players.length) {
                    return false;
                }

                return gameState.players.every((player, index) => {
                    const cells = rows[index].querySelectorAll('td');
                    if (cells[0].innerText.trim() !== player.name) {
                        return false;
                    }

                    for (let round = 1; round <= 4; round++) {
                        const hasCube = Object.entries(gameState.cubePlacements).some(
                            ([key, colors]) =>
                                parseInt(key.split('-')[0]) === round &&
                                colors.includes(player.color)
                        );
                        const showsScore = cells[round].innerText.trim() !== '-';
                        if (hasCube !== showsScore) {
                            return false;
                        }
                    }

                    return cells[5].innerText.trim() === goalInputs[index].value;
                });
            })
        )
        .toBe(true);
    await page.waitForLoadState('networkidle');
}

/**
 * Wait until the Game End table's rows match the current roster.
 *
 * Use this after anything that changes who is playing -- adding a seat, renaming someone.
 * A rebuild triggered by that change is still in flight, and until it lands the table
 * shows the previous roster, so comparing the rendered names against gameState is what
 * tells the two apart. Chains into waitForTablesSettled() for the scores themselves.
 */
async function waitForGameEndTableInSync(page) {
    await expect
        .poll(() =>
            page.evaluate(() => {
                const rendered = Array.from(
                    document.querySelectorAll('#gameEndTableBody tr')
                ).map((row) => row.querySelector('.player-name-display').textContent.trim());
                const roster = gameState.players.map((p) => p.name);
                return (
                    rendered.length === roster.length &&
                    rendered.every((name, i) => name === roster[i])
                );
            })
        )
        .toBe(true);
    await waitForTablesSettled(page);
}

/**
 * Rename the player in the given seat (1-indexed).
 *
 * The change handler rebuilds the game end table asynchronously, so wait for that rebuild
 * to finish before returning.
 */
async function renamePlayer(page, seat, name) {
    await attemptRename(page, seat, name);
    await expect(
        page.locator('#gameEndTableBody .player-name-display').nth(seat - 1)
    ).toHaveText(name);
}

/**
 * Type a name into a seat without requiring the app to accept it. Use this when the name
 * may be rejected (a duplicate, say); renamePlayer would time out waiting for a rename
 * that never happens.
 */
async function attemptRename(page, seat, name) {
    const input = page.locator('.player-name-input').nth(seat - 1);
    await input.fill(name);
    await input.blur();
    await waitForGameEndTableInSync(page);
}

/** Read the Round End Scoring table as { playerName: { r1, r2, r3, r4, total } }. */
async function readRoundScoreTable(page, expectedRows) {
    await waitForTablesSettled(page);
    if (expectedRows !== undefined) {
        await expect(page.locator('#scoreTableBody tr')).toHaveCount(expectedRows);
    }
    return page.evaluate(() => {
        const result = {};
        document.querySelectorAll('#scoreTableBody tr').forEach((row) => {
            const cells = row.querySelectorAll('td');
            const name = cells[0].innerText.trim();
            const num = (cell) => (cell.innerText.trim() === '-' ? null : Number(cell.innerText.trim()));
            result[name] = {
                r1: num(cells[1]),
                r2: num(cells[2]),
                r3: num(cells[3]),
                r4: num(cells[4]),
                total: num(cells[5]),
            };
        });
        return result;
    });
}

/** Read the Game End table's Round Goals column as { playerName: value }. */
async function readRoundGoalsColumn(page) {
    await waitForTablesSettled(page);
    return page.evaluate(() => {
        const result = {};
        document.querySelectorAll('#gameEndTableBody tr').forEach((row) => {
            const name = row.querySelector('.player-name-display').textContent.trim();
            const input = row.querySelector('.game-end-input[data-field="roundGoals"]');
            result[name] = input.value;
        });
        return result;
    });
}

/** Fill one player's game-end inputs. Only the supplied fields are written. */
async function fillGameEndRow(page, seat, fields) {
    for (const [field, value] of Object.entries(fields)) {
        const input = page.locator(`.game-end-input[data-player="${seat}"][data-field="${field}"]`);
        await input.fill(String(value));
        await input.blur();
    }
}

/** Click Calculate and return the final rankings in displayed order. */
async function calculateGameEnd(page) {
    await page.click('#calculateBtn');
    await expect(page.locator('#resultsSection')).toBeVisible();
    await expect(page.locator('#rankingsList .ranking-item').first()).toBeVisible();

    return page.evaluate(() =>
        Array.from(document.querySelectorAll('#rankingsList .ranking-item')).map((item) => ({
            rank: item.querySelector('.rank-label').textContent.trim(),
            name: item.querySelector('.player-name').textContent.trim(),
            score: Number(item.querySelector('.player-score').textContent.trim().match(/\d+/)[0]),
        }))
    );
}

// gameState is declared with `let` at the top level of a classic script, so it lives in
// the global declarative scope rather than on `window` -- reference it unqualified.
/** Read the whole cube layout as { 'round-score-position': [colors] }. */
async function readCubePlacements(page) {
    return page.evaluate(() => gameState.cubePlacements);
}

/** Read the player roster as [{ name, color }] in seat order. */
async function readPlayers(page) {
    return page.evaluate(() => gameState.players.map((p) => ({ name: p.name, color: p.color })));
}

module.exports = {
    GREEN_BOARD,
    GREEN_POSITIONS,
    DEFAULT_NAMES,
    openApp,
    setPlayerCount,
    switchToGreenSide,
    placeCube,
    renamePlayer,
    attemptRename,
    waitForGameEndTableInSync,
    waitForTablesSettled,
    readRoundScoreTable,
    readRoundGoalsColumn,
    fillGameEndRow,
    calculateGameEnd,
    readCubePlacements,
    readPlayers,
};
