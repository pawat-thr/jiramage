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
	ProjectKeys     []string
	LabelKeys       []string
	FixedStatus     string
	TeamFromDate    string
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
	if raw := os.Getenv("JIRA_PROJECT"); raw != "" {
		for _, k := range strings.Split(raw, ",") {
			k = strings.TrimSpace(strings.ToUpper(k))
			if k != "" {
				cfg.ProjectKeys = append(cfg.ProjectKeys, k)
			}
		}
	}
	cfg.FixedStatus = "SIT DEPLOYED"
	if raw := os.Getenv("JIRA_FIXED_STATUS"); raw != "" {
		cfg.FixedStatus = strings.TrimSpace(raw)
	}
	cfg.TeamFromDate = "2024-05-01"
	if raw := os.Getenv("JIRA_TEAM_FROM"); raw != "" {
		cfg.TeamFromDate = strings.TrimSpace(raw)
	}
	if raw := os.Getenv("JIRA_LABELS"); raw != "" {
		for _, k := range strings.Split(raw, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				cfg.LabelKeys = append(cfg.LabelKeys, k)
			}
		}
	}
	return cfg, nil
}
