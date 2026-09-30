#!/usr/bin/env bash

# ==========================================
cat << "EOF"
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
EOF
# ==========================================

COMMAND=$1
INSTALL_DIR="$HOME/.local/bin"
DATA_DIR="$HOME/.snippet-vault"

install_or_update() {
    echo "Building snippet-vault and snippet-vault-server..."
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$DATA_DIR"
    
    CGO_ENABLED=0 go build -ldflags="-w -s" -o "$INSTALL_DIR/snippet-vault" cmd/tui/main.go
    CGO_ENABLED=0 go build -ldflags="-w -s" -o "$INSTALL_DIR/snippet-vault-server" cmd/server/main.go
    
    echo "Installed successfully to $INSTALL_DIR!"
    echo "Please ensure '$INSTALL_DIR' is in your PATH."
}

uninstall() {
    echo "Removing binaries..."
    rm -f "$INSTALL_DIR/snippet-vault" "$INSTALL_DIR/snippet-vault-server"
    
    read -p "Do you want to permanently delete your database and data directory at $DATA_DIR? [y/N] " confirm
    if [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]]; then
        rm -rf "$DATA_DIR"
        echo "Data directory removed."
    else
        echo "Data directory preserved."
    fi
}

case "$COMMAND" in
    install|update)
        install_or_update
        ;;
    uninstall)
        uninstall
        ;;
    *)
        echo "Usage: $0 {install|update|uninstall}"
        exit 1
        ;;
esac
