$ErrorActionPreference = 'Stop'

# Detect arch
$arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
if ($arch -ne "amd64") {
    Write-Error "Only 64-bit Windows is supported."
    exit 1
}

$url = "https://github.com/priyabratasahoo21/kibuild-mcp/releases/latest/download/kibuild-mcp-windows-amd64.exe"
$installDir = "$env:LOCALAPPDATA\Programs\kibuild-mcp"
$exePath = "$installDir\kibuild-mcp.exe"

# Create dir
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
}

Write-Host "Downloading kibuild-mcp..."
Invoke-RestMethod -Uri $url -OutFile $exePath

# Update PATH if not present
$userPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*kibuild-mcp*") {
    Write-Host "Adding $installDir to User PATH..."
    [System.Environment]::SetEnvironmentVariable("PATH", $userPath + ";$installDir", "User")
    Write-Host "Please restart your terminal/IDE to apply PATH changes."
}

Write-Host "Success! kibuild-mcp installed to $exePath"
