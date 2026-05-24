package ui

import (
	"fmt"
	"strings"
)

func (m Model) teamView() string {
	var b strings.Builder

	if m.teamLoading {
		b.WriteString("  " + m.spinner.View() + "  Loading team issues…\n")
		return b.String()
	}
	if !m.teamLoaded {
		return b.String()
	}

	visible := m.teamFilteredIssues()
	searchActive := m.state == stateTeamSearch

	hasBadge := m.teamNameFilter != "" || m.teamStatusFilter != "" || m.teamSearchQuery != "" || searchActive

	if m.teamNameFilter != "" || m.teamStatusFilter != "" {
		var badges []string
		if m.teamNameFilter != "" {
			badges = append(badges, filterBadgeStyle.Render(" "+shortName(m.teamNameFilter)+" "))
		}
		if m.teamStatusFilter != "" {
			badges = append(badges, filterBadgeStyle.Render(" "+m.teamStatusFilter+" "))
		}
		b.WriteString("  " + strings.Join(badges, "  ") + "\n")
	}
	if m.teamSearchQuery != "" {
		b.WriteString("  " + filterBadgeStyle.Render(" / "+m.teamSearchQuery+" ") + "\n")
	}
	if searchActive {
		b.WriteString("  " + m.textInput.View() + "\n")
	}
	if hasBadge {
		b.WriteString("\n")
	}

	if len(visible) == 0 {
		if len(m.teamIssues) == 0 {
			b.WriteString("  " + dimStyle.Render("No team issues found.") + "\n")
		} else {
			b.WriteString("  " + dimStyle.Render("No issues match the current filters.") + "\n")
		}
		return b.String()
	}

	col1, col2, col3, col4 := 13, 36, 16, 22
	b.WriteString(headerStyle.Render(fmt.Sprintf("  %-*s %-*s %-*s %-*s",
		col1, "KEY", col2, "SUMMARY", col3, "STATUS", col4, "ASSIGNEE")) + "\n")

	extraH := 0
	if searchActive {
		extraH = 2
	}
	maxH := m.listHeight(hasBadge) - extraH
	if maxH < 3 {
		maxH = 3
	}
	start, end := scrollWindow(m.teamCursor, len(visible), maxH)
	for i, iss := range visible[start:end] {
		b.WriteString(renderIssueRow(iss, (i+start) == m.teamCursor, col1, col2, col3, col4, "assignee") + "\n")
	}
	if len(visible) > maxH {
		b.WriteString("  " + dimStyle.Render(fmt.Sprintf("%d–%d of %d  ↑↓ scroll", start+1, end, len(visible))) + "\n")
	}
	return b.String()
}
