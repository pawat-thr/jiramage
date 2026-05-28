# Jiramage — Claude Context File

This file gives a new Claude instance full context about this project: what it is, every decision made, the current state of the code, known issues, and what to work on next. Read this before touching anything.

---

## What This Project Is

**Jiramage** (from "Jira Management") is a terminal UI (TUI) dashboard for Jira Cloud. It lets you browse your tasks, monitor your team's workload, reassign issues, transition issue statuses, and track how many cards each member fixed in a sprint — all without leaving the terminal.

- **Language**: Go 1.26
- **TUI Framework**: Bubbletea v1.2.4 (Elm architecture: Model / Update / View)
- **Styling**: charmbracelet/lipgloss v1.0.0
- **Widgets**: charmbracelet/bubbles v0.20.0 (spinner, textinput)
- **Binary name**: `jiramage`
- **Module name** (go.mod): `jira-dashboard` (legacy name, do NOT rename — it would break all imports)
- **Go toolchain**: 1.26.2 (darwin/arm64) — `go.mod` declares `go 1.26`
- **Version**: v0.1.3 (hardcoded in `ui/model.go` as `appVersion`)
- **Credit**: "by MpLab" (hardcoded as `appCredit`)

---

## Owner / User Context

- Developer: Pawat T. (`pawat.t@orbitdigital.co.th`)
- Company: Orbit Digital
- Jira instance: `https://orbitdigital.atlassian.net`
- Team emails in `.env`: tanawat.k, peradone.c, taninchot.p — all at orbitdigital.co.th
- The app is intended to be **org-agnostic** — just changing `.env` should make it work for any Jira Cloud org

---

## Project Structure

```
jiramage/                   ← project root (binary built here)
├── main.go                 ← entry point
├── go.mod                  ← module: jira-dashboard, go 1.26
├── go.sum
├── jiramage               ← compiled binary (run with ./jiramage)
├── .env                    ← real credentials, never commit
├── .env.example            ← template for new users
├── README.md
├── .claude/docs/claude-context.md  ← this file
├── config/
│   └── config.go           ← Config struct + LoadConfig()
├── jira/
│   ├── types.go            ← Issue, User, Transition, StatusCategory types
│   ├── client.go           ← JiraClient, HTTP helpers (get/post/put), all accessor methods
│   └── api.go              ← GetMyIssues, GetTeamIssues, GetFixedIssues, GetAllFixedIssues,
│                              SearchUsers, AssignIssue, GetTransitions, DoTransition,
│                              projectFilter(), labelFilter()
└── ui/
    ├── model.go            ← Model struct, NewModel, NewErrorModel, Init, viewState/tabView iota
    ├── messages.go         ← all tea.Msg types + tea.Cmd factories
    ├── update.go           ← Update() — the main state machine handler
    ├── view.go             ← View() dispatcher + mainView + helpBar + renderIssueRow
    │                          + visualTruncate() + visualPad()
    ├── splash.go           ← buildLogo() + splashView()
    ├── styles.go           ← all lipgloss colors and styles
    ├── helpers.go          ← filter helpers, formatDuration, scrollWindow, openURL,
    │                          emailUsername(), listHeight(), etc.
    ├── keys.go             ← handleMyTasksKey, handleTeamKey, handleDashboardKey, handleFixedKey
    ├── modals.go           ← statusFilterModal, nameFilterModal, reassignSearchView,
    │                          reassignPickView, transitionPickView
    ├── page_mytasks.go     ← myTasksView()
    ├── page_team.go        ← teamView()
    ├── page_dashboard.go   ← dashboardView()
    └── page_fixed.go       ← fixedView(), fixedSummaryView(), fixedDrillView()
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
JIRA_PROJECT=DMP,APP,DX
JIRA_LABELS=R6.1-SIT
JIRA_FIXED_STATUS=SIT DEPLOYED
JIRA_TEAM_FROM=2024-05-01
```

