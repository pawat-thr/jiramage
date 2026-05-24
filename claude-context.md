# Jiramage — Claude Context File

This file gives a new Claude instance full context about this project: what it is, every decision made, the current state of the code, known issues, and what to work on next. Read this before touching anything.

---

## What This Project Is

**Jiramage** (from "Jira Management") is a terminal UI (TUI) dashboard for Jira Cloud. It lets you browse your tasks, monitor your team's workload, reassign issues, and transition issue statuses — all without leaving the terminal.

- **Language**: Go 1.26
- **TUI Framework**: Bubbletea v1.2.4 (Elm architecture: Model / Update / View)
- **Styling**: charmbracelet/lipgloss v1.0.0
- **Widgets**: charmbracelet/bubbles v0.20.0 (spinner, textinput)
- **Binary name**: `jiramage`
- **Binary output location**: `build/jiramage`
- **Module name** (go.mod): `jira-dashboard` (legacy name, do NOT rename — it would break all imports)
- **Go toolchain**: 1.26.2 (darwin/arm64) — `go.mod` declares `go 1.26`
- **Version**: v0.1.0 (hardcoded in `ui/model.go` as `appVersion`)
- **Credit**: "by MpLab" (hardcoded as `appCredit`)

---

## Owner / User Context

- Developer: Pawat T. (`pawat.t@orbitdigital.co.th`)
- Company: Orbit Digital
- Jira instance: `https://orbitdigital.atlassian.net`
- Team emails in `.env`: jirayu.j, rawipas.t, tanawat.k, peradone.c, taninchot.p — all at orbitdigital.co.th
- The app is intended to be **org-agnostic** — just changing `.env` should make it work for any Jira Cloud org

---

## Project Structure

```
jiramage/                   ← project root
├── main.go                 ← entry point
├── go.mod                  ← module: jira-dashboard, go 1.26
├── go.sum
├── build/
│   └── jiramage            ← compiled binary (run with ./build/jiramage)
├── .env                    ← real credentials, never commit
├── .env.example            ← template for new users
├── README.md
├── claude-context.md       ← this file
├── config/
│   └── config.go           ← Config struct + LoadConfig()
├── jira/
│   ├── types.go            ← Issue, User, Transition types
│   ├── client.go           ← JiraClient, HTTP helpers (get/post/put)
│   └── api.go              ← GetMyIssues, GetTeamIssues, SearchUsers, AssignIssue, GetTransitions, DoTransition
└── ui/
    ├── model.go            ← Model struct, NewModel, NewErrorModel, Init
    ├── messages.go         ← all tea.Msg types + tea.Cmd factories
    ├── update.go           ← Update() — the main state machine handler
    ├── view.go             ← View() dispatcher + mainView + helpBar + renderIssueRow
    ├── splash.go           ← buildLogo() + splashView()
    ├── styles.go           ← all lipgloss colors and styles
    ├── helpers.go          ← filter helpers, formatDuration, scrollWindow, openURL, etc.
    ├── keys.go             ← handleMyTasksKey, handleTeamKey, handleDashboardKey
    ├── modals.go           ← statusFilterModal, nameFilterModal, reassignSearchView, reassignPickView, transitionPickView
    ├── page_mytasks.go     ← myTasksView()
    ├── page_team.go        ← teamView()
    └── page_dashboard.go   ← dashboardView()
```

---

## Package Dependency Graph

```
config  ←  jira  ←  ui  ←  main
```

`config` has no internal imports. `jira` imports `config`. `ui` imports `jira` and `config` (via jira). Never import `ui` from `jira` or `config` — that creates a cycle.

---

## Configuration (.env)

```env
JIRA_URL=https://yourcompany.atlassian.net
JIRA_EMAIL=you@company.com
JIRA_TOKEN=your_api_token_here
TEAM_EMAILS=a@company.com,b@company.com,c@company.com
REFRESH_INTERVAL=5m
```

