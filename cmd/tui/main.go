package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"snippet-vault-go/internal/server"
	"snippet-vault-go/internal/tui"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var Version = "dev"

func getDBPath() string {
	// 1. Check for manual environment override
	if env := os.Getenv("SNIPPET_VAULT_DB"); env != "" {
		return env
	}

	// 2. If it's a dev build, use a local DB in the current directory
	if Version == "dev" {
		return "vault-dev.db"
	}

	// 3. Otherwise, use the production DB in the user's home directory
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".snippet-vault", "vault.db")
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "update":
			fmt.Println("To automatically update Snippet Vault to the latest version, run:")
			if runtime.GOOS == "windows" {
				fmt.Println("irm https://raw.githubusercontent.com/iammohdzaki/snippet-vault-go/main/scripts/install.ps1 | iex")
			} else {
				fmt.Println("curl -sSL https://raw.githubusercontent.com/iammohdzaki/snippet-vault-go/main/scripts/install.sh | bash")
			}
			return
		case "uninstall":
			fmt.Println("To safely uninstall Snippet Vault, run:")
			if runtime.GOOS == "windows" {
				fmt.Println("irm https://raw.githubusercontent.com/iammohdzaki/snippet-vault-go/main/scripts/uninstall.ps1 | iex")
			} else {
				fmt.Println("curl -sSL https://raw.githubusercontent.com/iammohdzaki/snippet-vault-go/main/scripts/uninstall.sh | bash")
			}
			return
		case "version", "-v", "--version":
			fmt.Printf("Snippet Vault %s\n", Version)
			return
		}
	}

	// Build the server with io.Discard so logs don't corrupt the TUI
	srv := server.New(":8080", io.Discard, getDBPath())

	// Start the server in a background goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// If port 8080 is already in use by another instance, we just ignore it
			// and let the TUI connect to the existing server.
		}
	}()

	// Give the background socket a tiny head start before Bubble Tea fires Init()
	time.Sleep(15 * time.Millisecond)

	// Initialize the HTTP client pointing to our local server
	client := &tui.APIClient{BaseURL: "http://localhost:8080"}

	// Initialize our Bubble Tea Model
	initialModel := tui.NewAppModel(client)

	// Start the terminal application
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting TUI: %v", err)
		os.Exit(1)
	}

	// User pressed 'q' -> Gracefully shut down the background server
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
