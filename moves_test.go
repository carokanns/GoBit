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

	//ml := moveList{}
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

func Test_move_cmp(t *testing.T) {
	type args struct {
		fr    int
		to    int
		pc    int // 12 bits
		cp    int
		pr    int
		epSq  int
		castl castlings
		eval1 int
		eval2 int
		want  bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"", args{A1, A2, wR, empty, empty, 0, castlings(shortW | shortB), 123, 3333, true}},
		{"", args{D4, D5, bR, wQ, empty, E3, castlings(shortW | longB), 0, 2, true}},
	}
	var m1, m2, m3 move
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m1.packMove(tt.args.fr, tt.args.to, tt.args.pc, tt.args.cp, tt.args.pr, tt.args.epSq, tt.args.castl)
			m1.packEval(tt.args.eval1)
			m2.packMove(tt.args.fr, tt.args.to, tt.args.pc, tt.args.cp, tt.args.pr, tt.args.epSq, tt.args.castl)
			m2.packEval(tt.args.eval2)
			m3.packMove(tt.args.fr+1, tt.args.to+1, tt.args.pc, tt.args.cp, tt.args.pr, tt.args.epSq, tt.args.castl)
			m3.packEval(tt.args.eval2)

			if got := m1.cmpFrTo(m2); got != tt.args.want {
				t.Errorf("move.cmp() = %v, want %v", got, tt.args.want)
			}
			if got := m1.cmpFrTo(m3); got == tt.args.want {
				t.Errorf("move.cmp() = %v, want %v", got, !tt.args.want)
			}
		})
	}
}

func Test_move_packMove(t *testing.T) {
	type args struct {
		fr    int
		to    int
		pc    int // 12 bit
		cp    int
		pr    int
		epSq  int
		castl castlings
	}
	tests := []struct {
		name string
		args args
	}{
		{"", args{A1, A2, wR, empty, empty, 0, castlings(shortW | shortB)}},
		{"", args{D4, D5, bR, wQ, empty, E3, castlings(shortW | longB)}},
	}
	var m move
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.packMove(tt.args.fr, tt.args.to, tt.args.pc, tt.args.cp, tt.args.pr, tt.args.epSq, tt.args.castl)
			if m.fr() != tt.args.fr {
				t.Errorf("%v: want fr=%v. Got %v ", tt.name, tt.args.fr, m.fr())
			}
			if m.to() != tt.args.to {
				t.Errorf("%v: want to=%v. Got %v ", tt.name, tt.args.to, m.to())
			}
			if m.pc() != tt.args.pc {
				t.Errorf("%v: want pc=%v. Got %v ", tt.name, tt.args.pc, m.pc())
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

func Test_moveList_remove(t *testing.T) {
	tests := []struct {
		name string
		cnt  int
		ix   int
	}{
		{"5 3", 5, 3},
		{"5 0", 5, 0},
		{"5 4", 5, 4},
		{"1 0", 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ml moveList
			for i := 0; i < tt.cnt; i++ {
				ml.add(move(i))
			}
			ml.remove(tt.ix)
			if len(ml) != tt.cnt-1 {
				t.Errorf("%v: we should have %v moves but have %v", tt.name, tt.cnt-1, len(ml))
			}
		})
	}
}

func Test_move_packEval(t *testing.T) {
	type args struct {
		fr    int
		to    int
		pc    int
		cp    int
		pr    int
		epSq  int
		castl castlings

		score int
	}
	tests := []struct {
		name string
		args args
	}{
		{"", args{A1, A2, wR, empty, empty, 0, castlings(shortW | shortB), -1111}},
		{"", args{D4, D5, bR, wQ, empty, E3, castlings(shortW | longB), 222}},
		{"", args{E2, E4, wP, empty, empty, 0, castlings(shortW | shortB), -129}},
		{"", args{G1, F3, wN, empty, empty, E6, castlings(shortW | longB), 169}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mv := noMove
			mv.packMove(tt.args.fr, tt.args.to, tt.args.pc, tt.args.cp, tt.args.pr, tt.args.epSq, tt.args.castl)
			mv.packEval(-2999) // set garbage in eval
			mv.packEval(tt.args.score)
			if mv.fr() != tt.args.fr {
				t.Errorf("%v: want fr=%v. Got %v ", tt.name, tt.args.fr, mv.fr())
			}
			if mv.to() != tt.args.to {
				t.Errorf("%v: want to=%v. Got %v ", tt.name, tt.args.to, mv.to())
			}
			if mv.pc() != tt.args.pc {
				t.Errorf("%v: want pc=%v. Got %v ", tt.name, tt.args.pc, mv.pc())
			}
			if mv.cp() != tt.args.cp {
				t.Errorf("%v: want cp=%v. Got %v ", tt.name, tt.args.cp, mv.cp())
			}
			if mv.pr() != tt.args.pr {
				t.Errorf("%v: want pr=%v. Got %v ", tt.name, tt.args.pr, mv.pr())
			}
			if mv.ep() != tt.args.epSq {
				t.Errorf("%v: want ep=%v. Got %v ", tt.name, tt.args.epSq, mv.ep())
			}
			if mv.castl() != castlings(tt.args.castl) {
				t.Errorf("%v: want castl=%v. Got %v ", tt.name, tt.args.castl, mv.castl())
			}

			if mv.eval() != tt.args.score {
				t.Errorf("%v: want score=%v. Got %v ", tt.name, tt.args.score, mv.eval())
			}

		})
	}
}
