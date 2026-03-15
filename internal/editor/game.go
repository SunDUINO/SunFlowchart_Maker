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
 * ║  Plik / File: game.go                                        ║
 * ║                                                              ║
 * ║  Licencja / License: MIT                                     ║
 * ║  Rok / Year: 2025-2026                                       ║
 * ╚══════════════════════════════════════════════════════════════╝
 */

package editor

import (
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/user/flowchart/internal/export"
	"github.com/user/flowchart/internal/model"
	"github.com/user/flowchart/internal/render"
	"github.com/user/flowchart/internal/ui"
)

type Game struct {
	diagram *model.Diagram
	history History
	toolbar *ui.Toolbar
	props   *ui.PropPanel
	dialog  *ui.SaveDialog

	camX, camY float64

	selected []*model.Node
	selEdge  *model.Edge

	dragNode *model.Node
	dragOffX float64
	dragOffY float64

	resizeNode   *model.Node
	resizeStartW float64
	resizeStartH float64
	resizeStartX float64
	resizeStartY float64

	banding        bool
	bandX1, bandY1 float64
	bandX2, bandY2 float64

	panning              bool
	panStartX, panStartY float64
	panCamX, panCamY     float64

	connectFrom *model.Node
	hoverNode   *model.Node
	hoverEdge   *model.Edge
	propDragging bool // dragging the properties panel

	editLegend bool   // editing legend text
	legendBuf  string // temp buffer while editing

	showGrid bool // whether to draw grid dots
	snapGrid bool // whether to snap nodes to grid

	editNode  *model.Node
	editEdge  *model.Edge
	editField int
	editBuf   string

	time      float64
	statusMsg string
	statusTTL float64
	screenW   int
	screenH   int
	initDone  bool
}

func NewGame() *Game {
	g := &Game{
		diagram: model.NewDiagram(),
		toolbar: ui.NewToolbar(),
		screenW: 1400,
		screenH: 860,
	}
	model.LoadSample(g.diagram)
	g.history.Push(g.diagram)
	g.showGrid = true
	g.snapGrid = true
	g.toolbar.ShowGrid = true
	g.toolbar.SnapGrid = true
	// Init panels at default size — they reposition on first Draw()
	g.props = ui.NewPropPanel(1400)
	g.dialog = ui.NewSaveDialog(1400, 860)
	return g
}

func (g *Game) ensureInit() {
	if g.initDone {
		return
	}
	if g.screenW == 0 {
		return
	}
	// Reposition panels to actual screen size
	g.props = ui.NewPropPanel(float64(g.screenW))
	g.dialog = ui.NewSaveDialog(float64(g.screenW), float64(g.screenH))
	g.initDone = true
}

