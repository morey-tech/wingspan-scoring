# Wingspan Goal Tiles SVG Sprite Sheet

## Overview
This directory contains the complete SVG sprite sheet extracted from the official Wingspan RulePop reference site, including all 46 goal tile definitions and supporting icon symbols.

## File
- **wingspan-sprites.svg** (62 KB) - Complete SVG sprite sheet with all symbols
  - ✓ **XML Validated** - Confirmed valid XML structure
  - ✓ **All References Valid** - 254 internal references verified
  - ✓ **All External Files Present** - 23 image dependencies available

## Usage

### In HTML
To use a goal tile in your HTML:

```html
<svg viewBox="0 0 500 500" width="100" height="100">
  <use href="/static/images/svg/wingspan-sprites.svg#g-birds-in-forest"/>
</svg>
```

### Styling
The SVG uses CSS custom properties (variables) for colors:
- `--paper` - Background color
- `--corner-color` - Corner triangle color
- `--dark-brown` - Text color
- `--taupe`, `--fish-blue`, `--seed-orange`, etc. - Food/resource colors

Define these in your CSS to match your app's theme.

## Goal Tiles Inventory

### Standard Goal Tiles (34 tiles)

#### Habitat Goals (3)
| ID | Description |
|----|-------------|
| `g-birds-in-forest` | Birds in forest habitat |
| `g-birds-in-grassland` | Birds in grassland habitat |
| `g-birds-in-wetland` | Birds in wetland habitat |

#### Egg Location Goals (7)
| ID | Description |
|----|-------------|
| `g-eggs-in-forest` | Eggs in forest habitat |
| `g-eggs-in-grassland` | Eggs in grassland habitat |
| `g-eggs-in-wetland` | Eggs in wetland habitat |
| `g-eggs-in-bowl` | Eggs in bowl nests |
| `g-eggs-in-cavity` | Eggs in cavity nests |
| `g-eggs-in-ground` | Eggs in ground nests |
| `g-eggs-in-platform` | Eggs in platform nests |

#### Nest Type with Eggs Goals (4)
| ID | Description |
|----|-------------|
| `g-bowl-with-egg` | Bowl nest birds with eggs |
| `g-cavity-with-egg` | Cavity nest birds with eggs |
| `g-ground-with-egg` | Ground nest birds with eggs |
| `g-platform-with-egg` | Platform nest birds with eggs |

#### Set Collection Goals (1)
| ID | Description |
|----|-------------|
| `g-egg-habitat-sets` | Sets of eggs across all three habitats |

#### Bird Count Goals (2)
| ID | Description |
|----|-------------|
| `g-total-birds` | Total birds played |
| `g-birds-in-one-row` | Birds in one habitat row |

#### Food Cost Goals (3)
| ID | Description |
|----|-------------|
| `g-inv-in-cost` | Invertebrates in food costs |
| `g-rodent-fish-in-cost` | Rodents + fish in food costs |
| `g-fruit-seed-in-cost` | Fruit + seeds in food costs |

#### Point Value Goals (2)
| ID | Description |
|----|-------------|
| `g-birds-3pt-or-under` | Birds worth ≤3 points |
| `g-birds-over-4pt` | Birds worth >4 points |

#### Bird Power Goals (2)
| ID | Description |
|----|-------------|
| `g-brown-powers` | Birds with brown "when activated" powers |
| `g-white-no-powers` | Birds with white or no powers |

#### Bird Features Goals (4)
| ID | Description |
|----|-------------|
| `g-beak-lt` | Birds with beaks pointing left |
| `g-beak-rt` | Birds with beaks pointing right |
| `g-birds-with-tucked-cards` | Birds with tucked cards |
| `g-eggless-birds` | Birds with no eggs |

#### Layout Goals (2)
| ID | Description |
|----|-------------|
| `g-filled-columns` | Completely filled columns |
| `g-cubes-on-play-bird` | Action cubes on "play a bird" row |

#### Resource Goals (2)
| ID | Description |
|----|-------------|
| `g-food-owned` | Food tokens in personal supply |
| `g-birds-in-hand` | Bird cards in hand at end of round |

#### Other Goals (2)
| ID | Description |
|----|-------------|
| `g-birds-food-cost` | Total food cost of played birds |
| `g-no-goal` | No goal (Oceania expansion) |

### Duet Mode Goal Tiles (12 tiles)

#### Duet Habitat Goals (3)
| ID | Description |
|----|-------------|
| `g-duets-in-forest` | Tokens in forest spaces |
| `g-duets-in-grassland` | Tokens in grassland spaces |
| `g-duets-in-wetland` | Tokens in wetland spaces |

#### Duet Placement Goals (5)
| ID | Description |
|----|-------------|
| `g-duets-in-a-row` | Tokens in any one horizontal row |
| `g-rows-with-duets` | Horizontal rows with at least one token |
| `g-duets-in-interior` | Tokens not on edge of map |
| `g-duets-on-edge` | Tokens on edge of map |
| `g-duets-on-pairs` | Tokens on pairs of matching symbols |

