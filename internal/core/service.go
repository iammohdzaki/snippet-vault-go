package core

import (
	"errors"
	"strings"
	"time"
)

type SnippetService struct {
	repo SnippetRepository
}

func NewSnippetService(repo SnippetRepository) *SnippetService {
	return &SnippetService{repo: repo}
}

func (s *SnippetService) CreateSnippet(snippet *Snippet) error {
	if strings.TrimSpace(snippet.Title) == "" || strings.TrimSpace(snippet.Code) == "" {
		return errors.New("title and code are required")
	}

	if strings.TrimSpace(snippet.Language) == "" {
		snippet.Language = "plaintext"
	}

	now := time.Now().UTC()
	snippet.CreatedAt = now
	snippet.UpdatedAt = now

	return s.repo.Save(snippet)
}

func (s *SnippetService) GetAllSnippets() ([]Snippet, error) {
	return s.repo.GetAll()
}

func (s *SnippetService) UpdateSnippet(id string, snippet *Snippet) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id is required")
	}
	if strings.TrimSpace(snippet.Title) == "" || strings.TrimSpace(snippet.Code) == "" {
		return errors.New("title and code are required")
	}

	if strings.TrimSpace(snippet.Language) == "" {
		snippet.Language = "plaintext"
	}

	snippet.ID = id
	snippet.UpdatedAt = time.Now().UTC()

	return s.repo.Update(snippet)
}

func (s *SnippetService) DeleteSnippet(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id is required")
	}
	return s.repo.Delete(id)
}