func (g *Game) Update() error {
	const dt = 1.0 / 60.0
	g.time += dt
	g.diagram.Update(dt)
	if g.statusTTL > 0 {
		g.statusTTL -= dt
	}
	// Track hovered node for connect tool highlighting
	if g.toolbar.ActiveTool == ui.ToolConnect {
		mx, my := ebiten.CursorPosition()
		wx := float64(mx) - g.camX - ui.SidebarW
		wy := float64(my) - g.camY
		g.hoverNode = g.diagram.NodeAtExpanded(wx, wy, 12)
	} else {
		g.hoverNode = nil
	}
	// Track hovered edge for select tool
	if g.toolbar.ActiveTool == ui.ToolSelect || g.toolbar.ActiveTool == ui.ToolDelete {
		mx, my := ebiten.CursorPosition()
		wx := float64(mx) - g.camX - ui.SidebarW
		wy := float64(my) - g.camY
		if g.diagram.NodeAt(wx, wy) == nil {
			g.hoverEdge = g.diagram.EdgeNear(wx, wy, 28)
		} else {
			g.hoverEdge = nil
		}
	} else {
		g.hoverEdge = nil
	}
	g.handleKeyboard()
	g.handleMouse()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.screenW = screen.Bounds().Dx()
	g.screenH = screen.Bounds().Dy()
	g.ensureInit()
	// Reposition dialog if screen was resized
	if g.dialog != nil && !g.dialog.Visible {
		g.dialog.Reposition(float64(g.screenW), float64(g.screenH))
	}

	tw := ui.SidebarW
	ox := g.camX + tw
	oy := g.camY

	render.DrawGrid(screen, ox, oy, g.diagram.BgColor, g.showGrid)

	// Edges
	offsets := g.diagram.EdgeOffsets()
	for _, e := range g.diagram.Edges {
		from := g.diagram.NodeByID(e.FromID)
		to := g.diagram.NodeByID(e.ToID)
		if from == nil || to == nil {
			continue
		}
		isHover := g.hoverEdge != nil && g.hoverEdge.ID == e.ID && !e.Selected
		render.DrawEdge(screen, e, offsetNode(from, ox, oy), offsetNode(to, ox, oy), offsets[e.ID], isHover)
	}

	// Connect preview + highlights
	if g.connectFrom != nil {
		cx, cy := ebiten.CursorPosition()
		render.DrawEdgePreview(screen, offsetNode(g.connectFrom, ox, oy), float64(cx), float64(cy))
		// Pulse ring on source node
		render.DrawConnectSource(screen, offsetNode(g.connectFrom, ox, oy), g.time)
	}
	// Hover highlight in Connect mode
	if g.hoverNode != nil && g.hoverNode != g.connectFrom {
		render.DrawConnectHover(screen, offsetNode(g.hoverNode, ox, oy))
	}

	// Nodes
	sorted := make([]*model.Node, len(g.diagram.Nodes))
	copy(sorted, g.diagram.Nodes)
	sort.Slice(sorted, func(i, j int) bool {
		return !sorted[i].Selected && sorted[j].Selected
	})
	for _, n := range sorted {
		tn := offsetNode(n, ox, oy)
		render.DrawNode(screen, tn)
		cx, cy := tn.Centre()
		render.DrawNodeLabel(screen, int(tn.Shape), tn.Label, tn.Sub, cx, cy, tn.TextColor)
		if n.Selected && len(g.selected) == 1 {
			render.DrawResizeHandle(screen, tn)
		}
	}

	if g.banding {
		render.DrawSelectionRect(screen, g.bandX1+ox, g.bandY1+oy, g.bandX2+ox, g.bandY2+oy)
	}

	g.toolbar.Draw(screen)

	// Props panel
	if g.props != nil {
		if len(g.selected) == 1 {
			g.props.Visible = true
			g.props.Draw(screen, g.selected[0], g.time)
		} else if g.selEdge != nil {
			g.props.Visible = true
			g.props.DrawEdgeProps(screen, g.selEdge)
		} else {
			g.props.Visible = false
		}
	}

	msg := ""
	if g.statusTTL > 0 {
		msg = g.statusMsg
	}
	ui.DrawStatusBar(screen, g.toolbar.ActiveTool, len(g.diagram.Nodes), len(g.diagram.Edges), msg)

	if g.editNode != nil {
		ui.DrawTextEditor(screen, g.editNode, g.editField, g.editBuf, g.camX, g.camY)
	}
	if g.editEdge != nil {
		from := g.diagram.NodeByID(g.editEdge.FromID)
		to := g.diagram.NodeByID(g.editEdge.ToID)
		if from != nil && to != nil {
			mx := (from.X+from.W/2+to.X+to.W/2)/2 + ox
			my := (from.Y+from.H/2+to.Y+to.H/2)/2 + oy
			ui.DrawEdgeLabelEditor(screen, mx, my, g.editBuf)
		}
	}

	if g.dialog != nil {
		g.dialog.Draw(screen)
	}

	// Legend panel — pinned bottom-left
	legendText := g.diagram.Legend
	if g.editLegend {
		legendText = g.legendBuf
	}
	render.DrawLegend(screen, legendText, ui.SidebarW, ui.StatusbarH, g.editLegend)
}

func (g *Game) Layout(ow, oh int) (int, int) { return ow, oh }

// ─── Keyboard ────────────────────────────────────────────────────────────────

