package layout

import "testing"

func footprint(shape Shape, h, w int, cx, cy float64) Footprint {
	return NewFootprint(Item{Shape: shape, Height: h, Width: w}, cx, cy)
}

func TestRectRectCollide(t *testing.T) {
	gap := 2
	rect := footprint(ShapeRectangle, 10, 10, 0, 0)

	tests := []struct {
		name      string
		otherCX   float64
		wantHit   bool
	}{
		{name: "overlapping", otherCX: 5, wantHit: true},
		{name: "touching at gap", otherCX: 12, wantHit: false},
		{name: "separated beyond gap", otherCX: 20, wantHit: false},
		{name: "too close", otherCX: 11, wantHit: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			other := footprint(ShapeRectangle, 10, 10, tt.otherCX, 0)
			if got := Collides(rect, other, gap); got != tt.wantHit {
				t.Fatalf("Collides() = %v, want %v", got, tt.wantHit)
			}
		})
	}
}

func TestCircleCircleCollide(t *testing.T) {
	gap := 2
	circle := footprint(ShapeCircle, 10, 10, 0, 0)

	tests := []struct {
		name    string
		otherCX float64
		wantHit bool
	}{
		{name: "overlapping", otherCX: 5, wantHit: true},
		{name: "touching at gap", otherCX: 12, wantHit: false},
		{name: "separated beyond gap", otherCX: 20, wantHit: false},
		{name: "too close", otherCX: 11, wantHit: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			other := footprint(ShapeCircle, 10, 10, tt.otherCX, 0)
			if got := Collides(circle, other, gap); got != tt.wantHit {
				t.Fatalf("Collides() = %v, want %v", got, tt.wantHit)
			}
		})
	}
}

func TestRectCircleCollide(t *testing.T) {
	gap := 2
	rect := footprint(ShapeRectangle, 10, 10, 0, 0)

	tests := []struct {
		name    string
		otherCX float64
		wantHit bool
	}{
		{name: "overlapping", otherCX: 5, wantHit: true},
		{name: "touching at gap", otherCX: 12, wantHit: false},
		{name: "separated beyond gap", otherCX: 20, wantHit: false},
		{name: "too close", otherCX: 11, wantHit: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			other := footprint(ShapeCircle, 10, 10, tt.otherCX, 0)
			if got := Collides(rect, other, gap); got != tt.wantHit {
				t.Fatalf("Collides() = %v, want %v", got, tt.wantHit)
			}
		})
	}

	t.Run("circle on left", func(t *testing.T) {
		leftCircle := footprint(ShapeCircle, 10, 10, -12, 0)
		if Collides(rect, leftCircle, gap) {
			t.Fatal("expected separated placement on left")
		}
	})
}

func TestEllipseCollideUsesInflatedBBox(t *testing.T) {
	gap := 2
	ellipse := footprint(ShapeEllipse, 12, 16, 0, 0)
	rect := footprint(ShapeRectangle, 12, 16, 0, 0)

	// Conservative ellipse test matches dimension-box behavior for identical bboxes.
	if !Collides(ellipse, rect, gap) {
		t.Fatal("identical centers should collide")
	}

	// width 16 → center distance 16+gap to touch without overlapping
	other := footprint(ShapeRectangle, 12, 16, 18, 0)
	if Collides(ellipse, other, gap) {
		t.Fatal("touching at gap should not collide")
	}
	if !Collides(ellipse, footprint(ShapeRectangle, 12, 16, 17, 0), gap) {
		t.Fatal("expected collision when too close")
	}

	t.Run("ellipse-ellipse", func(t *testing.T) {
		a := footprint(ShapeEllipse, 10, 14, 0, 0)
		b := footprint(ShapeEllipse, 10, 14, 16, 0) // width 14 + gap 2
		if Collides(a, b, gap) {
			t.Fatal("touching at gap should not collide")
		}
	})

	t.Run("ellipse-circle conservative", func(t *testing.T) {
		e := footprint(ShapeEllipse, 10, 10, 0, 0)
		c := footprint(ShapeCircle, 10, 10, 12, 0)
		if Collides(e, c, gap) {
			t.Fatal("touching at gap should not collide")
		}
	})
}

func TestSquareCollideSameAsRectangle(t *testing.T) {
	gap := 2
	square := footprint(ShapeSquare, 10, 10, 0, 0)
	rect := footprint(ShapeRectangle, 10, 10, 12, 0)

	if Collides(square, rect, gap) != Collides(
		footprint(ShapeRectangle, 10, 10, 0, 0),
		rect,
		gap,
	) {
		t.Fatal("square should collide like rectangle")
	}
}
