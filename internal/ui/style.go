package ui

import "github.com/charmbracelet/lipgloss"

var (
	sectionStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	subsectionStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("111"))
	infoStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	debugStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	successStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	activeStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	dimStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	summaryTableStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
)
