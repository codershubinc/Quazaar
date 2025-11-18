package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// ANSI color codes for visual formatting
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

// LogLevel represents the severity of the log message
type LogLevel int

const (
	INFO LogLevel = iota
	WARNING
	ERROR
	FATAL
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// getColor returns the ANSI color code for the log level
func (l LogLevel) getColor() string {
	switch l {
	case INFO:
		return colorGreen
	case WARNING:
		return colorYellow + colorBold
	case ERROR:
		return colorRed + colorBold
	case FATAL:
		return colorPurple + colorBold
	default:
		return colorWhite
	}
}

// getEmoji returns the emoji for the log level
func (l LogLevel) getEmoji() string {
	switch l {
	case INFO:
		return "✅"
	case WARNING:
		return "⚠️"
	case ERROR:
		return "❌"
	case FATAL:
		return "💀"
	default:
		return "📝"
	}
}

// isColorSupported checks if the terminal supports colors
func isColorSupported() bool {
	term := os.Getenv("TERM")
	return strings.Contains(term, "color") || strings.Contains(term, "xterm") || os.Getenv("COLORTERM") != ""
}

// RequestLogHandler is a middleware that logs HTTP requests with customizable messages
func RequestLogHandler(level LogLevel, message string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Call the next handler
			next(rw, r)

			// Calculate duration
			duration := time.Since(start)

			// Format log message based on level with colors
			var logMessage string
			if isColorSupported() {
				logMessage = fmt.Sprintf("%s%s %s%s | %s %s | Status: %d | Duration: %v | IP: %s | User-Agent: %s",
					level.getColor(),
					level.getEmoji(),
					level.String(),
					colorReset,
					r.Method,
					r.URL.Path,
					rw.statusCode,
					duration,
					getClientIP(r),
					r.Header.Get("User-Agent"),
				)
			} else {
				logMessage = fmt.Sprintf("%s [%s] %s | %s %s | Status: %d | Duration: %v | IP: %s | User-Agent: %s",
					level.getEmoji(),
					level.String(),
					message,
					r.Method,
					r.URL.Path,
					rw.statusCode,
					duration,
					getClientIP(r),
					r.Header.Get("User-Agent"),
				)
			}

			// Log the message with request details
			log.Println(logMessage)
		}
	}
}

// LogMessage logs a custom message with specified level and visual formatting
func LogMessage(level LogLevel, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)

	var finalMessage string
	if level == WARNING && isColorSupported() {
		// Create a boxed warning for better visibility
		boxedMessage := createBoxedWarning(formattedMessage)
		finalMessage = fmt.Sprintf("%s%s%s", level.getColor(), boxedMessage, colorReset)
	} else {
		if isColorSupported() {
			finalMessage = fmt.Sprintf("%s%s %s:%s %s",
				level.getColor(),
				level.getEmoji(),
				level.String(),
				colorReset,
				formattedMessage)
		} else {
			finalMessage = fmt.Sprintf("%s %s: %s",
				level.getEmoji(),
				level.String(),
				formattedMessage)
		}
	}

	log.Println(finalMessage)
}

// createBoxedWarning creates a visually appealing boxed warning message
func createBoxedWarning(message string) string {
	content := fmt.Sprintf("⚠️  WARNING: %s", message)
	width := len(content) + 4 // Add padding

	// Create top border
	topBorder := "+" + strings.Repeat("-", width-2) + "+"

	// Create content line with padding
	contentLine := fmt.Sprintf("| %-*s |", width-4, content)

	// Create bottom border
	bottomBorder := "+" + strings.Repeat("-", width-2) + "+"

	// Combine all lines
	return fmt.Sprintf("\n%s\n%s\n%s", topBorder, contentLine, bottomBorder)
}

// LogRequest logs an HTTP request with custom message and level
func LogRequest(level LogLevel, message string, r *http.Request, statusCode int, duration time.Duration) {
	var logMessage string

	if isColorSupported() {
		logMessage = fmt.Sprintf("%s%s %s%s | %s %s | Status: %d | Duration: %v | IP: %s",
			level.getColor(),
			level.getEmoji(),
			level.String(),
			colorReset,
			r.Method,
			r.URL.Path,
			statusCode,
			duration,
			getClientIP(r),
		)
	} else {
		logMessage = fmt.Sprintf("%s [%s] %s | %s %s | Status: %d | Duration: %v | IP: %s",
			level.getEmoji(),
			level.String(),
			message,
			r.Method,
			r.URL.Path,
			statusCode,
			duration,
			getClientIP(r),
		)
	}

	log.Println(logMessage)
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take the first IP if multiple are present
		if idx := strings.Index(xff, ","); idx > 0 {
			return strings.TrimSpace(xff[:idx])
		}
		return xff
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// Example usage functions for different scenarios

// LogAuthAttempt logs authentication attempts with visual warnings
func LogAuthAttempt(success bool, username string, ip string) {
	if success {
		LogMessage(INFO, "🔐 Authentication successful for user: %s from IP: %s", username, ip)
	} else {
		LogMessage(WARNING, "🚫 Authentication failed for user: %s from IP: %s", username, ip)
	}
}

// LogAPIError logs API errors with context and visual alerts
func LogAPIError(endpoint string, err error, statusCode int) {
	LogMessage(ERROR, "💥 API Error in %s: %v (Status: %d)", endpoint, err, statusCode)
}

// LogSecurityEvent logs security-related events with high visibility
func LogSecurityEvent(event string, details string, ip string) {
	LogMessage(WARNING, "🛡️  Security Event: %s | Details: %s | IP: %s", event, details, ip)
}

// LogPerformanceWarning logs performance warnings with timing alerts
func LogPerformanceWarning(endpoint string, duration time.Duration, threshold time.Duration) {
	LogMessage(WARNING, "⏱️  Performance Warning: %s took %v (threshold: %v)", endpoint, duration, threshold)
}

// LogSystemStartup logs system startup events
func LogSystemStartup(component string, version string) {
	LogMessage(INFO, "🚀 %s v%s started successfully", component, version)
}

// LogDatabaseError logs database errors with critical alerts
func LogDatabaseError(operation string, err error) {
	LogMessage(ERROR, "💾 Database Error during %s: %v", operation, err)
}

// LogNetworkError logs network-related errors
func LogNetworkError(service string, err error) {
	LogMessage(ERROR, "🌐 Network Error with %s: %v", service, err)
}
