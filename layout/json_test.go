package layout

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestValidationErrorIs(t *testing.T) {
	err := newDuplicateIDError("a")

	if !errors.Is(err, &ValidationError{Code: ValidationDuplicateID}) {
		t.Fatal("expected errors.Is match by code")
	}
	if errors.Is(err, &ValidationError{Code: ValidationEmptyInput}) {
		t.Fatal("unexpected code match")
	}
	if !errors.Is(err, &ValidationError{}) {
		t.Fatal("expected errors.Is match for any ValidationError")
	}
	if errors.Is(err, &PlacementError{}) {
		t.Fatal("validation error should not match placement error type")
	}
}

func TestPlacementErrorIs(t *testing.T) {
	err := newCannotPlaceError("stuck")

	if !errors.Is(err, &PlacementError{Code: PlacementCannotPlace}) {
		t.Fatal("expected errors.Is match by code")
	}
	if !errors.Is(err, &PlacementError{}) {
		t.Fatal("expected errors.Is match for any PlacementError")
	}
}

func TestParamsJSONRoundTrip(t *testing.T) {
	want := Params{
		Gap: 3, WallWidth: 120, WallHeight: 96,
		Items: []Item{
			{ID: "main", Height: 12, Width: 10, Shape: ShapeRectangle, Centerpiece: true},
			{ID: "a", Height: 8, Width: 6, Shape: ShapeCircle},
		},
	}

	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got Params
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Gap != want.Gap || len(got.Items) != len(want.Items) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if got.WallWidth != want.WallWidth || got.WallHeight != want.WallHeight {
		t.Fatalf("wall = %.0fx%.0f, want %.0fx%.0f", got.WallWidth, got.WallHeight, want.WallWidth, want.WallHeight)
	}
	if got.Items[1].Shape != ShapeCircle {
		t.Fatalf("shape = %q, want circle", got.Items[1].Shape)
	}
}

func TestResultJSONRoundTrip(t *testing.T) {
	want, err := Layout(Params{
		Gap: 2,
		Items: []Item{
			{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
			{ID: "side", Height: 8, Width: 6, Shape: ShapeRectangle},
		},
	})
	if err != nil {
		t.Fatalf("layout: %v", err)
	}

	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got Result
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Anchor.ItemID != want.Anchor.ItemID {
		t.Fatalf("anchor = %+v, want %+v", got.Anchor, want.Anchor)
	}
	if len(got.Items) != len(want.Items) {
		t.Fatalf("item count = %d, want %d", len(got.Items), len(want.Items))
	}
	if got.Bounds != want.Bounds {
		t.Fatalf("bounds = %+v, want %+v", got.Bounds, want.Bounds)
	}

	side := got.Items[1]
	if side.ID != "side" || side.Direction == "" || len(side.AdjacentIDs) == 0 {
		t.Fatalf("side = %+v, want id, direction, and neighbors", side)
	}
}

func TestPlacementErrorJSONRoundTrip(t *testing.T) {
	err := newCannotPlaceError("stuck")
	data, errMarshal := json.Marshal(err)
	if errMarshal != nil {
		t.Fatalf("marshal: %v", errMarshal)
	}

	var decoded PlacementError
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.Code != PlacementCannotPlace || decoded.ItemID != "stuck" {
		t.Fatalf("decoded = %+v", decoded)
	}
}