#### Duet Symbol Goals (2)
| ID | Description |
|----|-------------|
| `g-duets-on-nests` | Tokens on nest symbols |
| `g-duets-on-food` | Tokens on food symbols |

#### Duet Count Goals (2)
| ID | Description |
|----|-------------|
| `g-total-duets` | Total tokens on map |
| `g-fewest-bonus-duets` | Fewest tokens on bonus spaces |

## Supporting Symbols

The sprite sheet also includes many supporting icon symbols used within goal tiles:

### Icons
- `i-bird` - Bird icon
- `i-egg` - Egg icon
- `i-bird-card` - Bird card icon
- `i-forest`, `i-grassland`, `i-wetland` - Habitat icons
- `i-nest-bowl`, `i-nest-cavity`, `i-nest-ground`, `i-nest-platform` - Nest type icons
- `i-fish`, `i-fruit`, `i-rodent`, `i-seed`, `i-invertebrate`, `i-nectar` - Food icons
- `i-beak-lt`, `i-beak-rt` - Beak direction icons
- `action-cube` - Action cube icon
- `goal-bkg` - Goal tile background with corners

### Food Resources
- `fish` - Fish detailed illustration
- `rat` - Rodent detailed illustration
- `wheat` - Seed detailed illustration
- `berries` - Fruit detailed illustration
- `slug` - Invertebrate detailed illustration
- `flower` - Nectar detailed illustration

### Other Elements
- `pt-3`, `pt-4` - Point value numbers
- `pt-feather` - Point feather icon
- `pt-lteq`, `pt-gt` - Comparison symbols (≤, >)
- Various other game icons and symbols

## Technical Details

### File Structure
```xml
<svg id="sprites" role="presentation">
  <defs>
    <!-- Basic shapes and backgrounds -->
    <rect id="habitat-bkg" .../>
    <circle id="disk" .../>
    <g id="goal-bkg">...</g>

    <!-- Food resource illustrations -->
    <g id="fish">...</g>
    <path id="rat" .../>
    <!-- ... more resources ... -->

    <!-- Goal tile symbols -->
    <symbol id="g-birds-in-forest" class="goal-tile" viewBox="0 0 500 500">
      <use href="#goal-bkg"/>
      <use href="#i-bird" transform="..."/>
      <text transform="...">in</text>
      <use href="#i-forest" transform="..."/>
    </symbol>
    <!-- ... 45 more goal tiles ... -->
  </defs>
</svg>
```

### ViewBox
All goal tiles use `viewBox="0 0 500 500"` (500x500 coordinate system)

### Transforms
- Icons within tiles are positioned and scaled using `transform` attributes
- Text elements use `transform="translate(x y)"` for positioning

## Source
Extracted from: https://wingspan.rulepop.com/
Date: 2025-11-07

## License
These SVG graphics are copyrighted materials from Wingspan by Stonemaier Games.
Use within this scoring application should comply with fair use guidelines.

## Integration Example

### Display All Goal Tiles
```html
<!DOCTYPE html>
<html>
<head>
<style>
  :root {
    --paper: #f5f1e8;
    --corner-color: #8b7355;
    --dark-brown: #3d2817;
    --taupe: #a89080;
    --fish-blue: #4a90a4;
    --seed-orange: #d97644;
    --fruit-red: #c44569;
    --invertebrate-green: #6b9d5f;
    --nectar-pink: #e88bb4;
  }
  .goal-tile {
    width: 100px;
    height: 100px;
    margin: 10px;
    display: inline-block;
  }
</style>
</head>
<body>
  <div class="goal-tile">
    <svg viewBox="0 0 500 500" width="100" height="100">
      <use href="/static/images/svg/wingspan-sprites.svg#g-birds-in-forest"/>
    </svg>
  </div>
  <!-- Repeat for other tiles -->
</body>
</html>
```

### JavaScript Access
```javascript
// Get list of all goal tile IDs
const goalTiles = [
  'g-birds-in-forest',
  'g-birds-in-grassland',
  'g-birds-in-wetland',
  // ... etc
];

// Dynamically create goal tile elements
goalTiles.forEach(tileId => {
  const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
  svg.setAttribute('viewBox', '0 0 500 500');
  svg.setAttribute('width', '100');
  svg.setAttribute('height', '100');

  const use = document.createElementNS('http://www.w3.org/2000/svg', 'use');
  use.setAttributeNS('http://www.w3.org/1999/xlink', 'href',
    `/static/images/svg/wingspan-sprites.svg#${tileId}`);

  svg.appendChild(use);
  document.body.appendChild(svg);
});
```

## Notes
- The sprite sheet is self-contained with all dependencies
- No external image references except for a few webp files (duet tokens, brush patterns)
- All goal tiles reference the common `goal-bkg` background
- CSS variables allow easy theming without modifying SVG
