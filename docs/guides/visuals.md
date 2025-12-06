# 🎯 Single-User Local Auth Token System - Visual Guide

## 🏗️ Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│              QUAZAAR LOCAL SERVER                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌────────────────────┐         ┌─────────────────────┐   │
│  │   HTTP Endpoints   │         │   WebSocket         │      │
│  ├────────────────────┤         ├─────────────────────┤   │
│  │ /api/register      │         │ /ws (need token)    │   │
│  │ /api/login         │         │                     │   │
│  │ /api/tokens/create │         │ Connect with:       │   │
│  │ /api/tokens/list   │         │ ?token=YOUR_TOKEN   │   │
│  │ /api/tokens/revoke │         │                     │   │
│  └─────────┬──────────┘         └────────┬────────────┘   │
│            │                             │                 │
│            └──────────────┬──────────────┘                 │
│                          │                                 │
│                   ┌──────▼──────┐                          │
│                   │ Auth Logic  │                          │
│                   ├─────────────┤                          │
│                   │ Check token │                          │
│                   │ Hash pwd    │                          │
│                   │ Gen token   │                          │
│                   └──────┬──────┘                          │
│                          │                                 │
│                   ┌──────▼──────────────┐                  │
│                   │   SQLite Database   │                  │
│                   ├─────────────────────┤                  │
│                   │  ┌────────────────┐ │                  │
│                   │  │ user (1 row)   │ │                  │
│                   │  ├────────────────┤ │                  │
│                   │  │ id             │ │                  │
│                   │  │ username       │ │                  │
│                   │  │ password_hash  │ │                  │
│                   │  └────────────────┘ │                  │
│                   │                     │                  │
│                   │  ┌────────────────┐ │                  │
│                   │  │ tokens (many)  │ │                  │
│                   │  ├────────────────┤ │                  │
│                   │  │ id             │ │                  │
│                   │  │ name           │ │                  │
│                   │  │ token          │ │                  │
│                   │  │ service        │ │                  │
│                   │  │ expires_at     │ │                  │
│                   │  │ active         │ │                  │
│                   │  └────────────────┘ │                  │
│                   │                     │                  │
│                   │ ~/.quazaar/         │                  │
│                   │ quazaar.db          │                  │
│                   └─────────────────────┘                  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 📊 Data Flow Diagram

### 1. First Time Setup

```
User
  │
  ├─ POST /api/register
  │  │ {"username":"admin", "password":"pass"}
  │  │
  │  └─> auth.RegisterUser()
  │       │ Hash password with bcrypt
  │       │ INSERT into user table (id=1 only)
  │       │
  │       └─> ✅ User created
  │
  └─ Store credentials securely
```

### 2. Login Flow

```
User
  │
  ├─ POST /api/login
  │  │ {"username":"admin", "password":"pass"}
  │  │
  │  └─> auth.LoginUser()
  │       │ Find user by username
  │       │ Compare password hash
  │       │
  │       └─> ✅ Login successful
  │
  └─ Use credentials for further requests
```

### 3. Create Service Token

```
User (with valid login)
  │
  ├─ POST /api/tokens/create
  │  │ {"name":"My App", "service":"mobile", "duration_hours":720}
  │  │ Header: Authorization: <login_token>
  │  │
  │  ├─> Check Authorization header ✅
  │  │
  │  └─> auth.CreateToken()
  │       │ Generate random 256-bit token (base64)
  │       │ Calculate expiry: now + 720 hours
  │       │ INSERT into tokens table
  │       │
  │       └─> ✅ Token created: abc123xyz789...
  │
  └─ Save token for service
```

### 4. Use Token in WebSocket

```
Service/Device (with token)
  │
  ├─ ws://localhost:8765/ws?token=abc123xyz789...
  │  │
  │  ├─> websocket.Handle()
  │  │   │ Extract token from query params
  │  │   │
  │  │   └─> auth.ValidateToken(token)
  │  │       │ SELECT from tokens WHERE token = ?
  │  │       │ Check: active = TRUE
  │  │       │ Check: expires_at > now
  │  │       │ Check: token not revoked
  │  │       │
  │  │       └─> ✅ Valid token
  │  │
  │  └─> ✅ WebSocket connection established
  │
  └─ Send/receive messages
```

### 5. Revoke Token

```
User (with valid login)
  │
  ├─ POST /api/tokens/revoke
  │  │ {"token":"abc123xyz789..."}
  │  │ Header: Authorization: <login_token>
  │  │
  │  ├─> Check Authorization header ✅
  │  │
  │  └─> auth.RevokeToken(token)
  │       │ UPDATE tokens SET active = FALSE WHERE token = ?
  │       │
  │       └─> ✅ Token revoked
  │
  └─ Token can't be used anymore
```

---

## 🔐 Security Flow

```
Password Entry
      │
      ▼
Bcrypt Hash ─────> (never reversible)
      │
      ├─ Store in database
      │
      └─ Compare with incoming password
         │ if match ─> ✅ Valid
         │ if no match ─> ❌ Invalid


Token Generation
      │
      ├─ Generate 32 random bytes
      │
      ├─ Base64 encode
      │
      ├─ Store in database: tokens table
      │
      └─ Return to client (only once!)
         │
         └─ Client stores securely
```

