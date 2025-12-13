$ErrorActionPreference = "Stop"

Write-Host "🚀 Quazaar Release Build Script" -ForegroundColor Cyan
Write-Host "==============================="

# Path definitions
$ProjectRoot = $PSScriptRoot
$ReleaseDir = Join-Path $ProjectRoot "release\windows"

# Source Paths
$HelperCsproj = Join-Path $ProjectRoot "cmd\windows-helper\QuazaarMedia.csproj"
$TrayCsproj = Join-Path $ProjectRoot "cmd\windows-tray\QuazaarTray.csproj"
$GoMain = Join-Path $ProjectRoot "cmd\server\main.go"

# Temp Build Paths
$SidecarBuildDir = Join-Path $ProjectRoot "temp\sidecar_build"
$TrayBuildDir = Join-Path $ProjectRoot "temp\tray_build"

# Destination Paths (for embedding/final)
$SidecarDest = Join-Path $ProjectRoot "internal\sidecar\QuazaarMedia.exe"
$FinalExe = Join-Path $ReleaseDir "quazaar.exe"
$FinalTray = Join-Path $ReleaseDir "QuazaarTray.exe"

# Ensure Release Directory Exists
if (!(Test-Path $ReleaseDir)) { New-Item -ItemType Directory -Force -Path $ReleaseDir | Out-Null }

# 1. Build Windows Helper (Sidecar)
Write-Host "`n[1/4] Compiling Windows Helper (.NET)..." -ForegroundColor Yellow
if (Test-Path $SidecarBuildDir) { Remove-Item -Recurse -Force $SidecarBuildDir | Out-Null }
dotnet publish $HelperCsproj -c Release -r win-x64 -p:PublishSingleFile=true --self-contained false -o $SidecarBuildDir
if ($LASTEXITCODE -ne 0) { throw "dotnet publish helper failed" }

# 2. Embed Sidecar
Write-Host "`n[2/4] Embedding Sidecar..." -ForegroundColor Yellow
Copy-Item -Force (Join-Path $SidecarBuildDir "QuazaarMedia.exe") $SidecarDest

# 3. Build Tray Application
Write-Host "`n[3/4] Compiling Tray App..." -ForegroundColor Yellow
if (Test-Path $TrayBuildDir) { Remove-Item -Recurse -Force $TrayBuildDir | Out-Null }
dotnet publish $TrayCsproj -c Release -r win-x64 -p:PublishSingleFile=true --self-contained false -o $TrayBuildDir
if ($LASTEXITCODE -ne 0) { throw "dotnet publish tray app failed" }

# 4. Build Go Binary
Write-Host "`n[4/4] Building Quazaar (Go)..." -ForegroundColor Yellow
go build -o $FinalExe $GoMain
if ($LASTEXITCODE -ne 0) { throw "go build failed" }

# 5. Copy Artifacts to Release Folder
Write-Host "`n[Finalizing] Copying files to release folder..." -ForegroundColor Yellow
Copy-Item -Force (Join-Path $TrayBuildDir "QuazaarTray.exe") $FinalTray

# Copy install script if it exists in root (optional)
# if (Test-Path "install.ps1") { Copy-Item -Force "install.ps1" $ReleaseDir }

Write-Host "`n[SUCCESS] Build Successful!" -ForegroundColor Green
Write-Host "Output Directory: $ReleaseDir"
Write-Host "   - quazaar.exe"
Write-Host "   - QuazaarTray.exe"
