package main

import (
	"fmt"
	"math"
	"strconv"
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
	frEngine = make(chan string)
	toEngine = make(chan bool)
	go root(toEngine, frEngine)

	return
}

//TODO root: Iterative Depening

//TODO root: Aspiration search
func root(toEngine chan bool, frEngine chan string) {
	var depth int
	var pv pvList
	var childPV pvList
	childPV.new()
	pv.new()
	b := &board
	ml := make(moveList, 0, 60)
	for _ = range toEngine {
		limits.startTime, limits.nextTime = time.Now(), time.Now()
		alpha, beta := minEval, maxEval

		cntNodes = 0
		killers.clear()
		ml.clear()
		pv.clear()
		genAndSort(b, &ml)
		bm := ml[0]
		bs := noScore
		depth = limits.depth
		for depth = 1; depth <= limits.depth; depth++ {
			bs = noScore
			for ix := range ml {
				mv := &ml[ix]
				childPV.clear()

				b.move(*mv)
				tell("info depth ", strconv.Itoa(depth), " currmove ", mv.String(), " currmovenumber ", strconv.Itoa(ix+1))
				score := -search(-beta, -alpha, depth-1, 1, &childPV, b)
				b.unmove(*mv)
				if limits.stop {
					break
				}
				mv.packEval(score)
				if score > bs {
					bs = score
					pv.catenate(*mv, &childPV)

					bm = *mv
					alpha = score

					t1 := time.Since(limits.startTime)
					tell(fmt.Sprintf("info score cp %v depth %v nodes %v time %v pv ", bm.eval(), depth, cntNodes, int(t1.Seconds()*1000)), pv.String())
				}
			}
			ml.sort()
		}

		t1 := time.Since(limits.startTime)
		tell(fmt.Sprintf("info score cp %v depth %v nodes %v time %v pv ", bm.eval(), depth-1, cntNodes, int(t1.Seconds()*1000)), pv.String())
		frEngine <- fmt.Sprintf("bestmove %v%v", sq2Fen[bm.fr()], sq2Fen[bm.to()])
	}
}

//TODO search: generate all moves and put captures first  (temporary)
//TODO search: qs
//TODO search: hash table/transposition table
//TODO search: history table and maybe counter move table
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
		return signEval(b.stm, evaluate(b))
	}
	pv.clear()
	ml := make(moveList, 0, 60)
	//genAndSort(b, &ml)
	genInOrder(b, &ml, ply)

	bm, bs := noMove, noScore
	childPV := make(pvList, 0, maxPly)
	for _, mv := range ml {
		if !b.move(mv) {
			continue
		}

		childPV.clear()

		score := -search(-beta, -alpha, depth-1, ply+1, &childPV, b)

		b.unmove(mv)

		if score > bs {
			bs = score
			pv.catenate(mv, &childPV)

			if score >= beta { // beta cutoff
				// add killer and update history
				if mv.cp() == empty && mv.pr() == empty {
					killers.add(mv, ply)
				}
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
	ml.clear()
	b.genAllLegals(ml)
	for ix, mv := range *ml {
		b.move(mv)
		v := evaluate(b)
		b.unmove(mv)

		v = signEval(b.stm, v)
		(*ml)[ix].packEval(v)
	}

	ml.sort()
}

// generate capture moves first, then killers, then non captures
func genInOrder(b *boardStruct, ml *moveList, ply int) {
	ml.clear()
	b.genAllCaptures(ml)
	noCaptIx := len(*ml)
	b.genAllNonCaptures(ml)

	if len(*ml)-noCaptIx > 2 {
		// place killers first among non captuers
		for ix := noCaptIx; ix < len(*ml); ix++ {
			mv := (*ml)[ix]
			if killers[ply].k1.cmpFrTo(mv) {
				(*ml)[ix], (*ml)[noCaptIx] = (*ml)[noCaptIx], (*ml)[ix]
			} else if killers[ply].k2.cmpFrTo(mv) {
				(*ml)[ix], (*ml)[noCaptIx+1] = (*ml)[noCaptIx+1], (*ml)[ix]
			}
		}
	}

}

func signEval(stm color, ev int) int {
	if stm == BLACK {
		return -ev
	}
	return ev
}

/////////////////  Killers ///////////////////////////////////////////////
// killerStruct holds the killer moves per ply
type killerStruct [maxPly]struct {
	k1 move
	k2 move
}

// Clear killer moves
func (k *killerStruct) clear() {
	for ply := 0; ply < maxPly; ply++ {
		k[ply].k1 = noMove
		k[ply].k2 = noMove
	}
}

// add killer 1 and 2 (Not inCheck, caaptures and promotions)
func (k *killerStruct) add(mv move, ply int) {
	if !k[ply].k1.cmpFrTo(mv) {
		k[ply].k2 = k[ply].k1
		k[ply].k1 = mv
	}
}

var killers killerStruct

///////////////////////////// history table //////////////////////////////////
type historyStruct [2][64][64]uint

func (h *historyStruct) inc(fr, to int, stm color, depth int) {
	h[stm][fr][to] += uint(depth * depth)
}

func (h *historyStruct) clear() {
	for fr := 0; fr < 64; fr++ {
		for to := 0; to < 64; to++ {
			h[0][fr][to] = 0
			h[1][fr][to] = 0
		}
	}
}

var history historyStruct
