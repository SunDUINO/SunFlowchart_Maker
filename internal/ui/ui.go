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
 * ║  Plik / File: ui.go                                          ║
 * ║                                                              ║
 * ║  Licencja / License: MIT                                     ║
 * ║  Rok / Year: 2025-2026                                       ║
 * ╚══════════════════════════════════════════════════════════════╝
 */

/*
 * ╔══════════════════════════════════════════════════════════════╗
 * ║         SunFlowChart Maker  v1.2.0                          ║
 * ║         Neon-dark flowchart editor — Go + Ebiten            ║
 * ╠══════════════════════════════════════════════════════════════╣
 * ║  Autor / Author:                                            ║
 * ║    Andrzej "Sunriver" Gromczyński                           ║
 * ║    Lothar TeaM                                              ║
 * ╠══════════════════════════════════════════════════════════════╣
 * ║  GitHub  : https://github.com/SunDUINO                      ║
 * ║  Forum   : https://forum.lothar-team.pl/                    ║
 * ╠══════════════════════════════════════════════════════════════╣
 * ║  Plik / File: ui.go                                        ║
 * ║  Opis / Desc: Interfejs / User interface panels            ║
 * ╠══════════════════════════════════════════════════════════════╣
 * ║  Licencja / License: MIT                                    ║
 * ║  Rok / Year: 2025-2026                                      ║
 * ╚══════════════════════════════════════════════════════════════╝
 */
package ui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/user/flowchart/internal/model"
	"github.com/user/flowchart/internal/render"
)

// ─── Tool enum ────────────────────────────────────────────────────────────────

type Tool int

const (
	ToolSelect Tool = iota
	ToolNode
	ToolTitle
	ToolConnect
	ToolPan
	ToolDelete
	ToolClear
	ToolLoad
)

const (
	SidebarW   = 68.0
	PropW      = 230.0
	BtnSize    = 46.0
	BtnPad     = 5.0
	StatusbarH = 22.0
)

// ─── Toolbar ─────────────────────────────────────────────────────────────────

type toolMeta struct {
	tool    Tool
	label   string // short label below icon
	tooltip string
	color   color.RGBA
}

var toolMetas = []toolMeta{
	{ToolSelect, "Wybierz", "Zaznacz i przesuń  [V]", color.RGBA{40, 80, 180, 255}},
	{ToolNode, "Węzeł", "Dodaj węzeł  [N]", color.RGBA{30, 130, 60, 255}},
	{ToolTitle, "Tytuł", "Dodaj tytuł  [T]", color.RGBA{160, 120, 20, 255}},
	{ToolConnect, "Połącz", "Narysuj strzałkę  [C]", color.RGBA{30, 100, 150, 255}},
	{ToolPan, "Przesuń", "Przesuń widok  [H/Spacja]", color.RGBA{80, 60, 150, 255}},
	{ToolDelete, "Usuń", "Usuń element  [X/Del]", color.RGBA{160, 30, 30, 255}},
	{ToolClear, "Czysc", "Wyczysc ekran", color.RGBA{140, 80, 20, 255}},
	{ToolLoad, "Wczytaj", "Wczytaj JSON  [Ctrl+O]", color.RGBA{20, 110, 100, 255}},
}

type Toolbar struct {
	ActiveTool  Tool
	ActiveShape model.Shape
	HoverTool   Tool
	hoverSet    bool
	ShowGrid    bool // mirrored from Game for drawing
	SnapGrid    bool
}

func NewToolbar() *Toolbar { return &Toolbar{HoverTool: -1} }

