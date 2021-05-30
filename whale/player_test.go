package whale

import (
	"reflect"
	"testing"
)

func TestActions(t *testing.T) {
	// test that each action can be printed
	for _, a := range ActionList() {
		if len(a.String()) == 0 {
			t.Errorf("Action %d has no valid representation", a)
		}
	}
}

func TestPlayer_AvailableActions(t *testing.T) {
	type fields struct {
		Water       int
		Cards       []Card
		BonusType   Bonus
		BonusPlayed bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []Action
	}{
		{name: "no actions no cards",
			fields: fields{
				Water:       0,
				Cards:       []Card{},
				BonusType:   BonusGhost,
				BonusPlayed: true,
			},
			want: []Action{Skip},
		},
		{name: "no actions 3 water",
			fields: fields{
				Water:       0,
				Cards:       []Card{Water, Water, Water},
				BonusType:   BonusGhost,
				BonusPlayed: true,
			},
			want: []Action{Skip},
		},
		{name: "PutWater",
			fields: fields{
				Water:       0,
				Cards:       []Card{Wave, Water},
				BonusType:   BonusGhost,
				BonusPlayed: true,
			},
			want: []Action{PutWater, Skip},
		},
		{name: "PutWaterDouble",
			fields: fields{
				Water:       0,
				Cards:       []Card{DoubleWave, Water},
				BonusType:   BonusGhost,
				BonusPlayed: true,
			},
			want: []Action{PutWaterDouble, Skip},
		},
		{name: "PutWaterDouble & PutTwoWater",
			fields: fields{
				Water:       0,
				Cards:       []Card{DoubleWave, Water, Water},
				BonusType:   BonusGhost,
				BonusPlayed: true,
			},
			want: []Action{PutTwoWater, PutWaterDouble, Skip},
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				Water:       tt.fields.Water,
				Cards:       tt.fields.Cards,
				BonusType:   tt.fields.BonusType,
				BonusPlayed: tt.fields.BonusPlayed,
			}
			if got := p.AvailableActions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Player.AvailableActions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_Play(t *testing.T) {
	type fields struct {
		Water       int
		Cards       []Card
		BonusType   Bonus
		BonusPlayed bool
	}
	type args struct {
		d *Deck
		a Action
	}
	type want struct {
		discarded []Card
		water     int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{name: "no actions no cards",
			fields: fields{
				Water:       0,
				Cards:       []Card{},
				BonusType:   BonusGhost,
				BonusPlayed: true,
			},
			args: args{NewDeck(), Skip},
			want: want{[]Card{}, 0},
		},
		{name: "no actions 3 water",
			fields: fields{
				Water:       0,
				Cards:       []Card{Water, Water, Water},
				BonusType:   BonusGhost,
				BonusPlayed: true,
			},
			args: args{NewDeck(), Skip},
			want: want{[]Card{}, 0},
		},
		{name: "PutWater",
			fields: fields{
				Water:       0,
				Cards:       []Card{Wave, Water},
				BonusType:   BonusGhost,
				BonusPlayed: true,
			},
			args: args{NewDeck(), PutWater},
			want: want{[]Card{Wave}, 1},
		},
		{name: "PutWaterDouble",
			fields: fields{
				Water:       0,
				Cards:       []Card{DoubleWave, Water},
				BonusType:   BonusGhost,
				BonusPlayed: true,
			},
			args: args{NewDeck(), PutWaterDouble},
			want: want{[]Card{DoubleWave}, 1},
		},
		{name: "PutTwoWater",
			fields: fields{
				Water:       0,
				Cards:       []Card{DoubleWave, Water, Water},
				BonusType:   BonusGhost,
				BonusPlayed: true,
			},
			args: args{NewDeck(), PutTwoWater},
			want: want{[]Card{DoubleWave}, 2},
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				Water:       tt.fields.Water,
				Cards:       tt.fields.Cards,
				BonusType:   tt.fields.BonusType,
				BonusPlayed: tt.fields.BonusPlayed,
			}
			p.Play(tt.args.d, tt.args.a)
			if got := (want{tt.args.d.discarded, p.Water}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Player.Play() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
