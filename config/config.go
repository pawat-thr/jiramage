package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"time"
)

type Config struct {
	JiraURL         string
	Email           string
	Token           string
	TeamEmails      []string
	RefreshInterval time.Duration
}

func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.Trim(strings.TrimSpace(v), `"'`)
		if k != "" && os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}

func LoadConfig() (*Config, error) {
	loadDotEnv(".env")

	cfg := &Config{
		JiraURL: os.Getenv("JIRA_URL"),
		Email:   os.Getenv("JIRA_EMAIL"),
		Token:   os.Getenv("JIRA_TOKEN"),
	}
	if cfg.JiraURL == "" || cfg.Email == "" || cfg.Token == "" {
		return nil, errors.New("missing config: set JIRA_URL, JIRA_EMAIL, JIRA_TOKEN in .env")
	}
	cfg.RefreshInterval = 5 * time.Minute
	if raw := os.Getenv("REFRESH_INTERVAL"); raw != "" {
		if d, err := time.ParseDuration(raw); err == nil && d > 0 {
			cfg.RefreshInterval = d
		}
	}
	if raw := os.Getenv("TEAM_EMAILS"); raw != "" {
		for _, e := range strings.Split(raw, ",") {
			e = strings.TrimSpace(e)
			if e != "" {
				cfg.TeamEmails = append(cfg.TeamEmails, e)
			}
		}
	}
	return cfg, nil
}