| Variable | Required | Description |
|---|---|---|
| `JIRA_URL` | ✓ | Jira Cloud base URL |
| `JIRA_EMAIL` | ✓ | Your Atlassian account email |
| `JIRA_TOKEN` | ✓ | Atlassian API token |
| `TEAM_EMAILS` | optional | Comma-separated — used for Team, Dashboard, and Fixed tabs |
| `REFRESH_INTERVAL` | optional | Go duration string (e.g. `5m`, `30s`, `1h`), default `5m` |
| `JIRA_PROJECT` | optional | Comma-separated project keys. Normalized to uppercase. Filters ALL fetched issues via JQL. Header shows `[DMP,APP,DX]` in teal when active. |
| `JIRA_LABELS` | optional | Comma-separated Jira label keys (e.g. `R6.1-SIT`). Used by Tab 4 Fixed to filter which cards count as fixed. If not set, Tab 4 shows a setup message. |
| `JIRA_FIXED_STATUS` | optional | Status name that means "fixed/deployed" (default: `SIT DEPLOYED`). Used in Tab 4 JQL. |
| `JIRA_TEAM_FROM` | optional | Date cutoff for Team and Dashboard tabs (default: `2024-05-01`). Adds `AND created >= "date"` to team JQL to exclude old cards. |

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
- My Issues: `[project in (...) AND] assignee = currentUser() ORDER BY key DESC`
- Team Issues: `[project in (...) AND] assignee in ("email1","email2",...) AND created >= "2024-05-01" ORDER BY key DESC`
- Fixed Issues (per member): `status changed to "SIT DEPLOYED" by "email" AND labels in ("R6.1-SIT") ORDER BY key DESC`
- `projectFilter()` in `jira/api.go` builds the `project in (...)` prefix; empty string when no projects set
- `labelFilter()` in `jira/api.go` builds the `labels in (...)` clause; empty string when no labels set
- `JIRA_TEAM_FROM` adds `AND created >= "date"` to team JQL only (not My Issues)

**Fields fetched**: `summary`, `status`, `priority`, `assignee`
- `status` always includes the nested `statusCategory` object (Jira returns it automatically)

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

Tabs (`tabView` iota): `tabMyTasks`, `tabTeam`, `tabDashboard`, `tabFixed`

---

## All Keyboard Shortcuts

### Global
| Key | Action |
|-----|--------|
| `q` / `ctrl+c` | Quit |
| `1` | Switch to My Tasks tab |
| `2` | Switch to Team tab |
| `3` | Switch to Dashboard tab |
| `4` | Switch to Fixed tab |
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `esc` | Cancel / close modal / clear search |

### My Tasks tab
| Key | Action |
|-----|--------|
| `enter` | Open issue in browser |
| `h` | Toggle active-only mode |
| `f` | Open status filter modal |
| `/` | Open search input |
| `t` | Fetch transitions for selected issue → open transition picker |
| `a` | Open reassign search for selected issue |
| `r` | Refresh my issues |

### Team tab
| Key | Action |
|-----|--------|
| `enter` | Open issue in browser |
| `h` | Toggle active-only mode |
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
| `h` | Toggle active-only mode |
| `s` | Open status filter modal |
| `r` | Refresh team issues |

### Fixed tab (summary level)
| Key | Action |
|-----|--------|
| `enter` | Drill into selected member's fixed cards |
| `r` | Refresh fixed issues |

### Fixed tab (drill-down level)
| Key | Action |
|-----|--------|
| `enter` | Open issue in browser |
| `esc` | Back to summary |
| `r` | Refresh (returns to summary) |

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
- Each letter is a 6-row `[6]string` defined in `buildLogo()` char map
- Below logo: subtitle, version (`v0.1.3`), credit (`by MpLab`), hint text
- Vertically centered based on terminal height

### My Tasks (`ui/page_mytasks.go`)
- Columns: KEY | SUMMARY | STATUS | PRIORITY — SUMMARY column is dynamic width based on terminal
- Color-coded priority: Highest=red, High=orange, Medium=yellow, Low=blue, Lowest=dimGrey
- Color-coded status (see styles section)
- Scrollable list with scroll indicator; `─` separator line under column header
- Active filter badges shown above list (teal `● active only` badge + green filter badges)
- Search input shown inline when `stateMySearch` is active
- Filtering: AND logic — hideCompleted AND status filter AND search query

### Team (`ui/page_team.go`)
- Columns: KEY | SUMMARY | STATUS | ASSIGNEE — SUMMARY column is dynamic width
- Assignee displayed as short name from parentheses: `"John Smith (johnsmith)"` → `"johnsmith"`
  - Implemented in `helpers.go:shortName()` using `findLast()` byte scanner
- Only fetches issues `created >= JIRA_TEAM_FROM` (default `2024-05-01`) to exclude stale old cards
- Loading spinner on first load (lazy — only fetches when tab first visited)
- Dual filter badges: name badge + status badge, plus teal `● active only` badge
- Filtering: AND logic — hideCompleted AND name filter AND status filter AND search query

