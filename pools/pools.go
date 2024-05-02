package pools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://midgard.ninerealms.com/v2/pools"

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

type poolsMsg struct {
	pools []Pool
}

type Pool struct {
	Asset string `json:"asset"`
}

type Model struct {
	pools  []Pool
	Err    error
	cursor int
}

func (m Model) Init() tea.Cmd {
	return fetchPools
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	// var cmds []tea.Cmd
	switch msg := msg.(type) {
	case poolsMsg:
		m.pools = msg.pools
	case errMsg:
		m.Err = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.cursor++
			if m.cursor >= len(m.pools) {
				m.cursor = 0
			}
			return m, nil
		case "k", "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.pools) - 1
			}
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := fmt.Sprintf("Fetching %s...\n", url)
	if m.Err != nil {
		s += fmt.Sprintf("something went wrong: %s", m.Err)
		// } else if m.status != 0 {
		// 	s += fmt.Sprintf("%d %s", m.status, http.StatusText(m.status))
		// }
	} else {

		// log the number of pools there are
		s += fmt.Sprintf("Found %d pools\n", len(m.pools))

		for i, pool := range m.pools {

			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			s += fmt.Sprintf("%s %s\n", cursor, pool.Asset)
		}

	}

	return s + "\n"
}

func fetchPools() tea.Msg {
	resp, err := http.Get("https://midgard.ninerealms.com/v2/pools")
	if err != nil {
		return errMsg{err}
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errMsg{err}
	}

	var pools []Pool
	err = json.Unmarshal(data, &pools)
	if err != nil {
		return errMsg{err}
	}

	return poolsMsg{pools}

}
