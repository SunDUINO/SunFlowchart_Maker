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
 * ║  Plik / File: history.go                                     ║
 * ║                                                              ║
 * ║  Licencja / License: MIT                                     ║
 * ║  Rok / Year: 2025-2026                                       ║
 * ╚══════════════════════════════════════════════════════════════╝
 */

package editor

import (
	"encoding/json"
	"image/color"

	"github.com/user/flowchart/internal/model"
)

// ─── Serialisable snapshot ───────────────────────────────────────────────────

type snapNode struct {
	ID, Shape, Anim             int
	X, Y, W, H, AnimSpeed       float64
	BorderW                     float32
	Label, Sub                  string
	Color, FillColor, TextColor [4]uint8
}

type snapEdge struct {
	ID, FromID, ToID int
	Label            string
	Color            [4]uint8
	Style            int
}

type snapshot struct {
	Nodes              []snapNode
	Edges              []snapEdge
	Legend             string
	BgR, BgG, BgB, BgA uint8
}

func toSnap(d *model.Diagram) snapshot {
	s := snapshot{
		Legend: d.Legend,
		BgR:    d.BgColor.R, BgG: d.BgColor.G,
		BgB: d.BgColor.B, BgA: d.BgColor.A,
	}
	for _, n := range d.Nodes {
		s.Nodes = append(s.Nodes, snapNode{
			ID: n.ID, Shape: int(n.Shape), Anim: int(n.Anim),
			X: n.X, Y: n.Y, W: n.W, H: n.H, AnimSpeed: n.AnimSpeed,
			BorderW: n.BorderW,
			Label:   n.Label, Sub: n.Sub,
			Color:     rgba(n.Color),
			FillColor: rgba(n.FillColor),
			TextColor: rgba(n.TextColor),
		})
	}
	for _, e := range d.Edges {
		s.Edges = append(s.Edges, snapEdge{
			ID: e.ID, FromID: e.FromID, ToID: e.ToID,
			Label: e.Label, Color: rgba(e.Color), Style: int(e.Style),
		})
	}
	return s
}

func applySnap(d *model.Diagram, s snapshot) {
	d.Nodes = d.Nodes[:0]
	d.Edges = d.Edges[:0]
	d.Legend = s.Legend
	if s.BgA > 0 {
		d.BgColor = color.RGBA{s.BgR, s.BgG, s.BgB, s.BgA}
	}
	maxID := 0
	for _, sn := range s.Nodes {
		n := &model.Node{
			ID: sn.ID, Shape: model.Shape(sn.Shape), Anim: model.Anim(sn.Anim),
			X: sn.X, Y: sn.Y, W: sn.W, H: sn.H, AnimSpeed: sn.AnimSpeed,
			BorderW: sn.BorderW,
			Label:   sn.Label, Sub: sn.Sub,
			Color:     fromRGBA(sn.Color),
			FillColor: fromRGBA(sn.FillColor),
			TextColor: fromRGBA(sn.TextColor),
		}
		d.Nodes = append(d.Nodes, n)
		if n.ID > maxID {
			maxID = n.ID
		}
	}
	for _, se := range s.Edges {
		e := &model.Edge{
			ID: se.ID, FromID: se.FromID, ToID: se.ToID,
			Label: se.Label, Color: fromRGBA(se.Color), Style: model.EdgeStyle(se.Style),
		}
		d.Edges = append(d.Edges, e)
		if e.ID > maxID {
			maxID = e.ID
		}
	}
	// patch nextID via JSON roundtrip trick: we expose it via a helper below
	setNextID(d, maxID+1)
}

func rgba(c color.RGBA) [4]uint8     { return [4]uint8{c.R, c.G, c.B, c.A} }
func fromRGBA(a [4]uint8) color.RGBA { return color.RGBA{a[0], a[1], a[2], a[3]} }

// setNextID patches the unexported nextID field by marshalling/unmarshalling
// the diagram through a helper struct.  Simpler than reflection.
func setNextID(d *model.Diagram, id int) {
	// We add a public helper in model to avoid exporting nextID.
	d.SetNextID(id)
}

// ─── History ─────────────────────────────────────────────────────────────────

const maxHistory = 60

type History struct {
	stack [][]byte // JSON-encoded snapshots
	pos   int      // current position (0 = oldest)
}

func (h *History) Push(d *model.Diagram) {
	data, err := json.Marshal(toSnap(d))
	if err != nil {
		return
	}
	// Truncate redo stack
	h.stack = append(h.stack[:h.pos], data)
	h.pos++
	if len(h.stack) > maxHistory {
		h.stack = h.stack[1:]
		h.pos = len(h.stack)
	}
}

func (h *History) Undo(d *model.Diagram) bool {
	if h.pos <= 1 {
		return false
	}
	h.pos--
	var s snapshot
	if err := json.Unmarshal(h.stack[h.pos-1], &s); err != nil {
		return false
	}
	applySnap(d, s)
	return true
}

func (h *History) Redo(d *model.Diagram) bool {
	if h.pos >= len(h.stack) {
		return false
	}
	var s snapshot
	if err := json.Unmarshal(h.stack[h.pos], &s); err != nil {
		return false
	}
	h.pos++
	applySnap(d, s)
	return true
}
