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
 * ║  Plik / File: model.go                                       ║
 * ║                                                              ║
 * ║  Licencja / License: MIT                                     ║
 * ║  Rok / Year: 2025-2026                                       ║
 * ╚══════════════════════════════════════════════════════════════╝
 */

package model

import (
	"image/color"
	"math"
)

// ─── Enums ────────────────────────────────────────────────────────────────────

type Shape int

const (
	ShapeRect Shape = iota
	ShapeRounded
	ShapeDiamond
	ShapeParallel
	ShapeOval
	ShapeTitle // special: large text, no border, draggable
)

func (s Shape) String() string {
	return [...]string{"Rect", "Rounded", "Diamond", "Parallel", "Oval", "Title"}[s]
}

type Anim int

const (
	AnimNone Anim = iota
	AnimGlow
	AnimPulse
	AnimBlink
	AnimFlash
	AnimSpinner
)

func (a Anim) String() string {
	return [...]string{"None", "Glow", "Pulse", "Blink", "Flash", "Spinner"}[a]
}

// EdgeStyle controls arrow routing.
type EdgeStyle int

const (
	EdgeCurve EdgeStyle = iota // smooth bezier
	EdgeElbow                  // orthogonal L-shape
)

// ─── Node ────────────────────────────────────────────────────────────────────

type Node struct {
	ID    int
	X, Y  float64
	W, H  float64
	Label string
	Sub   string

	Shape     Shape
	Color     color.RGBA
	FillColor color.RGBA
	TextColor color.RGBA
	BorderW   float32 // border thickness, default 1.8

	Anim      Anim
	AnimSpeed float64
	animPhase float64
	Selected  bool
}

func (n *Node) Update(dt float64) {
	if n.AnimSpeed == 0 {
		n.AnimSpeed = 0.8
	}
	n.animPhase = math.Mod(n.animPhase+n.AnimSpeed*dt, 1.0)
}

func (n *Node) GlowAlpha() uint8 {
	p := n.animPhase
	switch n.Anim {
	case AnimNone:
		return 50
	case AnimGlow:
		return 220
	case AnimPulse:
		v := (math.Sin(p*2*math.Pi) + 1) / 2
		return uint8(v * 240)
	case AnimBlink:
		if math.Mod(p*2, 1.0) < 0.5 {
			return 240
		}
		return 0
	case AnimFlash:
		v := math.Exp(-math.Mod(p, 1.0) * 3.5)
		return uint8(v * 240)
	case AnimSpinner:
		return 180
	}
	return 50
}

func (n *Node) SpinnerAngle() float64 { return n.animPhase * 2 * math.Pi }

func (n *Node) Contains(px, py float64) bool {
	return px >= n.X && px <= n.X+n.W && py >= n.Y && py <= n.Y+n.H
}

func (n *Node) Centre() (float64, float64) {
	return n.X + n.W/2, n.Y + n.H/2
}

func (n *Node) EdgePoint(tx, ty float64) (float64, float64) {
	cx, cy := n.Centre()
	dx, dy := tx-cx, ty-cy
	if dx == 0 && dy == 0 {
		return cx, cy
	}
	hw, hh := n.W/2, n.H/2
	if math.Abs(dx)*hh > math.Abs(dy)*hw {
		if dx > 0 {
			return cx + hw, cy + dy*(hw/dx)
		}
		return cx - hw, cy + dy*(-hw/dx)
	}
	if dy > 0 {
		return cx + dx*(hh/dy), cy + hh
	}
	return cx + dx*(-hh/dy), cy - hh
}

func (n *Node) BorderThickness() float32 {
	if n.BorderW <= 0 {
		return 1.8
	}
	return n.BorderW
}

// ─── Edge ────────────────────────────────────────────────────────────────────

type Edge struct {
	ID       int
	FromID   int
	ToID     int
	Label    string
	Color    color.RGBA
	Style    EdgeStyle
	Selected bool
}

// ─── Diagram ─────────────────────────────────────────────────────────────────

type Diagram struct {
	Nodes   []*Node
	Edges   []*Edge
	BgColor color.RGBA
	Legend  string // multi-line legend text, shown pinned to bottom-left
	nextID  int
}

func NewDiagram() *Diagram {
	return &Diagram{BgColor: color.RGBA{7, 9, 22, 255}}
}

func (d *Diagram) AddNode(n *Node) {
	d.nextID++
	n.ID = d.nextID
	d.Nodes = append(d.Nodes, n)
}

func (d *Diagram) AddEdge(e *Edge) {
	d.nextID++
	e.ID = d.nextID
	d.Edges = append(d.Edges, e)
}

func (d *Diagram) RemoveNode(id int) {
	ns := d.Nodes[:0]
	for _, n := range d.Nodes {
		if n.ID != id {
			ns = append(ns, n)
		}
	}
	d.Nodes = ns
	es := d.Edges[:0]
	for _, e := range d.Edges {
		if e.FromID != id && e.ToID != id {
			es = append(es, e)
		}
	}
	d.Edges = es
}

