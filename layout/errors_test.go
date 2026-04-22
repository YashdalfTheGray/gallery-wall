package layout

import (
	"encoding/json"
	"testing"
)

func TestValidationErrorError(t *testing.T) {
	err := newDuplicateIDError("a")
	if got := err.Error(); got != "duplicate id: a" {
		t.Fatalf("Error() = %q, want %q", got, "duplicate id: a")
	}
}

func TestIsValidationCode(t *testing.T) {
	err := newInvalidDimensionsError("main")

	if !IsValidationCode(err, ValidationInvalidDimensions) {
		t.Fatal("expected invalid_dimensions code")
	}
	if IsValidationCode(err, ValidationEmptyInput) {
		t.Fatal("unexpected empty_input match")
	}
	if IsValidationCode(nil, ValidationEmptyInput) {
		t.Fatal("nil should not match")
	}
}

func TestValidationErrorJSON(t *testing.T) {
	err := newMultipleCenterpiecesError([]string{"a", "b"})
	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("marshal: %v", marshalErr)
	}

	var decoded ValidationError
	if unmarshalErr := json.Unmarshal(data, &decoded); unmarshalErr != nil {
		t.Fatalf("unmarshal: %v", unmarshalErr)
	}

	if decoded.Code != ValidationMultipleCenterpiece {
		t.Fatalf("code = %q, want %q", decoded.Code, ValidationMultipleCenterpiece)
	}
	if len(decoded.ItemIDs) != 2 || decoded.ItemIDs[0] != "a" || decoded.ItemIDs[1] != "b" {
		t.Fatalf("itemIds = %v", decoded.ItemIDs)
	}
}

func TestValidationCodesDistinct(t *testing.T) {
	codes := []ValidationCode{
		ValidationEmptyInput,
		ValidationDuplicateID,
		ValidationInvalidDimensions,
		ValidationNoCenterpiece,
		ValidationMultipleCenterpiece,
		ValidationInvalidWall,
		ValidationCenterpieceExceedsWall,
	}
	seen := make(map[ValidationCode]struct{}, len(codes))
	for _, code := range codes {
		if _, ok := seen[code]; ok {
			t.Fatalf("duplicate validation code: %q", code)
		}
		seen[code] = struct{}{}
	}
}
