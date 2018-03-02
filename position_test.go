
package main

import (
	"testing"
)

func Test_int2Fen(t *testing.T) {

	tests := []struct {
		name string
		p12  int
		want string
	}{
		{"", wP, "P"},
		{"", bP, "p"},
		{"", wK, "K"},
		{"", bK, "k"},
		{"", wN, "N"},
		{"", bN, "n"},
		{"", empty, " "},
		{"", wQ, "Q"},
		{"", bQ, "q"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := int2Fen(tt.p12); got != tt.want {
				t.Errorf("int2Fen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boardStruct_allBB(t *testing.T) {
	tests := []struct {
		name string
		wBB  bitBoard
		bBB  bitBoard
	}{
		{"1", 0xff, 0x00FF},
		{"2", 0x0, 0xAF00},
		{"3", 0xF, 0x00FF},
		{"4", 0xFFFF000000000000, 0xFFFF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board.wbBB[WHITE] = tt.wBB
			board.wbBB[BLACK] = tt.bBB
			correct := tt.wBB | tt.bBB
			if got := board.allBB(); got != correct {
				t.Errorf("%v: should be %v but we got %v", tt.name, correct, got)
			}
		})
	}
}
