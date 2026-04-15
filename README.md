# Redwall

Mit Redwall können Wallpaper von https://www.reddit.com/r/wallpaper/ direkt als Desktop-Hintergrund gesetzt werden.
Aktuell ist das Projekt für KDE ausgelegt.

## Abhängigkeiten

Fyne braucht zum Bauen einen C-Compiler und OpenGL-Header. Auf Arch:
```
pacman -S mesa libglvnd gcc
```

## Starten

```
go run ./cmd/gui/
```

Unter Wayland kann es nötig sein, den Software-Renderer zu erzwingen:
```
FYNE_RENDERER=software go run ./cmd/gui/
```

## Bauen

```
go build -o redwall ./cmd/gui/
```
