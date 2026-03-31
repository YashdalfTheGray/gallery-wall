package layout

import "testing"

func TestLayoutCenterpieceOnly(t *testing.T) {
	result, err := Layout(Params{
		Gap: 2,
		Items: []Item{
			{ID: "main", Height: 36, Width: 24, Shape: ShapeRectangle, Centerpiece: true},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Anchor.ItemID != "main" {
		t.Fatalf("anchor id = %q, want main", result.Anchor.ItemID)
	}
	if len(result.Items) != 1 {
		t.Fatalf("item count = %d, want 1", len(result.Items))
	}

	main := result.Items[0]
	if main.CenterX != 0 || main.CenterY != 0 {
		t.Fatalf("centerpiece center = (%v, %v), want (0, 0)", main.CenterX, main.CenterY)
	}
	if main.Direction != "C" {
		t.Fatalf("direction = %q, want C", main.Direction)
	}

	assertFloat(t, "minX", result.Bounds.MinX, -12)
	assertFloat(t, "minY", result.Bounds.MinY, -18)
	assertFloat(t, "maxX", result.Bounds.MaxX, 12)
	assertFloat(t, "maxY", result.Bounds.MaxY, 18)
}

func TestLayoutValidationBeforePlacement(t *testing.T) {
	_, err := Layout(Params{
		Gap: 2,
		Items: []Item{
			{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
			{ID: "a", Height: 8, Width: 8, Shape: ShapeSquare},
		},
	})
	if !IsValidationCode(err, ValidationDuplicateID) {
		t.Fatalf("expected duplicate_id, got %v", err)
	}
}

func TestLayoutPlacesSecondItem(t *testing.T) {
	gap := 2
	result, err := Layout(Params{
		Gap: gap,
		Items: []Item{
			{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
			{ID: "side", Height: 10, Width: 10, Shape: ShapeSquare},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var main, side PlacedResult
	for _, item := range result.Items {
		switch item.ID {
		case "main":
			main = item
		case "side":
			side = item
		}
	}

	if main.ID == "" || side.ID == "" {
		t.Fatalf("missing items in result: %+v", result.Items)
	}

	placed := Cluster{
		NewPlacedItem(Item{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true}, main.CenterX, main.CenterY),
		NewPlacedItem(Item{ID: "side", Height: 10, Width: 10, Shape: ShapeSquare}, side.CenterX, side.CenterY),
	}
	if Collides(placed[0].Footprint, placed[1].Footprint, gap) {
		t.Fatal("placed items overlap")
	}
	if !Adjacent(placed[0].Footprint, placed[1].Footprint, gap) {
		t.Fatal("placed items are not adjacent")
	}
}

func TestLayoutPlacesMultipleItems(t *testing.T) {
	gap := 2
	result, err := Layout(Params{
		Gap: gap,
		Items: []Item{
			{ID: "main", Height: 20, Width: 20, Shape: ShapeRectangle, Centerpiece: true},
			{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare},
			{ID: "b", Height: 10, Width: 10, Shape: ShapeSquare},
			{ID: "c", Height: 10, Width: 10, Shape: ShapeSquare},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 4 {
		t.Fatalf("item count = %d, want 4", len(result.Items))
	}

	cluster := make(Cluster, 0, len(result.Items))
	for _, item := range result.Items {
		cluster = append(cluster, NewPlacedItem(Item{
			ID: item.ID, Height: item.Height, Width: item.Width, Shape: item.Shape,
		}, item.CenterX, item.CenterY))
	}

	for i := 0; i < len(cluster); i++ {
		for j := i + 1; j < len(cluster); j++ {
			if Collides(cluster[i].Footprint, cluster[j].Footprint, gap) {
				t.Fatalf("overlap between %s and %s", cluster[i].Item.ID, cluster[j].Item.ID)
			}
		}
	}

	for _, placed := range cluster {
		if placed.Item.ID == "main" {
			continue
		}
		neighbors := 0
		for _, other := range cluster {
			if placed.Item.ID == other.Item.ID {
				continue
			}
			if Adjacent(placed.Footprint, other.Footprint, gap) {
				neighbors++
			}
		}
		if neighbors == 0 {
			t.Fatalf("floater detected: %s", placed.Item.ID)
		}
	}
}

func TestPlaceItemsReturnsCannotPlaceError(t *testing.T) {
	placed := Cluster{anchorCenterpiece(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true,
	})}
	stuck := Item{ID: "stuck", Height: 10, Width: 10, Shape: ShapeSquare}

	_, err := placeItemFromCandidates(placed, stuck, 2, nil)
	if err == nil {
		t.Fatal("expected placement error")
	}
	if !IsPlacementCode(err, PlacementCannotPlace) {
		t.Fatalf("expected cannot_place_item, got %v", err)
	}
	pe, ok := err.(*PlacementError)
	if !ok || pe.ItemID != "stuck" {
		t.Fatalf("expected placement error for stuck, got %v", err)
	}
}

func TestPlaceItemsExpandsOutwardWhenSurrounded(t *testing.T) {
	gap := 2
	surrounded := Cluster{
		anchorCenterpiece(Item{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true}),
		NewPlacedItem(Item{ID: "r", Height: 10, Width: 10, Shape: ShapeSquare}, 12, 0),
		NewPlacedItem(Item{ID: "l", Height: 10, Width: 10, Shape: ShapeSquare}, -12, 0),
		NewPlacedItem(Item{ID: "t", Height: 10, Width: 10, Shape: ShapeSquare}, 0, -12),
		NewPlacedItem(Item{ID: "b", Height: 10, Width: 10, Shape: ShapeSquare}, 0, 12),
		NewPlacedItem(Item{ID: "tr", Height: 10, Width: 10, Shape: ShapeSquare}, 12, -12),
		NewPlacedItem(Item{ID: "tl", Height: 10, Width: 10, Shape: ShapeSquare}, -12, -12),
		NewPlacedItem(Item{ID: "br", Height: 10, Width: 10, Shape: ShapeSquare}, 12, 12),
		NewPlacedItem(Item{ID: "bl", Height: 10, Width: 10, Shape: ShapeSquare}, -12, 12),
	}

	stuck := Item{ID: "stuck", Height: 10, Width: 10, Shape: ShapeSquare}
	if len(GenerateCandidates(stuck, surrounded, gap, nil)) == 0 {
		t.Fatal("expected outward attachment candidates")
	}

	placed, err := placeItems(surrounded, []Item{stuck}, gap, nil)
	if err != nil {
		t.Fatalf("expected outward placement, got %v", err)
	}
	if len(placed) != len(surrounded)+1 {
		t.Fatalf("placed count = %d, want %d", len(placed), len(surrounded)+1)
	}
}

func TestPlacementErrorJSONCode(t *testing.T) {
	err := newCannotPlaceError("stuck")
	if !IsPlacementCode(err, PlacementCannotPlace) {
		t.Fatal("expected cannot_place_item code")
	}
	if err.ItemID != "stuck" {
		t.Fatalf("itemId = %q, want stuck", err.ItemID)
	}
}

func TestLayoutNonSquareCirclesRespectGap(t *testing.T) {
	// Regression: non-square circles use inscribed radius for collision, not the
	// full width×height box. Layout must not report collisions; render uses the
	// same inscribed circle (see web/src/svg.ts).
	gap := 2
	items := []Item{
		{ID: "main", Height: 13, Width: 15, Shape: ShapeRectangle, Centerpiece: true},
		{ID: "p01", Height: 12, Width: 10, Shape: ShapeEllipse},
		{ID: "p02", Height: 7, Width: 10, Shape: ShapeCircle},
		{ID: "p03", Height: 5, Width: 9, Shape: ShapeEllipse},
		{ID: "p04", Height: 7, Width: 10, Shape: ShapeCircle},
		{ID: "p05", Height: 6, Width: 10, Shape: ShapeRectangle},
		{ID: "p06", Height: 6, Width: 7, Shape: ShapeSquare},
		{ID: "p07", Height: 10, Width: 7, Shape: ShapeRectangle},
		{ID: "p08", Height: 5, Width: 10, Shape: ShapeCircle},
		{ID: "p09", Height: 6, Width: 5, Shape: ShapeCircle},
		{ID: "p10", Height: 8, Width: 9, Shape: ShapeSquare},
		{ID: "p11", Height: 8, Width: 8, Shape: ShapeEllipse},
		{ID: "p12", Height: 7, Width: 10, Shape: ShapeRectangle},
		{ID: "p13", Height: 6, Width: 9, Shape: ShapeEllipse},
		{ID: "p14", Height: 11, Width: 10, Shape: ShapeCircle},
		{ID: "p15", Height: 12, Width: 7, Shape: ShapeCircle},
		{ID: "p16", Height: 8, Width: 9, Shape: ShapeEllipse},
		{ID: "p17", Height: 10, Width: 10, Shape: ShapeEllipse},
		{ID: "p18", Height: 5, Width: 8, Shape: ShapeEllipse},
		{ID: "p19", Height: 7, Width: 8, Shape: ShapeCircle},
	}

	result, err := Layout(Params{Gap: gap, Items: items})
	if err != nil {
		t.Fatalf("layout: %v", err)
	}
	if len(result.Items) != len(items) {
		t.Fatalf("item count = %d, want %d", len(result.Items), len(items))
	}

	byID := make(map[string]Item, len(items))
	for _, item := range items {
		byID[item.ID] = item
	}

	cluster := make(Cluster, 0, len(result.Items))
	for _, placed := range result.Items {
		item := byID[placed.ID]
		cluster = append(cluster, NewPlacedItem(item, placed.CenterX, placed.CenterY))
	}

	for i := 0; i < len(cluster); i++ {
		for j := i + 1; j < len(cluster); j++ {
			if Collides(cluster[i].Footprint, cluster[j].Footprint, gap) {
				t.Fatalf("overlap between %s and %s", cluster[i].Item.ID, cluster[j].Item.ID)
			}
		}
	}
}
