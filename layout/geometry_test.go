package layout

import (
	"math"
	"testing"
)

func TestNewPlacedFootprint(t *testing.T) {
	item := Item{ID: "main", Height: 20, Width: 10, Shape: ShapeRectangle, Centerpiece: true}
	pf := NewPlacedFootprint(item, 5, -3)

	if pf.Item.ID != "main" {
		t.Fatalf("item id = %q, want main", pf.Item.ID)
	}
	if pf.Footprint.CenterX != 5 || pf.Footprint.CenterY != -3 {
		t.Fatalf("center = (%v, %v), want (5, -3)", pf.Footprint.CenterX, pf.Footprint.CenterY)
	}
	if pf.Footprint.Height != 20 || pf.Footprint.Width != 10 {
		t.Fatalf("dims = %dx%d, want 20x10", pf.Footprint.Height, pf.Footprint.Width)
	}
}

func TestFootprintBBoxRect(t *testing.T) {
	tests := []struct {
		name           string
		shape          Shape
		height, width  int
		centerX, centerY float64
		wantX, wantY, wantW, wantH float64
	}{
		{
			name: "rectangle centered at origin",
			shape: ShapeRectangle,
			height: 20, width: 40,
			centerX: 0, centerY: 0,
			wantX: -20, wantY: -10, wantW: 40, wantH: 20,
		},
		{
			name: "square offset",
			shape: ShapeSquare,
			height: 12, width: 12,
			centerX: 6, centerY: 4,
			wantX: 0, wantY: -2, wantW: 12, wantH: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := NewFootprint(Item{
				ID: "a", Height: tt.height, Width: tt.width, Shape: tt.shape,
			}, tt.centerX, tt.centerY)

			x, y, w, h := fp.BBox()
			assertFloat(t, "x", x, tt.wantX)
			assertFloat(t, "y", y, tt.wantY)
			assertFloat(t, "width", w, tt.wantW)
			assertFloat(t, "height", h, tt.wantH)
		})
	}
}

func TestFootprintCircleRadius(t *testing.T) {
	tests := []struct {
		name          string
		height, width int
		wantRadius    float64
	}{
		{name: "square bbox", height: 10, width: 10, wantRadius: 5},
		{name: "wider than tall", height: 10, width: 16, wantRadius: 5},
		{name: "taller than wide", height: 14, width: 8, wantRadius: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := NewFootprint(Item{
				ID: "c", Height: tt.height, Width: tt.width, Shape: ShapeCircle,
			}, 0, 0)
			assertFloat(t, "radius", fp.CircleRadius(), tt.wantRadius)
		})
	}
}

func TestFootprintEllipseSemiAxes(t *testing.T) {
	fp := NewFootprint(Item{
		ID: "e", Height: 12, Width: 16, Shape: ShapeEllipse,
	}, 0, 0)

	semiY, semiX := fp.EllipseSemiAxes()
	assertFloat(t, "semiY", semiY, 6)
	assertFloat(t, "semiX", semiX, 8)
}

func TestFootprintBBoxAllShapes(t *testing.T) {
	// Dimension box is the same for every shape at the same center and dims.
	item := Item{ID: "a", Height: 12, Width: 16}
	shapes := []Shape{ShapeRectangle, ShapeSquare, ShapeCircle, ShapeEllipse}

	var refX, refY, refW, refH float64
	for i, shape := range shapes {
		item.Shape = shape
		x, y, w, h := NewFootprint(item, 2, -1).BBox()
		if i == 0 {
			refX, refY, refW, refH = x, y, w, h
			continue
		}
		assertFloat(t, shape.String()+".x", x, refX)
		assertFloat(t, shape.String()+".y", y, refY)
		assertFloat(t, shape.String()+".w", w, refW)
		assertFloat(t, shape.String()+".h", h, refH)
	}
}

func assertFloat(t *testing.T, label string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("%s = %v, want %v", label, got, want)
	}
}
