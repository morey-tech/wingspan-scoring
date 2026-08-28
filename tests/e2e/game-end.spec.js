const { test, expect } = require('@playwright/test');
const {
    GREEN_POSITIONS,
    openApp,
    setPlayerCount,
    switchToGreenSide,
    placeCube,
    renamePlayer,
    attemptRename,
    readRoundGoalsColumn,
    fillGameEndRow,
    calculateGameEnd,
} = require('./helpers');

const NAMES = ['Ada', 'Grace', 'Katherine', 'Dorothy'];

/** Fetch the most recently saved game straight from the API. */
async function fetchLatestGame(page) {
    const response = await page.request.get('/api/games?limit=1');
    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.games.length).toBeGreaterThan(0);
    return body.games[0];
}

test.describe('round goals stay in the game end table', () => {
    for (const numPlayers of [2, 3, 4]) {
        test(`${numPlayers} players keep their round goals when the table is rebuilt`, async ({
            page,
        }) => {
            await openApp(page);
            await setPlayerCount(page, numPlayers);
            await switchToGreenSide(page);
            const names = NAMES.slice(0, numPlayers);
            for (let seat = 1; seat <= numPlayers; seat++) {
                await renamePlayer(page, seat, names[seat - 1]);
            }

            // Round 4 only: 1st=7, 2nd=4, 3rd=2, 4th=0.
            for (let seat = 1; seat <= numPlayers; seat++) {
                await placeCube(page, 4, GREEN_POSITIONS[seat - 1], names[seat - 1]);
            }

            const expected = [7, 4, 2, 0].slice(0, numPlayers);
            const before = await readRoundGoalsColumn(page);
            for (let seat = 1; seat <= numPlayers; seat++) {
                expect(before[names[seat - 1]]).toBe(String(expected[seat - 1]));
            }

            // Renaming a player rebuilds every game-end row. The round goals are
            // derived from the cubes, so they must come back unchanged.
            await renamePlayer(page, 1, 'Ada Lovelace');
            const renamed = ['Ada Lovelace', ...names.slice(1)];

            const after = await readRoundGoalsColumn(page);
            for (let seat = 1; seat <= numPlayers; seat++) {
                expect(
                    after[renamed[seat - 1]],
                    `round goals for ${renamed[seat - 1]} after the table was rebuilt`
                ).toBe(String(expected[seat - 1]));
            }
        });
    }
});