func (tb *Toolbar) Draw(dst *ebiten.Image) {
	h := float64(dst.Bounds().Dy())
	vector.FillRect(dst, 0, 0, SidebarW, float32(h), color.RGBA{7, 10, 26, 252}, false)
	vector.StrokeRect(dst, SidebarW-1, 0, 1, float32(h), 1, color.RGBA{25, 45, 100, 255}, false)

	// Logo badge
	render.FillRoundRect(dst, 6, 6, SidebarW-12, 30, 6, color.RGBA{20, 50, 130, 220})
	render.DrawTextCentered(dst, "FC", SidebarW/2, 21, color.RGBA{120, 190, 255, 255})

	// Tool buttons
	for i, tm := range toolMetas {
		y := float64(44 + i*int(BtnSize+BtnPad))
		active := tm.tool == tb.ActiveTool
		tb.drawToolBtn(dst, 6, y, BtnSize, BtnSize, tm, active)
	}

	// Separator before shape picker
	sepY := float32(44 + len(toolMetas)*int(BtnSize+BtnPad) + 4)
	vector.FillRect(dst, 6, sepY, SidebarW-12, 1, color.RGBA{35, 55, 120, 200}, false)

	// Shape picker
	if tb.ActiveTool == ToolNode {
		render.DrawTextCentered(dst, "Kształt", SidebarW/2, float64(sepY)+12,
			color.RGBA{80, 110, 180, 200})
		shapeDrawers := []struct {
			s model.Shape
			f func(*ebiten.Image, float64, float64, color.RGBA)
		}{
			{model.ShapeRect, drawShapeRect},
			{model.ShapeRounded, drawShapeRounded},
			{model.ShapeDiamond, drawShapeDiamond},
			{model.ShapeParallel, drawShapeParallel},
			{model.ShapeOval, drawShapeOval},
		}
		for j, sd := range shapeDrawers {
			by := float64(sepY) + 18 + float64(j)*38
			active := tb.ActiveShape == sd.s
			col := color.RGBA{40, 80, 160, 255}
			bg := color.RGBA{10, 18, 50, 200}
			if active {
				bg = col
			}
			render.FillRoundRect(dst, 6, float32(by), BtnSize, 32, 5, bg)
			render.StrokeRoundRect(dst, 6, float32(by), BtnSize, 32, 5, 1.2,
				color.RGBA{50, 100, 200, active2alpha(active)})
			ic := color.RGBA{160, 200, 255, 220}
			if active {
				ic = color.RGBA{255, 255, 255, 255}
			}
			sd.f(dst, SidebarW/2, by+16, ic)
		}
	}

	// Grid/Snap toggle buttons — above status bar
	screenH := float64(dst.Bounds().Dy())
	btnH := 36.0
	gap := 4.0
	snapY := screenH - StatusbarH - btnH - 4
	gridY := snapY - btnH - gap
	drawGridToggleBtn(dst, 6, gridY, BtnSize, btnH, "Siatka", tb.ShowGrid, color.RGBA{30, 80, 120, 255})
	drawGridToggleBtn(dst, 6, snapY, BtnSize, btnH, "Snap", tb.SnapGrid, color.RGBA{80, 60, 120, 255})
}

func active2alpha(active bool) uint8 {
	if active {
		return 255
	}
	return 130
}

func (tb *Toolbar) drawToolBtn(dst *ebiten.Image, x, y, w, h float64, tm toolMeta, active bool) {
	base := tm.color
	bg := color.RGBA{base.R / 5, base.G / 5, base.B / 5, 210}
	border := color.RGBA{base.R / 2, base.G / 2, base.B / 2, 160}
	iconCol := color.RGBA{160, 190, 255, 200}

	if active {
		bg = base
		border = color.RGBA{
			clampC(int(base.R) + 70),
			clampC(int(base.G) + 70),
			clampC(int(base.B) + 70),
			255,
		}
		iconCol = color.RGBA{255, 255, 255, 255}
		// glow rings
		for i := 3; i >= 1; i-- {
			e := float32(i) * 2.5
			a := uint8(22 - uint8(i)*6)
			render.StrokeRoundRect(dst,
				float32(x)-e, float32(y)-e, float32(w)+e*2, float32(h)+e*2,
				8+e, 1, color.RGBA{border.R, border.G, border.B, a})
		}
	}

	// Button background
	render.FillRoundRect(dst, float32(x), float32(y), float32(w), float32(h), 7, bg)
	render.StrokeRoundRect(dst, float32(x), float32(y), float32(w), float32(h), 7, 1.4, border)

	// Geometric icon in upper 2/3 of button
	iconCX := x + w/2
	iconCY := y + h*0.42
	drawToolIcon(dst, tm.tool, iconCX, iconCY, iconCol)

	// Short label in lower part
	labelCol := color.RGBA{100, 130, 200, 180}
	if active {
		labelCol = color.RGBA{255, 255, 255, 220}
	}
	render.DrawTextCentered(dst, tm.label, x+w/2, y+h-7, labelCol)
}

