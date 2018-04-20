package main

import (
	"fmt"
	"math"
)


//TODO search limits: start clock and testing for movetime
//TODO search limits: counting nodes and testing for limit.nodes
//TODO search limits: depth and other limits
//TODO search limits: time per game w/wo increments
//TODO search limits: time per x moves and after x moves w/wo increments
type searchLimits struct {
	depth    int
	nodes    uint64
	moveTime int // in milliseconds
	infinite bool
	//////////////// current //////////
	stop bool
}

var limits searchLimits

func (s *searchLimits) init() {
	s.depth = 9999
	s.nodes = math.MaxUint64
	s.moveTime = 99999999999
	s.infinite = false
}

func (s *searchLimits) setStop(st bool) {
	s.stop = st
}
func (s *searchLimits) setDepth(d int) {
	s.depth = d
}
func (s *searchLimits) setMoveTime(m int) {
	s.moveTime = m
}
func (s *searchLimits) setInfinite(b bool) {
	s.infinite = b
}

func engine() (toEngine chan bool, frEngine chan string) {
	fmt.Println("info string Hello from engine")
	frEngine = make(chan string)
	toEngine = make(chan bool)
	go root(toEngine, frEngine)

	return
}

//TODO root: Iterative Depening
//TODO root: Aspiration search
func root(toEngine chan bool, frEngine chan string) {
	b := &board
	ml := moveList{}
	for _ = range toEngine {
		tell("info string engine got go! X")
		ml = moveList{}
		genAndSort(b, &ml)

		for _, mv := range ml {
			b.move(mv)
			score := -search(b)
			b.unmove(mv)

			mv.packEval(adjEval(b, score))
		}
		ml.sort()
		tell("info score cp ", fmt.Sprintf("%v", ml[0].eval()), " depth 1 pv ", ml[0].String())
		frEngine <- fmt.Sprintf("bestmove %v%v", sq2Fen[ml[0].fr()], sq2Fen[ml[0].to()])
	}
}

//TODO search: the 'stop' command
//TODO search: time handling basic movetime
//TODO search: generate all moves and put captures first  (temporary)
//TODO search: alpha Beta
//TODO search: qs
//TODO search: killer moves
//TODO search: hash table
//TODO search: history table or counter move table
//TODO search: move generation. More fast and accurate
//TODO search: Null Move
//TODO search: Late Move Reduction
//TODO search: Internal Iterative Depening
//TODO search: Delta Pruning
//TODO search: more complicated time handling schemes
//TODO search: other reductions and extensions

func search(b *boardStruct) int {

	return evaluate(b)
}

func genAndSort(b *boardStruct, ml *moveList) {
	b.genAllMoves(ml)

	for ix, mv := range *ml {
		b.move(mv)
		v := evaluate(b)
		b.unmove(mv)
		v = adjEval(b, v)
		(*ml)[ix].packEval(v)
	}

	ml.sort()
}

func adjEval(b *boardStruct, ev int) int {
	if b.stm == BLACK {
		return -ev
	}
	return ev
}
