#!/usr/bin/env bash

INSTALL_DIR="$HOME/.local/bin"
DATA_DIR="$HOME/.snippet-vault"

echo "Uninstalling Snippet Vault..."
rm -f "$INSTALL_DIR/snippet-vault" "$INSTALL_DIR/snippet-vault-server"

echo "Binaries removed from $INSTALL_DIR."

read -p "Do you want to permanently delete your snippets database at $DATA_DIR? [y/N] " confirm
if [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]]; then
    rm -rf "$DATA_DIR"
    echo "Data directory removed."
else
    echo "Data directory preserved. You can reinstall later without losing your snippets."
fi

echo "Uninstall complete."
