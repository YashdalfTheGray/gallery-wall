# Layout algorithm

Organic gallery-wall placement: one **centerpiece** at `(0, 0)`, every other frame attaches edge-to-edge (with gap) into a single connected **blob**.

**Entry point:** [`Layout()`](layout.go) · **Spec:** [`docs/DESIGN.md`](../docs/DESIGN.md)

---

## Pipeline

```mermaid
flowchart LR
  IN[Params] --> V[Validate]
  V -->|ok| A[Anchor centerpiece]
  A --> P[Place loop]
  P --> C[Candidates]
  C --> F[Filter]
  F --> S[Score and pick best]
  S --> P
  P -->|done| R[buildResult]
  R --> OUT[Result]
  V -->|err| E[ValidationError]
  S -->|none| E2[PlacementError]
```

---

## Source map

| Stage        | File                                                | Role                                    |
| ------------ | --------------------------------------------------- | --------------------------------------- |
| Types        | [`types.go`](types.go)                              | `Item`, `Params`, `Shape`               |
| Validate     | [`validate.go`](validate.go)                        | Input + wall checks                     |
| Errors       | [`errors.go`](errors.go)                            | Typed error codes                       |
| Order        | [`order.go`](order.go)                              | Find centerpiece, sort by size          |
| Geometry     | [`geometry.go`](geometry.go)                        | `Footprint`, bboxes, half-extents       |
| Collision    | [`collision.go`](collision.go)                      | Overlap test (shape-aware)              |
| Adjacency    | [`adjacency.go`](adjacency.go)                      | Touching within gap                     |
| Candidates   | [`candidates.go`](candidates.go)                    | Side + corner attachments, 1-unit slide |
| Scoring      | [`score.go`](score.go)                              | Blob heuristics, pick lowest score      |
| Orchestrator | [`layout.go`](layout.go)                            | `Layout`, placement loop, output        |
| Wall         | [`wall.go`](wall.go)                                | Optional bounding box filter            |
| State        | [`placed.go`](placed.go)                            | `PlacedItem`, `Cluster`, `Bounds`       |
| Output       | [`result.go`](result.go) · [`output.go`](output.go) | `Result`, neighbors, direction          |
| Quality      | [`quality.go`](quality.go)                          | Test helpers (aspect, spread, …)        |

---

## Example — 12 paintings

```json
{
  "gap": 2,
  "items": [
    {
      "id": "main",
      "height": 16,
      "width": 14,
      "shape": "rectangle",
      "centerpiece": true
    },
    { "id": "p01", "height": 10, "width": 8, "shape": "rectangle" },
    { "id": "p02", "height": 8, "width": 8, "shape": "square" },
    { "id": "p03", "height": 12, "width": 9, "shape": "rectangle" },
    { "id": "p04", "height": 6, "width": 6, "shape": "circle" },
    { "id": "p05", "height": 9, "width": 7, "shape": "ellipse" },
    { "id": "p06", "height": 11, "width": 8, "shape": "rectangle" },
    { "id": "p07", "height": 7, "width": 10, "shape": "rectangle" },
    { "id": "p08", "height": 8, "width": 6, "shape": "square" },
    { "id": "p09", "height": 10, "width": 10, "shape": "square" },
    { "id": "p10", "height": 5, "width": 7, "shape": "rectangle" },
    { "id": "p11", "height": 9, "width": 6, "shape": "rectangle" }
  ]
}
```

**Result:** 38×40 blob · 11 diagonal placements · centerpiece `main` at origin.

---

## End-to-end sequence

```mermaid
sequenceDiagram
  participant Client
  participant Layout as layout.go
  participant Validate as validate.go
  participant Order as order.go
  participant Cand as candidates.go
  participant Score as score.go
  participant Result as result.go

  Client->>Layout: Layout(params)
  Layout->>Validate: Validate(params)
  Validate-->>Layout: ok

  Layout->>Order: findCenterpiece(items)
  Order-->>Layout: main + 11 remaining

  Layout->>Order: anchorCenterpiece(main) → (0,0)
  Layout->>Order: sortForPlacement(remaining)

  loop each item largest-first
    Layout->>Cand: GenerateCandidates(item, cluster, gap, wall)
    Note over Cand: side and corner attach, slide, dedupe
    Cand->>Cand: filter collision, neighbor, wall
    Cand-->>Layout: candidates
    Layout->>Score: BestCandidate(candidates)
    Score-->>Layout: lowest score + tiebreak
    Layout->>Layout: append to cluster
  end

  Layout->>Result: buildResult(cluster)
  Result-->>Client: Result (positions, bounds, neighbors)
```

