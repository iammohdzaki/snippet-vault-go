package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"snippet-vault-go/internal/core"
)

type APIClient struct {
	BaseURL string
}

func (c *APIClient) FetchSnippets() ([]core.Snippet, error) {
	resp, err := http.Get(c.BaseURL + "/snippets")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var snippets []core.Snippet
	err = json.NewDecoder(resp.Body).Decode(&snippets)
	return snippets, err
}

func (c *APIClient) CreateSnippet(snippet core.Snippet) error {
	body, err := json.Marshal(snippet)
	if err != nil {
		return err
	}

	resp, err := http.Post(c.BaseURL+"/snippets", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create snippet, status: %d", resp.StatusCode)
	}

	return nil
}

func (c *APIClient) UpdateSnippet(id string, snippet core.Snippet) error {
	body, err := json.Marshal(snippet)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, c.BaseURL+"/snippets/"+id, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update snippet, status: %d", resp.StatusCode)
	}
	return nil
}

func (c *APIClient) DeleteSnippet(id string) error {
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+"/snippets/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete snippet, status: %d", resp.StatusCode)
	}
	return nil
}
