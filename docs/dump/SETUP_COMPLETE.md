# ✅ Single-User Local Auth Token System - Setup Complete!

## 📦 What You Now Have

### Database Files

- ✅ `utils/db/db.go` - Database initialization & table creation
- ✅ `utils/auth/auth.go` - Core auth functions
- ✅ `utils/auth/handlers.go` - HTTP API endpoints

### Documentation

- ✅ `docs/AUTH_SYSTEM.md` - Complete usage guide & examples
- ✅ `main_example.go` - Example integration

---

## 🚀 How to Use

### Step 1: Install Dependencies

```bash
go mod tidy
```

### Step 2: Update Your main.go

Copy the structure from `main_example.go` to your actual `main.go`:

```go
import (
    "Quazaar/utils/auth"
    "Quazaar/utils/db"
    // ... other imports
)

func main() {
    // Initialize database
    if err := db.Init(); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer db.CloseDB()

    // Setup HTTP routes
    http.HandleFunc("/api/register", auth.HandleRegister)
    http.HandleFunc("/api/login", auth.HandleLogin)
    http.HandleFunc("/api/tokens/create", auth.HandleCreateToken)
    http.HandleFunc("/api/tokens/list", auth.HandleListTokens)
    http.HandleFunc("/api/tokens/revoke", auth.HandleRevokeToken)
    http.HandleFunc("/ws", websocket.Handle)

    // ... rest of your code ...
}
```

### Step 3: Build & Run

```bash
go build -o quazaar && ./quazaar
```

**Expected Output**:

```
🚀 Hello Quazaar Server ...
✅ Database connected at /home/youruser/.quazaar/quazaar.db
✅ User table ready
✅ Tokens table ready
✅ Indexes created
✅ Database tables ready
✅ WebSocket service token: abc123...
```

---

## 📝 Quick Start

### 1. Register User (First Time Only)

```bash
curl -X POST http://localhost:8765/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

### 2. Login

```bash
curl -X POST http://localhost:8765/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

### 3. Create Token

```bash
AUTH_TOKEN="your_token_from_login"
curl -X POST http://localhost:8765/api/tokens/create \
  -H "Content-Type: application/json" \
  -H "Authorization: $AUTH_TOKEN" \
  -d '{"name":"My Mobile App","service":"mobile","duration_hours":720}'
```

### 4. Use Token in WebSocket

```javascript
const token = "your_token_here";
const ws = new WebSocket(`ws://localhost:8765/ws?token=${token}`);
```

---

## 💾 Database Structure

**Single-user design** - only 1 user can exist:

- `user` table: Username & password (one row max)
- `tokens` table: Unlimited tokens for different services

**Example tokens table**:

```
| name            | service   | active | expires_at |
|-----------------|-----------|--------|-----------|
| Mobile App      | mobile    | true   | 2025-12-16 |
| Web Client      | websocket | true   | never      |
| Python Script   | api       | false  | 2025-11-20 |
```

---

## 🔑 Key Functions

### In `utils/auth/auth.go`:

| Function                               | Purpose                    |
| -------------------------------------- | -------------------------- |
| `RegisterUser(username, password)`     | Register single user       |
| `LoginUser(username, password)`        | Verify credentials         |
| `GenerateToken()`                      | Create random token string |
| `CreateToken(name, service, duration)` | Create service token       |
| `ValidateToken(token)`                 | Check if token is valid    |
| `RevokeToken(token)`                   | Disable a token            |
| `RevokeAllTokens()`                    | Disable all tokens         |
| `GetAllTokens()`                       | List all tokens            |
| `CleanExpiredTokens()`                 | Remove expired tokens      |

---

## 🗂️ Files Overview

```
Quazaar/
├── utils/
│   ├── db/
│   │   └── db.go              ✅ Database setup & tables
│   ├── auth/
│   │   ├── auth.go            ✅ Core auth logic
│   │   └── handlers.go        ✅ HTTP endpoints
│   ├── websocket/
│   │   └── handler.go         (Update to use tokens)
│   └── poller/
│       └── handler.go
├── main.go                    📝 (Update with new routes)
├── main_example.go            📋 (Reference implementation)
├── docs/
│   └── AUTH_SYSTEM.md         📖 Complete guide
└── go.mod
```

---

## 🔧 Update WebSocket to Use Tokens

In `utils/websocket/handler.go`, check the token:

```go
func Handle(w http.ResponseWriter, r *http.Request) {
    // Get token from query parameter or header
    token := r.URL.Query().Get("token")
    if token == "" {
        token = r.Header.Get("Authorization")
    }

    // Validate token
    valid, err := auth.ValidateToken(token)
    if !valid || err != nil {
        log.Printf("❌ Unauthorized: %v", err)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    log.Printf("✅ WebSocket connection authorized")

    // Continue with WebSocket upgrade...
    conn, err := CreateWebSocketConnection(w, r)
    // ... rest of code
}
```

---

## 📚 Full Documentation

See `docs/AUTH_SYSTEM.md` for:

- ✅ Complete API endpoints with examples
- ✅ cURL commands for testing
- ✅ JavaScript WebSocket examples
- ✅ Go code examples
- ✅ Database structure
- ✅ Troubleshooting

---

## ⚡ Next Steps

1. **Update main.go** with the new routes (copy from `main_example.go`)
2. **Update websocket/handler.go** to validate tokens
3. **Test the API** with the cURL examples
4. **Create tokens** for your services
5. **Use tokens** to connect to WebSocket

---

## 🎯 Summary

✅ **Single-user design**: Perfect for local server
✅ **Multiple tokens**: Different services can have separate tokens  
✅ **Token expiry**: Optional expiration (0 = never)
✅ **Token revocation**: Disable tokens without deletion
✅ **Automatic cleanup**: Expired tokens deleted hourly
✅ **Secure**: bcrypt password hashing, random tokens
✅ **Local database**: SQLite file at `~/.quazaar/quazaar.db`

**You're ready to go!** 🚀
