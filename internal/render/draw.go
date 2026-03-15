package render

// draw.go — Ebiten v2.9.x
// FillPath/StrokePath use *DrawPathOptions{ColorScale, AntiAlias}

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/user/flowchart/internal/model"
)

const glowLayers = 4

func drawOpts(col color.RGBA) *vector.DrawPathOptions {
	var cs ebiten.ColorScale
	cs.ScaleWithColor(col)
	return &vector.DrawPathOptions{ColorScale: cs, AntiAlias: true}
}

// ─── Node ────────────────────────────────────────────────────────────────────

func DrawNode(dst *ebiten.Image, n *model.Node) {
	alpha := n.GlowAlpha()
	// Glow layers — size proportional, alpha fades outward
	for i := glowLayers; i >= 1; i-- {
		exp := float32(i) * 2.4
		// Keep alpha strong enough to be visible
		a := uint8(float64(alpha) * float64(i) / float64(glowLayers))
		drawShape(dst, n, exp, color.RGBA{n.Color.R, n.Color.G, n.Color.B, a}, color.RGBA{}, false)
	}
	drawShape(dst, n, 0, n.Color, n.FillColor, true)
	if n.Selected {
		sel := color.RGBA{255, 240, 80, 220}
		drawShape(dst, n, -1, sel, color.RGBA{}, false)
		drawShape(dst, n, 2.5, sel, color.RGBA{}, false)
	}
	if n.Anim == model.AnimSpinner {
		drawSpinner(dst, n)
	}
}

func drawShape(dst *ebiten.Image, n *model.Node, expand float32, border, fill color.RGBA, doFill bool) {
	x := float32(n.X) - expand
	y := float32(n.Y) - expand
	w := float32(n.W) + expand*2
	h := float32(n.H) + expand*2
	// Use actual border thickness for the main stroke, thin for glow/selection rings
	var lw float32
	if expand == 0 {
		lw = n.BorderThickness()
	} else {
		lw = 1.0
	}

	switch n.Shape {
	case model.ShapeRect:
		if doFill {
			vector.FillRect(dst, x, y, w, h, fill, true)
		}
		vector.StrokeRect(dst, x, y, w, h, lw, border, true)
	case model.ShapeRounded:
		r := float32(9) + expand/2
		if r < 3 {
			r = 3
		}
		if doFill {
			fillRoundRect(dst, x, y, w, h, r, fill)
		}
		strokeRoundRect(dst, x, y, w, h, r, lw, border)
	case model.ShapeDiamond:
		if doFill {
			fillDiamond(dst, x, y, w, h, fill)
		}
		strokeDiamond(dst, x, y, w, h, lw, border)
	case model.ShapeParallel:
		skew := h * 0.18
		if doFill {
			fillPara(dst, x, y, w, h, skew, fill)
		}
		strokePara(dst, x, y, w, h, skew, lw, border)
	case model.ShapeOval:
		if doFill {
			fillEllipse(dst, x+w/2, y+h/2, w/2, h/2, fill)
		}
		strokeEllipse(dst, x+w/2, y+h/2, w/2, h/2, lw, border)
	case model.ShapeTitle:
		// Title node: subtle underline only, no fill, no box
		if doFill {
			// very faint bg
			vector.FillRect(dst, x, y, w, h, color.RGBA{border.R / 12, border.G / 12, border.B / 12, 80}, true)
		}
		// underline
		vector.StrokeLine(dst, x+4, y+h-2, x+w-4, y+h-2, lw, border, true)
	}
}

func drawSpinner(dst *ebiten.Image, n *model.Node) {
	cx := float32(n.X + n.W/2)
	cy := float32(n.Y + n.H/2)
	r := float32(math.Min(n.W, n.H)/2) + 8
	start := n.SpinnerAngle()
	arcLen := math.Pi * 1.4
	var path vector.Path
	for i := 0; i <= 48; i++ {
		angle := start + float64(i)/48.0*arcLen
		px := cx + r*float32(math.Cos(angle))
		py := cy + r*float32(math.Sin(angle))
		if i == 0 {
			path.MoveTo(px, py)
		} else {
			path.LineTo(px, py)
		}
	}
	// Gradient-like: draw twice, thick+faded then thin+bright
	vector.StrokePath(dst, &path,
		&vector.StrokeOptions{Width: 5, LineCap: vector.LineCapRound},
		drawOpts(color.RGBA{n.Color.R, n.Color.G, n.Color.B, 60}))
	vector.StrokePath(dst, &path,
		&vector.StrokeOptions{Width: 2.5, LineCap: vector.LineCapRound},
		drawOpts(color.RGBA{n.Color.R, n.Color.G, n.Color.B, 230}))
}