---

## 📈 User <-> Tokens Relationship

```
User (1)
  ├─ id: 1
  ├─ username: "admin"
  ├─ password_hash: "$2a$10$..."
  │
  └─ owns many Tokens (*)
      │
      ├─ Token 1
      │  ├─ name: "Mobile App"
      │  ├─ service: "mobile"
      │  ├─ expires_at: 2025-12-16
      │  └─ active: true
      │
      ├─ Token 2
      │  ├─ name: "Web Client"
      │  ├─ service: "websocket"
      │  ├─ expires_at: null (never)
      │  └─ active: true
      │
      ├─ Token 3
      │  ├─ name: "Old Device"
      │  ├─ service: "mobile"
      │  ├─ expires_at: 2025-11-15
      │  └─ active: false (revoked)
      │
      └─ Token N...
```

---

## 🎯 Token Lifecycle

```
CREATE
  │
  ├─ name: Set by user
  ├─ token: Generated randomly
  ├─ service: Set by user
  ├─ expires_at: Calculated
  ├─ created_at: Now
  ├─ last_used: Null
  └─ active: True
  │
  ▼
ACTIVE & USABLE
  │
  ├─ Can authenticate requests
  ├─ last_used: Updated on each use
  │
  ├─ If expires_at passed ──> EXPIRED (can't use)
  │ If revoked ──> REVOKED (can't use)
  │
  ▼
EXPIRED or REVOKED
  │
  └─ Can be cleaned up by maintenance job
  │  (or kept for audit trail)
  │
  ▼
DELETE (optional)
  │
  └─ Removed from database
```

---

## 🔄 Request Lifecycle

### Valid Token Request

```
Client Request
  │ Authorization: token123
  │
  ▼
Server Handler
  │ Extract token
  │
  ▼
Database Query
  │ SELECT * FROM tokens WHERE token = 'token123'
  │
  ▼
Validation Checks
  │ ✅ Token exists
  │ ✅ active = true
  │ ✅ not expired
  │
  ▼
Update last_used ──────> (async)
  │
  ▼
✅ Process Request
  │ WebSocket connection established
  │ or API call proceeds
```

### Invalid Token Request

```
Client Request
  │ Authorization: invalid_token
  │
  ▼
Server Handler
  │ Extract token
  │
  ▼
Database Query
  │ SELECT * FROM tokens WHERE token = 'invalid_token'
  │ (returns: no rows)
  │
  ▼
Validation Fails
  │ ❌ Token not found
  │
  ▼
❌ Return 401 Unauthorized
  │
  ▼
Connection Rejected
  │ Client must get new token
```

---

## 📞 API Endpoints Map

```
┌──────────────────────────────────────────┐
│        Authentication Endpoints          │
├──────────────────────────────────────────┤
│                                          │
│ POST /api/register                       │
│  └─ First time setup (one user only)     │
│                                          │
│ POST /api/login                          │
│  └─ Get auth token for API access        │
│                                          │
├──────────────────────────────────────────┤
│         Service Token Endpoints          │
│      (require Authorization header)      │
├──────────────────────────────────────────┤
│                                          │
│ POST /api/tokens/create                  │
│  └─ Create token for a service           │
│     (e.g., mobile app, web client)       │
│                                          │
│ GET /api/tokens/list                     │
│  └─ List all created tokens              │
│     (useful for management)               │
│                                          │
│ POST /api/tokens/revoke                  │
│  └─ Disable a token                      │
│     (can't be used anymore)               │
│                                          │
└──────────────────────────────────────────┘
        │
        │ (authenticated requests use these tokens)
        │
        ▼
┌──────────────────────────────────────────┐
│          Service Endpoints               │
│   (require Service Token in query/auth)  │
├──────────────────────────────────────────┤
│                                          │
│ WS /ws?token=YOUR_TOKEN                  │
│  └─ WebSocket real-time connection       │
│     (control devices, get status)        │
│                                          │
└──────────────────────────────────────────┘
```

---

## 🎭 Typical User Journey

```
Day 1: Setup
  │
  ├─ 1. Register user
  │     curl POST /api/register → ✅ admin user created
  │
  ├─ 2. Login
  │     curl POST /api/login → ✅ get auth token
  │
  ├─ 3. Create tokens for services
  │     curl POST /api/tokens/create → ✅ Mobile App token
  │     curl POST /api/tokens/create → ✅ Web Client token
  │
  └─ 4. List tokens
        curl GET /api/tokens/list → ✅ see all tokens

Day 2+: Usage
  │
  ├─ Mobile App connects
  │  ws://localhost:8765/ws?token=MOBILE_TOKEN
  │  → ✅ Real-time control
  │
  ├─ Web Client connects
  │  ws://localhost:8765/ws?token=WEB_TOKEN
  │  → ✅ Real-time control
  │
  ├─ Add new device? Create new token
  │  curl POST /api/tokens/create → ✅ New Device token
  │
  ├─ Stop using old device? Revoke token
  │  curl POST /api/tokens/revoke → ✅ Token disabled
  │
  └─ Clean up expired tokens
     (automatic daily)
```

---

_This guide helps you visualize how the entire single-user token system works together!_
