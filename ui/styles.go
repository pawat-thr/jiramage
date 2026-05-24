package ui

import "github.com/charmbracelet/lipgloss"

var (
	purple  = lipgloss.Color("#7C3AED")
	white   = lipgloss.Color("#FFFFFF")
	grey    = lipgloss.Color("#9CA3AF")
	green   = lipgloss.Color("#10B981")
	red     = lipgloss.Color("#EF4444")
	blue    = lipgloss.Color("#3B82F6")
	orange  = lipgloss.Color("#F97316")
	yellow  = lipgloss.Color("#F59E0B")
	dimGrey = lipgloss.Color("#6B7280")
	teal    = lipgloss.Color("#0EA5E9")

	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(purple)
	dimStyle   = lipgloss.NewStyle().Foreground(grey)
	errorStyle = lipgloss.NewStyle().Foreground(red)
	okStyle    = lipgloss.NewStyle().Foreground(green)
	keyStyle   = lipgloss.NewStyle().Foreground(green).Bold(true)

	selStyle = lipgloss.NewStyle().
			Background(purple).
			Foreground(white).
			Bold(true)

	tabActiveStyle = lipgloss.NewStyle().
			Background(purple).
			Foreground(white).
			Bold(true).
			Padding(0, 2)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(grey).
				Padding(0, 2)

	filterBadgeStyle = lipgloss.NewStyle().
				Background(green).
				Foreground(white).
				Bold(true).
				Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Foreground(grey).
			Underline(true)

	modalBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(0, 2)

	priorityColor = map[string]lipgloss.Color{
		"Highest": red,
		"High":    orange,
		"Medium":  yellow,
		"Low":     blue,
		"Lowest":  dimGrey,
	}

	statusColor = map[string]lipgloss.Color{
		"To Do":              dimGrey,
		"In Progress":        blue,
		"In Review":          purple,
		"Done":               green,
		"REJECTED":           red,
		"READY TO TEST":      teal,
		"Tier 1 in progress": orange,
		"TIER 2":             yellow,
		"Need More Info":     yellow,
	}

	memberColors = []lipgloss.Color{
		lipgloss.Color("#7C3AED"),
		lipgloss.Color("#0EA5E9"),
		lipgloss.Color("#10B981"),
		lipgloss.Color("#F97316"),
		lipgloss.Color("#EF4444"),
		lipgloss.Color("#F59E0B"),
	}
)
