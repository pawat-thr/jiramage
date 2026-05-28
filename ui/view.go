package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"jira-dashboard/jira"
)

func (m Model) View() string {
	switch m.state {
	case stateSplash:
		return m.splashView()
	case stateLoading:
		return "\n  " + m.spinner.View() + "  Loading Jira issues…\n"
	case stateError:
		if m.client == nil {
			return m.configErrorView()
		}
		return "\n  " + errorStyle.Render("✗ "+m.err.Error()) +
			"\n\n  " + dimStyle.Render("r retry   q quit") + "\n"
	default:
		return m.mainView()
	}
}

func (m Model) configErrorView() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(red).
		Padding(1, 3)

	var s string
	s += errorStyle.Render("✗  Configuration Error") + "\n\n"
	s += dimStyle.Render(m.err.Error()) + "\n\n"
	s += titleStyle.Render("Create a .env file in the project directory:") + "\n\n"
	s += lipgloss.NewStyle().Foreground(teal).Render("  JIRA_URL=https://yourcompany.atlassian.net") + "\n"
	s += lipgloss.NewStyle().Foreground(teal).Render("  JIRA_EMAIL=you@company.com") + "\n"
	s += lipgloss.NewStyle().Foreground(teal).Render("  JIRA_TOKEN=your_api_token") + "\n"
	s += lipgloss.NewStyle().Foreground(teal).Render("  TEAM_EMAILS=a@co.com,b@co.com") + "\n\n"
	s += dimStyle.Render("Generate a token at: id.atlassian.com → Security → API tokens")

	content := box.Render(s)
	topPad := 0
	h := lipgloss.Height(content)
	if m.height > h+4 {
		topPad = (m.height - h) / 2
	}
	return strings.Repeat("\n", topPad) + "  " + content + "\n\n" +
		"  " + dimStyle.Render("q quit") + "\n"
}

func (m Model) mainView() string {
	var b strings.Builder

	// header
	left := titleStyle.Render("  Jiramage Dashboard")
	if m.projectLabel != "" {
		left += "  " + lipgloss.NewStyle().Foreground(teal).Bold(true).Render("["+m.projectLabel+"]")
	}
	ts := m.myLastUpdated
	switch m.activeTab {
	case tabTeam, tabDashboard:
		ts = m.teamLastUpdated
	case tabFixed:
		ts = m.fixedLastUpdated
	}
	right := dimStyle.Render(fmt.Sprintf("updated %s  auto-refresh %s  ", ts.Format("15:04:05"), formatDuration(m.refreshInterval)))
	pad := ""
	if m.width > 0 {
		gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
		if gap > 0 {
			pad = strings.Repeat(" ", gap)
		}
	}
	b.WriteString(lipgloss.NewStyle().Background(lipgloss.Color("#1F1035")).Render(left + pad + right))
	b.WriteString("\n")

	// tab bar
	tabs := []struct {
		label string
		tab   tabView
	}{
		{"  1  My Tasks  ", tabMyTasks},
		{"  2  Team  ", tabTeam},
		{"  3  Dashboard  ", tabDashboard},
		{"  4  Fixed  ", tabFixed},
	}
	for _, t := range tabs {
		if m.activeTab == t.tab {
			b.WriteString(tabActiveStyle.Render(t.label))
		} else {
			b.WriteString(tabInactiveStyle.Render(t.label))
		}
	}
	b.WriteString("\n\n")

	// page content
	switch m.activeTab {
	case tabMyTasks:
		b.WriteString(m.myTasksView())
	case tabTeam:
		b.WriteString(m.teamView())
	case tabDashboard:
		b.WriteString(m.dashboardView())
	case tabFixed:
		b.WriteString(m.fixedView())
	}

	// modals
	switch m.state {
	case stateMyStatusFilter:
		b.WriteString(m.statusFilterModal("Filter by Status", m.myFilterOptions, m.myFilterCursor, m.myStatusFilter))
	case stateTeamNameFilter:
		b.WriteString(m.nameFilterModal())
	case stateTeamStatusFilter:
		b.WriteString(m.statusFilterModal("Filter by Status", m.teamStatusOptions, m.teamStatusFilterCursor, m.teamStatusFilter))
	case stateReassignSearch:
		b.WriteString(m.reassignSearchView())
	case stateReassignPick:
		b.WriteString(m.reassignPickView())
	case stateTransitionPick:
		b.WriteString(m.transitionPickView())
	}

	// notification
	if m.statusMsg != "" {
		if m.statusOK {
			b.WriteString("  " + okStyle.Render(m.statusMsg) + "\n")
		} else {
			b.WriteString("  " + errorStyle.Render(m.statusMsg) + "\n")
		}
	}

	b.WriteString(m.helpBar())
	return b.String()
}