test.describe('final scoring', () => {
    const scenarios = [
        {
            numPlayers: 2,
            // Seat scores chosen so the totals are unambiguous.
            rows: [
                { birdPoints: 40, bonusCards: 10, eggs: 12, cachedFood: 3, tuckedCards: 5, nectarForest: 6, nectarGrassland: 2, nectarWetland: 4, unusedFood: 2 },
                { birdPoints: 35, bonusCards: 8, eggs: 14, cachedFood: 2, tuckedCards: 4, nectarForest: 3, nectarGrassland: 7, nectarWetland: 1, unusedFood: 5 },
            ],
            // Round goals from the green board, round 4 only: 7 and 4.
            roundGoals: [7, 4],
            // nectar: forest 6>3 -> 5/2 ; grassland 2<7 -> 2/5 ; wetland 4>1 -> 5/2
            nectar: [[5, 2, 5], [2, 5, 2]],
        },
        {
            numPlayers: 3,
            rows: [
                { birdPoints: 44, bonusCards: 12, eggs: 10, cachedFood: 1, tuckedCards: 6, nectarForest: 9, nectarGrassland: 1, nectarWetland: 4, unusedFood: 2 },
                { birdPoints: 38, bonusCards: 9, eggs: 16, cachedFood: 4, tuckedCards: 2, nectarForest: 5, nectarGrassland: 7, nectarWetland: 2, unusedFood: 4 },
                { birdPoints: 30, bonusCards: 6, eggs: 11, cachedFood: 0, tuckedCards: 3, nectarForest: 1, nectarGrassland: 3, nectarWetland: 8, unusedFood: 1 },
            ],
            roundGoals: [7, 4, 2],
            // forest 9>5>1 -> 5/2/0 ; grassland 7>3>1 -> seat2 5, seat3 2, seat1 0
            // wetland 8>4>2 -> seat3 5, seat1 2, seat2 0
            nectar: [[5, 0, 2], [2, 5, 0], [0, 2, 5]],
        },
        {
            numPlayers: 4,
            rows: [
                { birdPoints: 48, bonusCards: 14, eggs: 13, cachedFood: 2, tuckedCards: 7, nectarForest: 10, nectarGrassland: 2, nectarWetland: 5, unusedFood: 3 },
                { birdPoints: 40, bonusCards: 11, eggs: 15, cachedFood: 3, tuckedCards: 5, nectarForest: 6, nectarGrassland: 8, nectarWetland: 3, unusedFood: 5 },
                { birdPoints: 33, bonusCards: 8, eggs: 12, cachedFood: 1, tuckedCards: 4, nectarForest: 2, nectarGrassland: 4, nectarWetland: 9, unusedFood: 1 },
                { birdPoints: 26, bonusCards: 5, eggs: 9, cachedFood: 0, tuckedCards: 2, nectarForest: 0, nectarGrassland: 0, nectarWetland: 0, unusedFood: 0 },
            ],
            roundGoals: [7, 4, 2, 0],
            // forest 10>6>2, grassland 8>4>2, wetland 9>5>3; seat 4 has none anywhere.
            nectar: [[5, 0, 2], [2, 5, 0], [0, 2, 5], [0, 0, 0]],
        },
    ];

    for (const scenario of scenarios) {
        const { numPlayers } = scenario;

        test(`${numPlayers} players: totals, ranking and nectar`, async ({ page }) => {
            await openApp(page);
            await setPlayerCount(page, numPlayers);
            await switchToGreenSide(page);
            const names = NAMES.slice(0, numPlayers);
            for (let seat = 1; seat <= numPlayers; seat++) {
                await renamePlayer(page, seat, names[seat - 1]);
            }

            // One scored round is enough to give every player a round-goal total.
            for (let seat = 1; seat <= numPlayers; seat++) {
                await placeCube(page, 4, GREEN_POSITIONS[seat - 1], names[seat - 1]);
            }

            for (let seat = 1; seat <= numPlayers; seat++) {
                await fillGameEndRow(page, seat, scenario.rows[seat - 1]);
            }

            const rankings = await calculateGameEnd(page);
            expect(rankings).toHaveLength(numPlayers);

            const wantTotals = scenario.rows.map((row, i) => {
                const base =
                    row.birdPoints +
                    row.bonusCards +
                    row.eggs +
                    row.cachedFood +
                    row.tuckedCards +
                    scenario.roundGoals[i];
                const nectar = scenario.nectar[i].reduce((a, b) => a + b, 0);
                return base + nectar;
            });

            // Seats were built in descending score order, so the ranking should match.
            expect(rankings.map((r) => r.name)).toEqual(names);
            expect(rankings.map((r) => r.score)).toEqual(wantTotals);
            expect(rankings[0].rank).toContain('Winner');

            // The totals column in the table must agree with the rankings.
            const tableTotals = await page.evaluate(() =>
                Array.from(document.querySelectorAll('#gameEndTableBody tr')).map((row) => ({
                    name: row.querySelector('.player-name-display').textContent.trim(),
                    total: Number(row.querySelector('.total-display').textContent.trim()),
                }))
            );
            for (let seat = 1; seat <= numPlayers; seat++) {
                const row = tableTotals.find((r) => r.name === names[seat - 1]);
                expect(row.total, `table total for ${names[seat - 1]}`).toBe(
                    wantTotals[seat - 1]
                );
            }

            // The saved game must carry the same totals and the per-round breakdown
            // that the Game History page renders.
            const saved = await fetchLatestGame(page);
            expect(saved.numPlayers).toBe(numPlayers);
            expect(saved.winnerName).toBe(names[0]);
            expect(saved.winnerScore).toBe(wantTotals[0]);

            expect(saved.roundBreakdown, 'per-round goal breakdown was not saved').toBeTruthy();
            for (let seat = 1; seat <= numPlayers; seat++) {
                const breakdown = saved.roundBreakdown[names[seat - 1]];
                expect(breakdown, `no saved breakdown for ${names[seat - 1]}`).toBeTruthy();
                expect(breakdown.round4).toBe(scenario.roundGoals[seat - 1]);
            }
        });
    }
});

