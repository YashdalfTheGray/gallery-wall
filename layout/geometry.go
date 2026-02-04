package layout

// Footprint is the shape-aware geometry of a frame at a position.
// CenterX and CenterY are the center of the dimension bounding box.
type Footprint struct {
	Shape   Shape
	Height  int
	Width   int
	CenterX float64
	CenterY float64
}

// PlacedFootprint is an alias for PlacedItem used in geometry helpers.
type PlacedFootprint = PlacedItem

// NewPlacedFootprint builds a placed footprint for an item at a center.
func NewPlacedFootprint(item Item, centerX, centerY float64) PlacedFootprint {
	return NewPlacedItem(item, centerX, centerY)
}

// NewFootprint builds a footprint from item dimensions and center.
func NewFootprint(item Item, centerX, centerY float64) Footprint {
	return Footprint{
		Shape:   item.Shape,
		Height:  item.Height,
		Width:   item.Width,
		CenterX: centerX,
		CenterY: centerY,
	}
}

// BBox returns the top-left corner and size of the item's dimension box.
func (f Footprint) BBox() (x, y, width, height float64) {
	halfW := float64(f.Width) / 2
	halfH := float64(f.Height) / 2
	return f.CenterX - halfW, f.CenterY - halfH, float64(f.Width), float64(f.Height)
}

// CircleRadius returns the inscribed circle radius for circle footprints.
func (f Footprint) CircleRadius() float64 {
	minSide := f.Height
	if f.Width < minSide {
		minSide = f.Width
	}
	return float64(minSide) / 2
}

// EllipseSemiAxes returns the vertical and horizontal semi-axes for ellipse footprints.
func (f Footprint) EllipseSemiAxes() (semiY, semiX float64) {
	return float64(f.Height) / 2, float64(f.Width) / 2
}

// IsRectLike reports whether the footprint uses rectangular collision geometry.
func (f Footprint) IsRectLike() bool {
	return f.Shape == ShapeSquare || f.Shape == ShapeRectangle
}

// HalfExtents returns shape-aware horizontal and vertical half-extents from center.
func (f Footprint) HalfExtents() (halfX, halfY float64) {
	switch f.Shape {
	case ShapeCircle:
		r := f.CircleRadius()
		return r, r
	case ShapeEllipse:
		semiY, semiX := f.EllipseSemiAxes()
		return semiX, semiY
	default:
		return float64(f.Width) / 2, float64(f.Height) / 2
	}
}

// HalfExtents returns shape-aware half-extents for an item at the origin.
func halfExtents(item Item) (halfX, halfY float64) {
	return NewFootprint(item, 0, 0).HalfExtents()
}

