package tui

import (
	"whale/game"

	tea "github.com/charmbracelet/bubbletea"
)

// Human gives human controls returns true if player turn is finished
func (m *model) Human(msg tea.Msg, player *game.Player) bool {
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
			return true
		}
	}
	return false
}
