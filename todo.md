# Todo List - December 6, 2025

- [x] **Implement Volume Control**

  - File: `internal/system/volume/volume.go`, `internal/system/volume/handle_volume.go`
  - Task: Implemented `CurrentSystemVolume()`, `SickSystemSetVolume()`, `IncreaseSystemVolume()`, `DecreaseSystemVolume()` with pactl integration.
  - Task: Added WebSocket handlers (`HandleVolumeWS`) and integrated into system router.
  - Task: Updated `temp/web/index.html` with volume controls for testing.

- [x] **Integrate brightnessctl**

  - Source: `docs/dump/todo.txt`
  - Task: Investigate and implement brightness control integration.
  - Implementation:
    - Created `internal/system/brightness/brightness.go` with core logic.
    - Created `internal/system/brightness/handle_brightness.go` for WebSocket handling.
    - Registered handler in `internal/system/system_wesocket_handler.go`.
    - Updated `temp/web/index.html` with brightness controls.
    - Documented in `docs/features/brightness-control.md`.

- [ ] **Refactor Spotify Logging**

  - File: `internal/spotify/spotify.go`
  - Task: Replace `fmt.Println` with `helpers.LogMessage` for consistent logging.

- [ ] **Update Web Test Client**

  - File: `temp/web/index.html`
  - Task: Replace hardcoded device ID with dynamic input or token handling.

- [ ] **Verify Spotify Auth Flow**

  - Task: Verify the redirection logic when tokens are missing (currently just logs a warning).

- [x] **Add TODO comment to websocket handler**
  - File: `internal/websocket/handler.go`
  - Task: Add TODO note for future unified switch-case refactoring (legacy support).