---

## Step 1 — Validate

[`Validate()`](validate.go) · [`validateWall()`](wall.go)

```mermaid
flowchart TD
  A[items non-empty?] -->|no| E1[empty_input]
  A -->|yes| B[unique ids?]
  B -->|no| E2[duplicate_id]
  B -->|yes| C["dimensions positive?"]
  C -->|no| E3[invalid_dimensions]
  C -->|yes| D[exactly 1 centerpiece?]
  D -->|0| E4[no_centerpiece]
  D -->|2 or more| E5[multiple_centerpieces]
  D -->|1| W{wall set?}
  W -->|partial| E6[invalid_wall]
  W -->|yes| F[centerpiece fits wall?]
  F -->|no| E7[centerpiece_exceeds_wall]
  F -->|ok| OK[proceed]
  W -->|no| OK
```

---

## Step 2 — Anchor + sort

[`findCenterpiece()`](order.go) · [`sortForPlacement()`](order.go) · [`anchorCenterpiece()`](order.go)

| Order | ID       | Area | Why first                          |
| ----: | -------- | ---: | ---------------------------------- |
|     — | **main** |  224 | Centerpiece → `(0, 0)` immediately |
|     1 | **p03**  |  108 | Largest satellite                  |
|     2 | **p09**  |  100 |                                    |
|     3 | **p06**  |   88 |                                    |
|     4 | **p01**  |   80 |                                    |
|     5 | **p07**  |   70 |                                    |
|     6 | **p02**  |   64 |                                    |
|     7 | **p05**  |   63 |                                    |
|     8 | **p11**  |   54 |                                    |
|     9 | **p08**  |   48 |                                    |
|    10 | **p04**  |   36 |                                    |
|    11 | **p10**  |   35 | Smallest last — fills gaps         |

Large frames first → compact core; small frames last → tuck into corners.

---

## Step 3 — Placement loop (trace)

Each iteration: **generate → filter → score → append**.

### After centerpiece

```
Cluster: [main @ (0,0)]
```

### Placement 1 — `p03` (12×9)

```mermaid
flowchart LR
  subgraph candidates ["Candidates from main"]
    L[west side]
    R[east side]
    T[north side]
    B[south side]
    NW[corner NW]
    NE[corner NE]
    SW[corner SW]
    SE[corner SE]
  end
  candidates --> F[filter + score]
  F --> WIN["p03 at west side"]
```

| Field     | Value       |
| --------- | ----------- |
| Position  | `(-14, -2)` |
| Direction | W           |
| Touches   | main        |

### Placements 2–4 (cluster grows)

|   # | ID  | Position    | Dir | New adjacency  |
| --: | --- | ----------- | --- | -------------- |
|   2 | p09 | `(-2, 15)`  | S   | main           |
|   3 | p06 | `(13, 2)`   | E   | main           |
|   4 | p01 | `(-13, 11)` | SW  | main, p03, p09 |

```mermaid
graph TD
  main((main))
  p03 --- main
  p09 --- main
  p06 --- main
  p01 --- main
  p01 --- p03
  p01 --- p09
```

### Placements 5–11 (corners fill)

|   # | ID  | Position     | Dir |
| --: | --- | ------------ | --- |
|   5 | p07 | `(14, -9)`   | NE  |
|   6 | p02 | `(3, -14)`   | N   |
|   7 | p05 | `(-7, -14)`  | NW  |
|   8 | p11 | `(8, 14)`    | SE  |
|   9 | p08 | `(16, 14)`   | SE  |
|  10 | p04 | `(-16, -13)` | NW  |
|  11 | p10 | `(12, -17)`  | NE  |

Final bounds: **38 × 40** · all 12 connected · no floaters.

```mermaid
graph TD
  main((main))
  p03 --- main
  p04 --- p03
  p05 --- main
  p05 --- p03
  p01 --- main
  p01 --- p03
  p01 --- p09
  p09 --- main
  p09 --- p11
  p11 --- main
  p11 --- p08
  p06 --- main
  p06 --- p07
  p07 --- main
  p07 --- p02
  p07 --- p10
  p02 --- main
  p02 --- p10
```

---

## Candidate generation

[`GenerateCandidates()`](candidates.go) · [`sideCandidates()`](candidates.go) · [`cornerCandidates()`](candidates.go)

