$DataDir = "$env:USERPROFILE\.snippet-vault"
$BinDir = "$DataDir\bin"

Write-Host "Uninstalling Snippet Vault..." -ForegroundColor Cyan

if (Test-Path "$BinDir\snippet-vault.exe") { Remove-Item "$BinDir\snippet-vault.exe" -Force }
if (Test-Path "$BinDir\snippet-vault-server.exe") { Remove-Item "$BinDir\snippet-vault-server.exe" -Force }

Write-Host "Binaries removed from $BinDir."

$confirm = Read-Host "Do you want to permanently delete your snippets database at $DataDir? [y/N]"
if ($confirm -match "^[yY](es)?$") {
    Remove-Item $DataDir -Recurse -Force
    Write-Host "Data directory removed."
} else {
    Write-Host "Data directory preserved. You can reinstall later without losing your snippets."
}

Write-Host "Uninstall complete." -ForegroundColor Green
