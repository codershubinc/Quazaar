# 🔐 Single-User Token System - Quick Reference

## 📊 Database Schema

### `user` Table

```
id              INTEGER PRIMARY KEY (only 1 allowed)
username        TEXT (your username)
password_hash   TEXT (hashed password)
created_at      DATETIME
```

### `tokens` Table

```
id              INTEGER PRIMARY KEY
name            TEXT (e.g., "Mobile App", "Web Client")
token           TEXT (the actual token string)
service         TEXT (e.g., "websocket", "api", "mobile")
expires_at      DATETIME (NULL = never expires)
created_at      DATETIME
last_used       DATETIME
active          BOOLEAN (TRUE = can use, FALSE = revoked)
```

---

## 🚀 Usage Examples

### 1️⃣ Register User (First Time Only)

**Method**: POST
**URL**: `http://localhost:8765/api/register`
**Content-Type**: `application/json`

**Request**:

```json
{
  "username": "admin",
  "password": "secure_password_123"
}
```

**Response** (201 Created):

```json
{
  "success": true,
  "message": "User registered successfully",
  "user_id": 1,
  "username": "admin"
}
```

---

### 2️⃣ Login User

**Method**: POST
**URL**: `http://localhost:8765/api/login`
**Content-Type**: `application/json`

**Request**:

```json
{
  "username": "admin",
  "password": "secure_password_123"
}
```

**Response** (200 OK):

```json
{
  "success": true,
  "message": "Login successful",
  "user_id": 1,
  "username": "admin"
}
```

---

### 3️⃣ Create Token for a Service

Create tokens for different devices/services. Each token is independent.

**Method**: POST
**URL**: `http://localhost:8765/api/tokens/create`
**Content-Type**: `application/json`
**Authorization**: Your login token (or any active token)

**Request**:

```json
{
  "name": "My Mobile App",
  "service": "mobile",
  "duration_hours": 720
}
```

**Parameters**:

- `name`: Human-readable name (e.g., "iPhone", "Web App", "Python Script")
- `service`: Service identifier (e.g., "websocket", "api", "mobile")
- `duration_hours`: How long token lasts (0 = never expires)

**Response** (201 Created):

```json
{
  "id": 1,
  "name": "My Mobile App",
  "token": "abc123xyz789_base64_encoded_token...",
  "service": "mobile",
  "expires_at": "2025-12-16T10:30:00Z",
  "created_at": "2025-11-16T10:30:00Z",
  "last_used": null,
  "active": true
}
```

**Use this token**:

```bash
# Connect to WebSocket with token
ws://localhost:8765/ws?token=abc123xyz789_base64_encoded_token

# Or in header
Authorization: abc123xyz789_base64_encoded_token
```

---

### 4️⃣ List All Tokens

**Method**: GET
**URL**: `http://localhost:8765/api/tokens/list`
**Authorization**: Your token

**Response** (200 OK):

```json
{
  "success": true,
  "tokens": [
    {
      "id": 1,
      "name": "My Mobile App",
      "token": "abc123...",
      "service": "mobile",
      "expires_at": "2025-12-16T10:30:00Z",
      "created_at": "2025-11-16T10:30:00Z",
      "last_used": "2025-11-16T11:00:00Z",
      "active": true
    },
    {
      "id": 2,
      "name": "Web Client",
      "token": "xyz789...",
      "service": "websocket",
      "expires_at": null,
      "created_at": "2025-11-16T12:00:00Z",
      "last_used": null,
      "active": true
    }
  ],
  "count": 2
}
```

---

### 5️⃣ Revoke a Token

