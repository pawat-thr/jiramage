package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func buildLogo() [6]string {
	glyphs := map[rune][6]string{
		'J': {
			`     ██╗`,
			`     ██║`,
			`     ██║`,
			`██   ██║`,
			`╚█████╔╝`,
			` ╚════╝ `,
		},
		'I': {
			`██╗`,
			`██║`,
			`██║`,
			`██║`,
			`██║`,
			`╚═╝`,
		},
		'R': {
			`██████╗ `,
			`██╔══██╗`,
			`██████╔╝`,
			`██╔══██╗`,
			`██║  ██║`,
			`╚═╝  ╚═╝`,
		},
		'A': {
			` █████╗ `,
			`██╔══██╗`,
			`███████║`,
			`██╔══██║`,
			`██║  ██║`,
			`╚═╝  ╚═╝`,
		},
		'M': {
			`███╗   ███╗`,
			`████╗ ████║`,
			`██╔████╔██║`,
			`██║╚██╔╝██║`,
			`██║ ╚═╝ ██║`,
			`╚═╝     ╚═╝`,
		},
		'G': {
			` ██████╗ `,
			`██╔════╝ `,
			`██║  ███╗`,
			`██║   ██║`,
			`╚██████╔╝`,
			` ╚═════╝ `,
		},
		'E': {
			`███████╗`,
			`██╔════╝`,
			`█████╗  `,
			`██╔══╝  `,
			`███████╗`,
			`╚══════╝`,
		},
	}
	var rows [6]string
	for i, ch := range "JIRAMAGE" {
		g, ok := glyphs[ch]
		if !ok {
			continue
		}
		sep := ""
		if i > 0 {
			sep = "  "
		}
		for r := 0; r < 6; r++ {
			rows[r] += sep + g[r]
		}
	}
	return rows
}

func (m Model) splashView() string {
	logo := buildLogo()
	subtitle := `  ──── jira management dashboard ────   `

	logoStyle := lipgloss.NewStyle().Foreground(purple).Bold(true)
	subStyle := lipgloss.NewStyle().Foreground(dimGrey)
	versionStyle := lipgloss.NewStyle().Foreground(teal).Bold(true)
	creditStyle := lipgloss.NewStyle().Foreground(grey)
	hintStyle := lipgloss.NewStyle().Foreground(dimGrey)

	var content strings.Builder
	for _, line := range logo {
		content.WriteString("  " + logoStyle.Render(line) + "\n")
	}
	content.WriteString(subStyle.Render(subtitle) + "\n")
	content.WriteString("\n")
	content.WriteString("  " + versionStyle.Render(appVersion) + "  " + creditStyle.Render(appCredit) + "\n")
	content.WriteString("\n")
	content.WriteString("  " + hintStyle.Render("auto-starting in 2s  •  press enter to skip") + "\n")

	lines := strings.Count(content.String(), "\n") + 1
	topPad := 0
	if m.height > lines+4 {
		topPad = (m.height - lines) / 2
	}
	return strings.Repeat("\n", topPad) + content.String()
}
