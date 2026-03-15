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
 * ║  Plik / File: edges.go                                       ║
 * ║                                                              ║
 * ║  Licencja / License: MIT                                     ║
 * ║  Rok / Year: 2025-2026                                       ║
 * ╚══════════════════════════════════════════════════════════════╝
 */

package render


import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/user/flowchart/internal/model"
)

// DrawEdge renders one edge. offset is the lateral separation index for
// parallel edges between the same pair of nodes (0 = centre, ±1, ±2...).
func DrawEdge(dst *ebiten.Image, e *model.Edge, from, to *model.Node, offset int, hover bool) {
	tcx, tcy := to.Centre()
	fcx, fcy := from.Centre()
	sx, sy := from.EdgePoint(tcx, tcy)
	ex, ey := to.EdgePoint(fcx, fcy)

	col := e.Color
	if col.A == 0 {
		col = color.RGBA{80, 160, 255, 190}
	}
	if e.Selected {
		col = color.RGBA{255, 240, 80, 230}
	} else if hover {
		// Brighten on hover
		col = color.RGBA{
			clampCol(int(col.R) + 60),
			clampCol(int(col.G) + 60),
			clampCol(int(col.B) + 60),
			230,
		}
	}

	// Lateral offset to separate parallel edges (alternating ±)
	lateralPx := float64(offset) * 14.0

	switch e.Style {
	case model.EdgeElbow:
		drawElbow(dst, sx, sy, ex, ey, col, lateralPx)
	default:
		drawCurve(dst, sx, sy, ex, ey, col, lateralPx)
	}

	if e.Label != "" {
		lx, ly := labelPos(sx, sy, ex, ey, lateralPx)
		DrawEdgeLabel(dst, e.Label, lx, ly)
	}
}

func DrawEdgePreview(dst *ebiten.Image, from *model.Node, tx, ty float64) {
	sx, sy := from.EdgePoint(tx, ty)
	drawCurve(dst, sx, sy, tx, ty, color.RGBA{120, 200, 255, 120}, 0)
}

// labelPos returns a point offset perpendicular from midpoint.
func labelPos(sx, sy, ex, ey, lateralPx float64) (float64, float64) {
	mx := sx + (ex-sx)*0.45
	my := sy + (ey-sy)*0.45
	dx, dy := ex-sx, ey-sy
	l := math.Sqrt(dx*dx + dy*dy)
	if l > 0 {
		perp := 12.0 + math.Abs(lateralPx)
		mx += (-dy / l) * perp
		my += (dx / l) * perp
	}
	return mx, my
}

func drawCurve(dst *ebiten.Image, sx, sy, ex, ey float64, col color.RGBA, lateralPx float64) {
	dx := ex - sx
	dy := ey - sy
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist < 2 {
		return
	}

	// Perpendicular unit vector for lateral offset
	var px, py float64
	if dist > 0 {
		px, py = -dy/dist, dx/dist
	}

	// Apply lateral offset to control points — creates a bowed curve
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

	glow := color.RGBA{col.R, col.G, col.B, 40}
	bezierStroke(dst, sx, sy, c1x, c1y, c2x, c2y, ex, ey, 4.5, glow)
	bezierStroke(dst, sx, sy, c1x, c1y, c2x, c2y, ex, ey, 1.6, col)
	arrowhead(dst, ex, ey, c2x, c2y, col)
}

func drawElbow(dst *ebiten.Image, sx, sy, ex, ey float64, col color.RGBA, lateralPx float64) {
	dx := ex - sx
	dy := ey - sy
	var px, py float64
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist > 0 {
		px, py = -dy/dist, dx/dist
	}

	var pts [][2]float64
	if math.Abs(ex-sx) >= math.Abs(ey-sy) {
		mx := (sx+ex)/2 + px*lateralPx
		pts = [][2]float64{
			{sx, sy},
			{mx, sy + py*lateralPx},
			{mx, ey + py*lateralPx},
			{ex, ey},
		}
	} else {
		my := (sy+ey)/2 + py*lateralPx
		pts = [][2]float64{
			{sx, sy},
			{sx + px*lateralPx, my},
			{ex + px*lateralPx, my},
			{ex, ey},
		}
	}
	glow := color.RGBA{col.R, col.G, col.B, 40}
	strokePolyline(dst, pts, 4.5, glow)
	strokePolyline(dst, pts, 1.6, col)
	last := pts[len(pts)-2]
	arrowhead(dst, ex, ey, last[0], last[1], col)
}

func strokePolyline(dst *ebiten.Image, pts [][2]float64, lw float32, col color.RGBA) {
	for i := 0; i < len(pts)-1; i++ {
		a, b := pts[i], pts[i+1]
		var p vector.Path
		p.MoveTo(float32(a[0]), float32(a[1]))
		p.LineTo(float32(b[0]), float32(b[1]))
		vector.StrokePath(dst, &p, &vector.StrokeOptions{
			Width:    lw,
			LineCap:  vector.LineCapRound,
			LineJoin: vector.LineJoinRound,
		}, drawOpts(col))
	}
}

func clampCol(v int) uint8 {
	if v > 255 { return 255 }
	if v < 0 { return 0 }
	return uint8(v)
}