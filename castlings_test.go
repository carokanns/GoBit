package main

import (
	"testing"
)

func Test_castlings_off(t *testing.T) {

	tests := []struct {
		name      string
		initCastl uint
		setCastl  uint
		want      uint
	}{
		{"shortW 1", longW | shortW | longB | shortB, shortW, longW | shortB | longB},
		{"shortW 2", longW | longB | shortB, shortW, longW | shortB | longB},
		{"shortW 3", longB | shortB, shortW, shortB | longB},
		{"longW 1", longW | shortW | longB | shortB, longW, shortW | shortB | longB},
		{"longW 2", shortW | longB | shortB, longW, shortW | shortB | longB},
		{"longW 3", longB | shortB, longW, shortB | longB},
		{"shortB 1", longW | shortW | longB | shortB, shortB, longW | shortW | longB},
		{"shortB 2", longW | longB | shortW, shortB, longW | shortW | longB},
		{"shortB 3", shortW | longW, shortB, shortW | longW},
		{"longB 1", longW | shortW | longB | shortB, longB, longW | shortW | shortB},
		{"longB 2", longW | shortW | shortB, longB, longW | shortW | shortB},
		{"longB 3", longW | shortW, longB, shortW | longW},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board.castlings = castlings(tt.initCastl)
			board.off(tt.setCastl)
			if board.castlings != castlings(tt.want) {
				t.Errorf("%v: We want %v but we got %v", tt.name, tt.want, board.castlings)
			}

		})
	}
}

func Test_castlings_on(t *testing.T) {

	tests := []struct {
		name      string
		initCastl uint
		setCastl  uint
		want      uint
	}{
		{"shortW 1", longW | shortW | longB | shortB, shortW, longW | shortW | longB | shortB},
		{"shortW 2", longW | longB | shortB, shortW, longW | shortW | longB | shortB},
		{"shortW 3", 0, shortW, shortW},
		{"longW 1", longW | shortW | longB | shortB, longW, longW | shortW | longB | shortB},
		{"longW 2", shortW | longB | shortB, longW, longW | shortW | longB | shortB},
		{"longW 3", 0, longW, longW},
		{"shortB 1", longW | shortW | longB | shortB, shortB, longW | shortW | longB | shortB},
		{"shortB 2", longW | longB | shortW, shortB, longW | shortW | longB | shortB},
		{"shortB 3", 0, shortB, shortB},
		{"longB 1", longW | shortW | longB | shortB, longB, longW | shortW | longB | shortB},
		{"longB 2", longW | shortW | shortB, longB, longW | shortW | longB | shortB},
		{"longB 3", 0, longB, longB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board.castlings = castlings(tt.initCastl)
			board.on(tt.setCastl)
			if board.castlings != castlings(tt.want) {
				t.Errorf("%v: We want %v but we got %v", tt.name, tt.want, board.castlings)
			}

		})
	}
}

func Test_parseCastlings(t *testing.T) {

	tests := []struct {
		name string
		fen  string
		want castlings
	}{
		{"", "KQkq", castlings(shortW | longW | shortB | longB)},
		{"", "-", castlings(0)},
		{"", "K ", castlings(shortW)},
		{"", "Qkq", castlings(longW | shortB | longB)},
		{"", "kq", castlings(shortB | longB)},
		{"", "Kk", castlings(shortW | shortB)},
		{"", "KQkqrr", castlings(shortW | longW | shortB | longB)},
		{"", "kQq", castlings(longW | shortB | longB)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseCastlings(tt.fen); got != tt.want {
				t.Errorf("parseCastlings(%v) = %v, want %v", tt.fen, got, tt.want)
			}
		})
	}
}

func Test_castlings_String(t *testing.T) {
	tests := []struct {
		name string
		c    castlings
		want string
	}{
		{"", 0, "-"},
		{"", castlings(shortW | longB), "Kq"},
		{"", castlings(shortW | longB), "Kq"},
		{"", castlings(shortW | longW), "KQ"},
		{"", castlings(shortB | longB), "kq"},
		{"", castlings(shortW | longB | shortB | longW), "KQkq"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("castlings.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
