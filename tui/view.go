package tui

import (
	"fmt"
)

// View is printing the game status to the terminal
func (m model) View() string {
	player := m.game.CurentPlayer()
	// The header
	s := blue("~WHALE GAME") + "     (q to quit)\n"
	for lvl := 5; lvl > 0; lvl-- {
		for i, p := range m.game.Players {
			w := ""
			// spacing between players
			if i > 0 {
				w += " "
			}
			if p.Water >= lvl {
				switch lvl {
				case 5:
					w += `C\    /Ↄ`
				case 4:
					w += `  \  /  `
				default:
					w += "   ||   "
				}
			} else {
				w += "        "
			}
			s += ColorizeWater(lvl, w)
		}
		s += "\n"
	}

	for i := range m.game.Players {
		s += ColorizeWhale(i, ` (॰  \_/ `)
	}
	s += "\n"
	// show the active player in white
	for i := range m.game.Players {
		pNum := fmt.Sprintf("   P%d    ", i)
		if i == m.game.CurentPlayerIndex() {
			s += pNum
		} else {
			s += grey(pNum)
		}
	}
	s += "\n\n"

	if m.end {
		s += "        * * *\n"
		s += fmt.Sprintf("    PLAYER %d WINS !\n", m.game.CurentPlayerIndex())
		s += "        * * *\n"
	} else if m.game.CurentPlayerIndex() == m.humamIdx {
		s += fmt.Sprintf("Player:%d ", m.game.CurentPlayerIndex())
		s += player.String()
		s += "\n"
		// Iterate over our actions
		m.actions = player.AvailableActions()
		for i, choice := range m.actions {
			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			// Render the row
			s += fmt.Sprintf("%s %s\n", cursor, choice.String())
		}
	}

	// The footer
	for i := 0; i < m.game.PlayerCount(); i++ {
		s += lightBlue("~ ") + blue("~ ") + verylightBlue("~ ") + darkBlue("~  ")
	}
	s += "\n"

	// Send the UI for rendering
	return s
}
