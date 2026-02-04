package layout

// PlacedItem is a frame positioned on the wall with a cached footprint.
type PlacedItem struct {
	Item      Item
	Footprint Footprint
}

// NewPlacedItem places an item at the given center.
func NewPlacedItem(item Item, centerX, centerY float64) PlacedItem {
	return PlacedItem{
		Item:      item,
		Footprint: NewFootprint(item, centerX, centerY),
	}
}

// Cluster is the growing set of placed frames.
type Cluster []PlacedItem

// Bounds is an axis-aligned bounding box in centerpiece-relative coordinates.
type Bounds struct {
	MinX float64 `json:"minX"`
	MinY float64 `json:"minY"`
	MaxX float64 `json:"maxX"`
	MaxY float64 `json:"maxY"`
}

// Bounds returns the union bounding box of all placed item dimension boxes.
func (c Cluster) Bounds() Bounds {
	if len(c) == 0 {
		return Bounds{}
	}

	var b Bounds
	for i, placed := range c {
		x, y, w, h := placed.Footprint.BBox()
		right := x + w
		bottom := y + h

		if i == 0 {
			b = Bounds{MinX: x, MinY: y, MaxX: right, MaxY: bottom}
			continue
		}
		if x < b.MinX {
			b.MinX = x
		}
		if y < b.MinY {
			b.MinY = y
		}
		if right > b.MaxX {
			b.MaxX = right
		}
		if bottom > b.MaxY {
			b.MaxY = bottom
		}
	}
	return b
}
