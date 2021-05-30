package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"whale/whale"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type model struct {
	game    *whale.Game
	actions []whale.Action // items on the to-do list
	cursor  int            // which to-do list item our cursor is pointing at
	end     bool           // indicates end of game
}

var (
	term = termenv.ColorProfile()
)

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type TickEvent time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return TickEvent{}
}

func (m model) Init() tea.Cmd {
	return tick
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Bold().Styled
}

var (
	VeryLightBlue = makeFgStyle("45")
	LightBlue     = makeFgStyle("39")
	Blue          = makeFgStyle("27")
	DarkBlue      = makeFgStyle("20")
	Grey          = makeFgStyle("#888888")
)

// ColorizeWhale the whale base on index
func ColorizeWhale(i int, s string) string {
	switch i {
	case 0:
		return VeryLightBlue(s)
	case 1:
		return LightBlue(s)
	case 2:
		return Blue(s)
	case 3:
		return DarkBlue(s)
	default:
		return Grey(s)
	}
}

// ColorizeWater the water base on index
func ColorizeWater(i int, s string) string {
	switch i {
	case 0:
		return DarkBlue(s)
	case 1:
		return Blue(s)
	case 2:
		return LightBlue(s)
	case 3:
		return VeryLightBlue(s)
	case 4:
		return Grey(s)
	default:
		return s
	}
}

func (m model) View() string {
	player := m.game.CurentPlayer()
	// The header
	s := Blue("~WHALE GAME") + "     (q to quit)\n"
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
	for i := range m.game.Players {
		s += Grey(fmt.Sprintf("   P%d    ", i))
	}
	s += "\n\n"

	if m.end {
		s += "        * * *\n"
		s += fmt.Sprintf("    PLAYER %d WINS !\n", m.game.CurentPlayerIndex())
		s += "        * * *\n"
	} else {
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
		s += LightBlue("~ ") + Blue("~ ") + VeryLightBlue("~ ") + DarkBlue("~  ")
	}
	s += "\n"

	// Send the UI for rendering
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	player := m.game.CurentPlayer()
	m.actions = player.AvailableActions()

	if m.end {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

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
			if len(m.actions) > 0 {
				player.Play(m.game.Deck, m.actions[m.cursor])
			}
			player.AddCard(m.game.Deck.Pick())
			if m.game.NextPlayer() == nil {
				m.end = true
			}
			m.cursor = 0
		}
	}
	return m, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var initialModel = model{
		game: whale.NewGame(4),
	}

	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
