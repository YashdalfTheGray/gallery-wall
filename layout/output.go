package layout

import "slices"

func adjacentIDsFor(placed Cluster, item PlacedItem, gap int) []string {
	ids := make([]string, 0)
	for _, other := range placed {
		if other.Item.ID == item.Item.ID {
			continue
		}
		if Adjacent(item.Footprint, other.Footprint, gap) {
			ids = append(ids, other.Item.ID)
		}
	}
	slices.Sort(ids)
	return ids
}
