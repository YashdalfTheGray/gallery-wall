package layout

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParamsWallBounds(t *testing.T) {
	p := Params{WallWidth: 100, WallHeight: 60}
	wall, ok := p.WallBounds()
	if !ok {
		t.Fatal("expected wall bounds")
	}
	if wall.MinX != -50 || wall.MaxX != 50 || wall.MinY != -30 || wall.MaxY != 30 {
		t.Fatalf("wall = %+v", wall)
	}

	unbounded := Params{Gap: 2}
	if _, ok := unbounded.WallBounds(); ok {
		t.Fatal("expected no wall when dimensions omitted")
	}
}

func TestValidateWallRequiresBothDimensions(t *testing.T) {
	err := Validate(Params{
		WallWidth: 100,
		Items: []Item{
			{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
		},
	})
	if !IsValidationCode(err, ValidationInvalidWall) {
		t.Fatalf("got %v, want invalid_wall", err)
	}
}

func TestValidateCenterpieceExceedsWall(t *testing.T) {
	err := Validate(Params{
		WallWidth: 10, WallHeight: 10,
		Items: []Item{
			{ID: "main", Height: 20, Width: 20, Shape: ShapeSquare, Centerpiece: true},
		},
	})
	if !IsValidationCode(err, ValidationCenterpieceExceedsWall) {
		t.Fatalf("got %v, want centerpiece_exceeds_wall", err)
	}
}

func TestLayoutRespectsWallBounds(t *testing.T) {
	result, err := Layout(Params{
		Gap: 2, WallWidth: 40, WallHeight: 40,
		Items: []Item{
			{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
			{ID: "a", Height: 8, Width: 8, Shape: ShapeSquare},
			{ID: "b", Height: 8, Width: 8, Shape: ShapeSquare},
			{ID: "c", Height: 8, Width: 8, Shape: ShapeSquare},
			{ID: "d", Height: 8, Width: 8, Shape: ShapeSquare},
		},
	})
	if err != nil {
		t.Fatalf("layout: %v", err)
	}

	wall, _ := Params{WallWidth: 40, WallHeight: 40}.WallBounds()
	for _, item := range result.Items {
		fp := NewFootprint(Item{
			ID: item.ID, Height: item.Height, Width: item.Width, Shape: item.Shape,
		}, item.CenterX, item.CenterY)
		if !footprintFitsWall(fp, wall) {
			t.Fatalf("%s at (%.1f, %.1f) exceeds wall %+v", item.ID, item.CenterX, item.CenterY, wall)
		}
	}
}

func TestLayoutCannotPlaceWhenWallTooSmall(t *testing.T) {
	_, err := Layout(Params{
		Gap: 2, WallWidth: 24, WallHeight: 24,
		Items: []Item{
			{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
			{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare},
			{ID: "b", Height: 10, Width: 10, Shape: ShapeSquare},
			{ID: "c", Height: 10, Width: 10, Shape: ShapeSquare},
			{ID: "d", Height: 10, Width: 10, Shape: ShapeSquare},
			{ID: "e", Height: 10, Width: 10, Shape: ShapeSquare},
		},
	})
	if !IsPlacementCode(err, PlacementCannotPlace) {
		t.Fatalf("got %v, want cannot_place_item", err)
	}
}

func TestLayoutTwentyFiveWithinRecordedBounds(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "twentyfive_params.json"))
	if err != nil {
		t.Fatalf("read params: %v", err)
	}
	var params Params
	if err := json.Unmarshal(data, &params); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	params.WallWidth = 54
	params.WallHeight = 60

	result, err := Layout(params)
	if err != nil {
		t.Fatalf("layout: %v", err)
	}
	if len(result.Items) != 25 {
		t.Fatalf("placed %d items, want 25", len(result.Items))
	}
}

func TestGenerateCandidatesRejectsOutsideWall(t *testing.T) {
	placed := Cluster{anchorCenterpiece(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true,
	})}
	item := Item{ID: "side", Height: 10, Width: 10, Shape: ShapeSquare}
	wall := Bounds{MinX: -12, MinY: -12, MaxX: 12, MaxY: 12}

	cands := GenerateCandidates(item, placed, 2, &wall)
	for _, c := range cands {
		fp := c.footprint()
		if !footprintFitsWall(fp, wall) {
			t.Fatalf("candidate at (%.0f, %.0f) outside wall", c.CenterX, c.CenterY)
		}
	}

	far := Candidate{Item: item, CenterX: 30, CenterY: 0}
	if isValidCandidate(far, placed, 2, &wall) {
		t.Fatal("expected far candidate to be rejected by wall")
	}
}

func TestFootprintFitsWallEdge(t *testing.T) {
	wall := Bounds{MinX: -10, MinY: -10, MaxX: 10, MaxY: 10}
	fp := NewFootprint(Item{Height: 10, Width: 10, Shape: ShapeSquare}, 0, 0)
	if !footprintFitsWall(fp, wall) {
		t.Fatal("centered 10x10 should fit in 20x20 wall")
	}
	fpEdge := NewFootprint(Item{Height: 10, Width: 10, Shape: ShapeSquare}, 10, 0)
	if footprintFitsWall(fpEdge, wall) {
		t.Fatal("shifted frame should not fit")
	}
}
