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
 * ║  Plik / File: text.go                                        ║
 * ║                                                              ║
 * ║  Licencja / License: MIT                                     ║
 * ║  Rok / Year: 2025-2026                                       ║
 * ╚══════════════════════════════════════════════════════════════╝
 */
 
package render



import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	ebitentext "github.com/hajimehoshi/ebiten/v2/text/v2"
)

var bitmapFace *ebitentext.GoXFace

func init() {
	bitmapFace = ebitentext.NewGoXFace(bitmapfont.Face)
}

const lineHeight = 14.0

// DrawTextCentered draws multi-line text centred at (cx, cy).
func DrawTextCentered(dst *ebiten.Image, text string, cx, cy float64, col color.RGBA) {
	lines := strings.Split(text, "\n")
	totalH := float64(len(lines)) * lineHeight
	startY := cy - totalH/2

	for i, line := range lines {
		if line == "" {
			continue
		}
		w, _ := ebitentext.Measure(line, bitmapFace, lineHeight)
		op := &ebitentext.DrawOptions{}
		op.GeoM.Translate(cx-w/2, startY+float64(i)*lineHeight)
		op.ColorScale.ScaleWithColor(col)
		ebitentext.Draw(dst, line, bitmapFace, op)
	}
}

// DrawTextAt draws text at (x, y) top-left.
func DrawTextAt(dst *ebiten.Image, text string, x, y float64, col color.RGBA) {
	if text == "" {
		return
	}
	op := &ebitentext.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(col)
	ebitentext.Draw(dst, text, bitmapFace, op)
}

// DrawNodeLabel draws label + optional sub centred inside a node.
func DrawNodeLabel(dst *ebiten.Image, shape int, label, sub string, cx, cy float64, textCol color.RGBA) {
	if shape == int(5) { // ShapeTitle = 5
		DrawTitleLabel(dst, label, cx, cy, textCol)
		return
	}
	if sub == "" {
		DrawTextCentered(dst, label, cx, cy-lineHeight/2, textCol)
		return
	}
	DrawTextCentered(dst, label, cx, cy-lineHeight, textCol)
	subCol := color.RGBA{textCol.R, textCol.G, textCol.B, 170}
	DrawTextCentered(dst, sub, cx, cy+4, subCol)
}

// DrawEdgeLabel draws a small labelled box on an edge midpoint.
func DrawEdgeLabel(dst *ebiten.Image, text string, mx, my float64) {
	if text == "" {
		return
	}
	w, _ := ebitentext.Measure(text, bitmapFace, lineHeight)
	bw := w + 10
	bh := lineHeight + 4
	drawFilledRect(dst, mx-bw/2, my-bh/2, bw, bh, color.RGBA{10, 14, 35, 210})
	drawStrokeRect(dst, mx-bw/2, my-bh/2, bw, bh, 1, color.RGBA{60, 100, 180, 180})
	DrawTextCentered(dst, text, mx, my-lineHeight/2+2, color.RGBA{180, 200, 255, 230})
}

// ─── Internal rect helpers ────────────────────────────────────────────────────

func drawFilledRect(dst *ebiten.Image, x, y, w, h float64, col color.RGBA) {
	if w < 1 || h < 1 {
		return
	}
	img := ebiten.NewImage(int(w)+1, int(h)+1)
	img.Fill(col)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	dst.DrawImage(img, op)
}

func drawStrokeRect(dst *ebiten.Image, x, y, w, h, lw float64, col color.RGBA) {
	drawFilledRect(dst, x, y, w, lw, col)
	drawFilledRect(dst, x, y+h-lw, w, lw, col)
	drawFilledRect(dst, x, y, lw, h, col)
	drawFilledRect(dst, x+w-lw, y, lw, h, col)
}

// DrawTitleLabel draws a large bold-style title centred in a title node.
func DrawTitleLabel(dst *ebiten.Image, label string, cx, cy float64, col color.RGBA) {
	if label == "" {
		return
	}
	// Draw slightly offset to simulate bold
	for _, off := range [][2]float64{{-0.5, 0}, {0.5, 0}, {0, -0.5}} {
		op := &ebitentext.DrawOptions{}
		w, _ := ebitentext.Measure(label, bitmapFace, lineHeight)
		op.GeoM.Translate(cx-w/2+off[0], cy-lineHeight/2+off[1])
		op.ColorScale.ScaleWithColor(color.RGBA{col.R, col.G, col.B, col.A / 3})
		ebitentext.Draw(dst, label, bitmapFace, op)
	}
	// Main text
	op := &ebitentext.DrawOptions{}
	w, _ := ebitentext.Measure(label, bitmapFace, lineHeight)
	op.GeoM.Translate(cx-w/2, cy-lineHeight/2)
	op.ColorScale.ScaleWithColor(col)
	ebitentext.Draw(dst, label, bitmapFace, op)
}

// DrawLegend draws the pinned legend panel at bottom-left of screen.
func DrawLegend(dst *ebiten.Image, text string, sidebarW, statusbarH float64, editing bool) {
	if text == "" && !editing {
		return
	}

	sh := float64(dst.Bounds().Dy())
	x := sidebarW + 8
	maxW := 320.0
	lineH := lineHeight + 2

	// Count lines
	lines := splitLines(text)
	visLines := lines
	if len(visLines) == 0 {
		visLines = []string{""}
	}
	h := float64(len(visLines))*lineH + 20

	y := sh - statusbarH - h - 8

	// Background
	bgCol := color.RGBA{8, 12, 35, 220}
	borderCol := color.RGBA{50, 80, 160, 200}
	if editing {
		borderCol = color.RGBA{255, 180, 50, 230}
	}
	drawFilledRect(dst, x, y, maxW, h, bgCol)
	drawStrokeRect(dst, x, y, maxW, h, 1.5, borderCol)

	// Title bar
	drawFilledRect(dst, x, y, maxW, 16, color.RGBA{20, 40, 100, 200})
	titleCol := color.RGBA{120, 160, 255, 230}
	hint := "  Legenda  [L]=edytuj"
	if editing {
		hint = "  Legenda  [Enter]=zatwierdz  [Esc]=anuluj"
		titleCol = color.RGBA{255, 200, 80, 255}
	}
	DrawTextAt(dst, hint, x+4, y+2, titleCol)

	// Content
	for i, line := range visLines {
		ly := y + 18 + float64(i)*lineH
		DrawTextAt(dst, line, x+6, ly, color.RGBA{180, 210, 255, 220})
	}

	// Cursor at end if editing
	if editing {
		lastLine := ""
		if len(visLines) > 0 {
			lastLine = visLines[len(visLines)-1]
		}
		w2, _ := ebitentext.Measure(lastLine, bitmapFace, lineHeight)
		ly := y + 18 + float64(len(visLines)-1)*lineH
		drawFilledRect(dst, x+6+w2, ly, 2, lineH-2, color.RGBA{255, 200, 80, 200})
	}
}

func splitLines(text string) []string {
	if text == "" {
		return nil
	}
	var lines []string
	cur := ""
	for _, r := range text {
		if r == '\n' {
			lines = append(lines, cur)
			cur = ""
		} else {
			cur += string(r)
		}
	}
	lines = append(lines, cur)
	return lines
}