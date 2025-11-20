# Creating User on Quazaar Server

This guide shows how to create and manage users on your Quazaar server using curl commands.

## Prerequisites

- Quazaar server running (default: `http://localhost:8765`)
- curl installed on your system

## User Registration

### Create a New User

Quazaar supports **single-user mode** - only one user can be registered per server instance.

```bash
curl -X POST http://localhost:8765/api/v0.1/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your_secure_password"
  }'
```

**Response (Success):**

```json
{
  "success": true,
  "message": "User registered successfully"
}
```

**Response (User Already Exists):**

```json
{
  "error": "User already exists"
}
```

### Important Notes

- Only **one user** can be registered per server
- Password is hashed using bcrypt (cost: 10)
- Username must be unique
- Both username and password are required

## User Login

### Authenticate and Get Token

```bash
curl -X POST http://localhost:8765/api/v0.1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your_secure_password"
  }'
```

**Response (Success):**

```json
{
  "success": true,
  "message": "Login successful",
  "token": "$2a$10$UPqFIdvGUihieuLPwGYyreU5L8WjM/LuUQv36b8L0ZZM3Vh4ybJWq",
  "tokenType": "deviceId",
  "username": "quazaar_admin"
}
```

**Response (Invalid Credentials):**

```json
{
  "error": "Invalid username or password"
}
```