func (g *Game) handleKeyboard() {
	// Dialog input takes full priority
	if g.dialog != nil && g.dialog.Visible {
		g.handleDialogKeys()
		return
	}
	// Text editing takes full priority
	if g.editLegend {
		g.handleLegendInput()
		return
	}
	if g.editNode != nil || g.editEdge != nil {
		g.handleTextInput()
		return
	}

	ctrlHeld := ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta)

	// Undo / Redo
	if ctrlHeld && inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			if g.history.Redo(g.diagram) {
				g.flash("Redo")
				g.clearSelection()
			}
		} else {
			if g.history.Undo(g.diagram) {
				g.flash("Undo")
				g.clearSelection()
			}
		}
		return
	}
	if ctrlHeld && inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.openDialog(false)
		return
	}
	if ctrlHeld && inpututil.IsKeyJustPressed(ebiten.KeyO) {
		if err := loadDiagram(g.diagram); err != nil {
			g.flash("Blad wczytania: " + err.Error())
		} else {
			g.clearSelection()
			g.flash("Wczytano flowchart.json")
		}
		return
	}
	if ctrlHeld && inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.openDialog(true)
		return
	}
	if ctrlHeld && inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.clearSelection()
		for _, n := range g.diagram.Nodes {
			n.Selected = true
			g.selected = append(g.selected, n)
		}
		return
	}
	// Block other shortcuts when Ctrl held
	if ctrlHeld {
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeyV) {
		g.toolbar.ActiveTool = ui.ToolSelect
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) || inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.toolbar.ActiveTool = ui.ToolNode
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) || inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.toolbar.ActiveTool = ui.ToolTitle
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) || inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.toolbar.ActiveTool = ui.ToolConnect
	}
	if inpututil.IsKeyJustPressed(ebiten.Key5) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.toolbar.ActiveTool = ui.ToolPan
	}
	if inpututil.IsKeyJustPressed(ebiten.Key6) {
		g.toolbar.ActiveTool = ui.ToolDelete
	}
	if inpututil.IsKeyJustPressed(ebiten.Key7) {
		g.toolbar.ActiveTool = ui.ToolClear
	}
	// L = edit legend
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.legendBuf = g.diagram.Legend
		g.editLegend = true
	}
	// G = toggle grid visibility (check Shift first for snap)
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			g.snapGrid = !g.snapGrid
			g.toolbar.SnapGrid = g.snapGrid
			if g.snapGrid {
				g.flash("Snap: wlaczony")
			} else {
				g.flash("Snap: wylaczony")
			}
		} else {
			g.showGrid = !g.showGrid
			g.toolbar.ShowGrid = g.showGrid
			if g.showGrid {
				g.flash("Siatka: widoczna")
			} else {
				g.flash("Siatka: ukryta")
			}
		}
	}
	// F5 = save, F6 = export (alternative to Ctrl+S/E in case of focus issues)
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		g.openDialog(false)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF6) {
		g.openDialog(true)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		if len(g.selected) == 1 {
			g.startEditNode(g.selected[0], 0)
		} else if g.selEdge != nil {
			g.startEditEdge(g.selEdge)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) || inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		g.deleteSelected()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.clearSelection()
		if g.connectFrom != nil {
			g.connectFrom = nil
			g.flash("Połączenie anulowane")
		}
	}

	nudge := 4.0
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		nudge = 1.0
	}
	dx, dy := 0.0, 0.0
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		dx = -nudge
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		dx = nudge
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		dy = -nudge
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		dy = nudge
	}
	if (dx != 0 || dy != 0) && len(g.selected) > 0 {
		g.snapshot()
		for _, n := range g.selected {
			n.X += dx
			n.Y += dy
		}
	}
}

func (g *Game) handleTextInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.commitEdit()
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) && g.editNode != nil {
		field := g.editField
		n := g.editNode
		g.commitEdit()
		if len(g.selected) == 1 && g.selected[0] == n {
			g.startEditNode(n, 1-field)
		}
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.editNode = nil
		g.editEdge = nil
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		runes := []rune(g.editBuf)
		if len(runes) > 0 {
			g.editBuf = string(runes[:len(runes)-1])
		}
		return
	}
	for _, r := range ebiten.AppendInputChars(nil) {
		g.editBuf += string(r)
	}
}

func (g *Game) handleDialogKeys() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.dialog.Visible = false
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.executeDialog()
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		runes := []rune(g.dialog.Filename)
		if len(runes) > 0 {
			g.dialog.Filename = string(runes[:len(runes)-1])
		}
		return
	}
	for _, r := range ebiten.AppendInputChars(nil) {
		g.dialog.Filename += string(r)
	}
}

// ─── Mouse ───────────────────────────────────────────────────────────────────

