package tui

import (
	"whale/game"

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
	// only human player is playing
	if m.game.CurentPlayerIndex() != m.humamIdx {
		if _, ok := msg.(TickEvent); ok {
			// always plays the first available action
			player.Play(m.game.Deck, m.actions[0], nil)
			if m.game.CurentPlayer().IsWinner() {
				m.end = true
			}
			player.AddCard(m.game.Deck.Pick())
			m.cursor = 0
			m.game.NextPlayer()
			return m, tick
		}
		// wait for tick
		return m, nil
	}

	// human is playing skip tick event
	if _, ok := msg.(TickEvent); ok {
		return m, nil
	}
	if str, ok := msg.(tea.KeyMsg); ok {
		// Cool, what was the actual key pressed?
		switch str.String() {

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.actions)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			selectMode := 0
			if len(m.actions) > 0 {
				// case of bonuses
				actionsSelect1 := []game.Action{game.PlayPirat, game.PlayGhost}
				actionsSelect2 := []game.Action{game.PlayPiranha}
				for _, a := range actionsSelect1 {
					if m.actions[m.cursor] == a {
						selectMode = 1
					}
				}
				for _, a := range actionsSelect2 {
					if m.actions[m.cursor] == a {
						selectMode = 2
					}
				}
				if selectMode == 0 {
					player.Play(m.game.Deck, m.actions[m.cursor], nil)
				} else {
					switch selectMode {
					case 1:
						if m.actions[m.cursor] == game.PlayPirat {
							m.selectedPlayers = player.OtherPlayersWithWater()
						} else {
							m.selectedPlayers = player.OtherPlayers()
						}
						player.Play(m.game.Deck, m.actions[m.cursor],
							[]*game.Player{&m.game.Players[m.selectedPlayers[0]]})
					case 2:
						m.selectedPlayers = player.OtherPlayers()
						player.Play(m.game.Deck, m.actions[m.cursor],
							[]*game.Player{
								&m.game.Players[m.selectedPlayers[0]],
								&m.game.Players[m.selectedPlayers[1]],
							})
					default:
						panic("unexpected select mode")
					}
				}
			}
			player.AddCard(m.game.Deck.Pick())
			if m.game.NextPlayer() == nil {
				m.end = true
			}
			m.cursor = 0
		}
	}
	return m, tick
}
