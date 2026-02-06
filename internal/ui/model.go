package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"mathOps/internal/todo"
)

type mode int

const (
	modeList mode = iota
	modeAddName
	modeAddDesc
	modeEditName
	modeEditDesc
)

type model struct {
	store     todo.Store
	todos     []todo.ToDo
	cursor    int
	mode      mode
	nameInput textinput.Model
	descInput textinput.Model
	status    string
	theme     Theme
	width     int
	height    int
}

func newModel(store todo.Store, todos []todo.ToDo) model {
	name := textinput.New()
	name.Placeholder = "Task name"
	name.Focus()
	name.CharLimit = 80
	name.Width = 40

	desc := textinput.New()
	desc.Placeholder = "Description"
	desc.CharLimit = 200
	desc.Width = 60

	return model{
		store:     store,
		todos:     todos,
		cursor:    0,
		mode:      modeList,
		nameInput: name,
		descInput: desc,
		status:    "",
		theme:     NewTheme(),
		width:     0,
		height:    0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.mode == modeAddName || m.mode == modeAddDesc || m.mode == modeEditName || m.mode == modeEditDesc {
			if msg.String() == "enter" || msg.String() == "esc" {
				return m.handleKey(msg)
			}
			if m.mode == modeAddName || m.mode == modeEditName {
				var cmd tea.Cmd
				m.nameInput, cmd = m.nameInput.Update(msg)
				return m, cmd
			}
			var cmd tea.Cmd
			m.descInput, cmd = m.descInput.Update(msg)
			return m, cmd
		}
		return m.handleKey(msg)
	}

	if m.mode == modeAddName || m.mode == modeEditName {
		var cmd tea.Cmd
		m.nameInput, cmd = m.nameInput.Update(msg)
		return m, cmd
	}

	if m.mode == modeAddDesc || m.mode == modeEditDesc {
		var cmd tea.Cmd
		m.descInput, cmd = m.descInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch m.mode {
	case modeList:
		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.todos)-1 {
				m.cursor++
			}
		case "a":
			m.status = ""
			m.nameInput.SetValue("")
			m.descInput.SetValue("")
			m.mode = modeAddName
			m.nameInput.Focus()
			m.descInput.Blur()
		case "e":
			if len(m.todos) == 0 {
				m.status = "No tasks to edit."
				break
			}
			current := m.todos[m.cursor]
			m.status = ""
			m.nameInput.SetValue(current.Name)
			m.descInput.SetValue(current.Description)
			m.mode = modeEditName
			m.nameInput.Focus()
			m.descInput.Blur()
		case "c":
			m = toggleTodo(m)
		case "d":
			m = deleteTodo(m)
		}
	case modeAddName:
		switch key {
		case "esc":
			m.mode = modeList
			m.nameInput.Blur()
		case "enter":
			if m.nameInput.Value() == "" {
				m.status = "Name is required."
				break
			}
			m.mode = modeAddDesc
			m.nameInput.Blur()
			m.descInput.Focus()
		}
	case modeAddDesc:
		switch key {
		case "esc":
			m.mode = modeList
			m.descInput.Blur()
		case "enter":
			m = createTodo(m)
		}
	case modeEditName:
		switch key {
		case "esc":
			m.mode = modeList
			m.nameInput.Blur()
		case "enter":
			if m.nameInput.Value() == "" {
				m.status = "Name is required."
				break
			}
			m.mode = modeEditDesc
			m.nameInput.Blur()
			m.descInput.Focus()
		}
	case modeEditDesc:
		switch key {
		case "esc":
			m.mode = modeList
			m.descInput.Blur()
		case "enter":
			m = updateTodo(m)
		}
	}

	return m, nil
}

func toggleTodo(m model) model {
	if len(m.todos) == 0 {
		m.status = "No tasks to toggle."
		return m
	}

	m.todos[m.cursor].IsDone = !m.todos[m.cursor].IsDone
	if err := m.store.Save(m.todos); err != nil {
		m.status = fmt.Sprintf("Save failed: %v", err)
		return m
	}

	m.status = "Updated task."
	return m
}

