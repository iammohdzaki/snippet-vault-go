package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"snippet-vault-go/internal/server"
	"snippet-vault-go/internal/tui"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var Version = "dev"

func main() {
	// Build the server with io.Discard so logs don't corrupt the TUI
	srv := server.New(":8080", io.Discard)

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
