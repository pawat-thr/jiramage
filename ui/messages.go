package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"jira-dashboard/jira"
)

type issuesMsg struct{ issues []jira.Issue }
type teamIssuesMsg struct{ issues []jira.Issue }
type usersMsg struct{ users []jira.User }
type assignedMsg struct{ issueKey string }
type tickMsg time.Time
type teamTickMsg time.Time
type splashDoneMsg struct{}
type errMsg struct{ err error }
type transitionsMsg struct {
	transitions []jira.Transition
	issueKey    string
}
type transitionedMsg struct{ issueKey string }

func fetchMyIssues(c *jira.JiraClient) tea.Cmd {
	return func() tea.Msg {
		issues, err := c.GetMyIssues()
		if err != nil {
			return errMsg{err}
		}
		return issuesMsg{issues}
	}
}

func fetchTeamIssues(c *jira.JiraClient) tea.Cmd {
	return func() tea.Msg {
		issues, err := c.GetTeamIssues()
		if err != nil {
			return errMsg{err}
		}
		return teamIssuesMsg{issues}
	}
}

func searchUsers(c *jira.JiraClient, q string) tea.Cmd {
	return func() tea.Msg {
		users, err := c.SearchUsers(q)
		if err != nil {
			return errMsg{err}
		}
		return usersMsg{users}
	}
}

func assignIssue(c *jira.JiraClient, issueKey, accountID string) tea.Cmd {
	return func() tea.Msg {
		if err := c.AssignIssue(issueKey, accountID); err != nil {
			return errMsg{err}
		}
		return assignedMsg{issueKey}
	}
}

func myTimerCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func teamTimerCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return teamTickMsg(t) })
}

func fetchTransitions(c *jira.JiraClient, issueKey string) tea.Cmd {
	return func() tea.Msg {
		ts, err := c.GetTransitions(issueKey)
		if err != nil {
			return errMsg{err}
		}
		return transitionsMsg{ts, issueKey}
	}
}

func doTransition(c *jira.JiraClient, issueKey, transitionID string) tea.Cmd {
	return func() tea.Msg {
		if err := c.DoTransition(issueKey, transitionID); err != nil {
			return errMsg{err}
		}
		return transitionedMsg{issueKey}
	}
}

func splashTimer() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return splashDoneMsg{} })
}
