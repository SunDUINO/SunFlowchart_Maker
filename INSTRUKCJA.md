# SunFlowChart Editor — v. 1.1.0 Instrukcja\
# SunRiver / Lothar TeaM


---

## Pasek narzędzi (lewy)

| Ikona | Nazwa | Klawisz | Opis |
|-------|-------|---------|------|
| ↖ | Wybierz | `V` / `1` | Zaznacz, przeciągaj, zmieniaj rozmiar węzłów |
| □+ | Węzeł | `N` / `2` | Kliknij na canvas → dodaje węzeł (od razu edycja etykiety) |
| T̲ | Tytuł | `T` / `3` | Dodaje specjalny węzeł tytułowy (złoty, bez ramki) |
| ⤷ | Połącz | `C` / `4` | Kliknij węzeł źródłowy, potem docelowy → strzałka |
| ✥ | Przesuń | `H` / `5` | Przeciągaj canvas (też kółko myszy) |
| ✕ | Usuń | `X` / `6` | Kliknij węzeł lub strzałkę → usuwa |
| ⌫ | Wyczyść | `7` | Czyści cały ekran (Ctrl+Z cofa) |
| 📁 | Wczytaj | `Ctrl+O` | Otwiera dialog wczytania pliku JSON |

### Przyciski na dole paska
| Przycisk | Klawisz | Opis |
|----------|---------|------|
| Siatka (dioda) | `G` | Pokazuje/ukrywa siatkę punktów |
| Snap (dioda) | `Shift+G` | Włącza/wyłącza przyciąganie do siatki (24px) |

---

## Kształty węzłów
Aktywne gdy narzędzie **Węzeł** — wybierz kształt w dolnej części paska:

`▬` Prostokąt · `▭` Zaokrąglony · `◆` Romb · `▱` Równoległobok · `⬭` Owal

---

## Edycja węzłów

| Akcja | Efekt |
|-------|-------|
| Lewy klik | Zaznacz węzeł |
| Prawy klik (1×) | Zaznacz + panel właściwości |
| Prawy klik (2×) | Edycja etykiety |
| `F2` | Edycja etykiety zaznaczonego węzła |
| `Tab` w edycji | Przełącz etykieta ↔ pod-etykieta |
| `Enter` | Zatwierdź tekst |
| `Esc` | Anuluj edycję / odznacz |
| `Del` | Usuń zaznaczone |
| Strzałki ↑↓←→ | Nudge o 4px (+ `Shift` = 1px) |
| Żółty uchwyt (róg) | Przeciągnij → zmień rozmiar węzła |
| `Shift+klik` | Dodaj do zaznaczenia |
| Przeciąganie po canvas | Zaznaczenie prostokątne |
| `Ctrl+A` | Zaznacz wszystko |

---

## Panel właściwości (prawy bok — gdy węzeł zaznaczony)

- **Kolor** — 12 kolorów neonowych → kolor ramki i glow
- **Tło** — 8 kolorów tła canvas
- **Animacja** — Off · Glow · Pulse · Blink · Flash · Spin
- **Ramka** — grubość 1px / 2px / 3px / 4px

### Gdy zaznaczona strzałka
- Kolor strzałki (8 kolorów)
- Styl: `~ Krzywa` lub `| Łamana`
- Prawy klik na strzałkę → edycja etykiety strzałki

---

## Połączenia (strzałki)

1. Wybierz narzędzie **Połącz** (`C`)
2. Kliknij węzeł źródłowy — zielony pulsujący pierścień
3. Najedź na węzeł docelowy — żółte podświetlenie
4. Kliknij węzeł docelowy → strzałka gotowa
5. `Esc` = anuluj rysowanie

Wiele strzałek między tymi samymi węzłami jest automatycznie rozdzielanych.

---

## Tytuł i Legenda

| Element | Jak dodać | Jak edytować |
|---------|-----------|--------------|
| **Tytuł** | Narzędzie Tytuł (`T`), kliknij canvas | Prawy klik 2× lub `F2` |
| **Legenda** | Klawisz `L` lub klik panelu legenda | `Shift+Enter` = nowa linia, `Enter` = zatwierdź |

Legenda wyświetlana jest w lewym dolnym rogu ekranu i eksportowana do SVG.

---

## Zapis i Export

| Skrót | Akcja |
|-------|-------|
| `Ctrl+S` / `F5` | Zapisz diagram → dialog z nazwą pliku → `nazwa.json` |
| `Ctrl+E` / `F6` | Eksport → dialog → wybierz PNG lub SVG |
| `Ctrl+O` | Wczytaj → dialog z nazwą pliku → `nazwa.json` |

Pliki zapisywane są w katalogu z którego uruchomiono program.

---

## Undo / Redo

| Skrót | Akcja |
|-------|-------|
| `Ctrl+Z` | Cofnij |
| `Ctrl+Shift+Z` | Ponów |

---

## Widok i nawigacja

| Akcja | Efekt |
|-------|-------|
| Kółko myszy | Przewijanie pionowe |
| Środkowy przycisk myszy (drag) | Swobodne przesuwanie canvas |
| Narzędzie Przesuń + drag | Przesuwanie canvas |
| `G` | Siatka widoczna/ukryta |
| `Shift+G` | Snap do siatki wł./wył. |

---

## Skróty klawiszowe — podsumowanie

```
V/1  Wybierz       N/2  Węzeł        T/3  Tytuł
C/4  Połącz        H/5  Przesuń      6    Usuń
7    Wyczyść       L    Legenda      G    Siatka
Shift+G  Snap      F2   Etykieta     Del  Usuń
F5   Zapis         F6   Export       Ctrl+O  Wczytaj
Ctrl+Z  Cofnij     Ctrl+Shift+Z  Ponów
```

