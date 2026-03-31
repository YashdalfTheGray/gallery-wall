package layout

import "math"

// Layout places all frames organically around the centerpiece.
//
// The centerpiece center is anchored at (0, 0). Returns [ValidationError] for
// bad input or [PlacementError] when a frame cannot be attached to the cluster.
func Layout(params Params) (Result, error) {
	if err := Validate(params); err != nil {
		return Result{}, err
	}

	centerpiece, remaining, err := findCenterpiece(params.Items)
	if err != nil {
		return Result{}, err
	}

	placed, err := placeItems(Cluster{anchorCenterpiece(centerpiece)}, remaining, params.Gap, wallPtr(params))
	if err != nil {
		return Result{}, err
	}

	return buildResult(placed, centerpiece.ID, params.Gap), nil
}

func wallPtr(params Params) *Bounds {
	if wall, ok := params.WallBounds(); ok {
		return &wall
	}
	return nil
}

func placeItems(placed Cluster, remaining []Item, gap int, wall *Bounds) (Cluster, error) {
	for _, item := range sortForPlacement(remaining) {
		var err error
		placed, err = placeItem(placed, item, gap, wall)
		if err != nil {
			return nil, err
		}
	}
	return placed, nil
}

func placeItem(placed Cluster, item Item, gap int, wall *Bounds) (Cluster, error) {
	return placeItemFromCandidates(placed, item, gap, GenerateCandidates(item, placed, gap, wall))
}

func placeItemFromCandidates(placed Cluster, item Item, gap int, candidates []Candidate) (Cluster, error) {
	best, ok := BestCandidate(candidates, placed, gap)
	if !ok {
		return nil, newCannotPlaceError(item.ID)
	}
	return append(placed, NewPlacedItem(item, best.CenterX, best.CenterY)), nil
}

func buildResult(placed Cluster, centerpieceID string, gap int) Result {
	items := make([]PlacedResult, 0, len(placed))
	for _, p := range placed {
		items = append(items, toPlacedResult(p, adjacentIDsFor(placed, p, gap)))
	}

	return Result{
		Anchor: Anchor{
			ItemID:  centerpieceID,
			CenterX: 0,
			CenterY: 0,
		},
		Items:  items,
		Bounds: placed.Bounds(),
	}
}

func toPlacedResult(placed PlacedItem, adjacentIDs []string) PlacedResult {
	cx := placed.Footprint.CenterX
	cy := placed.Footprint.CenterY
	halfW := float64(placed.Item.Width) / 2
	halfH := float64(placed.Item.Height) / 2

	return PlacedResult{
		ID:               placed.Item.ID,
		CenterX:          cx,
		CenterY:          cy,
		X:                cx - halfW,
		Y:                cy - halfH,
		Width:            placed.Item.Width,
		Height:           placed.Item.Height,
		Shape:            placed.Item.Shape,
		OffsetFromAnchor: math.Hypot(cx, cy),
		Direction:        directionFromAnchor(cx, cy),
		AdjacentIDs:      adjacentIDs,
	}
}

func directionFromAnchor(cx, cy float64) string {
	if cx == 0 && cy == 0 {
		return "C"
	}

	angle := math.Atan2(cy, cx) * 180 / math.Pi
	if angle < 0 {
		angle += 360
	}

	switch {
	case angle >= 337.5 || angle < 22.5:
		return "E"
	case angle < 67.5:
		return "SE"
	case angle < 112.5:
		return "S"
	case angle < 157.5:
		return "SW"
	case angle < 202.5:
		return "W"
	case angle < 247.5:
		return "NW"
	case angle < 292.5:
		return "N"
	default:
		return "NE"
	}
}