func (g *Game) handleMouse() {
	mx, my := ebiten.CursorPosition()
	fmx, fmy := float64(mx), float64(my)
	wx := fmx - g.camX - ui.SidebarW
	wy := fmy - g.camY

	_, sy := ebiten.Wheel()
	if sy != 0 {
		g.camY += sy * 28
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.onLeftDown(fmx, fmy, wx, wy)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.onLeftUp(wx, wy)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.onLeftHeld(fmx, fmy, wx, wy)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.onRightDown(wx, wy)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		g.panning = true
		g.panStartX, g.panStartY = fmx, fmy
		g.panCamX, g.panCamY = g.camX, g.camY
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
		g.panning = false
	}
	if g.panning && ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		g.camX = g.panCamX + fmx - g.panStartX
		g.camY = g.panCamY + fmy - g.panStartY
	}
}

func (g *Game) onLeftDown(fmx, fmy, wx, wy float64) {
	// Commit any active text edit when clicking elsewhere
	if g.editNode != nil || g.editEdge != nil {
		g.commitEdit()
	}
	if g.dialog != nil && g.dialog.Visible {
		if g.dialog.ClickOK(fmx, fmy) {
			g.executeDialog()
		} else if g.dialog.ClickCancel(fmx, fmy) {
			g.dialog.Visible = false
		} else if g.dialog.ClickPNG(fmx, fmy) {
			g.dialog.ExportSVG = false
		} else if g.dialog.ClickSVG(fmx, fmy) {
			g.dialog.ExportSVG = true
		}
		return
	}

	// Grid/Snap buttons first — they're in the toolbar zone but handled separately
	gridT, snapT := g.toolbar.ClickGridBtn(fmx, fmy, float64(g.screenH))
	if gridT {
		g.showGrid = !g.showGrid
		g.toolbar.ShowGrid = g.showGrid
		if g.showGrid {
			g.flash("Siatka: widoczna  [G]")
		} else {
			g.flash("Siatka: ukryta  [G]")
		}
		return
	}
	if snapT {
		g.snapGrid = !g.snapGrid
		g.toolbar.SnapGrid = g.snapGrid
		if g.snapGrid {
			g.flash("Snap: wlaczony  [Shift+G]")
		} else {
			g.flash("Snap: wylaczony  [Shift+G]")
		}
		return
	}

	if g.toolbar.Click(fmx, fmy) {
		return
	}

	// Legend panel click — bottom left
	if g.legendHitTest(fmx, fmy) {
		g.legendBuf = g.diagram.Legend
		g.editLegend = true
		return
	}

	// Props panel title bar drag
	if g.props != nil && g.props.Visible && g.props.TitleBarHitTest(fmx, fmy) {
		g.props.StartDrag(fmx, fmy)
		g.propDragging = true
		return
	}

	// Props panel clicks
	if g.props != nil && g.props.HitTest(fmx, fmy) {
		if len(g.selected) == 1 {
			n := g.selected[0]
			if col, ok := g.props.ClickNodeColour(fmx, fmy); ok {
				g.snapshot()
				n.Color = col
				n.FillColor = model.DarkFill(col)
				return
			}
			if col, ok := g.props.ClickBgColour(fmx, fmy); ok {
				g.snapshot()
				g.diagram.BgColor = col
				return
			}
			if anim, ok := g.props.ClickAnim(fmx, fmy); ok {
				g.snapshot()
				n.Anim = anim
				return
			}
			if thick, ok := g.props.ClickBorder(fmx, fmy); ok {
				g.snapshot()
				n.BorderW = thick
				return
			}
		} else if g.selEdge != nil {
			e := g.selEdge
			if col, ok := g.props.ClickEdgeColour(fmx, fmy); ok {
				g.snapshot()
				e.Color = col
				return
			}
			if style, ok := g.props.ClickEdgeStyle(fmx, fmy); ok {
				g.snapshot()
				e.Style = style
				return
			}
		}
		return
	}

	switch g.toolbar.ActiveTool {
	case ui.ToolSelect:
		if len(g.selected) == 1 {
			n := g.selected[0]
			hx := n.X + n.W - 5
			hy := n.Y + n.H - 5
			if wx >= hx && wx <= hx+10 && wy >= hy && wy <= hy+10 {
				g.resizeNode = n
				g.resizeStartW = n.W
				g.resizeStartH = n.H
				g.resizeStartX = wx
				g.resizeStartY = wy
				return
			}
		}
		g.selectClick(wx, wy)
	case ui.ToolNode:
		g.snapshot()
		g.addNode(wx, wy)
	case ui.ToolTitle:
		g.snapshot()
		g.addTitleNode(wx, wy)
	case ui.ToolConnect:
		n := g.diagram.NodeAtExpanded(wx, wy, 12)
		if n != nil {
			if g.connectFrom == nil {
				g.connectFrom = n
				g.flash("Kliknij węzeł docelowy...")
			} else if g.connectFrom != n {
				g.snapshot()
				g.diagram.AddEdge(&model.Edge{
					FromID: g.connectFrom.ID,
					ToID:   n.ID,
					Color:  color.RGBA{80, 155, 255, 190},
				})
				g.connectFrom = nil
				g.flash("")
			}
		} else {
			g.connectFrom = nil
		}
	case ui.ToolPan:
		g.panning = true
		g.panStartX, g.panStartY = fmx, fmy
		g.panCamX, g.panCamY = g.camX, g.camY
	case ui.ToolDelete:
		n := g.diagram.NodeAt(wx, wy)
		if n != nil {
			g.snapshot()
			g.diagram.RemoveNode(n.ID)
			g.clearSelection()
		} else if e := g.diagram.EdgeNear(wx, wy, 28); e != nil {
			g.snapshot()
			g.diagram.RemoveEdge(e.ID)
		}
	case ui.ToolClear:
		g.snapshot()
		g.diagram.Nodes = nil
		g.diagram.Edges = nil
		g.clearSelection()
		g.flash("Ekran wyczyszczony  (Ctrl+Z cofa)")
	case ui.ToolLoad:
		g.openLoadDialog()
	}
}

