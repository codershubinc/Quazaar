# Project Improvements & Recommendations

Based on a review of the codebase and the current `todo.md`, here is a detailed list of recommended technical improvements to enhance stability, maintainability, and performance.

## 1. Graceful Shutdown & Context Propagation

**Current State:**
The server uses `http.ListenAndServe` directly, which blocks until the process is killed. There is no mechanism to catch interrupt signals (SIGINT/SIGTERM), meaning the server and background tasks (like the media poller) are terminated abruptly.

**Recommendation:**
- Implement `os/signal` handling to catch interrupt signals.
- Use `context.WithCancel` to create a root context that is canceled upon receiving a signal.
- Pass this context to background workers (e.g., `internal/poller`) so they can stop their work cleanly.
- Use `server.Shutdown(ctx)` to gracefully stop the HTTP server, allowing active requests to finish within a timeout period.

## 2. Unified Logging System

**Current State:**
The codebase mixes `fmt.Println`, `log.Println`, and `fmt.Printf`. This makes it difficult to control log levels (info, error, debug) or parse logs programmatically.

**Recommendation:**
- Adopt a structured logging library like `log/slog` (standard in Go 1.21+) or `zerolog`.
- Create a dedicated `internal/logger` package.
- Replace all `fmt` and `log` print calls with the new logger (e.g., `logger.Info`, `logger.Error`).
- Configure the logger to output JSON in production and text in development.

## 3. Configuration Management

**Current State:**
Environment variables are accessed directly via `os.Getenv` in multiple places (e.g., `cmd/server/main.go`). `godotenv.Load()` is called in `main`, but variables are not centralized.

**Recommendation:**
- Create an `internal/config` package.
- Define a `Config` struct that holds all application configuration (Port, IP, Spotify credentials, etc.).
- Implement a `Load()` function that reads environment variables once at startup and populates the struct.
- Pass this config object or access the global singleton throughout the application, ensuring type safety and centralized defaults.

## 4. Streaming File Uploads (Memory Optimization)

**Current State:**
The `internal/fileshare` handler uses `r.ParseMultipartForm(100 << 20)`. This parses the entire file upload into memory (up to 100MB) or temporary files before processing. This can cause high memory spikes for large file transfers.

**Recommendation:**
- Switch to using `r.MultipartReader()` for streaming uploads.
- Read the file stream part-by-part.
- Update `pkg/helpers/store_file.go` to accept an `io.Reader` instead of a `multipart.File`.
- Use `io.Copy` to stream data directly from the network request to the file on disk.

## 5. Build Optimization

**Recommendation:**
- As noted in `todo.md`, use linker flags to reduce binary size: `-ldflags="-s -w"`.
- Consider using `-trimpath` to remove file system paths from the executable.

## 6. Code Consistency & Cleanup

- **Remove Internal Monologues/Comments:** Ensure production code doesn't contain "scratchpad" comments.
- **Error Handling:** Ensure all errors are logged properly (using the unified logger) rather than just printed or ignored.
- **Context Usage:** Pass `context.Context` as the first argument to functions that involve I/O or long-running processes.

## 7. Future Considerations (Phase 4)

- **Unit Tests:** Add `_test.go` files, particularly for the `auth` and `helpers` packages.
- **Linter:** Integrate `golangci-lint` to catch static analysis issues early.
- **Makefile:** Create a `Makefile` to standardize commands like `make build`, `make run`, and `make test`.
- **Contribution Standards:** Add a `.github/PULL_REQUEST_TEMPLATE.md` to standardize contributions (Added).