func drawGridToggleBtn(dst *ebiten.Image, x, y, w, h float64, label string, active bool, base color.RGBA) {
	bg := color.RGBA{base.R / 4, base.G / 4, base.B / 4, 210}
	border := color.RGBA{base.R / 2, base.G / 2, base.B / 2, 160}
	tc := color.RGBA{120, 150, 210, 180}
	if active {
		bg = base
		border = color.RGBA{clampC(int(base.R) + 60), clampC(int(base.G) + 60), clampC(int(base.B) + 60), 255}
		tc = color.RGBA{255, 255, 255, 240}
	}
	render.FillRoundRect(dst, float32(x), float32(y), float32(w), float32(h), 6, bg)
	render.StrokeRoundRect(dst, float32(x), float32(y), float32(w), float32(h), 6, 1.2, border)
	// LED indicator dot
	dotCol := color.RGBA{80, 200, 80, 220}
	if !active {
		dotCol = color.RGBA{180, 60, 60, 180}
	}
	vector.FillCircle(dst, float32(x+w-8), float32(y+8), 4, dotCol, true)
	render.DrawTextCentered(dst, label, x+w/2, y+h/2+2, tc)
}

func (tb *Toolbar) Click(mx, my float64) bool {
	if mx >= SidebarW {
		return false
	}
	// Don't swallow clicks in the grid/snap button zone at the bottom
	// (those are handled separately by ClickGridBtn)
	for i, tm := range toolMetas {
		y := float64(44 + i*int(BtnSize+BtnPad))
		if inRect(mx, my, 10, y, BtnSize, BtnSize) {
			tb.ActiveTool = tm.tool
			return true
		}
	}
	if tb.ActiveTool == ToolNode {
		shapes := []model.Shape{model.ShapeRect, model.ShapeRounded, model.ShapeDiamond, model.ShapeParallel, model.ShapeOval}
		sepY := float64(44 + len(toolMetas)*int(BtnSize+BtnPad) + 4)
		for j, sh := range shapes {
			by := sepY + 18 + float64(j)*38
			if inRect(mx, my, 6, by, BtnSize, 32) {
				tb.ActiveShape = sh
				return true
			}
		}
	}
	// Only consume if NOT in the grid/snap zone — let ClickGridBtn handle those
	return false
}

func (tb *Toolbar) ClickGridBtn(mx, my float64, screenH float64) (gridToggle, snapToggle bool) {
	if mx >= SidebarW {
		return false, false
	}
	btnH := 36.0
	gap := 4.0
	snapY := screenH - StatusbarH - btnH - 4
	gridY := snapY - btnH - gap
	if inRect(mx, my, 6, gridY, BtnSize, btnH) {
		return true, false
	}
	if inRect(mx, my, 6, snapY, BtnSize, btnH) {
		return false, true
	}
	return false, false
}

// ─── Properties panel ────────────────────────────────────────────────────────

type PropPanel struct {
	X, Y, W, H  float64
	Visible     bool
	colSwatches []swatch
	bgSwatches  []swatch
	animBtns    []animBtn
	styleBtns   []styleBtn
	borderBtns  []borderBtn
	// drag state
	dragging bool
	dragOffX float64
	dragOffY float64
}

type swatch struct {
	x, y, w, h float64
	col        color.RGBA
}
type animBtn struct {
	x, y, w, h float64
	anim       model.Anim
	label      string
}
type styleBtn struct {
	x, y, w, h float64
	style      model.EdgeStyle
	label      string
}
type borderBtn struct {
	x, y, w, h float64
	thick      float32
	label      string
}

func NewPropPanel(screenW float64) *PropPanel {
	px := screenW - PropW - 4
	pp := &PropPanel{X: px, Y: 60, W: PropW, H: 580}
	pp.rebuild()
	return pp
}

