param (
    [Parameter(Mandatory=$true)]
    [ValidateSet("install", "update", "uninstall")]
    [string]$Command
)

$DataDir = "$env:USERPROFILE\.snippet-vault"
$BinDir = "$DataDir\bin"

function Install-Update {
    Write-Host "Building vault and vault-server..."
    if (!(Test-Path $BinDir)) {
        New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
    }
    
    $env:CGO_ENABLED="0"
    go build -ldflags="-w -s" -o "$BinDir\vault.exe" cmd/tui/main.go
    go build -ldflags="-w -s" -o "$BinDir\vault-server.exe" cmd/server/main.go
    
    Write-Host "Installed successfully to $BinDir!"
    Write-Host "Please ensure '$BinDir' is added to your system PATH environment variable to use 'vault' from anywhere."
}

function Uninstall {
    Write-Host "Removing binaries..."
    if (Test-Path "$BinDir\vault.exe") { Remove-Item "$BinDir\vault.exe" -Force }
    if (Test-Path "$BinDir\vault-server.exe") { Remove-Item "$BinDir\vault-server.exe" -Force }
    
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
