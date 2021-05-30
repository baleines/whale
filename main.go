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
	game     *whale.Game
	actions  []whale.Action // items on the to-do list
	cursor   int            // which to-do list item our cursor is pointing at
	end      bool           // indicates end of game
	humamIdx int            // indicate the index of the human playing
}

var (
	term = termenv.ColorProfile()
)

// TickEvent indicates that the timer has ticked.
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
	verylightBlue = makeFgStyle("45")
	lightBlue     = makeFgStyle("39")
	blue          = makeFgStyle("27")
	darkBlue      = makeFgStyle("20")
	grey          = makeFgStyle("#888888")
)

// ColorizeWhale the whale base on index
func ColorizeWhale(i int, s string) string {
	switch i {
	case 0:
		return verylightBlue(s)
	case 1:
		return lightBlue(s)
	case 2:
		return blue(s)
	case 3:
		return darkBlue(s)
	default:
		return grey(s)
	}
}

// ColorizeWater the water base on index
func ColorizeWater(i int, s string) string {
	switch i {
	case 0:
		return darkBlue(s)
	case 1:
		return blue(s)
	case 2:
		return lightBlue(s)
	case 3:
		return verylightBlue(s)
	case 4:
		return grey(s)
	default:
		return s
	}
}

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
			player.Play(m.game.Deck, m.actions[0])
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
	return m, tick
}

func main() {
	rand.Seed(time.Now().UnixNano())
	const nbPlayers = 4
	humanIdx := rand.Int() % nbPlayers
	var initialModel = model{
		game:     whale.NewGame(nbPlayers),
		actions:  []whale.Action{},
		cursor:   0,
		end:      false,
		humamIdx: humanIdx,
	}

	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
