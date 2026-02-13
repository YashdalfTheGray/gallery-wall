package layout

// WallBounds returns the axis-aligned wall rectangle centered on the anchor (0, 0).
// The second value is false when no wall constraint is configured.
func (p Params) WallBounds() (Bounds, bool) {
	if p.WallWidth <= 0 || p.WallHeight <= 0 {
		return Bounds{}, false
	}
	halfW := p.WallWidth / 2
	halfH := p.WallHeight / 2
	return Bounds{
		MinX: -halfW,
		MinY: -halfH,
		MaxX: halfW,
		MaxY: halfH,
	}, true
}

func footprintFitsWall(fp Footprint, wall Bounds) bool {
	x, y, w, h := fp.BBox()
	return x >= wall.MinX && y >= wall.MinY && x+w <= wall.MaxX && y+h <= wall.MaxY
}

func validateWall(params Params) error {
	hasWidth := params.WallWidth > 0
	hasHeight := params.WallHeight > 0
	if !hasWidth && !hasHeight {
		return nil
	}
	if !hasWidth || !hasHeight {
		return newInvalidWallError("wall width and height must both be set")
	}

	wall, _ := params.WallBounds()
	for _, item := range params.Items {
		if !item.Centerpiece {
			continue
		}
		fp := NewFootprint(item, 0, 0)
		if !footprintFitsWall(fp, wall) {
			return newCenterpieceExceedsWallError(item.ID)
		}
		return nil
	}
	return nil
}
