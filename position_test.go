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

func Test_boardStruct_setSq(t *testing.T) {
	board.newGame()
	tests := []struct {
		name string
		p12  int
		sq   int
	}{
		{"Ke4", fen2Int("K"), E4},
		{"pe3", fen2Int("p"), E3},
		{"pe5", fen2Int("p"), E5},
		{"Qb6", fen2Int("Q"), B6},
		{"nh6", fen2Int("n"), H6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := board.count[tt.p12]
			board.setSq(tt.p12, tt.sq)
			if board.sq[tt.sq] != tt.p12 {
				t.Errorf("%v: board should contain %v on sq=%v. Got %v", tt.name, tt.p12, tt.sq, board.sq[tt.sq])
			}
			pc := piece(tt.p12)
			sd := p12Color(tt.p12)
			if !board.wbBB[sd].test(uint(tt.sq)) {
				t.Errorf("%v: wbBB[%v] on sq=%v should be set to 1 but is 0", tt.name, sd, tt.sq)
			}
			if !board.pieceBB[pc].test(uint(tt.sq)) {
				t.Errorf("%v: pieceBB[%v] on sq=%v should be set to 1 but is 0", tt.name, pc, tt.sq)
			}
			if board.count[tt.p12] != count+1 {
				t.Errorf("%v: count[%v] should be %v. Got %v", tt.name, tt.p12, count+1, board.count[tt.p12])
			}
		})
	}
}
