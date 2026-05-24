package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"jira-dashboard/config"
	"jira-dashboard/jira"
	"jira-dashboard/ui"
)

func main() {
	cfg, err := config.LoadConfig()

	var m ui.Model
	if err != nil {
		m = ui.NewErrorModel(err)
	} else {
		m = ui.NewModel(jira.NewJiraClient(cfg))
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Fatal:", err)
		os.Exit(1)
	}
}
