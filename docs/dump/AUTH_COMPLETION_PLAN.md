# 🔐 Authentication System Completion Plan

## Current Status

### ✅ What's Already Implemented

- **User Management**

  - User registration (`POST /api/register`)
  - User login with automatic auth token creation (`POST /api/login`)
  - Password hashing with bcrypt
  - Single-user system (local server)

- **Token Management**

  - Token generation (cryptographically secure)
  - Token creation for services (`POST /api/tokens/create`)
  - Token validation (`ValidateToken()` function)
  - Token revocation (`POST /api/tokens/revoke`)
  - Token listing (`GET /api/tokens/list`)
  - Automatic token cleanup (expired tokens)

- **Database Layer**
  - SQLite database with users and tokens tables
  - Proper indexes for performance
  - Token expiration tracking
  - Last used timestamp tracking

### ❌ What's Missing

1. **No endpoint protection** - WebSocket and player control endpoints are completely open
2. **No authentication middleware** - Every handler manually checks auth (code duplication)
3. **No token extraction helpers** - No support for `Bearer` token format
4. **Missing utility endpoints** - No password change, token refresh, or user info endpoints
5. **Incomplete testing** - Test scripts don't verify protected endpoints work correctly

---

## 🎯 Completion Tasks

### Task 1: Create Authentication Middleware

**File:** `utils/auth/middleware.go`

**Purpose:** Centralized authentication logic to protect endpoints without code duplication

**Features:**

```go
// RequireAuth wraps HTTP handlers to enforce authentication
func RequireAuth(next http.HandlerFunc) http.HandlerFunc

// RequireAuthWithRole wraps handlers with role-based access (future)
func RequireAuthWithRole(role string, next http.HandlerFunc) http.HandlerFunc

// ExtractToken gets token from Authorization header or query parameter
func ExtractToken(r *http.Request) string

// Supports both formats:
// - Authorization: <token>
// - Authorization: Bearer <token>
// - Query param: ?token=<token>
```

**Benefits:**

- ✅ Centralized auth logic (DRY principle)
- ✅ Consistent error responses
- ✅ Easy to add/remove protection from endpoints
- ✅ Better logging and debugging

---

### Task 2: Protect WebSocket Endpoint

**File:** `utils/websocket/handler.go`

**Current State:** Completely open - anyone can connect

**Required Changes:**

```go
func Handle(res http.ResponseWriter, req *http.Request) {
    // ✅ ADD: Token validation before WebSocket upgrade
    token := auth.ExtractToken(req)
    valid, err := auth.ValidateToken(token)
    if !valid || err != nil {
        log.Printf("❌ Unauthorized WebSocket attempt: %v", err)
        http.Error(res, "Unauthorized", http.StatusUnauthorized)
        return
    }

    log.Printf("✅ WebSocket authenticated for token: %s...", token[:10])

    // Continue with existing connection logic...
    conn, err := CreateWebSocketConnection(res, req)
    // ... rest of code
}
```

**Security Impact:** **CRITICAL** - Prevents unauthorized remote control access

---

### Task 3: Protect Player Control Endpoints

**File:** `main.go`

**Current State:** All player endpoints are open

**Required Changes:**

```go
// Protected player control endpoints (require auth)
http.HandleFunc("/api/v0.1/player/play-pause", auth.RequireAuth(player.HandlePlayPause))
http.HandleFunc("/api/v0.1/player/next", auth.RequireAuth(player.HandleNext))
http.HandleFunc("/api/v0.1/player/previous", auth.RequireAuth(player.HandlePrevious))
http.HandleFunc("/api/v0.1/player/play", auth.RequireAuth(player.HandlePlay))
http.HandleFunc("/api/v0.1/player/pause", auth.RequireAuth(player.HandlePause))

// Optional: Also protect info endpoints for consistency
http.HandleFunc("/api/v0.1/player/info", auth.RequireAuth(player.HandleGetPlayerInfo))
http.HandleFunc("/api/v0.1/player/list", auth.RequireAuth(player.HandleGetActivePlayers))
```

**Decision Required:**

- **Option A:** Protect only POST endpoints (control commands) - allows public player info viewing
- **Option B:** Protect all player endpoints (GET + POST) - maximum security
- **Recommendation:** Option B for local server security

---

### Task 4: Add Missing Utility Endpoints

**File:** `utils/auth/handlers.go`

**New Endpoints:**

#### 1. Change Password

```go
// POST /api/password/change
// Body: {"old_password": "...", "new_password": "..."}
// Requires: Authorization header with valid token
func HandleChangePassword(w http.ResponseWriter, r *http.Request)
```

#### 2. Refresh Token

