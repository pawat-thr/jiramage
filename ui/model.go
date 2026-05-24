package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"jira-dashboard/jira"
)

type viewState int
type tabView int

const (
	appName    = "jiramage"
	appVersion = "v0.1.0"
	appCredit  = "by MpLab"
)

const (
	stateSplash viewState = iota
	stateLoading
	stateMain
	stateMyStatusFilter
	stateTeamNameFilter
	stateTeamStatusFilter
	stateMySearch
	stateTeamSearch
	stateTransitionPick
	stateReassignSearch
	stateReassignPick
	stateError
)

const (
	tabMyTasks  tabView = iota
	tabTeam
	tabDashboard
)

type Model struct {
	client    *jira.JiraClient
	state     viewState
	activeTab tabView

	// page 1 – My Tasks
	myIssues        []jira.Issue
	myCursor        int
	myLastUpdated   time.Time
	myStatusFilter  string
	myFilterOptions []string
	myFilterCursor  int

	// page 2 – Team list
	teamIssues             []jira.Issue
	teamCursor             int
	teamLastUpdated        time.Time
	teamLoaded             bool
	teamLoading            bool
	teamNameFilter         string
	teamNameOptions        []string
	teamNameFilterCursor   int
	teamStatusFilter       string
	teamStatusOptions      []string
	teamStatusFilterCursor int

	// page 3 – Dashboard
	dashboardCursor int // 0 = All, 1..N = member index

	// search
	mySearchQuery   string
	teamSearchQuery string

	// transitions
	transitions      []jira.Transition
	transitionCursor int

	// refresh
	refreshInterval time.Duration

	// shared
	err       error
	spinner   spinner.Model
	textInput textinput.Model
	users     []jira.User
	userCursor int
	statusMsg  string
	statusOK   bool
	width      int
	height     int
}

func NewModel(client *jira.JiraClient) Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(purple)

	ti := textinput.New()
	ti.Placeholder = "Type name or email..."
	ti.CharLimit = 80
	ti.Width = 36

	return Model{client: client, state: stateSplash, spinner: sp, textInput: ti, refreshInterval: client.RefreshInterval()}
}

// NewErrorModel creates a model that shows a config error immediately on startup.
func NewErrorModel(err error) Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(purple)
	return Model{state: stateError, err: err, spinner: sp}
}

func (m Model) Init() tea.Cmd {
	if m.client == nil {
		return nil // config error — nothing to fetch
	}
	return tea.Batch(m.spinner.Tick, fetchMyIssues(m.client), splashTimer())
}
