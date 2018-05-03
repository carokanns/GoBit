package main

import (
	"fmt"
	"testing"
	"time"
)

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

func Benchmark_root(b *testing.B) {
	/*	initFenSq2Int()
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



func fakeRoot(){
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
		killersx.clearx()
		ml.clear()
		genAndSort( b, &ml)

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


