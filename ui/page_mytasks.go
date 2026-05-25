package ui

import (
	"fmt"
	"strings"
)

func (m Model) myTasksView() string {
	var b strings.Builder
	visible := m.myFilteredIssues()
	searchActive := m.state == stateMySearch

	hasBadge := m.hideCompleted || m.myStatusFilter != "" || m.mySearchQuery != "" || searchActive
	if m.hideCompleted {
		b.WriteString("  " + modeBadgeStyle.Render(" ● active only ") + "\n")
	}
	if m.myStatusFilter != "" {
		b.WriteString("  " + filterBadgeStyle.Render(" Status: "+m.myStatusFilter+" ") + "\n")
	}
	if m.mySearchQuery != "" {
		b.WriteString("  " + filterBadgeStyle.Render(" / "+m.mySearchQuery+" ") + "\n")
	}
	if searchActive {
		b.WriteString("  " + m.textInput.View() + "\n")
	}
	if hasBadge {
		b.WriteString("\n")
	}

	if len(visible) == 0 {
		msg := "No issues assigned to you."
		if hasBadge {
			msg = "No issues match the current filters."
		}
		b.WriteString("  " + dimStyle.Render(msg) + "\n")
		return b.String()
	}

	col1, col2, col3, col4 := 13, 44, 16, 11
	b.WriteString(headerStyle.Render(fmt.Sprintf("  %-*s %-*s %-*s %-*s",
		col1, "KEY", col2, "SUMMARY", col3, "STATUS", col4, "PRIORITY")) + "\n")

	extraH := 0
	if searchActive {
		extraH = 2
	}
	maxH := m.listHeight(hasBadge) - extraH
	if maxH < 3 {
		maxH = 3
	}
	start, end := scrollWindow(m.myCursor, len(visible), maxH)
	for i, iss := range visible[start:end] {
		b.WriteString(renderIssueRow(iss, (i+start) == m.myCursor, col1, col2, col3, col4, "priority") + "\n")
	}
	if len(visible) > maxH {
		b.WriteString("  " + dimStyle.Render(fmt.Sprintf("%d–%d of %d  ↑↓ scroll", start+1, end, len(visible))) + "\n")
	}
	return b.String()
}
