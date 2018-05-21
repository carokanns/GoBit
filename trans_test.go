package main

import "testing"

func Test_trans_new(t *testing.T) {
	tests := []struct {
		name    string
		mB      int
		want    int
		wantErr string
	}{
		{"zero", 0, 1, ""},
		{"4096", 4096, 0, "max transtable size is 4GB (~4000 MB)"},
		{"257", 257, (256 << 20) / entrySize, ""},
		{"256", 256, (256 << 20) / entrySize, ""},
		{"255", 255, (128 << 20) / entrySize, ""},
		{"65", 65, (64 << 20) / entrySize, ""},
		{"64", 64, (64 << 20) / entrySize, ""},
		{"63", 63, (32 << 20) / entrySize, ""},
	}
	for _, tt := range tests {
		var trans transpStruct
		t.Run(tt.name, func(t *testing.T) {
			err := trans.new(tt.mB)
			if got := len(trans.tab); got != tt.want {
				t.Errorf("want %v got %v", tt.want, got)
			}
			str := ""
			if err != nil {
				str = err.Error()
			}

			if str != tt.wantErr {
				t.Errorf("want error=%#v. Got %#v", tt.wantErr, str)
			}
		})
	}
}

func Test_flipSide(t *testing.T) {
	tests := []struct {
		name string
		key  uint64
		want uint64
	}{
		{"", 123, ^uint64(123)},
		{"", 1, 0xfffffffffffffffe},
		{"", 0, 0xffffffffffffffff},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := flipSide(tt.key); got != tt.want {
				t.Errorf("flipSide() = %x, want %x back=%x", got, tt.want, ^got)
			}
		})
	}
}

func Test_init(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test"}, // dummy case.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initKeys()
			for i := 0; i < 12*64; i++ {
				for j := i + 1; j < 12*64; j++ {
					if randPcSq[i] == randPcSq[j] {
						t.Errorf("index %v and %v has the same value %v", i, j, randPcSq[i])
					}
				}
				for j := 0; j < 8; j++ {
					if randPcSq[i] == randEp[j] {
						t.Errorf("randPcSq[%v] is equal to randEp[%v]: %v ", i, j, randPcSq[i])
					}
				}
				for j := 0; j < 4; j++ {
					if randPcSq[i] == randCastl[j] {
						t.Errorf("randPcSq[%v] is equal to randCastl[%v]: %v ", i, j, randPcSq[i])
					}
				}
			}
		})
	}
}