// ─── Grid ────────────────────────────────────────────────────────────────────

const GridSpacing = 24.0

func DrawGrid(dst *ebiten.Image, offX, offY float64, bgColor color.RGBA, showGrid bool) {
	dst.Fill(bgColor)
	if !showGrid {
		return
	}
	// Dot colour: slightly lighter than bg
	dot := color.RGBA{
		clampU8(int(bgColor.R) + 15),
		clampU8(int(bgColor.G) + 15),
		clampU8(int(bgColor.B) + 20),
		180,
	}
	if bgColor.R > 180 {
		dot = color.RGBA{180, 180, 190, 120}
	}
	w := float64(dst.Bounds().Dx())
	h := float64(dst.Bounds().Dy())
	sp := GridSpacing
	startX := math.Mod(offX, sp)
	startY := math.Mod(offY, sp)
	for x := startX; x < w; x += sp {
		for y := startY; y < h; y += sp {
			vector.FillRect(dst, float32(x)-0.5, float32(y)-0.5, 1.5, 1.5, dot, false)
		}
	}
}

func clampU8(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

func DrawSelectionRect(dst *ebiten.Image, x1, y1, x2, y2 float64) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	vector.FillRect(dst, float32(x1), float32(y1), float32(x2-x1), float32(y2-y1),
		color.RGBA{80, 150, 255, 18}, false)
	vector.StrokeRect(dst, float32(x1), float32(y1), float32(x2-x1), float32(y2-y1),
		1.4, color.RGBA{80, 150, 255, 180}, false)
}

// ─── Resize handle ────────────────────────────────────────────────────────────

func DrawResizeHandle(dst *ebiten.Image, n *model.Node) {
	const s float32 = 10
	x := float32(n.X+n.W) - s/2
	y := float32(n.Y+n.H) - s/2
	vector.FillRect(dst, x, y, s, s, color.RGBA{255, 240, 80, 230}, true)
	for i := 0; i < 3; i++ {
		o := float32(i) * 3
		vector.StrokeLine(dst, x+2+o, y+s-2, x+s-2, y+2+o, 1, color.RGBA{20, 20, 40, 200}, true)
	}
}

// ─── Shape helpers ───────────────────────────────────────────────────────────

func fillRoundRect(dst *ebiten.Image, x, y, w, h, r float32, col color.RGBA) {
	vector.FillPath(dst, rrectPath(x, y, w, h, r), &vector.FillOptions{}, drawOpts(col))
}

func strokeRoundRect(dst *ebiten.Image, x, y, w, h, r, lw float32, col color.RGBA) {
	vector.StrokePath(dst, rrectPath(x, y, w, h, r),
		&vector.StrokeOptions{Width: lw, LineJoin: vector.LineJoinRound}, drawOpts(col))
}

