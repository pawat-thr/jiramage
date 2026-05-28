package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		return m.handleKey(msg)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case issuesMsg:
		m.myIssues = msg.issues
		m.myLastUpdated = time.Now()
		m.myFilterOptions = buildStatusOptions(m.myIssues)
		if m.state == stateLoading {
			m.state = stateMain
		}
		m.err = nil
		if m.myCursor >= len(m.myFilteredIssues()) {
			m.myCursor = max(0, len(m.myFilteredIssues())-1)
		}
		return m, myTimerCmd(m.refreshInterval)

	case teamIssuesMsg:
		m.teamIssues = msg.issues
		m.teamLastUpdated = time.Now()
		m.teamLoaded = true
		m.teamLoading = false
		m.teamNameOptions = buildNameOptions(m.teamIssues)
		m.teamStatusOptions = buildStatusOptions(m.teamIssues)
		m.dashboardCursor = 0
		for i, opt := range m.teamNameOptions {
			if opt == m.teamNameFilter {
				m.dashboardCursor = i + 1
				break
			}
		}
		if m.teamCursor >= len(m.teamFilteredIssues()) {
			m.teamCursor = max(0, len(m.teamFilteredIssues())-1)
		}
		return m, teamTimerCmd(m.refreshInterval)

	case fixedIssuesMsg:
		m.fixedByMember = msg.byMember
		m.fixedLastUpdated = time.Now()
		m.fixedLoaded = true
		m.fixedLoading = false
		return m, fixedTimerCmd(m.refreshInterval)

	case transitionsMsg:
		m.transitions = msg.transitions
		m.transitionCursor = 0
		m.state = stateTransitionPick

	case transitionedMsg:
		m.statusMsg = fmt.Sprintf("✓ %s moved", msg.issueKey)
		m.statusOK = true
		m.state = stateMain
		if m.activeTab == tabTeam {
			return m, fetchTeamIssues(m.client)
		}
		return m, fetchMyIssues(m.client)

	case usersMsg:
		m.users = msg.users
		m.userCursor = 0
		m.state = stateReassignPick

	case assignedMsg:
		m.statusMsg = fmt.Sprintf("✓ %s reassigned", msg.issueKey)
		m.statusOK = true
		m.state = stateMain
		if m.activeTab == tabTeam {
			return m, fetchTeamIssues(m.client)
		}
		return m, fetchMyIssues(m.client)

	case splashDoneMsg:
		if m.state == stateSplash {
			if len(m.myIssues) > 0 {
				m.state = stateMain
			} else {
				m.state = stateLoading
			}
		}

	case tickMsg:
		return m, fetchMyIssues(m.client)

	case teamTickMsg:
		return m, fetchTeamIssues(m.client)

	case fixedTickMsg:
		return m, fetchAllFixed(m.client)

	case errMsg:
		m.teamLoading = false
		m.fixedLoading = false
		if m.state == stateLoading {
			m.state = stateError
			m.err = msg.err
		} else {
			m.statusMsg = "Error: " + msg.err.Error()
			m.statusOK = false
			m.state = stateMain
		}
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch m.state {

	case stateSplash:
		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ":
			if len(m.myIssues) > 0 {
				m.state = stateMain
			} else {
				m.state = stateLoading
			}
		}
		return m, nil

	case stateMain:
		m.statusMsg = ""
		switch key {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.activeTab = tabMyTasks
			return m, nil
		case "2":
			m.activeTab = tabTeam
			if !m.teamLoaded && !m.teamLoading {
				m.teamLoading = true
				return m, fetchTeamIssues(m.client)
			}
			return m, nil
		case "3":
			m.activeTab = tabDashboard
			if !m.teamLoaded && !m.teamLoading {
				m.teamLoading = true
				return m, fetchTeamIssues(m.client)
			}
			return m, nil
		case "4":
			m.activeTab = tabFixed
			if !m.fixedLoaded && !m.fixedLoading {
				m.fixedLoading = true
				return m, fetchAllFixed(m.client)
			}
			return m, nil
		}
		switch m.activeTab {
		case tabMyTasks:
			return m.handleMyTasksKey(key)
		case tabTeam:
			return m.handleTeamKey(key)
		case tabDashboard:
			return m.handleDashboardKey(key)
		case tabFixed:
			return m.handleFixedKey(key)
		}

	case stateMyStatusFilter:
		total := len(m.myFilterOptions) + 1
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = stateMain
		case "up", "k":
			if m.myFilterCursor > 0 {
				m.myFilterCursor--
			}
		case "down", "j":
			if m.myFilterCursor < total-1 {
				m.myFilterCursor++
			}
		case "enter":
			if m.myFilterCursor == 0 {
				m.myStatusFilter = ""
			} else {
				m.myStatusFilter = m.myFilterOptions[m.myFilterCursor-1]
			}
			m.myCursor = 0
			m.state = stateMain
		}

	case stateTeamNameFilter:
		total := len(m.teamNameOptions) + 1
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = stateMain
		case "up", "k":
			if m.teamNameFilterCursor > 0 {
				m.teamNameFilterCursor--
			}
		case "down", "j":
			if m.teamNameFilterCursor < total-1 {
				m.teamNameFilterCursor++
			}
		case "enter":
			if m.teamNameFilterCursor == 0 {
				m.teamNameFilter = ""
			} else {
				m.teamNameFilter = m.teamNameOptions[m.teamNameFilterCursor-1]
			}
			m.teamCursor = 0
			m.state = stateMain
		}

	case stateTeamStatusFilter:
		total := len(m.teamStatusOptions) + 1
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = stateMain
		case "up", "k":
			if m.teamStatusFilterCursor > 0 {
				m.teamStatusFilterCursor--
			}
		case "down", "j":
			if m.teamStatusFilterCursor < total-1 {
				m.teamStatusFilterCursor++
			}
		case "enter":
			if m.teamStatusFilterCursor == 0 {
				m.teamStatusFilter = ""
			} else {
				m.teamStatusFilter = m.teamStatusOptions[m.teamStatusFilterCursor-1]
			}
			m.teamCursor = 0
			m.state = stateMain
		}

	case stateMySearch:
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.mySearchQuery = ""
			m.textInput.Reset()
			m.state = stateMain
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			m.mySearchQuery = m.textInput.Value()
			m.myCursor = 0
			return m, cmd
		}

	case stateTeamSearch:
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.teamSearchQuery = ""
			m.textInput.Reset()
			m.state = stateMain
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			m.teamSearchQuery = m.textInput.Value()
			m.teamCursor = 0
			return m, cmd
		}

	case stateTransitionPick:
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = stateMain
		case "up", "k":
			if m.transitionCursor > 0 {
				m.transitionCursor--
			}
		case "down", "j":
			if m.transitionCursor < len(m.transitions)-1 {
				m.transitionCursor++
			}
		case "enter":
			visible := m.currentVisible()
			if len(m.transitions) > 0 && len(visible) > 0 {
				iss := visible[m.currentCursor()]
				t := m.transitions[m.transitionCursor]
				m.state = stateMain
				return m, doTransition(m.client, iss.Key, t.ID)
			}
		}

	case stateReassignSearch:
		visible := m.currentVisible()
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = stateMain
		case "enter":
			q := strings.TrimSpace(m.textInput.Value())
			if q != "" && len(visible) > 0 {
				return m, searchUsers(m.client, q)
			}
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

	case stateReassignPick:
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = stateMain
		case "up", "k":
			if m.userCursor > 0 {
				m.userCursor--
			}
		case "down", "j":
			if m.userCursor < len(m.users)-1 {
				m.userCursor++
			}
		case "enter":
			visible := m.currentVisible()
			if len(m.users) > 0 && len(visible) > 0 {
				iss := visible[m.currentCursor()]
				m.state = stateMain
				return m, assignIssue(m.client, iss.Key, m.users[m.userCursor].AccountID)
			}
		}

	case stateError:
		switch key {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			if m.client != nil {
				m.state = stateLoading
				return m, tea.Batch(m.spinner.Tick, fetchMyIssues(m.client))
			}
		}
	}

	return m, nil
}
