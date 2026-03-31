package layout

import (
	"math"
	"sort"
)

const (
	weightCompactness     = 0.35
	weightBalance         = 0.6
	weightBlobSmoothness  = 1.8
	weightConcavity       = 0.5
	weightLocalContinuity = 0.9
	weightCollinearity    = 1.5
	weightAxisAttachment  = 1.0
	weightAxisOnly        = 1.5
	weightDiagonalFill    = 2.0
	weightQuadrantFill    = 1.5
	weightClusterArea     = 0.25
	targetAspectRatio     = 1.4
	maxBlobAspectRatio    = 2.2
	axisEpsilon           = 1.0
)

// ScoreCandidate returns a lower-is-better score for placing cand on placed cluster.
// Exported for tests and algorithm tuning; [Layout] uses this internally.
func ScoreCandidate(cand Candidate, placed Cluster, gap int) float64 {
	cluster := clusterWithCandidate(placed, cand)

	compactness := scoreCompactness(cand)
	balance := scoreBalance(cluster)
	blob := scoreBlobSmoothness(cluster)
	concavity := scoreConcavity(cluster)
	continuity := scoreLocalContinuity(cand, placed, gap)
	collinearity := scoreCollinearity(cluster)
	axisAttach := scoreAxisAttachment(cand, cluster)
	axisOnly := scoreAxisOnlyPlacement(cand, cluster)
	diagonal := scoreDiagonalFill(cand)
	quadrant := scoreQuadrantFill(cand, cluster)
	clusterArea := scoreClusterArea(cluster)

	return weightCompactness*compactness +
		weightBalance*balance +
		weightBlobSmoothness*blob +
		weightConcavity*concavity +
		weightCollinearity*collinearity +
		weightAxisAttachment*axisAttach +
		weightAxisOnly*axisOnly +
		weightClusterArea*clusterArea -
		weightLocalContinuity*continuity -
		weightDiagonalFill*diagonal -
		weightQuadrantFill*quadrant
}

// BestCandidate selects the lowest-scoring valid candidate.
func BestCandidate(candidates []Candidate, placed Cluster, gap int) (Candidate, bool) {
	if len(candidates) == 0 {
		return Candidate{}, false
	}

	best := candidates[0]
	bestScore := ScoreCandidate(best, placed, gap)

	for _, cand := range candidates[1:] {
		score := ScoreCandidate(cand, placed, gap)
		if score < bestScore || (score == bestScore && candidateLess(cand, best, placed)) {
			best = cand
			bestScore = score
		}
	}

	return best, true
}

func clusterWithCandidate(placed Cluster, cand Candidate) Cluster {
	cluster := append(Cluster(nil), placed...)
	cluster = append(cluster, NewPlacedItem(cand.Item, cand.CenterX, cand.CenterY))
	return cluster
}

func scoreCompactness(cand Candidate) float64 {
	return math.Hypot(cand.CenterX, cand.CenterY)
}

func scoreBalance(cluster Cluster) float64 {
	var quads [4]float64
	axisExt := 0.0

	for _, placed := range cluster {
		if placed.Item.Centerpiece {
			continue
		}
		area := float64(placed.Item.Height * placed.Item.Width)
		if q, ok := quadrantIndex(placed.Footprint.CenterX, placed.Footprint.CenterY); ok {
			quads[q] += area
			continue
		}
		axisExt += area
	}

	if axisExt > 0 && quadMass(quads) == 0 {
		return axisExt * 0.8
	}

	var sum float64
	for _, mass := range quads {
		sum += mass
	}
	if sum == 0 {
		return 0
	}

	avg := sum / 4.0
	penalty := 0.0
	for _, mass := range quads {
		penalty += math.Abs(mass - avg)
	}
	return penalty / avg
}

func quadMass(quads [4]float64) float64 {
	sum := 0.0
	for _, mass := range quads {
		sum += mass
	}
	return sum
}

func scoreBlobSmoothness(cluster Cluster) float64 {
	b := cluster.Bounds()
	spreadX := b.MaxX - b.MinX
	spreadY := b.MaxY - b.MinY
	if spreadY < 1 {
		spreadY = 1
	}
	if spreadX < 1 {
		spreadX = 1
	}

	aspect := spreadX / spreadY
	diff := aspect - targetAspectRatio
	penalty := diff * diff
	if aspect > maxBlobAspectRatio {
		penalty += (aspect - maxBlobAspectRatio) * 3
	}
	if aspect < 1.0/maxBlobAspectRatio {
		penalty += (1.0/maxBlobAspectRatio - aspect) * 3
	}
	return penalty
}

