package layout

import "math"

// Adjacent reports whether two footprints do not overlap and are separated
// by at most gap (touching within the allowed margin counts as connected).
func Adjacent(a, b Footprint, gap int) bool {
	if overlaps(a, b) {
		return false
	}
	return minSeparation(a, b) <= float64(gap)
}

func overlaps(a, b Footprint) bool {
	return minSeparation(a, b) < 0
}

func minSeparation(a, b Footprint) float64 {
	switch {
	case a.IsRectLike() && b.IsRectLike():
		return rectRectMinSep(a, b)
	case a.Shape == ShapeCircle && b.Shape == ShapeCircle:
		return circleCircleMinSep(a, b)
	case a.IsRectLike() && b.Shape == ShapeCircle:
		return rectCircleMinSep(a, b)
	case a.Shape == ShapeCircle && b.IsRectLike():
		return rectCircleMinSep(b, a)
	default:
		return rectRectMinSep(a, b)
	}
}

func rectRectMinSep(a, b Footprint) float64 {
	ax, ay, aw, ah := a.BBox()
	bx, by, bw, bh := b.BBox()

	aRight := ax + aw
	aBottom := ay + ah
	bRight := bx + bw
	bBottom := by + bh

	overlapX := minFloat(aRight, bRight) - maxFloat(ax, bx)
	overlapY := minFloat(aBottom, bBottom) - maxFloat(ay, by)

	if overlapX > 0 && overlapY > 0 {
		return -minFloat(overlapX, overlapY)
	}

	gapX := maxFloat(ax-bRight, bx-aRight)
	gapY := maxFloat(ay-bBottom, by-aBottom)

	if gapX > 0 && gapY > 0 {
		return math.Hypot(gapX, gapY)
	}
	if gapX > 0 {
		return gapX
	}
	return gapY
}

func circleCircleMinSep(a, b Footprint) float64 {
	dx := a.CenterX - b.CenterX
	dy := a.CenterY - b.CenterY
	return math.Hypot(dx, dy) - a.CircleRadius() - b.CircleRadius()
}

func rectCircleMinSep(rect, circle Footprint) float64 {
	rx, ry, rw, rh := rect.BBox()
	cx, cy := circle.CenterX, circle.CenterY
	r := circle.CircleRadius()

	closestX := clamp(cx, rx, rx+rw)
	closestY := clamp(cy, ry, ry+rh)
	dx := cx - closestX
	dy := cy - closestY

	return math.Hypot(dx, dy) - r
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
