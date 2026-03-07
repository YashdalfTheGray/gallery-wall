package layout

import (
	"cmp"
	"slices"
)

func findCenterpiece(items []Item) (Item, []Item, error) {
	var centerpiece Item
	var found bool
	rest := make([]Item, 0, len(items))

	for _, item := range items {
		if item.Centerpiece {
			if found {
				return Item{}, nil, newMultipleCenterpiecesError(centerpieceIDs(items))
			}
			centerpiece = item
			found = true
			continue
		}
		rest = append(rest, item)
	}

	if !found {
		return Item{}, nil, newNoCenterpieceError()
	}
	return centerpiece, rest, nil
}

func centerpieceIDs(items []Item) []string {
	ids := make([]string, 0)
	for _, item := range items {
		if item.Centerpiece {
			ids = append(ids, item.ID)
		}
	}
	return ids
}

func sortForPlacement(items []Item) []Item {
	sorted := append([]Item(nil), items...)
	slices.SortFunc(sorted, func(a, b Item) int {
		areaA := a.Height * a.Width
		areaB := b.Height * b.Width
		if areaA != areaB {
			return cmp.Compare(areaB, areaA)
		}

		maxA := a.Height
		if a.Width > maxA {
			maxA = a.Width
		}
		maxB := b.Height
		if b.Width > maxB {
			maxB = b.Width
		}
		if maxA != maxB {
			return cmp.Compare(maxB, maxA)
		}

		return cmp.Compare(a.ID, b.ID)
	})
	return sorted
}

func anchorCenterpiece(item Item) PlacedItem {
	return NewPlacedItem(item, 0, 0)
}
