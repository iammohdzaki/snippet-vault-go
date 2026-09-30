package core

import "errors"

var ErrNotFound = errors.New("snippet not found")

type SnippetRepository interface {
	Save(snippet *Snippet) error
	GetAll() ([]Snippet, error)
	Update(snippet *Snippet) error
	Delete(id string) error
}
