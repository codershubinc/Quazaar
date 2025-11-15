# 📊 Testing Setup - Visual Overview

## Problem → Solution

```
BEFORE (❌ Broken)
═══════════════════

User (Web UI)
    ↓
Register → ✅ Works
    ↓
Login → ✅ Returns user info (NO token!)
    ↓
Try Create Token
    ↓
❌ No token to send
    ↓
401 Unauthorized


AFTER (✅ Fixed)
════════════════

User (Web UI)
    ↓
Register → ✅ Works
    ↓
Login → ✅ Returns user info + TOKEN! ✅
    ↓
Create Token (using token from login)
    ↓
✅ Token received!
    ↓
Create Service Token → ✅ Success!
    ↓
List Tokens → ✅ Works
    ↓
Connect WebSocket → ✅ Authenticated!
```

---

## API Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    QUAZAAR AUTH FLOW                        │
└─────────────────────────────────────────────────────────────┘

Step 1: REGISTER (One Time Only)
┌──────────────────────────────────────────┐
│ POST /api/register                       │
│ {username: "testuser", password: "..."}  │
└──────────────────────────────────────────┘
             ↓
        ✅ Success
        User created in database
             ↓

Step 2: LOGIN (Get Auth Token)
┌──────────────────────────────────────────┐
│ POST /api/login                          │
│ {username: "testuser", password: "..."}  │
└──────────────────────────────────────────┘
             ↓
        ✅ Response:
        {
          "token": "a1b2c3d4e5f6...",  ← USE THIS!
          "user_id": 1,
          "username": "testuser"
        }
             ↓
    Copy token for next steps
             ↓

Step 3: CREATE SERVICE TOKEN (Use Auth Token)
┌────────────────────────────────────────────────────┐
│ POST /api/tokens/create                           │
│ Headers: Authorization: a1b2c3d4e5f6...           │
│ Body: {                                            │
│   "name": "Mobile App",                           │
│   "service": "mobile",                            │
│   "duration_hours": 720                           │
│ }                                                  │
└────────────────────────────────────────────────────┘
             ↓
        ✅ Response:
        {
          "token": "x9y8z7a6b5c4...",  ← USE FOR WebSocket!
          "name": "Mobile App",
          "expires_at": "2025-12-17...",
          "active": true
        }
             ↓

Step 4: CONNECT WEBSOCKET (Use Service Token)
┌─────────────────────────────────────────────────────┐
│ ws://server:8765/ws?token=x9y8z7a6b5c4...          │
└─────────────────────────────────────────────────────┘
             ↓
        ✅ WebSocket Connected!
        Real-time communication established
             ↓
```

---

## Token Types

```
┌─────────────────────────────────────────────┐
│           TOKEN TYPES                       │
├─────────────────────────────────────────────┤
│                                             │
│ 1. AUTH TOKEN (from login)                 │
│    ├─ Purpose: Authenticate API requests   │
│    ├─ Duration: 7 days                     │
│    ├─ Usage: Authorization header          │
│    └─ Scope: Create/manage service tokens  │
│                                             │
│ 2. SERVICE TOKEN (created via API)         │
│    ├─ Purpose: Connect to WebSocket/API    │
│    ├─ Duration: User defined (0=never)     │
│    ├─ Usage: Query parameter or header     │
│    └─ Scope: Single service/device         │
│                                             │
└─────────────────────────────────────────────┘
```

---

## Testing Methods

```
┌────────────────────────────────────────────────────┐
│            HOW TO TEST                             │
├────────────────────────────────────────────────────┤
│                                                    │
│ METHOD 1: WEB UI (Easiest)                        │
│ ├─ URL: http://192.168.1.109:8765/auth.html      │
│ ├─ Register → Login → Create Token → Test         │
│ └─ No coding needed!                              │
│                                                    │
│ METHOD 2: Test Script (Fast)                      │
│ ├─ Command: bash test_api.sh                      │
│ ├─ Automatically tests all endpoints              │
│ └─ Shows tokens to use                            │
│                                                    │
│ METHOD 3: cURL (Manual)                           │
│ ├─ Copy-paste commands from API_TESTING_GUIDE.md │
│ ├─ Full control and transparency                  │
│ └─ Good for debugging                             │
│                                                    │
│ METHOD 4: Browser DevTools (Debugging)            │
│ ├─ F12 → Network tab                             │
│ ├─ See headers, request body, response           │
│ └─ Debug exact issues                             │
│                                                    │
└────────────────────────────────────────────────────┘
```

---

## File Structure

```
~/Github/Quazaar/
├── main.go                          ← Server setup + routes
├── go.mod
├── utils/
│   ├── db/
│   │   └── db.go                   ← Database initialization
│   ├── auth/
│   │   ├── auth.go                 ← Core auth logic
│   │   └── handlers.go             ← ✅ FIXED: Token in login!
│   └── websocket/
│       └── handler.go               ← WebSocket server
├── temp/web/
│   ├── index.html                  ← WebSocket test UI
│   └── auth.html                   ← ✅ FIXED: Token extraction!
├── docs/
│   ├── AUTH_SYSTEM.md              ← API reference
│   └── ... (other docs)
├── test_api.sh                      ← ✅ NEW: Test script
├── API_TESTING_GUIDE.md             ← ✅ NEW: Testing guide
├── TESTING_SETUP.md                 ← ✅ NEW: Setup guide
├── FIX_SUMMARY.md                   ← ✅ NEW: What was fixed
├── start.sh                         ← ✅ NEW: Quick start
└── .quazaar/
    └── quazaar.db                   ← SQLite database (auto-created)
