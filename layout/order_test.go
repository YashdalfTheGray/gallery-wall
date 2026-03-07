package layout

import (
	"testing"
)

func TestFindCenterpiece(t *testing.T) {
	items := []Item{
		{ID: "left", Height: 8, Width: 8, Shape: ShapeSquare, Centerpiece: false},
		{ID: "main", Height: 24, Width: 18, Shape: ShapeRectangle, Centerpiece: true},
		{ID: "right", Height: 8, Width: 8, Shape: ShapeSquare, Centerpiece: false},
	}

	centerpiece, rest, err := findCenterpiece(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if centerpiece.ID != "main" {
		t.Fatalf("centerpiece id = %q, want main", centerpiece.ID)
	}
	if len(rest) != 2 || rest[0].ID != "left" || rest[1].ID != "right" {
		t.Fatalf("rest = %#v", rest)
	}

	t.Run("no centerpiece", func(t *testing.T) {
		_, _, err := findCenterpiece([]Item{{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare}})
		if !IsValidationCode(err, ValidationNoCenterpiece) {
			t.Fatalf("expected no_centerpiece, got %v", err)
		}
	})

	t.Run("multiple centerpieces", func(t *testing.T) {
		_, _, err := findCenterpiece([]Item{
			{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
			{ID: "b", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
		})
		if !IsValidationCode(err, ValidationMultipleCenterpiece) {
			t.Fatalf("expected multiple_centerpieces, got %v", err)
		}
	})
}

func TestSortForPlacement(t *testing.T) {
	items := []Item{
		{ID: "small", Height: 8, Width: 8, Shape: ShapeSquare},
		{ID: "large", Height: 20, Width: 10, Shape: ShapeRectangle},
		{ID: "mid", Height: 10, Width: 10, Shape: ShapeSquare},
		{ID: "tie-b", Height: 12, Width: 5, Shape: ShapeRectangle},
		{ID: "tie-a", Height: 10, Width: 6, Shape: ShapeRectangle},
	}

	sorted := sortForPlacement(items)
	got := []string{
		sorted[0].ID,
		sorted[1].ID,
		sorted[2].ID,
		sorted[3].ID,
		sorted[4].ID,
	}
	want := []string{"large", "mid", "small", "tie-b", "tie-a"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order = %v, want %v", got, want)
		}
	}

	// Deterministic for identical sort keys.
	dup := sortForPlacement([]Item{
		{ID: "b", Height: 10, Width: 10, Shape: ShapeSquare},
		{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare},
	})
	if dup[0].ID != "a" || dup[1].ID != "b" {
		t.Fatalf("id tiebreak order = [%s %s], want [a b]", dup[0].ID, dup[1].ID)
	}
}

func TestAnchorCenterpiece(t *testing.T) {
	item := Item{ID: "main", Height: 36, Width: 24, Shape: ShapeRectangle, Centerpiece: true}
	placed := anchorCenterpiece(item)

	if placed.Footprint.CenterX != 0 || placed.Footprint.CenterY != 0 {
		t.Fatalf("center = (%v, %v), want (0, 0)", placed.Footprint.CenterX, placed.Footprint.CenterY)
	}
	if placed.Item.ID != "main" {
		t.Fatalf("item id = %q, want main", placed.Item.ID)
	}
}
