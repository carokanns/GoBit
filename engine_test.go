package main

import (
	"fmt"
	"testing"
	"time"
)

///////////////////////////////////////// BENCHMARKS /////////////////////////////////////////
func Benchmark_root(b *testing.B) {
	/*	initFen2Sq()
		initMagic()
		initAtksKings()
		initAtksKnights()
	*/
	limits.init()
	limits.setDepth(8)
	handleNewgame()
	handlePosition("position startpos moves d2d4 d7d5 c1f4 g8f6 e2e3 e7e6 b1d2 c7c5 c2c3 b8c6 g1f3 f8e7 f1d3 c8d7")

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		fakeRoot()
	}
}
//  times
/* only a/b nothing else
Benchmark_root-4: 192.911s 	       1	192848735500 ns/op	13250269384 B/op	28233303 allocs/op
*/
/* a/b  sort captures first, followed by killers
Benchmark_root-4 10.313s   	       1	10248174300 ns/op	4391864352 B/op	 9352403 allocs/op
*/
/* a/b only killers first no other sorting
Benchmark_root-4 241.245s   	       1	241180949700 ns/op	79804119512 B/op	169935345 allocs/op
*/

func fakeRoot() {
	var pv pvList
	var childPV pvList
	childPV.new()
	b := &board
	ml := make(moveList, 0, 60)
	limits.startTime, limits.nextTime = time.Now(), time.Now()
	alpha, beta := minEval, maxEval
	bm, bs := noMove, noScore
	depth := limits.depth
	cntNodes = 0
	killers.clear()
	ml.clear()
	genAndSort(b, &ml)

	for ix := range ml {
		mv := &ml[ix]
		childPV.clear()

		b.move(*mv)
		tell("info currmove ", mv.String())
		score := -search(-beta, -alpha, depth-1, 1, &childPV, b)
		b.unmove(*mv)
		mv.packEval(signEval(b.stm, score))
		if score > bs {
			bs = score
			pv.clear()
			pv.catenate(*mv, &childPV)

			bm = *mv
			alpha = score
			tell(fmt.Sprintf("info score cp %v depth %v nodes %v pv ", bs, depth, cntNodes), pv.String())
		}
	}
	ml.sort()
	tell(fmt.Sprintf("info score cp %v depth %v nodes %v pv ", bm.eval(), depth, cntNodes), pv.String())
	fmt.Printf("bestmove %v%v", sq2Fen[ml[0].fr()], sq2Fen[ml[0].to()])
}


//////////////////////////////////// TESTS ////////////////////////////////////////////////////////
func Test_genAndSort(t *testing.T) {
	tests := []struct {
		name    string
		pos     string
		want1Mv string
		want1Ev int
	}{
		{"", "position startpos ", "e2e4", 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ml moveList
			genAndSort(&board, &ml)
			if tt.want1Mv != trim(ml[0].String()) {
				t.Errorf("%v: %#v should be best move. Got %#v", tt.name, tt.want1Mv, trim(ml[0].String()))
			}
			if tt.want1Ev != ml[0].eval() {
				t.Errorf("%v: %v should be best score. Got %v", tt.name, tt.want1Ev, ml[0].eval())
			}

		})
	}
}

