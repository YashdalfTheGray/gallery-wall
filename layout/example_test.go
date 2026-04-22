package layout_test

import (
	"encoding/json"
	"fmt"

	"github.com/yashdalfthegray/gallery-wall/layout"
)

func ExampleLayout() {
	result, err := layout.Layout(layout.Params{
		Gap: 2,
		Items: []layout.Item{
			{ID: "main", Height: 20, Width: 16, Shape: layout.ShapeRectangle, Centerpiece: true},
			{ID: "left", Height: 10, Width: 8, Shape: layout.ShapeRectangle},
			{ID: "right", Height: 10, Width: 8, Shape: layout.ShapeRectangle},
		},
	})
	if err != nil {
		panic(err)
	}

	main := result.Items[0]
	fmt.Printf("centerpiece at (%.0f, %.0f)\n", main.CenterX, main.CenterY)
	fmt.Printf("cluster bounds %.0f×%.0f\n",
		result.Bounds.MaxX-result.Bounds.MinX,
		result.Bounds.MaxY-result.Bounds.MinY)
	// Output:
	// centerpiece at (0, 0)
	// cluster bounds 36×20
}

func ExampleLayout_validationError() {
	_, err := layout.Layout(layout.Params{Items: nil})
	if layout.IsValidationCode(err, layout.ValidationEmptyInput) {
		fmt.Println("invalid input")
	}
	// Output:
	// invalid input
}

func ExampleLayout_jsonRoundTrip() {
	params := layout.Params{
		Gap: 2,
		Items: []layout.Item{
			{ID: "main", Height: 10, Width: 10, Shape: layout.ShapeSquare, Centerpiece: true},
			{ID: "side", Height: 8, Width: 6, Shape: layout.ShapeRectangle},
		},
	}

	encoded, _ := json.Marshal(params)
	var decoded layout.Params
	_ = json.Unmarshal(encoded, &decoded)

	result, err := layout.Layout(decoded)
	if err != nil {
		panic(err)
	}

	out, _ := json.Marshal(result)
	fmt.Printf("placed %d items, json ok=%t\n", len(result.Items), len(out) > 0)
	// Output:
	// placed 2 items, json ok=true
}
