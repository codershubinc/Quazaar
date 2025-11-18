# Spotify Integration Guide

## Overview

Quazaar integrates with Spotify's Web API using OAuth 2.0 Authorization Code Flow to enable secure access to user's Spotify data and playback control.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     Spotify Integration Flow                     │
└─────────────────────────────────────────────────────────────────┘

1. User Authentication (OAuth 2.0)
   ┌──────────┐         ┌──────────┐         ┌───────────┐
   │  Client  │────────▶│ Quazaar  │────────▶│  Spotify  │
   │          │         │  Server  │         │   OAuth   │
   └──────────┘         └──────────┘         └───────────┘
        │                     │                     │
        │  /spotify/login     │                     │
        │────────────────────▶│                     │
        │                     │  Redirect to Auth   │
        │                     │────────────────────▶│
        │                     │                     │
        │                 User authorizes app       │
        │◀─────────────────────────────────────────│
        │                     │                     │
        │  /spotify/callback  │                     │
        │────────────────────▶│                     │
        │                     │  Exchange code      │
        │                     │────────────────────▶│
        │                     │  Return tokens      │
        │                     │◀────────────────────│
        │                     │                     │
        │  Success response   │                     │
        │◀────────────────────│                     │

2. Token Management
   ┌──────────────────────────────────────────────────────┐
   │  tokens/store.go - Token Storage & Retrieval         │
   │    - GetSpotifyRefreshToken()                        │
   │    - GetSpotifyAccessToken()                         │
   │    - SetSpotifyRefreshToken()                        │
   │    - In-memory cache + SQLite persistence            │
   └──────────────────────────────────────────────────────┘
               │
               ▼
   ┌──────────────────────────────────────────────────────┐
   │  tokens/exchange.go - Token Exchange & Refresh       │
   │    - ExchangeCodeForToken()                          │
   │    - RefreshAccessToken()                            │
   │    - ValidateToken()                                 │
   │    - RevokeToken()                                   │
   └──────────────────────────────────────────────────────┘

3. API Integration
   ┌──────────────────────────────────────────────────────┐
   │  auth/user.go - User Profile Operations              │
   │    - Login() - Initiates OAuth flow                  │
   │    - Callback() - Handles OAuth callback             │
   │    - GetUserProfile() - Fetches user data            │
   └──────────────────────────────────────────────────────┘
               │
               ▼
   ┌──────────────────────────────────────────────────────┐
   │  devices/manager.go - Device Management              │
   │    - getUserDevices() - Fetch available devices      │
   └──────────────────────────────────────────────────────┘
```

## File Structure

```
internal/spotify/
├── config.go              # API base URL configuration
├── spotify.go             # Initialization & token check
├── auth/
│   ├── api.go            # HTTP handlers (GetUser)
│   └── user.go           # OAuth flow (Login, Callback, GetUserProfile)
├── devices/
│   ├── api.go            # HTTP handlers (GetUserDevices)
│   └── manager.go        # Device fetching logic
└── tokens/
    ├── exchange.go       # Token operations (exchange, refresh, validate)
    └── store.go          # Token storage & caching
```

## Environment Variables

Required configuration in `.env` file:

```bash
# Spotify API Credentials (from Spotify Developer Dashboard)
SPOTIFY_CLIENT_ID=your_client_id_here
SPOTIFY_CLIENT_SECRET=your_client_secret_here

# OAuth Redirect URI (must match Spotify app settings)
SPOTIFY_REDIRECT_URI=http://127.0.0.1:8765/api/v0.1/spotify/callback

