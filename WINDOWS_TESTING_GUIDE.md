🧪 WINDOWS MEDIA INFO TESTING GUIDE
=====================================

✅ What we have:
1. Cross-compiled Windows executable: quazaar-windows.exe
2. Mock test program: test_windows_media.go
3. Windows-specific media info code: utils/mediaInfoWindows.go

🖥️  TESTING OPTIONS:

Option 1 - Mock Test (Current Linux Environment)
-----------------------------------------------
go run test_windows_media.go
# Simulates Windows behavior without real Windows

Option 2 - Real Windows Testing
-------------------------------
1. Copy quazaar-windows.exe to Windows machine
2. Run: .\quazaar-windows.exe
3. Test endpoints:
   curl http://localhost:8765/api/v0.1/player/info
   curl http://localhost:8765/api/v0.1/player/list

Option 3 - Wine Testing (Experimental)
--------------------------------------
wine quazaar-windows.exe
# May work for basic functionality

📋 WINDOWS-SPECIFIC FEATURES:
- Uses Windows Media Foundation APIs
- Supports Groove Music, Windows Media Player
- Compatible with UWP media apps
- Direct Windows API integration

🔧 BUILD VERIFICATION:
$(file quazaar-windows.exe)

🎯 READY TO TEST ON WINDOWS!