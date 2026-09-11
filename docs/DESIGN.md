# Gallery Wall Layout — Design Document

## Overview

This project lays out a physical gallery wall. Frames vary in size and shape. One frame is the centerpiece. The centerpiece is the visual and geometric anchor. All other frames attach to form one connected cluster.

The core is a Go library (`layout` package). It is published as a standalone module. A TypeScript and Vite web app runs the same library in the browser via WASM.

Current architecture: [ARCHITECTURE.md](ARCHITECTURE.md). Implementation walkthrough: [layout/ALGORITHM.md](../layout/ALGORITHM.md).

## Scope

### In scope (v1)

- Pure layout algorithm: positions for N frames from specs and gap
- Validation of structural input (IDs, centerpiece count, positive dimensions)
- Shape-aware collision detection
- Connected blob placement centered on the centerpiece
- Optional wall bounding box (`wallWidth` and `wallHeight`, centered on anchor)
- Output for humans who hang a real wall
- Browser WASM bridge and static web UI (implemented in this repo)

### Out of scope (v1)

- HTTP API server
- Unit conversion (inches, cm, or abstract)
- Rotation (all frames are axis-aligned)
- Drag-and-drop editing on the preview canvas
- Image rendering inside frames

## Input

### Item

Each frame has these fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | yes | Stable identifier, unique in the input set |
| `height` | int | yes | Bounding box height (> 0) |
| `width` | int | yes | Bounding box width (> 0) |
| `shape` | enum | yes | `square`, `rectangle`, `circle`, `ellipse` |
| `centerpiece` | bool | yes | Exactly one item must be `true` |

Dimensions are the source of truth. The `shape` field controls collision geometry and rendering. The library does not validate shape against dimensions. The UI may warn about mismatches (for example `square` with 10×12). The library accepts any positive height and width.

### Layout parameters

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `gap` | int | yes | — | Minimum separation between frames |
| `wallWidth` | float | no | 0 (none) | Wall width. Both wall fields are required if either is set. |
| `wallHeight` | float | no | 0 (none) | Wall height. The cluster must fit inside. |

Wall bounds are centered on the anchor `(0, 0)`. Omit or zero both dimensions for unbounded placement.

All numeric values share the same implicit unit (inches, cm, or abstract). The library performs no unit conversion.

### Example input

```json
{
  "gap": 2,
  "items": [
    { "id": "main",   "height": 36, "width": 24, "shape": "rectangle", "centerpiece": true  },
    { "id": "left",   "height": 14, "width": 11, "shape": "rectangle", "centerpiece": false },
    { "id": "right",  "height": 14, "width": 11, "shape": "rectangle", "centerpiece": false },
    { "id": "top",    "height": 10, "width": 10, "shape": "circle",     "centerpiece": false },
    { "id": "bottom", "height": 12, "width": 16, "shape": "ellipse",    "centerpiece": false }
  ]
}
```

## Output

All positions are relative to the centerpiece anchor. The geometric center of the centerpiece frame is `(0, 0)`. Negative coordinates mean left of or above the centerpiece center. Positive means right of or below.

### Per-item output

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Same as input |
| `centerX` | float | Horizontal offset of frame center from centerpiece center |
| `centerY` | float | Vertical offset of frame center from centerpiece center |
| `x` | float | Top-left X of bounding box (derived: `centerX - width/2`) |
| `y` | float | Top-left Y of bounding box (derived: `centerY - height/2`) |
| `width` | int | Echoed from input |
| `height` | int | Echoed from input |
| `shape` | enum | Echoed from input |
| `offsetFromAnchor` | float | Distance from frame center to centerpiece center |
| `direction` | string | Compass label from centerpiece to frame (`N`, `NE`, `E`, …) |
| `adjacentIds` | string[] | IDs of frames adjacent within `gap` |

Additional hanging helpers may be added if they derive from the core positions.

### Layout metadata

| Field | Type | Description |
|-------|--------|-------------|
| `bounds` | object | Axis-aligned bounding box of the full cluster: `{ minX, minY, maxX, maxY }` in centerpiece-relative coordinates |

### Example output

```json
{
  "anchor": { "itemId": "main", "centerX": 0, "centerY": 0 },
  "items": [
    {
      "id": "main",
      "centerX": 0,
      "centerY": 0,
      "x": -12,
      "y": -18,
      "width": 24,
      "height": 36,
      "shape": "rectangle",
      "offsetFromAnchor": 0,
      "direction": "C"
    },
    {
      "id": "left",
      "centerX": -20.5,
      "centerY": 0,
      "x": -26,
      "y": -7,
      "width": 11,
      "height": 14,
      "shape": "rectangle",
      "offsetFromAnchor": 20.5,
      "direction": "W"
    }
  ],
  "bounds": { "minX": -26, "minY": -18, "maxX": 26, "maxY": 22 }
}
```

