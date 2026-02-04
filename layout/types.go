package layout

// Shape describes the collision footprint and rendering geometry for a frame.
// Dimensions (height, width) are always the source of truth; Shape selects
// how those dimensions are interpreted geometrically.
type Shape string

const (
	ShapeSquare    Shape = "square"
	ShapeRectangle Shape = "rectangle"
	ShapeCircle    Shape = "circle"
	ShapeEllipse   Shape = "ellipse"
)

// String returns the JSON/API representation of the shape.
func (s Shape) String() string {
	return string(s)
}

// Item describes a single frame to place on the gallery wall.
type Item struct {
	ID          string `json:"id"`
	Height      int    `json:"height"`
	Width       int    `json:"width"`
	Shape       Shape  `json:"shape"`
	Centerpiece bool   `json:"centerpiece"`
}

// Params are the inputs to [Layout].
// Gap is the minimum space between adjacent frames. WallWidth and WallHeight
// optionally constrain the cluster to a rectangle centered on the anchor; omit
// or set to zero for no wall limit.
type Params struct {
	Gap        int     `json:"gap"`
	WallWidth  float64 `json:"wallWidth,omitempty"`
	WallHeight float64 `json:"wallHeight,omitempty"`
	Items      []Item  `json:"items"`
}
