$ErrorActionPreference = "Stop"

$Repo = "iammohdzaki/snippet-vault-go"
$DataDir = "$env:USERPROFILE\.snippet-vault"
$BinDir = "$DataDir\bin"

Write-Host "Installing Snippet Vault..." -ForegroundColor Cyan

# Detect Architecture
$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
$OsName = "windows"

# Fetch latest release
$ApiUrl = "https://api.github.com/repos/$Repo/releases/latest"
$Release = Invoke-RestMethod -Uri $ApiUrl
$Tag = $Release.tag_name

Write-Host "Found latest version: $Tag"

# Construct download URL
$ZipName = "snippet-vault_${OsName}_${Arch}.zip"
$DownloadUrl = "https://github.com/$Repo/releases/download/$Tag/$ZipName"
$TempZip = "$env:TEMP\$ZipName"
$TempExt = "$env:TEMP\snippet-vault-ext"

Write-Host "Downloading from $DownloadUrl..."
Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempZip

# Extract
if (Test-Path $TempExt) { Remove-Item $TempExt -Recurse -Force }
New-Item -ItemType Directory -Path $TempExt | Out-Null
Expand-Archive -Path $TempZip -DestinationPath $TempExt -Force

# Install
if (!(Test-Path $BinDir)) { New-Item -ItemType Directory -Force -Path $BinDir | Out-Null }
Move-Item -Path "$TempExt\snippet-vault.exe" -Destination "$BinDir\snippet-vault.exe" -Force
Move-Item -Path "$TempExt\snippet-vault-server.exe" -Destination "$BinDir\snippet-vault-server.exe" -Force

# Cleanup
Remove-Item $TempZip -Force
Remove-Item $TempExt -Recurse -Force

# Add to PATH
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notmatch [regex]::Escape($BinDir)) {
    Write-Host "Adding $BinDir to your User PATH..."
    $NewPath = if ($UserPath.EndsWith(";")) { "$UserPath$BinDir" } else { "$UserPath;$BinDir" }
    [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")
    Write-Host "Successfully added to PATH! (Restart your terminal to use 'snippet-vault')" -ForegroundColor Green
} else {
    Write-Host "Successfully installed! ($BinDir is already in your PATH)" -ForegroundColor Green
}
