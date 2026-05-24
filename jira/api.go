package jira

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func (c *JiraClient) searchAll(jql string, fields []string) ([]Issue, error) {
	const pageSize = 100
	var all []Issue
	var nextPageToken string
	for {
		params := map[string]interface{}{
			"jql":        jql,
			"maxResults": pageSize,
			"fields":     fields,
		}
		if nextPageToken != "" {
			params["nextPageToken"] = nextPageToken
		}
		body, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		var result searchResult
		if err := c.post("/rest/api/3/search/jql", string(body), &result); err != nil {
			return nil, err
		}
		all = append(all, result.Issues...)
		if result.NextPageToken == "" || len(result.Issues) == 0 {
			break
		}
		nextPageToken = result.NextPageToken
	}
	return all, nil
}

func (c *JiraClient) GetMyIssues() ([]Issue, error) {
	return c.searchAll(
		"assignee = currentUser() ORDER BY key DESC",
		[]string{"summary", "status", "priority", "assignee"},
	)
}

func (c *JiraClient) GetTeamIssues() ([]Issue, error) {
	emails := append([]string{c.cfg.Email}, c.cfg.TeamEmails...)
	if len(emails) == 0 {
		return nil, nil
	}
	quoted := make([]string, len(emails))
	for i, e := range emails {
		quoted[i] = `"` + e + `"`
	}
	jql := fmt.Sprintf("assignee in (%s) ORDER BY key DESC", strings.Join(quoted, ","))
	return c.searchAll(jql, []string{"summary", "status", "priority", "assignee"})
}

func (c *JiraClient) SearchUsers(query string) ([]User, error) {
	path := fmt.Sprintf("/rest/api/3/user/search?query=%s&maxResults=10", url.QueryEscape(query))
	var users []User
	if err := c.get(path, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *JiraClient) AssignIssue(issueKey, accountID string) error {
	payload := fmt.Sprintf(`{"accountId":"%s"}`, accountID)
	return c.put("/rest/api/3/issue/"+issueKey+"/assignee", payload)
}

func (c *JiraClient) GetTransitions(issueKey string) ([]Transition, error) {
	var result transitionsResult
	if err := c.get("/rest/api/3/issue/"+issueKey+"/transitions", &result); err != nil {
		return nil, err
	}
	return result.Transitions, nil
}

func (c *JiraClient) DoTransition(issueKey, transitionID string) error {
	payload := fmt.Sprintf(`{"transition":{"id":"%s"}}`, transitionID)
	return c.post("/rest/api/3/issue/"+issueKey+"/transitions", payload, nil)
}
