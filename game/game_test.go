package game

import (
	"fmt"
	"testing"
)

func TestGame_PlayerCount(t *testing.T) {
	type fields struct {
		Deck        *Deck
		Players     int
		Round       int
		playerIndex int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "2 players",
			fields: fields{
				Deck:        NewDeck(),
				Players:     2,
				Round:       0,
				playerIndex: 0,
			},
			want: 2,
		},
		{
			name: "3 players",
			fields: fields{
				Deck:        NewDeck(),
				Players:     3,
				Round:       0,
				playerIndex: 0,
			},
			want: 3,
		},
		{
			name: "4 players",
			fields: fields{
				Deck:        NewDeck(),
				Players:     4,
				Round:       0,
				playerIndex: 0,
			},
			want: 4,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			players, playersInfo := Table(tt.fields.Players)
			g := &Game{
				Deck:        tt.fields.Deck,
				Players:     players,
				Round:       tt.fields.Round,
				playerIndex: tt.fields.playerIndex,
				playersInfo: playersInfo,
			}
			if got := g.PlayerCount(); got != tt.want {
				t.Errorf("Game.PlayerCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGame_CurentPlayer(t *testing.T) {
	nbPlayer := 4
	game := NewGame(nbPlayer)
	// NOTE : could use a lib but for now avoid to much deps
	expect := func(num int) {
		if game.CurentPlayerIndex() != num {
			t.Error("unexpected player index")
		}
	}
	for i := 0; i < nbPlayer; i++ {
		expect(i)
		p := game.NextPlayer()
		if p == nil {
			t.Error("player should not be nil")
		}
		p = game.CurentPlayer()
		if p == nil {
			t.Error("player should not be nil")
		}
	}
	// should be back to first player
	expect(0)
}

func TestGame_State(t *testing.T) {
	nbPlayer := 4
	game := NewGame(nbPlayer)

	expect := func(p *Player, pI PlayerInfo) {
		if len(p.Cards) != pI.CardCount {
			t.Error("unexpected state card count", p, pI)
		}
		if p.Water != pI.Water {
			t.Error("unexpected state water", p, pI)
		}
	}
	for i := 0; i < nbPlayer; i++ {
		p := game.CurentPlayer()
		pI := game.State().PlayersInfo[i]
		expect(p, pI)
		_ = game.NextPlayer()
	}
}

func TestGame_EntireGame(t *testing.T) {
	nbPlayer := 4
	game := NewGame(nbPlayer)
	roundLimit := 1000

	// NOTE : could use a lib but for now avoid to much deps
	for !game.CurentPlayer().IsWinner() {
		p := game.CurentPlayer()
		// always pick the first action
		a := p.AvailableActions()[0]
		fmt.Println(a)
		p.Play(game.Deck, a, nil)
		p.AddCard(game.Deck.Pick())
		// always pick the first action
		p = game.NextPlayer()
		if p == nil {
			// game finished
			break
		}
		if game.Round > roundLimit {
			t.Error("Game did not finish")
			break
		}
		for _, p := range game.Players {
			fmt.Println(p)
		}
	}

}
