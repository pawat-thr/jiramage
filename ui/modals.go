package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) statusFilterModal(title string, options []string, cursor int, active string) string {
	var s string
	s += titleStyle.Render(title) + "\n\n"

	check := " "
	if active == "" {
		check = "✓"
	}
	if cursor == 0 {
		s += selStyle.Render("▶ "+fmt.Sprintf("%-22s", "All")) + "\n"
	} else {
		s += fmt.Sprintf("  %s %-22s\n", check, "All")
	}

	for i, opt := range options {
		sColor := statusColor[opt]
		if sColor == "" {
			sColor = grey
		}
		label := lipgloss.NewStyle().Foreground(sColor).Bold(true).Render(opt)
		check = " "
		if active == opt {
			check = "✓"
		}
		if cursor == i+1 {
			s += selStyle.Render("▶ "+fmt.Sprintf("%-22s", opt)) + "\n"
		} else {
			s += fmt.Sprintf("  %s %s\n", check, label)
		}
	}
	s += "\n" + dimStyle.Render("↑↓ navigate   enter apply   esc cancel")
	return modalBorder.Render(s) + "\n\n"
}

func (m Model) nameFilterModal() string {
	var s string
	s += titleStyle.Render("Filter by Name") + "\n\n"

	check := " "
	if m.teamNameFilter == "" {
		check = "✓"
	}
	if m.teamNameFilterCursor == 0 {
		s += selStyle.Render("▶ "+fmt.Sprintf("%-22s", "All members")) + "\n"
	} else {
		s += fmt.Sprintf("  %s %-22s\n", check, "All members")
	}

	for i, opt := range m.teamNameOptions {
		nick := shortName(opt)
		color := memberColors[i%len(memberColors)]
		label := lipgloss.NewStyle().Foreground(color).Bold(true).Render(nick)
		check = " "
		if m.teamNameFilter == opt {
			check = "✓"
		}
		if m.teamNameFilterCursor == i+1 {
			s += selStyle.Render("▶ "+fmt.Sprintf("%-22s", nick)) + "\n"
		} else {
			s += fmt.Sprintf("  %s %s\n", check, label)
		}
	}
	s += "\n" + dimStyle.Render("↑↓ navigate   enter apply   esc cancel")
	return modalBorder.Render(s) + "\n\n"
}

func (m Model) reassignSearchView() string {
	visible := m.currentVisible()
	if len(visible) == 0 {
		return ""
	}
	iss := visible[m.currentCursor()]
	return modalBorder.Render(fmt.Sprintf("%s\n\n%s\n\n%s",
		titleStyle.Render("Reassign  "+iss.Key),
		m.textInput.View(),
		dimStyle.Render("enter search   esc cancel"),
	)) + "\n\n"
}

func (m Model) reassignPickView() string {
	visible := m.currentVisible()
	if len(visible) == 0 {
		return ""
	}
	iss := visible[m.currentCursor()]
	var s string
	s += titleStyle.Render("Assign  "+iss.Key+"  →  select user") + "\n\n"
	if len(m.users) == 0 {
		s += dimStyle.Render("No users found.") + "\n"
	} else {
		for i, u := range m.users {
			label := fmt.Sprintf("%-28s  %s", u.DisplayName, u.EmailAddress)
			if i == m.userCursor {
				s += selStyle.Render("▶ "+label) + "\n"
			} else {
				s += "  " + label + "\n"
			}
		}
	}
	s += "\n" + dimStyle.Render("↑↓ navigate   enter assign   esc cancel")
	return modalBorder.Render(s) + "\n\n"
}

func (m Model) transitionPickView() string {
	visible := m.currentVisible()
	if len(visible) == 0 {
		return ""
	}
	iss := visible[m.currentCursor()]
	var s string
	s += titleStyle.Render("Move  "+iss.Key+"  →  select status") + "\n\n"
	if len(m.transitions) == 0 {
		s += dimStyle.Render("No transitions available.") + "\n"
	} else {
		for i, t := range m.transitions {
			sColor := statusColor[t.Name]
			if sColor == "" {
				sColor = grey
			}
			if i == m.transitionCursor {
				s += selStyle.Render("▶ "+fmt.Sprintf("%-22s", t.Name)) + "\n"
			} else {
				s += "  " + lipgloss.NewStyle().Foreground(sColor).Render(t.Name) + "\n"
			}
		}
	}
	s += "\n" + dimStyle.Render("↑↓ navigate   enter apply   esc cancel")
	return modalBorder.Render(s) + "\n\n"
}
