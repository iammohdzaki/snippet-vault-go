package tui

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"snippet-vault-go/internal/core"
)

// LinksModel orchestrates the UI components
type LinksModel struct {
	client     *APIClient
	state      state
	focus      focus
	err        error
	
	width      int
	height     int

	spinner    spinner.Model
	list       list.Model
	form       LinkFormModel

	links      []core.Link
	
	confirmDelete bool
	statusMsg     string
}

type linkItem struct {
	core.Link
}

func (i linkItem) Title() string       { return i.Link.Title }
func (i linkItem) Description() string { return i.Link.URL }
func (i linkItem) FilterValue() string { return i.Link.Title + " " + i.Link.URL }

type linksLoadedMsg []core.Link
type linkSavedMsg struct{}
type linkDeletedMsg struct{}

func NewLinksModel(client *APIClient) LinksModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Padding(0, 0, 0, 2)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().Padding(0, 0, 0, 2).Foreground(subTextColor)
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(activeBorder).Foreground(activeBorder).Padding(0, 0, 0, 1)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(activeBorder).Foreground(subTextColor).Padding(0, 0, 0, 1)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowTitle(false)
	
	l.Styles.TitleBar = lipgloss.NewStyle().Margin(0, 0, 1, 0)
	l.Styles.Title = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#282a36")).
			Background(primaryColor).
			Padding(0, 1)
	
	l.Styles.NoItems = lipgloss.NewStyle().Margin(0, 0, 0, 1).Foreground(subTextColor)
	l.Styles.PaginationStyle = lipgloss.NewStyle().Margin(0, 0, 0, 1).Foreground(subTextColor)
	l.Styles.HelpStyle = lipgloss.NewStyle().Margin(0, 0, 0, 1)
	
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.FilterInput.Prompt = "Search: "

	return LinksModel{
		client:   client,
		state:    stateLoading,
		focus:    focusList,
		spinner:  s,
		list:     l,
		form:     NewLinkFormModel(),
	}
}

func (m LinksModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchLinksCmd(),
	)
}

func (m LinksModel) fetchLinksCmd() tea.Cmd {
	return func() tea.Msg {
		links, err := m.client.FetchLinks()
		if err != nil {
			return errMsg{err}
		}
		return linksLoadedMsg(links)
	}
}

func (m LinksModel) saveLinkCmd(link core.Link) tea.Cmd {
	return func() tea.Msg {
		if link.ID != "" {
			err := m.client.UpdateLink(link.ID, link)
			if err != nil {
				return errMsg{err}
			}
		} else {
			err := m.client.CreateLink(link)
			if err != nil {
				return errMsg{err}
			}
		}
		return linkSavedMsg{}
	}
}

func (m LinksModel) deleteLinkCmd(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.DeleteLink(id)
		if err != nil {
			return errMsg{err}
		}
		return linkDeletedMsg{}
	}
}