# Spotify API Base URL (optional, defaults to https://api.spotify.com/v1)
SPOTIFY_API_BASE_URL=https://api.spotify.com/v1
```

## API Endpoints

### 1. OAuth Login
**Endpoint:** `GET /api/v0.1/spotify/login`  
**Authentication:** Required (Quazaar token)  
**Description:** Initiates Spotify OAuth flow, redirects user to Spotify authorization page

**Response:**
- `307 Temporary Redirect` to Spotify authorization page

**Flow:**
1. Generates random state parameter for CSRF protection
2. Constructs authorization URL with:
   - `client_id`: Your Spotify app client ID
   - `response_type`: "code"
   - `redirect_uri`: Callback URL
   - `scope`: Requested permissions
   - `state`: Random string for security
3. Redirects user to Spotify for authorization

**Scopes Requested:**
- `user-read-private` - Read user's private info
- `user-read-email` - Read user's email
- `user-read-currently-playing` - Read current playback
- `user-read-playback-state` - Read playback state
- `user-read-recently-played` - Read recently played
- `user-follow-read` - Read following
- `user-library-read` - Read saved library
- `user-modify-playback-state` - Control playback

---

### 2. OAuth Callback
**Endpoint:** `GET /api/v0.1/spotify/callback`  
**Authentication:** Public (no Quazaar token required)  
**Description:** Handles OAuth callback from Spotify, exchanges code for tokens

**Query Parameters:**
- `code` (required): Authorization code from Spotify
- `state` (required): State parameter for CSRF validation

**Response:**
```json
{
  "success": true,
  "message": "Spotify authentication successful",
  "expires_in": 3600
}
```

**Error Responses:**
- `400 Bad Request` - Missing code or state
- `500 Internal Server Error` - Token exchange failed

**Flow:**
1. Validates state parameter (CSRF protection)
2. Exchanges authorization code for tokens via Spotify API
3. Stores refresh token in database
4. Returns success response with token expiry

---

### 3. Get User Profile
**Endpoint:** `GET /api/v0.1/spotify/user`  
**Authentication:** Required (Quazaar token)  
**Description:** Fetches authenticated user's Spotify profile

**Response:**
```json
{
  "id": "user_spotify_id",
  "display_name": "John Doe",
  "email": "john@example.com",
  "country": "US",
  "product": "premium",
  "followers": {
    "total": 123
  },
  "images": [
    {
      "url": "https://...",
      "height": 300,
      "width": 300
    }
  ]
}
```

**Error Responses:**
- `500 Internal Server Error` - Failed to get access token or user profile

---

### 4. Get User Devices
**Endpoint:** `GET /api/v0.1/spotify/devices`  
**Authentication:** Required (Quazaar token)  
**Description:** Lists user's available Spotify playback devices

**Response:**
```json
{
  "devices": [
    {
      "id": "device_id",
      "is_active": true,
      "is_private_session": false,
      "is_restricted": false,
      "name": "My Computer",
      "type": "Computer",
      "volume_percent": 50
    }
  ]
}
```

**Error Responses:**
- `405 Method Not Allowed` - Only GET is supported
- `500 Internal Server Error` - Failed to retrieve devices

---

## OAuth 2.0 Flow Details

### Step-by-Step Process

#### 1. **User Initiates Login**
```
Client → GET /api/v0.1/spotify/login
```
- User clicks "Login with Spotify" button
- Quazaar generates random state parameter
- User is redirected to Spotify authorization page

#### 2. **User Authorizes Application**
```
User → Authorizes on Spotify → Redirected back to Quazaar
```
- User sees Spotify permission request
- User accepts or denies
- Spotify redirects to callback URL with authorization code

#### 3. **Token Exchange**
```
Quazaar → POST https://accounts.spotify.com/api/token
```
Request Body:
```
grant_type=authorization_code
code=<authorization_code>
redirect_uri=<callback_url>
```

Headers:
```
Authorization: Basic <base64(client_id:client_secret)>
Content-Type: application/x-www-form-urlencoded
```

Response:
```json
{
  "access_token": "NgCXRK...MzYjw",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "NgAagA...Um_SHo",
  "scope": "user-read-private user-read-email..."
}
```

#### 4. **Token Storage**
- Refresh token stored in SQLite database
- Access token cached in memory with expiry time
- Tokens associated with user account

#### 5. **Automatic Token Refresh**
When access token expires (after 1 hour):
```
Quazaar → POST https://accounts.spotify.com/api/token
```
Request Body:
```
grant_type=refresh_token
refresh_token=<stored_refresh_token>
```

Response:
```json
{
  "access_token": "NgCXRK...NewToken",
  "token_type": "Bearer",
  "expires_in": 3600
}
```
Note: New refresh token may or may not be included

---

## Token Management

### Caching Strategy

**In-Memory Cache:**
- `spotifyClientRefreshToken` - Long-lived refresh token
- `spotifyClientAccessToken` - Short-lived access token (1 hour)
- `spotifyAccessTokenExpiry` - Timestamp for token expiration

**Database Storage:**
- Refresh tokens persisted in SQLite
- Stored with token type "spotify"
- Associated with token name "spotifyClientRefreshToken"

### Token Lifecycle

```
┌─────────────────────────────────────────────────────┐
│              Token Lifecycle                        │
└─────────────────────────────────────────────────────┘

1. Initial Authorization
   ┌──────────┐
   │ No Token │
   └────┬─────┘
        │
        ▼
   ┌──────────────┐
   │ OAuth Login  │
   └────┬─────────┘
        │
        ▼
   ┌──────────────────┐
   │ Get Auth Code    │
   └────┬─────────────┘
        │
        ▼
   ┌──────────────────────┐
   │ Exchange for Tokens  │
   └────┬─────────────────┘
        │
        ▼
   ┌─────────────────────────┐
   │ Store Refresh Token     │
   │ Cache Access Token      │
   └────┬────────────────────┘

2. Regular Use
   ┌──────────────────┐
   │ API Call Needed  │
   └────┬─────────────┘
        │
        ▼
   ┌─────────────────────┐
   │ Check Cache Valid?  │
   └────┬────────┬───────┘
        │ Yes    │ No
        ▼        ▼
   ┌────────┐  ┌──────────────────┐
   │ Return │  │ Refresh Token    │
   │ Cached │  │ Update Cache     │
   └────────┘  └────┬─────────────┘
                    ▼
               ┌────────────┐
               │ Return New │
               └────────────┘

