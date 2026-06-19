# layout

Arrange picture frames of mixed sizes and shapes around a centerpiece. The
centerpiece sits at `(0, 0)`; every other frame attaches edge-to-edge (with a
configurable gap) to form one connected cluster.

This directory is a **standalone Go module**. You can depend on it without
pulling in the web app, WASM build, or CLI tools in the parent repository.

## Install

```bash
go get github.com/yashdalfthegray/gallery-wall/layout@latest
```

Or pin a release:

```bash
go get github.com/yashdalfthegray/gallery-wall/layout@v1.0.0
```

Releases are tagged on the parent repository as `layout/vX.Y.Z` (for example
`layout/v1.0.0`).

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/yashdalfthegray/gallery-wall/layout"
)

func main() {
	result, err := layout.Layout(layout.Params{
		Gap: 2,
		Items: []layout.Item{
			{ID: "main", Height: 20, Width: 16, Shape: layout.ShapeRectangle, Centerpiece: true},
			{ID: "left", Height: 10, Width: 8, Shape: layout.ShapeRectangle},
			{ID: "right", Height: 10, Width: 8, Shape: layout.ShapeRectangle},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("placed %d frames\n", len(result.Items))
	fmt.Printf("cluster bounds: %.0f×%.0f\n",
		result.Bounds.MaxX-result.Bounds.MinX,
		result.Bounds.MaxY-result.Bounds.MinY)
}
```

### Validate without placing

```go
if err := layout.Validate(params); err != nil {
	// handle layout.ValidationError
}
```

### JSON

`Params`, `Result`, and error types marshal/unmarshal JSON for HTTP APIs and
config files. See `example_test.go` for round-trip examples.

## API surface

| Category | Symbols |
|----------|---------|
| Entry points | `Layout`, `Validate` |
| Input | `Params`, `Item`, `Shape` |
| Output | `Result`, `PlacedResult`, `Anchor`, `Bounds` |
| Errors | `ValidationError`, `PlacementError`, `IsValidationCode`, `IsPlacementCode` |

Geometry, candidate generation, and scoring helpers are exported for tests and
tuning; prefer `Layout` in application code. Full godoc is in `doc.go`.

## Algorithm

See [ALGORITHM.md](./ALGORITHM.md) for the placement walkthrough.

## Development

From this directory:

```bash
go test .
go test -run Example -v
```

From the repository root (requires [go.work](../go.work) in a full clone):

```bash
go test ./layout/...
```
