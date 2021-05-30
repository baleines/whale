package whale

import "fmt"

type Bonus int

const (
	BonusUnknown Bonus = iota
	BonusPirat
	BonusGhost
	BonusRound
	BonusPiranha
)

var bonusList = [...]Bonus{
	BonusPirat,
	BonusGhost,
	BonusRound,
	BonusPiranha,
}

func BonusList() []Bonus {
	return bonusList[:]
}

type Action int

const (
	UnknownAction Action = iota
	Skip
	PutWater
	PutWaterDouble
	PutTwoWater
	PlayPirat
	PlayGhost
	PlayRound
	PlayPiranha
)

// String converts enum to string for printing
func (a Action) String() string {
	return [...]string{
		"UnknownAction",
		"..zZZz..",
		" ~ + ||",
		" ≈ + ||",
		" ≈ + || + ||",
		"PlayPirat",
		"PlayGhost",
		"PlayRound",
		"PlayPiranha",
	}[a]
}

var actionList = [...]Action{
	Skip,
	PutWater,
	PutWaterDouble,
	PutTwoWater,
	PlayPirat,
	PlayGhost,
	PlayRound,
	PlayPiranha,
}

// ActionList return the list of all actions
func ActionList() []Action {
	return actionList[:]
}

type Player struct {
	// Water level of the player
	Water int
	// Cards in player hand
	Cards []Card
	// Type of bonus for player
	BonusType Bonus
	// Where the bonus is played or not
	BonusPlayed bool
}

func NewPlayer(bonus Bonus) *Player {
	return &Player{
		Water:       0,
		Cards:       []Card{},
		BonusType:   bonus,
		BonusPlayed: false,
	}
}

func (p *Player) String() string {
	s := fmt.Sprintf("Water:%d Cards:[", p.Water)
	for i, c := range p.Cards {
		if i > 0 {
			s += ","
		}
		s += c.String()
	}
	s += "]"
	return s
}

func (p *Player) IsWinner() bool {
	return p.Water >= 5
}

func (p *Player) AvailableActions() []Action {
	actions := []Action{}
	// actions
	hasSingle, hasDouble := false, false
	var water int
	for _, card := range p.Cards {

		switch card {
		case Wave:
			hasSingle = true
		case DoubleWave:
			hasDouble = true
		case Water:
			water++
		}
	}
	if hasDouble && (water >= 2) {
		actions = append(actions, PutTwoWater)
	}
	if hasDouble && (water >= 1) {
		actions = append(actions, PutWaterDouble)
	}
	if hasSingle && (water >= 1) {
		actions = append(actions, PutWater)
	}
	// Skip is always available
	// But always before bonuses
	// This ensures always doing the first action finishes the game
	actions = append(actions, Skip)
	// bonus
	if !p.BonusPlayed && p.BonusType != BonusUnknown {
		var action Action
		switch p.BonusType {
		case BonusPirat:
			action = PlayPirat
		case BonusGhost:
			action = PlayGhost
		case BonusRound:
			action = PlayRound
		case BonusPiranha:
			action = PlayPiranha
		}
		actions = append(actions, action)
	}
	return actions
}

func (p *Player) Play(d *Deck, a Action) {
	switch a {
	case Skip:
		// no op
	case PutWater:
		if !(p.GetCard(Wave)) {
			panic("invalid action")
		}
		d.Discard(Wave)
		if !(p.GetCard(Water)) {
			panic("invalid action")
		}
		p.Water++
	case PutWaterDouble:
		if !(p.GetCard(DoubleWave)) {
			panic("invalid action")
		}
		d.Discard(DoubleWave)
		if !(p.GetCard(Water)) {
			panic("invalid action")
		}
		p.Water++
	case PutTwoWater:
		if !(p.GetCard(Water)) {
			panic("invalid action")
		}
		if !(p.GetCard(Water)) {
			panic("invalid action")
		}
		p.Water += 2
		d.Discard(DoubleWave)
	}
}

func (p *Player) AddCard(c Card) {
	p.Cards = append(p.Cards, c)
}

func (p *Player) GetCard(c Card) bool {
	for i := range p.Cards {
		if c == p.Cards[i] {
			p.Cards = append(p.Cards[:i], p.Cards[i+1:]...)
			return true
		}
	}
	return false
}
