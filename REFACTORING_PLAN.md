# Refactoring & Improvement Plan

This document outlines the planned technical improvements for the Quazaar codebase, along with estimated time efforts for completion.

## Overview

**Total Estimated Duration:** 10-12 Days

The goal of this phase is to improve code maintainability, testability, and robustness by introducing standard patterns like Dependency Injection, Structured Logging, and Centralized Configuration.

## Detailed Tasks & Estimates

### 1. Refactor Configuration Management (1 Day)

- **Goal:** Centralize all configuration (Env vars, defaults) into a single package.
- **Tasks:**
  - Create `pkg/config` package.
  - Define a `Config` struct with tags for validation.
  - Implement a `Load()` function to parse environment variables.
  - Replace usage of `os.Getenv` and `godotenv` in `main.go` and other files.

### 2. Structured Logging (1 Day)

- **Goal:** Replace standard `log` and `fmt` printing with a structured logger (`log/slog`).
- **Tasks:**
  - Initialize a global or injectable logger.
  - Replace all `log.Println`, `log.Fatal`, `fmt.Println` with `slog.Info`, `slog.Error`, etc.
  - ensure logs are machine-parsable (JSON) in production.

### 3. Standardized API Responses (1 Day)

- **Goal:** Ensure all API endpoints return a consistent JSON structure.
- **Tasks:**
  - Create `pkg/response` or similar utility package.
  - Define standard response functions: `JSON(w, code, data)` and `Error(w, code, message)`.
  - Refactor existing handlers to use these helpers instead of manual `json.NewEncoder`.

### 4. Modularize Routing (1 Day)

- **Goal:** Decompose `router.go` into domain-specific route setup functions.
- **Tasks:**
  - Create `internal/player/routes.go`, `internal/spotify/routes.go`, `internal/auth/routes.go`.
  - Move route registration from `internal/api/router.go` to these respective packages.
  - Update `api.SetupRoutes` to call these sub-setup functions.

### 5. Dependency Injection for Handlers (2-3 Days)

- **Goal:** Remove global variables (like `auth.DB`) and enable better testing.
- **Tasks:**
  - Define a `Server` or `App` struct holding dependencies (Config, DB, Logger).
  - Refactor handlers to be methods of `Server` or structs that accept dependencies.
  - Remove global `DB` variables in `auth` package.

### 6. Interface-Based System Info (2 Days)

- **Goal:** Decouple system information logic from the hardware for testing.
- **Tasks:**
  - Define interfaces for `CPU`, `RAM`, `Disk` info retrievers.
  - Make the current implementations satisfy these interfaces.
  - Inject these interfaces into the usage handlers.
  - Allow for Mock implementations for unit tests.

### 7. Explicit Context Usage (1 Day)

- **Goal:** Ensure long-running operations can be cancelled / timed out.
- **Tasks:**
  - Review database calls and switch to `QueryContext` / `ExecContext`.
  - Ensure `request.Context()` is passed down from handlers to the DB layer.

### 8. Graceful Shutdown (1 Day)

- **Goal:** Allow the server to finish active requests before stopping.
- **Tasks:**
  - Implement `os/signal` handling in `main.go`.
  - Use `server.Shutdown(ctx)` with a timeout.
  - Ensure background pollers (e.g., media poller) handle context cancellation.

## Schedule Suggestion

| Day         | Task                              |
| :---------- | :-------------------------------- |
| **Day 1**   | Refactor Configuration Management |
| **Day 2**   | Structured Logging                |
| **Day 3**   | Standardized API Responses        |
| **Day 4**   | Modularize Routing                |
| **Day 5-6** | Dependency Injection for Handlers |
| **Day 7**   | Explicit Context Usage            |
| **Day 8-9** | Interface-Based System Info       |
| **Day 10**  | Graceful Shutdown                 |