### Dashboard (`ui/page_dashboard.go`)
- Shows one row per team member + "All" aggregate row
- Uses the same `GetTeamIssues()` data as Tab 2 (includes `JIRA_TEAM_FROM` date filter)
- Horizontal bar chart: bar width proportional to count, max 52 chars, scales to terminal width
- Member colors cycle through: purple → teal → green → orange → red → yellow (`memberColors`)
- `enter` on a member row switches to Team tab and sets name filter to that member
- `s` opens status filter; `h` toggles active-only mode

### Fixed (`ui/page_fixed.go`) — added v0.1.3
- **Purpose**: Tracks how many cards each team member fixed (transitioned to `JIRA_FIXED_STATUS`), even after those cards were reassigned to QA
- **Requires**: `JIRA_LABELS` set in `.env`. If not set, shows a setup message
- **JQL used**: `status changed to "SIT DEPLOYED" by "email" AND labels in ("R6.1-SIT")` — one query per team member (including logged-in user)
- **Two levels**:
  - **Summary**: bar chart showing each member's fix count (like Dashboard). Press `enter` on a member to drill down
  - **Drill-down**: list of that member's fixed cards — KEY, SUMMARY, current STATUS, current ASSIGNEE (usually QA after reassignment)
- Lazy-loaded on first visit; auto-refreshes on same interval as other tabs
- `fixedTimerCmd` starts when first `fixedIssuesMsg` is received
- Member order: `cfg.Email` (logged-in user) first, then `cfg.TeamEmails` — consistent with `client.AllMembers()`
- Display name in summary uses `emailUsername(email)` helper (e.g. `pawat.t` from `pawat.t@orbitdigital.co.th`)

### Active-Only Mode (`h` key) — added v0.1.1
- **Default ON** — hides issues in Jira's "done" status category
- Uses `issue.Fields.Status.StatusCategory.Key == "done"` — org-agnostic
- Applies to My Tasks, Team, and Dashboard

### Project Filter (`JIRA_PROJECT`) — added v0.1.1
- Filters ALL fetched issues at the JQL level via `projectFilter()`
- Header shows `[APP,DX]` in teal when active

### Search (`/` key)
- Inline text input, case-insensitive substring match on `key + " " + summary`
- Real-time filtering; `esc` clears and exits

### Status Transitions (`t` key)
- Async fetch of transitions, modal picker, applies via POST
- On success: shows "✓ {key} moved", refreshes issue list

### Reassign (`a` key)
- Two-step: search users → pick user → PUT assignee
- On success: shows "✓ {key} reassigned", refreshes issue list

### Auto-Refresh
- Three independent timers: `myTimerCmd`, `teamTimerCmd`, `fixedTimerCmd`
- All use `m.refreshInterval` from `cfg.RefreshInterval`
- Fixed timer only starts after Tab 4 is first visited and data loads

### Help Bar (`ui/view.go:helpBar`)
- **Two rows** to avoid overflow on narrow terminals:
  - Row 1: global navigation (`1`–`4`, `↑↓`, `r`, `q`) — always visible, never wraps
  - Row 2: tab-specific shortcuts
- Both rows have consistent `  ` (2-space) left margin

### Column Layout
- SUMMARY column (`col2`) is **dynamic**: `col2 = m.width - col1 - col3 - col4 - 6`
  - Minimum 20 chars; expands to fill available terminal width
  - Applied in Tab 1, Tab 2, and Tab 4 drill-down
- `─────` separator line after column header row in all list views
- Status names truncated to `col3` width to prevent column overflow

### Thai-Safe Summary Truncation (`ui/view.go:visualTruncate`, `visualPad`)
- Thai text contains zero-width combining characters (vowels/tone marks above consonants)
- These count as 1 rune but 0 terminal columns, so rune-based `%-*s` padding misaligns columns
- Fix: `visualTruncate(s, maxW)` and `visualPad(s, w)` use `lipgloss.Width()` (which handles zero-width chars) for all width calculations
- `renderIssueRow` in `view.go` uses these for the SUMMARY column

---

## Data Types (`jira/types.go`)

