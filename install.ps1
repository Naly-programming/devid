# devid PowerShell installer
# Usage: irm https://raw.githubusercontent.com/Naly-programming/devid/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$Repo = "Naly-programming/devid"
$InstallDir = "$env:LOCALAPPDATA\Programs\devid"

# Detect architecture
$Arch = if ([System.Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Error "32-bit Windows is not supported"
    exit 1
}

Write-Host "Fetching latest devid version..."

# Get latest version
try {
    $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
    $Version = $Release.tag_name -replace '^v', ''
} catch {
    Write-Error "Failed to fetch latest release: $_"
    exit 1
}

Write-Host "Installing devid v$Version (windows/$Arch)..."

# Download
$Url = "https://github.com/$Repo/releases/download/v$Version/devid_${Version}_windows_${Arch}.zip"
$TmpDir = New-Item -ItemType Directory -Path "$env:TEMP\devid-install-$([System.Guid]::NewGuid())" -Force
$ZipPath = Join-Path $TmpDir "devid.zip"

try {
    Invoke-WebRequest -Uri $Url -OutFile $ZipPath -UseBasicParsing
} catch {
    Write-Error "Download failed: $_"
    Remove-Item -Recurse -Force $TmpDir
    exit 1
}

# Extract
Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

# Install
New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
$BinaryPath = Join-Path $TmpDir "devid.exe"
$DestPath = Join-Path $InstallDir "devid.exe"
Copy-Item -Path $BinaryPath -Destination $DestPath -Force

# Cleanup
Remove-Item -Recurse -Force $TmpDir

# Add to PATH if not already there
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    Write-Host ""
    Write-Host "Added $InstallDir to your user PATH."
    Write-Host "Restart your terminal for the change to take effect."
}

Write-Host ""
Write-Host "devid v$Version installed to $DestPath"
Write-Host ""
Write-Host "Get started:"
Write-Host "  devid init"