Disable a token (can't be used anymore, but stays in database).

**Method**: POST
**URL**: `http://localhost:8765/api/tokens/revoke`
**Content-Type**: `application/json`
**Authorization**: Your token

**Request**:

```json
{
  "token": "abc123xyz789_the_token_to_revoke"
}
```

**Response** (200 OK):

```json
{
  "success": true,
  "message": "Token revoked successfully"
}
```

---

## 💻 Testing with cURL

### Register:

```bash
curl -X POST http://localhost:8765/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

### Login:

```bash
curl -X POST http://localhost:8765/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

### Create Token:

```bash
TOKEN="your_auth_token_here"
curl -X POST http://localhost:8765/api/tokens/create \
  -H "Content-Type: application/json" \
  -H "Authorization: $TOKEN" \
  -d '{"name":"My App","service":"websocket","duration_hours":720}'
```

### List Tokens:

```bash
TOKEN="your_auth_token_here"
curl -X GET http://localhost:8765/api/tokens/list \
  -H "Authorization: $TOKEN"
```

### Revoke Token:

```bash
TOKEN="your_auth_token_here"
REVOKE_TOKEN="token_to_revoke"
curl -X POST http://localhost:8765/api/tokens/revoke \
  -H "Content-Type: application/json" \
  -H "Authorization: $TOKEN" \
  -d "{\"token\":\"$REVOKE_TOKEN\"}"
```

---

## 🔌 WebSocket with Token

**Connect to WebSocket**:

### Method 1: Token in URL

```javascript
const token = "your_token_here";
const ws = new WebSocket(`ws://localhost:8765/ws?token=${token}`);

ws.onopen = () => console.log("✅ Connected!");
ws.onmessage = (event) => console.log("📨 Message:", event.data);
ws.onerror = (error) => console.error("❌ Error:", error);
```

### Method 2: Token in Header (for compatible clients)

```javascript
const token = "your_token_here";
const ws = new WebSocket("ws://localhost:8765/ws");

// Send token first
ws.onopen = () => {
  ws.send(JSON.stringify({ type: "auth", token: token }));
};
```

---

## 📝 Code Examples

### Programmatically Create Token (in Go):

```go
package main

import (
	"Quazaar/utils/auth"
	"log"
	"time"
)

func main() {
	// Create token valid for 7 days
	duration := 7 * 24 * time.Hour
	token, err := auth.CreateToken("My Device", "mobile", &duration)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Token: %s", token.Token)
	log.Printf("Expires: %s", token.ExpiresAt)
}
```

### Validate Token (in Go):

```go
valid, err := auth.ValidateToken("your_token_string")
if valid && err == nil {
	log.Println("✅ Token is valid!")
} else {
	log.Printf("❌ Token invalid: %v", err)
}
```

### List All Tokens (in Go):

```go
tokens, err := auth.GetAllTokens()
if err != nil {
	log.Fatal(err)
}

for _, t := range tokens {
	log.Printf("Token: %s | Service: %s | Active: %v", t.Name, t.Service, t.Active)
}
```

---

## 🧹 Automatic Cleanup

Expired tokens are automatically cleaned up:

- Every 1 hour automatically
- Via `auth.CleanExpiredTokens()` in code
- In database via `DELETE FROM tokens WHERE expires_at < NOW()`

---

## 🔍 View Database

**View all tokens**:

```bash
sqlite3 ~/.quazaar/quazaar.db "SELECT id, name, service, active, expires_at FROM tokens;"
```

**View user**:

```bash
sqlite3 ~/.quazaar/quazaar.db "SELECT id, username FROM user;"
```

**View active tokens**:

```bash
sqlite3 ~/.quazaar/quazaar.db "SELECT id, name, service FROM tokens WHERE active = TRUE;"
```

---

## 🎯 Flow Diagram

```
1. User Registration
   └─> /api/register → User stored in database

2. User Login
   └─> /api/login → Validates credentials

3. Create Service Token
   └─> /api/tokens/create → New token created
   └─> Token stored in database

4. Use Token in WebSocket
   └─> ws://localhost:8765/ws?token=YOUR_TOKEN
   └─> Server validates token
   └─> Connection established

5. Revoke Token
   └─> /api/tokens/revoke → Token marked inactive
   └─> Can't use anymore
```

---

## ⚠️ Important Notes

1. **Single User Only**: Only 1 user can be registered (due to `CHECK (id = 1)` constraint)
2. **Multiple Tokens**: One user can have many tokens for different services
3. **Token Expiry**: Optional - set `duration_hours: 0` for tokens that never expire
4. **Token Revocation**: Revoking a token doesn't delete it, just marks it inactive
5. **Security**: Store tokens securely, treat them like passwords
6. **Database Location**: `~/.quazaar/quazaar.db` (in your home directory)

---

## 🐛 Troubleshooting

**"user already registered"**

- User already exists, use login instead

**"invalid or expired token"**

- Token doesn't exist, expired, or revoked
- Try creating a new token

**"Token has expired"**

- Create new token with longer duration

**"unauthorized"**

- Missing Authorization header
- Invalid token

---

_Last Updated: November 16, 2025_
