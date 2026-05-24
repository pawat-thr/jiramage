package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"

	"jira-dashboard/jira"
)

func uniqueSorted(items []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range items {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}

func buildStatusOptions(issues []jira.Issue) []string {
	var names []string
	for _, iss := range issues {
		names = append(names, iss.Fields.Status.Name)
	}
	return uniqueSorted(names)
}

func buildNameOptions(issues []jira.Issue) []string {
	var names []string
	for _, iss := range issues {
		if iss.Fields.Assignee != nil {
			names = append(names, iss.Fields.Assignee.DisplayName)
		}
	}
	return uniqueSorted(names)
}

func (m Model) myFilteredIssues() []jira.Issue {
	var out []jira.Issue
	q := strings.ToLower(m.mySearchQuery)
	for _, iss := range m.myIssues {
		if m.myStatusFilter != "" && iss.Fields.Status.Name != m.myStatusFilter {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(iss.Key+" "+iss.Fields.Summary), q) {
			continue
		}
		out = append(out, iss)
	}
	return out
}

func (m Model) teamFilteredIssues() []jira.Issue {
	var out []jira.Issue
	q := strings.ToLower(m.teamSearchQuery)
	for _, iss := range m.teamIssues {
		if m.teamNameFilter != "" {
			name := "Unassigned"
			if iss.Fields.Assignee != nil {
				name = iss.Fields.Assignee.DisplayName
			}
			if name != m.teamNameFilter {
				continue
			}
		}
		if m.teamStatusFilter != "" && iss.Fields.Status.Name != m.teamStatusFilter {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(iss.Key+" "+iss.Fields.Summary), q) {
			continue
		}
		out = append(out, iss)
	}
	return out
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d%time.Hour == 0 {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dm", int(d.Minutes()))
}

func (m Model) currentCursor() int {
	if m.activeTab == tabTeam {
		return m.teamCursor
	}
	return m.myCursor
}

func (m Model) currentVisible() []jira.Issue {
	if m.activeTab == tabTeam {
		return m.teamFilteredIssues()
	}
	return m.myFilteredIssues()
}

func (m Model) listHeight(hasBadge bool) int {
	if m.height == 0 {
		return 20
	}
	fixed := 7
	if hasBadge {
		fixed += 2
	}
	h := m.height - fixed
	if h < 3 {
		h = 3
	}
	return h
}

func scrollWindow(cursor, total, maxVisible int) (start, end int) {
	if total <= maxVisible {
		return 0, total
	}
	start = cursor - maxVisible/2
	if start < 0 {
		start = 0
	}
	end = start + maxVisible
	if end > total {
		end = total
		start = end - maxVisible
	}
	return start, end
}

func shortName(name string) string {
	if start := len(name) - 1; start >= 0 {
		if s := findLast(name, '('); s != -1 {
			if e := findLast(name, ')'); e > s {
				return name[s+1 : e]
			}
		}
	}
	return name
}

func findLast(s string, b byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func openURL(u string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{u}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", u}
	default:
		cmd = "xdg-open"
		args = []string{u}
	}
	exec.Command(cmd, args...).Start() //nolint:errcheck
}
