# SunFlowChart Maker

> Neon-dark flowchart editor written in Go + Ebiten

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![Ebiten](https://img.shields.io/badge/Ebiten-v2.9.9-orange?style=flat)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)
![Version](https://img.shields.io/badge/Version-1.3.0-blue?style=flat)

**Author:** Andrzej "Sunriver" Gromczyński / Lothar TeaM

**Forum:** [forum.lothar-team.pl](https://forum.lothar-team.pl/)

**GitHub:** [github.com/SunDUINO](https://github.com/SunDUINO)



---

## Features

- **Drag & drop** node placement and movement
- **5 node shapes** — Rectangle, Rounded, Diamond, Parallelogram, Oval
- **Special Title node** — large draggable heading with underline
- **6 animations** — Glow, Pulse, Blink, Flash, Spinner, Off
- **Smart edge routing** — bezier curves or orthogonal elbows, parallel edges auto-separated
- **Neon glow effect** on all nodes
- **Color palette** — 12 neon node colors + 8 background colors
- **Border thickness** — 1px / 2px / 3px / 4px
- **Edge labels** — editable on right-click
- **Legend panel** — pinned to bottom-left, multi-line, exported to SVG
- **Snap to grid** — 24px grid with toggle
- **Undo / Redo** — full history
- **Save / Load** — JSON format
- **Export PNG** — renders to file
- **Export SVG** — 2× scaled, with bezier curves and legend

---

## Requirements

- Go 1.22+
- On Linux: `sudo apt install libgl1-mesa-dev xorg-dev`

---

## Build & Run

```bash
# Clone
git clone https://github.com/SunDUINO/SunFlowchart_Maker
cd SunFlowchart_Maker

# Download dependencies
go mod tidy

# Run
go run ./src

# Build binary (Windows)
go build -o sunflowchart.exe ./src

# Build binary (Linux)
GOOS=linux GOARCH=amd64 go build -o sunflowchart ./src
```

---

## Project Structure

```
SunFlowchart_Maker/
├── src/
│   └── main.go              # Entry point
├── internal/
│   ├── model/
│   │   ├── model.go         # Node, Edge, Diagram data model
│   │   ├── sample.go        # Sample startup diagram
│   │   └── helpers.go       # Model helper methods
│   ├── render/
│   │   ├── draw.go          # Node rendering, glow, grid
│   │   ├── edges.go         # Arrow routing and rendering
│   │   └── text.go          # Text, labels, legend panel
│   ├── editor/
│   │   ├── game.go          # Main game loop, input handling
│   │   ├── history.go       # Undo/redo snapshot system
│   │   └── serial.go        # JSON save & load
│   ├── export/
│   │   └── png.go           # PNG and SVG export
│   └── ui/
│       ├── ui.go            # Toolbar, properties panel, dialogs
│       └── icons.go         # Geometric toolbar icons
├── go.mod
└── README.md
```

---

## Controls

### Tools

| Key | Tool | Description |
|-----|------|-------------|
| `V` / `1` | Select | Select and drag nodes |
| `N` / `2` | Node | Click canvas to add a node |
| `T` / `3` | Title | Add a large draggable title node |
| `C` / `4` | Connect | Click source node, then target node |
| `H` / `5` | Pan | Drag to scroll canvas |
| `X` / `6` | Delete | Click node or edge to remove |
| `7` | Clear | Clear the entire canvas |
| `Ctrl+O` | Load | Open load dialog |

### Editing

| Key | Action |
|-----|--------|
| `F2` | Edit selected node label |
| `Tab` | Switch label ↔ sub-label during edit |
| `Enter` | Confirm text |
| `Esc` | Cancel / deselect |
| `Del` | Delete selected |
| `Arrow keys` | Nudge node 4px (`+Shift` = 1px) |
| `Ctrl+A` | Select all |
| `Shift+Click` | Multi-select |
| Right-click drag | Rubber-band selection |

### Mouse

| Action | Effect |
|--------|--------|
| Left click | Select node |
| Right click (canvas) | Deselect all |
| Right click (node 1×) | Select + properties panel |
| Right click (node 2×) | Edit label |
| Right click (node tool) | Delete node under cursor |
| Right click (edge) | Edit edge label |
| Yellow corner handle | Resize node |
| Middle mouse drag | Free pan |
| Scroll wheel | Vertical pan |

### File

| Key | Action |
|-----|--------|
| `Ctrl+S` / `F5` | Save as JSON |
| `Ctrl+E` / `F6` | Export PNG or SVG |
| `Ctrl+O` | Load JSON |
| `Ctrl+Z` | Undo |
| `Ctrl+Shift+Z` | Redo |

### View

| Key | Action |
|-----|--------|
| `G` | Toggle grid visibility |
| `Shift+G` | Toggle snap to grid |
| `L` | Edit legend |

---

## Node Animations

| Name | Effect |
|------|--------|
| **Off** | No animation, static dim glow |
| **Glow** | Constant bright neon glow |
| **Pulse** | Slow breathing effect |
| **Blink** | Hard on/off flicker |
| **Flash** | Quick bright flash, then fades |
| **Spin** | Rotating arc around the node |

---

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/hajimehoshi/ebiten/v2` | v2.9.9 | Game engine / rendering |
| `github.com/hajimehoshi/bitmapfont/v3` | v3.3.0 | Built-in bitmap font |
| `golang.org/x/image` | v0.37.0 | Image utilities |

---

## License

MIT License — see [LICENSE](LICENSE) for details.

```
Copyright (c) 2026 Andrzej "Sunriver" Gromczyński / Lothar TeaM
```