test('unused food breaks a tie on total score', async ({ page }) => {
    await openApp(page);
    await setPlayerCount(page, 3);
    const names = NAMES.slice(0, 3);
    for (let seat = 1; seat <= 3; seat++) {
        await renamePlayer(page, seat, names[seat - 1]);
    }
    await page.uncheck('#oceania');

    // Seats 1 and 2 tie on 60; seat 2 holds more unused food and must win.
    await fillGameEndRow(page, 1, { birdPoints: 60, unusedFood: 1 });
    await fillGameEndRow(page, 2, { birdPoints: 60, unusedFood: 8 });
    await fillGameEndRow(page, 3, { birdPoints: 50, unusedFood: 9 });

    const rankings = await calculateGameEnd(page);

    expect(rankings[0].name).toBe(names[1]);
    expect(rankings[0].rank).toContain('Winner');
    expect(rankings[1].name).toBe(names[0]);
    expect(rankings[2].name).toBe(names[2]);
});

test('two players sharing a name are not merged into one competitor', async ({ page }) => {
    await openApp(page);
    await setPlayerCount(page, 3);
    await switchToGreenSide(page);

    await renamePlayer(page, 1, 'Alex');
    await attemptRename(page, 2, 'Alex');
    await renamePlayer(page, 3, 'Katherine');

    const players = await page.evaluate(() => gameState.players.map((p) => p.name));

    // Either the app keeps the seats distinguishable, or it rejects the duplicate.
    // What it must not do is quietly collapse two players into a single scoring entity.
    expect(new Set(players).size, `duplicate player names: ${players.join(', ')}`).toBe(3);
});

test('a finished 4-player game shows up on the history page', async ({ page }) => {
    await openApp(page);
    await setPlayerCount(page, 4);
    await switchToGreenSide(page);
    const names = NAMES.slice(0, 4);
    for (let seat = 1; seat <= 4; seat++) {
        await renamePlayer(page, seat, names[seat - 1]);
    }

    // Round 4 places: 7 / 4 / 2 / 0.
    for (let seat = 1; seat <= 4; seat++) {
        await placeCube(page, 4, GREEN_POSITIONS[seat - 1], names[seat - 1]);
    }

    await page.uncheck('#oceania');
    const birdPoints = [80, 70, 60, 50];
    for (let seat = 1; seat <= 4; seat++) {
        await fillGameEndRow(page, seat, { birdPoints: birdPoints[seat - 1] });
    }

    const rankings = await calculateGameEnd(page);
    const winnerScore = rankings[0].score;
    expect(winnerScore).toBe(birdPoints[0] + 7);

    await page.goto('/history');
    const card = page.locator('.game-card').first();
    await expect(card).toBeVisible();
    await expect(card.locator('.game-badge').first()).toHaveText('4 Players');
    await expect(card.locator('.winner-name')).toHaveText(names[0]);
    await expect(card.locator('.winner-score')).toHaveText(`${winnerScore} points`);

    // Every player must be listed with the total they finished on.
    await card.locator('.btn-view-details').click();
    const rows = page.locator('.details-scores .details-table tbody tr');
    await expect(rows).toHaveCount(4);
    for (let seat = 1; seat <= 4; seat++) {
        const row = rows.nth(seat - 1);
        await expect(row).toContainText(names[seat - 1]);
        await expect(row.locator('.total-col')).toHaveText(String(rankings[seat - 1].score));
    }

    // The per-round goal breakdown is only rendered when it was saved with the game.
    const breakdown = page.locator('.details-round-breakdown');
    await expect(breakdown, 'round goals breakdown missing from the saved game').toBeVisible();

    const breakdownRows = breakdown.locator('tbody tr');
    await expect(breakdownRows).toHaveCount(4);
    const round4Points = [7, 4, 2, 0];
    for (let seat = 1; seat <= 4; seat++) {
        const cells = breakdownRows.nth(seat - 1).locator('td');
        await expect(cells.nth(0)).toContainText(names[seat - 1]);
        // Columns are player, round 1-4, total. Only round 4 was played.
        await expect(cells.nth(4)).toHaveText(String(round4Points[seat - 1]));
        await expect(cells.nth(5)).toHaveText(String(round4Points[seat - 1]));
    }
});