// rebuild recalculates all element positions from current pp.X / pp.Y.
// Call this after dragging the panel.
func (pp *PropPanel) rebuild() {
	px := pp.X
	sw, sg := 24.0, 5.0
	cols := 4

	pp.colSwatches = pp.colSwatches[:0]
	for i, c := range model.Palette {
		pp.colSwatches = append(pp.colSwatches, swatch{
			x: px + 10 + float64(i%cols)*(sw+sg),
			y: pp.Y + 80 + float64(i/cols)*(sw+sg),
			w: sw, h: sw, col: c,
		})
	}

	pp.bgSwatches = pp.bgSwatches[:0]
	for i, c := range model.BgPalette {
		pp.bgSwatches = append(pp.bgSwatches, swatch{
			x: px + 10 + float64(i%cols)*(sw+sg),
			y: pp.Y + 210 + float64(i/cols)*(sw+sg),
			w: sw, h: sw, col: c,
		})
	}

	anims := []struct {
		a model.Anim
		l string
	}{
		{model.AnimNone, "Off"}, {model.AnimGlow, "Glow"}, {model.AnimPulse, "Pulse"},
		{model.AnimBlink, "Blink"}, {model.AnimFlash, "Flash"}, {model.AnimSpinner, "Spin"},
	}
	aw, ah, ag := 56.0, 24.0, 5.0
	pp.animBtns = pp.animBtns[:0]
	for i, a := range anims {
		pp.animBtns = append(pp.animBtns, animBtn{
			x: px + 10 + float64(i%3)*(aw+ag), y: pp.Y + 310 + float64(i/3)*(ah+ag),
			w: aw, h: ah, anim: a.a, label: a.l,
		})
	}

	borders := []struct {
		t float32
		l string
	}{{1, "1px"}, {2, "2px"}, {3, "3px"}, {4, "4px"}}
	bw, bh, bg2 := 42.0, 24.0, 5.0
	pp.borderBtns = pp.borderBtns[:0]
	for i, b := range borders {
		pp.borderBtns = append(pp.borderBtns, borderBtn{
			x: px + 10 + float64(i)*(bw+bg2), y: pp.Y + 385,
			w: bw, h: bh, thick: b.t, label: b.l,
		})
	}

	pp.styleBtns = []styleBtn{
		{x: px + 10, y: pp.Y + 430, w: 90, h: 24, style: model.EdgeCurve, label: "~ Krzywa"},
		{x: px + 115, y: pp.Y + 430, w: 90, h: 24, style: model.EdgeElbow, label: "| Lamana"},
	}
}

func (pp *PropPanel) Draw(dst *ebiten.Image, n *model.Node, t float64) {
	if !pp.Visible || n == nil {
		return
	}
	pp.drawBg(dst, pp.H)

	// Drag title bar
	pp.drawTitleBar(dst, "Wlasciwosci wezla")

	render.DrawTextAt(dst, "Kolor:", pp.X+10, pp.Y+50, color.RGBA{80, 110, 180, 200})
	vector.FillRect(dst, float32(pp.X+60), float32(pp.Y+44), 28, 16, n.Color, true)
	vector.StrokeRect(dst, float32(pp.X+60), float32(pp.Y+44), 28, 16, 1, color.RGBA{200, 200, 255, 160}, true)
	pp.drawSwatches(dst, pp.colSwatches, n.Color)

	render.DrawTextAt(dst, "Tlo:", pp.X+10, pp.Y+178, color.RGBA{80, 110, 180, 200})
	pp.drawSwatches(dst, pp.bgSwatches, color.RGBA{})

	render.DrawTextAt(dst, "Animacja:", pp.X+10, pp.Y+278, color.RGBA{80, 110, 180, 200})
	phase := math.Mod(t*n.AnimSpeed, 1.0)
	pv := (math.Sin(phase*2*math.Pi) + 1) / 2
	dotR := float32(3 + pv*4)
	dotA := uint8(60 + pv*190)
	vector.FillCircle(dst, float32(pp.X+pp.W-14), float32(pp.Y+283), dotR,
		color.RGBA{n.Color.R, n.Color.G, n.Color.B, dotA}, true)
	pp.drawAnimBtns(dst, n.Anim)

	render.DrawTextAt(dst, "Ramka:", pp.X+10, pp.Y+358, color.RGBA{80, 110, 180, 200})
	pp.drawBorderBtns(dst, n.BorderThickness())

	render.DrawTextAt(dst, "Ksztalt: "+n.Shape.String(), pp.X+10, pp.Y+420, color.RGBA{80, 110, 180, 200})

	render.DrawTextAt(dst, "F2=etykieta  Tab=podtytul", pp.X+10, pp.Y+545, color.RGBA{55, 80, 130, 180})
	render.DrawTextAt(dst, "Del=usun  Esc=odznacz", pp.X+10, pp.Y+558, color.RGBA{55, 80, 130, 180})
}

