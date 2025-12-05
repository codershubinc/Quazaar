# Authentication System

Quazaar uses a single-user token-based authentication system. This ensures that only the owner can control the server, while allowing multiple devices and services to connect securely.

## Overview

- **Single User**: Only one user account is allowed (the owner).
- **Multiple Tokens**: You can generate multiple tokens for different devices (e.g., "Mobile App", "Web Dashboard").
- **Revocable**: Tokens can be revoked individually if a device is lost or compromised.
- **Expiration**: Tokens can have an expiration time or be permanent.

## Database Schema

The authentication system uses two main tables in the SQLite database:

### `user` Table

Stores the single user account.

| Column          | Type     | Description            |
| :-------------- | :------- | :--------------------- |
| `id`            | INTEGER  | Primary Key (always 1) |
| `username`      | TEXT     | The username           |
| `password_hash` | TEXT     | Bcrypt hashed password |
| `created_at`    | DATETIME | Creation timestamp     |

### `tokens` Table

Stores access tokens for devices and services.

| Column       | Type     | Description                           |
| :----------- | :------- | :------------------------------------ |
| `id`         | INTEGER  | Primary Key                           |
| `name`       | TEXT     | Device/Service name (e.g., "Pixel 7") |
| `token`      | TEXT     | The token string                      |
| `service`    | TEXT     | Service type (e.g., "mobile", "web")  |
| `expires_at` | DATETIME | Expiration time (NULL = never)        |
| `created_at` | DATETIME | Creation timestamp                    |
| `last_used`  | DATETIME | Last usage timestamp                  |
| `active`     | BOOLEAN  | TRUE if valid, FALSE if revoked       |

## API Endpoints

### 1. Register User

Register the owner account. This can only be done once.

- **URL**: `/api/register`
- **Method**: `POST`
- **Body**:

```json
{
  "username": "admin",
  "password": "secure_password"
}
```

### 2. Login

Login to get the initial session or token.

- **URL**: `/api/login`
- **Method**: `POST`
- **Body**:

```json
{
  "username": "admin",
  "password": "secure_password"
}
```

### 3. Create Token

Generate a new token for a specific device or service.

- **URL**: `/api/tokens/create`
- **Method**: `POST`
- **Headers**: `Authorization: <your_token>`
- **Body**:

```json
{
  "name": "Living Room Tablet",
  "service": "mobile",
  "duration_hours": 720
}
```

> **Note**: Set `duration_hours` to `0` for a non-expiring token.

### 4. List Tokens

View all active and inactive tokens.

- **URL**: `/api/tokens/list`
- **Method**: `GET`
- **Headers**: `Authorization: <your_token>`

### 5. Revoke Token

Invalidate a specific token.

- **URL**: `/api/tokens/revoke`
- **Method**: `POST`
- **Headers**: `Authorization: <your_token>`
- **Body**:

```json
{
  "token": "token_string_to_revoke"
}
```

## WebSocket Authentication

To connect to the WebSocket endpoint, you can pass the token in the query string or headers.

### Query String

```text
ws://localhost:8765/ws?token=YOUR_TOKEN_HERE
```

### Header (if supported by client)

```text
Authorization: YOUR_TOKEN_HERE
```

## Security Best Practices

1. **Strong Password**: Use a strong password for the main account.
2. **Device-Specific Tokens**: Create a separate token for each device.
3. **Revoke Unused Tokens**: Periodically check the token list and revoke unused ones.
4. **HTTPS**: Always run Quazaar behind a reverse proxy with HTTPS in production to protect tokens in transit.
