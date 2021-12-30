package game

import "math/rand"

// Game holds data related to a game
type Game struct {
	// deck of cards
	Deck *Deck
	// slice holding players
	Players []Player
	// number of current round
	Round int
	// index of playing player
	playerIndex int
	// public infos of all players
	playersInfo *[]PlayerInfo
}

// State data available to all players
type State struct {
	// info available to all players
	PlayersInfo []PlayerInfo
	// number of current round
	Round int
	// index of playing player
	PlayerIndex int
}

// number of cards in the initial hand
const intialCardCount = 3

// Table creates nbPlayers players with random bonuses
func Table(nbPlayers int) ([]Player, *[]PlayerInfo) {
	if nbPlayers < 2 {
		panic("invalid player count must be more than 2")
	}
	if nbPlayers > 4 {
		panic("max player is 4")
	}
	bonuses := BonusList()
	rand.Shuffle(len(bonuses), func(i, j int) {
		bonuses[i], bonuses[j] = bonuses[j], bonuses[i]
	})
	players := make([]Player, nbPlayers)
	// lazy intialisation of players info could be improved
	playersInfo := make([]PlayerInfo, nbPlayers)
	for i := 0; i < nbPlayers; i++ {
		players[i] = *NewPlayer(bonuses[i], &playersInfo)
	}
	return players, &playersInfo
}

// NewGame creates a new game with nbPlayers players.
func NewGame(nbPlayers int) *Game {
	deck := NewDeck()
	deck.Shuffle()
	players, playersInfo := Table(nbPlayers)

	// draw cards
	for i := 0; i < nbPlayers*intialCardCount; i++ {
		players[i%nbPlayers].AddCard(deck.Pick())
	}
	game := Game{
		Deck:        deck,
		Players:     players,
		Round:       0,
		playerIndex: 0,
		playersInfo: playersInfo,
	}
	// update players info after drawing cards
	game.UpdatePlayersInfo()
	return &game
}

func (g *Game) UpdatePlayersInfo() {
	for i := range g.Players {
		(*g.playersInfo)[i] = g.Players[i].Info()
	}
}

// PlayerCount is the number of players in the game
func (g *Game) PlayerCount() int {
	return len(g.Players)
}

// CurentPlayer is player to play
func (g *Game) CurentPlayer() *Player {
	return &g.Players[g.playerIndex]
}

// CurentPlayerIndex index is the index of player to play
func (g *Game) CurentPlayerIndex() int {
	return g.playerIndex
}

// NextPlayer returns the next player to play
func (g *Game) NextPlayer() *Player {
	// update player infos
	g.UpdatePlayersInfo()

	// end game
	if g.Players[g.playerIndex].IsWinner() {
		return nil
	}

	g.playerIndex++

	if g.playerIndex == len(g.Players) {
		g.Round++
		g.playerIndex = 0
	}
	return &g.Players[g.playerIndex]
}

// State return information about the game available to all players
func (g *Game) State() State {
	return State{
		PlayersInfo: *g.playersInfo,
		Round:       g.Round,
		PlayerIndex: g.playerIndex,
	}
}