func (g *Game) onLeftUp(_, _ float64) {
	g.panning = false
	g.propDragging = false
	if g.props != nil {
		g.props.StopDrag()
	}
	if g.resizeNode != nil {
		g.snapshot()
		g.resizeNode = nil
	}
	if g.dragNode != nil {
		g.snapshot()
		g.dragNode = nil
	}
	if g.banding {
		g.banding = false
		g.applyBand()
	}
}

func (g *Game) onLeftHeld(fmx, fmy, wx, wy float64) {
	if g.dialog != nil && g.dialog.Visible {
		return
	}
	// Props panel drag
	if g.propDragging && g.props != nil {
		g.props.UpdateDrag(fmx, fmy)
		return
	}
	if g.panning {
		g.camX = g.panCamX + fmx - g.panStartX
		g.camY = g.panCamY + fmy - g.panStartY
		return
	}
	if g.resizeNode != nil {
		g.resizeNode.W = math.Max(40, g.resizeStartW+(wx-g.resizeStartX))
		g.resizeNode.H = math.Max(24, g.resizeStartH+(wy-g.resizeStartY))
		if g.snapGrid {
			g.resizeNode.W = snapTo(g.resizeNode.W)
			g.resizeNode.H = snapTo(g.resizeNode.H)
		}
		return
	}
	if g.dragNode != nil {
		if len(g.selected) > 1 {
			dx := wx - g.dragNode.X - g.dragOffX
			dy := wy - g.dragNode.Y - g.dragOffY
			for _, n := range g.selected {
				n.X += dx
				n.Y += dy
			}
		} else {
			newX := wx - g.dragOffX
			newY := wy - g.dragOffY
			if g.snapGrid {
				newX = snapTo(newX)
				newY = snapTo(newY)
			}
			g.dragNode.X = newX
			g.dragNode.Y = newY
		}
		return
	}
	if g.banding {
		g.bandX2, g.bandY2 = wx, wy
	}
}

func (g *Game) onRightDown(wx, wy float64) {
	// In Node/Title tool — right click deletes node under cursor
	if g.toolbar.ActiveTool == ui.ToolNode || g.toolbar.ActiveTool == ui.ToolTitle {
		n := g.diagram.NodeAt(wx, wy)
		if n != nil {
			g.snapshot()
			g.diagram.RemoveNode(n.ID)
			g.clearSelection()
			g.flash("Wezel usuniety")
			return
		}
	}

	// Cancel connect mode on right click anywhere
	if g.connectFrom != nil {
		g.connectFrom = nil
		g.flash("Anulowano")
		return
	}

	n := g.diagram.NodeAt(wx, wy)
	if n != nil {
		if n.Selected && len(g.selected) == 1 && g.selected[0] == n {
			g.startEditNode(n, 0)
		} else {
			g.clearSelection()
			n.Selected = true
			g.selected = []*model.Node{n}
		}
		return
	}
	e := g.diagram.EdgeNear(wx, wy, 28)
	if e != nil {
		g.clearSelection()
		e.Selected = true
		g.selEdge = e
		g.startEditEdge(e)
		return
	}
	// Right click on empty canvas = clear selection
	g.clearSelection()
}

