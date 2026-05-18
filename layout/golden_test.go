package layout

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestGoldenCenterpieceOnly(t *testing.T) {
	runGoldenTest(t, "centerpiece_only")
}

func TestGoldenCenterpieceTwoSymmetric(t *testing.T) {
	runGoldenTest(t, "centerpiece_two_symmetric")
}

func TestGoldenCenterpieceSixMixed(t *testing.T) {
	runGoldenTest(t, "centerpiece_six_mixed")
}

func TestGoldenGapSensitivity(t *testing.T) {
	items := []Item{
		{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
		{ID: "side", Height: 10, Width: 10, Shape: ShapeSquare},
	}

	close, err := Layout(Params{Gap: 2, Items: items})
	if err != nil {
		t.Fatalf("gap 2 layout: %v", err)
	}
	wide, err := Layout(Params{Gap: 6, Items: items})
	if err != nil {
		t.Fatalf("gap 6 layout: %v", err)
	}

	closeSide := itemByID(close, "side")
	wideSide := itemByID(wide, "side")

	closeDist := math.Hypot(closeSide.CenterX, closeSide.CenterY)
	wideDist := math.Hypot(wideSide.CenterX, wideSide.CenterY)
	if wideDist <= closeDist {
		t.Fatalf("gap 6 distance %v should exceed gap 2 distance %v", wideDist, closeDist)
	}
}

func runGoldenTest(t *testing.T, name string) {
	t.Helper()

	params := loadGoldenParams(t, name)
	want := loadGoldenResult(t, name)

	got, err := Layout(params)
	if err != nil {
		t.Fatalf("layout error: %v", err)
	}

	assertResultsEqual(t, got, want)
}

func loadGoldenParams(t *testing.T, name string) Params {
	t.Helper()
	path := filepath.Join("testdata", name+"_params.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read params: %v", err)
	}
	var params Params
	if err := json.Unmarshal(data, &params); err != nil {
		t.Fatalf("unmarshal params: %v", err)
	}
	return params
}

func loadGoldenResult(t *testing.T, name string) Result {
	t.Helper()
	path := filepath.Join("testdata", name+"_result.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	return result
}

func assertResultsEqual(t *testing.T, got, want Result) {
	t.Helper()

	if got.Anchor != want.Anchor {
		t.Fatalf("anchor: got %+v, want %+v", got.Anchor, want.Anchor)
	}
	if got.Bounds != want.Bounds {
		t.Fatalf("bounds: got %+v, want %+v", got.Bounds, want.Bounds)
	}
	if len(got.Items) != len(want.Items) {
		t.Fatalf("item count: got %d, want %d", len(got.Items), len(want.Items))
	}

	wantByID := make(map[string]PlacedResult, len(want.Items))
	for _, item := range want.Items {
		wantByID[item.ID] = item
	}

	for _, item := range got.Items {
		expected, ok := wantByID[item.ID]
		if !ok {
			t.Fatalf("unexpected item %q", item.ID)
		}
		assertPlacedResultEqual(t, item.ID, item, expected)
	}
}

func assertPlacedResultEqual(t *testing.T, id string, got, want PlacedResult) {
	t.Helper()

	if got.ID != want.ID || got.Shape != want.Shape || got.Width != want.Width || got.Height != want.Height {
		t.Fatalf("%s identity: got %+v, want %+v", id, got, want)
	}
	assertFloat(t, id+".centerX", got.CenterX, want.CenterX)
	assertFloat(t, id+".centerY", got.CenterY, want.CenterY)
	assertFloat(t, id+".x", got.X, want.X)
	assertFloat(t, id+".y", got.Y, want.Y)
	assertFloat(t, id+".offsetFromAnchor", got.OffsetFromAnchor, want.OffsetFromAnchor)
	if got.Direction != want.Direction {
		t.Fatalf("%s direction: got %q, want %q", id, got.Direction, want.Direction)
	}
	if len(got.AdjacentIDs) != len(want.AdjacentIDs) {
		t.Fatalf("%s adjacentIds: got %v, want %v", id, got.AdjacentIDs, want.AdjacentIDs)
	}
	for i := range got.AdjacentIDs {
		if got.AdjacentIDs[i] != want.AdjacentIDs[i] {
			t.Fatalf("%s adjacentIds: got %v, want %v", id, got.AdjacentIDs, want.AdjacentIDs)
		}
	}
}

func itemByID(result Result, id string) PlacedResult {
	for _, item := range result.Items {
		if item.ID == id {
			return item
		}
	}
	panic("item not found: " + id)
}
