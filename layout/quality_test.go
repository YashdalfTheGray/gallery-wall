package layout

import "testing"

func TestBlobQualityGoldenFixtures(t *testing.T) {
	fixtures := []string{
		"centerpiece_only",
		"centerpiece_two_symmetric",
		"centerpiece_six_mixed",
	}

	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			params := loadGoldenParams(t, name)
			result, err := Layout(params)
			if err != nil {
				t.Fatalf("layout error: %v", err)
			}

			centerpieceID := result.Anchor.ItemID
			if len(result.Items) > 1 && !AllItemsHaveNeighbor(result, centerpieceID) {
				t.Fatal("non-centerpiece item without neighbor")
			}

			if len(result.Items) > 1 {
				aspect := ClusterAspectRatio(result.Bounds)
				if aspect < minClusterAspectRatio || aspect > maxClusterAspectRatio {
					t.Fatalf("aspect ratio %v outside [%v, %v]", aspect, minClusterAspectRatio, maxClusterAspectRatio)
				}
			}

			if len(result.Items) >= 3 {
				fill := ClusterHullFillRatio(resultCluster(result))
				if fill < minHullFillRatio {
					t.Fatalf("hull fill ratio %v below minimum %v", fill, minHullFillRatio)
				}
			}

			if !ClusterHas2DSpread(result) {
				t.Fatal("cluster lacks 2D spread (blob should not be a single line)")
			}

			if !PassesBlobQuality(result, resultCluster(result), centerpieceID) {
				t.Fatal("failed blob quality checks")
			}
		})
	}
}

func TestClusterAspectRatio(t *testing.T) {
	bounds := Bounds{MinX: -10, MinY: -5, MaxX: 14, MaxY: 5}
	got := ClusterAspectRatio(bounds)
	assertFloat(t, "aspect", got, 2.4)
}

func TestAllItemsHaveNeighbor(t *testing.T) {
	result := Result{
		Anchor: Anchor{ItemID: "main"},
		Items: []PlacedResult{
			{ID: "main", AdjacentIDs: []string{"side"}},
			{ID: "side", AdjacentIDs: []string{"main"}},
		},
	}
	if !AllItemsHaveNeighbor(result, "main") {
		t.Fatal("expected all items to have neighbors")
	}
}