func scoreCollinearity(cluster Cluster) float64 {
	if len(cluster) < 3 {
		return 0
	}

	maxRow := maxAxisCount(cluster, true)
	maxCol := maxAxisCount(cluster, false)
	n := float64(len(cluster))

	penalty := 0.0
	if maxRow/n > 0.5 {
		penalty += (maxRow/n - 0.5) * n * 3
	}
	if maxCol/n > 0.5 {
		penalty += (maxCol/n - 0.5) * n * 3
	}
	return penalty
}

func scoreAxisOnlyPlacement(cand Candidate, cluster Cluster) float64 {
	if IsDiagonalPlacement(cand.CenterX, cand.CenterY) {
		return 0
	}
	if math.Abs(cand.CenterX) <= axisEpsilon && math.Abs(cand.CenterY) <= axisEpsilon {
		return 0
	}

	penalty := 6.0
	if len(cluster) >= 3 {
		penalty += 4.0
	}
	if clusterHasDiagonalPlacement(cluster) {
		penalty += 6.0
	}
	return penalty
}

func clusterHasDiagonalPlacement(cluster Cluster) bool {
	for _, placed := range cluster {
		if placed.Item.Centerpiece {
			continue
		}
		if IsDiagonalPlacement(placed.Footprint.CenterX, placed.Footprint.CenterY) {
			return true
		}
	}
	return false
}

func scoreDiagonalFill(cand Candidate) float64 {
	if IsDiagonalPlacement(cand.CenterX, cand.CenterY) {
		return 3.0
	}
	return 0
}

func scoreQuadrantFill(cand Candidate, cluster Cluster) float64 {
	counts := [4]int{}
	for _, placed := range cluster {
		if q, ok := quadrantIndex(placed.Footprint.CenterX, placed.Footprint.CenterY); ok {
			counts[q]++
		}
	}

	q, ok := quadrantIndex(cand.CenterX, cand.CenterY)
	if !ok {
		return 0
	}
	return 2.5 / (1.0 + float64(counts[q]))
}

func scoreClusterArea(cluster Cluster) float64 {
	b := cluster.Bounds()
	width := b.MaxX - b.MinX
	height := b.MaxY - b.MinY
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return width * height / 100.0
}

// IsDiagonalPlacement reports whether a center sits off both axes from the anchor.
func IsDiagonalPlacement(centerX, centerY float64) bool {
	return math.Abs(centerX) > axisEpsilon && math.Abs(centerY) > axisEpsilon
}

func quadrantIndex(centerX, centerY float64) (int, bool) {
	if math.Abs(centerX) <= axisEpsilon || math.Abs(centerY) <= axisEpsilon {
		return 0, false
	}
	switch {
	case centerX > 0 && centerY < 0:
		return 0, true // NE
	case centerX < 0 && centerY < 0:
		return 1, true // NW
	case centerX > 0 && centerY > 0:
		return 2, true // SE
	default:
		return 3, true // SW
	}
}

func scoreAxisAttachment(cand Candidate, cluster Cluster) float64 {
	cy := int64(math.Round(cand.CenterY))
	cx := int64(math.Round(cand.CenterX))

	rowCount := 0
	colCount := 0
	for _, placed := range cluster {
		if int64(math.Round(placed.Footprint.CenterY)) == cy {
			rowCount++
		}
		if int64(math.Round(placed.Footprint.CenterX)) == cx {
			colCount++
		}
	}

	penalty := 0.0
	if rowCount >= 3 {
		penalty += float64(rowCount-2) * 8
	}
	if colCount >= 3 {
		penalty += float64(colCount-2) * 8
	}
	return penalty
}

func maxAxisCount(cluster Cluster, useY bool) float64 {
	counts := make(map[int64]int)
	for _, placed := range cluster {
		var key int64
		if useY {
			key = int64(math.Round(placed.Footprint.CenterY))
		} else {
			key = int64(math.Round(placed.Footprint.CenterX))
		}
		counts[key]++
	}

	maxCount := 0
	for _, count := range counts {
		if count > maxCount {
			maxCount = count
		}
	}
	return float64(maxCount)
}