func (pp *PropPanel) DrawEdgeProps(dst *ebiten.Image, e *model.Edge) {
	if !pp.Visible || e == nil {
		return
	}
	pp.drawBg(dst, 260)
	pp.drawTitleBar(dst, "Wlasciwosci strzalki")

	render.DrawTextAt(dst, "Kolor:", pp.X+10, pp.Y+30, color.RGBA{80, 110, 180, 200})
	vector.FillRect(dst, float32(pp.X+60), float32(pp.Y+24), 28, 16, e.Color, true)
	vector.StrokeRect(dst, float32(pp.X+60), float32(pp.Y+24), 28, 16, 1, color.RGBA{200, 200, 255, 160}, true)
	// First 8 colour swatches only
	pp.drawSwatches(dst, pp.colSwatches[:8], e.Color)

	render.DrawTextAt(dst, "Styl:", pp.X+10, pp.Y+128, color.RGBA{80, 110, 180, 200})
	pp.drawStyleBtns(dst, e.Style)

	render.DrawTextAt(dst, "Etykieta: \""+e.Label+"\"", pp.X+10, pp.Y+170, color.RGBA{80, 120, 180, 200})
	render.DrawTextAt(dst, "PrawyKlik=edytuj  Del=usun", pp.X+10, pp.Y+185, color.RGBA{55, 80, 130, 180})
	render.DrawTextAt(dst, "Aby usunac etykiete: edytuj->Enter", pp.X+10, pp.Y+198, color.RGBA{55, 80, 130, 160})
}

func (pp *PropPanel) drawBg(dst *ebiten.Image, h float64) {
	render.FillRoundRect(dst, float32(pp.X-6), float32(pp.Y-6), float32(pp.W+12), float32(h+12), 10,
		color.RGBA{7, 10, 28, 242})
	render.StrokeRoundRect(dst, float32(pp.X-6), float32(pp.Y-6), float32(pp.W+12), float32(h+12), 10, 1,
		color.RGBA{35, 60, 140, 200})
}

// drawTitleBar draws the draggable title strip at the top of the panel.
func (pp *PropPanel) drawTitleBar(dst *ebiten.Image, title string) {
	// Highlighted bar — signals it's draggable
	vector.FillRect(dst, float32(pp.X-6), float32(pp.Y-6), float32(pp.W+12), 24,
		color.RGBA{20, 40, 100, 200}, false)
	render.DrawTextAt(dst, ":: "+title, pp.X+6, pp.Y+2, color.RGBA{120, 170, 255, 255})
	// Drag hint dots on right side
	for i := 0; i < 3; i++ {
		vector.FillCircle(dst,
			float32(pp.X+pp.W-10), float32(pp.Y+4+float64(i)*5), 1.5,
			color.RGBA{80, 120, 200, 180}, true)
	}
}

// TitleBarHitTest returns true if (mx,my) is on the title bar.
func (pp *PropPanel) TitleBarHitTest(mx, my float64) bool {
	return pp.Visible && inRect(mx, my, pp.X-6, pp.Y-6, pp.W+12, 24)
}

// StartDrag begins dragging the panel from screen point (mx,my).
func (pp *PropPanel) StartDrag(mx, my float64) {
	pp.dragging = true
	pp.dragOffX = mx - pp.X
	pp.dragOffY = my - pp.Y
}

// UpdateDrag moves the panel to follow (mx,my).
func (pp *PropPanel) UpdateDrag(mx, my float64) {
	if !pp.dragging {
		return
	}
	pp.X = mx - pp.dragOffX
	pp.Y = my - pp.dragOffY
	pp.rebuild()
}

// StopDrag ends dragging.
func (pp *PropPanel) StopDrag() {
	pp.dragging = false
}

func (pp *PropPanel) drawSwatches(dst *ebiten.Image, swatches []swatch, active color.RGBA) {
	for _, s := range swatches {
		isActive := s.col == active
		if isActive {
			render.FillRoundRect(dst, float32(s.x)-3, float32(s.y)-3, float32(s.w)+6, float32(s.h)+6, 5,
				color.RGBA{s.col.R, s.col.G, s.col.B, 90})
		}
		render.FillRoundRect(dst, float32(s.x), float32(s.y), float32(s.w), float32(s.h), 4, s.col)
		if isActive {
			render.StrokeRoundRect(dst, float32(s.x), float32(s.y), float32(s.w), float32(s.h), 4, 2,
				color.RGBA{255, 255, 255, 230})
		}
	}
}

func (pp *PropPanel) drawAnimBtns(dst *ebiten.Image, current model.Anim) {
	for _, ab := range pp.animBtns {
		active := ab.anim == current
		bg, border, tc := smallBtnColors(active)
		render.FillRoundRect(dst, float32(ab.x), float32(ab.y), float32(ab.w), float32(ab.h), 5, bg)
		render.StrokeRoundRect(dst, float32(ab.x), float32(ab.y), float32(ab.w), float32(ab.h), 5, 1, border)
		render.DrawTextCentered(dst, ab.label, ab.x+ab.w/2, ab.y+ab.h/2, tc)
	}
}