```mermaid
flowchart TD
  A[For each placed frame] --> B[4 side attachments]
  A --> C[4 corner attachments]
  B --> D["Slide 1 unit along edge"]
  C --> D
  D --> E[Dedupe by rounded position]
  E --> F{Valid?}
  F -->|Collides| X1[reject]
  F -->|No neighbor| X2[reject floater]
  F -->|Outside wall| X3[reject]
  F -->|ok| G[keep]
```

**Attachment offset** uses shape-aware half-extents from [`HalfExtents()`](geometry.go):

| Shape              | Extent                  |
| ------------------ | ----------------------- |
| rectangle / square | half width, half height |
| circle             | radius = min(w,h)/2     |
| ellipse            | semi-axes w/2, h/2      |

Gap `g` is added between extents: `anchor_half + g + item_half`.

---

## Collision & adjacency

[`Collides()`](collision.go) · [`Adjacent()`](adjacency.go)

```mermaid
flowchart LR
  subgraph collision ["Collides?"]
    R[rect vs rect AABB]
    C[circle vs circle distance]
    E[ellipse conservative AABB]
  end
  subgraph adjacent ["Adjacent?"]
    N[not colliding and gap within g]
  end
  collision -->|yes| REJ[reject candidate]
  adjacent -->|no| REJ
  adjacent -->|yes| ACC[valid attachment]
```

---

## Scoring

[`ScoreCandidate()`](score.go) · [`BestCandidate()`](score.go)

**Lower score wins.** Rewards (−) pull toward good blobs; penalties (+) push away.

```mermaid
flowchart TB
  SC[ScoreCandidate]
  SC --> PEN["Penalties +"]
  SC --> REW["Rewards -"]
  PEN --> P1[compactness]
  PEN --> P2[balance]
  PEN --> P3[blob smoothness]
  PEN --> P4[concavity]
  PEN --> P5[collinearity]
  PEN --> P6[axis attachment]
  PEN --> P7[axis-only placement]
  PEN --> P8[cluster area]
  REW --> R1[local continuity]
  REW --> R2[diagonal fill]
  REW --> R3[quadrant fill]
```

| Signal                   | Intent                              |
| ------------------------ | ----------------------------------- |
| Compactness              | Stay near anchor                    |
| Balance                  | Even mass in quadrants              |
| Blob smoothness          | Wide-round silhouette (~1.4 aspect) |
| Concavity                | Avoid L-shaped bays                 |
| Collinearity             | No horizontal/vertical combs        |
| Diagonal / quadrant fill | Corner blobs, less wasted wall      |
| Local continuity         | Prefer slots touching 2+ neighbors  |

Tiebreak in [`candidateLess()`](score.go): prefer diagonal → spread off flat axis → stable x/y order.

---

## Output

[`buildResult()`](layout.go) · [`toPlacedResult()`](layout.go) · [`adjacentIDsFor()`](output.go)

```mermaid
flowchart LR
  C[Cluster] --> T[toPlacedResult]
  T --> P["centerX, centerY, direction"]
  C --> A[adjacentIDsFor]
  A --> N["sorted neighbor ids"]
  C --> B[cluster Bounds]
  B --> R[result bounds]
```

**Coordinate system:** +X right, +Y down · anchor at centerpiece center.

---

## Optional wall

[`Params.WallWidth / WallHeight`](types.go) · [`footprintFitsWall()`](wall.go)

Wall is a rectangle **centered on anchor**. Candidates whose dimension bbox would cross the edge are filtered out before scoring. Omitted = unbounded.

---

## Errors

| When          | Code                | File                     |
| ------------- | ------------------- | ------------------------ |
| Bad input     | `validation_*`      | [`errors.go`](errors.go) |
| No valid slot | `cannot_place_item` | [`errors.go`](errors.go) |

Use `errors.Is(err, &ValidationError{Code: …})` or [`IsValidationCode()`](errors.go).

---

## Run the example

[`example_test.go`](example_test.go) uses small inline datasets (self-contained runnable docs). Golden regression uses [`testdata/`](testdata/).

```bash
go test ./layout/ -run Example -v
go test ./layout/ -run Golden -v
```

**Visualize** with [`cmd/gallery-svg/`](../cmd/gallery-svg/) — needs a **result** JSON (positions), not params:

```bash
go run ./cmd/gallery-svg/
go run ./cmd/gallery-svg/ \
  layout/testdata/centerpiece_six_mixed_result.json \
  layout/testdata/six_mixed.svg
```
