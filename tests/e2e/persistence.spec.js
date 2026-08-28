/**
 * A game in progress is held in localStorage and must survive a page refresh.
 *
 * This is the scenario that used to separate 2-player games from larger ones: the restore
 * was gated on the player-count selector, which starts at the server-rendered default, so
 * any game with a different number of players was silently discarded on load and reset to
 * a fresh 2-player game. Cubes, names, goals and the chosen board side all vanished.
 */
const { test, expect } = require('@playwright/test');
const {
    openApp,
    setPlayerCount,
    switchToGreenSide,
    placeCube,
    renamePlayer,
    readCubePlacements,
    readPlayers,
} = require('./helpers');

// A game in progress lives in localStorage and must survive a page refresh. This is the
// scenario that separates 2-player games (which restore) from 3- and 4-player games.
test.describe('game state survives a page reload', () => {
    for (const numPlayers of [2, 3, 4]) {
        test(`${numPlayers} players`, async ({ page }) => {
            await openApp(page);
            await setPlayerCount(page, numPlayers);
            await switchToGreenSide(page);

            // Give each player a distinctive name so a silent reset is unmistakable.
            const names = ['Ada', 'Grace', 'Katherine', 'Dorothy'].slice(0, numPlayers);
            for (let seat = 1; seat <= numPlayers; seat++) {
                await renamePlayer(page, seat, names[seat - 1]);
            }

            // Round 1: everyone takes a different place.
            const positions = ['1st', '2nd', '3rd', 'other'];
            for (let seat = 1; seat <= numPlayers; seat++) {
                await placeCube(page, 1, positions[seat - 1], names[seat - 1]);
            }

            const placementsBefore = await readCubePlacements(page);
            const playersBefore = await readPlayers(page);
            expect(Object.keys(placementsBefore)).toHaveLength(numPlayers);

            await page.reload();
            await expect(page.locator('#playerList .player-item').first()).toBeVisible();

            // The player-count selector must still report the saved game's size.
            await expect(page.locator('#numPlayers')).toHaveValue(String(numPlayers));
            await expect(page.locator('#playerList .player-item')).toHaveCount(numPlayers);

            expect(await readPlayers(page)).toEqual(playersBefore);
            expect(await readCubePlacements(page)).toEqual(placementsBefore);

            // The green side must still be showing -- the mode is part of the game.
            await expect(page.locator('.green-track').first()).toBeVisible();

            // And the cubes must actually be drawn, not just present in state.
            const renderedCubes = await page.locator('.green-track .placed-cube').count();
            expect(renderedCubes).toBe(numPlayers);
        });
    }
});

test.describe('game end table follows the roster', () => {
    for (const numPlayers of [2, 3, 4]) {
        test(`${numPlayers} players keep their names after a reload`, async ({ page }) => {
            await openApp(page);
            await setPlayerCount(page, numPlayers);

            const names = ['Ada', 'Grace', 'Katherine', 'Dorothy'].slice(0, numPlayers);
            for (let seat = 1; seat <= numPlayers; seat++) {
                await renamePlayer(page, seat, names[seat - 1]);
            }

            await page.reload();
            await expect(page.locator('#gameEndTableBody tr')).toHaveCount(numPlayers);

            const rowNames = await page
                .locator('#gameEndTableBody .player-name-display')
                .allInnerTexts();
            expect(rowNames.map((n) => n.trim())).toEqual(names);
        });
    }
});
