package game

import "testing"

func TestCards(t *testing.T) {
	// test that each cards can be printed
	for _, c := range CardList() {
		if len(c.String()) == 0 {
			t.Errorf("Card %d has no valid representation", c)
		}
	}
}
