package layout

import (
	"fmt"
	"math"
)

const (
	candidateStep = 1
)

// Candidate is a possible placement for an unplaced item.
type Candidate struct {
	Item    Item
	CenterX float64
	CenterY float64
}

func (c Candidate) footprint() Footprint {
	return NewFootprint(c.Item, c.CenterX, c.CenterY)
}

// GenerateCandidates returns valid attachment positions for item against placed cluster.
// Exported for tests; [Layout] uses this internally. Pass nil wall for no wall limit.
func GenerateCandidates(item Item, placed Cluster, gap int, wall *Bounds) []Candidate {
	raw := make([]Candidate, 0)
	for _, anchor := range placed {
		raw = append(raw, sideCandidates(anchor, item, gap)...)
		raw = append(raw, cornerCandidates(anchor, item, gap)...)
	}
	return filterCandidates(placed, gap, wall, dedupeCandidates(raw))
}

func sideCandidates(anchor PlacedItem, item Item, gap int) []Candidate {
	a := anchor.Footprint
	aHX, aHY := a.HalfExtents()
	pHX, pHY := halfExtents(item)
	g := float64(gap)

	out := make([]Candidate, 0)
	for _, centerY := range slideValues(a.CenterY, aHY, pHY) {
		out = append(out, Candidate{
			Item:    item,
			CenterX: a.CenterX - (aHX + g + pHX),
			CenterY: centerY,
		})
		out = append(out, Candidate{
			Item:    item,
			CenterX: a.CenterX + (aHX + g + pHX),
			CenterY: centerY,
		})
	}

	for _, centerX := range slideValues(a.CenterX, aHX, pHX) {
		out = append(out, Candidate{
			Item:    item,
			CenterX: centerX,
			CenterY: a.CenterY - (aHY + g + pHY),
		})
		out = append(out, Candidate{
			Item:    item,
			CenterX: centerX,
			CenterY: a.CenterY + (aHY + g + pHY),
		})
	}

	return out
}

type cornerKind int

const (
	cornerTopLeft cornerKind = iota
	cornerTopRight
	cornerBottomLeft
	cornerBottomRight
)

func cornerCandidates(anchor PlacedItem, item Item, gap int) []Candidate {
	out := make([]Candidate, 0)
	for _, kind := range []cornerKind{
		cornerTopLeft,
		cornerTopRight,
		cornerBottomLeft,
		cornerBottomRight,
	} {
		out = append(out, cornerCandidatesForKind(anchor, item, gap, kind)...)
	}
	return out
}

func cornerCandidatesForKind(anchor PlacedItem, item Item, gap int, kind cornerKind) []Candidate {
	a := anchor.Footprint
	aHX, aHY := a.HalfExtents()
	pHX, pHY := halfExtents(item)
	g := float64(gap)

	var baseX, baseY float64
	switch kind {
	case cornerTopLeft:
		baseX = a.CenterX - (aHX + g + pHX)
		baseY = a.CenterY - (aHY + g + pHY)
	case cornerTopRight:
		baseX = a.CenterX + (aHX + g + pHX)
		baseY = a.CenterY - (aHY + g + pHY)
	case cornerBottomLeft:
		baseX = a.CenterX - (aHX + g + pHX)
		baseY = a.CenterY + (aHY + g + pHY)
	case cornerBottomRight:
		baseX = a.CenterX + (aHX + g + pHX)
		baseY = a.CenterY + (aHY + g + pHY)
	}

	out := make([]Candidate, 0)
	for _, centerX := range slideValues(baseX, aHY, pHY) {
		out = append(out, Candidate{Item: item, CenterX: centerX, CenterY: baseY})
	}
	for _, centerY := range slideValues(baseY, aHX, pHX) {
		if centerY == baseY {
			continue
		}
		out = append(out, Candidate{Item: item, CenterX: baseX, CenterY: centerY})
	}
	return out
}

func slideValues(center, halfA, halfP float64) []float64 {
	start := math.Round(center - (halfA + halfP))
	end := math.Round(center + (halfA + halfP))
	count := int(end-start) + 1
	if count <= 0 {
		return []float64{center}
	}

	values := make([]float64, 0, count)
	for v := start; v <= end; v += candidateStep {
		values = append(values, v)
	}
	return values
}

func filterCandidates(placed Cluster, gap int, wall *Bounds, candidates []Candidate) []Candidate {
	out := make([]Candidate, 0, len(candidates))
	for _, cand := range candidates {
		if !isValidCandidate(cand, placed, gap, wall) {
			continue
		}
		out = append(out, cand)
	}
	return out
}

func isValidCandidate(cand Candidate, placed Cluster, gap int, wall *Bounds) bool {
	fp := cand.footprint()

	for _, p := range placed {
		if Collides(fp, p.Footprint, gap) {
			return false
		}
	}

	adjacent := false
	for _, p := range placed {
		if Adjacent(fp, p.Footprint, gap) {
			adjacent = true
			break
		}
	}
	if !adjacent {
		return false
	}

	if wall != nil && !footprintFitsWall(fp, *wall) {
		return false
	}
	return true
}

func dedupeCandidates(candidates []Candidate) []Candidate {
	seen := make(map[string]struct{}, len(candidates))
	out := make([]Candidate, 0, len(candidates))

	for _, cand := range candidates {
		key := candidateKey(cand)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, cand)
	}
	return out
}

func candidateKey(c Candidate) string {
	x := int64(math.Round(c.CenterX / candidateStep))
	y := int64(math.Round(c.CenterY / candidateStep))
	return fmt.Sprintf("%d,%d", x, y)
}