func (m Model) helpBar() string {
	sep := dimStyle.Render("  ")

	line1 := "  " + strings.Join([]string{
		keyStyle.Render("1") + dimStyle.Render(" my tasks"),
		keyStyle.Render("2") + dimStyle.Render(" team"),
		keyStyle.Render("3") + dimStyle.Render(" dashboard"),
		keyStyle.Render("4") + dimStyle.Render(" fixed"),
		keyStyle.Render("↑↓") + dimStyle.Render(" navigate"),
		keyStyle.Render("r") + dimStyle.Render(" refresh"),
		keyStyle.Render("q") + dimStyle.Render(" quit"),
	}, sep)

	hHint := dimStyle.Render(" show all")
	if !m.hideCompleted {
		hHint = dimStyle.Render(" active only")
	}

	var row2 []string
	switch m.activeTab {
	case tabMyTasks:
		row2 = []string{
			keyStyle.Render("enter") + dimStyle.Render(" open"),
			keyStyle.Render("h") + hHint,
			keyStyle.Render("f") + dimStyle.Render(" filter status"),
			keyStyle.Render("/") + dimStyle.Render(" search"),
			keyStyle.Render("t") + dimStyle.Render(" transition"),
			keyStyle.Render("a") + dimStyle.Render(" reassign"),
		}
	case tabTeam:
		row2 = []string{
			keyStyle.Render("enter") + dimStyle.Render(" open"),
			keyStyle.Render("h") + hHint,
			keyStyle.Render("f") + dimStyle.Render(" filter name"),
			keyStyle.Render("s") + dimStyle.Render(" filter status"),
			keyStyle.Render("/") + dimStyle.Render(" search"),
			keyStyle.Render("t") + dimStyle.Render(" transition"),
			keyStyle.Render("a") + dimStyle.Render(" reassign"),
		}
	case tabDashboard:
		row2 = []string{
			keyStyle.Render("enter") + dimStyle.Render(" view team"),
			keyStyle.Render("h") + hHint,
			keyStyle.Render("s") + dimStyle.Render(" filter status"),
		}
	case tabFixed:
		if m.fixedDrill != "" {
			row2 = []string{
				keyStyle.Render("enter") + dimStyle.Render(" open"),
				keyStyle.Render("esc") + dimStyle.Render(" back"),
			}
		} else {
			row2 = []string{
				keyStyle.Render("enter") + dimStyle.Render(" see cards"),
			}
		}
	}

	line2 := "  " + strings.Join(row2, sep)
	return line1 + "\n" + line2 + "\n"
}

// visualTruncate truncates s to at most maxW terminal columns, appending "…" when cut.
// Uses lipgloss.Width so Thai combining marks (0-width) are handled correctly.
func visualTruncate(s string, maxW int) string {
	if lipgloss.Width(s) <= maxW {
		return s
	}
	var out []rune
	w := 0
	for _, r := range s {
		cw := lipgloss.Width(string(r))
		if w+cw > maxW-1 {
			break
		}
		out = append(out, r)
		w += cw
	}
	return string(out) + "…"
}

// visualPad pads s with spaces to exactly w terminal columns.
func visualPad(s string, w int) string {
	sw := lipgloss.Width(s)
	if sw >= w {
		return s
	}
	return s + strings.Repeat(" ", w-sw)
}

func renderIssueRow(iss jira.Issue, selected bool, col1, col2, col3, col4 int, col4Mode string) string {
	sum := visualPad(visualTruncate(iss.Fields.Summary, col2), col2)

	statusName := iss.Fields.Status.Name
	if len(statusName) > col3 {
		statusName = statusName[:col3-1] + "…"
	}

	col4Val := iss.Fields.Priority.Name
	if col4Mode == "assignee" {
		col4Val = "Unassigned"
		if iss.Fields.Assignee != nil {
			col4Val = shortName(iss.Fields.Assignee.DisplayName)
		}
	}
	if len(col4Val) > col4 {
		col4Val = col4Val[:col4-1] + "…"
	}

	if selected {
		return selStyle.Render(fmt.Sprintf("▶ %-*s %s %-*s %-*s",
			col1, iss.Key, sum, col3, statusName, col4, col4Val))
	}

	sColor := statusColor[iss.Fields.Status.Name]
	if sColor == "" {
		sColor = grey
	}
	var c4Color lipgloss.Color
	if col4Mode == "priority" {
		c4Color = priorityColor[iss.Fields.Priority.Name]
		if c4Color == "" {
			c4Color = grey
		}
	} else {
		c4Color = teal
	}

	kStr := lipgloss.NewStyle().Foreground(purple).Bold(true).Render(fmt.Sprintf("%-*s", col1, iss.Key))
	sStr := lipgloss.NewStyle().Foreground(sColor).Render(fmt.Sprintf("%-*s", col3, statusName))
	c4Str := lipgloss.NewStyle().Foreground(c4Color).Render(fmt.Sprintf("%-*s", col4, col4Val))
	return fmt.Sprintf("  %s %s %s %s", kStr, sum, sStr, c4Str)
}
