# Authentication API

## `/api/v0.1/signup` [POST]

**Description:** Registers a new user with the system.

**Required Body (JSON):**

```json
{
  "username": "user123", // Min 3 chars
  "password": "pass123" // Min 6 chars
}
```

**Response Data (JSON):**

```json
{
  "success": true,
  "message": "User registered successfully",
  "username": "user123"
}
```

**Example cURL:**

```bash
curl -X POST http://localhost:8766/api/v0.1/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user123",
    "password": "pass123"
  }'
```

## `/api/v0.1/login` [POST]

**Description:** Authenticates a user and returns a session token.

**Required Body (JSON):**

```json
{
  "username": "user123",
  "password": "pass123"
}
```

**Response Data (JSON):**

```json
{
  "success": true,
  "message": "Login successful",
  "username": "user123",
  "token": "a1b2c3d4e5f6...",
  "tokenType": "deviceId"
}
```

**Example cURL:**

```bash
curl -X POST http://localhost:8766/api/v0.1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user123",
    "password": "pass123"
  }'
```

## `/api/v0.1/auth/refresh` [POST]

**Description:** Refreshes an existing authentication token.

**Required Body (JSON):**

```json
{
  "token": "old_token_string"
}
```

**Response Data (JSON):**

```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "token": "new_token_string",
  "tokenType": "deviceId"
}
```

## `/api/v0.1/auth/change-password` [POST]

**Description:** Changes the password for the authenticated user.

**Required Body (JSON):**

```json
{
  "token": "current_valid_token",
  "old_password": "old_pass123",
  "new_password": "new_pass456" // Min 6 chars
}
```

**Response Data (JSON):**

```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

## `/api/v0.1/auth/user` [GET]

**Description:** Retrieves information about the currently authenticated user.

**Required Headers:**

- `Authorization`: `Bearer <token>`
- OR Query Param: `?token=<token>`

**Response Data (JSON):**

```json
{
  "success": true,
  "user": {
    "id": 1,
    "username": "user123",
    "name": "User Name",
    "created_at": "2023-10-27T10:00:00Z"
  }
}
```

## `/api/v0.1/auth/logout` [POST]

**Description:** Logs out the user by invalidating the provided token.

**Required Body (JSON):**

```json
{
  "token": "token_to_invalidate"
}
```

**Response Data (JSON):**

```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

## `/api/v0.1/auth/tokens` [GET]

**Description:** Lists all active tokens for the user.

**Required Headers:**

- `Authorization`: `Bearer <token>`
- OR Query Param: `?token=<token>`

**Response Data (JSON):**

```json
{
  "success": true,
  "tokens": [
    {
      "id": 1,
      "name": "Web Client",
      "token": "abc...",
      "service": "websocket",
      "created_at": "2023-10-27T10:00:00Z",
      "active": true
    }
  ]
}
```
