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
type LinkService struct {
	repo LinkRepository
}

func NewLinkService(repo LinkRepository) *LinkService {
	return &LinkService{repo: repo}
}

func (s *LinkService) CreateLink(link *Link) error {
	if strings.TrimSpace(link.Title) == "" || strings.TrimSpace(link.URL) == "" {
		return errors.New("title and url are required")
	}
	now := time.Now().UTC()
	link.CreatedAt = now
	link.UpdatedAt = now
	return s.repo.SaveLink(link)
}

func (s *LinkService) GetAllLinks() ([]Link, error) {
	return s.repo.GetAllLinks()
}

func (s *LinkService) UpdateLink(id string, link *Link) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id is required")
	}
	if strings.TrimSpace(link.Title) == "" || strings.TrimSpace(link.URL) == "" {
		return errors.New("title and url are required")
	}
	link.ID = id
	link.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateLink(link)
}

func (s *LinkService) DeleteLink(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id is required")
	}
	return s.repo.DeleteLink(id)
}
