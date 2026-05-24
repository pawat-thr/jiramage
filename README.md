# jiramage

A terminal dashboard for Jira Cloud — browse your tasks, monitor your team's workload, and reassign issues without leaving the terminal.

Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) · Go · v0.1.0 · by MpLab

---

## Features

| Page | Key | What it does |
|------|-----|--------------|
| **1 · My Tasks** | `1` | Issues assigned to you, filterable by status |
| **2 · Team** | `2` | Issues for all team members, filterable by name & status |
| **3 · Dashboard** | `3` | Performance meter — task count bar per member |

- Auto-refreshes every **5 minutes**
- Paginated fetch — retrieves **all** issues, not just the first 100
- Press **Enter** on any issue to open it in the browser
- Press **a** to reassign an issue to another user
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
```

> Generate an API token at **id.atlassian.com → Security → API tokens**

**4. Build & run**

```bash
go build -o build/jiramage .
./build/jiramage
```

---

## Keyboard Shortcuts

### Global

| Key | Action |
|-----|--------|
| `1` | My Tasks page |
| `2` | Team page |
| `3` | Dashboard page |
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `r` | Refresh current page |
| `q` | Quit |

### My Tasks (page 1)

| Key | Action |
|-----|--------|
| `Enter` | Open issue in browser |
| `f` | Filter by status |
| `a` | Reassign issue |

### Team (page 2)

| Key | Action |
|-----|--------|
| `Enter` | Open issue in browser |
| `f` | Filter by member name |
| `s` | Filter by status |
| `a` | Reassign issue |

### Dashboard (page 3)

| Key | Action |
|-----|--------|
| `Enter` | Jump to Team page filtered to selected member |
| `s` | Filter bars by status |

---

## Project Structure

```
.
├── main.go               — entry point
├── .env                  — your credentials (never commit this)
├── .env.example          — template
├── config/
│   └── config.go         — loads .env and Config struct
├── jira/
│   ├── types.go          — Issue, User, Status, Priority types
│   ├── client.go         — HTTP client (get / post / put)
│   └── api.go            — GetMyIssues, GetTeamIssues, SearchUsers, AssignIssue
└── ui/
    ├── model.go           — Model struct, NewModel, Init
    ├── messages.go        — Bubbletea message types & commands
    ├── styles.go          — Colors and lipgloss styles
    ├── update.go          — Update() and key handler dispatch
    ├── keys.go            — Per-page key handlers
    ├── view.go            — View(), mainView(), renderIssueRow()
    ├── splash.go          — Splash screen & logo
    ├── page_mytasks.go    — My Tasks page view
    ├── page_team.go       — Team page view
    ├── page_dashboard.go  — Dashboard page view
    ├── modals.go          — Filter & reassign modals
    └── helpers.go         — Scroll, filter, shortName, openURL utilities
```

---

## Switching to a Different Jira Org

Only `.env` needs to change:

```env
JIRA_URL=https://newcompany.atlassian.net
JIRA_EMAIL=you@newcompany.com
JIRA_TOKEN=new_token
TEAM_EMAILS=person1@newcompany.com,person2@newcompany.com
```

All API calls use standard Jira Cloud REST API v3 — no code changes needed.

---

## License

MIT