func (m LinksModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case spinner.TickMsg:
		return m, nil
	case tea.KeyMsg:
		if m.confirmDelete {
			if msg.String() == "y" || msg.String() == "Y" {
				m.confirmDelete = false
				
				if i := m.list.SelectedItem(); i != nil {
					item := i.(linkItem)
					m.statusMsg = fmt.Sprintf("Deleting '%s'...", item.Title())
					cmds = append(cmds, m.deleteLinkCmd(item.ID))
					return m, tea.Batch(cmds...)
				}
			} else if msg.String() == "n" || msg.String() == "N" || msg.String() == "esc" {
				m.confirmDelete = false
				m.statusMsg = "Deletion cancelled."
				cmds = append(cmds, clearStatusCmd())
			}
			return m, tea.Batch(cmds...)
		}

		if msg.String() == "tab" && m.list.FilterState() != list.Filtering {
			if m.focus == focusList {
				m.focus = focusEditor
				m.form.Focus()
			} else {
				m.focus = focusList
				m.form.Blur()
			}
			return m, nil
		}

		if m.focus == focusEditor {
			if msg.String() == "esc" {
				m.focus = focusList
				m.form.Blur()
				m.updateViewportContent() // RESTORE the selected item view!
				m.statusMsg = "" // Clear the "Creating..." message
				return m, nil
			} else if msg.String() == "ctrl+s" {
				newLink := m.form.GetLink()
				m.statusMsg = fmt.Sprintf("Saving '%s'...", newLink.Title)
				cmds = append(cmds, m.saveLinkCmd(newLink))
				m.focus = focusList
				m.form.Blur()
				return m, tea.Batch(cmds...)
			} else if msg.String() == "ctrl+d" {
				if m.list.SelectedItem() != nil {
					m.confirmDelete = true
				}
				return m, nil
			}
			
			m.form, cmd = m.form.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		if m.focus == focusList {
			if msg.String() == "ctrl+n" && m.list.FilterState() != list.Filtering {
				m.form.SetLink(nil)
				m.focus = focusEditor
				m.form.Focus()
				m.statusMsg = "Creating new link. Fill in details."
				return m, nil
			}
			if msg.String() == "c" && m.list.FilterState() != list.Filtering {
				if i := m.list.SelectedItem(); i != nil {
					item := i.(linkItem)
					clipboard.WriteAll(item.URL)
					m.statusMsg = "Copied URL to clipboard!"
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
			if msg.String() == "o" && m.list.FilterState() != list.Filtering {
				// We can add "open in browser" here. For now just copy to clipboard
				if i := m.list.SelectedItem(); i != nil {
					item := i.(linkItem)
					clipboard.WriteAll(item.URL)
					m.statusMsg = "Copied URL to open in browser!"
					cmds = append(cmds, clearStatusCmd())
				}
				return m, tea.Batch(cmds...)
			}
			
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
		
		m.form, cmd = m.form.Update(msg)
		cmds = append(cmds, cmd)

	case linksLoadedMsg:
		m.state = stateLoaded
		m.links = msg
		items := make([]list.Item, len(msg))
		for i, s := range msg {
			items[i] = linkItem{s}
		}
		cmd = m.list.SetItems(items)
		cmds = append(cmds, cmd)
		m.updateViewportContent()

	case linkSavedMsg:
		m.statusMsg = "Link saved! Reloading..."
		cmds = append(cmds, clearStatusCmd(), m.fetchLinksCmd())

	case linkDeletedMsg:
		m.statusMsg = "Link deleted! Reloading..."
		cmds = append(cmds, clearStatusCmd(), m.fetchLinksCmd())

	case clearStatusMsg:
		m.statusMsg = ""

	case errMsg:
		m.state = stateError
		m.err = msg.err
	}

	switch m.state {
	case stateLoading:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case stateLoaded:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
		m.form, cmd = m.form.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *LinksModel) updateSizes() {
	if m.width == 0 || m.height == 0 {
		return
	}
	paneH, paneV := basePaneStyle.GetFrameSize()
	
	targetHeight := m.height - paneV - 2 
	availableWidth := m.width
	
	listWidth := (availableWidth / 3) - paneH
	vpWidth := (availableWidth - (availableWidth / 3)) - paneH

	
	m.list.SetSize(listWidth-2, targetHeight)
	m.form.SetSize(vpWidth-2, targetHeight)
}

func (m *LinksModel) updateViewportContent() {
	i := m.list.SelectedItem()
	if i == nil {
		m.form.SetLink(nil)
		return
	}
	item := i.(linkItem)
	m.form.SetLink(&item.Link)
}

func (m LinksModel) View() string {
	if m.state == stateLoading {
		return fmt.Sprintf("%s Loading links...", m.spinner.View())
	}
	if m.state == stateError {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}
	
	paneH, paneV := basePaneStyle.GetFrameSize()
	targetHeight := m.height - paneV - 2
	availableWidth := m.width
	
	listWidth := (availableWidth / 3) - paneH
	vpWidth := (availableWidth - (availableWidth / 3)) - paneH

	listStyle := basePaneStyle.Height(targetHeight).Width(listWidth)
	vpStyle := basePaneStyle.Height(targetHeight).Width(vpWidth)

	if m.focus == focusList {
		listStyle = activePaneStyle.Height(targetHeight).Width(listWidth)
	} else if m.focus == focusEditor {
		vpStyle = activePaneStyle.Height(targetHeight).Width(vpWidth)
	}

	listView := listStyle.Render(m.list.View())
	vpView := vpStyle.Render(m.form.View())
	split := lipgloss.JoinHorizontal(lipgloss.Top, listView, vpView)

	var helpText string
	if m.confirmDelete {
		helpText = errorStyle.Render("Are you sure you want to delete this link? (y/N)")
	} else if m.focus == focusEditor {
		helpText = "tab/shift+tab: cycle fields • esc: back to list • ctrl+s: save • ctrl+d: delete"
	} else {
		addBtn := lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("ctrl+n: new link")
		helpText = fmt.Sprintf("tab: edit link • %s • c: copy • o: open/copy • d: delete • q: quit • left/right: switch tabs", addBtn)
	}
	
	if m.statusMsg != "" {
		helpText = fmt.Sprintf("%s  |  %s", helpText, activeTextStyle.Render(m.statusMsg))
	}

	help := lipgloss.NewStyle().Foreground(subTextColor).Render("\n  " + helpText)
	return lipgloss.JoinVertical(lipgloss.Left, split, help)
}
