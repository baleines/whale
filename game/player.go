package game

import "fmt"

// Bonus is the bonus for one player it's unique
type Bonus int

// All bonuses as enum
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

// String converts enum to string for printing
func (b Bonus) String() string {
	return [...]string{
		"BonusUnknown",
		"BonusPirat",
		"BonusGhost",
		"BonusRound",
		"BonusPiranha",
	}[b]
}

// BonusList give the list of all available bonuses
func BonusList() []Bonus {
	return bonusList[:]
}

// Action represents all possible actions in the game
type Action int

// All possible actions as enum
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

// Player is holding player data
type Player struct {
	// index of player in the game
	ID int
	// Cards in player hand
	Cards []Card
	// Water level of the player
	Water int
	// Type of bonus for player
	BonusType Bonus
	// Where the bonus is played or not
	BonusPlayed bool
	// Public infos from other players
	PlayersInfo *[]PlayerInfo
}

// PlayerInfo represents public player data
type PlayerInfo struct {
	// index of player in the game
	ID int
	// CardCount number of cards of player
	CardCount int
	// Water level of the player
	Water int
	// Type of bonus for player
	BonusType Bonus
	// Where the bonus is played or not
	BonusPlayed bool
}

// NewPlayer creates a player with the given bonus
func NewPlayer(id int, bonus Bonus, playersInfo *[]PlayerInfo) *Player {
	return &Player{
		ID:          id,
		Water:       0,
		Cards:       []Card{},
		BonusType:   bonus,
		BonusPlayed: false,
		PlayersInfo: playersInfo,
	}
}

// String exports player details as string
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

// IsWinner returns true is the player won
func (p *Player) IsWinner() bool {
	return p.Water >= 5
}

// OtherPlayersWithWater returns a slice of ids of players having water
func (p *Player) OtherPlayersWithWater() []int {
	playersWithWater := make([]int, 0)
	for i := range *p.PlayersInfo {
		if (*p.PlayersInfo)[i].Water > 0 && (*p.PlayersInfo)[i].ID != p.ID {
			playersWithWater = append(playersWithWater, (*p.PlayersInfo)[i].ID)
		}
	}
	return playersWithWater
}

// playersFromIDs is an internal method to retrieve players pointers form IDs
func (p *Player) playersFromIDs(ids []int, allPlayers []*Player) []*Player {
	players := make([]*Player, 0)
	for i := range allPlayers {
		for j := range ids {
			if allPlayers[i].ID == ids[j] {
				players = append(players, allPlayers[i])
				break
			}
		}
	}
	if len(players) != len(ids) {
		panic("unexpected result length")
	}
	return players
}

// OtherPlayers returns a slice of ids of other players
func (p *Player) OtherPlayers() []int {
	otherPlayers := make([]int, 0)
	for i := range *p.PlayersInfo {
		if (*p.PlayersInfo)[i].ID != p.ID {
			otherPlayers = append(otherPlayers, (*p.PlayersInfo)[i].ID)
		}
	}
	return otherPlayers
}

// AllPlayers returns a slice of ids of all players ordered
func (p *Player) AllPlayers() []int {
	playersIDs := make([]int, 0)
	for i := range *p.PlayersInfo {
		playersIDs = append(playersIDs, (*p.PlayersInfo)[i].ID)
	}
	return playersIDs
}

// AvailableActions list action available for player
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
	// bonuses Not implemented yet
	if !p.BonusPlayed && p.BonusType != BonusUnknown {
		switch p.BonusType {
		case BonusPirat:
			if len(p.OtherPlayersWithWater()) > 0 {
				actions = append(actions, PlayPirat)
			}
		case BonusGhost:
			actions = append(actions, PlayGhost)
		case BonusRound:
			actions = append(actions, PlayRound)
		case BonusPiranha:
			actions = append(actions, PlayPiranha)
		default:
			fmt.Println("unexpected bonus:", p.BonusType)
			panic("invalid bonus")
		}
	}
	return actions
}

// Play makes the given action and discard cards to deck
// note: there is coupling here we need players to make the actions
func (p *Player) Play(d *Deck, a Action, ids []int, allPlayers []*Player) {
	selectedPlayers := p.playersFromIDs(ids, allPlayers)
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
	case PlayPirat:
		if !p.PlayBonus(BonusPirat, selectedPlayers) {
			panic("invalid action")
		}
	case PlayGhost:
		if !p.PlayBonus(BonusGhost, selectedPlayers) {
			panic("invalid action")
		}
	case PlayRound:
		if !p.PlayBonus(BonusRound, selectedPlayers) {
			panic("invalid action")
		}
	case PlayPiranha:
		if !p.PlayBonus(BonusPiranha, selectedPlayers) {
			panic("invalid action")
		}
	}
}

// AddCard adds the given card to player's hand
func (p *Player) AddCard(c Card) {
	p.Cards = append(p.Cards, c)
}

// GetCard removes the given card from player's hand
// returns false if the card is not in player's hand
func (p *Player) GetCard(c Card) bool {
	for i := range p.Cards {
		if c == p.Cards[i] {
			p.Cards = append(p.Cards[:i], p.Cards[i+1:]...)
			return true
		}
	}
	return false
}

func playBonus(p *Player, b Bonus, targetedPlayers []*Player) bool {
	switch b {
	case BonusPirat:
		p.BonusPlayed = true
		targetedPlayers[0].Water--
		p.Water++
	case BonusGhost:
		p.BonusPlayed = true
		// swap cards
		p.Cards, targetedPlayers[0].Cards = targetedPlayers[0].Cards, p.Cards
	case BonusPiranha:
		p.BonusPlayed = true
		// remove 2 water
		targetedPlayers[0].Water--
		targetedPlayers[1].Water--
		// TODO fix deck add back discarded water and check 0
	case BonusRound:
		p.BonusPlayed = true
		// keep first value
		waters := targetedPlayers[0].Water
		// swap waters form one player with the next one
		for i := 0; i < len(targetedPlayers)-1; i++ {
			targetedPlayers[i].Water = targetedPlayers[i+1].Water
		}
		// get the last player water to the first one
		targetedPlayers[len(targetedPlayers)-1].Water = waters
	default:
		return false
	}
	return true
}

func (p *Player) PlayBonus(b Bonus, targetedPlayers []*Player) bool {
	if p.BonusType != b {
		return false
	}
	return playBonus(p, b, targetedPlayers)
}

// Info return info on player available to all players
func (p *Player) Info() PlayerInfo {
	return PlayerInfo{
		ID:          p.ID,
		CardCount:   len(p.Cards),
		Water:       p.Water,
		BonusType:   p.BonusType,
		BonusPlayed: p.BonusPlayed,
	}
}
