package tui

import (
	"math/rand"
	"time"
	"whale/whale"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	game     *whale.Game
	actions  []whale.Action // items on the to-do list
	cursor   int            // which to-do list item our cursor is pointing at
	end      bool           // indicates end of game
	humamIdx int            // indicate the index of the human playing
}

// NewWhale returns a model for a one player whale game
func NewWhale() tea.Model {
	rand.Seed(time.Now().UnixNano())
	const nbPlayers = 4
	humanIdx := rand.Int() % nbPlayers
	return &model{
		game:     whale.NewGame(nbPlayers),
		actions:  []whale.Action{},
		cursor:   0,
		end:      false,
		humamIdx: humanIdx,
	}

}

// TickEvent indicates that the timer has ticked.
type TickEvent time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return TickEvent{}
}

// Init initialize the model
func (m model) Init() tea.Cmd {
	// tick is returned in case no user action is needed next
	return tick
}
