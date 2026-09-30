package tui

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"snippet-vault-go/internal/core"
)

type state int

const (
	stateLoading state = iota
	stateLoaded
	stateError
)

type focus int

const (
	focusList focus = iota
	focusEditor
)

// AppModel orchestrates the UI components
type AppModel struct {
	client     *APIClient
	state      state
	focus      focus
	err        error
	
	width      int
	height     int

	spinner    spinner.Model
	list       list.Model
	editor     EditorModel

	snippets   []core.Snippet
	
	confirmDelete bool
	statusMsg     string
}

// snippetItem wraps core.Snippet to implement list.DefaultItem
type snippetItem struct {
	core.Snippet
}

func (i snippetItem) Title() string       { return i.Snippet.Title }
func (i snippetItem) Description() string { return "Language: " + i.Snippet.Language }
func (i snippetItem) FilterValue() string { return i.Snippet.Title }

// Custom messages
type snippetsLoadedMsg []core.Snippet
type errMsg struct{ err error }
type clearStatusMsg struct{}

// NewAppModel creates the root application model
func NewAppModel(client *APIClient) AppModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	// Customize list delegate to reduce spacing
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Padding(0, 0, 0, 1)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().Padding(0, 0, 0, 1).Foreground(subTextColor)
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(activeBorder).Foreground(activeBorder).Padding(0, 0, 0, 1)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(activeBorder).Foreground(subTextColor).Padding(0, 0, 0, 1)

	// Initial placeholder list
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "+ Add Snippet (ctrl+n)"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.FilterInput.Prompt = "Search: "
	
	// Tighter margins for the list
	l.Styles.PaginationStyle = lipgloss.NewStyle().PaddingLeft(2)

	return AppModel{
		client:   client,
		state:    stateLoading,
		focus:    focusList,
		spinner:  s,
		list:     l,
		editor:   NewEditorModel(),
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchSnippetsCmd(),
	)
}

func (m AppModel) fetchSnippetsCmd() tea.Cmd {
	return func() tea.Msg {
		snippets, err := m.client.FetchSnippets()
		if err != nil {
			return errMsg{err}
		}
		return snippetsLoadedMsg(snippets)
	}
}

type snippetSavedMsg struct{}

func (m AppModel) saveSnippetCmd(snippet core.Snippet) tea.Cmd {
	return func() tea.Msg {
		if snippet.ID != "" {
			err := m.client.UpdateSnippet(snippet.ID, snippet)
			if err != nil {
				return errMsg{err}
			}
		} else {
			err := m.client.CreateSnippet(snippet)
			if err != nil {
				return errMsg{err}
			}
		}
		return snippetSavedMsg{}
	}
}

type snippetDeletedMsg struct{}

func (m AppModel) deleteSnippetCmd(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.DeleteSnippet(id)
		if err != nil {
			return errMsg{err}
		}
		return snippetDeletedMsg{}
	}
}

