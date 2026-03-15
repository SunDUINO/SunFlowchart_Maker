package ui

// icons.go — draws toolbar icons as pure geometry, no font required.
// Each icon fits in a ~20x20 box centred at (cx, cy).

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawToolIcon draws the icon for a given Tool at centre (cx,cy) in the given colour.
func drawToolIcon(dst *ebiten.Image, t Tool, cx, cy float64, col color.RGBA) {
	switch t {
	case ToolSelect:
		drawIconSelect(dst, cx, cy, col)
	case ToolNode:
		drawIconNode(dst, cx, cy, col)
	case ToolTitle:
		drawIconTitle(dst, cx, cy, col)
	case ToolConnect:
		drawIconConnect(dst, cx, cy, col)
	case ToolPan:
		drawIconPan(dst, cx, cy, col)
	case ToolDelete:
		drawIconDelete(dst, cx, cy, col)
	case ToolClear:
		drawIconClear(dst, cx, cy, col)
	case ToolLoad:
		drawIconLoad(dst, cx, cy, col)
	}
}

// ── Select: classic arrow cursor ─────────────────────────────────────────────
func drawIconSelect(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	// Arrow: shaft going bottom-left to top-right, with arrowhead
	x, y := float32(cx), float32(cy)
	var p vector.Path
	// Cursor body — pointing up-left
	p.MoveTo(x-7, y+8)  // bottom of shaft
	p.LineTo(x-7, y-7)  // top-left (tip)
	p.LineTo(x+8, y+8)  // bottom-right
	p.LineTo(x+2, y+4)  // inner notch
	p.LineTo(x+5, y+10) // right barb
	p.LineTo(x+2, y+10)
	p.LineTo(x-1, y+4)  // back to notch
	p.Close()
	strokePath(dst, &p, 1.5, col)
}

// ── Node: rounded rectangle with a + inside ──────────────────────────────────
func drawIconNode(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	// Rounded rect outline
	var p vector.Path
	addRoundRect(&p, x-9, y-6, 18, 12, 3)
	strokePath(dst, &p, 1.5, col)
	// Plus sign inside
	vector.StrokeLine(dst, x, y-3, x, y+3, 1.5, col, true)
	vector.StrokeLine(dst, x-3, y, x+3, y, 1.5, col, true)
}

// ── Connect: line with arrowhead ─────────────────────────────────────────────
func drawIconConnect(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	// Small node at start
	var ps vector.Path
	addRoundRect(&ps, x-10, y-4, 8, 8, 2)
	strokePath(dst, &ps, 1.2, col)
	// Small node at end
	var pe vector.Path
	addRoundRect(&pe, x+2, y-4, 8, 8, 2)
	strokePath(dst, &pe, 1.2, col)
	// Arrow line connecting them
	vector.StrokeLine(dst, x-2, y, x+2, y, 1.5, col, true)
	// Arrowhead pointing right
	var pa vector.Path
	pa.MoveTo(x+2, y)
	pa.LineTo(x-1, y-2.5)
	pa.LineTo(x-1, y+2.5)
	pa.Close()
	fillPath(dst, &pa, col)
}

// ── Pan: open hand with four direction arrows ─────────────────────────────────
func drawIconPan(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	// Four arrows: up, down, left, right
	// Centre dot
	vector.FillRect(dst, x-1, y-1, 2, 2, col, true)
	// Up arrow
	drawSmallArrow(dst, x, y-3, x, y-8, col)
	// Down arrow
	drawSmallArrow(dst, x, y+3, x, y+8, col)
	// Left arrow
	drawSmallArrow(dst, x-3, y, x-8, y, col)
	// Right arrow
	drawSmallArrow(dst, x+3, y, x+8, y, col)
}

func drawSmallArrow(dst *ebiten.Image, sx, sy, ex, ey float32, col color.RGBA) {
	vector.StrokeLine(dst, sx, sy, ex, ey, 1.5, col, true)
	dx, dy := ex-sx, ey-sy
	l := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	if l < 1 {
		return
	}
	ux, uy := dx/l, dy/l
	px, py := -uy, ux
	size := float32(3.5)
	var p vector.Path
	p.MoveTo(ex, ey)
	p.LineTo(ex-ux*size+px*size*0.5, ey-uy*size+py*size*0.5)
	p.LineTo(ex-ux*size-px*size*0.5, ey-uy*size-py*size*0.5)
	p.Close()
	fillPath(dst, &p, col)
}

// ── Delete: X cross ───────────────────────────────────────────────────────────
func drawIconDelete(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	vector.StrokeLine(dst, x-7, y-7, x+7, y+7, 2, col, true)
	vector.StrokeLine(dst, x+7, y-7, x-7, y+7, 2, col, true)
}