func Test_Trans(t *testing.T) {
	tests := []struct {
		fr        int
		to        int
		pc        int
		cp        int
		ep        int
		castl     uint
		key       uint64
		depth     int
		ply       int
		score     int
		scoreType int
		comment   string
	}{ //  mv     ep           castlings           key              de  pl sc st
		{A1, A2, wR, empty, E6, shortW | longW | shortB | longB, 0x0fffffffffffffff, 14, 5, 1, 0xf, "first entry"},
		{A1, A2, wR, bQ, E6, shortW | longW | shortB | longB, 0x1fffffffffffffff, 13, 5, 2, 0xf, "second entry"},
		{A1, A2, bR, wQ, E6, shortW | longW | shortB | longB, 0x2fffffffffffffff, 10, 5, 3, 0xf, "third entry should be replaced by fifth"},
		{A1, A2, wR, empty, E6, shortW | longW | shortB | longB, 0x3fffffffffffffff, 11, 5, 4, 0xf, "fourth entry"},
		{A1, A2, wR, wQ, E6, shortW | longW | shortB | longB, 0x4fffffffffffffff, 15, 5, 5, 0xf, "fifth entry replaces the third"},
	}

	// Store in transTab
	handlePosition("position fen 8/6kp/5p2/3n2pq/3N1n1R/1P3P2/P6P/4QK2 w - - 2 2")
	trans.new(128)
	mv := noMove
	// fill all 4 entries with the same key but different lock. The fifth replavces the third
	for _, e := range tests {
		mv.packMove(e.fr, e.to, e.pc, bQ, empty, e.ep, castlings(e.castl))
		key := e.key
		depth := e.depth
		ply := e.ply
		sc := e.score
		scoreType := e.scoreType

		idx := key & uint64(trans.mask)
		t.Log("idx", idx)
		trans.store(key, mv, depth, ply, sc, scoreType)
	}

	// retrieve from TransTab
	for ix, e := range tests {
		key := e.key
		depth := e.depth
		ply := e.ply
		sc := e.score
		scType := e.scoreType
		//values from retrieve
		transMove := noMove
		transScore := noScore
		transScType := 0

		ok := false
		if transMove, transScore, transScType, ok = trans.retrieve(key, depth, ply); ok {
			// retrieve found ok
			switch ix {
			case 2: // special case
				t.Errorf("case %v: Shouldn't find its values. (mv %v, transMv %v), (sc %v, transSc %v), (st %v, transSt %v\n", ix+1, uint(mv), uint(transMove), sc, transScore, scType, transScType)
				t.Logf("case %v: %v", ix+1, e.comment)
			case 4: // special case
				if transMove != mv || transScore != sc || transScType != scType {
					t.Errorf("case %v: Did not overwrite case 3. (mv %v, transMv %v), (sc %v, transSc %v), (st %v, transSt %v\n", ix+1, uint(mv), uint(transMove), sc, transScore, scType, transScType)
					t.Logf("case %v not ok. %v", ix+1, e.comment)
				}
			default: // other cases
				if transMove == mv && transScore == sc && transScType == scType {
					// found the right valuse. fine.
				} else {
					t.Errorf("case %v: values not ok. (mv %v, transMv %v), (sc %v, transSc %v), (st %v, transSt %v\n", ix+1, uint(mv), uint(transMove), sc, transScore, scType, transScType)
				}
			}
		} else if ix == 2 {
			t.Logf("case %v. Coundn't find the entry and that's ok: %v", ix+1, e.comment)
		} else {
			t.Errorf("case %v: couldn't find the entry", ix+1)
		}
		// check more things
		if transMove != noMove {
			if transMove.fr() != e.fr {
				t.Errorf("case %v: want fr=%v. got %v", e.fr, transMove.fr(), ix+1)
			}
			if transMove.to() != e.to {
				t.Errorf("case %v: want to=%v. got %v", e.to, transMove.to(), ix+1)
			}
			if transMove.pc() != e.pc {
				t.Errorf("case %v: want pc=%v. got %v", e.pc, transMove.pc(), ix+1)
			}
			if transMove.ep(pcColor(e.pc)) != e.ep {
				t.Errorf("case %v: want ep=%v. got %v", e.ep, transMove.ep(pcColor(e.pc)), ix+1)
			}
			if transMove.castl() != castlings(e.castl) {
				t.Errorf("case %v: want castl=%v. got %v", castlings(e.castl), transMove.castl(), ix+1)
			}
			if transMove.cp() != bQ {
				t.Errorf("case %v: want cp=%v. got %v", e.cp, transMove.cp(), ix+1)
			}

		}
	}
}

func Test_transpStruct_retrieve(t *testing.T) {
	type args struct {
		fullKey uint64
		depth   int
		ply     int
	}
	tests := []struct {
		name       string
		t          transpStruct
		args       args
		wantMv     move
		wantSc     int
		wantScType int
		wantOk     bool
	}{
		{"not in tab", trans, args{0x0101010101010101, 4, 4}, noMove, noScore, 0, false},
		{"in tab", trans, args{0x0404040404040404, 4, 1}, 40, 4, 0, true},
		{"lower depth", trans, args{0x0404040404040404, 1, 1}, 40, 4, 0, true},
		{"higher depth", trans, args{0x0404040404040404, 6, 1}, 40, 4, 0, false},
	}
	// make some stores
	trans.store(0x0fffffffffffffff, 10, 1, 1, 1, 0)
	trans.store(0x0202020202020202, 20, 2, 2, 2, 0)
	trans.store(0x0303030303030303, 30, 3, 1, 3, 0)
	trans.store(0x0404040404040404, 40, 4, 3, 4, 0)
	trans.store(0x0606060606060606, 60, 6, 4, 5, 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transMv, transSc, transScType, found := tt.t.retrieve(tt.args.fullKey, tt.args.depth, tt.args.ply)
			if transMv != tt.wantMv {
				t.Errorf("retrieve() mv = %v, want %v", transMv, tt.wantMv)
			}
			if transSc != tt.wantSc {
				t.Errorf("retrieve() Score = %v, want %v", transSc, tt.wantSc)
			}
			if transScType != tt.wantScType {
				t.Errorf("retrieve() ScoreType = %v, want %v", transScType, tt.wantScType)
			}
			if found != tt.wantOk {
				t.Errorf("retrieve() Found = %v, want %v", found, tt.wantOk)
			}
		})
	}
}

func Test_transpStruct_lock(t *testing.T) {
	type args struct {
		fullKey uint64
	}
	tests := []struct {
		name string
		t    *transpStruct
		args args
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.lock(tt.args.fullKey); got != tt.want {
				t.Errorf("transpStruct.lock() = %v, want %v", got, tt.want)
			}
		})
	}
}