## Coordinate system

### Anchor

The centerpiece center is the origin `(0, 0)`. +X is right. +Y is down (canvas and screen convention, matches typical frontend rendering).

### Internal vs output

During placement, the algorithm works in centerpiece-relative space. Negative coordinates are valid (frames above or left of the centerpiece). No global shift forces all coordinates positive. That would hide measure-from-centerpiece semantics.

If internal math uses a temporary translation, discard it before output. Output is always centerpiece-relative.

### Hanging interpretation

Installer workflow:

1. Mark the centerpiece center on the wall (for example at eye level, centered on the wall).
2. For each other frame, measure horizontally and vertically from that mark to the frame center (`centerX`, `centerY`).
3. Use `width` and `height` to align the frame around that center point.

## Shape geometry

Dimensions define an axis-aligned bounding box (AABB). Shape determines the collision footprint within that box:

| Shape | Collision geometry |
|-------|-------------------|
| `rectangle` | Solid axis-aligned rectangle h×w |
| `square` | Same as rectangle (label is cosmetic) |
| `circle` | Circle inscribed in bbox: diameter = `min(height, width)`, centered in bbox |
| `ellipse` | Ellipse inscribed in bbox: semi-axes `height/2`, `width/2` |

### Collision with gap

Two frames collide if their shape-aware footprints overlap after each is expanded outward by `gap/2`. For v1:

- Rectangles: AABB separation test with gap padding on all sides.
- Circles: center distance < `r1 + r2 + gap`.
- Ellipses: conservative test using inflated AABB (ellipse bbox + gap). Exact ellipse-ellipse overlap may be added later if inflated AABB is too loose.

### Adjacency (connectivity)

Two frames are adjacent if they do not overlap and the gap between footprints is ≤ `gap`. Connectivity uses shape-aware geometry, not raw bbox overlap alone.

## Layout algorithm

### High-level flow

```
Layout(items, gap) → LayoutResult | error

1. Validate input
2. Identify and anchor centerpiece at (0, 0)
3. Sort remaining items for placement order
4. For each remaining item:
     a. Generate candidate attachment positions on cluster perimeter
     b. Filter to collision-free candidates
     c. Score survivors
     d. Place at best-scoring position
     e. If no valid candidate → error
5. Compute output (positions, bounds, helpers)
6. Return result
```

### Step 1 — Validate

| Rule | Error |
|------|-------|
| `len(items) == 0` | `empty input` |
| Duplicate IDs | `duplicate id: {id}` |
| Any `height <= 0` or `width <= 0` | `invalid dimensions: {id}` |
| Zero centerpieces | `no centerpiece specified` |
| Multiple centerpieces | `multiple centerpieces: {ids}` |
| Only one of `wallWidth` / `wallHeight` set | `invalid_wall` |
| Centerpiece footprint exceeds wall | `centerpiece_exceeds_wall` |

The library does not validate shape against dimensions.

### Step 2 — Anchor centerpiece

Place the centerpiece at `centerX = 0, centerY = 0`. Initialize `placed` with the centerpiece. Initialize cluster geometry from the centerpiece footprint.

### Step 3 — Placement order

Sort remaining items before placement:

1. Area descending (`height × width`). Large frames establish the blob skeleton.
2. Max side descending (`max(height, width)`). Tiebreaker.
3. ID ascending. Deterministic tiebreaker for tests.

The centerpiece is never reordered. It is already placed.

### Step 4 — Candidate generation

For each unplaced item `P`, generate candidate positions by attaching to the perimeter of the current cluster.

Perimeter attachment means `P` is adjacent (within `gap`) to at least one placed item and does not overlap any placed item.

#### 4a — Find attachment edges

For each placed item `A`, enumerate candidate positions for `P` along each side of `A`:

- Left of A: `P.centerX = A.centerX - (A.width/2 + gap + P.width/2)` (adjusted for shape geometry)
- Right of A: symmetric
- Above A: `P.centerY = A.centerY - (A.height/2 + gap + P.height/2)`
- Below A: symmetric

For circles and ellipses, use shape-aware footprint radius or extents instead of raw bbox half-sides.

Also generate corner attachments (diagonal of two sides) for L-shaped growth and smoother silhouettes.

Slide along each attachment edge in steps of 1 unit to produce candidate center positions.

#### 4b — Filter

Remove any candidate where:

- `P` overlaps any placed item (shape-aware collision with gap), or
- `P` is not adjacent to at least one placed item (floater)

#### 4c — Score

Among surviving candidates, compute a score. Lower is better. Weighted sum of:

