package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) fixedView() string {
	var b strings.Builder

	if !m.client.HasLabels() {
		b.WriteString("\n  " + dimStyle.Render("Set JIRA_LABELS=R6.1-SIT in .env to enable this tab.") + "\n")
		return b.String()
	}

	if m.fixedLoading {
		b.WriteString("  " + m.spinner.View() + "  Loading fixed issues…\n")
		return b.String()
	}
	if !m.fixedLoaded {
		return b.String()
	}

	if m.fixedDrill != "" {
		return m.fixedDrillView()
	}
	return m.fixedSummaryView()
}

func (m Model) fixedSummaryView() string {
	var b strings.Builder

	members := m.client.AllMembers()
	labelStr := m.client.LabelLabel()
	statusStr := m.client.FixedStatusName()

	b.WriteString("  " + titleStyle.Render("Fixed by Team") +
		"  " + modeBadgeStyle.Render(" "+labelStr+" ") +
		"  " + dimStyle.Render("→ "+statusStr) + "\n\n")

	total := 0
	counts := make(map[string]int)
	for _, email := range members {
		n := len(m.fixedByMember[email])
		counts[email] = n
		total += n
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

	allSel := m.fixedCursor == 0
	allBar := strings.Repeat("█", barW)
	if allSel {
		b.WriteString(selStyle.Render(fmt.Sprintf("▶ %-14s [%-*s] %4d", "All", barW, allBar, total)) + "\n\n")
	} else {
		allBarStr := lipgloss.NewStyle().Foreground(white).Render(allBar)
		b.WriteString(fmt.Sprintf("  %-14s [%s] %4d\n\n", "All", allBarStr, total))
	}

	for i, email := range members {
		count := counts[email]
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
		nick := emailUsername(email)
		sel := m.fixedCursor == i+1

		if sel {
			b.WriteString(selStyle.Render(fmt.Sprintf("▶ %-14s [%-*s] %4d  %3d%%", nick, barW, bar, count, pct)) + "\n")
		} else {
			nameStr := lipgloss.NewStyle().Foreground(color).Bold(true).Render(fmt.Sprintf("%-14s", nick))
			barStr := lipgloss.NewStyle().Foreground(color).Render(bar)
			b.WriteString(fmt.Sprintf("  %s [%s] %4d  %3d%%\n", nameStr, barStr, count, pct))
		}
	}

	b.WriteString("\n" + sepLine + "\n")
	b.WriteString("  " + dimStyle.Render(fmt.Sprintf("total %d cards fixed  |  enter → see cards", total)) + "\n")
	return b.String()
}

func (m Model) fixedDrillView() string {
	var b strings.Builder

	nick := emailUsername(m.fixedDrill)
	labelStr := m.client.LabelLabel()
	issues := m.fixedByMember[m.fixedDrill]

	b.WriteString("  " + titleStyle.Render("Fixed by "+nick) +
		"  " + modeBadgeStyle.Render(" "+labelStr+" ") + "\n\n")

	if len(issues) == 0 {
		b.WriteString("  " + dimStyle.Render("No fixed issues found.") + "\n")
		return b.String()
	}

	col1 := 12
	col3 := 16
	col4 := 22
	col2 := m.width - col1 - col3 - col4 - 6
	if col2 < 20 {
		col2 = 20
	}
	sepW := col1 + col2 + col3 + col4 + 3

	b.WriteString(dimStyle.Render(fmt.Sprintf(
		"  %-*s %-*s %-*s %-*s", col1, "KEY", col2, "SUMMARY", col3, "STATUS", col4, "ASSIGNEE")) + "\n")
	b.WriteString("  " + dimStyle.Render(strings.Repeat("─", sepW)) + "\n")

	listH := m.listHeight(false) - 1
	start, end := scrollWindow(m.fixedDrillCursor, len(issues), listH)

	for i, iss := range issues[start:end] {
		idx := start + i
		b.WriteString(renderIssueRow(iss, idx == m.fixedDrillCursor, col1, col2, col3, col4, "assignee") + "\n")
	}

	if len(issues) > listH {
		b.WriteString("  " + dimStyle.Render(fmt.Sprintf("(%d-%d of %d)  ↑↓ scroll", start+1, end, len(issues))) + "\n")
	}

	return b.String()
}
