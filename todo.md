# Todo List

## 📅 Dec 6, 2025 - Phase 1: Stability

- [x] **Implement Volume Control**
- [x] **Integrate brightnessctl**
- [x] **Add TODO comment to websocket handler**
- [x] **Fix go.mod Version**
  - Issue: `go 1.25.3` is invalid.
  - Fix: Downgrade to `1.23.0` or `1.22`.
- [ ] **Unified Logger**
  - Issue: Inconsistent logging (`fmt.Println` vs `log.Println`).
  - Fix: Adopt `slog` or `zerolog` (Go 1.21+). Replace `fmt.Println`.
- [ ] **Config Struct**
  - Issue: Redundant `godotenv.Load()`.
  - Fix: Create `internal/config/config.go` to load ENV once.
- [x] **Fix Network Listener**
  - Issue: Hardcoded `127.0.0.1`.
  - Fix: Change default bind address to `0.0.0.0:8765`.

## 📅 Dec 7, 2025 - Phase 2: Systems Engineering

- [ ] **Graceful Shutdown**
  - Fix: Implement `signal.Notify`, `context.WithCancel`, and `server.Shutdown()`.
- [ ] **Streaming File Uploads**
  - Issue: `ParseMultipartForm` uses too much RAM.
  - Fix: Rewrite `HandleTempFileShareAccept` to use `io.Copy`.
- [ ] **Context Propagation**
  - Fix: Pass `ctx` to poller.

## 🔮 Future / Upcoming - Phase 3: Professionalism

- [ ] **Unit Tests**
  - Fix: Write `_test.go` for Auth logic.
- [ ] **Makefile**
  - Fix: Add Makefile for build/run/test.
- [ ] **Linter**
  - Fix: Run `golangci-lint` and fix warnings.
- [ ] **Update Web Test Client** (Dynamic deviceId/token)
- [ ] **Verify Spotify Auth Flow**