func Test_see(t *testing.T) {
	tests := []struct {
		name   string
		fr, to int
		fen    string
		want   int
	}{

		// Pawns
		{"D4E5", D4, E5, "rnbqkbnr/pppp1ppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2", 100},
		{"F5E4", F5, E4, "rnbqkbnr/ppppp2p/8/5pB1/3PP3/8/PPP2PPP/RN1QKBNR b KQkq - 0 3", 100},
		{"D4E5 b", D4, E5, "rnbqkbnr/pp4pp/2p2p2/3pp3/3PPP2/2P5/PP4PP/RNBQKBNR w KQkq - 0 5", 100},
		{"F4E5", F4, E5, "rnbqkbnr/pp4pp/2p2p2/3pp3/3PPP2/2P5/PP4PP/RNBQKBNR w KQkq - 0 5", 100},
		{"E4D5", E4, D5, "rnbqkbnr/pp4pp/2p2p2/3pp3/3PPP2/2P5/PP4PP/RNBQKBNR w KQkq - 0 5", 0},
		{"D5E4", D5, E4, "rnbqkbnr/pp4pp/2p2p2/3pp3/3PPP2/2P5/PP4PP/RNBQKBNR b KQkq - 0 5", 100},
		{"E5D4", E5, D4, "rnbqkbnr/pp4pp/2p2p2/3pp3/3PPP2/2P5/PP4PP/RNBQKBNR b KQkq - 0 5", 0},
		{"E5F4", E5, F4, "rnbqkbnr/pp4pp/2p2p2/3pp3/3PPP2/2P5/PP4PP/RNBQKBNR b KQkq - 0 5", 0},
		//King
		{"Bxf7", C4, F7, "rnbqk2r/pppp1pp1/5n1p/2b1p1N1/2B1P3/8/PPPP1PPP/RNBQK2R w KQkq - 0 5", 100},

		// Mixed
		{"Bxh7", D3, H7, "r2q1rk1/pp1bbppp/2n1pn2/2pp2N1/3P1B2/2PBP3/PP1N1PPP/R2QK2R w KQ - 6 9", -250},
		{"Nxh7", G5, H7, "r2q1rk1/pp1bbppp/2n1pn2/2pp2N1/3P1B2/2PBP3/PP1N1PPP/R2QK2R w KQ - 6 9", -225},
		{"Nxf7", G5, F7, "r2q1rk1/pp1bbppp/2n1pn2/2pp2N1/3P1B2/2PBP3/PP1N1PPP/R2QK2R w KQ - 6 9", -225},
		{"Nxe6", G5, E6, "r2q1rk1/pp1bbppp/2n1pn2/2pp2N1/3P1B2/2PBP3/PP1N1PPP/R2QK2R w KQ - 6 9", -225},
		{"d5xf3 a", D5, F3, "2kr1bnr/pb3ppp/1p2p3/2pqn3/8/P1PP1NP1/1P3PBP/RNBQ1RK1 b - - 0 10", -275},
		{"d5xf3 b", D5, F3, "r3kbnr/pb3ppp/1p2p3/2pqn3/8/2PP1NP1/PP3PBP/RNBQ1RK1 b kq - 4 9", -275},
		{"d5xd3", D5, D3, "r3kbnr/pb3ppp/1p2p3/2pqn3/8/2PP1NP1/PP3PBP/RNBQ1RK1 b kq - 4 9", 100},

		{"Dxg7", G4, G7, "r3kbnr/pb3ppp/1p2p3/2pqn3/3N2Q1/2PP2P1/PP3PBP/RNB2RK1 w kq - 4 3", -850},
		{"Dxe6 a", G4, E6, "r3kbnr/pb3ppp/1p2p3/2pqn3/3N2Q1/2PP2P1/PP3PBP/RNB2RK1 w kq - 4 3", -850},
		{"Dxe6 b", G4, E6, "r3kbnr/pb4pp/1p2pp2/2pqn3/3N2Q1/2PP2P1/PP3PBP/RNB2RK1 w kq - 0 4", 100},
		{"Nxe6", D4, E6, "r3kbnr/pb4pp/1p2pp2/2pqn3/3N2Q1/2PP2P1/PP3PBP/RNB2RK1 w kq - 0 4", 100},

		{"e5xf3", E5, F3, "r3kbnr/pb3ppp/1p2p3/2pqn3/8/2PP1NP1/PP3PBP/RNBQ1RK1 b kq - 4 9", 325},
		{"e5xd3", E5, D3, "r3kbnr/pb3ppp/1p2p3/2pqn3/8/2PP1NP1/PP3PBP/RNBQ1RK1 b kq - 4 9", 100},
		{"fxe4", F5, E4, "rnbqkbnr/ppppp2p/8/5pp1/3PP3/8/PPP2PPP/RNBQKBNR b KQkq - 0 3", 100},
		{"w ungarded pawn", B5, C4, "rnbqkbnr/p1pppppp/8/1p6/2P5/8/PP1PPPPP/RNBQKBNR b KQkq - 0 2", 100},

		{"f3xe5", F3, E5, "r3kbnr/pb3ppp/1p2p3/2pqn3/8/2PP1NP1/PP3PBP/RNBQ1RK1 w kq - 4 9", 0},
		{"Bxg5 a", C1, G5, "rnbqkbnr/pppp3p/4p3/5pp1/3PP3/8/PPP2PPP/RNBQKBNR w KQkq - 0 3", -250},
		{"Bxg5 b", C1, G5, "rnbqkbnr/ppppp2p/8/5pp1/3PP3/8/PPP2PPP/RNBQKBNR w KQkq - 0 3", 100},
		{"exf5", E4, F5, "rnbqkbnr/ppppp2p/8/5pp1/3PP3/8/PPP2PPP/RNBQKBNR w KQkq - 0 3", 100},
		{"ungarded pawn", C4, B5, "rnbqkbnr/p1pppppp/8/1p6/2P5/8/PP1PPPPP/RNBQKBNR w KQkq - 0 2", 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlePosition("position fen " + tt.fen)
			if got := see(tt.fr, tt.to, &board); got != tt.want {
				t.Errorf("see() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initQs(t *testing.T) {
	type args struct {
		fr int
		to int
	}
	tests := []struct {
		name string
		pos  string
		want []args
	}{
		{"startpos", "position startpos", []args{}},
		{"dxe5", "position startpos moves d2d4 e7e5", []args{{D4, E5}}},
		{"QxQ delayed", "position fen r1b1kbnr/ppp2ppp/2n5/3pp1q1/3P2Q1/4PN2/PPP2PPP/RNB1KB1R w KQkq - 4 5", []args{{D4, E5}, {F3, E5}, {F3, G5}} },
		{"ungarded pawn", "position fen rnbqkbnr/p1pppppp/8/1p6/2P5/8/PP1PPPPP/RNBQKBNR w KQkq - 0 2", []args{{C4, B5}, } },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &board
			handlePosition(tt.pos)
			wantMl := make(moveList, 0, 30)
			ml := make(moveList, 0, 30)
			for _, strMv := range tt.want {
				fr := strMv.fr
				to := strMv.to
				mv := move(0)
				mv.packMove(fr, to, b.sq[fr], b.sq[to], empty, b.ep, b.castlings)
				wantMl.add(mv)
			}
			initQS(&ml, b)

			if len(wantMl) == 0 {
				if len(ml) != 0 {
					t.Errorf("%v: wantMl len=0 got len=%v", tt.name, len(ml))
				}
			} else {
				if len(ml) < 1 {
					t.Errorf("%v: want len(ml) > 0 got len=%v", tt.name, len(ml))
				} else {
					if wantMl[0] != ml[0] {
						t.Errorf("%v: wantMl[0] = %v; got %v  (%v)", tt.name, wantMl[0], ml[0], ml.String())
					}
					if len(wantMl) > 1 && wantMl[1] != ml[1] {
						t.Errorf("%v: wantMl[1] = %v; got %v", tt.name, wantMl[1], ml[1])
					}
					if len(wantMl) > 2 && wantMl[2] != ml[2] {
						t.Errorf("%v: wantMl[2] = %v; got %v", tt.name, wantMl[2], ml[2])
					}
				}
			}

		})
	}
}

func Test_QS(t *testing.T) {
	tests := []struct {
		name    string
		comment string
		pos     string

		wantEval  int
		wantDelta int
	}{
		{"startpos", "startpos must be 0", "position startpos", 0, 0},
		{"dxe5", "win a pawn on e5", "position startpos moves d2d4 e7e5", 80, 120},
		{"QxQ delayed", "win a Queen", "position fen r1b1kbnr/ppp2ppp/2n5/3pp1q1/3P2Q1/4PN2/PPP2PPP/RNB1KB1R w KQkq - 4 5", 900, 50},

		{"", "After fxe4 and Bxg5 it is material equal", "position fen rnbqkbnr/ppppp2p/8/5pp1/3PP3/8/PPP2PPP/RNBQKBNR b KQkq - 0 3", 50, 10},
		{"", "fxe4 and black is pawn up", "position fen rnbqkbnr/ppppp2p/8/5pp1/3PP3/8/PPPK1PPP/RNBQ1BNR b kq - 1 3", 80, 20},
		{"", "equal after fxe4", "position fen rnbqkbnr/ppppp2p/8/5pB1/3PP3/8/PPP2PPP/RN1QKBNR b KQkq - 0 3", -60, 15},
		{"", "Pawn can take unguarded pawn", "position fen rnbqkbnr/p1pppppp/8/1p6/2P5/8/PP1PPPPP/RNBQKBNR w KQkq - 0 2", pieceVal[wP], 20},
		// Pawn
		{"", "Pawn captures guarded pawn", "position fen rnbqkbnr/ppp1pppp/8/3p4/2P5/8/PP1PPPPP/RNBQKBNR w KQkq - 0 1", 0, 20},
		{"", "Pawn captures unguarded knight. Now under with queen and pawn", "position fen rnb1kbnr/ppp1pppp/8/3n4/2P5/8/PP1PPPPP/RNBQKBNR w KQkq - 0 1", pieceVal[wP] + pieceVal[wQ], 20},
		{"", "White is Q-P down", "position fen rnbqkbnr/p1pppppp/8/1p6/2P5/8/PP1PPPPP/RNB1KBNR w KQkq - 0 2", -pieceVal[wQ] + pieceVal[wP], 20},
		{"", "Black is Q-P down", "position fen rnb1kbnr/p1pppppp/8/1p6/2P5/8/PP1PPPPP/RNBQKBNR b KQkq - 0 2", -pieceVal[wQ] + pieceVal[wP], 20},
		
		// Promotions
		{"", "Pawn promotes and no captures", "position fen 8/1k2P3/8/q7/8/8/6K1/8 w - - 0 1", 30, 30},
		{"", "Pawn promotes with capture. B under. This special case is not working prperly", "position fen 3b1n2/1k2P3/8/q7/8/8/6K1/8 w - - 0 1", pieceVal[wP] - pieceVal[wB] - pieceVal[wN], 30},

		// Knight
		{"", "White is up knight-pawn", "position fen rnbqkbnr/ppp1pppp/8/3p4/1N6/8/PP1PPPPP/RNBQKBNR w KQkq - 0 1", pieceVal[wN] - pieceVal[wP], 20},
		{"", "White is up knight and queen", "position fen rnb1kbnr/ppp1pppp/8/3n4/5N2/8/PP1PPPPP/RNBQKBNR w KQkq - 0 1", pieceVal[wN] + pieceVal[wQ], 20},
		{"", "Bl Knigh captures guarded W queen", "position fen rnbqkbnr/ppp1pppp/8/3n4/5N2/8/PP1PPPPP/RNBQKBNR b KQkq - 0 1", pieceVal[wN], 20},
		// Bishop
		{"", "White is up bishop and queen", "position fen rnb1kbnr/ppp1pppp/8/3p4/2B5/8/PP1PPPPP/RNBQKBNR w KQkq - 0 1", pieceVal[wB] + pieceVal[wQ], 30},
		{"", "Bishop captures guarded knight", "position fen rnbqkbnr/ppp1pppp/8/3n4/4B3/8/PP1PPPPP/RNBQKBNR w KQkq - 0 1", 25, 10},
		{"", "Black is up a bishop", "position fen rnbqkbnr/ppp1pppp/8/3b4/4Q3/8/PP1PPPPP/RNBQKBNR b KQkq - 0 1", pieceVal[wB], 25},
		// Rook
		{"", "White is up Rook and queen", "position fen rnb1kbnr/ppp1pppp/8/3p4/3R4/8/PP1PPPPP/RNBQKBNR w KQkq - 0 1", pieceVal[wR] + pieceVal[wQ], 30},
		{"", "White is up Rook vs knight", "position fen rnbqkbnr/ppp1pppp/8/3n4/3R4/8/PP1PPPPP/RNBQKBNR w KQkq - 0 1", pieceVal[wR] - pieceVal[wN], 20},
		{"", "Black is up a Rook", "position fen rnbqkbnr/ppp1pppp/8/3r1B2/8/8/PP1PPPPP/RNBQKBNR b KQkq - 0 1", pieceVal[wR], 15},
		// Queen
		{"", "White is up a Queen vs a Pawn", "position fen rnbqkbnr/ppp1pppp/8/3p4/4Q3/8/PP1PPPPP/RNBQKBNR w KQkq - 0 1", pieceVal[wQ] - pieceVal[wP], 20},
		{"", "White is up 2 Queens", "position fen rnb1kbnr/ppp1pppp/8/3n4/2Q5/8/PP1PPPPP/RNBQKBNR w KQkq - 0 1", pieceVal[wQ] * 2, 50},
		{"", "Black Queen captures guarded W queen and is up a Queen", "position fen rnbqkbnr/ppp1pppp/8/3q4/4Q3/8/PP1PPPPP/RNBQKBNR b KQkq - 0 1", pieceVal[wQ], 60},
		// King
		{"", "White is Queen up after King captures unguarded pawn", "position fen rnb1kbnr/ppp1pppp/8/3p4/2K5/8/PP1PPPPP/RNBQ1BNR w KQkq - 0 1", pieceVal[wQ], 60},
		{"", "Black is knight up", "position fen rnbqkbnr/ppp1pppp/8/3n4/3K4/8/PP1PPPPP/RNBQ1BNR w KQkq - 0 1", -pieceVal[wN], 60},
		{"", "Bl King captures unguarded W knight. It's equal", "position fen rnbq1bnr/ppp1pppp/8/3k4/4N3/8/PP1PPPPP/RNBQKBNR b KQkq - 0 1", -45, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &board
			handlePosition(tt.pos)
			val := qs(maxEval, b)

			if val < tt.wantEval-tt.wantDelta || val > tt.wantEval+tt.wantDelta {
				t.Errorf("%v:  %v Should be in [%v,  %v] %v", tt.name, val, tt.wantEval-tt.wantDelta, tt.wantEval+tt.wantDelta, tt.comment)
			}

		})
	}
}

