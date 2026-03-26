/*
 * ╔══════════════════════════════════════════════════════════════╗
 * ║         SunFlowChart Maker  v1.2.0                        	  ║
 * ║         Neon-dark flowchart editor — Go + Ebiten             ║
 * ╠══════════════════════════════════════════════════════════════╣
 * ║  Autor / Author:                                             ║
 * ║    Andrzej "Sunriver" Gromczyński                            ║
 * ║    Lothar TeaM                                               ║
 * ╠══════════════════════════════════════════════════════════════╣
 * ║  GitHub  : https://github.com/SunDUINO                       ║
 * ║  Forum   : https://forum.lothar-team.pl/                     ║
 * ║                                                              ║
 * ║  Plik / File: png.go                                         ║
 * ║                                                              ║
 * ║  Licencja / License: MIT                                     ║
 * ║  Rok / Year: 2025-2026                                       ║
 * ╚══════════════════════════════════════════════════════════════╝
 */

package export

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/user/flowchart/internal/model"
	"github.com/user/flowchart/internal/render"
)

const pad = 40.0

func bounds(d *model.Diagram) (minX, minY, maxX, maxY float64) {
	minX, minY = math.MaxFloat64, math.MaxFloat64
	maxX, maxY = -math.MaxFloat64, -math.MaxFloat64
	for _, n := range d.Nodes {
		if n.X < minX {
			minX = n.X
		}
		if n.Y < minY {
			minY = n.Y
		}
		if n.X+n.W > maxX {
			maxX = n.X + n.W
		}
		if n.Y+n.H > maxY {
			maxY = n.Y + n.H
		}
	}
	return
}

