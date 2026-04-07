package layout

import "math"

const (
	minClusterAspectRatio = 0.9
	maxClusterAspectRatio = 2.2
	minHullFillRatio      = 0.12
	minUniqueAxes         = 2
)

// ClusterAspectRatio returns cluster bbox width divided by height.
func ClusterAspectRatio(bounds Bounds) float64 {
	width := bounds.MaxX - bounds.MinX
	height := bounds.MaxY - bounds.MinY
	if height < 1 {
		height = 1
	}
	if width < 1 {
		width = 1
	}
	return width / height
}

// ClusterHullFillRatio returns convex hull area divided by cluster bbox area.
func ClusterHullFillRatio(cluster Cluster) float64 {
	bounds := cluster.Bounds()
	bboxArea := (bounds.MaxX - bounds.MinX) * (bounds.MaxY - bounds.MinY)
	if bboxArea <= 0 {
		return 1
	}

	hullArea := convexHullArea(clusterBBoxCorners(cluster))
	if hullArea <= 0 {
		return 0
	}
	return hullArea / bboxArea
}

// AllItemsHaveNeighbor reports whether every non-centerpiece item lists a neighbor.
func AllItemsHaveNeighbor(result Result, centerpieceID string) bool {
	for _, item := range result.Items {
		if item.ID == centerpieceID {
			continue
		}
		if len(item.AdjacentIDs) == 0 {
			return false
		}
	}
	return true
}

// PassesBlobQuality checks silhouette heuristics for a layout result.
func PassesBlobQuality(result Result, cluster Cluster, centerpieceID string) bool {
	if len(result.Items) > 1 && !AllItemsHaveNeighbor(result, centerpieceID) {
		return false
	}

	if len(result.Items) > 1 {
		aspect := ClusterAspectRatio(result.Bounds)
		if aspect < minClusterAspectRatio || aspect > maxClusterAspectRatio {
			return false
		}
	}

	if len(result.Items) >= 3 && ClusterHullFillRatio(cluster) < minHullFillRatio {
		return false
	}

	if len(result.Items) >= 4 && !ClusterHas2DSpread(result) {
		return false
	}

	return true
}

// ClusterHas2DSpread reports whether item centers occupy multiple rows and columns.
func ClusterHas2DSpread(result Result) bool {
	if len(result.Items) < 4 {
		return true
	}

	rows := make(map[int64]struct{})
	cols := make(map[int64]struct{})
	for _, item := range result.Items {
		rows[int64(math.Round(item.CenterY))] = struct{}{}
		cols[int64(math.Round(item.CenterX))] = struct{}{}
	}
	return len(rows) >= minUniqueAxes && len(cols) >= minUniqueAxes
}

// CountDiagonalPlacements returns frames whose centers are off both axes from the anchor.
func CountDiagonalPlacements(result Result) int {
	count := 0
	for _, item := range result.Items {
		if IsDiagonalPlacement(item.CenterX, item.CenterY) {
			count++
		}
	}
	return count
}

func resultCluster(result Result) Cluster {
	cluster := make(Cluster, 0, len(result.Items))
	for _, item := range result.Items {
		cluster = append(cluster, NewPlacedItem(Item{
			ID:     item.ID,
			Height: item.Height,
			Width:  item.Width,
			Shape:  item.Shape,
		}, item.CenterX, item.CenterY))
	}
	return cluster
}
