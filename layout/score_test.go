package layout

import "testing"

func TestScoreCompactness(t *testing.T) {
	near := Candidate{Item: Item{ID: "a"}, CenterX: 12, CenterY: 0}
	far := Candidate{Item: Item{ID: "a"}, CenterX: 30, CenterY: 0}

	if scoreCompactness(near) >= scoreCompactness(far) {
		t.Fatalf("near = %v, far = %v; want nearer to score lower", scoreCompactness(near), scoreCompactness(far))
	}
}

func TestScoreBalance(t *testing.T) {
	centerpiece := anchorCenterpiece(Item{
		ID: "main", Height: 20, Width: 20, Shape: ShapeRectangle, Centerpiece: true,
	})
	item := Item{ID: "sat", Height: 10, Width: 10, Shape: ShapeSquare}

	balanced := Cluster{
		centerpiece,
		NewPlacedItem(item, -12, -12),
		NewPlacedItem(item, 12, -12),
		NewPlacedItem(item, 12, 12),
		NewPlacedItem(item, -12, 12),
	}

	lopsided := Cluster{
		centerpiece,
		NewPlacedItem(item, 12, -12),
		NewPlacedItem(item, 12, 12),
		NewPlacedItem(item, 24, -12),
		NewPlacedItem(item, 24, 12),
	}

	if scoreBalance(balanced) >= scoreBalance(lopsided) {
		t.Fatalf("balanced = %v, lopsided = %v; want balanced to score lower", scoreBalance(balanced), scoreBalance(lopsided))
	}
}

func TestScoreBlobSmoothness(t *testing.T) {
	ideal := Cluster{
		anchorCenterpiece(Item{ID: "main", Height: 10, Width: 14, Shape: ShapeRectangle, Centerpiece: true}),
	}
	tall := Cluster{
		anchorCenterpiece(Item{ID: "main", Height: 40, Width: 10, Shape: ShapeRectangle, Centerpiece: true}),
	}

	if scoreBlobSmoothness(ideal) >= scoreBlobSmoothness(tall) {
		t.Fatalf("ideal = %v, tall = %v; want wide-round cluster closer to target aspect",
			scoreBlobSmoothness(ideal), scoreBlobSmoothness(tall))
	}
}

func TestScoreConcavity(t *testing.T) {
	single := Cluster{
		anchorCenterpiece(Item{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true}),
	}

	lShape := Cluster{
		anchorCenterpiece(Item{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true}),
		NewPlacedItem(Item{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare}, 12, 0),
		NewPlacedItem(Item{ID: "b", Height: 10, Width: 10, Shape: ShapeSquare}, 0, 12),
	}

	if scoreConcavity(single) >= scoreConcavity(lShape) {
		t.Fatalf("single = %v, lShape = %v; want tighter silhouette to score lower",
			scoreConcavity(single), scoreConcavity(lShape))
	}
}

func TestScoreLocalContinuity(t *testing.T) {
	gap := 2
	placed := Cluster{
		anchorCenterpiece(Item{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true}),
		NewPlacedItem(Item{ID: "right", Height: 10, Width: 10, Shape: ShapeSquare}, 12, 0),
		NewPlacedItem(Item{ID: "below", Height: 10, Width: 10, Shape: ShapeSquare}, 0, 12),
	}

	oneNeighbor := Candidate{Item: Item{ID: "new", Height: 10, Width: 10, Shape: ShapeSquare}, CenterX: 24, CenterY: 0}
	twoNeighbors := Candidate{Item: Item{ID: "new", Height: 10, Width: 10, Shape: ShapeSquare}, CenterX: 12, CenterY: 12}

	if scoreLocalContinuity(twoNeighbors, placed, gap) <= scoreLocalContinuity(oneNeighbor, placed, gap) {
		t.Fatalf("two-neighbor = %v, one-neighbor = %v; want more continuity to score higher reward",
			scoreLocalContinuity(twoNeighbors, placed, gap),
			scoreLocalContinuity(oneNeighbor, placed, gap),
		)
	}
}

func TestBestCandidate(t *testing.T) {
	gap := 2
	placed := Cluster{anchorCenterpiece(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true,
	})}
	item := Item{ID: "side", Height: 10, Width: 10, Shape: ShapeSquare}

	candidates := []Candidate{
		{Item: item, CenterX: 30, CenterY: 0},
		{Item: item, CenterX: 12, CenterY: 0},
		{Item: item, CenterX: -12, CenterY: 0},
	}

	best, ok := BestCandidate(candidates, placed, gap)
	if !ok {
		t.Fatal("expected best candidate")
	}
	if best.CenterX != 12 && best.CenterX != -12 {
		t.Fatalf("best centerX = %v, want 12 or -12 (compact balanced attach)", best.CenterX)
	}
}

func TestBestCandidateTiebreakByPosition(t *testing.T) {
	placed := Cluster{anchorCenterpiece(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true,
	})}
	item := Item{ID: "side", Height: 10, Width: 10, Shape: ShapeSquare}

	symmetric := []Candidate{
		{Item: item, CenterX: 12, CenterY: 0},
		{Item: item, CenterX: -12, CenterY: 0},
	}

	first, _ := BestCandidate(symmetric, placed, 2)
	second, _ := BestCandidate([]Candidate{symmetric[1], symmetric[0]}, placed, 2)

	if first.CenterX != second.CenterX {
		t.Fatalf("tiebreak unstable: got %v and %v", first.CenterX, second.CenterX)
	}
	if first.CenterX != -12 {
		t.Fatalf("tiebreak centerX = %v, want -12", first.CenterX)
	}
}

func TestScoreCandidateLowerIsBetter(t *testing.T) {
	placed := Cluster{anchorCenterpiece(Item{
		ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true,
	})}
	item := Item{ID: "side", Height: 10, Width: 10, Shape: ShapeSquare}

	near := Candidate{Item: item, CenterX: 12, CenterY: 0}
	far := Candidate{Item: item, CenterX: 30, CenterY: 0}

	if ScoreCandidate(near, placed, 2) >= ScoreCandidate(far, placed, 2) {
		t.Fatal("expected nearer candidate to have lower total score")
	}
}