func rrectPath(x, y, w, h, r float32) *vector.Path {
	if r > w/2 {
		r = w / 2
	}
	if r > h/2 {
		r = h / 2
	}
	const steps = 8
	p := &vector.Path{}
	corners := []struct {
		cx, cy     float32
		startAngle float64
	}{
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
	return p
}

func fillDiamond(dst *ebiten.Image, x, y, w, h float32, col color.RGBA) {
	cx, cy := x+w/2, y+h/2
	p := &vector.Path{}
	p.MoveTo(cx, y); p.LineTo(x+w, cy); p.LineTo(cx, y+h); p.LineTo(x, cy)
	p.Close()
	vector.FillPath(dst, p, &vector.FillOptions{}, drawOpts(col))
}

func strokeDiamond(dst *ebiten.Image, x, y, w, h, lw float32, col color.RGBA) {
	cx, cy := x+w/2, y+h/2
	p := &vector.Path{}
	p.MoveTo(cx, y); p.LineTo(x+w, cy); p.LineTo(cx, y+h); p.LineTo(x, cy)
	p.Close()
	vector.StrokePath(dst, p, &vector.StrokeOptions{Width: lw, LineJoin: vector.LineJoinMiter}, drawOpts(col))
}

func fillPara(dst *ebiten.Image, x, y, w, h, skew float32, col color.RGBA) {
	p := &vector.Path{}
	p.MoveTo(x+skew, y); p.LineTo(x+w, y); p.LineTo(x+w-skew, y+h); p.LineTo(x, y+h)
	p.Close()
	vector.FillPath(dst, p, &vector.FillOptions{}, drawOpts(col))
}

func strokePara(dst *ebiten.Image, x, y, w, h, skew, lw float32, col color.RGBA) {
	p := &vector.Path{}
	p.MoveTo(x+skew, y); p.LineTo(x+w, y); p.LineTo(x+w-skew, y+h); p.LineTo(x, y+h)
	p.Close()
	vector.StrokePath(dst, p, &vector.StrokeOptions{Width: lw, LineJoin: vector.LineJoinMiter}, drawOpts(col))
}

func fillEllipse(dst *ebiten.Image, cx, cy, rx, ry float32, col color.RGBA) {
	vector.FillPath(dst, ellipsePath(cx, cy, rx, ry), &vector.FillOptions{}, drawOpts(col))
}

func strokeEllipse(dst *ebiten.Image, cx, cy, rx, ry, lw float32, col color.RGBA) {
	vector.StrokePath(dst, ellipsePath(cx, cy, rx, ry), &vector.StrokeOptions{Width: lw}, drawOpts(col))
}

func ellipsePath(cx, cy, rx, ry float32) *vector.Path {
	p := &vector.Path{}
	const steps = 64
	for i := 0; i <= steps; i++ {
		a := float64(i) / float64(steps) * 2 * math.Pi
		px := cx + rx*float32(math.Cos(a))
		py := cy + ry*float32(math.Sin(a))
		if i == 0 {
			p.MoveTo(px, py)
		} else {
			p.LineTo(px, py)
		}
	}
	p.Close()
	return p
}

// ─── Exported round rect helpers (used by ui package) ────────────────────────

func FillRoundRect(dst *ebiten.Image, x, y, w, h, r float32, col color.RGBA) {
	fillRoundRect(dst, x, y, w, h, r, col)
}

func StrokeRoundRect(dst *ebiten.Image, x, y, w, h, r, lw float32, col color.RGBA) {
	strokeRoundRect(dst, x, y, w, h, r, lw, col)
}

// shared bezier stroke used by edges.go
func bezierStroke(dst *ebiten.Image, sx, sy, c1x, c1y, c2x, c2y, ex, ey float64, lw float32, col color.RGBA) {
	var p vector.Path
	p.MoveTo(float32(sx), float32(sy))
	p.CubicTo(float32(c1x), float32(c1y), float32(c2x), float32(c2y), float32(ex), float32(ey))
	vector.StrokePath(dst, &p, &vector.StrokeOptions{
		Width: lw, LineCap: vector.LineCapRound, LineJoin: vector.LineJoinRound,
	}, drawOpts(col))
}

func arrowhead(dst *ebiten.Image, tx, ty, fromX, fromY float64, col color.RGBA) {
	dx, dy := tx-fromX, ty-fromY
	l := math.Sqrt(dx*dx + dy*dy)
	if l < 1 {
		return
	}
	ux, uy := dx/l, dy/l
	px, py := -uy, ux
	size, wing := 10.0, 4.5
	var p vector.Path
	p.MoveTo(float32(tx), float32(ty))
	p.LineTo(float32(tx-ux*size+px*wing), float32(ty-uy*size+py*wing))
	p.LineTo(float32(tx-ux*size-px*wing), float32(ty-uy*size-py*wing))
	p.Close()
	vector.FillPath(dst, &p, &vector.FillOptions{}, drawOpts(col))
}

// ─── Connect tool highlights ──────────────────────────────────────────────────

// DrawConnectSource draws a pulsing ring around the source node when connecting.
func DrawConnectSource(dst *ebiten.Image, n *model.Node, t float64) {
	pulse := float32((math.Sin(t*6)+1)/2) * 6
	exp := float32(6) + pulse
	col := color.RGBA{80, 220, 120, 200}
	// Two rings
	StrokeRoundRect(dst,
		float32(n.X)-exp, float32(n.Y)-exp,
		float32(n.W)+exp*2, float32(n.H)+exp*2,
		12+exp, 2, col)
	StrokeRoundRect(dst,
		float32(n.X)-exp*0.5, float32(n.Y)-exp*0.5,
		float32(n.W)+exp, float32(n.H)+exp,
		10+exp*0.5, 1, color.RGBA{col.R, col.G, col.B, 120})
}

// DrawConnectHover draws a hover ring showing a node can be clicked as target.
func DrawConnectHover(dst *ebiten.Image, n *model.Node) {
	col := color.RGBA{255, 220, 60, 200}
	StrokeRoundRect(dst,
		float32(n.X)-8, float32(n.Y)-8,
		float32(n.W)+16, float32(n.H)+16,
		12, 2, col)
	// Corner dots at four corners
	corners := [][2]float32{
		{float32(n.X) - 3, float32(n.Y) - 3},
		{float32(n.X+n.W) - 3, float32(n.Y) - 3},
		{float32(n.X) - 3, float32(n.Y+n.H) - 3},
		{float32(n.X+n.W) - 3, float32(n.Y+n.H) - 3},
	}
	for _, c := range corners {
		vector.FillRect(dst, c[0], c[1], 6, 6, col, true)
	}
}