- `JIRA_URL`, `JIRA_EMAIL`, `JIRA_TOKEN` are required — app shows error screen if missing
- `TEAM_EMAILS` is optional — comma-separated; used for Team tab and Dashboard tab
- `REFRESH_INTERVAL` is optional — Go duration string (e.g. `5m`, `30s`, `1h`), default `5m`
- `.env` is parsed by `config/config.go:loadDotEnv()` — no `export` prefix, supports `#` comments, no library dependency
- Real env vars take precedence over `.env` values (only sets if `os.Getenv(k) == ""`)

---

## Jira API Details

**Auth**: HTTP Basic Auth — `base64(email:token)` in `Authorization` header

**Endpoints used**:

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/rest/api/3/search/jql` | Paginated issue search |
| GET | `/rest/api/3/user/search?query=...&maxResults=10` | Search users for reassign |
| PUT | `/rest/api/3/issue/{key}/assignee` | Reassign issue |
| GET | `/rest/api/3/issue/{key}/transitions` | List available status transitions |
| POST | `/rest/api/3/issue/{key}/transitions` | Apply a status transition |

**Pagination** (`jira/api.go:searchAll`):
- Uses `POST /rest/api/3/search/jql` (NOT the old `GET /rest/api/2/search`)
- Page size: 100 per request
- Uses `nextPageToken` cursor field in response (NOT `startAt` — that's the old API and causes HTTP 400)
- Loops until `nextPageToken == ""` or `len(issues) == 0`

**JQL queries**:
- My Issues: `assignee = currentUser() ORDER BY key DESC`
- Team Issues: `assignee in ("email1","email2",...) ORDER BY key DESC` — includes the logged-in user's email too

**Fields fetched**: `summary`, `status`, `priority`, `assignee`

---

## State Machine

All view states are defined in `ui/model.go` as `viewState` iota:

```
stateSplash          → shown for 2s on startup (or press Enter to skip)
stateLoading         → spinner while fetching My Issues on first load
stateMain            → normal operation (tab-based)
stateMyStatusFilter  → status filter modal open (My Tasks tab)
stateTeamNameFilter  → name filter modal open (Team tab)
stateTeamStatusFilter→ status filter modal open (Team tab)
stateMySearch        → text input active for My Tasks search
stateTeamSearch      → text input active for Team search
stateTransitionPick  → transition picker modal open
stateReassignSearch  → text input for user search (reassign flow step 1)
stateReassignPick    → user picker modal (reassign flow step 2)
stateError           → config error (client == nil) OR runtime fetch error
```

Tabs (`tabView` iota): `tabMyTasks`, `tabTeam`, `tabDashboard`

---

## All Keyboard Shortcuts

### Global
| Key | Action |
|-----|--------|
| `q` / `ctrl+c` | Quit |
| `1` | Switch to My Tasks tab |
| `2` | Switch to Team tab |
| `3` | Switch to Dashboard tab |
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `esc` | Cancel / close modal / clear search |

### My Tasks tab
| Key | Action |
|-----|--------|
| `enter` | Open issue in browser |
| `f` | Open status filter modal |
| `/` | Open search input |
| `t` | Fetch transitions for selected issue → open transition picker |
| `a` | Open reassign search for selected issue |
| `r` | Refresh my issues |

### Team tab
| Key | Action |
|-----|--------|
| `enter` | Open issue in browser |
| `f` | Open name filter modal |
| `s` | Open status filter modal |
| `/` | Open search input |
| `t` | Fetch transitions for selected issue |
| `a` | Open reassign search for selected issue |
| `r` | Refresh team issues |

### Dashboard tab
| Key | Action |
|-----|--------|
| `enter` | Navigate to Team tab filtered to selected member |
| `s` | Open status filter modal |
| `r` | Refresh team issues |

### Error screen (config error)
| Key | Action |
|-----|--------|
| `q` | Quit (only key that works here) |

### Runtime error screen
| Key | Action |
|-----|--------|
| `r` | Retry fetching data |
| `q` | Quit |

---

## Features in Detail

### Splash Screen (`ui/splash.go`)
- Shows on startup for 2 seconds, then auto-advances to loading
- Press `enter` or `space` to skip immediately
- Displays ASCII block-letter logo spelling **JIRAMAGE** using Unicode box-drawing chars
- Each letter (J, I, R, A, M, G, E) is a 6-row `[6]string` defined in `buildLogo()` char map
- Below logo: subtitle, version (`v0.1.0`), credit (`by MpLab`), hint text
- Vertically centered based on terminal height

### My Tasks (`ui/page_mytasks.go`)
- Columns: KEY | SUMMARY | STATUS | PRIORITY
- Color-coded priority: Highest=red, High=orange, Medium=yellow, Low=blue, Lowest=dimGrey
- Color-coded status (see styles section below)
- Scrollable list with scroll indicator when list > visible height
- Active filter badges shown above list
- Search input shown inline when `stateMySearch` is active
- Filtering: AND logic — status filter AND search query both apply

### Team (`ui/page_team.go`)
- Columns: KEY | SUMMARY | STATUS | ASSIGNEE
- Assignee displayed as short name extracted from parentheses: `"John Smith (johnsmith)"` → `"johnsmith"`
  - Implemented in `helpers.go:shortName()` using `findLast()` byte scanner
- Loading spinner on first load (lazy — only fetches when tab first visited)
- Dual filter badges: name badge + status badge
- Filtering: AND logic — name filter AND status filter AND search query

### Dashboard (`ui/page_dashboard.go`)
- Shows one row per team member + "All" aggregate row
- Horizontal bar chart: bar width proportional to count, max 52 chars, min 10 chars, scales to terminal width
- Member colors cycle through: purple → teal → green → orange → red → yellow (`memberColors` slice in styles.go)
- `enter` on a member row switches to Team tab and sets name filter to that member
- `s` opens status filter for dashboard-level status filtering (filters what counts in the bars)
- Total task count shown at bottom

### Search (`/` key)
- Opens inline text input (placeholder: "Search by title or key…")
- Case-insensitive substring match on `key + " " + summary`
- Filters in real-time as user types (no API call — filters already-fetched slice)
- `esc` clears query and exits search mode
- Previous query pre-filled when re-opening search

### Status Transitions (`t` key)
- Fetches `GET /rest/api/3/issue/{key}/transitions` asynchronously
- Shows spinner-like wait (no separate spinner, just blocked input)
- Opens `transitionPickView` modal listing available transitions
- Color-coded by status name (using same `statusColor` map)
- `enter` calls `POST /rest/api/3/issue/{key}/transitions`
- On success: shows "✓ {key} moved", refreshes issue list

### Reassign (`a` key)
- Step 1: Opens text input modal "Reassign {KEY}" — user types name or email
- `enter` fires `GET /rest/api/3/user/search?query=...` (max 10 results)
- Step 2: Shows user list "Assign {KEY} → select user" — name (28 chars) + email
- `enter` fires `PUT /rest/api/3/issue/{key}/assignee`
- On success: shows "✓ {key} reassigned", refreshes issue list

### Auto-Refresh
- Two independent timers: `myTimerCmd(d)` and `teamTimerCmd(d)` using `tea.Tick`
- Duration from `m.refreshInterval` which comes from `client.RefreshInterval()` → `cfg.RefreshInterval`
- Both timers start when the app loads
- Timer messages: `tickMsg` (my tasks) and `teamTickMsg` (team)
- Header shows: `auto-refresh {duration}` (e.g. "auto-refresh 5m")

### Config Error Screen (`ui/view.go:configErrorView`)
- Shown when `m.client == nil` (meaning `config.LoadConfig()` returned an error)
- Displays error message + `.env` template with required keys
- Rounded red border box, vertically centered
- Only `q` works — user must fix config and restart

---

## Data Types (`jira/types.go`)

```go
type Issue struct {
    Key    string
    Fields IssueFields  // Summary, Status, Priority, Assignee
}

