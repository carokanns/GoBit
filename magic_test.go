package main

import (
	"testing"
)

func Test_sMagic_atks(t *testing.T) {
	initFen2Sq()
	initMagic()
	handleNewgame()
	tests := []struct {
		name     string
		position string
		m        *sMagic
		want     bitBoard
	}{
		{"", "startpos", &mRookTab[A1], createBitBoard(A2, B1)},
		{"", "startpos moves a2a4 b7b6", &mRookTab[A1], createBitBoard(A2, A3, A4, B1)},
		{"", "startpos", &mRookTab[H1], createBitBoard(H2, G1)},
		{"", "startpos moves h2h4 b7b6", &mRookTab[H1], createBitBoard(H2, H3, H4, G1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlePosition("position " + tt.position)
			if got := tt.m.atks(board.allBB()); got != tt.want {
				t.Errorf("sMagic.atks() = \n%v\nwant \n%v", got.Stringln(), tt.want.Stringln())
			}
		})
	}
}

func createBitBoard(bits ...int) bitBoard {
	BB := bitBoard(0)
	for _, b := range bits {
		BB.set(b)
	}
	return BB
}
