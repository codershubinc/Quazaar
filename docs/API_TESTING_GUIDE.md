# 🧪 API Testing Guide

## Quick Start - For Testing

### 1. Run the Test Script

```bash
cd ~/Github/Quazaar
bash test_api.sh
```

This will automatically:

- Register a test user
- Login to get auth token
- Create a service token
- List all tokens

### 2. Manual Testing with cURL

#### Register User

```bash
curl -X POST http://192.168.1.109:8765/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

#### Login

```bash
curl -X POST http://192.168.1.109:8765/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

Response includes `token` - copy this for next steps.

#### Create Service Token

```bash
# Replace TOKEN_HERE with actual token from login
curl -X POST http://192.168.1.109:8765/api/tokens/create \
  -H "Content-Type: application/json" \
  -H "Authorization: TOKEN_HERE" \
  -d '{
    "name":"Mobile App",
    "service":"mobile",
    "duration_hours":720
  }'
```

#### List Tokens

```bash
curl -X GET http://192.168.1.109:8765/api/tokens/list \
  -H "Authorization: TOKEN_HERE"
```

#### Revoke Token

```bash
curl -X POST http://192.168.1.109:8765/api/tokens/revoke \
  -H "Content-Type: application/json" \
  -H "Authorization: TOKEN_HERE" \
  -d '{"token":"SERVICE_TOKEN_HERE"}'
```

### 3. Web UI Testing

Open `http://192.168.1.109:8765/auth.html` and:

1. Register a user
2. Login
3. Create service tokens
4. List and manage tokens

### 4. WebSocket Testing

Once you have a service token:

```
ws://192.168.1.109:8765/ws?token=YOUR_SERVICE_TOKEN
```

Open `http://192.168.1.109:8765/index.html` and use the WebSocket test interface.

---

## Troubleshooting

### 401 Unauthorized Error

**Problem:** Getting `401 Unauthorized` when creating tokens

**Solution:**

1. Make sure you logged in first with correct credentials
2. Copy the `token` from login response
3. Paste it in the `Authorization` header (exactly, with no modifications)
4. Check browser console (Developer Tools → Console) for debug info

### Token Not Working

**Problem:** Token says it's invalid

**Solution:**

1. Make sure token hasn't expired
2. Check that token is active (not revoked)
3. Get a fresh token by logging in again

### Registration Fails

**Problem:** User already exists error

**Solution:**

- This is expected! Only one user can register (single-user design)
- Just login with existing credentials instead

---

## Response Examples

### Login Response

```json
{
  "success": true,
  "message": "Login successful",
  "user_id": 1,
  "username": "testuser",
  "token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9...",
  "auth_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9..."
}
```

### Create Token Response

```json
{
  "id": 1,
  "name": "Mobile App",
  "token": "x9y8z7a6b5c4d3e2f1g0h9i8j7k6l5m4n3o2p1...",
  "service": "mobile",
  "expires_at": "2025-12-17T15:30:45Z",
  "created_at": "2025-11-16T15:30:45Z",
  "active": true
}
```

### List Tokens Response

```json
{
  "success": true,
  "tokens": [
    {
      "id": 1,
      "name": "Mobile App",
      "token": "...",
      "service": "mobile",
      "expires_at": "2025-12-17T15:30:45Z",
      "created_at": "2025-11-16T15:30:45Z",
      "last_used": "2025-11-16T15:35:00Z",
      "active": true
    }
  ],
  "count": 1
}
```

---

## API Endpoints Summary

| Method | Endpoint             | Auth        | Purpose                    |
| ------ | -------------------- | ----------- | -------------------------- |
| POST   | `/api/register`      | No          | Create user account (once) |
| POST   | `/api/login`         | No          | Authenticate and get token |
| POST   | `/api/tokens/create` | Yes         | Create service token       |
| GET    | `/api/tokens/list`   | Yes         | List all tokens            |
| POST   | `/api/tokens/revoke` | Yes         | Revoke a token             |
| WS     | `/ws`                | Token param | WebSocket connection       |

---

## Debug Logging

The server logs all API requests. Watch the terminal where the server runs:

```
✅ User registered via API: testuser
✅ User logged in via API: testuser (auth token created)
🔍 Create token attempt - Auth header: a1b2c3d4e5...
✅ Token created via API: Mobile App
✅ Listed 1 tokens via API
```

If you get errors, check the server logs for detailed messages.