func deleteTodo(m model) model {
	if len(m.todos) == 0 {
		m.status = "No tasks to delete."
		return m
	}

	m.todos = append(m.todos[:m.cursor], m.todos[m.cursor+1:]...)
	if m.cursor >= len(m.todos) && m.cursor > 0 {
		m.cursor--
	}

	if err := m.store.Save(m.todos); err != nil {
		m.status = fmt.Sprintf("Save failed: %v", err)
		return m
	}

	m.status = "Deleted task."
	return m
}

func createTodo(m model) model {
	name := m.nameInput.Value()
	desc := m.descInput.Value()

	if name == "" {
		m.status = "Name is required."
		return m
	}

	m.todos = append(m.todos, todo.New(name, desc))
	if err := m.store.Save(m.todos); err != nil {
		m.status = fmt.Sprintf("Save failed: %v", err)
		return m
	}

	m.status = "Added task."
	m.mode = modeList
	m.descInput.Blur()
	return m
}

func updateTodo(m model) model {
	if len(m.todos) == 0 {
		m.status = "No tasks to edit."
		m.mode = modeList
		m.descInput.Blur()
		return m
	}

	name := m.nameInput.Value()
	desc := m.descInput.Value()

	if name == "" {
		m.status = "Name is required."
		return m
	}

	m.todos[m.cursor].Name = name
	m.todos[m.cursor].Description = desc

	if err := m.store.Save(m.todos); err != nil {
		m.status = fmt.Sprintf("Save failed: %v", err)
		return m
	}

	m.status = "Updated task."
	m.mode = modeList
	m.descInput.Blur()
	return m
}

func (m model) View() string {
	banner := m.theme.BannerStyle.Render(Banner())
	title := m.theme.Subtitle.Render("Cappuccin Tokyo Night")
	body := m.renderBody()
	status := m.renderStatus()
	footer := m.renderFooter()

	content := lipgloss.JoinVertical(lipgloss.Left, banner, title, body, status, footer)
	return m.theme.AppStyle.Render(content)
}

func (m model) renderBody() string {
	if m.mode == modeAddName {
		label := m.theme.Title.Render("New task")
		return lipgloss.JoinVertical(lipgloss.Left, label, m.nameInput.View(), m.theme.Muted.Render("Enter: next  Esc: cancel"))
	}

	if m.mode == modeAddDesc {
		label := m.theme.Title.Render("Description")
		return lipgloss.JoinVertical(lipgloss.Left, label, m.descInput.View(), m.theme.Muted.Render("Enter: save  Esc: cancel"))
	}

	if m.mode == modeEditName {
		label := m.theme.Title.Render("Edit task name")
		return lipgloss.JoinVertical(lipgloss.Left, label, m.nameInput.View(), m.theme.Muted.Render("Enter: next  Esc: cancel"))
	}

	if m.mode == modeEditDesc {
		label := m.theme.Title.Render("Edit description")
		return lipgloss.JoinVertical(lipgloss.Left, label, m.descInput.View(), m.theme.Muted.Render("Enter: save  Esc: cancel"))
	}

	if len(m.todos) == 0 {
		return m.theme.Muted.Render("No tasks yet. Press 'a' to add one.")
	}

	lines := make([]string, 0, len(m.todos))
	for i, t := range m.todos {
		cursor := "  "
		if i == m.cursor {
			cursor = m.theme.Cursor.Render("> ")
		}

		status := "[ ]"
		style := m.theme.Item
		if t.IsDone {
			status = "[x]"
			style = m.theme.ItemDone
		}

		line := fmt.Sprintf("%s%s %s", cursor, status, t.Name)
		if t.Description != "" {
			line = fmt.Sprintf("%s - %s", line, t.Description)
		}

		lines = append(lines, style.Render(line))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m model) renderStatus() string {
	if m.status == "" {
		return ""
	}

	lower := strings.ToLower(m.status)
	if strings.Contains(lower, "failed") || strings.Contains(lower, "error") {
		return m.theme.Error.Render(m.status)
	}

	return m.theme.Muted.Render(m.status)
}

func (m model) renderFooter() string {
	return m.theme.Help.Render("a add  e edit  c toggle  d delete  q quit")
}