func (pp *PropPanel) drawBorderBtns(dst *ebiten.Image, current float32) {
	for _, bb := range pp.borderBtns {
		active := math.Abs(float64(bb.thick-current)) < 0.5
		bg, border, tc := smallBtnColors(active)
		render.FillRoundRect(dst, float32(bb.x), float32(bb.y), float32(bb.w), float32(bb.h), 5, bg)
		render.StrokeRoundRect(dst, float32(bb.x), float32(bb.y), float32(bb.w), float32(bb.h), 5, 1, border)
		render.DrawTextCentered(dst, bb.label, bb.x+bb.w/2, bb.y+bb.h/2, tc)
	}
}

func (pp *PropPanel) drawStyleBtns(dst *ebiten.Image, current model.EdgeStyle) {
	for _, sb := range pp.styleBtns {
		active := sb.style == current
		bg, border, tc := smallBtnColors(active)
		render.FillRoundRect(dst, float32(sb.x), float32(sb.y), float32(sb.w), float32(sb.h), 5, bg)
		render.StrokeRoundRect(dst, float32(sb.x), float32(sb.y), float32(sb.w), float32(sb.h), 5, 1, border)
		render.DrawTextCentered(dst, sb.label, sb.x+sb.w/2, sb.y+sb.h/2, tc)
	}
}

func smallBtnColors(active bool) (bg, border, tc color.RGBA) {
	if active {
		return color.RGBA{25, 55, 140, 255}, color.RGBA{80, 155, 255, 255}, color.RGBA{255, 255, 255, 255}
	}
	return color.RGBA{12, 18, 50, 255}, color.RGBA{35, 60, 130, 180}, color.RGBA{140, 170, 230, 200}
}

func (pp *PropPanel) ClickNodeColour(mx, my float64) (color.RGBA, bool) {
	for _, s := range pp.colSwatches {
		if inRect(mx, my, s.x, s.y, s.w, s.h) {
			return s.col, true
		}
	}
	return color.RGBA{}, false
}

func (pp *PropPanel) ClickBgColour(mx, my float64) (color.RGBA, bool) {
	for _, s := range pp.bgSwatches {
		if inRect(mx, my, s.x, s.y, s.w, s.h) {
			return s.col, true
		}
	}
	return color.RGBA{}, false
}

func (pp *PropPanel) ClickEdgeColour(mx, my float64) (color.RGBA, bool) {
	for _, s := range pp.colSwatches[:8] {
		if inRect(mx, my, s.x, s.y, s.w, s.h) {
			return s.col, true
		}
	}
	return color.RGBA{}, false
}

func (pp *PropPanel) ClickAnim(mx, my float64) (model.Anim, bool) {
	for _, ab := range pp.animBtns {
		if inRect(mx, my, ab.x, ab.y, ab.w, ab.h) {
			return ab.anim, true
		}
	}
	return 0, false
}

func (pp *PropPanel) ClickBorder(mx, my float64) (float32, bool) {
	for _, bb := range pp.borderBtns {
		if inRect(mx, my, bb.x, bb.y, bb.w, bb.h) {
			return bb.thick, true
		}
	}
	return 0, false
}

func (pp *PropPanel) ClickEdgeStyle(mx, my float64) (model.EdgeStyle, bool) {
	for _, sb := range pp.styleBtns {
		if inRect(mx, my, sb.x, sb.y, sb.w, sb.h) {
			return sb.style, true
		}
	}
	return 0, false
}

func (pp *PropPanel) HitTest(mx, my float64) bool {
	return pp.Visible && inRect(mx, my, pp.X-6, pp.Y-6, pp.W+12, pp.H+12)
}

// ─── Save dialog ─────────────────────────────────────────────────────────────

type SaveDialog struct {
	Visible   bool
	X, Y      float64
	W, H      float64
	Filename  string
	IsExport  bool
	IsLoad    bool // load mode
	ExportSVG bool
	Active    bool
}

func NewSaveDialog(screenW, screenH float64) *SaveDialog {
	w, h := 400.0, 170.0
	return &SaveDialog{
		X: screenW/2 - w/2, Y: screenH/2 - h/2,
		W: w, H: h,
		Filename: "flowchart",
	}
}

