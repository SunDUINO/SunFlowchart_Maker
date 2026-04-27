/*
 * ╔════════════════════════════════════════════════════════════════╗
 * ║ SunFlowchart_Maker  v1.3.0                                     ║
 * ║ Plik / File: main.go                                           ║
 * ╠════════════════════════════════════════════════════════════════╣
 * ║ Autor / Author:                                                ║
 * ║   SunRiver                                                     ║
 * ║   Lothar TeaM                                                  ║
 * ╠════════════════════════════════════════════════════════════════╣
 * ║ GitHub  : github.com/user/flowchart                            ║
 * ║ WWW     : https://lothar-team.pl                               ║
 * ║ Forum   : https://forum.lothar-team.pl                         ║
 * ║                                                                ║
 * ║ Licencja / License: MIT                                        ║
 * ║ Rok / Year: 2026                                               ║
 * ╚════════════════════════════════════════════════════════════════╝
 */

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/user/flowchart/internal/editor"
)

var version = "1.3.0"

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