func clearStatusCmd() tea.Cmd {
	return tea.Tick(time.Second*3, func(_ time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 1. Handle Global Escapes and Quits
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		
		// If we are confirming deletion, intercept all keys
		if m.confirmDelete {
			if msg.String() == "y" || msg.String() == "Y" {
				m.confirmDelete = false
				
				if i := m.list.SelectedItem(); i != nil {
					item := i.(snippetItem)
					m.statusMsg = fmt.Sprintf("Deleting '%s'...", item.Title())
					cmds = append(cmds, m.deleteSnippetCmd(item.ID))
					return m, tea.Batch(cmds...)
				}
			} else if msg.String() == "n" || msg.String() == "N" || msg.String() == "esc" {
				m.confirmDelete = false
				m.statusMsg = "Deletion cancelled."
				cmds = append(cmds, clearStatusCmd())
			}
			return m, tea.Batch(cmds...)
		}

		// 2. Handle Tab (Focus Switching)
		if msg.String() == "tab" && m.list.FilterState() != list.Filtering {
			if m.focus == focusList {
				m.focus = focusEditor
				m.editor.Focus()
			} else {
				m.focus = focusList
				m.editor.Blur()
			}
			return m, nil
		}

		// 3. Handle Editor Focus
		if m.focus == focusEditor {
			if msg.String() == "esc" {
				// Escape returns focus to the list
				m.focus = focusList
				m.editor.Blur()
				return m, nil
			} else if msg.String() == "ctrl+s" {
				newSnippet := m.editor.GetSnippet()
				m.statusMsg = fmt.Sprintf("Saving '%s'...", newSnippet.Title)
				cmds = append(cmds, m.saveSnippetCmd(newSnippet))
				return m, tea.Batch(cmds...)
			} else if msg.String() == "ctrl+d" {
				if m.list.SelectedItem() != nil {
					m.confirmDelete = true
				}
				return m, nil
			} else if msg.String() == "ctrl+b" {
				clipboard.WriteAll(m.editor.GetSnippet().Code)
				m.statusMsg = "Copied code to clipboard!"
				cmds = append(cmds, clearStatusCmd())
				return m, tea.Batch(cmds...)
			}
			
			// Let editor handle other typing
			m.editor, cmd = m.editor.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		// 4. Handle List Focus
		if m.focus == focusList {
			if msg.String() == "ctrl+n" && m.list.FilterState() != list.Filtering {
				// Create new snippet
				m.editor.SetSnippet(nil)
				m.focus = focusEditor
				m.editor.Focus()
				m.statusMsg = "Creating new snippet. Fill in Title, Lang, and Code."
				return m, nil
			}
			if msg.String() == "q" && m.list.FilterState() != list.Filtering {
				return m, tea.Quit
			}
			if msg.String() == "c" && m.list.FilterState() != list.Filtering {
				if i := m.list.SelectedItem(); i != nil {
					item := i.(snippetItem)
					clipboard.WriteAll(item.Code)
					m.statusMsg = "Copied code to clipboard!"
					cmds = append(cmds, clearStatusCmd())
				}
				return m, tea.Batch(cmds...)
			}
			if msg.String() == "d" && m.list.FilterState() != list.Filtering {
				if m.list.SelectedItem() != nil {
					m.confirmDelete = true
				}
				return m, nil
			}
			
			// Forward key to list
			oldIdx := m.list.Index()
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
			if m.list.Index() != oldIdx {
				m.updateViewportContent()
			}
			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()
		
		// Window size updates need to reach editor so it resizes
		m.editor, cmd = m.editor.Update(msg)
		cmds = append(cmds, cmd)

	case snippetsLoadedMsg:
		m.state = stateLoaded
		m.snippets = msg
		items := make([]list.Item, len(msg))
		for i, s := range msg {
			items[i] = snippetItem{s}
		}
		cmd = m.list.SetItems(items)
		cmds = append(cmds, cmd)
		m.updateViewportContent()

	case snippetSavedMsg:
		m.statusMsg = "Snippet saved! Reloading..."
		cmds = append(cmds, clearStatusCmd(), m.fetchSnippetsCmd())

	case snippetDeletedMsg:
		m.statusMsg = "Snippet deleted! Reloading..."
		cmds = append(cmds, clearStatusCmd(), m.fetchSnippetsCmd())

	case clearStatusMsg:
		m.statusMsg = ""

	case errMsg:
		m.state = stateError
		m.err = msg.err
	}

	// Route normal non-key messages
	switch m.state {
	case stateLoading:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case stateLoaded:
		// Forward non-key messages to both
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
		m.editor, cmd = m.editor.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *AppModel) updateSizes() {
	if m.width == 0 || m.height == 0 {
		return
	}
	
	h, v := appStyle.GetFrameSize()
	availableWidth := m.width - h
	// Subtract lines for Header (3) and Help text (2)
	availableHeight := m.height - v - 5

	listWidth := availableWidth / 3
	vpWidth := availableWidth - listWidth

	paneH, paneV := basePaneStyle.GetFrameSize()

	m.list.SetSize(listWidth-paneH, availableHeight-paneV)
	m.editor.SetSize(vpWidth-paneH, availableHeight-paneV)
}

func (m *AppModel) updateViewportContent() {
	i := m.list.SelectedItem()
	if i == nil {
		m.editor.SetSnippet(nil)
		return
	}
	item := i.(snippetItem)
	m.editor.SetSnippet(&item.Snippet)
}

func (m AppModel) View() string {
	if m.state == stateLoading {
		return appStyle.Render(fmt.Sprintf("%s Loading snippets...", m.spinner.View()))
	}
	if m.state == stateError {
		return appStyle.Render(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}
	
	// HEADER
	header := titleStyle.Render("Snippet Vault")
	
	// Calculate dynamic height and width for panes
	h, v := appStyle.GetFrameSize()
	paneH, paneV := basePaneStyle.GetFrameSize()
	
	targetHeight := m.height - v - paneV - 5
	availableWidth := m.width - h
	
	listWidth := (availableWidth / 3) - paneH
	vpWidth := (availableWidth - (availableWidth / 3)) - paneH

	// PANES
	listStyle := basePaneStyle.Height(targetHeight).Width(listWidth)
	vpStyle := basePaneStyle.Height(targetHeight).Width(vpWidth)

	if m.focus == focusList {
		listStyle = activePaneStyle.Height(targetHeight).Width(listWidth)
	} else if m.focus == focusEditor {
		vpStyle = activePaneStyle.Height(targetHeight).Width(vpWidth)
	}

	listView := listStyle.Render(m.list.View())
	vpView := vpStyle.Render(m.editor.View())
	split := lipgloss.JoinHorizontal(lipgloss.Top, listView, vpView)

	// FOOTER (Help & Status)
	var helpText string
	if m.confirmDelete {
		helpText = errorStyle.Render("Are you sure you want to delete this snippet? (y/N)")
	} else if m.focus == focusEditor {
		helpText = "tab/shift+tab: cycle fields • esc: back to list • ctrl+s: save • ctrl+b: copy • ctrl+d: delete"
	} else {
		helpText = "tab: edit snippet • ctrl+n: new snippet • c: copy • d: delete • q: quit"
	}
	
	if m.statusMsg != "" {
		helpText = fmt.Sprintf("%s  |  %s", helpText, activeTextStyle.Render(m.statusMsg))
	}

	help := lipgloss.NewStyle().Foreground(subTextColor).Render("\n  " + helpText)
	
	// JOIN ALL
	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left, header, split, help))
}