// PNG Exporter
func ExportPNG(d *model.Diagram, filename string) (string, error) {
	if len(d.Nodes) == 0 {
		return "", fmt.Errorf("diagram jest pusty")
	}
	minX, minY, maxX, maxY := bounds(d)
	ox := minX - pad
	oy := minY - pad
	w := int(maxX - minX + pad*2 + 1)
	h := int(maxY - minY + pad*2 + 1)
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}

	canvas := ebiten.NewImage(w, h)
	canvas.Fill(d.BgColor)
	render.DrawGrid(canvas, 0, 0, d.BgColor, true)

	for _, n := range d.Nodes {
		n.X -= ox
		n.Y -= oy
	}
	offsets := d.EdgeOffsets()
	for _, e := range d.Edges {
		from := d.NodeByID(e.FromID)
		to := d.NodeByID(e.ToID)
		if from != nil && to != nil {
			render.DrawEdge(canvas, e, from, to, offsets[e.ID], false)
		}
	}
	for _, n := range d.Nodes {
		render.DrawNode(canvas, n)
		cx, cy := n.Centre()
		render.DrawNodeLabel(canvas, int(n.Shape), n.Label, n.Sub, cx, cy, n.TextColor)
	}
	for _, n := range d.Nodes {
		n.X += ox
		n.Y += oy
	}

	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for py := 0; py < h; py++ {
		for px := 0; px < w; px++ {
			rgba.Set(px, py, canvas.At(px, py))
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if err := png.Encode(f, rgba); err != nil {
		return "", err
	}
	return filename, nil
}

// SVG scale factor — 2x makes it larger and clearer
const svgScale = 2.0
// SVG Exporter 
func ExportSVG(d *model.Diagram, filename string) error {
	if len(d.Nodes) == 0 {
		return fmt.Errorf("diagram jest pusty")
	}
	minX, minY, maxX, maxY := bounds(d)
	ox := minX - pad
	oy := minY - pad

	// Base dimensions
	baseW := maxX - minX + pad*2
	baseH := maxY - minY + pad*2

	// Extra height for legend
	legendLines := 0
	if d.Legend != "" {
		legendLines = len(strings.Split(d.Legend, "\n")) + 2
	}
	legendH := float64(legendLines) * 16.0
	if legendH > 0 {
		legendH += 20
	}

	// Scaled output size
	svgW := baseW * svgScale
	svgH := (baseH + legendH) * svgScale

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	bg := d.BgColor
	fmt.Fprintf(f, `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="%s"/>
`, svgW, svgH, svgW, svgH, svgColor(bg))

	// Scale group
	fmt.Fprintf(f, `<g transform="scale(%.1f)">`, svgScale)

	// Arrow marker
	fmt.Fprintln(f, `<defs><marker id="arrow" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
<path d="M0,0 L0,6 L6,3 z" fill="#50a0ff"/></marker></defs>`)

	// Edges
	offsets := d.EdgeOffsets()
	for _, e := range d.Edges {
		from := d.NodeByID(e.FromID)
		to := d.NodeByID(e.ToID)
		if from == nil || to == nil {
			continue
		}
		col := e.Color
		if col.A == 0 {
			col = color.RGBA{80, 160, 255, 190}
		}

		tcx, tcy := to.Centre()
		fcx, fcy := from.Centre()
		sx, sy := from.EdgePoint(tcx, tcy)
		ex, ey := to.EdgePoint(fcx, fcy)

		// Apply lateral offset for parallel edges
		lateralPx := float64(offsets[e.ID]) * 14.0
		dx, dy := ex-sx, ey-sy
		dist := math.Sqrt(dx*dx + dy*dy)
		var px, py float64
		if dist > 0 {
			px, py = -dy/dist, dx/dist
		}
		cp := dist * 0.35
		var c1x, c1y, c2x, c2y float64
		if math.Abs(dx) > math.Abs(dy) {
			c1x = sx + cp + px*lateralPx
			c1y = sy + py*lateralPx
			c2x = ex - cp + px*lateralPx
			c2y = ey + py*lateralPx
		} else {
			c1x = sx + px*lateralPx
			c1y = sy + cp + py*lateralPx
			c2x = ex + px*lateralPx
			c2y = ey - cp + py*lateralPx
		}

		fmt.Fprintf(f, `<path d="M%.1f,%.1f C%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="none" stroke="%s" stroke-width="1.6" marker-end="url(#arrow)"/>`,
			sx-ox, sy-oy, c1x-ox, c1y-oy, c2x-ox, c2y-oy, ex-ox, ey-oy, svgColor(col))
		if e.Label != "" {
			lx := (sx+ex)/2 - ox + (-dy/dist)*12
			ly := (sy+ey)/2 - oy + (dx/dist)*12
			fmt.Fprintf(f, `<text x="%.1f" y="%.1f" fill="#aac4ff" font-size="9" font-family="monospace" text-anchor="middle">%s</text>`,
				lx, ly, svgEscape(e.Label))
		}
		fmt.Fprintln(f)
	}

	// Nodes
	for _, n := range d.Nodes {
		x := n.X - ox
		y := n.Y - oy
		fill := svgColor(n.FillColor)
		stroke := svgColor(n.Color)
		lw := n.BorderThickness()

		switch n.Shape {
		case model.ShapeRect:
			fmt.Fprintf(f, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" stroke="%s" stroke-width="%.1f"/>`,
				x, y, n.W, n.H, fill, stroke, lw)
		case model.ShapeRounded:
			fmt.Fprintf(f, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="9" fill="%s" stroke="%s" stroke-width="%.1f"/>`,
				x, y, n.W, n.H, fill, stroke, lw)
		case model.ShapeDiamond:
			cx, cy := x+n.W/2, y+n.H/2
			fmt.Fprintf(f, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="%s" stroke-width="%.1f"/>`,
				cx, y, x+n.W, cy, cx, y+n.H, x, cy, fill, stroke, lw)
		case model.ShapeOval:
			fmt.Fprintf(f, `<ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f" fill="%s" stroke="%s" stroke-width="%.1f"/>`,
				x+n.W/2, y+n.H/2, n.W/2, n.H/2, fill, stroke, lw)
		case model.ShapeParallel:
			skew := n.H * 0.18
			fmt.Fprintf(f, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="%s" stroke-width="%.1f"/>`,
				x+skew, y, x+n.W, y, x+n.W-skew, y+n.H, x, y+n.H, fill, stroke, lw)
		case model.ShapeTitle:
			// Title: just underline
			fmt.Fprintf(f, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="2"/>`,
				x+4, y+n.H-2, x+n.W-4, y+n.H-2, stroke)
			// Large title text
			tc := svgColor(n.TextColor)
			fmt.Fprintf(f, `<text x="%.1f" y="%.1f" fill="%s" font-size="16" font-family="monospace" font-weight="bold" text-anchor="middle">%s</text>`,
				x+n.W/2, y+n.H/2+6, tc, svgEscape(n.Label))
			fmt.Fprintln(f)
			continue
		default:
			fmt.Fprintf(f, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" stroke="%s" stroke-width="%.1f"/>`,
				x, y, n.W, n.H, fill, stroke, lw)
		}

		tc := svgColor(n.TextColor)
		if n.Sub != "" {
			fmt.Fprintf(f, `<text x="%.1f" y="%.1f" fill="%s" font-size="10" font-family="monospace" text-anchor="middle">%s</text>`,
				x+n.W/2, y+n.H/2, tc, svgEscape(n.Label))
			fmt.Fprintf(f, `<text x="%.1f" y="%.1f" fill="%s" font-size="8" font-family="monospace" text-anchor="middle" opacity="0.8">%s</text>`,
				x+n.W/2, y+n.H/2+12, tc, svgEscape(n.Sub))
		} else {
			fmt.Fprintf(f, `<text x="%.1f" y="%.1f" fill="%s" font-size="10" font-family="monospace" text-anchor="middle" dominant-baseline="middle">%s</text>`,
				x+n.W/2, y+n.H/2, tc, svgEscape(n.Label))
		}
		fmt.Fprintln(f)
	}

	// Legend — below diagram
	if d.Legend != "" {
		ly := baseH + 10
		lx := 8.0
		lw2 := baseW - 16
		// Legend box
		fmt.Fprintf(f, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="rgba(8,12,35,0.9)" stroke="rgba(50,80,160,0.8)" stroke-width="1" rx="4"/>`,
			lx, ly, lw2, legendH-4)
		fmt.Fprintf(f, `<text x="%.1f" y="%.1f" fill="#7890cc" font-size="9" font-family="monospace" font-weight="bold">Legenda</text>`,
			lx+6, ly+12)
		for i, line := range strings.Split(d.Legend, "\n") {
			fmt.Fprintf(f, `<text x="%.1f" y="%.1f" fill="#b0c8ff" font-size="9" font-family="monospace">%s</text>`,
				lx+6, ly+24+float64(i)*14, svgEscape(line))
		}
		fmt.Fprintln(f)
	}

	fmt.Fprintln(f, `</g>`)
	fmt.Fprintln(f, "</svg>")
	return nil
}

func svgColor(c color.RGBA) string {
	return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", c.R, c.G, c.B, float64(c.A)/255)
}

func svgEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