type User struct {
    AccountID    string  // used for AssignIssue
    DisplayName  string  // shown in UI
    EmailAddress string  // shown in reassign picker
}

type Transition struct {
    ID   string  // passed to DoTransition
    Name string  // shown in UI (e.g. "In Progress", "Done")
}
```

---

## Styles & Colors (`ui/styles.go`)

```go
purple  = "#7C3AED"
white   = "#FFFFFF"
grey    = "#9CA3AF"
green   = "#10B981"
red     = "#EF4444"
blue    = "#3B82F6"
orange  = "#F97316"
yellow  = "#F59E0B"
dimGrey = "#6B7280"
teal    = "#0EA5E9"
```

**Status color map** (hardcoded — see Known Issues):
```
"To Do"              → dimGrey
"In Progress"        → blue
"In Review"          → purple
"Done"               → green
"REJECTED"           → red
"READY TO TEST"      → teal
"Tier 1 in progress" → orange
"TIER 2"             → yellow
"Need More Info"     → yellow
```
Any unknown status defaults to `grey`.

**Priority color map**:
```
"Highest" → red
"High"    → orange
"Medium"  → yellow
"Low"     → blue
"Lowest"  → dimGrey
```

---

## Message / Command Flow (`ui/messages.go`)

All async operations follow the Bubbletea pattern: return a `tea.Cmd` from `Update()`, the command runs in a goroutine and returns a `tea.Msg`, which triggers another `Update()` call.

| Msg type | Triggered by | Result |
|----------|-------------|--------|
| `issuesMsg` | `fetchMyIssues()` | Updates `m.myIssues`, sets last updated time |
| `teamIssuesMsg` | `fetchTeamIssues()` | Updates `m.teamIssues`, sets last updated time |
| `usersMsg` | `searchUsers()` | Populates user list, enters `stateReassignPick` |
| `assignedMsg` | `assignIssue()` | Shows success msg, refreshes issues |
| `transitionsMsg` | `fetchTransitions()` | Sets `m.transitions`, enters `stateTransitionPick` |
| `transitionedMsg` | `doTransition()` | Shows success msg, refreshes issues |
| `tickMsg` | `myTimerCmd(d)` | Re-fetches my issues, restarts timer |
| `teamTickMsg` | `teamTimerCmd(d)` | Re-fetches team issues, restarts timer |
| `splashDoneMsg` | `splashTimer()` | Transitions from `stateSplash` to `stateLoading` |
| `errMsg` | Any failed cmd | Sets `m.err`, enters `stateError` |

---

## Known Issues (Things to Fix)

### 1. Hardcoded status colors are org-specific
**File**: `ui/styles.go` lines 61–71

The `statusColor` map contains statuses specific to Orbit Digital (`"READY TO TEST"`, `"Tier 1 in progress"`, `"TIER 2"`, `"Need More Info"`, `"REJECTED"`). Any other Jira org using different status names will get grey for everything.

**Fix**: For unknown statuses, cycle through a palette instead of returning grey. Example approach:
```go
func getStatusColor(name string) lipgloss.Color {
    if c, ok := statusColor[name]; ok {
        return c
    }
    // hash name to pick a deterministic color from memberColors
    h := 0
    for _, b := range name { h += int(b) }
    return memberColors[h % len(memberColors)]
}
```

### 2. No quick way to clear active filters
**Current**: Must open filter modal → select "All" → confirm.
**Fix**: Add `c` key on My Tasks and Team tabs to instantly clear all active filters (status filter, name filter, search query) for that tab. One keystroke.

### 3. `SearchUsers` requires Browse Users permission
`GET /rest/api/3/user/search` may be restricted on some Jira Cloud instances by admin policy. Not a code bug but worth noting in README so users understand why reassign search might return no results.

### 4. Two loading patterns are inconsistent
Initial load uses `stateLoading` (full screen spinner). Tab switches use `teamLoading: bool` flag inside `stateMain`. Both work but the code is asymmetric. Low priority — only matters if refactoring.

### 5. `last updated` timestamp not displayed in header
`m.myLastUpdated` and `m.teamLastUpdated` are tracked but the header only shows `auto-refresh {interval}`. Showing the actual last-fetched time would help users know if data is stale.

---

## Recommended Additions (Prioritized)

These were reviewed and agreed as valuable for a future version. Implement in this order:

### High Priority

**1. Dynamic status colors** (fixes Known Issue #1)
- Makes app truly org-agnostic
- Small change in `ui/styles.go` or `ui/helpers.go`
- Replace the grey fallback in `statusColor` lookup with a deterministic color from `memberColors`

**2. `c` key to clear filters** (fixes Known Issue #2)
- Add to `handleMyTasksKey` and `handleTeamKey` in `ui/keys.go`
- Resets `myStatusFilter`, `mySearchQuery` (or `teamNameFilter`, `teamStatusFilter`, `teamSearchQuery`)
- Add `c` to the help bar in `ui/view.go:helpBar()`

### Medium Priority

**3. `g` / `G` to jump to top / bottom of list**
- Vim-standard. Add to `handleMyTasksKey` and `handleTeamKey`
- `g` → set cursor to 0, `G` → set cursor to `len(filtered)-1`

**4. `JIRA_PROJECT` env var to pre-filter issues**
- Useful for teams with many projects — filter to only show e.g. `DX-*` tickets
- In `config/config.go`: add `ProjectKey string` field
- In `jira/api.go`: if `cfg.ProjectKey != ""`, prepend `project = "DX" AND ` to JQL

**5. Show issue URL in status bar when cursor moves**
- Show the browse URL at the bottom so user can preview before pressing enter
- Low implementation cost, useful for power users

### Low Priority

**6. `JIRA_PROJECT` multi-value support**
- Allow `JIRA_PROJECT=DX,DMP` to filter multiple projects

**7. Display `last updated` time in header**
- Already tracked in `m.myLastUpdated` / `m.teamLastUpdated`
- Just add to the header string in `ui/view.go:mainView()`

---

## How to Build & Run

```bash
# build (output goes to build/ folder)
go build -o build/jiramage .

