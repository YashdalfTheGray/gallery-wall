package layout

import (
	"math"
	"testing"
)

func TestHalfExtents(t *testing.T) {
	tests := []struct {
		name       string
		item       Item
		wantHalfX  float64
		wantHalfY  float64
	}{
		{
			name: "rectangle",
			item: Item{Shape: ShapeRectangle, Height: 12, Width: 20},
			wantHalfX: 10, wantHalfY: 6,
		},
		{
			name: "circle wider than tall",
			item: Item{Shape: ShapeCircle, Height: 10, Width: 16},
			wantHalfX: 5, wantHalfY: 5,
		},
		{
			name: "ellipse",
			item: Item{Shape: ShapeEllipse, Height: 12, Width: 16},
			wantHalfX: 8, wantHalfY: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hx, hy := halfExtents(tt.item)
			assertFloat(t, "halfX", hx, tt.wantHalfX)
			assertFloat(t, "halfY", hy, tt.wantHalfY)
		})
	}
}

func TestSideAttachmentPositions(t *testing.T) {
	gap := 2
	anchor := NewPlacedItem(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeRectangle, Centerpiece: true,
	}, 0, 0)
	item := Item{ID: "left", Height: 10, Width: 10, Shape: ShapeRectangle}

	t.Run("left aligned", func(t *testing.T) {
		cands := sideCandidates(anchor, item, gap)
		found := false
		for _, c := range cands {
			if c.CenterX == -12 && c.CenterY == 0 {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected left attachment at (-12, 0)")
		}
	})

	t.Run("right aligned", func(t *testing.T) {
		found := false
		for _, c := range sideCandidates(anchor, item, gap) {
			if c.CenterX == 12 && c.CenterY == 0 {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected right attachment at (12, 0)")
		}
	})

	t.Run("above aligned", func(t *testing.T) {
		found := false
		for _, c := range sideCandidates(anchor, item, gap) {
			if c.CenterX == 0 && c.CenterY == -12 {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected above attachment at (0, -12)")
		}
	})

	t.Run("below aligned", func(t *testing.T) {
		found := false
		for _, c := range sideCandidates(anchor, item, gap) {
			if c.CenterX == 0 && c.CenterY == 12 {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected below attachment at (0, 12)")
		}
	})
}

func TestSideAttachmentSlideCount(t *testing.T) {
	anchor := NewPlacedItem(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeRectangle, Centerpiece: true,
	}, 0, 0)
	item := Item{ID: "side", Height: 10, Width: 10, Shape: ShapeRectangle}

	cands := sideCandidates(anchor, item, 2)
	// 21 Y slides × 2 horizontal sides + 21 X slides × 2 vertical sides = 84
	if len(cands) != 84 {
		t.Fatalf("side candidate count = %d, want 84", len(cands))
	}
}

func TestCornerAttachment(t *testing.T) {
	gap := 2
	anchor := NewPlacedItem(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeRectangle, Centerpiece: true,
	}, 0, 0)
	item := Item{ID: "corner", Height: 10, Width: 10, Shape: ShapeRectangle}

	cands := cornerCandidates(anchor, item, gap)
	found := false
	for _, c := range cands {
		if c.CenterX == -12 && c.CenterY == -12 {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected top-left corner attachment at (-12, -12)")
	}
}

func TestCornerSlideProducesMultipleCandidates(t *testing.T) {
	anchor := NewPlacedItem(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeRectangle, Centerpiece: true,
	}, 0, 0)
	item := Item{ID: "corner", Height: 10, Width: 10, Shape: ShapeRectangle}

	cands := cornerCandidates(anchor, item, 2)
	if len(cands) < 4 {
		t.Fatalf("expected multiple corner slide candidates, got %d", len(cands))
	}
}

func TestGenerateCandidatesTwoItemCluster(t *testing.T) {
	placed := Cluster{anchorCenterpiece(Item{
		ID: "main", Height: 20, Width: 20, Shape: ShapeRectangle, Centerpiece: true,
	})}
	item := Item{ID: "side", Height: 10, Width: 10, Shape: ShapeRectangle}

	cands := GenerateCandidates(item, placed, 2, nil)
	if len(cands) == 0 {
		t.Fatal("expected non-empty candidates for simple cluster")
	}
}

func TestFilterRejectsCollisions(t *testing.T) {
	placed := Cluster{anchorCenterpiece(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeRectangle, Centerpiece: true,
	})}
	item := Item{ID: "side", Height: 10, Width: 10, Shape: ShapeRectangle}

	cands := filterCandidates(placed, 2, nil, []Candidate{{
		Item: item, CenterX: 0, CenterY: 0, // overlaps centerpiece
	}})
	if len(cands) != 0 {
		t.Fatalf("expected collision to be filtered, got %d", len(cands))
	}
}

func TestFilterRejectsFloaters(t *testing.T) {
	placed := Cluster{anchorCenterpiece(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeRectangle, Centerpiece: true,
	})}
	item := Item{ID: "side", Height: 10, Width: 10, Shape: ShapeRectangle}

	cands := filterCandidates(placed, 2, nil, []Candidate{{
		Item: item, CenterX: 100, CenterY: 0, // far away, not adjacent
	}})
	if len(cands) != 0 {
		t.Fatalf("expected floater to be filtered, got %d", len(cands))
	}
}

func TestDedupeCandidates(t *testing.T) {
	item := Item{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare}
	raw := []Candidate{
		{Item: item, CenterX: 12, CenterY: 0},
		{Item: item, CenterX: 12, CenterY: 0},
		{Item: item, CenterX: 12.0000005, CenterY: 0},
	}
	deduped := dedupeCandidates(raw)
	if len(deduped) != 1 {
		t.Fatalf("deduped count = %d, want 1", len(deduped))
	}
}

func TestGenerateCandidatesAllValid(t *testing.T) {
	placed := Cluster{anchorCenterpiece(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeRectangle, Centerpiece: true,
	})}
	item := Item{ID: "side", Height: 10, Width: 10, Shape: ShapeRectangle}
	gap := 2

	for _, c := range GenerateCandidates(item, placed, gap, nil) {
		if !isValidCandidate(c, placed, gap, nil) {
			t.Fatalf("invalid candidate leaked: %+v", c)
		}
	}
}

func TestCircleSideAttachmentUsesRadius(t *testing.T) {
	gap := 2
	anchor := NewPlacedItem(Item{
		ID: "main", Height: 10, Width: 16, Shape: ShapeCircle, Centerpiece: true,
	}, 0, 0)
	item := Item{ID: "side", Height: 10, Width: 10, Shape: ShapeRectangle}

	found := false
	for _, c := range sideCandidates(anchor, item, gap) {
		// circle r=5, rect half=5, gap=2 → center distance 12
		if math.Abs(c.CenterX-(-12)) < 1e-9 && math.Abs(c.CenterY) < 1e-9 {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected shape-aware left attachment at (-12, 0)")
	}
}