| Component | Weight | Description |
|-----------|--------|-------------|
| Compactness | w1 | Minimize average distance from `P.center` to centerpiece center |
| Balance | w2 | Reduce asymmetry of placed area across quadrants relative to centerpiece |
| Blob smoothness | w3 | Penalize long thin protrusions; target cluster aspect ~1.3–1.6 (wider than tall) |
| Concavity penalty | w4 | Penalize deep inward bays (convex hull area vs cluster bbox area) |
| Local continuity | w5 | Reward candidates adjacent to 2+ placed items |

Initial weights (tuned during testing):

```
w1 = 1.0   compactness
w2 = 0.8   balance
w3 = 0.6   blob smoothness
w4 = 0.5   concavity
w5 = 0.4   local continuity (reward → subtract from score)
```

Exact weights are implementation constants. Adjust via tests against golden fixtures.

#### 4d — Place

Select the lowest-scoring candidate. Add `P` to `placed`. Update cluster geometry.

#### 4e — Failure

If no candidate survives filtering for item `P`:

```
error: cannot place item "{id}" — no valid connected position
```

There is no fallback. The library does not drop items or place disconnected frames.

### Step 5 — Output

For each placed item, emit `centerX`, `centerY`, derived `x`, `y`, `offsetFromAnchor`, `direction`, and echo `width`, `height`, `shape`, `id`. Compute cluster `bounds` from the union of all item bboxes.

## Algorithm properties

| Property | Guarantee |
|----------|-----------|
| Centerpiece | Always at `(0, 0)` |
| Connectivity | Every non-centerpiece item is adjacent to ≥1 placed item |
| No overlap | All pairs satisfy shape-aware collision with `gap` |
| Determinism | Same input → same output (fixed sort and tiebreak rules) |
| Failure | Explicit error if any item cannot be placed |

## Go package structure

```
layout/                         (nested module: github.com/.../gallery-wall/layout)
  types.go          Item, Shape, Params, Bounds
  validate.go       input and wall validation
  wall.go           wall bounds helpers
  geometry.go       footprints, bbox
  collision.go      shape-aware collision with gap
  adjacency.go      connectivity within gap
  placed.go         PlacedItem, Cluster
  order.go          centerpiece anchor, placement sort
  candidates.go     perimeter attachment, 1-unit slide
  score.go          blob heuristics, BestCandidate
  layout.go         Layout() entry point
  result.go         Result, PlacedResult, Anchor
  output.go         top-left coords, direction, adjacent IDs
  errors.go         typed validation and placement errors
  quality.go        regression metrics (aspect, hull fill, neighbors)
  doc.go            package documentation
  ALGORITHM.md      implementation walkthrough
  testdata/         golden JSON fixtures
docs/
  DESIGN.md         this document
  ARCHITECTURE.md   repo structure and data flow
```

Public API:

```go
func Layout(params Params) (Result, error)
func Validate(params Params) error
```

External install: `go get github.com/yashdalfthegray/gallery-wall/layout@v1.0.0`

## Test plan

### Validation tests

- Empty input → error
- Duplicate IDs → error
- Zero or negative dimensions → error
- Zero or multiple centerpieces → error

### Geometry tests

- Rectangle-rectangle collision with gap
- Circle-circle collision with gap
- Mixed shape collision (circle next to rectangle)
- Adjacency detection (touching within gap = connected; gap+1 = not connected)

### Layout tests

| Fixture | Asserts |
|---------|---------|
| Single centerpiece only | At origin, bounds match dimensions |
| Centerpiece + 1 item | Adjacent, no overlap, balanced placement |
| Centerpiece + 2 symmetric | Roughly mirrored placement |
| Mixed sizes (1 large + 6 small) | Connected cluster, no floaters, compact blob |
| Golden regression | Fixed input → exact positions (JSON fixture) |
| Unplaceable item | Returns error with item ID |
| Gap sensitivity | gap=2 vs gap=6 produces wider spacing |

### Blob quality tests (heuristic)

- Cluster convex hull fill ratio above a minimum threshold for standard fixtures
- No placed item farther from centerpiece than `2×` the largest item's max side (configurable sanity bound for tests)
- Every non-centerpiece item has ≥1 adjacent neighbor

## Future extensions

- `layoutMode` enum (organic blob vs symmetric grid vs radial)
- Wall reflow when cluster exceeds bounds (currently rejects at candidate generation)
- Frame rotation
- HTTP API wrapping `Layout()`
- Hanging guide PDF (cut sheet with measurements from anchor)

## Locked implementation preferences

| Preference | Decision |
|------------|----------|
| Target silhouette | Mostly round blob, wider than tall (typical wall proportions) |
| Attachment step size | 1 unit (precise candidate search) |
| Agent context | [.cursor/rules/](../.cursor/rules/) and [AGENTS.md](../AGENTS.md) |
