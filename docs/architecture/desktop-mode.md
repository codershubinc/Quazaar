# Desktop Mode Architecture & Setup

Use this layout to add a system-tray desktop mode while keeping the existing server entrypoint unchanged.

## Target Layout

```text
internal/
  desktop/
    desktop.go              # Interface + shared wiring
    desktop_linux.go        # //go:build linux
    desktop_windows.go      # //go:build windows
    desktop_darwin.go       # //go:build darwin
    icon/
      icon.go               # Interface + helpers
      icon_linux.go         # Linux icon loading
      icon_windows.go       # Windows icon loading
      icon_darwin.go        # macOS icon loading
    tray/
      tray.go               # Interface + shared tray glue
      tray_linux.go         # Linux tray using systray/ayatana
      tray_windows.go       # Windows tray
      tray_darwin.go        # macOS tray
cmd/
  desktop/
    main.go                 # Desktop entrypoint (starts server + tray)
assets/
  images/
    icon.ico                # Tray icon (fallback to embedded if missing)
```

## Responsibilities

- `internal/desktop/desktop.go`: define a small interface (Init(url string), Quit()) and hold shared helpers.
- OS-specific files implement the interface; Go selects them via build tags.
- `icon/` handles loading an icon file first, then falls back to embedded bytes.
- `tray/` wires menu items (Open UI, Status, Quit) and browser launching.
- `cmd/desktop/main.go` mirrors [cmd/server/main.go](../cmd/server/main.go) but starts the HTTP server in a goroutine, then blocks in the tray.

## Project-Wide Changes Needed

- Create `cmd/desktop/main.go` that starts the server in a goroutine, then blocks in the tray init.
- Add `internal/desktop/**` with shared interface plus OS-tagged files (`*_linux.go`, `*_windows.go`, `*_darwin.go`) for tray and icon loading; identical function signatures across OS builds.
- Keep existing OS-specific patterns for media info (e.g., `media_info_linux.go`, `media_info_windows.go`, `media_info_darwin.go`); callers keep using the same APIs.
- Move internal-only helpers/models out of `pkg/` into `internal/` to avoid exporting unintended APIs (e.g., `internal/shared/helpers`, `internal/shared/models` or colocate per domain).
- Ensure assets contain a tray icon at `assets/images/icon.ico`; keep an embedded fallback in code.

## Apply Across Repo (by area)

- Entry points: keep server at [cmd/server](../cmd/server) and add desktop at `cmd/desktop/main.go`; everything else remains in `internal/**`.
- Desktop stack: `internal/desktop/desktop.go` defines the interface; OS-specific impls live in `*_linux.go`, `*_windows.go`, `*_darwin.go` under `desktop/`, `desktop/icon/`, and `desktop/tray/` with identical exported signatures.
- Media info: keep platform files like `media_info_linux.go`, `media_info_windows.go`, `media_info_darwin.go` with build tags; the API surface stays the same for callers.
- Shared code location: relocate app-only helpers/models from [`pkg`](../pkg) into [`internal`](../internal) (e.g., `internal/shared/helpers`, `internal/shared/models`, or per-domain folders) to avoid accidental external imports.
- Assets: place tray icon at `assets/images/icon.ico`; code should attempt file load first, then fallback bytes.

## Linux Notes

- Dependency: `github.com/getlantern/systray`.
- System libs (typical): `libayatana-appindicator3-dev` and `libgtk-3-dev`.
- Build tag on Linux implementation: `//go:build linux` and `// +build linux` on the first lines.
- Browser open helper should call `xdg-open`.

## Build & Run

- Fetch dependency: `go get github.com/getlantern/systray`.
- Build desktop binary: `go build -o quazaar-desktop ./cmd/desktop`.
- Run desktop mode: `./quazaar-desktop` (server auto-starts; tray stays resident).

## Testing Checklist

- Server reachable at `http://127.0.0.1:<LOCAL_HOST_PORT>` when tray is running.
- Tray menu opens browser to the local URL.
- Quit menu cleanly stops the process.
- Icon loads from assets; if missing, embedded fallback is shown.
- Cross-compile sanity: linux, windows, darwin all build with their tagged files.