# run
./build/jiramage
```

No flags. All config via `.env`. Run from the project root directory so `.env` is found.

---

## What Was Built in This Session (History)

Tracking major decisions so you understand why things are the way they are:

1. **Started as single-file** → refactored into `config/`, `jira/`, `ui/` packages with many split files
2. **Pagination was broken** — originally used `startAt` (GET v2 API pattern). Fixed to use `nextPageToken` with `POST /rest/api/3/search/jql` (v3 API)
3. **App was named "MP-JIRA"** → renamed to "Jiramage" everywhere. The go module is still `jira-dashboard` (internal only, doesn't affect UX)
4. **Logo uses character map approach** — `buildLogo()` in `splash.go` defines each letter individually as `[6]string` of Unicode box chars, assembled row-by-row. DO NOT go back to hardcoded rows — they're error-prone
5. **`shortName()`** uses `findLast()` byte scanner instead of `strings.LastIndex` (which takes a string arg, not a byte)
6. **`cfg` fields are private** — exposed via `BaseURL()` and `RefreshInterval()` methods on `JiraClient` to avoid leaking config internals into the ui package
7. **Self user included in team** — `GetTeamIssues()` prepends `c.cfg.Email` to `TeamEmails` so the logged-in user appears in Team tab and Dashboard too
8. **Error model** — `NewErrorModel(err)` creates a Model with `client == nil` and `state = stateError`. `Init()` guards `if m.client == nil { return nil }` to avoid nil panics. `configErrorView()` checks `m.client == nil` to distinguish config errors from runtime errors
9. **Go version updated** — `go.mod` bumped from `go 1.22` to `go 1.26` to match local toolchain (Go 1.26.2 darwin/arm64)
10. **Build output moved** — binary now outputs to `build/jiramage` instead of project root. Build command: `go build -o build/jiramage .`

---

## File Quick Reference

| If you need to... | Look in... |
|---|---|
| Change env var names or defaults | `config/config.go` |
| Change JQL queries or API calls | `jira/api.go` |
| Add a new Jira API endpoint | `jira/api.go` + `jira/types.go` + `jira/client.go` |
| Add a new keyboard shortcut | `ui/keys.go` + `ui/view.go` (helpBar) |
| Add a new view state | `ui/model.go` (iota) + `ui/update.go` + `ui/view.go` |
| Add a new modal | `ui/modals.go` + `ui/view.go` (modal switch) |
| Change colors or styles | `ui/styles.go` |
| Change the splash screen | `ui/splash.go` |
| Change the header or tabs | `ui/view.go:mainView()` |
| Add a new tea.Msg or tea.Cmd | `ui/messages.go` |
| Change filter logic | `ui/helpers.go` |
| Change how pages look | `ui/page_mytasks.go`, `ui/page_team.go`, `ui/page_dashboard.go` |
