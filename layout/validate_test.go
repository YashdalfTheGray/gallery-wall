package layout

import "testing"

func TestValidateEmptyInput(t *testing.T) {
	tests := []struct {
		name   string
		params Params
	}{
		{
			name:   "nil items",
			params: Params{Gap: 2, Items: nil},
		},
		{
			name:   "empty slice",
			params: Params{Gap: 2, Items: []Item{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.params)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !IsValidationCode(err, ValidationEmptyInput) {
				t.Fatalf("code mismatch: %v", err)
			}
		})
	}
}

func TestValidateDuplicateID(t *testing.T) {
	params := Params{
		Gap: 2,
		Items: []Item{
			{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
			{ID: "a", Height: 8, Width: 8, Shape: ShapeSquare, Centerpiece: false},
		},
	}

	err := Validate(params)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsValidationCode(err, ValidationDuplicateID) {
		t.Fatalf("code mismatch: %v", err)
	}
	ve := err.(*ValidationError)
	if ve.ItemID != "a" {
		t.Fatalf("itemId = %q, want %q", ve.ItemID, "a")
	}
}

func TestValidateInvalidDimensions(t *testing.T) {
	tests := []struct {
		name string
		item Item
	}{
		{
			name: "zero height",
			item: Item{ID: "a", Height: 0, Width: 10, Shape: ShapeSquare, Centerpiece: true},
		},
		{
			name: "zero width",
			item: Item{ID: "b", Height: 10, Width: 0, Shape: ShapeSquare, Centerpiece: true},
		},
		{
			name: "negative height",
			item: Item{ID: "c", Height: -1, Width: 10, Shape: ShapeSquare, Centerpiece: true},
		},
		{
			name: "negative width",
			item: Item{ID: "d", Height: 10, Width: -1, Shape: ShapeSquare, Centerpiece: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(Params{Gap: 2, Items: []Item{tt.item}})
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !IsValidationCode(err, ValidationInvalidDimensions) {
				t.Fatalf("code mismatch: %v", err)
			}
			ve := err.(*ValidationError)
			if ve.ItemID != tt.item.ID {
				t.Fatalf("itemId = %q, want %q", ve.ItemID, tt.item.ID)
			}
		})
	}
}

func TestValidateCenterpiece(t *testing.T) {
	t.Run("no centerpiece", func(t *testing.T) {
		err := Validate(Params{
			Gap: 2,
			Items: []Item{
				{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: false},
			},
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !IsValidationCode(err, ValidationNoCenterpiece) {
			t.Fatalf("code mismatch: %v", err)
		}
	})

	t.Run("multiple centerpieces", func(t *testing.T) {
		err := Validate(Params{
			Gap: 2,
			Items: []Item{
				{ID: "a", Height: 10, Width: 10, Shape: ShapeSquare, Centerpiece: true},
				{ID: "b", Height: 8, Width: 8, Shape: ShapeSquare, Centerpiece: true},
			},
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !IsValidationCode(err, ValidationMultipleCenterpiece) {
			t.Fatalf("code mismatch: %v", err)
		}
		ve := err.(*ValidationError)
		if len(ve.ItemIDs) != 2 || ve.ItemIDs[0] != "a" || ve.ItemIDs[1] != "b" {
			t.Fatalf("itemIds = %v", ve.ItemIDs)
		}
	})

	t.Run("exactly one centerpiece", func(t *testing.T) {
		err := Validate(Params{
			Gap: 2,
			Items: []Item{
				{ID: "main", Height: 24, Width: 18, Shape: ShapeRectangle, Centerpiece: true},
				{ID: "left", Height: 10, Width: 8, Shape: ShapeRectangle, Centerpiece: false},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
