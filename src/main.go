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
 * ║  Plik / File: main.go                                        ║
 * ║                                                              ║
 * ║  Licencja / License: MIT                                     ║
 * ║  Rok / Year: 2025-2026                                       ║
 * ╚══════════════════════════════════════════════════════════════╝
 */
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/user/flowchart/internal/editor"
)

var version ="1.3.0";

func main() {
	ebiten.SetWindowSize(1400, 860)
	ebiten.SetWindowTitle("SunFlowChart Editor — Neon Dark" + " v." + version)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)

	g := editor.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