func scoreConcavity(cluster Cluster) float64 {
	b := cluster.Bounds()
	bboxArea := (b.MaxX - b.MinX) * (b.MaxY - b.MinY)
	if bboxArea <= 0 {
		return 0
	}

	hullArea := convexHullArea(clusterBBoxCorners(cluster))
	if hullArea <= 0 {
		return 1
	}

	fillRatio := hullArea / bboxArea
	if fillRatio > 1 {
		fillRatio = 1
	}
	return 1 - fillRatio
}

func clusterBBoxCorners(cluster Cluster) []point {
	pts := make([]point, 0, len(cluster)*4)
	for _, placed := range cluster {
		x, y, w, h := placed.Footprint.BBox()
		right := x + w
		bottom := y + h
		pts = append(pts,
			point{x: x, y: y},
			point{x: right, y: y},
			point{x: x, y: bottom},
			point{x: right, y: bottom},
		)
	}
	return pts
}

func scoreLocalContinuity(cand Candidate, placed Cluster, gap int) float64 {
	fp := cand.footprint()
	neighbors := 0
	for _, p := range placed {
		if Adjacent(fp, p.Footprint, gap) {
			neighbors++
		}
	}
	if neighbors < 2 {
		return 0
	}
	return float64(neighbors - 1)
}

func candidateLess(a, b Candidate, placed Cluster) bool {
	aDiag := IsDiagonalPlacement(a.CenterX, a.CenterY)
	bDiag := IsDiagonalPlacement(b.CenterX, b.CenterY)
	if aDiag != bDiag {
		return aDiag
	}

	if clusterSpreadIsFlat(placed, true) {
		if math.Abs(a.CenterY) != math.Abs(b.CenterY) {
			return math.Abs(a.CenterY) > math.Abs(b.CenterY)
		}
	}
	if clusterSpreadIsFlat(placed, false) {
		if math.Abs(a.CenterX) != math.Abs(b.CenterX) {
			return math.Abs(a.CenterX) > math.Abs(b.CenterX)
		}
	}
	if a.CenterX != b.CenterX {
		return a.CenterX < b.CenterX
	}
	return a.CenterY < b.CenterY
}

func clusterSpreadIsFlat(placed Cluster, horizontal bool) bool {
	if len(placed) < 2 {
		return false
	}
	b := placed.Bounds()
	spreadX := b.MaxX - b.MinX
	spreadY := b.MaxY - b.MinY
	if spreadY < 1 {
		spreadY = 1
	}
	if spreadX < 1 {
		spreadX = 1
	}
	if horizontal {
		return spreadX/spreadY > 1.8
	}
	return spreadY/spreadX > 1.8
}

type point struct {
	x float64
	y float64
}

func clusterCenters(cluster Cluster) []point {
	pts := make([]point, 0, len(cluster))
	for _, placed := range cluster {
		pts = append(pts, point{
			x: placed.Footprint.CenterX,
			y: placed.Footprint.CenterY,
		})
	}
	return pts
}

func convexHullArea(pts []point) float64 {
	if len(pts) < 3 {
		return 0
	}

	hull := convexHull(pts)
	if len(hull) < 3 {
		return 0
	}
	return polygonArea(hull)
}

func convexHull(pts []point) []point {
	sorted := append([]point(nil), pts...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].x == sorted[j].x {
			return sorted[i].y < sorted[j].y
		}
		return sorted[i].x < sorted[j].x
	})

	cross := func(o, a, b point) float64 {
		return (a.x-o.x)*(b.y-o.y) - (a.y-o.y)*(b.x-o.x)
	}

	lower := make([]point, 0, len(sorted))
	for _, p := range sorted {
		for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
			lower = lower[:len(lower)-1]
		}
		lower = append(lower, p)
	}

	upper := make([]point, 0, len(sorted))
	for i := len(sorted) - 1; i >= 0; i-- {
		p := sorted[i]
		for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
			upper = upper[:len(upper)-1]
		}
		upper = append(upper, p)
	}

	return append(lower[:len(lower)-1], upper[:len(upper)-1]...)
}

func polygonArea(hull []point) float64 {
	if len(hull) < 3 {
		return 0
	}

	area := 0.0
	for i := 0; i < len(hull); i++ {
		j := (i + 1) % len(hull)
		area += hull[i].x*hull[j].y - hull[j].x*hull[i].y
	}
	return math.Abs(area) / 2
}
