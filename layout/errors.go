package layout

import (
	"errors"
	"fmt"
)

// ValidationCode identifies a validation failure for machine-readable handling.
// Integrators should compare Code (or use [IsValidationCode] / errors.Is), not Message.
type ValidationCode string

const (
	ValidationEmptyInput          ValidationCode = "empty_input"
	ValidationDuplicateID         ValidationCode = "duplicate_id"
	ValidationInvalidDimensions   ValidationCode = "invalid_dimensions"
	ValidationNoCenterpiece       ValidationCode = "no_centerpiece"
	ValidationMultipleCenterpiece ValidationCode = "multiple_centerpieces"
	ValidationInvalidWall         ValidationCode = "invalid_wall"
	ValidationCenterpieceExceedsWall ValidationCode = "centerpiece_exceeds_wall"
)

// ValidationError describes invalid layout input.
type ValidationError struct {
	Code    ValidationCode `json:"code"`
	Message string         `json:"message"`
	ItemID  string         `json:"itemId,omitempty"`
	ItemIDs []string       `json:"itemIds,omitempty"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Is reports whether target matches this validation error.
// Pass &ValidationError{Code: ValidationDuplicateID} to match by code;
// pass &ValidationError{} to match any validation error.
func (e *ValidationError) Is(target error) bool {
	t, ok := target.(*ValidationError)
	if !ok {
		return false
	}
	if t.Code == "" {
		return true
	}
	return e.Code == t.Code
}

// IsValidationCode reports whether err is a [ValidationError] with the given code.
func IsValidationCode(err error, code ValidationCode) bool {
	var ve *ValidationError
	if !errors.As(err, &ve) {
		return false
	}
	return ve.Code == code
}

func newEmptyInputError() *ValidationError {
	return &ValidationError{
		Code:    ValidationEmptyInput,
		Message: "empty input",
	}
}

func newDuplicateIDError(id string) *ValidationError {
	return &ValidationError{
		Code:    ValidationDuplicateID,
		Message: fmt.Sprintf("duplicate id: %s", id),
		ItemID:  id,
	}
}

func newInvalidDimensionsError(id string) *ValidationError {
	return &ValidationError{
		Code:    ValidationInvalidDimensions,
		Message: fmt.Sprintf("invalid dimensions: %s", id),
		ItemID:  id,
	}
}

func newNoCenterpieceError() *ValidationError {
	return &ValidationError{
		Code:    ValidationNoCenterpiece,
		Message: "no centerpiece specified",
	}
}

func newMultipleCenterpiecesError(ids []string) *ValidationError {
	return &ValidationError{
		Code:    ValidationMultipleCenterpiece,
		Message: fmt.Sprintf("multiple centerpieces: %v", ids),
		ItemIDs: append([]string(nil), ids...),
	}
}

func newInvalidWallError(message string) *ValidationError {
	return &ValidationError{
		Code:    ValidationInvalidWall,
		Message: message,
	}
}

func newCenterpieceExceedsWallError(itemID string) *ValidationError {
	return &ValidationError{
		Code:    ValidationCenterpieceExceedsWall,
		Message: fmt.Sprintf("centerpiece %s does not fit in wall bounds", itemID),
		ItemID:  itemID,
	}
}

// PlacementCode identifies a placement failure for machine-readable handling.
type PlacementCode string

const (
	PlacementCannotPlace PlacementCode = "cannot_place_item"
)

// PlacementError describes a failure while placing a frame.
type PlacementError struct {
	Code    PlacementCode `json:"code"`
	Message string        `json:"message"`
	ItemID  string        `json:"itemId"`
}

func (e *PlacementError) Error() string {
	return e.Message
}

// Is reports whether target matches this placement error.
// Pass &PlacementError{Code: PlacementCannotPlace} to match by code;
// pass &PlacementError{} to match any placement error.
func (e *PlacementError) Is(target error) bool {
	t, ok := target.(*PlacementError)
	if !ok {
		return false
	}
	if t.Code == "" {
		return true
	}
	return e.Code == t.Code
}

// IsPlacementCode reports whether err is a [PlacementError] with the given code.
func IsPlacementCode(err error, code PlacementCode) bool {
	var pe *PlacementError
	if !errors.As(err, &pe) {
		return false
	}
	return pe.Code == code
}

func newCannotPlaceError(itemID string) *PlacementError {
	return &PlacementError{
		Code:    PlacementCannotPlace,
		Message: fmt.Sprintf(`cannot place item "%s" — no valid connected position`, itemID),
		ItemID:  itemID,
	}
}
