$appName = "Quazaar"
$installDir = "$env:LOCALAPPDATA\$appName"

# Paths for shortcuts
$desktopDir = [Environment]::GetFolderPath("Desktop")
$startMenuDir = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs"

$desktopShortcut = Join-Path $desktopDir "$appName.lnk"
$startMenuShortcut = Join-Path $startMenuDir "$appName.lnk"

Write-Host "Installing $appName..."

# 0. Stop existing processes
$processes = Get-Process -Name "QuazaarTray", "quazaar" -ErrorAction SilentlyContinue
if ($processes) {
    Write-Host "Stopping running Quazaar processes..."
    $processes | Stop-Process -Force
    Start-Sleep -Seconds 1
}

# 1. Create Install Directory
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
    Write-Host "Created directory: $installDir"
}

# 2. Check Source Files
$filesToCopy = @("quazaar.exe", "QuazaarTray.exe")
foreach ($file in $filesToCopy) {
    if (Test-Path $file) {
        Copy-Item $file -Destination $installDir -Force
        Write-Host "Copied $file"
    } else {
        Write-Error "❌ Missing source file: $file. Please build the project first."
        exit 1
    }
}

# Optional Icon
if (Test-Path "icon.ico") {
    $iconSource = Convert-Path "icon.ico"
    $iconDest = Join-Path $installDir "icon.ico"
    
    # Check if it's a PNG (Magic bytes: 89 50 4E 47)
    $bytes = Get-Content $iconSource -Encoding Byte -TotalCount 4
    if ($bytes[0] -eq 0x89 -and $bytes[1] -eq 0x50 -and $bytes[2] -eq 0x4E -and $bytes[3] -eq 0x47) {
        Write-Host "⚠️  Detected PNG file renamed as ICO. Converting to real ICO..."
        
        try {
            Add-Type -AssemblyName System.Drawing
            $pngBytes = [System.IO.File]::ReadAllBytes($iconSource)
            $ms = New-Object System.IO.MemoryStream
            $bw = New-Object System.IO.BinaryWriter($ms)

            # Write ICO Header (0, 1, 1)
            $bw.Write([Int16]0); $bw.Write([Int16]1); $bw.Write([Int16]1)

            # Write Icon Directory Entry
            $img = [System.Drawing.Image]::FromStream((New-Object System.IO.MemoryStream(,$pngBytes)))
            $w = $img.Width; $h = $img.Height
            if ($w -ge 256) { $w = 0 }; if ($h -ge 256) { $h = 0 }
            
            $bw.Write([byte]$w); $bw.Write([byte]$h); $bw.Write([byte]0); $bw.Write([byte]0)
            $bw.Write([Int16]0); $bw.Write([Int16]32)
            $bw.Write([int]$pngBytes.Length)
            $bw.Write([int]22) # Offset (6 header + 16 entry)
            
            # Write PNG Data
            $bw.Write($pngBytes)
            
            [System.IO.File]::WriteAllBytes($iconDest, $ms.ToArray())
            $ms.Dispose()
            Write-Host "✅ Converted and copied icon.ico"
        } catch {
            Write-Error "Failed to convert icon: $_"
            Copy-Item $iconSource -Destination $installDir -Force
        }
    } else {
        Copy-Item $iconSource -Destination $installDir -Force
        Write-Host "Copied icon.ico"
    }
}

# Function to create shortcut
function Create-Shortcut {
    param(
        [string]$path,
        [string]$target,
        [string]$arguments,
        [string]$workDir,
        [string]$iconPath
    )
    try {
        $WshShell = New-Object -ComObject WScript.Shell
        $Shortcut = $WshShell.CreateShortcut($path)
        $Shortcut.TargetPath = $target
        $Shortcut.Arguments = $arguments
        $Shortcut.WorkingDirectory = $workDir
        if ($iconPath -and (Test-Path $iconPath)) {
            $Shortcut.IconLocation = $iconPath
        } else {
            $Shortcut.IconLocation = "shell32.dll,13" # Default 'Globe' icon
        }
        $Shortcut.Save()
        Write-Host "✅ Shortcut created at: $path"
    } catch {
        Write-Error "❌ Failed to create shortcut at $path : $_"
    }
}

$installedIcon = Join-Path $installDir "icon.ico"

# 3. Create Shortcuts
# Desktop
Create-Shortcut -path $desktopShortcut -target "$installDir\QuazaarTray.exe" -arguments "" -workDir $installDir -iconPath $installedIcon

# Start Menu (Makes it searchable)
Create-Shortcut -path $startMenuShortcut -target "$installDir\QuazaarTray.exe" -arguments "" -workDir $installDir -iconPath $installedIcon

Write-Host "`n🎉 Installation Complete!"
Write-Host "👉 You can now find '$appName' in your Start Menu and Windows Search." 