```go
type Issue struct {
    Key    string
    Fields IssueFields  // Summary, Status, Priority, Assignee
}

type IssueFields struct {
    Summary  string
    Status   Status
    Priority Priority
    Assignee *User
}

type Status struct {
    Name           string
    StatusCategory StatusCategory
}

type StatusCategory struct {
    Key string  // "new" | "indeterminate" | "done"
}

type User struct {
    AccountID    string
    DisplayName  string
    EmailAddress string
}

type Transition struct {
    ID   string
    Name string
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

**Named styles**:
- `filterBadgeStyle` — green background, status/name/search filter badges
- `modeBadgeStyle` — teal background, `● active only` mode indicator and label badges in Tab 4
- `selStyle` — purple background, selected row highlight
- `keyStyle` — green bold, keyboard shortcut labels in helpBar
- `headerStyle` — grey foreground, underline — column header row
- `modalBorder` — purple rounded border, used by all modals

**Status color map**:
```
"To Do"              → dimGrey
"In Progress"        → blue
"In Review"          → purple
"Done"               → green
"REJECTED"           → red
"READY TO TEST"      → teal
"SIT DEPLOYED"       → orange
"SIT TESTING"        → teal
"Tier 1 in progress" → orange
"TIER 2"             → yellow
"Need More Info"     → yellow
```
Any unknown status defaults to `grey`.

**Priority color map**:
```
"Highest" → red    "High" → orange    "Medium" → yellow
"Low"     → blue   "Lowest" → dimGrey
```

---

## Message / Command Flow (`ui/messages.go`)

| Msg type | Triggered by | Result |
|----------|-------------|--------|
| `issuesMsg` | `fetchMyIssues()` | Updates `m.myIssues`, sets last updated time |
| `teamIssuesMsg` | `fetchTeamIssues()` | Updates `m.teamIssues`, sets last updated time |
| `fixedIssuesMsg` | `fetchAllFixed()` | Updates `m.fixedByMember`, sets `fixedLastUpdated` |
| `usersMsg` | `searchUsers()` | Populates user list, enters `stateReassignPick` |
| `assignedMsg` | `assignIssue()` | Shows success msg, refreshes issues |
| `transitionsMsg` | `fetchTransitions()` | Sets `m.transitions`, enters `stateTransitionPick` |
| `transitionedMsg` | `doTransition()` | Shows success msg, refreshes issues |
| `tickMsg` | `myTimerCmd(d)` | Re-fetches my issues, restarts timer |
| `teamTickMsg` | `teamTimerCmd(d)` | Re-fetches team issues, restarts timer |
| `fixedTickMsg` | `fixedTimerCmd(d)` | Re-fetches all fixed issues, restarts timer |
| `splashDoneMsg` | `splashTimer()` | Transitions from `stateSplash` to `stateLoading` |
| `errMsg` | Any failed cmd | Sets status message or `stateError` |

---

## Client Methods (`jira/client.go`)

| Method | Returns | Purpose |
|--------|---------|---------|
| `BaseURL()` | string | Jira instance URL for building browse links |
| `RefreshInterval()` | time.Duration | From cfg — used by all timers |
| `ProjectLabel()` | string | e.g. `"APP,DX"` for header display |
| `LabelLabel()` | string | e.g. `"R6.1-SIT"` for Tab 4 display |
| `HasLabels()` | bool | Whether `JIRA_LABELS` is configured |
| `FixedStatusName()` | string | e.g. `"SIT DEPLOYED"` for Tab 4 display |
| `TeamFromDate()` | string | e.g. `"2024-05-01"` |
| `AllMembers()` | []string | `[cfg.Email] + cfg.TeamEmails` — ordered member list for Tab 4 |

---

## Known Issues (Things to Fix)

### 1. Hardcoded status colors are org-specific
The `statusColor` map in `ui/styles.go` contains Orbit Digital statuses. Other orgs get grey for unknown status names. Fix: hash unknown names to `memberColors` palette using `StatusCategory` as fallback.

### 2. No quick way to clear active filters
Must open filter modal → "All" → confirm. Fix: add `c` key in `handleMyTasksKey` and `handleTeamKey` to clear all filters for that tab in one keystroke.

### 3. `SearchUsers` requires Browse Users permission
`GET /rest/api/3/user/search` may be restricted by Jira admin. Not a code bug — worth noting in README.

### 4. Two loading patterns are inconsistent
Initial load uses `stateLoading` (full screen). Tab 2/3/4 switches use `teamLoading`/`fixedLoading` bool inside `stateMain`. Asymmetric but functional. Low priority.

---

## Recommended Additions (Prioritized)

### High Priority

**1. Dynamic status colors** (fixes Known Issue #1)
- Replace grey fallback in `renderIssueRow` with hash-to-`memberColors`
- `StatusCategory` already on every issue

**2. `c` key to clear filters** (fixes Known Issue #2)
- Add to `handleMyTasksKey` and `handleTeamKey` in `ui/keys.go`
- Add `c` entry to Row 2 of `helpBar()`

### Medium Priority

**3. `g` / `G` to jump to top / bottom**
- `g` → cursor 0, `G` → cursor `len(filtered)-1`
- Add to `handleMyTasksKey`, `handleTeamKey`, `handleFixedKey`

**4. Tab 4 date range filter**
- Currently shows all time for the configured labels
- Could add `JIRA_FIXED_FROM` similar to `JIRA_TEAM_FROM`

---

## How to Build & Run

```bash
# build
go build -o jiramage .

