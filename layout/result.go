package layout

// Anchor identifies the layout origin on the centerpiece.
type Anchor struct {
	ItemID  string  `json:"itemId"`
	CenterX float64 `json:"centerX"`
	CenterY float64 `json:"centerY"`
}

// PlacedResult is a positioned frame in layout output.
// X and Y are the top-left of the frame bounding box; CenterX and CenterY are
// relative to the anchor. Direction is a compass label from the anchor (C, N, NE, …).
type PlacedResult struct {
	ID               string  `json:"id"`
	CenterX          float64 `json:"centerX"`
	CenterY          float64 `json:"centerY"`
	X                float64 `json:"x"`
	Y                float64 `json:"y"`
	Width            int     `json:"width"`
	Height           int     `json:"height"`
	Shape            Shape   `json:"shape"`
	OffsetFromAnchor float64  `json:"offsetFromAnchor"`
	Direction        string   `json:"direction"`
	AdjacentIDs      []string `json:"adjacentIds,omitempty"`
}

// Result is the output of [Layout].
type Result struct {
	Anchor Anchor         `json:"anchor"`
	Items  []PlacedResult `json:"items"`
	Bounds Bounds         `json:"bounds"`
}
