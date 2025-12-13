Set WshShell = CreateObject("WScript.Shell")
' Start Quazaar hidden (0)
WshShell.Run "quazaar.exe", 0, False
' Wait a bit for server to start
WScript.Sleep 1000
' Open Browser
WshShell.Run "http://localhost:8765/"