```go
// POST /api/tokens/refresh
// Body: {"token": "..."}
// Returns: New token with extended expiration
func HandleRefreshToken(w http.ResponseWriter, r *http.Request)
```

#### 3. Get User Info

```go
// GET /api/user/me
// Requires: Authorization header with valid token
// Returns: Current user information (username, created_at, etc.)
func HandleGetUserInfo(w http.ResponseWriter, r *http.Request)
```

#### 4. Revoke All Tokens

```go
// POST /api/tokens/revoke-all
// Requires: Authorization header with valid token
// Revokes all tokens except the one being used
func HandleRevokeAllTokens(w http.ResponseWriter, r *http.Request)
```

---

### Task 5: Update Test Scripts

**File:** `tests/test_auth.sh`

**Add Tests For:**

1. ✅ Protected endpoints return 401 without token
2. ✅ Protected endpoints return 401 with invalid token
3. ✅ Protected endpoints return 200 with valid token
4. ✅ WebSocket connection fails without token
5. ✅ WebSocket connection succeeds with valid token
6. ✅ Player control commands require authentication
7. ✅ Password change functionality
8. ✅ Token refresh functionality

**Example Test:**

```bash
# Test protected endpoint without token (should fail)
echo "Testing protected endpoint without token..."
RESPONSE=$(curl -s -X POST http://localhost:8765/api/v0.1/player/play-pause)
if [[ $RESPONSE == *"Unauthorized"* ]]; then
    echo "✅ Correctly rejected unauthorized request"
else
    echo "❌ SECURITY ISSUE: Endpoint is not protected!"
fi

# Test protected endpoint with valid token (should succeed)
echo "Testing protected endpoint with valid token..."
RESPONSE=$(curl -s -X POST http://localhost:8765/api/v0.1/player/play-pause \
  -H "Authorization: $AUTH_TOKEN")
if [[ $RESPONSE == *"success"* ]]; then
    echo "✅ Authenticated request succeeded"
else
    echo "❌ Authenticated request failed"
fi
```

---

### Task 6: Documentation

**File:** `docs/AUTH_COMPLETE.md`

**Content:**

- ✅ Complete API reference with authentication examples
- ✅ Middleware usage patterns
- ✅ Token format specifications (`Bearer` vs plain)
- ✅ WebSocket authentication examples
- ✅ Security best practices
- ✅ Troubleshooting guide
- ✅ Example client implementations (curl, JavaScript, Python)

---

## 🔧 Implementation Order

### Phase 1: Core Protection (Priority: CRITICAL)

1. Create `utils/auth/middleware.go` with `RequireAuth()` and `ExtractToken()`
2. Protect WebSocket endpoint in `utils/websocket/handler.go`
3. Protect player control endpoints in `main.go`

**Estimated Time:** 30-45 minutes  
**Why First:** Closes critical security vulnerability

### Phase 2: Utility Endpoints (Priority: HIGH)

4. Add password change endpoint
5. Add token refresh endpoint
6. Add user info endpoint
7. Add revoke all tokens endpoint

**Estimated Time:** 45-60 minutes  
**Why Second:** Improves usability and token management

### Phase 3: Testing & Documentation (Priority: MEDIUM)

8. Update test scripts with protected endpoint tests
9. Create comprehensive documentation

**Estimated Time:** 30-45 minutes  
**Why Last:** Validates implementation and helps users

---

## 🎨 Design Decisions

### Token Format Support

**Support Both Formats:**

```go
// Format 1: Plain token (current)
Authorization: abc123xyz789...

// Format 2: Bearer token (industry standard)
Authorization: Bearer abc123xyz789...

// Format 3: Query parameter (for WebSocket compatibility)
ws://localhost:8765/ws?token=abc123xyz789...
```

**Implementation:**

```go
func ExtractToken(r *http.Request) string {
    // Try Authorization header first
    authHeader := r.Header.Get("Authorization")
    if authHeader != "" {
        // Support "Bearer <token>" format
        if strings.HasPrefix(authHeader, "Bearer ") {
            return strings.TrimPrefix(authHeader, "Bearer ")
        }
        // Support plain token format
        return authHeader
    }

    // Fallback to query parameter (for WebSocket)
    return r.URL.Query().Get("token")
}
```

### Endpoint Protection Strategy

**Recommendation: Protect ALL player endpoints**

**Reasoning:**

