package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) dashboardView() string {
	var b strings.Builder

	if m.teamLoading {
		b.WriteString("  " + m.spinner.View() + "  Loading team issues…\n")
		return b.String()
	}
	if !m.teamLoaded {
		return b.String()
	}

	b.WriteString("  " + titleStyle.Render("Team Performance Meter") + "\n")

	if m.teamStatusFilter != "" {
		b.WriteString("  " + filterBadgeStyle.Render(" "+m.teamStatusFilter+" ") + "  ")
		b.WriteString(dimStyle.Render("(counts filtered by status)") + "\n")
	}
	b.WriteString("\n")

	counts := map[string]int{}
	total := 0
	for _, iss := range m.teamIssues {
		if m.teamStatusFilter != "" && iss.Fields.Status.Name != m.teamStatusFilter {
			continue
		}
		name := "Unassigned"
		if iss.Fields.Assignee != nil {
			name = iss.Fields.Assignee.DisplayName
		}
		counts[name]++
		total++
	}

	maxCount := 1
	for _, c := range counts {
		if c > maxCount {
			maxCount = c
		}
	}

	barW := m.width - 38
	if barW < 10 {
		barW = 10
	}
	if barW > 52 {
		barW = 52
	}

	sepLine := "  " + dimStyle.Render(strings.Repeat("─", barW+24))
	b.WriteString(sepLine + "\n\n")

	allSel := m.dashboardCursor == 0
	allBar := strings.Repeat("█", barW)
	if allSel {
		b.WriteString(selStyle.Render(fmt.Sprintf("▶ %-14s [%-*s] %4d  %3d%%", "All", barW, allBar, total, 100)) + "\n\n")
	} else {
		allBarStr := lipgloss.NewStyle().Foreground(white).Render(allBar)
		b.WriteString(fmt.Sprintf("  %-14s [%s] %4d  %3d%%\n\n", "All", allBarStr, total, 100))
	}

	for i, opt := range m.teamNameOptions {
		count := counts[opt]
		pct := 0
		if total > 0 {
			pct = count * 100 / total
		}
		filled := 0
		if maxCount > 0 {
			filled = count * barW / maxCount
		}
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barW-filled)

		color := memberColors[i%len(memberColors)]
		nick := shortName(opt)
		sel := m.dashboardCursor == i+1

		if sel {
			b.WriteString(selStyle.Render(fmt.Sprintf("▶ %-14s [%-*s] %4d  %3d%%", nick, barW, bar, count, pct)) + "\n")
		} else {
			nameStr := lipgloss.NewStyle().Foreground(color).Bold(true).Render(fmt.Sprintf("%-14s", nick))
			barStr := lipgloss.NewStyle().Foreground(color).Render(bar)
			b.WriteString(fmt.Sprintf("  %s [%s] %4d  %3d%%\n", nameStr, barStr, count, pct))
		}
	}

	b.WriteString("\n" + sepLine + "\n")
	b.WriteString("  " + dimStyle.Render(fmt.Sprintf("total %d tasks  |  enter → open member tasks in page 2", total)) + "\n")
	return b.String()
}
