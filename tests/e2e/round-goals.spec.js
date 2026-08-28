/**
 * Round goal scoring, for 2, 3 and 4 players.
 *
 * The green side of the goal mat is competitive: players place a cube on the box for the
 * place they finished, and the box's printed value is what they score. Players who tie
 * add together the points for every place they collectively occupy and divide evenly,
 * rounding down. The blue side is flat -- 1 point per item, capped at 5.
 *
 * The browser sends the printed value of the box a cube sits on as that player's "count"
 * and lets the backend re-derive places from the ordering, so these tests are what proves
 * that translation is right at every player count.
 */
const { test, expect } = require('@playwright/test');
const {
    GREEN_BOARD,
    GREEN_POSITIONS,
    openApp,
    setPlayerCount,
    switchToGreenSide,
    placeCube,
    renamePlayer,
    readRoundScoreTable,
    readRoundGoalsColumn,
} = require('./helpers');

const NAMES = ['Ada', 'Grace', 'Katherine', 'Dorothy'];

/** Set up a game with `numPlayers` distinctly-named players on the green side. */
async function startGreenGame(page, numPlayers) {
    await openApp(page);
    await setPlayerCount(page, numPlayers);
    await switchToGreenSide(page);
    for (let seat = 1; seat <= numPlayers; seat++) {
        await renamePlayer(page, seat, NAMES[seat - 1]);
    }
    return NAMES.slice(0, numPlayers);
}

test.describe('green side, every player takes a different place', () => {
    for (const numPlayers of [2, 3, 4]) {
        test(`${numPlayers} players score the printed value of their box`, async ({ page }) => {
            const names = await startGreenGame(page, numPlayers);

            for (let round = 1; round <= 4; round++) {
                for (let seat = 1; seat <= numPlayers; seat++) {
                    await placeCube(page, round, GREEN_POSITIONS[seat - 1], names[seat - 1]);
                }
            }

            const table = await readRoundScoreTable(page, numPlayers);

            for (let seat = 1; seat <= numPlayers; seat++) {
                const name = names[seat - 1];
                const expected = {
                    r1: GREEN_BOARD[1][seat - 1],
                    r2: GREEN_BOARD[2][seat - 1],
                    r3: GREEN_BOARD[3][seat - 1],
                    r4: GREEN_BOARD[4][seat - 1],
                };
                expected.total = expected.r1 + expected.r2 + expected.r3 + expected.r4;
                expect(table[name], `round scores for ${name}`).toEqual(expected);
            }

            // The Game End "Round Goals" column mirrors the Total column above it.
            const roundGoals = await readRoundGoalsColumn(page);
            for (let seat = 1; seat <= numPlayers; seat++) {
                const name = names[seat - 1];
                expect(roundGoals[name], `round goals cell for ${name}`).toBe(
                    String(table[name].total)
                );
            }
        });
    }
});

test.describe('green side ties', () => {
    // Tied players split the sum of the places they occupy, rounded down.
    const scenarios = [
        {
            name: '2 players tie for 1st in round 4',
            numPlayers: 2,
            round: 4,
            // seat -> place box
            places: ['1st', '1st'],
            want: [5, 5], // (7+4)/2
        },
        {
            name: '3 players, two tie for 1st in round 3',
            numPlayers: 3,
            round: 3,
            places: ['1st', '1st', '3rd'],
            want: [4, 4, 2], // (6+3)/2, then 3rd place outright
        },
        {
            name: '3 players, two tie for 2nd in round 4',
            numPlayers: 3,
            round: 4,
            places: ['1st', '2nd', '2nd'],
            want: [7, 3, 3], // (4+2)/2
        },
        {
            name: '3 players all tied in round 3',
            numPlayers: 3,
            round: 3,
            places: ['1st', '1st', '1st'],
            want: [3, 3, 3], // (6+3+2)/3
        },
        {
            name: '4 players, three tie for 2nd in round 4',
            numPlayers: 4,
            round: 4,
            places: ['1st', '2nd', '2nd', '2nd'],
            want: [7, 2, 2, 2], // (4+2+0)/3
        },
        {
            name: '4 players, two tie for 3rd in round 3',
            numPlayers: 4,
            round: 3,
            places: ['1st', '2nd', '3rd', '3rd'],
            want: [6, 3, 1, 1], // (2+0)/2
        },
        {
            name: '4 players all tied in round 4',
            numPlayers: 4,
            round: 4,
            places: ['1st', '1st', '1st', '1st'],
            want: [3, 3, 3, 3], // (7+4+2+0)/4
        },
    ];

    for (const scenario of scenarios) {
        test(scenario.name, async ({ page }) => {
            const names = await startGreenGame(page, scenario.numPlayers);

            for (let seat = 1; seat <= scenario.numPlayers; seat++) {
                await placeCube(page, scenario.round, scenario.places[seat - 1], names[seat - 1]);
            }

            const table = await readRoundScoreTable(page, scenario.numPlayers);
            const key = `r${scenario.round}`;

            for (let seat = 1; seat <= scenario.numPlayers; seat++) {
                const name = names[seat - 1];
                expect(table[name][key], `round ${scenario.round} points for ${name}`).toBe(
                    scenario.want[seat - 1]
                );
            }
        });
    }
});

test.describe('blue side', () => {
    for (const numPlayers of [2, 3, 4]) {
        test(`${numPlayers} players score one point per item`, async ({ page }) => {
            await openApp(page);
            await setPlayerCount(page, numPlayers);
            const names = NAMES.slice(0, numPlayers);
            for (let seat = 1; seat <= numPlayers; seat++) {
                await renamePlayer(page, seat, names[seat - 1]);
            }

            // Each seat takes a different count so nobody's score can be confused
            // with anyone else's: seat 1 -> 5, seat 2 -> 4, and so on.
            const counts = [5, 4, 3, 2];

            for (let round = 1; round <= 4; round++) {
                for (let seat = 1; seat <= numPlayers; seat++) {
                    await placeCube(page, round, String(counts[seat - 1]), names[seat - 1]);
                }
            }

            const table = await readRoundScoreTable(page, numPlayers);

            for (let seat = 1; seat <= numPlayers; seat++) {
                const name = names[seat - 1];
                const perRound = counts[seat - 1];
                expect(table[name], `blue-side scores for ${name}`).toEqual({
                    r1: perRound,
                    r2: perRound,
                    r3: perRound,
                    r4: perRound,
                    total: perRound * 4,
                });
            }
        });
    }
});
