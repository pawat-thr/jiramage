package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"jira-dashboard/config"
)

type JiraClient struct {
	cfg  *config.Config
	http *http.Client
}

func NewJiraClient(cfg *config.Config) *JiraClient {
	return &JiraClient{
		cfg:  cfg,
		http: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *JiraClient) BaseURL() string {
	return c.cfg.JiraURL
}

func (c *JiraClient) RefreshInterval() time.Duration {
	return c.cfg.RefreshInterval
}

func (c *JiraClient) authHeader() string {
	raw := c.cfg.Email + ":" + c.cfg.Token
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}

func (c *JiraClient) get(path string, out interface{}) error {
	req, err := http.NewRequest("GET", c.cfg.JiraURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.authHeader())
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return json.Unmarshal(body, out)
}

func (c *JiraClient) post(path string, payload string, out interface{}) error {
	req, err := http.NewRequest("POST", c.cfg.JiraURL+path, strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.authHeader())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	if out != nil {
		return json.Unmarshal(body, out)
	}
	return nil
}

func (c *JiraClient) put(path string, payload string) error {
	req, err := http.NewRequest("PUT", c.cfg.JiraURL+path, strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.authHeader())
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
