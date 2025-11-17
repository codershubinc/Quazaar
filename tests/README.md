# Quazaar Test Scripts

This directory contains test scripts for testing various components of the Quazaar server.

## Available Test Scripts

### 1. Authentication Flow Test (`test_auth_complete.sh`)

Comprehensive test suite for the complete authentication system including all endpoints and edge cases.

**Features:**

- 30+ test cases covering all authentication scenarios
- Colored output for easy reading
- Tests signup, login, logout flow
- Tests all utility endpoints (refresh, change-password, user info, tokens)
- Tests protected endpoints with and without authentication
- Tests edge cases and error handling
- Automatic test result summary

**Usage:**

```bash
# Make sure the server is running first
./quazaar

# In another terminal, run the test script
cd tests
./test_auth_complete.sh
```

**Requirements:**

- Server running on `127.0.0.1:8765`
- `curl` command available
- `jq` for JSON parsing (optional, for prettier output)

**Test Coverage:**

1. **User Registration**

   - Signup with valid credentials
   - Signup with short username (validation)
   - Signup with short password (validation)

2. **User Login**

   - Login with correct credentials
   - Login with invalid credentials
   - Token generation and retrieval

3. **User Information**

   - Get user info with Bearer token
   - Get user info with query parameter token
   - Get user info without token (unauthorized)

4. **Token Management**

   - List active tokens
   - Token refresh with valid token
   - Old token invalidation after refresh
   - Token validation after logout

5. **Protected Endpoints**

   - WebSocket access with/without token
   - Player info endpoints with/without token
   - System info endpoints (WiFi, Bluetooth) with token

6. **Password Management**

   - Change password with valid old password
   - Login with new password
   - Login with old password (should fail)
   - Change password with wrong old password (validation)

7. **Logout**
   - Logout with valid token
   - Token invalidation after logout

**Output Example:**

```
================================================
  Quazaar Authentication Flow Test Suite
================================================

Test 1: Signup new user
✓ PASSED - Status: 201
{
  "success": true,
  "message": "User registered successfully",
  "username": "testuser_1700000000"
}

...

================================================
  Test Summary
================================================
Total Tests:  30
Passed:       30
Failed:       0

✓ All tests passed!
```

## Other Test Scripts

### 2. API Test (`test_api.sh`)

Basic API endpoint testing (original test script).

### 3. Auth Test (`test_auth.sh`)

Basic authentication testing (original test script).

### 4. Windows API Test (`test_windows_api.sh`)

Windows-specific API testing.

## Running All Tests

To run all test scripts:

```bash
cd tests
for script in test_*.sh; do
    echo "Running $script..."
    ./$script
    echo ""
done
```

## Writing New Tests

When writing new test scripts:

1. Follow the naming convention: `test_*.sh`
2. Make the script executable: `chmod +x test_*.sh`
3. Use colored output for better readability
4. Include test counters and summary
5. Test both success and failure cases
6. Document the test cases in this README

## Notes

- All test scripts expect the server to be running on `127.0.0.1:8765`
- Tests create temporary test users with timestamps
- Authentication tests automatically clean up by logging out
- Some tests may fail if the database is not initialized

## Troubleshooting

**Server not running:**

```
curl: (7) Failed to connect to 127.0.0.1 port 8765: Connection refused
```

Solution: Start the server with `./quazaar`

**Permission denied:**

```
bash: ./test_auth_complete.sh: Permission denied
```

Solution: Make the script executable with `chmod +x test_auth_complete.sh`

**jq command not found:**
The tests will still work, but JSON output won't be prettified.
Solution: Install jq with `sudo apt install jq` (Ubuntu/Debian) or `brew install jq` (macOS)
