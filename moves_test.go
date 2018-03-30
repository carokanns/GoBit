package main

import (
	"testing"
)

func Test_moveList_add(t *testing.T) {
	tests := []struct {
		name string
		ml   moveList
		mv   move
	}{
		{"", moveList{}, 0},
		{"", moveList{}, 1},
		{"", moveList{}, 2},
		{"", moveList{}, 3},
	}

	ml = moveList{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ml.add(tt.mv)
			ix := len(tt.ml) - 1
			if ix < 0 {
				t.Fatalf("wrong len %v: %v %v", ix, tt.ml, tt.ml[0])
			}
			testMv := tt.ml[ix]
			if testMv != tt.mv {
				t.Errorf("tt.ml.add() = %v. want %v", testMv, tt.mv)
			}
		})
	}
}

func Test_move_packMove(t *testing.T) {
	type args struct {
		fr    uint
		to    uint
		p12   uint
		cp    uint
		pr    uint
		epSq  uint
		castl uint
	}
	tests := []struct {
		name string
		args args
	}{
		{"", args{A1, A2, wR, empty, empty, 0, (shortW | shortB)}},
		{"", args{D4, D5, bR, wQ, empty, E3, (shortW | longB)}},
	}
	var m move
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.packMove(tt.args.fr, tt.args.to, tt.args.p12, tt.args.cp, tt.args.pr, tt.args.epSq, tt.args.castl)
			if m.fr() != tt.args.fr {
				t.Errorf("%v: want fr=%v. Got %v ", tt.name, tt.args.fr, m.fr())
			}
			if m.to() != tt.args.to {
				t.Errorf("%v: want to=%v. Got %v ", tt.name, tt.args.to, m.to())
			}
			if m.p12() != tt.args.p12 {
				t.Errorf("%v: want p12=%v. Got %v ", tt.name, tt.args.p12, m.p12())
			}
			if m.cp() != tt.args.cp {
				t.Errorf("%v: want cp=%v. Got %v ", tt.name, tt.args.cp, m.cp())
			}
			if m.pr() != tt.args.pr {
				t.Errorf("%v: want r=%v. Got %v ", tt.name, tt.args.pr, m.pr())
			}
			if m.ep() != tt.args.epSq {
				t.Errorf("%v: want epr=%v. Got %v ", tt.name, tt.args.epSq, m.ep())
			}
			if m.castl() != castlings(tt.args.castl) {
				t.Errorf("%v: want castl=%v. Got %v ", tt.name, tt.args.castl, m.castl())
			}
		})
	}
}
