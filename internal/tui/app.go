package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type activeTab int

const (
	tabSnippets activeTab = iota
	tabLinks
)

// AppModel is the root router for the application
type AppModel struct {
	activeTab     activeTab
	width         int
	height        int
	
	snippetsModel SnippetsModel
	linksModel    LinksModel
}

func NewAppModel(client *APIClient) AppModel {
	return AppModel{
		activeTab:     tabSnippets,
		snippetsModel: NewSnippetsModel(client),
		linksModel:    NewLinksModel(client),
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		m.snippetsModel.Init(),
		m.linksModel.Init(),
	)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global tab switching with left/right arrow keys
		// Only allow switching tabs if we are in the list view (not typing in editor)
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		
		canSwitchTabs := true
		if m.activeTab == tabSnippets && m.snippetsModel.focus == focusEditor {
			canSwitchTabs = false
		}
		if m.activeTab == tabLinks && m.linksModel.focus == focusEditor {
			canSwitchTabs = false
		}
		
		if canSwitchTabs {
			if msg.String() == "right" {
				m.activeTab = tabLinks
				m.updateSizes()
				return m, nil
			} else if msg.String() == "left" {
				m.activeTab = tabSnippets
				m.updateSizes()
				return m, nil
			}
		}
	
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()
		return m, nil
	}

	// Route updates to the active tab model
	if m.activeTab == tabSnippets {
		var sm tea.Model
		sm, cmd = m.snippetsModel.Update(msg)
		m.snippetsModel = sm.(SnippetsModel)
		cmds = append(cmds, cmd)
	} else if m.activeTab == tabLinks {
		var lm tea.Model
		lm, cmd = m.linksModel.Update(msg)
		m.linksModel = lm.(LinksModel)
		cmds = append(cmds, cmd)
	}

	// Route non-key messages to inactive tabs
	if _, isKey := msg.(tea.KeyMsg); !isKey {
		if m.activeTab != tabSnippets {
			var sm tea.Model
			sm, cmd = m.snippetsModel.Update(msg)
			m.snippetsModel = sm.(SnippetsModel)
			cmds = append(cmds, cmd)
		}
		if m.activeTab != tabLinks {
			var lm tea.Model
			lm, cmd = m.linksModel.Update(msg)
			m.linksModel = lm.(LinksModel)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *AppModel) updateSizes() {
	if m.width == 0 || m.height == 0 {
		return
	}
	
	h, v := appStyle.GetFrameSize()
	
	// Subtract header and tabs from the available height for the children
	// Header: 1 line text + 1 line margin = 2
	// Tabs: 1 line text + 1 line top border = 2
	// Extra safety margin for Windows terminal scrolling = 2
	// Total reserved: 6 lines + appStyle vertical padding
	childHeight := m.height - v - 6
	childWidth := m.width - h
	
	msg := tea.WindowSizeMsg{Width: childWidth, Height: childHeight}
	
	sm, _ := m.snippetsModel.Update(msg)
	m.snippetsModel = sm.(SnippetsModel)

	lm, _ := m.linksModel.Update(msg)
	m.linksModel = lm.(LinksModel)
}

func (m AppModel) View() string {
	// 1. App Header
	header := titleStyle.Render("Snippet Vault")

	// 2. Tab Bar
	tabs := []string{"Snippets", "Links"}
	var renderedTabs []string
	
	for i, t := range tabs {
		if activeTab(i) == m.activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, inactiveTabStyle.Render(t))
		}
	}
	
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	tabBar = tabBarStyle.Render(tabBar)

	// 3. Active Workspace
	var content string
	if m.activeTab == tabSnippets {
		content = m.snippetsModel.View()
	} else if m.activeTab == tabLinks {
		content = m.linksModel.View()
	}

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		header,
		tabBar,
		content,
	))
}
