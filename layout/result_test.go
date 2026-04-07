package layout

import (
	"slices"
	"testing"
)

func TestToPlacedResultDerivedFields(t *testing.T) {
	placed := NewPlacedItem(Item{
		ID: "left", Height: 14, Width: 11, Shape: ShapeRectangle,
	}, -20.5, 0)

	result := toPlacedResult(placed, nil)

	assertFloat(t, "x", result.X, -26)
	assertFloat(t, "y", result.Y, -7)
	assertFloat(t, "centerX", result.CenterX, -20.5)
	assertFloat(t, "centerY", result.CenterY, 0)
	assertFloat(t, "offsetFromAnchor", result.OffsetFromAnchor, 20.5)
	if result.Direction != "W" {
		t.Fatalf("direction = %q, want W", result.Direction)
	}
}

func TestDirectionFromAnchor(t *testing.T) {
	tests := []struct {
		cx, cy float64
		want   string
	}{
		{0, 0, "C"},
		{10, 0, "E"},
		{0, 10, "S"},
		{-10, 0, "W"},
		{0, -10, "N"},
	}

	for _, tt := range tests {
		if got := directionFromAnchor(tt.cx, tt.cy); got != tt.want {
			t.Fatalf("direction(%v, %v) = %q, want %q", tt.cx, tt.cy, got, tt.want)
		}
	}
}

func TestBuildResultBounds(t *testing.T) {
	placed := Cluster{
		NewPlacedItem(Item{ID: "main", Height: 20, Width: 40, Shape: ShapeRectangle, Centerpiece: true}, 0, 0),
		NewPlacedItem(Item{ID: "right", Height: 10, Width: 10, Shape: ShapeSquare}, 12, 0),
	}

	result := buildResult(placed, "main", 2)
	assertFloat(t, "minX", result.Bounds.MinX, -20)
	assertFloat(t, "maxX", result.Bounds.MaxX, 20)
	if result.Anchor.ItemID != "main" {
		t.Fatalf("anchor id = %q, want main", result.Anchor.ItemID)
	}
}

func TestAdjacentIDsFor(t *testing.T) {
	gap := 2
	placed := Cluster{
		anchorCenterpiece(Item{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true}),
		NewPlacedItem(Item{ID: "right", Height: 10, Width: 10, Shape: ShapeSquare}, 12, 0),
		NewPlacedItem(Item{ID: "below", Height: 10, Width: 10, Shape: ShapeSquare}, 0, 12),
		NewPlacedItem(Item{ID: "corner", Height: 10, Width: 10, Shape: ShapeSquare}, 12, 12),
	}

	mainIDs := adjacentIDsFor(placed, placed[0], gap)
	if !slices.Equal(mainIDs, []string{"below", "right"}) {
		t.Fatalf("main adjacent = %v, want [below right]", mainIDs)
	}

	cornerIDs := adjacentIDsFor(placed, placed[3], gap)
	if !slices.Equal(cornerIDs, []string{"below", "right"}) {
		t.Fatalf("corner adjacent = %v, want [below right]", cornerIDs)
	}

	rightIDs := adjacentIDsFor(placed, placed[1], gap)
	if !slices.Equal(rightIDs, []string{"corner", "main"}) {
		t.Fatalf("right adjacent = %v, want [corner main]", rightIDs)
	}
}

func TestLayoutIncludesAdjacentIDs(t *testing.T) {
	gap := 2
	result, err := Layout(Params{
		Gap: gap,
		Items: []Item{
			{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
			{ID: "right", Height: 10, Width: 10, Shape: ShapeSquare},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	byID := make(map[string]PlacedResult, len(result.Items))
	for _, item := range result.Items {
		byID[item.ID] = item
	}

	main := byID["main"]
	right := byID["right"]

	if !slices.Equal(main.AdjacentIDs, []string{"right"}) {
		t.Fatalf("main adjacent = %v, want [right]", main.AdjacentIDs)
	}
	if !slices.Equal(right.AdjacentIDs, []string{"main"}) {
		t.Fatalf("right adjacent = %v, want [main]", right.AdjacentIDs)
	}
}
