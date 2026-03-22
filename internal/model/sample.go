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
 * ║  Plik / File: sample.go                                      ║
 * ║                                                              ║
 * ║  Licencja / License: MIT                                     ║
 * ║  Rok / Year: 2025-2026                                       ║
 * ╚══════════════════════════════════════════════════════════════╝
 */

package model

import "image/color"

// LoadSample populates d with a compact version of the SunGo Macro PAD flowchart.
func LoadSample(d *Diagram) {
	add := func(x, y, w, h float64, label, sub string, shape Shape, col color.RGBA, anim Anim) *Node {
		n := &Node{
			X: x, Y: y, W: w, H: h,
			Label:     label,
			Sub:       sub,
			Shape:     shape,
			Color:     col,
			FillColor: DarkFill(col),
			TextColor: color.RGBA{220, 230, 255, 255},
			Anim:      anim,
			AnimSpeed: 0.8, // slower = less visual noise
		}
		d.AddNode(n)
		return n
	}

	lb := color.RGBA{80, 180, 255, 255}   // light blue
	ind := color.RGBA{63, 84, 186, 255}   // indigo
	red := color.RGBA{220, 50, 50, 255}   // red
	grn := color.RGBA{50, 200, 80, 255}   // green
	ora := color.RGBA{230, 140, 40, 255}  // orange
	wht := color.RGBA{200, 200, 220, 255} // white
	blu := color.RGBA{40, 120, 220, 255}  // blue

	// ── Main state ──────────────────────────────────────────────────────────
	main := add(460, 20, 160, 52, "STAN GŁÓWNY", "BEZCZYNNOŚĆ", ShapeRounded, ind, AnimGlow)

	// ── Linter branch ───────────────────────────────────────────────────────
	lintQ := add(60, 110, 110, 52, "Błędy\nLintera?", "", ShapeDiamond, lb, AnimNone)
	lintErr := add(40, 210, 130, 48, "BŁĄD LINTER", "Kolor: RED", ShapeRounded, red, AnimPulse)
	fixQ := add(52, 300, 106, 48, "Naprawiono?", "", ShapeDiamond, lb, AnimNone)

	// ── GUI branch ──────────────────────────────────────────────────────────
	guiQ := add(240, 110, 110, 52, "Które okno?", "", ShapeDiamond, lb, AnimNone)
	creator := add(195, 210, 128, 48, "KREATOR", "Ciągłe", ShapeRounded, wht, AnimGlow)
	dashboard := add(340, 210, 128, 48, "DASHBOARD", "Ciągłe", ShapeRounded, grn, AnimGlow)

	// ── Simple actions ───────────────────────────────────────────────────────
	cmdQ := add(490, 110, 110, 52, "Komenda?", "", ShapeDiamond, lb, AnimNone)
	fmtN := add(440, 210, 128, 48, "FORMATOWANIE", "Flash", ShapeRounded, ora, AnimFlash)
	impN := add(590, 210, 128, 48, "IMPORTOWANIE", "Spinner", ShapeRounded, lb, AnimSpinner)

	// ── Long processes ───────────────────────────────────────────────────────
	procQ := add(760, 110, 110, 52, "Proces?", "", ShapeDiamond, lb, AnimNone)
	runN := add(700, 210, 118, 48, "RUN", "Pulse", ShapeRounded, grn, AnimPulse)
	buildN := add(832, 210, 118, 48, "BUILD", "Blink", ShapeRounded, blu, AnimBlink)
	testN := add(964, 210, 118, 48, "TEST", "Pulse", ShapeRounded, wht, AnimPulse)

	// ── Results ──────────────────────────────────────────────────────────────
	runRes := add(710, 310, 100, 46, "Wynik\nRUN?", "", ShapeDiamond, lb, AnimNone)
	buildRes := add(836, 310, 100, 46, "Wynik\nBUILD?", "", ShapeDiamond, lb, AnimNone)
	testRes := add(962, 310, 100, 46, "Wynik\nTEST?", "", ShapeDiamond, lb, AnimNone)

	okN := add(740, 410, 120, 48, "WYNIK OK", "", ShapeRounded, grn, AnimFlash)
	failN := add(880, 410, 120, 48, "WYNIK BŁĄD", "", ShapeRounded, red, AnimFlash)

	// ── Edges ─────────────────────────────────────────────────────────────────
	conn := func(from, to *Node, label string) {
		d.AddEdge(&Edge{
			FromID: from.ID, ToID: to.ID,
			Label: label,
			Color: color.RGBA{80, 150, 255, 180},
		})
	}

	conn(main, lintQ, "Wykryto błędy")
	conn(main, guiQ, "Otwarcie okna")
	conn(main, cmdQ, "Akcja")
	conn(main, procQ, "Proces")

	conn(lintQ, lintErr, "TAK")
	conn(lintQ, main, "NIE")
	conn(lintErr, fixQ, "")
	conn(fixQ, main, "TAK")
	conn(fixQ, lintErr, "NIE")

	conn(guiQ, creator, "Creator")
	conn(guiQ, dashboard, "Dashboard")
	conn(creator, main, "Zamknięcie")
	conn(dashboard, main, "Zamknięcie")

	conn(cmdQ, fmtN, "GoFMT")
	conn(cmdQ, impN, "Import")
	conn(fmtN, main, "Koniec")
	conn(impN, main, "Koniec")

	conn(procQ, runN, "RUN")
	conn(procQ, buildN, "BUILD")
	conn(procQ, testN, "TEST")

	conn(runN, runRes, "")
	conn(buildN, buildRes, "")
	conn(testN, testRes, "")

	conn(runRes, okN, "OK")
	conn(runRes, failN, "FAIL")
	conn(buildRes, okN, "OK")
	conn(buildRes, failN, "FAIL")
	conn(testRes, okN, "OK")
	conn(testRes, failN, "FAIL")

	conn(okN, main, "")
	conn(failN, main, "")
}