// ── Clear: eraser — rectangle being wiped away ────────────────────────────────
func drawIconClear(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	// Small nodes (dots) being erased
	vector.FillRect(dst, x-9, y-6, 5, 5, col, true)
	vector.FillRect(dst, x+1, y-6, 5, 5, col, true)
	vector.FillRect(dst, x-9, y+2, 5, 5, col, true)
	// Eraser rectangle sweeping across
	var p vector.Path
	addRoundRect(&p, x-4, y-9, 14, 18, 2)
	strokePath(dst, &p, 1.5, color.RGBA{col.R, col.G, col.B, 200})
	// Fill eraser semi-transparent
	var pf vector.Path
	addRoundRect(&pf, x-4, y-9, 14, 18, 2)
	fillPath(dst, &pf, color.RGBA{col.R / 4, col.G / 4, col.B / 4, 160})
}

// ── Shape mini-icons ──────────────────────────────────────────────────────────

func drawShapeRect(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	var p vector.Path
	addRoundRect(&p, x-9, y-5, 18, 10, 0)
	strokePath(dst, &p, 1.5, col)
}

func drawShapeRounded(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	var p vector.Path
	addRoundRect(&p, x-9, y-5, 18, 10, 4)
	strokePath(dst, &p, 1.5, col)
}

func drawShapeDiamond(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	var p vector.Path
	p.MoveTo(x, y-7)
	p.LineTo(x+9, y)
	p.LineTo(x, y+7)
	p.LineTo(x-9, y)
	p.Close()
	strokePath(dst, &p, 1.5, col)
}

func drawShapeParallel(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	var p vector.Path
	p.MoveTo(x-7, y-5)
	p.LineTo(x+9, y-5)
	p.LineTo(x+7, y+5)
	p.LineTo(x-9, y+5)
	p.Close()
	strokePath(dst, &p, 1.5, col)
}

func drawShapeOval(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	var p vector.Path
	const steps = 32
	for i := 0; i <= steps; i++ {
		a := float64(i) / float64(steps) * 2 * math.Pi
		px := x + 9*float32(math.Cos(a))
		py := y + 5*float32(math.Sin(a))
		if i == 0 {
			p.MoveTo(px, py)
		} else {
			p.LineTo(px, py)
		}
	}
	p.Close()
	strokePath(dst, &p, 1.5, col)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func addRoundRect(p *vector.Path, x, y, w, h, r float32) {
	if r <= 0 {
		p.MoveTo(x, y)
		p.LineTo(x+w, y)
		p.LineTo(x+w, y+h)
		p.LineTo(x, y+h)
		p.Close()
		return
	}
	const steps = 6
	type corner struct {
		cx, cy     float32
		startAngle float64
	}
	corners := []corner{
		{x + r, y + r, math.Pi},
		{x + w - r, y + r, -math.Pi / 2},
		{x + w - r, y + h - r, 0},
		{x + r, y + h - r, math.Pi / 2},
	}
	first := true
	for _, c := range corners {
		for i := 0; i <= steps; i++ {
			angle := c.startAngle + float64(i)/float64(steps)*math.Pi/2
			px := c.cx + r*float32(math.Cos(angle))
			py := c.cy + r*float32(math.Sin(angle))
			if first {
				p.MoveTo(px, py)
				first = false
			} else {
				p.LineTo(px, py)
			}
		}
	}
	p.Close()
}

func strokePath(dst *ebiten.Image, p *vector.Path, lw float32, col color.RGBA) {
	var cs ebiten.ColorScale
	cs.ScaleWithColor(col)
	vector.StrokePath(dst, p, &vector.StrokeOptions{
		Width:    lw,
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinRound,
	}, &vector.DrawPathOptions{ColorScale: cs, AntiAlias: true})
}

func fillPath(dst *ebiten.Image, p *vector.Path, col color.RGBA) {
	var cs ebiten.ColorScale
	cs.ScaleWithColor(col)
	vector.FillPath(dst, p, &vector.FillOptions{},
		&vector.DrawPathOptions{ColorScale: cs, AntiAlias: true})
}

// ── Title: large T letter with underline ─────────────────────────────────────
func drawIconTitle(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	// Horizontal top bar
	vector.StrokeLine(dst, x-8, y-6, x+8, y-6, 2, col, true)
	// Vertical stem
	vector.StrokeLine(dst, x, y-6, x, y+6, 2, col, true)
	// Underline
	vector.StrokeLine(dst, x-8, y+8, x+8, y+8, 1.5, color.RGBA{col.R, col.G, col.B, 160}, true)
}

// ── Load: folder with arrow ───────────────────────────────────────────────────
func drawIconLoad(dst *ebiten.Image, cx, cy float64, col color.RGBA) {
	x, y := float32(cx), float32(cy)
	// Folder body
	var p vector.Path
	p.MoveTo(x-9, y-2)
	p.LineTo(x-9, y+7)
	p.LineTo(x+9, y+7)
	p.LineTo(x+9, y-4)
	p.LineTo(x+1, y-4)
	p.LineTo(x-1, y-6)
	p.LineTo(x-9, y-6)
	p.Close()
	strokePath(dst, &p, 1.5, col)
	// Arrow pointing down into folder
	vector.StrokeLine(dst, x+3, y-9, x+3, y+1, 1.5, col, true)
	drawSmallArrow(dst, x+3, y-1, x+3, y+4, col)
}