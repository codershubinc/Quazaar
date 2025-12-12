# Changelog v0.2.1-nightly "Pulsar"

**Release Date:** December 12, 2025
**Status:** 🔴 Nightly - Highly Unstable

> ⚠️ **Warning:** This is a nightly build containing major architectural changes to the Windows subsystem. It is highly unstable and may crash frequently. Use for testing and development purposes only.

## 🚀 Highlights

This release focuses on stabilizing the Windows experience and establishing a robust architecture for the Windows Helper sidecar. We've moved from a monolithic script to a modular, service-oriented architecture in C# and enhanced the Go sidecar manager with enterprise-grade error reporting and crash handling.

- **Modular Windows Helper:** Complete rewrite of the C# helper into a clean, service-based architecture.
- **Robust Sidecar Management:** The Go backend now handles sidecar crashes, errors, and lifecycle events with structured JSON logs.
- **Single Binary Distribution:** The Windows helper is now embedded directly into the Go binary using `//go:embed`.
- **New Build System:** Introduced `build.ps1` for one-step cross-compilation and bundling.

---

## 🛠️ Refactoring & Architecture

### Windows Helper (C#)

- **Service-Oriented Design:** Split the monolithic `Program.cs` into:
  - `MediaController.cs`: Main orchestrator.
  - `Services/PlaybackControlService.cs`: Handles Play, Pause, Next, Prev, Seek.
  - `Services/MediaStatusService.cs`: Handles metadata extraction and status reporting.
  - `Utilities/JsonHelper.cs`: Centralized JSON escaping and formatting.
  - `Models.cs`: Strongly typed command and response models.
- **Improved Serialization:** Switched to `System.Text.Json` with Source Generators for better performance and AOT compatibility.
- **Cleanup:** Removed unused legacy files and build artifacts.

### Sidecar Manager (Go)

- **Embedded Binary:** `QuazaarMedia.exe` is now bundled inside `quazaar.exe`.
- **Lifecycle Management:** Improved startup handshake ("ready" signal) and timeout handling.
- **Crash Handling:** New `HandleCrash` function captures and logs sidecar failures with:
  - Component name
  - Failure stage (extraction, pipe creation, handshake)
  - Timestamp
  - Error details
- **JSON Error Reporting:** Errors from the sidecar are now parsed as structured JSON, allowing for precise debugging.

---

## 🐛 Bug Fixes

- **Error Formatting:** Fixed a bug in `provider_windows.go` where error messages were malformed (`%!(EXTRA ...)`).
- **JSON Escaping:** Fixed issues where special characters in song titles or artist names could break the JSON response from the C# helper.
- **Build Script:** Fixed encoding issues in `build.ps1` that caused execution errors on some PowerShell configurations.
- **Null Safety:** Added comprehensive null checks in the C# helper to prevent crashes when media properties are missing.

---

## 📦 Build System

- Added `build.ps1` PowerShell script.
- **Capabilities:**
  1.  Compiles C# project (`cmd/windows-helper`) to a single `.exe`.
  2.  Copies the `.exe` to `internal/sidecar/`.
  3.  Builds the final Go binary (`quazaar.exe`) with the embedded sidecar.

---

## 📝 Known Issues

- **Sidecar Extraction:** On rare occasions, if `quazaar.exe` is restarted very quickly, the OS might lock `QuazaarMedia.exe` in the temp folder, causing an "Access Denied" error during extraction.
- **Multi-Session:** The current implementation selects the system's "active" media session. Support for selecting specific sessions (if multiple apps are playing) is planned for a future release.