- Local server = sensitive information (what you're listening to)
- Consistent security model (no confusion about which endpoints need auth)
- Easy to relax later if needed (harder to add protection after launch)
- Better audit trail (all access is logged with token info)

**Exception:** `/api/health` remains public (for monitoring tools)

---

## 📊 Security Impact

### Before Implementation

```
🔴 CRITICAL VULNERABILITIES:
- WebSocket: Anyone on network can connect and control media
- Player Control: No authentication required for play/pause/next/etc
- Token Creation: Works but can't be used to protect anything
- Information Leak: Anyone can see what you're playing
```

### After Implementation

```
🟢 SECURE:
- WebSocket: Token required, validated before connection
- Player Control: All control commands require valid token
- Token Management: Full lifecycle with refresh/revoke
- Access Control: All sensitive endpoints protected
- Audit Trail: All authenticated requests logged with token ID
```

---

## 🧪 Testing Checklist

### Manual Testing

- [ ] Register new user
- [ ] Login and receive auth token
- [ ] Create service token with auth token
- [ ] Access protected endpoint with service token (should succeed)
- [ ] Access protected endpoint without token (should fail with 401)
- [ ] Access protected endpoint with invalid token (should fail with 401)
- [ ] Access protected endpoint with expired token (should fail with 401)
- [ ] Connect to WebSocket with valid token (should succeed)
- [ ] Connect to WebSocket without token (should fail with 401)
- [ ] Change password with valid token
- [ ] Refresh token
- [ ] Revoke token and verify it can't be used
- [ ] Revoke all tokens

### Automated Testing

- [ ] Run `tests/test_auth.sh` - all tests pass
- [ ] Run `tests/test_api.sh` with auth - all tests pass
- [ ] Run with invalid tokens - all properly rejected
- [ ] WebSocket connection tests
- [ ] Rate limiting tests (if implemented)

---

## 🚀 Success Criteria

### Authentication is Complete When:

- [x] User registration works
- [x] User login returns auth token
- [x] Token creation requires authentication
- [ ] **WebSocket requires valid token**
- [ ] **All player control endpoints require valid token**
- [ ] **Middleware provides reusable auth checking**
- [ ] Password change endpoint works
- [ ] Token refresh endpoint works
- [ ] User info endpoint works
- [ ] All tests pass
- [ ] Documentation is complete

---

## 📝 Code Quality Standards

### Middleware Requirements

- ✅ Proper error handling with descriptive messages
- ✅ Logging for security events (auth failures, successful auth)
- ✅ Clean separation of concerns
- ✅ Minimal performance overhead
- ✅ Easy to use (one-line protection: `auth.RequireAuth(handler)`)

### Security Requirements

- ✅ Token validation on every protected request
- ✅ No token leakage in logs (only show first 10-20 chars)
- ✅ Proper HTTP status codes (401 for unauthorized, 403 for forbidden)
- ✅ Rate limiting considered (optional for v1.0)
- ✅ CORS headers for web client support

---

## 🎯 Deliverables

### Code Files

1. `utils/auth/middleware.go` - Authentication middleware
2. `utils/auth/handlers.go` - Updated with new endpoints
3. `utils/websocket/handler.go` - Protected WebSocket
4. `main.go` - Protected player endpoints
5. `tests/test_auth.sh` - Updated test suite

### Documentation Files

1. `docs/AUTH_COMPLETE.md` - Complete authentication guide
2. Update `README.md` - Add authentication section
3. Update API documentation with auth examples

---

## 📅 Timeline

**Target Completion: Today (November 16, 2025)**

- **Phase 1 (Core Protection):** 30-45 min ⏱️
- **Phase 2 (Utility Endpoints):** 45-60 min ⏱️
- **Phase 3 (Testing & Docs):** 30-45 min ⏱️

**Total Estimated Time:** 2-3 hours

---

## 🔄 Future Enhancements (Post-Completion)

### Short Term

- [ ] Rate limiting (e.g., 100 requests/min per token)
- [ ] Token scopes/permissions (e.g., read-only tokens)
- [ ] HTTPS/TLS support
- [ ] CORS configuration

### Medium Term

- [ ] Two-factor authentication (TOTP)
- [ ] Token usage analytics
- [ ] Audit logs
- [ ] IP whitelist/blacklist

### Long Term

- [ ] OAuth2 support for third-party apps
- [ ] Multi-user support (beyond single local user)
- [ ] Role-based access control (RBAC)
- [ ] API key management for integrations

---

## 🎉 Completion Verification

Run this command to verify all components are working:

```bash
# Full authentication system test
./tests/test_auth.sh --full

# Should output:
# ✅ User registration works
# ✅ User login works
# ✅ Token creation works
# ✅ Token validation works
# ✅ Protected endpoints require auth
# ✅ WebSocket requires auth
# ✅ Password change works
# ✅ Token refresh works
# ✅ Token revocation works
# 🎉 ALL TESTS PASSED - Authentication System Complete!
```

---

**Status:** Ready for implementation  
**Priority:** CRITICAL (security vulnerability)  
**Complexity:** Medium  
**Risk:** Low (well-defined scope)
