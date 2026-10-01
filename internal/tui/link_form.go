package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"snippet-vault-go/internal/core"
)

type LinkFormModel struct {
	TitleInput textinput.Model
	URLInput   textinput.Model
	DescInput  textinput.Model
	
	focusIndex int // 0: title, 1: url, 2: desc
	link       *core.Link
	width      int
	height     int
}

func NewLinkFormModel() LinkFormModel {
	ti := textinput.New()
	ti.Placeholder = "Link Title"
	ti.Prompt = "Title: "
	ti.CharLimit = 100
	ti.Width = 40

	ui := textinput.New()
	ui.Placeholder = "URL (https://...)"
	ui.Prompt = "URL:   "
	ui.CharLimit = 250
	ui.Width = 40

	di := textinput.New()
	di.Placeholder = "Description"
	di.Prompt = "Desc:  "
	di.CharLimit = 100
	di.Width = 40

	return LinkFormModel{
		TitleInput: ti,
		URLInput:   ui,
		DescInput:  di,
		focusIndex: -1,
	}
}

func (m *LinkFormModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	
	m.TitleInput.Width = width - 10
	m.URLInput.Width = width - 10
	m.DescInput.Width = width - 10
}

func (m *LinkFormModel) Update(msg tea.Msg) (LinkFormModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.focusIndex != -1 {
			switch msg.String() {
			case "tab", "shift+tab", "up", "down":
				s := msg.String()
				
				if s == "shift+tab" || s == "up" {
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

	if m.focusIndex == 0 {
		m.TitleInput, cmd = m.TitleInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focusIndex == 1 {
		m.URLInput, cmd = m.URLInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focusIndex == 2 {
		m.DescInput, cmd = m.DescInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return *m, tea.Batch(cmds...)
}

func (m *LinkFormModel) updateFocus() {
	m.TitleInput.Blur()
	m.URLInput.Blur()
	m.DescInput.Blur()

	if m.focusIndex == 0 {
		m.TitleInput.Focus()
	} else if m.focusIndex == 1 {
		m.URLInput.Focus()
	} else if m.focusIndex == 2 {
		m.DescInput.Focus()
	}
}

func (m *LinkFormModel) View() string {
	if m.focusIndex == -1 {
		// Read-only View! The user wants a beautiful UI with maybe an animation.
		// Let's create a stylized card for the link!
		if m.link == nil {
			return lipgloss.NewStyle().Foreground(subTextColor).Padding(2).Render("Select a link to view details.")
		}

		cardStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Width(m.width - 6).
			Height(m.height - 4)

		title := lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Render(m.link.Title)
		url := lipgloss.NewStyle().Foreground(activeBorder).Underline(true).Render(m.link.URL)
		desc := lipgloss.NewStyle().Foreground(textColor).Render(m.link.Description)
		
		// ASCII Art / Graphic "Animation" (a stylized web icon)
		webIcon := lipgloss.NewStyle().Foreground(primaryColor).Render(`
    🌐
   /  \
  /____\
 `)

		content := lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			url,
			"",
			desc,
			"",
			lipgloss.NewStyle().Foreground(subTextColor).Render("Press 'o' to open in browser!"),
		)
		
		layout := lipgloss.JoinHorizontal(lipgloss.Top, webIcon, "   ", content)

		return cardStyle.Render(layout)
	}

	// Edit View
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Render("Edit Link"),
		"",
		m.TitleInput.View(),
		m.URLInput.View(),
		m.DescInput.View(),
	)
}

func (m *LinkFormModel) SetLink(link *core.Link) {
	m.link = link
	if link == nil {
		m.TitleInput.SetValue("")
		m.URLInput.SetValue("")
		m.DescInput.SetValue("")
		return
	}
	m.TitleInput.SetValue(link.Title)
	m.URLInput.SetValue(link.URL)
	m.DescInput.SetValue(link.Description)
}

func (m *LinkFormModel) Focus() {
	m.focusIndex = 0
	m.updateFocus()
}

func (m *LinkFormModel) Blur() {
	m.focusIndex = -1
	m.updateFocus()
}

func (m *LinkFormModel) GetLink() core.Link {
	id := ""
	if m.link != nil {
		id = m.link.ID
	}
	return core.Link{
		ID:          id,
		Title:       m.TitleInput.Value(),
		URL:         m.URLInput.Value(),
		Description: m.DescInput.Value(),
	}
}
