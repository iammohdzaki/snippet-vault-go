package core

import "testing"

// 1. Create a simple Mock Repository
type MockRepo struct {
	SaveFunc   func(snippet *Snippet) error
	UpdateFunc func(snippet *Snippet) error // <-- Add *
	DeleteFunc func(id string) error
}

func (m *MockRepo) Save(snippet *Snippet) error   { return m.SaveFunc(snippet) }
func (m *MockRepo) GetAll() ([]Snippet, error)    { return nil, nil }
func (m *MockRepo) Update(snippet *Snippet) error { return m.UpdateFunc(snippet) } // <-- Add *
func (m *MockRepo) Delete(id string) error        { return m.DeleteFunc(id) }

// 2. Write the test function (Must start with 'Test' and take *testing.T)
func TestCreateSnippet(t *testing.T) {
	// Define our test cases
	tests := []struct {
		name        string
		input       Snippet
		expectError bool
	}{
		{
			name:        "Valid Snippet",
			input:       Snippet{Title: "A", Code: "B"},
			expectError: false,
		},
		{
			name:        "Missing Title",
			input:       Snippet{Code: "B"}, // Title is empty
			expectError: true,
		},
		{
			name:        "Missing Code",
			input:       Snippet{Title: "A"}, // Code is empty
			expectError: true,
		},
	}

	// Setup our service with the mock repo
	mockRepo := &MockRepo{
		SaveFunc: func(snippet *Snippet) error { return nil }, // Assume DB success
	}
	service := NewSnippetService(mockRepo)

	// Run the table tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := service.CreateSnippet(&tc.input)

			if tc.expectError && err == nil {
				t.Errorf("Expected an error for %s, but got none", tc.name)
			}
			if !tc.expectError && err != nil {
				t.Errorf("Did not expect an error for %s, but got: %v", tc.name, err)
			}
		})
	}
}
