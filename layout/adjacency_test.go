package layout

import "testing"

func TestAdjacentRectRect(t *testing.T) {
	gap := 2
	rect := footprint(ShapeRectangle, 10, 10, 0, 0)

	tests := []struct {
		name    string
		otherCX float64
		want    bool
	}{
		{name: "overlapping", otherCX: 5, want: false},
		{name: "touching at gap", otherCX: 12, want: true},
		{name: "beyond gap", otherCX: 13, want: false},
		{name: "well separated", otherCX: 20, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			other := footprint(ShapeRectangle, 10, 10, tt.otherCX, 0)
			if got := Adjacent(rect, other, gap); got != tt.want {
				t.Fatalf("Adjacent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdjacentCircleCircle(t *testing.T) {
	gap := 2
	circle := footprint(ShapeCircle, 10, 10, 0, 0)

	tests := []struct {
		name    string
		otherCX float64
		want    bool
	}{
		{name: "overlapping", otherCX: 5, want: false},
		{name: "touching at gap", otherCX: 12, want: true},
		{name: "beyond gap", otherCX: 13, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			other := footprint(ShapeCircle, 10, 10, tt.otherCX, 0)
			if got := Adjacent(circle, other, gap); got != tt.want {
				t.Fatalf("Adjacent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdjacentMixedShapes(t *testing.T) {
	gap := 2
	rect := footprint(ShapeRectangle, 10, 10, 0, 0)

	t.Run("rect and circle touching at gap", func(t *testing.T) {
		circle := footprint(ShapeCircle, 10, 10, 12, 0)
		if !Adjacent(rect, circle, gap) {
			t.Fatal("expected adjacent")
		}
		if Collides(rect, circle, gap) {
			t.Fatal("touching at gap should not collide")
		}
	})

	t.Run("rect and circle beyond gap", func(t *testing.T) {
		circle := footprint(ShapeCircle, 10, 10, 13, 0)
		if Adjacent(rect, circle, gap) {
			t.Fatal("expected not adjacent")
		}
	})

	t.Run("ellipse and rect touching at gap", func(t *testing.T) {
		ellipse := footprint(ShapeEllipse, 10, 10, 0, 0)
		other := footprint(ShapeRectangle, 10, 10, 12, 0)
		if !Adjacent(ellipse, other, gap) {
			t.Fatal("expected adjacent")
		}
	})

	t.Run("ellipse and circle beyond gap", func(t *testing.T) {
		ellipse := footprint(ShapeEllipse, 10, 10, 0, 0)
		circle := footprint(ShapeCircle, 10, 10, 13, 0)
		if Adjacent(ellipse, circle, gap) {
			t.Fatal("expected not adjacent")
		}
	})
}

func TestAdjacentConsistentWithCollides(t *testing.T) {
	gap := 2
	pairs := []struct {
		a, b Footprint
	}{
		{footprint(ShapeRectangle, 10, 10, 0, 0), footprint(ShapeRectangle, 10, 10, 12, 0)},
		{footprint(ShapeCircle, 10, 10, 0, 0), footprint(ShapeCircle, 10, 10, 12, 0)},
		{footprint(ShapeRectangle, 10, 10, 0, 0), footprint(ShapeCircle, 10, 10, 12, 0)},
	}

	for i, pair := range pairs {
		if Collides(pair.a, pair.b, gap) {
			t.Fatalf("pair %d: expected no collision at gap boundary", i)
		}
		if !Adjacent(pair.a, pair.b, gap) {
			t.Fatalf("pair %d: expected adjacent at gap boundary", i)
		}
	}
}
