package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"mathOps/internal/todo"
)

func Run(store todo.Store) error {
	todos, err := store.Load()
	if err != nil {
		return fmt.Errorf("load todos: %w", err)
	}

	m := newModel(store, todos)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("run tui: %w", err)
	}

	return nil
}