3. Token Expiry
   Access Token: 1 hour (auto-refresh)
   Refresh Token: Long-lived (months/years)
```

---

## Error Handling

### Common Errors

**1. Missing Environment Variables**
```
Error: SPOTIFY_CLIENT_ID not set
Solution: Add credentials to .env file
```

**2. Invalid Redirect URI**
```
Error: redirect_uri mismatch
Solution: Ensure SPOTIFY_REDIRECT_URI matches Spotify app settings
```

**3. Token Expired**
```
Error: Token validation failed with status: 401
Solution: Automatic refresh will be triggered
```

**4. Invalid Authorization Code**
```
Error: Token exchange failed: status 400
Reason: Code already used or expired
Solution: User must re-authorize
```

**5. Refresh Token Revoked**
```
Error: Token refresh failed: status 400
Solution: User must re-authenticate via /spotify/login
```

---

## Security Considerations

### 1. **CSRF Protection**
- Random state parameter generated for each login
- State validated in callback to prevent CSRF attacks

### 2. **Secure Token Storage**
- Refresh tokens encrypted in database
- Access tokens cached in memory only
- Tokens associated with authenticated users

### 3. **HTTPS Required**
- OAuth flow requires HTTPS in production
- Redirect URIs must use HTTPS (except localhost)

### 4. **Scope Limitation**
- Request only necessary scopes
- Users can see and revoke permissions anytime

### 5. **Token Rotation**
- Access tokens expire after 1 hour
- Automatic refresh using refresh token
- Refresh tokens can be rotated by Spotify

---

## Testing Guide

### 1. **Setup Spotify App**
1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create new app
3. Add redirect URI: `http://127.0.0.1:8765/api/v0.1/spotify/callback`
4. Copy Client ID and Client Secret

### 2. **Configure Quazaar**
```bash
# Create .env file
cat > .env << EOF
SPOTIFY_CLIENT_ID=your_client_id
SPOTIFY_CLIENT_SECRET=your_client_secret
SPOTIFY_REDIRECT_URI=http://127.0.0.1:8765/api/v0.1/spotify/callback
EOF
```

### 3. **Test OAuth Flow**
```bash
# 1. Start Quazaar server
go build -o quazaar ./cmd/server
./quazaar

# 2. Login to Quazaar (get auth token)
curl -X POST http://127.0.0.1:8765/api/v0.1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"your_user","password":"your_pass"}'

# Save the token from response
TOKEN="your_quazaar_token"

# 3. Initiate Spotify login (will redirect in browser)
curl -L http://127.0.0.1:8765/api/v0.1/spotify/login \
  -H "Authorization: Bearer $TOKEN"

# 4. After authorization, test user endpoint
curl http://127.0.0.1:8765/api/v0.1/spotify/user \
  -H "Authorization: Bearer $TOKEN"

# 5. Test devices endpoint
curl http://127.0.0.1:8765/api/v0.1/spotify/devices \
  -H "Authorization: Bearer $TOKEN"
```

---

## Troubleshooting

### Issue: "Refresh token not found in DB"
**Cause:** User hasn't completed OAuth flow  
**Solution:** Navigate to `/api/v0.1/spotify/login` to authenticate

### Issue: "Failed to get Spotify access token"
**Cause:** Refresh token expired or revoked  
**Solution:** Re-authenticate via login endpoint

### Issue: "Token exchange failed: status 400"
**Cause:** Invalid client credentials or redirect URI mismatch  
**Solution:** 
1. Verify SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET
2. Check redirect URI matches Spotify app settings
3. Ensure authorization code hasn't been used

### Issue: "state mismatch"
**Cause:** CSRF protection triggered  
**Solution:** This is normal security - user should retry login

---

## Future Enhancements

- [ ] Support for multiple users (multi-tenant)
- [ ] Automatic token refresh before expiry
- [ ] Playback control endpoints (play, pause, skip)
- [ ] Playlist management
- [ ] Search functionality
- [ ] Currently playing widget
- [ ] Webhook support for playback events
- [ ] Token encryption at rest
- [ ] Rate limiting for API calls
- [ ] Caching for user profile and devices

---

## References

- [Spotify Web API Documentation](https://developer.spotify.com/documentation/web-api)
- [OAuth 2.0 Authorization Code Flow](https://developer.spotify.com/documentation/general/guides/authorization-guide/#authorization-code-flow)
- [Spotify API Endpoints](https://developer.spotify.com/documentation/web-api/reference/)
- [Spotify Scopes](https://developer.spotify.com/documentation/general/guides/authorization-guide/#list-of-scopes)

---

## Support

For issues or questions:
- Check logs: Quazaar logs all Spotify API interactions
- Verify environment variables are set correctly
- Ensure Spotify app is configured properly
- Check network connectivity to Spotify API

**Created:** November 17, 2025  
**Version:** 0.0.1.3 (Beta)  
**Author:** Quazaar Development Team
