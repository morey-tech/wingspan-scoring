# Score Tracking Features

## Overview
The Wingspan Goals app now includes comprehensive score tracking functionality with colored player cubes, interactive scoring, and automatic calculations.

## Player Management

### Setup
- **Player Count**: Select 2-5 players
- **Player Names**: Editable names (default: "Player 1", "Player 2", etc.)
- **Cube Colors**: Each player assigned unique color
  - Player 1: Blue (#2196F3)
  - Player 2: Purple (#9C27B0)
  - Player 3: Green (#4CAF50)
  - Player 4: Red (#F44336)
  - Player 5: Yellow (#FFC107)

### Player List
- Visual cube color indicator for each player
- Editable name inputs
- Color label showing assigned color name
- Clean, card-style layout

## Interactive Score Tracking

### Score Box Interaction
1. **Click any score box** in the goal card
2. **Player selection menu appears** near the clicked box
3. **Select player** to place their cube in that box
4. **Click again** to move cube to different box
5. **Multiple cubes** can be placed in same box (for ties)

### Player Selection Menu
- Shows all players with their cube colors
- Currently placed players marked with ✓
- Click to toggle placement
- Automatically enforces: one cube per player per round
- Close button to dismiss

### Cube Placement Logic
- **One cube per player per round**: Clicking a new box removes cube from old box in that round
- **Multiple players per box**: Supports ties (2+ players can score the same)
- **Visual feedback**: Cubes pop in with animation
- **Persistent**: Saved to localStorage automatically

## Score Tracker Table

### Display Format
| Player | R1 | R2 | R3 | R4 | Total |
|--------|----|----|----|----|-------|
| 🔵 Alice | 5 | 2 | 6 | 7 | **20** |
| 🟣 Bob | 4 | 5 | 3 | 4 | **16** |
| 🟢 Charlie | 3 | 1 | 2 | 2 | **8** |

### Features
- **Running Totals**: Automatic calculation as cubes are placed
- **Winner Highlighting**: Yellow background for highest score
- **Visual Cube Indicators**: Colored cube next to each name
- **Real-time Updates**: Table updates immediately when cubes are placed
- **Clear Button**: "Clear All Cubes" button to reset scores

### Score Calculation
- Reads cube placements from score boxes
- Extracts score value from the box (0-7 depending on round/mode)
- Calculates total across all 4 rounds
- Highlights winner(s) automatically

## Session Persistence

### localStorage Integration
- **Auto-save**: Game state saved on every change
- **Auto-load**: State restored on page refresh
- **Saved Data**:
  - Player count, names, and colors
  - Cube placements for all rounds
  - Current scoring mode (blue/green)

### Data Structure
```javascript
{
  players: [
    { id: 0, name: "Alice", color: "blue", scores: [5, 2, 6, 7] },
    { id: 1, name: "Bob", color: "purple", scores: [4, 5, 3, 4] }
  ],
  cubePlacements: {
    "1-5": ["blue"],      // Round 1, score 5: Alice
    "1-4": ["purple"],    // Round 1, score 4: Bob
    "2-2": ["blue", "green"]  // Round 2, score 2: Alice & Charlie tied
  }
}
```

## Visual Design

### Cube Styling
- **Size**: 16x16px in score boxes, 24x24px in player list
- **Shape**: Rounded squares (4px radius)
- **Border**: Dark outline for visibility
- **Shadow**: Subtle drop shadow
- **Animation**: Pop-in effect when placed

### Score Boxes
- **Hover Effect**: Scale up slightly, add shadow
- **Cursor**: Pointer to indicate clickability
- **Cube Container**: Centered below score value
- **Spacing**: 3px gap between multiple cubes

### Score Table
- **Header**: Purple gradient (#667eea)
- **Rows**: Alternating hover effect
- **Winner Row**: Yellow highlight (#fff9c4)
- **Total Column**: Light gray background
- **Winner Total**: Bright yellow (#ffeb3b) with orange text

## User Workflow

### Typical Usage
1. **Before Game**:
   - Set up players (names and count)
   - Select expansions
   - Generate new goals

2. **During Game**:
   - Play physical Wingspan game
   - After each round, click score boxes
   - Place player cubes in earned score positions
   - Watch score tracker update automatically

3. **After Game**:
   - View final scores in tracker table
   - Winner highlighted automatically
   - Print card and scores for records
   - Clear cubes for next game

### Score Entry Methods

**Blue Side (Linear)**:
- Each player places cube in box matching their count
- Example: Player scores 3 items → place cube in "3" box
- Multiple players can score the same amount

**Green Side (Competitive)**:
- 1st place player → cube in 1ST PLACE box
- 2nd place player → cube in 2ND PLACE box
- 3rd place player → cube in 3RD PLACE box
- 4th/5th place → cube in "0" box
- Tied players → both cubes in same box

## Technical Implementation

### JavaScript State Management
- Global `gameState` object
- Functions for player management
- Cube placement logic
- Score calculation algorithms
- localStorage sync

### Event Handlers
- Score box click → show player menu
- Player menu item click → toggle placement
- Number of players change → reinitialize
- Clear button click → confirm and reset

### CSS Features
- Cube color classes (blue, purple, green, red, yellow)
- Hover and transition effects
- Responsive layout
- Print media queries
- Animation keyframes

## Accessibility

### Keyboard Support
- Tab navigation through player name inputs
- Enter to confirm name changes
- Escape to close player menu (future enhancement)

### Visual Feedback
- Clear hover states
- Distinct cube colors
- Selected state in menu
- Winner highlighting

## Browser Compatibility
- Modern browsers (Chrome, Firefox, Safari, Edge)
- Requires JavaScript enabled
- localStorage support
- CSS Grid and Flexbox

## Future Enhancements
- Undo/redo cube placements
- Keyboard shortcuts
- Touch/swipe gestures for mobile
- Export scores as CSV/JSON
- Game history tracking
- Statistics across multiple games
- Drag-and-drop cube placement
