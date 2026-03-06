package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type resultMsg struct {
	data any
	err  error
}

type SpinnerModel struct {
	spinner spinner.Model
	message string
	done    bool
	err     error
	result  any
	work    func() (any, error)
}

func NewSpinner(message string, work func() (any, error)) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	return SpinnerModel{
		spinner: s,
		message: message,
		work:    work,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.doWork())
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case resultMsg:
		m.done = true
		m.result = msg.data
		m.err = msg.err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m SpinnerModel) View() string {
	if m.done {
		return ""
	}
	return fmt.Sprintf("%s %s\n", m.spinner.View(), m.message)
}

func (m SpinnerModel) Result() (any, error) {
	return m.result, m.err
}

func (m SpinnerModel) doWork() tea.Cmd {
	return func() tea.Msg {
		data, err := m.work()
		return resultMsg{data: data, err: err}
	}
}

func RunWithSpinner(message string, work func() (any, error)) (any, error) {
	model := NewSpinner(message, work)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("running spinner: %w", err)
	}
	return finalModel.(SpinnerModel).Result()
}
