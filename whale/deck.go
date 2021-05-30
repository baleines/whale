package whale

import "math/rand"

// Deck represent the pile and discard pile
type Deck struct {
	remaining []Card
	discarded []Card
}

const (
	water       = 22
	waves       = 8
	doubleWaves = 4
)

// NewDeck creates a new deck base one cards distribution (const)
func NewDeck() *Deck {
	d := &Deck{}
	for i := 0; i < water; i++ {
		d.remaining = append(d.remaining, Water)
	}
	for i := 0; i < waves; i++ {
		d.remaining = append(d.remaining, Wave)
	}
	for i := 0; i < doubleWaves; i++ {
		d.remaining = append(d.remaining, DoubleWave)
	}
	// initialize slice to avoid nil vs empty confusion
	d.discarded = []Card{}
	return d
}

// Pick a card from remaining cards
func (d *Deck) Pick() Card {
	if len(d.remaining) == 0 {
		d.Shuffle()
	}
	if len(d.remaining) == 0 {
		panic("never suppose to pick on empty deck")
	}
	card := d.remaining[0]
	d.remaining = d.remaining[1:]
	return card
}

// Discard adds the card to the discard pile
func (d *Deck) Discard(c Card) {
	if c == UnknownCard {
		panic("should never discard unknown card")
	}
	d.discarded = append(d.discarded, c)
}

// Shuffle set the cards in random order and adds the discard pile if no more cards
func (d *Deck) Shuffle() {
	// get back remaining cards if needed
	if len(d.remaining) == 0 {
		d.remaining, d.discarded = append(d.remaining, d.discarded...), []Card{}
	}
	rand.Shuffle(len(d.remaining), func(i, j int) {
		d.remaining[i], d.remaining[j] = d.remaining[j], d.remaining[i]
	})
}
