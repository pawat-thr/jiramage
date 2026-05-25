package jira

type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary  string   `json:"summary"`
	Status   Status   `json:"status"`
	Priority Priority `json:"priority"`
	Assignee *User    `json:"assignee"`
}

type Status struct {
	Name           string         `json:"name"`
	StatusCategory StatusCategory `json:"statusCategory"`
}

type StatusCategory struct {
	Key string `json:"key"` // "new" | "indeterminate" | "done"
}

type Priority struct {
	Name string `json:"name"`
}

type User struct {
	AccountID    string `json:"accountId"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

type searchResult struct {
	Issues        []Issue `json:"issues"`
	NextPageToken string  `json:"nextPageToken"`
}

type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type transitionsResult struct {
	Transitions []Transition `json:"transitions"`
}

func (i Issue) BrowseURL(baseURL string) string {
	return baseURL + "/browse/" + i.Key
}
