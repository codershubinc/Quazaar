# 🧪 Testing Setup - Complete Guide

## Issue Fixed

**Problem:** `401 Unauthorized` when calling `/api/tokens/create`

**Root Cause:**

- Login endpoint wasn't returning the auth token
- Client couldn't get the token to use for subsequent API calls

**Solution:**

- Modified `HandleLogin` to create and return an auth token
- Updated web UI to properly extract and store the token
- Added comprehensive logging for debugging
- Created test scripts for manual testing

---

## ✅ What's Now Available for Testing

### 1. **Web UI**

- **URL:** `http://192.168.1.109:8765/auth.html`
- **Features:**
  - Register user
  - Login (now returns token!)
  - Create service tokens
  - List all tokens
  - Revoke tokens
  - Copy/paste helpers

### 2. **Automated Test Script**

```bash
cd ~/Github/Quazaar
bash test_api.sh
```

This script:

- Registers a test user
- Logs in to get auth token
- Creates a service token
- Lists all tokens
- Shows you the tokens to use

### 3. **Manual cURL Testing**

See `API_TESTING_GUIDE.md` for all cURL examples

### 4. **Health Check Endpoint** (New!)

```bash
curl http://192.168.1.109:8765/api/health
```

Response: `{"status":"healthy","service":"quazaar","version":"1.0"}`

### 5. **Server Logs**

Watch the terminal where server runs for debug info:

```
✅ User logged in via API: testuser (auth token created)
🔍 Create token attempt - Auth header: a1b2c3d4e5...
✅ Token created via API: Mobile App
```

---

## Testing Flow

### Step 1: Register User (Once Only)

```bash
curl -X POST http://192.168.1.109:8765/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

✅ Response:

```json
{
  "success": true,
  "message": "User registered successfully",
  "user_id": 1,
  "username": "testuser"
}
```

### Step 2: Login (Get Auth Token)

```bash
curl -X POST http://192.168.1.109:8765/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

✅ Response:

```json
{
  "success": true,
  "message": "Login successful",
  "user_id": 1,
  "username": "testuser",
  "token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z...",
  "auth_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z..."
}
```

**👉 Copy the `token` value - this is your auth token!**

### Step 3: Create Service Token (Use Auth Token)

```bash
AUTH_TOKEN="a1b2c3d4e5f6..." # From login response

curl -X POST http://192.168.1.109:8765/api/tokens/create \
  -H "Content-Type: application/json" \
  -H "Authorization: $AUTH_TOKEN" \
  -d '{
    "name":"Mobile App",
    "service":"mobile",
    "duration_hours":720
  }'
```

✅ Response:

```json
{
  "id": 1,
  "name": "Mobile App",
  "token": "x9y8z7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1q0r9s8t7u6v5w4...",
  "service": "mobile",
  "expires_at": "2025-12-17T15:30:45Z",
  "created_at": "2025-11-16T15:30:45Z",
  "active": true
}
```

**👉 This token can be used for WebSocket connections!**

### Step 4: List Tokens (Use Auth Token)

```bash
curl -X GET http://192.168.1.109:8765/api/tokens/list \
  -H "Authorization: $AUTH_TOKEN"
```

### Step 5: Use Service Token for WebSocket

```
ws://192.168.1.109:8765/ws?token=x9y8z7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1q0r9s8t7u6v5w4...
```

---

## Files Modified/Created

### Modified Files

- ✏️ `utils/auth/handlers.go` - Updated HandleLogin to return token + debug logging
- ✏️ `main.go` - Added /api/health endpoint and /auth route handler
- ✏️ `temp/web/auth.html` - Better token extraction and debug logging

### New Files

- 🆕 `test_api.sh` - Automated API test script
- 🆕 `API_TESTING_GUIDE.md` - Complete testing reference
- 🆕 `TESTING_SETUP.md` - This file

---

## Quick Reference

| What                       | How                                               |
| -------------------------- | ------------------------------------------------- |
| Check if server is running | `curl http://192.168.1.109:8765/api/health`       |
| Test via web UI            | Open `http://192.168.1.109:8765/auth.html`        |
| Test with script           | `bash test_api.sh`                                |
| View server logs           | Watch terminal where you ran `./quazaar`          |
| WebSocket test UI          | Open `http://192.168.1.109:8765/`                 |
| Token format               | Base64-encoded 256-bit random value (long string) |

---

## Debugging

### If you still get 401 Unauthorized:

1. **Check server is running:**

   ```bash
   curl http://192.168.1.109:8765/api/health
   ```

   Should return: `{"status":"healthy",...}`

2. **Check login works:**

   ```bash
   curl -X POST http://192.168.1.109:8765/api/login \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","password":"password123"}'
   ```

   Should include `"token":"..."` in response

3. **Check token is being passed:**
   Look at browser DevTools → Network tab

   - Should see `Authorization: [token]` header
   - Token should be long string (not empty)

4. **Check server logs:**
   In terminal running server, you should see:

   ```
   🔍 Create token attempt - Auth header: a1b2c3d4e5...
   ```

5. **If token validation fails:**
   Server logs will show:
   ```
   ❌ Token validation failed: [error]
   ```
   - Try logging in again to get fresh token
   - Check token hasn't expired

---

## Success Indicators

✅ Working system shows:

- Health check returns `{"status":"healthy"}`
- Registration succeeds (or says user already exists)
- Login returns auth token
- Can create service tokens
- Tokens appear in token list
- Can connect to WebSocket with service token

---

## Next Steps

1. **Run test script:**

   ```bash
   bash test_api.sh
   ```

2. **Or use web UI:**

   - Open `http://192.168.1.109:8765/auth.html`
   - Register → Login → Create Token → Test WebSocket

3. **Or test manually:**
   - Follow the curl examples in `API_TESTING_GUIDE.md`

---

## Files to Reference

- `API_TESTING_GUIDE.md` - All API examples and responses
- `test_api.sh` - Automated testing
- `docs/AUTH_SYSTEM.md` - Complete API documentation
- `main.go` - Server setup and routes

Good luck! 🚀
