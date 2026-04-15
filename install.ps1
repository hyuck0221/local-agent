#Requires -Version 5.1
# local-agent installer for Windows. Downloads the latest release binary and
# drops it into %LOCALAPPDATA%\Programs\local-agent, adding it to the user PATH.
$ErrorActionPreference = 'Stop'

$Repo    = if ($env:LOCAL_AGENT_REPO)    { $env:LOCAL_AGENT_REPO }    else { 'hyuck0221/local-agent' }
$Version = if ($env:LOCAL_AGENT_VERSION) { $env:LOCAL_AGENT_VERSION } else { 'latest' }

$arch = if ([Environment]::Is64BitOperatingSystem) {
  if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'amd64' }
} else { throw 'unsupported arch: local-agent requires 64-bit Windows' }

if ($Version -eq 'latest') {
  $rel = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
  $Version = $rel.tag_name
}

$stripped = $Version -replace '^v', ''
$asset = "local-agent_${stripped}_windows_${arch}.zip"
$url   = "https://github.com/$Repo/releases/download/$Version/$asset"

$tmp = Join-Path $env:TEMP "local-agent-$([guid]::NewGuid())"
New-Item -ItemType Directory -Path $tmp | Out-Null
try {
  Write-Host "Downloading $asset..."
  $zipPath = Join-Path $tmp $asset
  Invoke-WebRequest -Uri $url -OutFile $zipPath
  Expand-Archive -Path $zipPath -DestinationPath $tmp -Force

  $dest = Join-Path $env:LOCALAPPDATA 'Programs\local-agent'
  New-Item -ItemType Directory -Force -Path $dest | Out-Null
  Copy-Item (Join-Path $tmp 'local-agent.exe') (Join-Path $dest 'local-agent.exe') -Force

  $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
  if (-not ($userPath -split ';' | Where-Object { $_ -eq $dest })) {
    [Environment]::SetEnvironmentVariable('Path', "$userPath;$dest", 'User')
    Write-Host "Added $dest to user PATH (open a new terminal to use it)."
  }

  Write-Host "Installed local-agent $Version to $dest"
  Write-Host ""
  Write-Host "Next: local-agent start"
} finally {
  Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}
