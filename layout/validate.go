package layout

// Validate checks Params for structural input errors without running placement.
func Validate(params Params) error {
	if len(params.Items) == 0 {
		return newEmptyInputError()
	}

	seen := make(map[string]struct{}, len(params.Items))
	var centerpieceIDs []string

	for _, item := range params.Items {
		if _, ok := seen[item.ID]; ok {
			return newDuplicateIDError(item.ID)
		}
		seen[item.ID] = struct{}{}

		if item.Height <= 0 || item.Width <= 0 {
			return newInvalidDimensionsError(item.ID)
		}

		if item.Centerpiece {
			centerpieceIDs = append(centerpieceIDs, item.ID)
		}
	}

	switch len(centerpieceIDs) {
	case 0:
		return newNoCenterpieceError()
	case 1:
		return validateWall(params)
	default:
		return newMultipleCenterpiecesError(centerpieceIDs)
	}
}
