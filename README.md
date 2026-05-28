# jiramage

A terminal dashboard for Jira Cloud — browse your tasks, monitor your team's workload, track sprint fixes, and reassign issues without leaving the terminal.

Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) · Go · v0.1.3 · by MpLab

---

## Features

| Tab | Key | What it does |
|-----|-----|--------------|
| **1 · My Tasks** | `1` | Issues assigned to you, filterable by status |
| **2 · Team** | `2` | Team issues filterable by member & status (scoped by date) |
| **3 · Dashboard** | `3` | Task count bar chart per team member |
| **4 · Fixed** | `4` | Cards each member fixed this sprint (by label) |

- Auto-refreshes on configurable interval (default **5 minutes**)
- Paginated fetch — retrieves **all** matching issues, not just first 100
- Press **Enter** on any issue to open it in the browser
- Press **h** to toggle active-only mode (hides "done" category issues)
- Press **t** to transition an issue status directly from the terminal
- Press **a** to reassign an issue to another user
- Columns adapt to terminal width — summaries use all available space
- Thai and multi-byte text handled correctly (zero-width combining chars)
- Works with **any Jira Cloud org** — just change `.env`

---

## Requirements

- Go 1.26+
- A Jira Cloud account with an API token

---

## Setup

**1. Clone the repo**

```bash
git clone <repo-url>
cd jiramage
```

**2. Install dependencies**

```bash
go mod download
```

**3. Create your `.env` file**

```bash
cp .env.example .env
```

Edit `.env` with your credentials:

```env
JIRA_URL=https://yourcompany.atlassian.net
JIRA_EMAIL=you@company.com
JIRA_TOKEN=your_api_token_here
TEAM_EMAILS=teammate1@company.com,teammate2@company.com
REFRESH_INTERVAL=5m
JIRA_PROJECT=DMP,APP,DX
JIRA_LABELS=R6.1-SIT
JIRA_FIXED_STATUS=SIT DEPLOYED
JIRA_TEAM_FROM=2024-05-01
```

> Generate an API token at **id.atlassian.com → Security → API tokens**

**4. Build & run**

```bash
go build -o jiramage .
./jiramage
```

---

## Configuration

| Variable | Required | Default | Description |
|---|---|---|---|
| `JIRA_URL` | ✓ | — | Jira Cloud base URL |
| `JIRA_EMAIL` | ✓ | — | Your Atlassian account email |
| `JIRA_TOKEN` | ✓ | — | Atlassian API token |
| `TEAM_EMAILS` | — | — | Comma-separated team emails for Tab 2, 3, 4 |
| `REFRESH_INTERVAL` | — | `5m` | Auto-refresh interval (e.g. `30s`, `5m`, `1h`) |
| `JIRA_PROJECT` | — | — | Comma-separated project keys to filter all tabs (e.g. `APP,DX`) |
| `JIRA_LABELS` | — | — | Comma-separated labels for Tab 4 Fixed (e.g. `R6.1-SIT`). Tab 4 is disabled if not set. |
| `JIRA_FIXED_STATUS` | — | `SIT DEPLOYED` | Status name that means "fixed" — used in Tab 4 JQL |
| `JIRA_TEAM_FROM` | — | `2024-05-01` | Excludes issues created before this date from Tab 2 & 3 |

---

## Keyboard Shortcuts

### Global

| Key | Action |
|-----|--------|
| `1` | My Tasks tab |
| `2` | Team tab |
| `3` | Dashboard tab |
| `4` | Fixed tab |
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `h` | Toggle active-only mode (hide/show completed issues) |
| `r` | Refresh current tab |
| `q` / `Ctrl+C` | Quit |

### My Tasks (Tab 1)

| Key | Action |
|-----|--------|
| `Enter` | Open issue in browser |
| `f` | Filter by status |
| `/` | Search by title or key |
| `t` | Transition issue status |
| `a` | Reassign issue |

### Team (Tab 2)