func (sd *SaveDialog) Reposition(screenW, screenH float64) {
	sd.X = screenW/2 - sd.W/2
	sd.Y = screenH/2 - sd.H/2
}

func (sd *SaveDialog) Draw(dst *ebiten.Image) {
	if !sd.Visible {
		return
	}
	// Dim background
	vector.FillRect(dst, 0, 0, float32(dst.Bounds().Dx()), float32(dst.Bounds().Dy()),
		color.RGBA{0, 0, 0, 150}, false)

	// Panel
	render.FillRoundRect(dst, float32(sd.X), float32(sd.Y), float32(sd.W), float32(sd.H), 10,
		color.RGBA{9, 13, 34, 252})
	render.StrokeRoundRect(dst, float32(sd.X), float32(sd.Y), float32(sd.W), float32(sd.H), 10, 1.5,
		color.RGBA{60, 120, 255, 230})

	title := "Zapisz diagram"
	if sd.IsLoad {
		title = "Wczytaj diagram"
	} else if sd.IsExport {
		title = "Eksportuj diagram"
	}
	render.DrawTextCentered(dst, title, sd.X+sd.W/2, sd.Y+18, color.RGBA{130, 190, 255, 255})

	// Filename label
	render.DrawTextAt(dst, "Nazwa pliku:", sd.X+16, sd.Y+42, color.RGBA{80, 110, 180, 200})

	// Input field — highlighted border
	render.FillRoundRect(dst, float32(sd.X+16), float32(sd.Y+56), float32(sd.W-32), 28, 5,
		color.RGBA{12, 20, 55, 255})
	render.StrokeRoundRect(dst, float32(sd.X+16), float32(sd.Y+56), float32(sd.W-32), 28, 5, 2,
		color.RGBA{100, 180, 255, 255})

	// Filename text + blinking cursor
	display := sd.Filename + "|"
	render.DrawTextAt(dst, display, sd.X+24, sd.Y+63, color.RGBA{220, 235, 255, 255})

	// Extension hint
	ext := ".json"
	if sd.IsExport {
		if sd.ExportSVG {
			ext = ".svg"
		} else {
			ext = ".png"
		}
	}
	if sd.IsLoad {
		ext = ".json"
	}
	render.DrawTextAt(dst, ext, sd.X+sd.W-52, sd.Y+63, color.RGBA{100, 130, 200, 180})

	if sd.IsExport {
		render.DrawTextAt(dst, "Format:", sd.X+16, sd.Y+98, color.RGBA{80, 110, 180, 200})
		// PNG button
		pngActive := !sd.ExportSVG
		pbg, pborder, ptc := smallBtnColors(pngActive)
		render.FillRoundRect(dst, float32(sd.X+70), float32(sd.Y+92), 60, 22, 5, pbg)
		render.StrokeRoundRect(dst, float32(sd.X+70), float32(sd.Y+92), 60, 22, 5, 1, pborder)
		render.DrawTextCentered(dst, "PNG", sd.X+100, sd.Y+103, ptc)
		// SVG button
		sbg, sborder, stc := smallBtnColors(sd.ExportSVG)
		render.FillRoundRect(dst, float32(sd.X+140), float32(sd.Y+92), 60, 22, 5, sbg)
		render.StrokeRoundRect(dst, float32(sd.X+140), float32(sd.Y+92), 60, 22, 5, 1, sborder)
		render.DrawTextCentered(dst, "SVG", sd.X+170, sd.Y+103, stc)
	}

	// OK button
	btnY := sd.Y + sd.H - 38
	render.FillRoundRect(dst, float32(sd.X+16), float32(btnY), 90, 28, 6, color.RGBA{20, 80, 200, 255})
	render.StrokeRoundRect(dst, float32(sd.X+16), float32(btnY), 90, 28, 6, 1.5, color.RGBA{80, 160, 255, 230})
	okLabel := "OK  [Enter]"
	if sd.IsLoad {
		okLabel = "Wczytaj [Enter]"
	}
	render.DrawTextCentered(dst, okLabel, sd.X+61, btnY+14, color.RGBA{255, 255, 255, 255})

	// Cancel button
	render.FillRoundRect(dst, float32(sd.X+sd.W-106), float32(btnY), 90, 28, 6, color.RGBA{60, 15, 15, 255})
	render.StrokeRoundRect(dst, float32(sd.X+sd.W-106), float32(btnY), 90, 28, 6, 1.5, color.RGBA{180, 50, 50, 220})
	render.DrawTextCentered(dst, "Anuluj  [Esc]", sd.X+sd.W-61, btnY+14, color.RGBA{255, 180, 180, 255})
}