# run
./jiramage
```

No flags. All config via `.env`. Run from project root so `.env` is found.

---

## What Was Built (History)

1. **Started as single-file** → refactored into `config/`, `jira/`, `ui/` packages
2. **Pagination was broken** — fixed from `startAt` (v2) to `nextPageToken` with `POST /rest/api/3/search/jql` (v3)
3. **App was named "MP-JIRA"** → renamed to "Jiramage". Module is still `jira-dashboard` (internal only)
4. **Logo uses character map** — `buildLogo()` defines each letter as `[6]string`, assembled row-by-row. Do NOT go back to hardcoded rows
5. **`shortName()`** uses `findLast()` byte scanner instead of `strings.LastIndex` (wrong signature)
6. **`cfg` fields are private** — exposed via accessor methods on `JiraClient` to avoid leaking config into ui
7. **Self user included in team** — `GetTeamIssues()` prepends `c.cfg.Email` to `TeamEmails`; `AllMembers()` does the same for Tab 4
8. **Error model** — `NewErrorModel(err)` sets `client == nil`; `Init()` guards against nil; `configErrorView()` checks `client == nil` to distinguish config vs runtime errors
9. **Active-only mode (v0.1.1)** — `hideCompleted bool` on Model, default `true`. Uses `StatusCategory.Key == "done"` — org-agnostic. Two badge styles: `filterBadgeStyle` (green) vs `modeBadgeStyle` (teal) — visually distinct intentionally
10. **Project filter (v0.1.1)** — `JIRA_PROJECT` parsed into `cfg.ProjectKeys`. JQL prefix via `projectFilter()`. Applied to My Issues and Team Issues. `ProjectLabel()` exposes it for the header
11. **Fixed tab (v0.1.3)** — Tab 4 uses `status changed to "X" by "email" AND labels in (...)` JQL per team member. Captures cards the developer fixed even after reassigning to QA. `GetAllFixedIssues()` runs one query per member sequentially. Summary bar chart + drill-down card list. Three new env vars: `JIRA_LABELS`, `JIRA_FIXED_STATUS`, `JIRA_TEAM_FROM`
12. **Team date filter (v0.1.3)** — `JIRA_TEAM_FROM` adds `AND created >= "date"` to `GetTeamIssues()` JQL only. Prevents stale old cards cluttering Tab 2 and Tab 3
13. **Dynamic column widths (v0.1.3)** — SUMMARY column (`col2`) is `m.width - fixed_cols - 6`, minimum 20. Fills available terminal width. Applied in Tab 1, 2, and Tab 4 drill-down
14. **Two-row help bar (v0.1.3)** — Row 1: global nav always visible. Row 2: tab-specific shortcuts. Prevents overflow on narrow terminals. `listHeight` fixed increased from 7 to 8
15. **Thai-safe truncation (v0.1.3)** — `visualTruncate()` and `visualPad()` use `lipgloss.Width()` instead of rune count. Fixes column misalignment when Thai combining marks (zero-width chars) appear in summaries. All SUMMARY column rendering uses these functions

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
| Change column widths or row rendering | `ui/view.go:renderIssueRow`, `visualTruncate`, `visualPad` |
| Change how list pages look | `ui/page_mytasks.go`, `ui/page_team.go`, `ui/page_dashboard.go`, `ui/page_fixed.go` |
| Change Tab 4 Fixed logic | `ui/page_fixed.go`, `ui/keys.go:handleFixedKey`, `jira/api.go:GetAllFixedIssues` |
