package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update controls the changes made at each action
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	player := m.game.CurentPlayer()
	m.actions = player.AvailableActions()

	if m.end {
		return m, tea.Quit
	}
	if str, ok := msg.(tea.KeyMsg); ok {
		// Cool, what was the actual key pressed?
		switch str.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	// flag to know if the player is finished playing
	actionDone := false
	// play bot or human
	if m.game.CurentPlayerIndex() != m.humamIdx {
		if _, ok := msg.(TickEvent); !ok {
			// wait for tick
			return m, nil
		}
		actionDone = m.Bot(msg, player)
	} else {
		// human is playing skip tick event
		if _, ok := msg.(TickEvent); ok {
			return m, nil
		}
		actionDone = m.Human(msg, player)
	}
	// player turn is done
	if actionDone {
		if m.game.CurentPlayer().IsWinner() {
			m.end = true
		}
		player.AddCard(m.game.Deck.Pick())
		m.cursor = 0
		m.game.NextPlayer()
	}
	return m, tick
}