// ─── Selection ────────────────────────────────────────────────────────────────

func (g *Game) selectClick(wx, wy float64) {
	n := g.diagram.NodeAt(wx, wy)
	if n != nil {
		shift := ebiten.IsKeyPressed(ebiten.KeyShift)
		if !shift && !n.Selected {
			g.clearSelection()
		}
		if !g.inSel(n) {
			n.Selected = true
			g.selected = append(g.selected, n)
		}
		g.dragNode = n
		g.dragOffX = wx - n.X
		g.dragOffY = wy - n.Y
	} else {
		e := g.diagram.EdgeNear(wx, wy, 28)
		if e != nil {
			g.clearSelection()
			e.Selected = true
			g.selEdge = e
		} else {
			g.clearSelection()
			g.banding = true
			g.bandX1, g.bandY1 = wx, wy
			g.bandX2, g.bandY2 = wx, wy
		}
	}
}

func (g *Game) clearSelection() {
	for _, n := range g.selected {
		n.Selected = false
	}
	g.selected = nil
	if g.selEdge != nil {
		g.selEdge.Selected = false
		g.selEdge = nil
	}
}

func (g *Game) inSel(n *model.Node) bool {
	for _, s := range g.selected {
		if s == n {
			return true
		}
	}
	return false
}

func (g *Game) applyBand() {
	x1, y1 := g.bandX1, g.bandY1
	x2, y2 := g.bandX2, g.bandY2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if math.Abs(x2-x1) < 5 || math.Abs(y2-y1) < 5 {
		return
	}
	g.clearSelection()
	for _, n := range g.diagram.Nodes {
		if n.X+n.W > x1 && n.X < x2 && n.Y+n.H > y1 && n.Y < y2 {
			n.Selected = true
			g.selected = append(g.selected, n)
		}
	}
}

func (g *Game) deleteSelected() {
	if len(g.selected) == 0 && g.selEdge == nil {
		return
	}
	g.snapshot()
	for _, n := range g.selected {
		g.diagram.RemoveNode(n.ID)
	}
	if g.selEdge != nil {
		g.diagram.RemoveEdge(g.selEdge.ID)
	}
	g.clearSelection()
}

// ─── Node / edge creation ─────────────────────────────────────────────────────

func (g *Game) handleLegendInput() {
	// Shift+Enter = newline
	if ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.legendBuf += "\n"
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.snapshot()
		g.diagram.Legend = g.legendBuf
		g.editLegend = false
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.editLegend = false
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		runes := []rune(g.legendBuf)
		if len(runes) > 0 {
			g.legendBuf = string(runes[:len(runes)-1])
		}
		return
	}
	for _, r := range ebiten.AppendInputChars(nil) {
		g.legendBuf += string(r)
	}
}

func (g *Game) legendHitTest(mx, my float64) bool {
	if g.screenH == 0 {
		return false
	}
	// Match DrawLegend position: x=SidebarW+8, bottom above statusbar
	x := ui.SidebarW + 8
	sh := float64(g.screenH)
	maxW := 320.0
	// rough height estimate — enough to catch a click
	h := 80.0
	if g.diagram.Legend != "" {
		lines := len([]rune(g.diagram.Legend)) / 40 + 2
		h = float64(lines)*16 + 24
	}
	y := sh - ui.StatusbarH - h - 8
	return mx >= x && mx <= x+maxW && my >= y && my <= y+h
}

func (g *Game) addTitleNode(wx, wy float64) {
	col := color.RGBA{230, 190, 50, 255}
	n := &model.Node{
		X: wx - 160, Y: wy - 20, W: 320, H: 40,
		Label:     "Tytuł diagramu",
		Shape:     model.ShapeTitle,
		Color:     col,
		FillColor: model.DarkFill(col),
		TextColor: color.RGBA{255, 230, 100, 255},
		Anim:      model.AnimNone,
		AnimSpeed: 0.8,
	}
	g.diagram.AddNode(n)
	g.clearSelection()
	n.Selected = true
	g.selected = []*model.Node{n}
	g.startEditNode(n, 0)
}

