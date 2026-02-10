package layout

import "testing"

func TestNewPlacedItem(t *testing.T) {
	item := Item{ID: "main", Height: 20, Width: 10, Shape: ShapeRectangle, Centerpiece: true}
	placed := NewPlacedItem(item, 5, -3)

	if placed.Item.ID != "main" {
		t.Fatalf("item id = %q, want main", placed.Item.ID)
	}
	if placed.Footprint.CenterX != 5 || placed.Footprint.CenterY != -3 {
		t.Fatalf("center = (%v, %v), want (5, -3)", placed.Footprint.CenterX, placed.Footprint.CenterY)
	}
	if placed.Footprint.Height != 20 || placed.Footprint.Width != 10 {
		t.Fatalf("dims = %dx%d, want 20x10", placed.Footprint.Height, placed.Footprint.Width)
	}
}

func TestClusterBounds(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var c Cluster
		b := c.Bounds()
		if b != (Bounds{}) {
			t.Fatalf("empty bounds = %+v, want zero", b)
		}
	})

	t.Run("single item at origin", func(t *testing.T) {
		c := Cluster{NewPlacedItem(Item{
			ID: "main", Height: 20, Width: 40, Shape: ShapeRectangle, Centerpiece: true,
		}, 0, 0)}

		b := c.Bounds()
		assertFloat(t, "minX", b.MinX, -20)
		assertFloat(t, "minY", b.MinY, -10)
		assertFloat(t, "maxX", b.MaxX, 20)
		assertFloat(t, "maxY", b.MaxY, 10)
	})

	t.Run("two items", func(t *testing.T) {
		c := Cluster{
			NewPlacedItem(Item{ID: "main", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true}, 0, 0),
			NewPlacedItem(Item{ID: "right", Height: 8, Width: 8, Shape: ShapeSquare}, 12, 0),
		}

		b := c.Bounds()
		assertFloat(t, "minX", b.MinX, -5)
		assertFloat(t, "minY", b.MinY, -5)
		assertFloat(t, "maxX", b.MaxX, 16)
		assertFloat(t, "maxY", b.MaxY, 5)
	})
}
