package ui

import tea "github.com/charmbracelet/bubbletea"

func (m Model) handleMyTasksKey(key string) (Model, tea.Cmd) {
	visible := m.myFilteredIssues()
	switch key {
	case "up", "k":
		if m.myCursor > 0 {
			m.myCursor--
		}
	case "down", "j":
		if m.myCursor < len(visible)-1 {
			m.myCursor++
		}
	case "enter":
		if len(visible) > 0 {
			openURL(visible[m.myCursor].BrowseURL(m.client.BaseURL()))
		}
	case "a":
		if len(visible) > 0 {
			m.state = stateReassignSearch
			m.textInput.Reset()
			m.textInput.Focus()
			m.users = nil
		}
	case "f":
		m.state = stateMyStatusFilter
		m.myFilterCursor = 0
		for i, opt := range m.myFilterOptions {
			if opt == m.myStatusFilter {
				m.myFilterCursor = i + 1
				break
			}
		}
	case "/":
		m.state = stateMySearch
		m.textInput.Placeholder = "Search by title or key…"
		m.textInput.Reset()
		if m.mySearchQuery != "" {
			m.textInput.SetValue(m.mySearchQuery)
		}
		m.textInput.Focus()
	case "t":
		if len(visible) > 0 {
			return m, fetchTransitions(m.client, visible[m.myCursor].Key)
		}
	case "h":
		m.hideCompleted = !m.hideCompleted
		m.myCursor = 0
	case "r":
		m.state = stateLoading
		return m, tea.Batch(m.spinner.Tick, fetchMyIssues(m.client))
	}
	return m, nil
}

func (m Model) handleTeamKey(key string) (Model, tea.Cmd) {
	visible := m.teamFilteredIssues()
	switch key {
	case "up", "k":
		if m.teamCursor > 0 {
			m.teamCursor--
		}
	case "down", "j":
		if m.teamCursor < len(visible)-1 {
			m.teamCursor++
		}
	case "enter":
		if len(visible) > 0 {
			openURL(visible[m.teamCursor].BrowseURL(m.client.BaseURL()))
		}
	case "a":
		if len(visible) > 0 {
			m.state = stateReassignSearch
			m.textInput.Reset()
			m.textInput.Focus()
			m.users = nil
		}
	case "f":
		m.state = stateTeamNameFilter
		m.teamNameFilterCursor = 0
		for i, opt := range m.teamNameOptions {
			if opt == m.teamNameFilter {
				m.teamNameFilterCursor = i + 1
				break
			}
		}
	case "s":
		m.state = stateTeamStatusFilter
		m.teamStatusFilterCursor = 0
		for i, opt := range m.teamStatusOptions {
			if opt == m.teamStatusFilter {
				m.teamStatusFilterCursor = i + 1
				break
			}
		}
	case "/":
		m.state = stateTeamSearch
		m.textInput.Placeholder = "Search by title or key…"
		m.textInput.Reset()
		if m.teamSearchQuery != "" {
			m.textInput.SetValue(m.teamSearchQuery)
		}
		m.textInput.Focus()
	case "t":
		if len(visible) > 0 {
			return m, fetchTransitions(m.client, visible[m.teamCursor].Key)
		}
	case "h":
		m.hideCompleted = !m.hideCompleted
		m.teamCursor = 0
	case "r":
		m.teamLoading = true
		return m, fetchTeamIssues(m.client)
	}
	return m, nil
}

func (m Model) handleDashboardKey(key string) (Model, tea.Cmd) {
	total := len(m.teamNameOptions) + 1
	switch key {
	case "up", "k":
		if m.dashboardCursor > 0 {
			m.dashboardCursor--
		}
	case "down", "j":
		if m.dashboardCursor < total-1 {
			m.dashboardCursor++
		}
	case "enter":
		if m.dashboardCursor == 0 {
			m.teamNameFilter = ""
		} else {
			m.teamNameFilter = m.teamNameOptions[m.dashboardCursor-1]
		}
		m.teamCursor = 0
		m.activeTab = tabTeam
	case "s":
		m.state = stateTeamStatusFilter
		m.teamStatusFilterCursor = 0
		for i, opt := range m.teamStatusOptions {
			if opt == m.teamStatusFilter {
				m.teamStatusFilterCursor = i + 1
				break
			}
		}
	case "h":
		m.hideCompleted = !m.hideCompleted
		m.dashboardCursor = 0
	case "r":
		m.teamLoading = true
		return m, fetchTeamIssues(m.client)
	}
	return m, nil
}
