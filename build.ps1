$ErrorActionPreference = "Stop"

Write-Host "🚀 Quazaar Build Script" -ForegroundColor Cyan
Write-Host "======================="

# Path definitions
$ProjectRoot = $PSScriptRoot
$CsprojPath = Join-Path $ProjectRoot "cmd\windows-helper\QuazaarMedia.csproj"
$SidecarBuildDir = Join-Path $ProjectRoot "temp\sidecar_build"
$SidecarDest = Join-Path $ProjectRoot "internal\sidecar\QuazaarMedia.exe"
$GoMain = Join-Path $ProjectRoot "cmd\server\main.go"
$OutputExe = Join-Path $ProjectRoot "quazaar.exe"

# 1. Clean previous builds
if (Test-Path $SidecarBuildDir) { Remove-Item -Recurse -Force $SidecarBuildDir | Out-Null }

# 2. Build .NET Sidecar
Write-Host "`n[1/3] Compiling Windows Helper (.NET)..." -ForegroundColor Yellow
# Note: --self-contained false requires .NET Runtime on target machine. Set to true for standalone.
dotnet publish $CsprojPath -c Release -r win-x64 -p:PublishSingleFile=true --self-contained false -o $SidecarBuildDir
if ($LASTEXITCODE -ne 0) { throw "dotnet publish failed" }

# 3. Copy Sidecar
Write-Host "`n[2/3] Embedding Sidecar..." -ForegroundColor Yellow
Copy-Item -Force (Join-Path $SidecarBuildDir "QuazaarMedia.exe") $SidecarDest

# 3.5 Build Tray App
Write-Host "`n[2.5/3] Compiling Tray App..." -ForegroundColor Yellow
$TrayCsprojPath = Join-Path $ProjectRoot "cmd\windows-tray\QuazaarTray.csproj"
$TrayBuildDir = Join-Path $ProjectRoot "temp\tray_build"
$TrayDest = Join-Path $ProjectRoot "QuazaarTray.exe"
if (Test-Path $TrayBuildDir) { Remove-Item -Recurse -Force $TrayBuildDir | Out-Null }
dotnet publish $TrayCsprojPath -c Release -r win-x64 -p:PublishSingleFile=true --self-contained false -o $TrayBuildDir
if ($LASTEXITCODE -ne 0) { throw "dotnet publish tray app failed" }
Copy-Item -Force (Join-Path $TrayBuildDir "QuazaarTray.exe") $TrayDest

# 4. Build Go Binary
Write-Host "`n[3/3] Building Quazaar (Go)..." -ForegroundColor Yellow
go build -o $OutputExe $GoMain
if ($LASTEXITCODE -ne 0) { throw "go build failed" }

Write-Host "`n✅ Build Successful: $OutputExe" -ForegroundColor Green
