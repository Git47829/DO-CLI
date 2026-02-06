package ui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	AppStyle    lipgloss.Style
	BannerStyle lipgloss.Style
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	Cursor      lipgloss.Style
	Item        lipgloss.Style
	ItemDone    lipgloss.Style
	Muted       lipgloss.Style
	Error       lipgloss.Style
	Help        lipgloss.Style
}

func NewTheme() Theme {
	bg := lipgloss.Color("#1A1B26")
	fg := lipgloss.Color("#C0CAF5")
	muted := lipgloss.Color("#565F89")
	blue := lipgloss.Color("#7AA2F7")
	green := lipgloss.Color("#9ECE6A")
	red := lipgloss.Color("#F7768E")
	purple := lipgloss.Color("#BB9AF7")

	return Theme{
		AppStyle:    lipgloss.NewStyle().Background(bg).Foreground(fg).Padding(1, 2),
		BannerStyle: lipgloss.NewStyle().Foreground(blue),
		Title:       lipgloss.NewStyle().Foreground(fg).Bold(true),
		Subtitle:    lipgloss.NewStyle().Foreground(purple),
		Cursor:      lipgloss.NewStyle().Foreground(blue).Bold(true),
		Item:        lipgloss.NewStyle().Foreground(fg),
		ItemDone:    lipgloss.NewStyle().Foreground(green).Strikethrough(true),
		Muted:       lipgloss.NewStyle().Foreground(muted),
		Error:       lipgloss.NewStyle().Foreground(red).Bold(true),
		Help:        lipgloss.NewStyle().Foreground(muted),
	}
}
