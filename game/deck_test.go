package game

import "testing"

func TestDeck_PickAll(t *testing.T) {
	d := NewDeck()
	cardCount := len(d.remaining)
	for i := 0; i < cardCount; i++ {
		_ = d.Pick()
	}
	if len(d.remaining) != 0 {
		t.Errorf("all cards should have been taken")
	}
	if len(d.discarded) != 0 {
		t.Errorf("no card should have been discarded")
	}
}

func TestDeck_PickAndDiscardAll(t *testing.T) {
	d := NewDeck()
	cardCount := len(d.remaining)
	initialDeck := make([]Card, len(d.remaining))
	copy(initialDeck, d.remaining)
	for i := 0; i < cardCount; i++ {
		c := d.Pick()
		d.Discard(c)
	}
	if len(d.remaining) != 0 {
		t.Errorf("all cards should have been taken")
	}
	if len(d.discarded) != cardCount {
		t.Errorf("all cards should be in discard pile")
	}
	if len(d.discarded) != cardCount {
		t.Errorf("discarded pile should contain the same card count as the initial deck")
	}
	for i := 0; i < cardCount; i++ {
		if d.discarded[i] != initialDeck[i] {
			t.Errorf("discard pile should be the same as initial deck but %d != %d", d.discarded[i], initialDeck[cardCount-i-1])
		}
	}
}

func TestDeck_PickUntilShuffle(t *testing.T) {
	d := NewDeck()
	cardCount := len(d.remaining)
	for i := 0; i < cardCount; i++ {
		c := d.Pick()
		d.Discard(c)
	}

	// Pick one should trigger shuffle
	_ = d.Pick()
	if len(d.remaining) != (cardCount - 1) {
		t.Errorf("all cards should have been taken")
	}
}
