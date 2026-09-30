package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"snippet-vault-go/internal/core"
)

type EditorModel struct {
	TitleInput textinput.Model
	LangInput  textinput.Model
	CodeArea   textarea.Model
	
	focusIndex int // 0: title, 1: lang, 2: code
	snippet    *core.Snippet
	width      int
	height     int
}

func NewEditorModel() EditorModel {
	ti := textinput.New()
	ti.Placeholder = "Snippet Title"
	ti.Prompt = "Title: "
	ti.CharLimit = 100
	ti.Width = 30

	li := textinput.New()
	li.Placeholder = "Language (e.g. go, python)"
	li.Prompt = "Lang:  "
	li.CharLimit = 50
	li.Width = 30

	ta := textarea.New()
	ta.Placeholder = "Write your code here..."
	ta.Prompt = "" // Removes the default "┃ " prompt
	ta.ShowLineNumbers = true
	ta.CharLimit = 0

	// Style the active line (purple background) and line numbers
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(primaryColor).Foreground(lipgloss.Color("#282a36"))
	ta.FocusedStyle.CursorLineNumber = lipgloss.NewStyle().Background(primaryColor).Foreground(lipgloss.Color("#282a36")).Bold(true)
	ta.FocusedStyle.LineNumber = lipgloss.NewStyle().Foreground(subTextColor)
	ta.FocusedStyle.EndOfBuffer = lipgloss.NewStyle().Foreground(subTextColor)
	
	// Style the cursor to match the pink block in the screenshot
	ta.Cursor.Style = lipgloss.NewStyle().Background(secondaryColor).Foreground(lipgloss.Color("#282a36"))

	// When blurred, don't show the purple background highlight
	ta.BlurredStyle.CursorLine = lipgloss.NewStyle()
	ta.BlurredStyle.CursorLineNumber = lipgloss.NewStyle().Foreground(subTextColor)
	ta.BlurredStyle.LineNumber = lipgloss.NewStyle().Foreground(subTextColor)
	ta.BlurredStyle.EndOfBuffer = lipgloss.NewStyle().Foreground(subTextColor)

	// Set end of buffer character (the tildes)
	ta.EndOfBufferCharacter = '~'

	return EditorModel{
		TitleInput: ti,
		LangInput:  li,
		CodeArea:   ta,
		focusIndex: -1, // -1 means unfocused
	}
}

func (m *EditorModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	
	// Adjust CodeArea height: total height minus 2 rows for inputs, minus margins
	m.CodeArea.SetWidth(width)
	m.CodeArea.SetHeight(height - 4)
	
	m.TitleInput.Width = width - 10
	m.LangInput.Width = width - 10
}

func (m *EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.focusIndex != -1 { // If editor is active
			switch msg.String() {
			case "tab", "shift+tab":
				s := msg.String()
				
				if s == "shift+tab" {
					m.focusIndex--
					if m.focusIndex < 0 {
						m.focusIndex = 2
					}
				} else {
					m.focusIndex++
					if m.focusIndex > 2 {
						m.focusIndex = 0
					}
				}
				m.updateFocus()
				return *m, nil
			}
		}
	}

	// Route messages to the focused component
	if m.focusIndex == 0 {
		m.TitleInput, cmd = m.TitleInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focusIndex == 1 {
		m.LangInput, cmd = m.LangInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focusIndex == 2 {
		m.CodeArea, cmd = m.CodeArea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return *m, tea.Batch(cmds...)
}

func (m *EditorModel) updateFocus() {
	m.TitleInput.Blur()
	m.LangInput.Blur()
	m.CodeArea.Blur()

	if m.focusIndex == 0 {
		m.TitleInput.Focus()
	} else if m.focusIndex == 1 {
		m.LangInput.Focus()
	} else if m.focusIndex == 2 {
		m.CodeArea.Focus()
	}
}

func (m *EditorModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.TitleInput.View(),
		m.LangInput.View(),
		"", // empty line spacer
		m.CodeArea.View(),
	)
}

func (m *EditorModel) SetSnippet(snippet *core.Snippet) {
	m.snippet = snippet
	if snippet == nil {
		m.TitleInput.SetValue("")
		m.LangInput.SetValue("")
		m.CodeArea.SetValue("")
		return
	}
	m.TitleInput.SetValue(snippet.Title)
	m.LangInput.SetValue(snippet.Language)
	m.CodeArea.SetValue(snippet.Code)
}

func (m *EditorModel) Focus() {
	m.focusIndex = 0 // Focus Title first
	m.updateFocus()
}

func (m *EditorModel) Blur() {
	m.focusIndex = -1
	m.updateFocus()
}

func (m *EditorModel) GetSnippet() core.Snippet {
	id := ""
	if m.snippet != nil {
		id = m.snippet.ID
	}
	return core.Snippet{
		ID:       id,
		Title:    m.TitleInput.Value(),
		Language: m.LangInput.Value(),
		Code:     m.CodeArea.Value(),
	}
}
