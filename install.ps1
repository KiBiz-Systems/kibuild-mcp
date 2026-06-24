$ErrorActionPreference = 'Stop'

# ── Step 1: Binary ──────────────────────────────────────────────────────────

$arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
if ($arch -ne "amd64") {
    Write-Error "Only 64-bit Windows is supported."
    exit 1
}

$url = "https://github.com/priyabratasahoo21/kibuild-mcp/releases/latest/download/kibuild-mcp-windows-amd64.exe"
$installDir = "$env:LOCALAPPDATA\Programs\kibuild-mcp"
$exePath = "$installDir\kibuild-mcp.exe"

if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
}

Write-Host "Downloading kibuild-mcp..."
Invoke-RestMethod -Uri $url -OutFile $exePath

$userPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*kibuild-mcp*") {
    [System.Environment]::SetEnvironmentVariable("PATH", $userPath + ";$installDir", "User")
    $env:PATH = $env:PATH + ";$installDir"
}

$installedVersion = & $exePath --version 2>$null
Write-Host "✓ kibuild-mcp installed  ($installedVersion)"

# ── Step 2: Claude Code slash command ───────────────────────────────────────

$commandUrl = "https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/.claude/commands/setup-kibuild.md"
$commandDir = "$env:USERPROFILE\.claude\commands"
$commandFile = "$commandDir\setup-kibuild.md"

try {
    if (-not (Test-Path $commandDir)) {
        New-Item -ItemType Directory -Force -Path $commandDir | Out-Null
    }
    Invoke-RestMethod -Uri $commandUrl -OutFile $commandFile
    Write-Host "✓ /setup-kibuild command installed"
} catch {
    Write-Host "  (Could not download command file — skipping)"
}

# ── Step 3: Hand off to the native interactive setup ────────────────────────
# The binary writes the MCP config and verifies tools itself (one code path,
# identical on every OS). PowerShell's Read-Host reads from the console even
# when this script was loaded via `irm | iex`, so the prompts work directly.

Write-Host ""
& $exePath --setup
