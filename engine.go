package main

import (
	"fmt"
	"math"
	"time"
)

const (
	maxDepth = 100
	maxPly   = 100
)

var cntNodes uint64

//TODO search limits: start clock and testing for movetime
//TODO search limits: counting nodes and testing for limit.nodes
//TODO search limits: limit.depth

//TODO search limits: time per game w/wo increments
//TODO search limits: time per x moves and after x moves w/wo increments
type searchLimits struct {
	depth     int
	nodes     uint64
	moveTime  int // in milliseconds
	infinite  bool
	startTime time.Time
	nextTime  time.Time

	//////////////// current //////////
	stop bool
}

var limits searchLimits

func (s *searchLimits) init() {
	s.depth = 9999
	s.nodes = math.MaxUint64
	s.moveTime = 99999999999
	s.infinite = false
	s.stop = false
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

type pvList []move

func (pv *pvList) new() {
	*pv = make(pvList, 0, maxPly)
}

func (pv *pvList) add(mv move) {
	*pv = append(*pv, mv)
}

func (pv *pvList) clear() {
	*pv = (*pv)[:0]
}

func (pv *pvList) addPV(pv2 *pvList) {
	*pv = append(*pv, *pv2...)
}

func (pv *pvList) catenate(mv move, pv2 *pvList) {
	pv.clear()
	pv.add(mv)
	pv.addPV(pv2)
}

func (pv *pvList) String() string {
	s := ""
	for _, mv := range *pv {
		s += mv.String() + " "
	}
	return s
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
	var pv pvList
	var childPV pvList
	childPV.new()
	b := &board
	ml := make(moveList, 0, 60)
	for _ = range toEngine {
		limits.startTime, limits.nextTime = time.Now(), time.Now()
		alpha, beta := -minEval, maxEval
		bm, bs := noMove, noScore
		depth := limits.depth
		cntNodes = 0
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
		frEngine <- fmt.Sprintf("bestmove %v%v", sq2Fen[ml[0].fr()], sq2Fen[ml[0].to()])
	}
}

//TODO search: the 'stop' command
//TODO search: time handling basic movetime
//TODO search: alpha Beta

//TODO search: generate all moves and put captures first  (temporary)
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

func search(alpha, beta, depth, ply int, pv *pvList, b *boardStruct) int {
	cntNodes++
	if depth <= 0 {
		return signEval(b.stm,evaluate(b))
	}
	pv.clear()
	ml := make(moveList, 0, 60)
	genAndSort(b, &ml)

	bm, bs := noMove, noScore
	var childPV pvList
	for _, mv := range ml {
		childPV.clear()
		b.move(mv)
		score := -search(-beta, -alpha, depth-1, ply+1, &childPV, b)
		b.unmove(mv)
		if score > bs {
			bs = score
			pv.catenate(mv, &childPV)

			if score >= beta { // beta cutoff
				return score
			}
			if score > alpha {
				bm = mv
				_ = bm
				alpha = score
			}
		}
		if time.Since(limits.nextTime) >= time.Duration(time.Second) {
			t1 := time.Since(limits.startTime)
			tell(fmt.Sprintf("info time %v nodes %v nps %v", int(t1.Seconds()*1000), cntNodes, cntNodes/uint64(t1.Seconds())))
			limits.nextTime = time.Now()
		}

		if limits.stop {
			return alpha
		}
	}
	return bs
}

func genAndSort(b *boardStruct, ml *moveList) {

	b.genAllMoves(ml)

	for ix, mv := range *ml {
		b.move(mv)
		v := evaluate(b)
		b.unmove(mv)
		v = signEval(b.stm, v)
		(*ml)[ix].packEval(v)
	}

	ml.sort()
}

func signEval(stm color, ev int) int {
	if stm == BLACK {
		return -ev
	}
	return ev
}
