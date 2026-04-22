package layout

import (
	"fmt"
	"math"
	"testing"
)

func TestLayoutSixMixedFormsBlob(t *testing.T) {
	result, err := Layout(Params{
		Gap: 2,
		Items: []Item{
			{ID: "main", Height: 36, Width: 24, Shape: ShapeRectangle, Centerpiece: true},
			{ID: "left", Height: 14, Width: 11, Shape: ShapeRectangle},
			{ID: "right", Height: 14, Width: 11, Shape: ShapeRectangle},
			{ID: "top", Height: 10, Width: 10, Shape: ShapeCircle},
			{ID: "bottom", Height: 12, Width: 16, Shape: ShapeEllipse},
			{ID: "tl", Height: 8, Width: 8, Shape: ShapeSquare},
			{ID: "br", Height: 8, Width: 8, Shape: ShapeSquare},
		},
	})
	if err != nil {
		t.Fatalf("layout error: %v", err)
	}

	if !ClusterHas2DSpread(result) {
		t.Fatal("expected 2D blob spread for mixed sizes")
	}

	aspect := ClusterAspectRatio(result.Bounds)
	if aspect < 0.9 || aspect > 2.2 {
		t.Fatalf("aspect %v outside wide-round band", aspect)
	}

	nonZeroY := 0
	for _, item := range result.Items {
		if math.Abs(item.CenterY) > 1 {
			nonZeroY++
		}
	}
	if nonZeroY < 2 {
		t.Fatalf("expected multiple frames off the horizontal axis, got %d", nonZeroY)
	}

	diagonal := CountDiagonalPlacements(result)
	if diagonal < 2 {
		t.Fatalf("expected corner/diagonal fills, got %d diagonal placements", diagonal)
	}
}

func TestLayoutSixSmallVsLargeSameShape(t *testing.T) {
	gap := 2

	smallItems := []Item{{ID: "main", Height: 8, Width: 8, Shape: ShapeSquare, Centerpiece: true}}
	for i := 0; i < 6; i++ {
		smallItems = append(smallItems, Item{
			ID: fmt.Sprintf("s%d", i), Height: 8, Width: 8, Shape: ShapeSquare,
		})
	}

	largeItems := []Item{{ID: "main", Height: 24, Width: 18, Shape: ShapeRectangle, Centerpiece: true}}
	for i := 0; i < 6; i++ {
		largeItems = append(largeItems, Item{
			ID: fmt.Sprintf("L%d", i), Height: 24, Width: 18, Shape: ShapeRectangle,
		})
	}

	small, err := Layout(Params{Gap: gap, Items: smallItems})
	if err != nil {
		t.Fatalf("small layout: %v", err)
	}
	large, err := Layout(Params{Gap: gap, Items: largeItems})
	if err != nil {
		t.Fatalf("large layout: %v", err)
	}

	if !ClusterHas2DSpread(small) || !ClusterHas2DSpread(large) {
		t.Fatal("expected both small and large sets to form 2D blobs")
	}

	if CountDiagonalPlacements(small) < 2 || CountDiagonalPlacements(large) < 2 {
		t.Fatal("expected corner fills in same-size sets")
	}

	smallAspect := ClusterAspectRatio(small.Bounds)
	largeAspect := ClusterAspectRatio(large.Bounds)
	if math.Abs(smallAspect-largeAspect) > 0.6 {
		t.Fatalf("aspect mismatch: small=%.2f large=%.2f; same-size sets should scale similarly",
			smallAspect, largeAspect)
	}
}

func TestScoreDiagonalPreferredOverAxis(t *testing.T) {
	placed := Cluster{anchorCenterpiece(Item{
		ID: "main", Height: 20, Width: 20, Shape: ShapeSquare, Centerpiece: true,
	})}
	item := Item{ID: "side", Height: 10, Width: 10, Shape: ShapeSquare}

	axis := Candidate{Item: item, CenterX: 17, CenterY: 0}
	diag := Candidate{Item: item, CenterX: 17, CenterY: 12}

	if ScoreCandidate(axis, placed, 2) <= ScoreCandidate(diag, placed, 2) {
		t.Fatal("expected diagonal candidate to score better than pure axis extension")
	}
}

func TestScoreCollinearityPenalizesCombs(t *testing.T) {
	comb := Cluster{
		anchorCenterpiece(Item{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true}),
		NewPlacedItem(Item{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare}, 12, 0),
		NewPlacedItem(Item{ID: "b", Height: 10, Width: 10, Shape: ShapeSquare}, 24, 0),
		NewPlacedItem(Item{ID: "c", Height: 10, Width: 10, Shape: ShapeSquare}, 36, 0),
	}
	blob := Cluster{
		anchorCenterpiece(Item{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true}),
		NewPlacedItem(Item{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare}, 12, 0),
		NewPlacedItem(Item{ID: "b", Height: 10, Width: 10, Shape: ShapeSquare}, 0, 12),
		NewPlacedItem(Item{ID: "c", Height: 10, Width: 10, Shape: ShapeSquare}, -12, 0),
	}

	if scoreCollinearity(comb) <= scoreCollinearity(blob) {
		t.Fatalf("comb=%v blob=%v; want comb to score higher penalty",
			scoreCollinearity(comb), scoreCollinearity(blob))
	}
}
