# Changelog - v0.1.4

**Release Date:** November 20, 2025  
**Project Name:** Quazaar  
**Status:** 🟡 Beta - File Sharing Feature

---

## Overview

Patch release introducing secure file sharing functionality with temporary URI-based authentication and token management. This update adds the ability to securely transfer files between devices using time-limited access tokens.

---

## 🎉 New Features

### File Sharing System

- ✅ **Temporary URI Generation** - Create secure, time-limited file upload endpoints
- ✅ **Token-Based Authentication** - Device ID and token validation for secure file transfers
- ✅ **Multipart Form Upload** - Support for file uploads up to 100MB
- ✅ **Automatic Token Cleanup** - Tokens are automatically deleted after use
- ✅ **File Storage** - Uploaded files are securely stored using the helper system
- ✅ **Token Expiry** - Configurable expiration (default: 3600 seconds)

### File Sharing API Endpoints

- `GET /api/v0.1/fileshare/create-accept-uri` - Generate temporary upload URI
- `POST /api/v0.1/fileshare/acceptfile?deviceId={id}&token={token}` - Accept file upload

### Testing Tools

- ✅ **HTML Test Interface** - Interactive file upload test page
- ✅ **Real-time Progress Tracking** - Upload progress bar with percentage
- ✅ **Speed Monitoring** - Current and average upload speed (MB/s)
- ✅ **Dynamic Endpoint Detection** - Auto-detects server URL from page location
- ✅ **Query Parameter Support** - Pre-fill deviceId and token from URL

---

## 📦 New Files

### Internal Packages

```
internal/
└── fileshare/
    ├── api.go                  # File upload handlers
    └── file_url_handler.go     # Token generation and validation
```

### Testing Files

```
temp/
└── fileshare_test.html        # Interactive file sharing test interface
```

---

## 🔧 Technical Details

### Database Integration

- **File Share Tokens Table** - Stores temporary device tokens
- **Token Management** - CRUD operations for token lifecycle
  - `StoreFileShareDeviceToken(token, deviceId, expiry)`
  - `GetFileShareDeviceToken(token)`
  - `DeleteFileShareDeviceToken(token)`

### Security Features

- **Token-Based Access Control** - Each upload requires valid device ID + token pair
- **Single-Use Tokens** - Tokens are deleted after use (via defer)
- **Time-Limited Access** - Tokens expire after configured duration
- **Validation Layer** - Device ID verification before file acceptance

### File Upload Flow

1. Client requests temporary URI via `GET /api/v0.1/fileshare/create-accept-uri`
2. Server generates random device ID and token (16 chars each)
3. Server stores token in database with expiry time
4. Server returns JSON response with upload URI
5. Client uploads file to URI with device ID and token
6. Server validates credentials and accepts multipart form data
7. Server stores file and automatically deletes used token

---

## 🎨 HTML Test Interface Features

### User Interface

- Clean, modern design with responsive layout
- File selection with size display (MB)
- Pre-filled device ID and token fields
- Real-time upload progress visualization

### Progress Monitoring

- **Progress Bar** - Visual percentage indicator (0-100%)
- **Current Speed** - Instantaneous upload speed in MB/s
- **Average Speed** - Overall average upload speed
- **Time Elapsed** - Upload duration in seconds
- **Size Tracking** - Uploaded size vs total file size

### Smart Features

- Auto-populate fields from URL query parameters
- Dynamic server endpoint detection
- Success/error message display with auto-hide
- Form reset after successful upload
- Network error handling

---

## 🔄 Code Quality

### Helper Function Usage

- `helpers.GenerateRandomString(16)` - Secure token generation
- `helpers.StoreFile(filename, file)` - File storage abstraction
- `helpers.SendJsonDataToClient(w, status, data)` - Consistent JSON responses

### Error Handling

- Method validation (GET/POST enforcement)
- Multipart form parsing with memory limits
- Token validation and expiry checks
- File storage error handling
- Deferred token cleanup for reliability

---

## 🐛 Bug Fixes

- Fixed incomplete `deviceId` and `token` variable declaration in HTML
- Corrected endpoint URL construction for dynamic host detection
- Added proper error responses for invalid tokens

---

## 📝 API Documentation

### Request Temporary File Share URI

```http
GET /api/v0.1/fileshare/create-accept-uri
```

**Response:**

```json
{
  "acceptUri": "/api/v0.1/fileshare/acceptfile?deviceId={id}&token={token}",
  "message": "Temporary file share URI created successfully",
  "expiry": 3600,
  "time": "Mon, 02 Jan 2006 15:04:05 GMT"
}
```

### Upload File

```http
POST /api/v0.1/fileshare/acceptfile?deviceId={id}&token={token}
Content-Type: multipart/form-data

file: [binary data]
```

**Success Response:**

```
HTTP/1.1 200 OK
File acceptance authorized, file stored successfully
```

**Error Responses:**

- `401 Unauthorized` - Invalid or expired token
- `400 Bad Request` - Failed to parse form or missing file
- `500 Internal Server Error` - File storage failure

---

## 🧪 Testing

### Manual Testing

1. Start Quazaar server
2. Open `temp/fileshare_test.html` in browser
3. Optional: Add `?deviceId=xxx&token=yyy` to URL
4. Select file and click "Upload File"
5. Monitor real-time progress and speed
6. Verify success message and check stored file

### Command Line Testing

```bash
# Request temporary URI
curl http://localhost:8765/api/v0.1/fileshare/create-accept-uri

# Upload file
curl -X POST \
  "http://localhost:8765/api/v0.1/fileshare/acceptfile?deviceId=xxx&token=yyy" \
  -F "file=@/path/to/file.txt"
```

---

## 🔐 Security Considerations

- Tokens are single-use and auto-deleted
- Device ID must match stored value
- Time-limited access prevents stale tokens
- File size limited to 100MB per upload
- No directory traversal in filename handling

---

## 📊 Performance

- **Max Upload Size:** 100 MB (configurable)
- **Token Length:** 16 characters (random)
- **Default Expiry:** 3600 seconds (1 hour)
- **Memory Buffering:** 100 MB for multipart parsing

---

## 🎯 What's Next (v0.1.5)

### Planned Enhancements

- [ ] Multiple file upload support
- [ ] File type validation and filtering
- [ ] Download file endpoint
- [ ] File listing and management API
- [ ] Storage quota management
- [ ] File metadata (upload time, size, type)
- [ ] Thumbnail generation for images
- [ ] Direct device-to-device transfer
- [ ] File sharing history/logs
- [ ] Compression support

---

## 📚 Dependencies

- **Database:** SQLite (via `internal/db`)
- **HTTP:** Go standard library `net/http`
- **Helpers:** `pkg/helpers` (file storage, random generation, JSON responses)

---

## 🔗 Related Documentation

- [File Sharing API](/docs/FILESHARE_API.md) (if created)
- [Testing Guide](/docs/README_TESTING.md)
- [API Testing Guide](/docs/API_TESTING_GUIDE.md)

---

## ✅ Summary

v0.1.4 successfully introduces a secure, token-based file sharing system with:

- Complete file upload API with validation
- Time-limited, single-use access tokens
- Interactive HTML test interface with progress tracking
- Database integration for token management
- Proper error handling and security measures

This release provides the foundation for device-to-device file transfers and sets the stage for more advanced file management features in future versions.
