package layout

// Collides reports whether two footprints overlap when each is expanded
// outward by gap/2 (minimum center-to-edge separation is gap).
func Collides(a, b Footprint, gap int) bool {
	g := float64(gap)

	switch {
	case a.IsRectLike() && b.IsRectLike():
		return rectRectCollide(a, b, g)
	case a.Shape == ShapeCircle && b.Shape == ShapeCircle:
		return circleCircleCollide(a, b, g)
	case a.IsRectLike() && b.Shape == ShapeCircle:
		return rectCircleCollide(a, b, g)
	case a.Shape == ShapeCircle && b.IsRectLike():
		return rectCircleCollide(b, a, g)
	default:
		// Ellipse and mixed non-circle cases use a conservative inflated AABB test.
		return rectRectCollide(a, b, g)
	}
}

func rectRectCollide(a, b Footprint, gap float64) bool {
	ax, ay, aw, ah := a.BBox()
	bx, by, bw, bh := b.BBox()

	aRight := ax + aw
	aBottom := ay + ah
	bRight := bx + bw
	bBottom := by + bh

	separatedX := aRight+gap <= bx || bRight+gap <= ax
	separatedY := aBottom+gap <= by || bBottom+gap <= ay

	return !(separatedX || separatedY)
}

func circleCircleCollide(a, b Footprint, gap float64) bool {
	dx := a.CenterX - b.CenterX
	dy := a.CenterY - b.CenterY
	distSq := dx*dx + dy*dy
	minDist := a.CircleRadius() + b.CircleRadius() + gap
	return distSq < minDist*minDist
}

func rectCircleCollide(rect, circle Footprint, gap float64) bool {
	rx, ry, rw, rh := rect.BBox()
	half := gap / 2
	rx -= half
	ry -= half
	rw += gap
	rh += gap

	r := circle.CircleRadius() + half
	cx, cy := circle.CenterX, circle.CenterY

	closestX := clamp(cx, rx, rx+rw)
	closestY := clamp(cy, ry, ry+rh)
	dx := cx - closestX
	dy := cy - closestY

	return dx*dx+dy*dy < r*r
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