```

---

## Quick Start Sequence

```
┌─────────────────────────────────────────────────┐
│          QUICK START (3 minutes)               │
├─────────────────────────────────────────────────┤
│                                                 │
│ 1. BUILD (30 sec)                             │
│    cd ~/Github/Quazaar                        │
│    go build -o quazaar                        │
│                                                 │
│ 2. RUN (2 sec)                                │
│    ./quazaar                                  │
│    (Server starts on 192.168.1.109:8765)      │
│                                                 │
│ 3. TEST (2 min 28 sec)                        │
│    Option A: Open http://.../auth.html        │
│    Option B: bash test_api.sh                 │
│    Option C: curl examples from guide         │
│                                                 │
│ 4. RESULT ✅                                   │
│    See tokens and WebSocket working!          │
│                                                 │
└─────────────────────────────────────────────────┘
```

---

## HTTP Status Codes

```
┌──────┬─────────────────────────────────┐
│ Code │ Meaning                         │
├──────┼─────────────────────────────────┤
│ 200  │ ✅ Success (GET/LIST)           │
│ 201  │ ✅ Created (POST new resource)  │
│ 400  │ ❌ Bad request (invalid JSON)   │
│ 401  │ ❌ Unauthorized (no/bad token)  │
│ 409  │ ❌ Conflict (user exists)       │
│ 500  │ ❌ Server error                 │
└──────┴─────────────────────────────────┘

Before fix: Got 401 (no token sent)
After fix:  Got 201 (token sent successfully)
```

---

## Server Startup Checklist

```
┌─────────────────────────────────────────────────┐
│         SERVER STARTUP CHECKLIST               │
├─────────────────────────────────────────────────┤
│                                                 │
│ ✅ Database initialized at ~/.quazaar/...     │
│ ✅ Users table ready                           │
│ ✅ Tokens table ready                          │
│ ✅ Indexes created                             │
│ ✅ Server running at http://192.168.1.109:... │
│ ✅ API endpoints registered                    │
│ ✅ WebSocket handler ready                     │
│ ✅ Token cleanup scheduled (hourly)            │
│ ✅ Media poller started                        │
│                                                 │
│ Server is HEALTHY and ready for testing!      │
│                                                 │
└─────────────────────────────────────────────────┘
```

---

## Documentation Files

```
📚 DOCUMENTATION
├─ API_TESTING_GUIDE.md
│  └─ Complete API reference with cURL examples
├─ TESTING_SETUP.md
│  └─ Setup, troubleshooting, and debugging
├─ FIX_SUMMARY.md
│  └─ What was broken and how it was fixed
├─ docs/AUTH_SYSTEM.md
│  └─ Complete technical documentation
└─ docs/VISUAL_GUIDE.md
   └─ Architecture diagrams and data flows
```

---

## Success Indicators

```
IF YOU SEE THESE, IT'S WORKING ✅

Server Logs:
  ✅ Database connected at /home/.../.quazaar/quazaar.db
  ✅ Users table ready
  ✅ Tokens table ready
  ✅ Server running at http://192.168.1.109:8765
  ✅ User registered via API: testuser
  ✅ User logged in via API: testuser (auth token created)
  ✅ Token created via API: Mobile App

Web UI Response:
  ✅ "User registered successfully"
  ✅ "Login successful" + token displayed
  ✅ "Token created" + token displayed
  ✅ Token list shows "Active" badges

API Response:
  ✅ POST /api/register → 201 Created
  ✅ POST /api/login → 200 OK + token
  ✅ POST /api/tokens/create → 201 Created
  ✅ GET /api/tokens/list → 200 OK + token list
```

---

## Now You're Ready! 🚀

Everything is set up for testing. Pick your preferred method:

1. **Web UI** → Most user-friendly
2. **Test Script** → Fastest full test
3. **cURL** → Full control and transparency
4. **Browser DevTools** → Debug issues

All should now work without 401 errors! ✅
