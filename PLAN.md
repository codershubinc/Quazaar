# Implementation Plan

Based on the review of `todo.md` and the codebase, here is the plan to address pending tasks for Phase 2 (Systems Engineering) and Phase 3 (Dec 8 tasks).

## Phase 2: Systems Engineering (Dec 7 Backlog)

1. **Graceful Shutdown**
   - Implement `signal.Notify` to catch SIGINT/SIGTERM in `cmd/server/main.go`.
   - Use `context.WithCancel` for orchestration.
   - Add `server.Shutdown()` for HTTP server.
   - Refactor `internal/poller` to accept a context or a stop channel that is properly managed.

2. **Context Propagation**
   - Pass `ctx` to `poller` functions to support cancellation.
   - Ensure the poller stops when the context is canceled.

3. **Streaming File Uploads**
   - Rewrite `HandleTempFileShareAccept` in `internal/fileshare/api.go`.
   - Use `r.MultipartReader()` and `io.Copy` instead of `r.ParseMultipartForm` to handle large files with minimal memory usage.

## Dec 8 Tasks

4. **Unified Logger**
   - Create `internal/logger/logger.go` using `log/slog`.
   - Replace `fmt.Println` and `log.Println` across the codebase with the new logger.
   - Configure different log levels and formats (JSON/Text).
   - Ensure consistent logging format.

5. **Config Struct**
   - Create `internal/config/config.go` to load environment variables once.
   - Define a `Config` struct with typed fields.
   - Replace direct `os.Getenv` calls with config struct access.
   - Remove redundant `godotenv.Load()` calls.

## Execution Order

1. **Graceful Shutdown & Context Propagation** (Highest priority for stability)
2. **Unified Logger** (To improve observability during other changes)
3. **Config Struct** (To clean up configuration management)
4. **Streaming File Uploads** (Optimization)