| Key | Action |
|-----|--------|
| `Enter` | Open issue in browser |
| `f` | Filter by member name |
| `s` | Filter by status |
| `/` | Search by title or key |
| `t` | Transition issue status |
| `a` | Reassign issue |

### Dashboard (Tab 3)

| Key | Action |
|-----|--------|
| `Enter` | Jump to Team tab filtered to selected member |
| `s` | Filter bars by status |

### Fixed (Tab 4)

| Key | Action |
|-----|--------|
| `Enter` | Drill into member's fixed cards |
| `Esc` | Back to summary |
| `Enter` (on card) | Open issue in browser |

---

## Tab 4 — Fixed

Tab 4 answers the question: **"How many cards did each person fix this sprint?"**

The problem it solves: when a developer fixes a defect and redeploys, they typically reassign the card to QA. After reassignment the card no longer appears under the developer's name, so the Dashboard undercounts their work.

Tab 4 uses Jira's historical JQL operator to query who *transitioned* the card to the fixed status, regardless of current assignee:

```
status changed to "SIT DEPLOYED" by "dev@company.com"
AND labels in ("R6.1-SIT")
```

**Summary view** — bar chart of fix counts per member:
```
  Fixed by Team  [R6.1-SIT]  →  SIT DEPLOYED

  ──────────────────────────────────────────
  ▶ All              [████████████████]  20

    pawat.t          [██████░░░░░░░░░░]   6   30%
    tanawat.k        [████░░░░░░░░░░░░]   4   20%
    ...
  ──────────────────────────────────────────
  total 20 cards fixed  |  enter → see cards
```

Press `Enter` on any member to see their individual fixed cards (KEY, SUMMARY, current STATUS, current ASSIGNEE).

To change the sprint label, update `JIRA_LABELS=R6.2-SIT` in `.env`.

> **Note:** Tab 4 requires `JIRA_LABELS` to be set. Without it, the tab shows a setup message.

---

## Reassign

The reassign flow (`a` key) requires **Browse Users** permission on your Jira Cloud instance. If the user search returns no results, ask your Jira admin to check the "Anyone can search for user profiles" setting under **Settings → System → User management**.

---

## Project Structure

```
.
├── main.go               — entry point
├── .env                  — your credentials (never commit this)
├── .env.example          — template
├── config/
│   └── config.go         — loads .env, Config struct
├── jira/
│   ├── types.go          — Issue, User, Status, Priority, Transition types
│   ├── client.go         — HTTP client + accessor methods
│   └── api.go            — GetMyIssues, GetTeamIssues, GetAllFixedIssues,
│                           SearchUsers, AssignIssue, GetTransitions, DoTransition
└── ui/
    ├── model.go           — Model struct, viewState/tabView iota, NewModel, Init
    ├── messages.go        — Bubbletea message types & commands
    ├── styles.go          — Colors and lipgloss styles
    ├── update.go          — Update() and key handler dispatch
    ├── keys.go            — Per-tab key handlers
    ├── view.go            — View(), mainView(), helpBar(), renderIssueRow()
    ├── splash.go          — Splash screen & logo
    ├── helpers.go         — Filter helpers, scrollWindow, visualTruncate, emailUsername, etc.
    ├── page_mytasks.go    — My Tasks tab view
    ├── page_team.go       — Team tab view
    ├── page_dashboard.go  — Dashboard tab view
    ├── page_fixed.go      — Fixed tab view (summary + drill-down)
    └── modals.go          — Filter, transition & reassign modals
```

---

## Switching to a Different Jira Org

Only `.env` needs to change:

```env
JIRA_URL=https://newcompany.atlassian.net
JIRA_EMAIL=you@newcompany.com
JIRA_TOKEN=new_token
TEAM_EMAILS=person1@newcompany.com,person2@newcompany.com
JIRA_FIXED_STATUS=Done
JIRA_LABELS=sprint-42
```

All API calls use standard Jira Cloud REST API v3 — no code changes needed.

---

## License

MIT
