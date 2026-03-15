package editor

import (
	"encoding/json"
	"image/color"
	"os"

	"github.com/user/flowchart/internal/model"
)

type jsonNode struct {
	ID, Shape, Anim     int
	X, Y, W, H         float64
	AnimSpeed, BorderW  float64
	Label, Sub          string
	CR, CG, CB, CA      uint8
	FR, FG, FB, FA      uint8
	TR, TG, TB, TA      uint8
}

type jsonEdge struct {
	ID, FromID, ToID int
	Label            string
	Style            int
	CR, CG, CB, CA   uint8
}

type jsonDiagram struct {
	Nodes               []jsonNode
	Edges               []jsonEdge
	NextID              int
	BgR, BgG, BgB, BgA uint8
	Legend              string
}

func marshalDiagram(d *model.Diagram) ([]byte, error) {
	jd := jsonDiagram{
		BgR: d.BgColor.R, BgG: d.BgColor.G, BgB: d.BgColor.B, BgA: d.BgColor.A,
		Legend: d.Legend,
	}
	for _, n := range d.Nodes {
		jd.Nodes = append(jd.Nodes, jsonNode{
			ID: n.ID, Shape: int(n.Shape), Anim: int(n.Anim),
			X: n.X, Y: n.Y, W: n.W, H: n.H,
			AnimSpeed: n.AnimSpeed, BorderW: float64(n.BorderW),
			Label: n.Label, Sub: n.Sub,
			CR: n.Color.R, CG: n.Color.G, CB: n.Color.B, CA: n.Color.A,
			FR: n.FillColor.R, FG: n.FillColor.G, FB: n.FillColor.B, FA: n.FillColor.A,
			TR: n.TextColor.R, TG: n.TextColor.G, TB: n.TextColor.B, TA: n.TextColor.A,
		})
	}
	for _, e := range d.Edges {
		jd.Edges = append(jd.Edges, jsonEdge{
			ID: e.ID, FromID: e.FromID, ToID: e.ToID,
			Label: e.Label, Style: int(e.Style),
			CR: e.Color.R, CG: e.Color.G, CB: e.Color.B, CA: e.Color.A,
		})
	}
	return json.MarshalIndent(jd, "", "  ")
}

func unmarshalDiagram(d *model.Diagram, data []byte) error {
	var jd jsonDiagram
	if err := json.Unmarshal(data, &jd); err != nil {
		return err
	}
	d.Nodes = d.Nodes[:0]
	d.Edges = d.Edges[:0]
	d.BgColor = color.RGBA{jd.BgR, jd.BgG, jd.BgB, jd.BgA}
	if d.BgColor.A == 0 {
		d.BgColor = color.RGBA{7, 9, 22, 255}
	}
	d.Legend = jd.Legend
	for _, jn := range jd.Nodes {
		n := &model.Node{
			ID: jn.ID, Shape: model.Shape(jn.Shape), Anim: model.Anim(jn.Anim),
			X: jn.X, Y: jn.Y, W: jn.W, H: jn.H,
			AnimSpeed: jn.AnimSpeed, BorderW: float32(jn.BorderW),
			Label: jn.Label, Sub: jn.Sub,
			Color:     color.RGBA{jn.CR, jn.CG, jn.CB, jn.CA},
			FillColor: color.RGBA{jn.FR, jn.FG, jn.FB, jn.FA},
			TextColor: color.RGBA{jn.TR, jn.TG, jn.TB, jn.TA},
		}
		d.Nodes = append(d.Nodes, n)
	}
	for _, je := range jd.Edges {
		e := &model.Edge{
			ID: je.ID, FromID: je.FromID, ToID: je.ToID,
			Label: je.Label, Style: model.EdgeStyle(je.Style),
			Color: color.RGBA{je.CR, je.CG, je.CB, je.CA},
		}
		d.Edges = append(d.Edges, e)
	}
	d.SetNextID(jd.NextID)
	return nil
}

func writeFile(path string, data []byte) error { return os.WriteFile(path, data, 0644) }
func readFile(path string) ([]byte, error)      { return os.ReadFile(path) }

func loadDiagramFrom(d *model.Diagram, filename string) error {
	data, err := readFile(filename)
	if err != nil {
		return err
	}
	return unmarshalDiagram(d, data)
}