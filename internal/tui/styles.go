package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors (Dracula inspired for a cool hacker vibe)
	primaryColor   = lipgloss.Color("#bd93f9")
	secondaryColor = lipgloss.Color("#ff79c6")
	textColor      = lipgloss.Color("#f8f8f2")
	subTextColor   = lipgloss.Color("#6272a4")
	borderColor    = lipgloss.Color("#44475a")
	activeBorder   = lipgloss.Color("#8be9fd")
	errorColor     = lipgloss.Color("#ff5555")

	// Global Styles
	appStyle = lipgloss.NewStyle().Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#282a36")).
			Background(primaryColor).
			Padding(0, 2).
			MarginBottom(1)

	basePaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1)

	activePaneStyle = basePaneStyle.
			BorderForeground(activeBorder)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	activeTextStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	
	// Tab Styles
	activeTabStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, true, false, true).
			BorderForeground(activeBorder).
			Foreground(activeBorder).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, true, false, true).
			BorderForeground(borderColor).
			Foreground(subTextColor).
			Padding(0, 1)

	tabBarStyle = lipgloss.NewStyle().
			Padding(0, 0).
			Margin(0, 0)
)