func (sd *SaveDialog) ClickOK(mx, my float64) bool {
	btnY := sd.Y + sd.H - 38
	return inRect(mx, my, sd.X+16, btnY, 90, 28)
}

func (sd *SaveDialog) ClickCancel(mx, my float64) bool {
	btnY := sd.Y + sd.H - 38
	return inRect(mx, my, sd.X+sd.W-106, btnY, 90, 28)
}

func (sd *SaveDialog) ClickPNG(mx, my float64) bool {
	return sd.IsExport && inRect(mx, my, sd.X+70, sd.Y+92, 60, 22)
}

func (sd *SaveDialog) ClickSVG(mx, my float64) bool {
	return sd.IsExport && inRect(mx, my, sd.X+140, sd.Y+92, 60, 22)
}

func (sd *SaveDialog) HitTest(mx, my float64) bool {
	return sd.Visible && inRect(mx, my, sd.X, sd.Y, sd.W, sd.H)
}

// ─── Status bar ──────────────────────────────────────────────────────────────

func DrawStatusBar(dst *ebiten.Image, tool Tool, nodeCount, edgeCount int, msg string) {
	sw := float64(dst.Bounds().Dx())
	sh := float64(dst.Bounds().Dy())
	y := sh - StatusbarH
	vector.FillRect(dst, 0, float32(y), float32(sw), StatusbarH, color.RGBA{7, 9, 25, 235}, false)
	vector.StrokeRect(dst, 0, float32(y), float32(sw), 1, 1, color.RGBA{25, 45, 100, 200}, false)

	toolNames := map[Tool]string{
		ToolSelect: "Wybierz(V)", ToolNode: "Wezel(N)", ToolConnect: "Polacz(C)",
		ToolPan: "Przesun(H)", ToolDelete: "Usun(X)", ToolClear: "Wyczysc", ToolLoad: "Wczytaj",
	}
	status := fmt.Sprintf("  %s  |  Wezly:%d  Krawedzie:%d  |  F5=zapis  F6=export  T=tytul  L=legenda  G=siatka  S+G=snap",
		toolNames[tool], nodeCount, edgeCount)
	if msg != "" {
		status = "  " + msg
	}
	render.DrawTextAt(dst, status, SidebarW+8, y+5, color.RGBA{90, 120, 190, 215})
}

// ─── Text editor overlays ─────────────────────────────────────────────────────
func DrawTextEditor(dst *ebiten.Image, n *model.Node, field int, buf string, camX, camY float64) {
	if n == nil {
		return
	}
	fieldLabel := "Etykieta"
	if field == 1 {
		fieldLabel = "Podtytul"
	}
	text := fieldLabel + ": " + buf + "|"
	bw := float64(len([]rune(text))*7) + 24
	bh := 26.0
	wx := n.X + camX + SidebarW
	wy := n.Y + camY + n.H + 8
	render.FillRoundRect(dst, float32(wx), float32(wy), float32(bw), float32(bh), 5, color.RGBA{8, 12, 38, 245})
	render.StrokeRoundRect(dst, float32(wx), float32(wy), float32(bw), float32(bh), 5, 1.5, color.RGBA{80, 155, 255, 220})
	render.DrawTextAt(dst, text, wx+8, wy+7, color.RGBA{200, 220, 255, 255})
}

func DrawEdgeLabelEditor(dst *ebiten.Image, mx, my float64, buf string) {
	text := "Etykieta: " + buf + "|"
	bw := float64(len([]rune(text))*7) + 24
	bh := 26.0
	x := mx - bw/2
	y := my - bh/2 - 24
	render.FillRoundRect(dst, float32(x), float32(y), float32(bw), float32(bh), 5, color.RGBA{8, 12, 38, 245})
	render.StrokeRoundRect(dst, float32(x), float32(y), float32(bw), float32(bh), 5, 1.5, color.RGBA{255, 180, 50, 220})
	render.DrawTextAt(dst, text, x+8, y+7, color.RGBA{255, 220, 150, 255})
}

// ─── Helper ───────────────────────────────────────────────────────────────────

func inRect(px, py, x, y, w, h float64) bool {
	return px >= x && px <= x+w && py >= y && py <= y+h
}

func clampC(v int) uint8 {
	if v > 255 {
		return 255
	}
	if v < 0 {
		return 0
	}
	return uint8(v)
}