func (g *Game) addNode(wx, wy float64) *model.Node {
	shape := g.toolbar.ActiveShape
	w, h := defaultSizeFor(shape)
	col := color.RGBA{40, 120, 220, 255}
	nx, ny := wx-w/2, wy-h/2
	if g.snapGrid {
		nx = snapTo(nx)
		ny = snapTo(ny)
	}
	n := &model.Node{
		X: nx, Y: ny, W: w, H: h,
		Label: "Node", Shape: shape,
		Color: col, FillColor: model.DarkFill(col),
		TextColor: color.RGBA{220, 230, 255, 255},
		Anim: model.AnimNone, AnimSpeed: 0.8,
	}
	g.diagram.AddNode(n)
	g.clearSelection()
	n.Selected = true
	g.selected = []*model.Node{n}
	g.startEditNode(n, 0)
	return n
}

func defaultSizeFor(s model.Shape) (float64, float64) {
	switch s {
	case model.ShapeDiamond:
		return 120, 68
	case model.ShapeOval:
		return 110, 44
	default:
		return 130, 52
	}
}

// ─── Text editing ─────────────────────────────────────────────────────────────

func (g *Game) startEditNode(n *model.Node, field int) {
	g.editNode = n
	g.editEdge = nil
	g.editField = field
	if field == 0 {
		g.editBuf = n.Label
	} else {
		g.editBuf = n.Sub
	}
}

func (g *Game) startEditEdge(e *model.Edge) {
	g.editEdge = e
	g.editNode = nil
	g.editBuf = e.Label
}

func (g *Game) commitEdit() {
	if g.editNode != nil {
		g.snapshot()
		if g.editField == 0 {
			g.editNode.Label = g.editBuf
		} else {
			g.editNode.Sub = g.editBuf
		}
		g.editNode = nil
	} else if g.editEdge != nil {
		g.snapshot()
		g.editEdge.Label = g.editBuf
		g.editEdge = nil
	}
}

// ─── Save dialog ──────────────────────────────────────────────────────────────

func (g *Game) openLoadDialog() {
	if g.dialog == nil {
		return
	}
	g.dialog.IsLoad = true
	g.dialog.IsExport = false
	g.dialog.Filename = "flowchart"
	g.dialog.Visible = true
}

func (g *Game) openDialog(isExport bool) {
	if g.dialog == nil {
		return
	}
	g.dialog.IsLoad = false
	g.dialog.IsExport = isExport
	g.dialog.ExportSVG = false
	g.dialog.Filename = "flowchart"
	g.dialog.Visible = true
}

func (g *Game) executeDialog() {
	if g.dialog == nil {
		return
	}
	g.dialog.Visible = false
	name := g.dialog.Filename
	if name == "" {
		name = "flowchart"
	}

	// Load mode
	if g.dialog.IsLoad {
		if err := loadDiagramFrom(g.diagram, name+".json"); err != nil {
			g.flash("Blad wczytania: " + err.Error())
		} else {
			g.clearSelection()
			g.history.Push(g.diagram)
			g.flash("Wczytano: " + name + ".json")
		}
		return
	}
	if g.dialog.IsExport {
		if g.dialog.ExportSVG {
			if err := export.ExportSVG(g.diagram, name+".svg"); err != nil {
				g.flash("SVG błąd: " + err.Error())
			} else {
				g.flash("Zapisano: " + name + ".svg")
			}
		} else {
			if filename, err := export.ExportPNG(g.diagram, name+".png"); err != nil {
				g.flash("PNG błąd: " + err.Error())
			} else {
				g.flash("Zapisano: " + filename)
			}
		}
	} else {
		if err := saveDiagram(g.diagram, name+".json"); err != nil {
			g.flash("Błąd zapisu: " + err.Error())
		} else {
			g.flash("Zapisano: " + name + ".json")
		}
	}
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func (g *Game) snapshot() { g.history.Push(g.diagram) }

func (g *Game) flash(msg string) {
	g.statusMsg = msg
	g.statusTTL = 4.0
}

func offsetNode(n *model.Node, ox, oy float64) *model.Node {
	cp := *n
	cp.X += ox
	cp.Y += oy
	return &cp
}

// snapTo rounds v to nearest grid cell.
func snapTo(v float64) float64 {
	s := render.GridSpacing
	return math.Round(v/s) * s
}

func saveDiagram(d *model.Diagram, filename string) error {
	data, err := marshalDiagram(d)
	if err != nil {
		return err
	}
	return writeFile(filename, data)
}

func loadDiagram(d *model.Diagram) error {
	data, err := readFile("flowchart.json")
	if err != nil {
		return err
	}
	return unmarshalDiagram(d, data)
}