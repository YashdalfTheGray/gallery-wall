package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type fileResult struct {
	Items  []placed `json:"items"`
	Bounds bounds   `json:"bounds"`
}

type placed struct {
	ID      string  `json:"id"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Width   int     `json:"width"`
	Height  int     `json:"height"`
	Shape   string  `json:"shape"`
	CenterX float64 `json:"centerX"`
	CenterY float64 `json:"centerY"`
}

type bounds struct {
	MinX, MinY, MaxX, MaxY float64
}

func main() {
	inPath := "layout/testdata/twentyfive_result.json"
	outPath := "layout/testdata/twentyfive_diagram.svg"
	if len(os.Args) > 1 {
		inPath = os.Args[1]
	}
	if len(os.Args) > 2 {
		outPath = os.Args[2]
	}

	data, err := os.ReadFile(inPath)
	if err != nil {
		panic(err)
	}
	var res fileResult
	if err := json.Unmarshal(data, &res); err != nil {
		panic(err)
	}

	const (
		scale   = 8.0
		padding = 24.0
		gapVis  = 2.0 // layout gap shown as stroke between frames
	)

	b := res.Bounds
	wallW := (b.MaxX - b.MinX) * scale
	wallH := (b.MaxY - b.MinY) * scale
	svgW := wallW + padding*2
	svgH := wallH + padding*2

	items := append([]placed(nil), res.Items...)
	sort.Slice(items, func(i, j int) bool {
		// draw centerpiece last so it sits on top visually
		if items[i].ID == "main" {
			return false
		}
		if items[j].ID == "main" {
			return true
		}
		return items[i].ID < items[j].ID
	})

	var out strings.Builder
	fmt.Fprintf(&out, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`+"\n", svgW, svgH, svgW, svgH)
	out.WriteString(`<rect width="100%" height="100%" fill="#1a1a1a"/>` + "\n")
	fmt.Fprintf(&out, `<text x="%.0f" y="16" fill="#aaa" font-family="system-ui,sans-serif" font-size="12">Gallery wall — 25 frames (1 unit ≈ %.0f px) · gap=2</text>`+"\n", padding, scale)

	// anchor crosshair at (0,0)
	ax := padding + (0-b.MinX)*scale
	ay := padding + (0-b.MinY)*scale
	fmt.Fprintf(&out, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#444" stroke-width="1" stroke-dasharray="4 3"/>`+"\n", ax-10, ay, ax+10, ay)
	fmt.Fprintf(&out, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#444" stroke-width="1" stroke-dasharray="4 3"/>`+"\n", ax, ay-10, ax, ay+10)

	for _, item := range items {
		x := padding + (item.X-b.MinX)*scale
		y := padding + (item.Y-b.MinY)*scale
		w := float64(item.Width) * scale
		h := float64(item.Height) * scale

		fill, stroke, sw := frameColors(item)
		rx := 0.0
		if item.Shape == "circle" {
			rx = math.Min(w, h) / 2
		} else if item.Shape == "ellipse" {
			rx = w / 2
		} else if item.Shape == "square" || item.ID == "main" {
			rx = 3
		}

		if item.Shape == "circle" || item.Shape == "ellipse" {
			cx := x + w/2
			cy := y + h/2
			var erx, ery float64
			if item.Shape == "circle" {
				erx = rx - gapVis/2
				ery = erx
			} else {
				erx = w/2 - gapVis/2
				ery = h/2 - gapVis/2
			}
			fmt.Fprintf(&out, `<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="%s" stroke="%s" stroke-width="%.1f"/>`+"\n",
				cx, cy, erx, ery, fill, stroke, sw)
		} else {
			fmt.Fprintf(&out, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="%.1f" fill="%s" stroke="%s" stroke-width="%.1f"/>`+"\n",
				x, y, w, h, rx, fill, stroke, sw)
		}

		label := item.ID
		if item.ID == "main" {
			label = "main\n15×13"
		} else {
			label = fmt.Sprintf("%s\n%d×%d", item.ID, item.Width, item.Height)
		}
		fontSize := 10.0
		if w < 40 || h < 40 {
			fontSize = 8
		}
		if w < 28 || h < 28 {
			fontSize = 0 // too small for label
		}
		if fontSize > 0 {
			lines := strings.Split(label, "\n")
			lineH := fontSize * 1.2
			startY := y + h/2 - float64(len(lines)-1)*lineH/2
			for i, line := range lines {
				fmt.Fprintf(&out, `<text x="%.1f" y="%.1f" fill="%s" font-family="system-ui,sans-serif" font-size="%.0f" text-anchor="middle" dominant-baseline="middle">%s</text>`+"\n",
					x+w/2, startY+float64(i)*lineH, textColor(item), fontSize, escapeXML(line))
			}
		}
	}

	out.WriteString("</svg>\n")

	if err := os.WriteFile(outPath, []byte(out.String()), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("wrote %s (%.0f×%.0f px)\n", outPath, svgW, svgH)
}

func frameColors(item placed) (fill, stroke string, strokeWidth float64) {
	if item.ID == "main" {
		return "#4a90d9", "#1a4a7a", 3
	}
	switch item.Shape {
	case "circle":
		return "#d4edda", "#2d6a3e", 2
	case "ellipse":
		if item.Width > item.Height {
			return "#fce4ec", "#c2185b", 2
		}
		return "#fff3cd", "#856404", 2
	case "square":
		return "#f0e6ff", "#5a3d8a", 2
	default:
		if item.Width > item.Height {
			return "#d4e8f7", "#2b6cb0", 2
		}
		return "#ffe5d9", "#c45c26", 2
	}
}

func textColor(item placed) string {
	if item.ID == "main" {
		return "#fff"
	}
	return "#333"
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
