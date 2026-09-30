param (
    [Parameter(Mandatory=$true)]
    [ValidateSet("install", "update", "uninstall")]
    [string]$Command
)

# ==========================================
Write-Host @"
                     /$$                                 /$$                                        /$$   /$$
                    |__/                                | $$                                       | $$  | $$
  /$$$$$$$ /$$$$$$$  /$$  /$$$$$$   /$$$$$$   /$$$$$$  /$$$$$$        /$$    /$$ /$$$$$$  /$$   /$$| $$ /$$$$$$
 /$$_____/| $$__  $$| $$ /$$__  $$ /$$__  $$ /$$__  $$|_  $$_//$$$$$$|  $$  /$$/|____  $$| $$  | $$| $$|_  $$_/
|  $$$$$$ | $$  \ $$| $$| $$  \ $$| $$  \ $$| $$$$$$$$  | $$ |______/ \  $$/$$/  /$$$$$$$| $$  | $$| $$  | $$
 \____  $$| $$  | $$| $$| $$  | $$| $$  | $$| $$_____/  | $$ /$$       \  $$$/  /$$__  $$| $$  | $$| $$  | $$ /$$
 /$$$$$$$/| $$  | $$| $$| $$$$$$$/| $$$$$$$/|  $$$$$$$  |  $$$$/        \  $/  |  $$$$$$$|  $$$$$$/| $$  |  $$$$/
|_______/ |__/  |__/|__/| $$____/ | $$____/  \_______/   \___/           \_/    \_______/ \______/ |__/   \___/
                        | $$      | $$
                        | $$      | $$
                        |__/      |__/
"@ -ForegroundColor Cyan
# ==========================================

$DataDir = "$env:USERPROFILE\.snippet-vault"
$BinDir = "$DataDir\bin"

function Install-Update {
    Write-Host "Building snippet-vault and snippet-vault-server..."
    if (!(Test-Path $BinDir)) {
        New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
    }
    
    $env:CGO_ENABLED="0"
    go build -ldflags="-w -s" -o "$BinDir\snippet-vault.exe" cmd/tui/main.go
    go build -ldflags="-w -s" -o "$BinDir\snippet-vault-server.exe" cmd/server/main.go
    
    # Automatically add to User PATH if not present
    $UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($UserPath -notmatch [regex]::Escape($BinDir)) {
        Write-Host "Adding $BinDir to your User PATH..."
        $NewPath = if ($UserPath.EndsWith(";")) { "$UserPath$BinDir" } else { "$UserPath;$BinDir" }
        [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")
        Write-Host "Successfully added to PATH! (You may need to restart your terminal to use 'snippet-vault')" -ForegroundColor Green
    } else {
        Write-Host "$BinDir is already in your PATH." -ForegroundColor Green
    }

    Write-Host "Installed successfully to $BinDir!"
}

function Uninstall {
    Write-Host "Removing binaries..."
    if (Test-Path "$BinDir\snippet-vault.exe") { Remove-Item "$BinDir\snippet-vault.exe" -Force }
    if (Test-Path "$BinDir\snippet-vault-server.exe") { Remove-Item "$BinDir\snippet-vault-server.exe" -Force }
    
    $confirm = Read-Host "Do you want to permanently delete your database and data directory at $DataDir? [y/N]"
    if ($confirm -match "^[yY](es)?$") {
        Remove-Item $DataDir -Recurse -Force
        Write-Host "Data directory removed."
    } else {
        Write-Host "Data directory preserved."
    }
}

switch ($Command) {
    "install" { Install-Update }
    "update" { Install-Update }
    "uninstall" { Uninstall }
}
