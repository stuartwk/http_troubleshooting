package main

import (
	"fmt"
	"os"

	"github.com/stuartwk/http_troubleshooting/pools"

	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	poolsView sessionState = iota
	poolView
)

type model struct {
	state      sessionState
	poolsModel pools.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	switch m.state {
	case poolsView:
		newPool, newCmd := m.poolsModel.Update(msg)
		poolModel, ok := newPool.(pools.Model)
		if !ok {
			panic("could not asset pools.Model")
		}
		m.poolsModel = poolModel
		cmd = newCmd
		// return m, cmd
	default:
		return m, nil
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state {
	case poolsView:
		return m.poolsModel.View()
	default:
		return m.poolsModel.View()
	}
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