func (d *Diagram) RemoveEdge(id int) {
	es := d.Edges[:0]
	for _, e := range d.Edges {
		if e.ID != id {
			es = append(es, e)
		}
	}
	d.Edges = es
}

func (d *Diagram) NodeByID(id int) *Node {
	for _, n := range d.Nodes {
		if n.ID == id {
			return n
		}
	}
	return nil
}

func (d *Diagram) Update(dt float64) {
	for _, n := range d.Nodes {
		n.Update(dt)
	}
}

func (d *Diagram) NodeAt(wx, wy float64) *Node {
	for i := len(d.Nodes) - 1; i >= 0; i-- {
		if d.Nodes[i].Contains(wx, wy) {
			return d.Nodes[i]
		}
	}
	return nil
}

func (d *Diagram) EdgeNear(wx, wy, r float64) *Edge {
	offsets := d.EdgeOffsets()
	for _, e := range d.Edges {
		from := d.NodeByID(e.FromID)
		to := d.NodeByID(e.ToID)
		if from == nil || to == nil {
			continue
		}
		if edgePointNear(e, from, to, offsets[e.ID], wx, wy, r) {
			return e
		}
	}
	return nil
}

// edgePointNear returns true if (wx,wy) is within r pixels of the edge curve.
func edgePointNear(_ *Edge, from, to *Node, offset int, wx, wy, r float64) bool {
	tcx, tcy := to.Centre()
	fcx, fcy := from.Centre()
	sx, sy := from.EdgePoint(tcx, tcy)
	ex, ey := to.EdgePoint(fcx, fcy)

	lateralPx := float64(offset) * 14.0
	dx := ex - sx
	dy := ey - sy
	dist := math.Sqrt(dx*dx + dy*dy)
	var px, py float64
	if dist > 0 {
		px, py = -dy/dist, dx/dist
	}

	// Sample points along the bezier curve
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

	const steps = 20
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		u := 1 - t
		// Cubic bezier point
		bx := u*u*u*sx + 3*u*u*t*c1x + 3*u*t*t*c2x + t*t*t*ex
		by := u*u*u*sy + 3*u*u*t*c1y + 3*u*t*t*c2y + t*t*t*ey
		ddx, ddy := wx-bx, wy-by
		if ddx*ddx+ddy*ddy <= r*r {
			return true
		}
	}
	return false
}

// ─── Colour palette ───────────────────────────────────────────────────────────

var Palette = []color.RGBA{
	{63, 84, 186, 255},
	{220, 50, 50, 255},
	{50, 200, 80, 255},
	{230, 140, 40, 255},
	{80, 180, 255, 255},
	{40, 120, 220, 255},
	{220, 220, 240, 255},
	{0, 210, 210, 255},
	{160, 50, 210, 255},
	{210, 210, 50, 255},
	{220, 80, 140, 255},
	{50, 180, 150, 255},
}

var BgPalette = []color.RGBA{
	{7, 9, 22, 255},
	{12, 18, 40, 255},
	{20, 12, 30, 255},
	{8, 22, 12, 255},
	{30, 20, 10, 255},
	{20, 20, 20, 255},
	{240, 240, 240, 255},
	{255, 255, 255, 255},
}

func DarkFill(c color.RGBA) color.RGBA {
	return color.RGBA{c.R / 7, c.G / 7, c.B / 7, 215}
}

// NodeAtExpanded returns the topmost node within expanded hit area (for easier clicking).
func (d *Diagram) NodeAtExpanded(wx, wy, margin float64) *Node {
	for i := len(d.Nodes) - 1; i >= 0; i-- {
		n := d.Nodes[i]
		if wx >= n.X-margin && wx <= n.X+n.W+margin &&
			wy >= n.Y-margin && wy <= n.Y+n.H+margin {
			return n
		}
	}
	return nil
}

// EdgeOffsets computes the lateral offset index for each edge,
// so parallel edges between the same pair of nodes are spread apart.
// Returns a map from edge ID → offset (-1, 0, +1, -2, +2, ...).
func (d *Diagram) EdgeOffsets() map[int]int {
	// Count edges per unordered pair
	type pair struct{ a, b int }
	makePair := func(x, y int) pair {
		if x < y {
			return pair{x, y}
		}
		return pair{y, x}
	}

	// Group edge IDs by pair
	groups := map[pair][]int{}
	for _, e := range d.Edges {
		p := makePair(e.FromID, e.ToID)
		groups[p] = append(groups[p], e.ID)
	}

	result := make(map[int]int, len(d.Edges))
	for _, ids := range groups {
		n := len(ids)
		for i, id := range ids {
			// Spread: 0, +1, -1, +2, -2, ...
			var off int
			switch {
			case n == 1:
				off = 0
			case i%2 == 0:
				off = i / 2
			default:
				off = -(i/2 + 1)
			}
			result[id] = off
		}
	}
	return result
}
