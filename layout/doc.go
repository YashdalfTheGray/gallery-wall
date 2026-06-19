// Package layout arranges gallery-wall frames organically around a centerpiece.
//
// Install as a standalone module:
//
//	go get github.com/yashdalfthegray/gallery-wall/layout@latest
//
// See README.md in this directory for usage and release tags (layout/vX.Y.Z).
//
// Algorithm walkthrough: see ALGORITHM.md in this directory.
//
// The centerpiece is anchored at (0, 0). All other frames attach edge-to-edge
// (with a configurable gap) to form one connected cluster—a compact blob
// suitable for hanging on a real wall.
//
// # Entry point
//
// Call [Layout] with [Params] to produce a [Result] or a typed error.
// [Validate] checks input without placing frames.
//
// # Stable API (for HTTP integrators)
//
//   - [Layout], [Validate]
//   - [Params], [Item], [Shape] — optional WallWidth/WallHeight constrain the cluster
//   - [Result], [PlacedResult], [Anchor], [Bounds]
//   - [ValidationError], [ValidationCode], [PlacementError], [PlacementCode]
//   - [IsValidationCode], [IsPlacementCode] and errors.Is with partial errors
//
// # Advanced exports
//
// Geometry ([Footprint], [Collides], [Adjacent]), placement ([Candidate],
// [GenerateCandidates], [Cluster]), and scoring ([ScoreCandidate], [BestCandidate])
// are exported for tests and algorithm tuning. Prefer [Layout] for production use.
//
// Quality helpers ([ClusterAspectRatio], [ClusterHas2DSpread], etc.) support
// regression tests and fixture validation.
package layout
