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

//TODO search limits: count nodes and test for limit.nodes
//TODO search limits: limit.depth

//TODO search limits: time per game w/wo increments
//TODO search limits: time per x moves and after x moves w/wo increments
type searchLimits struct {
	depth     int
	nodes     uint64
	moveTime  int // in milliseconds
	infinite  bool
	startTime time.Time
	lastTime  time.Time

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

type ebfStruct []uint64

func (e *ebfStruct) new() {
	*e = make(ebfStruct, 0, maxDepth)
}
func (e *ebfStruct) add(nodes uint64) {
	*e = append(*e, nodes)
}
func (e *ebfStruct) clear() {
	*e = (*e)[:0]
}
func (e *ebfStruct) ebf() float64 {
	if len(*e) < 4 {
		return 0
	}
	ebf := 0.0
	prevNodes1 := float64((*e)[len(*e)-2])
	prevNodes2 := float64((*e)[len(*e)-3])
	prevNodes3 := float64((*e)[len(*e)-4])

	if prevNodes2 > 0.0 && prevNodes3 > 0.0 {
		ebf = (prevNodes2/prevNodes3 + prevNodes1/prevNodes2) / 2
	}
	fmt.Printf("ebf: %0.2f age=%v Used=%v Stored: %v Tried: %v Found: %v Prunes: %v Best: %v\n", ebf, trans.age, trans.cntUsed, trans.cStores, trans.cTried, trans.cFound, trans.cPrune, trans.cBest)
	return ebf
}

func engine() (toEngine chan bool, frEngine chan string) {
	frEngine = make(chan string)
	toEngine = make(chan bool)
	go root(toEngine, frEngine)

	return
}

//TODO root: Aspiration search
func root(toEngine chan bool, frEngine chan string) {
	var depth, alpha, beta int
	var ebfTab ebfStruct
	var pv pvList
	var childPV pvList
	var ml moveList
	childPV.new()
	pv.new()
	ml.new(60)
	ebfTab.new()
	b := &board
	for _ = range toEngine {
		limits.startTime, limits.lastTime = time.Now(), time.Now()
		cntNodes = 0
		ebfTab.clear()
		killers.clear()
		ml.clear()
		pv.clear()

		trans.initSearch() // incr age coounters=0

		genAndSort(0, b, &ml)
		bm := ml[0]
		bs := noScore
		depth = 0

		transDepth := 0

		for depth = 1; depth <= limits.depth && !limits.stop; depth++ {
			ml.sort()
			bs = noScore // bm keeps the best from prev iteration in case of immediate stop before first is done in this iterastion
			alpha, beta = minEval, maxEval
			for ix, mv := range ml {
				childPV.clear()
				b.move(mv)
				tell("info depth ", strconv.Itoa(depth), " currmove ", mv.String(), " currmovenumber ", strconv.Itoa(ix+1))
				score := -search(-beta, -alpha, depth-1, 1, &childPV, b)

				b.unmove(mv)

				////////////////////////////////////
				if !checkKey(b) {
					fmt.Println("fullkey=", b.fullKey(), "Key", b.key, mv.StringFull())
					fmt.Println("INVALID KEY AFTER UNMOVE ROOT")
				}
				///////////////////////////////////////

				if limits.stop {
					break
				}
				ml[ix].packEval(score)
				if score > bs {
					bs = score
					pv.catenate(mv, &childPV)

					bm = ml[ix]
					alpha = score
					transDepth = depth
					if depth >= 0 {
						trans.store(b.fullKey(), mv, transDepth, 0, score, scoreTypeLower)
					}

					t1 := time.Since(limits.startTime)
					tell(fmt.Sprintf("info score cp %v depth %v nodes %v time %v pv ", bm.eval(), depth, cntNodes, int(t1.Seconds()*1000)), pv.String())
				}
			}

			ebfTab.add(cntNodes)
		}
		ml.sort()

		trans.store(b.fullKey(), bm, transDepth, 0, bs, scoreType(bs, alpha, beta))

		// time, nps, ebf
		t1 := time.Since(limits.startTime)
		nps := float64(0)
		if t1.Seconds() != 0 {
			nps = float64(cntNodes) / t1.Seconds()
		}
		ebfTab.ebf()
		tell(fmt.Sprintf("info score cp %v depth %v nodes %v  time %v nps %v pv %v", bm.eval(), depth-1, cntNodes, int(t1.Seconds()*1000), uint(nps), pv.String()))
		frEngine <- fmt.Sprintf("bestmove %v%v", sq2Fen[bm.fr()], sq2Fen[bm.to()])
	}
}

//TODO search: Null Move

//TODO search: Late Move Reduction
//TODO search: Internal Iterative Depening
//TODO search: Delta Pruning
//TODO search: more complicated time handling schemes
//TODO search: other reductions and extensions
func search(alpha, beta, depth, ply int, pv *pvList, b *boardStruct) int {
	cntNodes++
	if depth <= 0 {
		//return signEval(b.stm, evaluate(b))
		return qs(beta, b)
	}
	pv.clear()

	pvNode := depth > 0 && beta != alpha+1

	transMove := noMove
	useTT := depth >= 0
	transDepth := depth
	inCheck := b.isAttacked(b.King[b.stm], b.stm.opp())

	if depth < 0 && inCheck {
		useTT = true
		transDepth = 0
	}

	if useTT {
		var transSc, scType int
		ok := false

		if transMove, transSc, scType, ok = trans.retrieve(b.fullKey(), transDepth, ply); ok && !pvNode {
			switch {
			case scType == scoreTypeLower && transSc >= beta:
				trans.cPrune++
				return transSc
			case scType == scoreTypeUpper && transSc <= alpha:
				trans.cPrune++
				return transSc
			case scType == scoreTypeBetween:
				trans.cPrune++
				return transSc
			}
		}
	}

	var childPV pvList
	childPV.new() // TODO? make it smaller for each depth maxDepth-ply
	/////////////////////////////////////// NULL MOVE /////////////////////////////////////////
	ev := signEval(b.stm, evaluate(b))
	// null-move pruning
	if !pvNode && depth > 0 && !isMateScore(beta) && !inCheck && !b.isAntiNullMove() && ev >= beta {
		nullMv := b.moveNull()
		sc := minEval
		if depth <= 3 { // static
			// if you don't beat me with 100 points,
			// then I think your position sucks
			sc = -qs(-beta+1, b)
		} else { // dynamic
			sc = -search(-beta, -beta+1, depth-3-1, ply, &childPV, b)
		}

		b.undoNull(nullMv)

		if sc >= beta {
			if useTT {
				trans.store(b.fullKey(), noMove, transDepth, ply, sc, scoreTypeLower)
			}
			return sc
		}
	}
	/////////////////////// NULL MOVE END //////////////////////

	bs, score := noScore, noScore
	bm := noMove

	var genInfo = genInfoStruct{sv: 0, ply: ply, transMove: transMove}
	next = nextNormal
	for mv, msg := next(&genInfo, b); mv != noMove; mv, msg = next(&genInfo, b) {
		_ = msg

		if !b.move(mv) {
			continue
		}

		childPV.clear()

		if pvNode && bm != noMove {
			score = -search(-alpha-1, -alpha, depth-1, ply+1, &childPV, b)
			if score > alpha {
				score = -search(-beta, -alpha, depth-1, ply+1, &childPV, b)
			}
		} else {
			score = -search(-beta, -alpha, depth-1, ply+1, &childPV, b)
		}

		b.unmove(mv)
		
		if score > bs {
			bs = score
			bm = mv
			pv.catenate(mv, &childPV)
			if score > alpha {
				alpha = score
				if useTT {
					trans.store(b.fullKey(), mv, transDepth, ply, score, scoreType(score, alpha, beta))
				}
			}

			if score >= beta { // beta cutoff
				// add killer and update history
				if mv.cp() == empty && mv.pr() == empty {
					killers.add(mv, ply)
					history.inc(mv.fr(), mv.to(), b.stm, depth)
				}
				if mv.cmp(transMove) {
					trans.cPrune++
				}
				return score
			}
		}

		tStep := time.Since(limits.lastTime) - time.Duration(time.Millisecond*200)
		if tStep >= 0 {
			limits.lastTime = time.Now().Add(-time.Duration(tStep))
			t1 := time.Since(limits.startTime)
			ms := uint64(t1.Nanoseconds() / 1000000)
			if t1.Seconds() > 1 {
				if ms%1000 <= 5 {
					tell(fmt.Sprintf("info time %v nodes %v nps %v", ms, cntNodes, cntNodes/uint64(t1.Seconds())))
				}
			}

			if ms >= uint64(limits.moveTime)-200 {
				//				fmt.Println("t1", uint64(t1.Nanoseconds()/1000000)-100, "limit", uint64(limits.moveTime))
				limits.stop = true
			}
		}

		if limits.stop {
			return alpha
		}
	}
	if bm.cmp(transMove) {
		trans.cBest++
	}
	return bs
}

// is this a position to avoid null move?
func (b *boardStruct) isAntiNullMove() bool {
	if b.wbBB[b.stm] == b.pieceBB[King]&b.wbBB[b.stm] {
		return true
	}
	return false
}

func initQS(ml *moveList, b *boardStruct) {
	ml.clear()
	b.genAllCaptures(ml)
}
func qs(beta int, b *boardStruct) int {
	ev := signEval(b.stm, evaluate(b))
	if ev >= beta {
		// we are good. No need to try captures
		return ev
	}
	bs := ev

	qsList := make(moveList, 0, 60)
	initQS(&qsList, b) // create attacks
	done := bitBoard(0)

	// move loop
	for _, mv := range qsList {
		fr := mv.fr()
		to := mv.to()

		// This works because we pick lower value pieces first
		if done.test(to) { // Don't do the same to-sw again
			continue
		}
		done.set(to)

		see := see(fr, to, b)

		if see == 0 && mv.cp() == empty {
			// must be a promotion that didn't captureand was not captured
			see = pieceVal[wQ] - pieceVal[wP]
		}

		if see <= 0 {
			continue // equal captures not interesting
		}

		sc := ev + see
		if sc > bs {
			bs = sc
			if sc >= beta {
				return sc
			}
		}
	}

	return bs
}

// see (Static Echange Evaluation)
// Start with the capture fr-to and find out all the other captures to to-sq
func see(fr, to int, b *boardStruct) int {
	var pVal = [16]int{100, -100, 325, -325, 325, -325, 500, -500, 950, -950, 10000, -10000, 0, 0, 0, 0}
	pc := b.sq[fr]
	cp := b.sq[to]
	cnt := 1
	us := pcColor(pc)
	them := us.opp()

	// All the attackers to the to-sq, but first remove the moving piece and use X-ray to the to-sq
	occ := b.allBB()
	occ.clr(fr)
	attackingBB :=
		mRookTab[to].atks(occ)&(b.pieceBB[Rook]|b.pieceBB[Queen]) |
			mBishopTab[to].atks(occ)&(b.pieceBB[Bishop]|b.pieceBB[Queen]) |
			(atksKnights[to] & b.pieceBB[Knight]) |
			(atksKings[to] & b.pieceBB[King]) |
			(b.wPawnAtksFr(to) & b.pieceBB[Pawn] & b.wbBB[BLACK]) |
			(b.bPawnAtksFr(to) & b.pieceBB[Pawn] & b.wbBB[WHITE])
	attackingBB &= occ

	if (attackingBB & b.wbBB[them]) == 0 { // 'they' have no attackers - good bye
		return abs(pVal[cp]) // always return score from 'our' point of view
	}

	// Now we continue to keep track of the material gain/loss for each capture
	// Always remove the last attacker and use x-ray to find possible new attackers

	lastAtkVal := abs(pVal[pc]) // save attacker piece value for later use
	var captureList [32]int
	captureList[0] = abs(pVal[cp])
	n := 1

	stm := them // change side to move

	for {
		cnt++

		var pt int
		switch { // select the least valuable attacker
		case (attackingBB & b.pieceBB[Pawn] & b.wbBB[stm]) != 0:
			pt = Pawn
		case (attackingBB & b.pieceBB[Knight] & b.wbBB[stm]) != 0:
			pt = Knight
		case (attackingBB & b.pieceBB[Bishop] & b.wbBB[stm]) != 0:
			pt = Bishop
		case (attackingBB & b.pieceBB[Rook] & b.wbBB[stm]) != 0:
			pt = Rook
		case (attackingBB & b.pieceBB[Queen] & b.wbBB[stm]) != 0:
			pt = Queen
		case (attackingBB & b.pieceBB[King] & b.wbBB[stm]) != 0:
			pt = King
		default:
			panic("Don't come here in see! ")
		}

		// now remove the pt above from the attackingBB and scan for new attackers by possible x-ray
		BB := attackingBB & (attackingBB & b.pieceBB[pt] & b.wbBB[stm])
		occ ^= (BB & -BB) // turn off the rightmost bit from BB in occ

		//  pick sliding attacks again (do it from to-sq)
		attackingBB |= mRookTab[to].atks(occ)&(b.pieceBB[Rook]|b.pieceBB[Queen]) |
			mBishopTab[to].atks(occ)&(b.pieceBB[Bishop]|b.pieceBB[Queen])
		attackingBB &= occ // but only attacking pieces

		captureList[n] = -captureList[n-1] + lastAtkVal
		n++

		// save the value of tha capturing piece to be used later
		lastAtkVal = pVal[pt2pc(pt, WHITE)] // using WHITE always gives positive integer
		stm = stm.opp()                     // next side to move

		if pt == King && (attackingBB&b.wbBB[stm]) != 0 { //NOTE: just changed stm-color above
			// if king capture and 'they' are atting we have to stop
			captureList[n] = pVal[wK]
			n++
			break
		}

		if attackingBB&b.wbBB[stm] == 0 { // if no more attackers
			break
		}

	}

	// find the optimal capture sequence and 'our' material value will be on top
	for n--; n != 0; n-- {
		captureList[n-1] = min(-captureList[n], captureList[n-1])
	}

	return captureList[0]
}

/* func genAndSort(b *boardStruct, ml *moveList) {
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
} */
func genAndSort(ply int, b *boardStruct, ml *moveList) {
	if ply > maxPly {
		panic("wtf maxply")
	}

	ml.clear()
	b.genAllLegals(ml)

	for ix, mv := range *ml {
		b.move(mv)
		v := evaluate(b)
		b.unmove(mv)
		if killers[ply].k1.cmp(mv) {
			v += 1000
		} else if killers[ply].k2.cmp(mv) {
			v += 900
		}

		v = signEval(b.stm, v)

		(*ml)[ix].packEval(v)
	}

	ml.sort()
}

// generate capture moves first, then killers, then non captures
func genInOrder(b *boardStruct, ml *moveList, ply int, transMove move) {
	ml.clear()
	b.genAllCaptures(ml)
	noCaptIx := len(*ml)
	b.genAllNonCaptures(ml)
	if transMove != noMove {
		for ix := 0; ix < len(*ml); ix++ {
			mv := (*ml)[ix]
			if transMove.cmp(mv) {
				(*ml)[ix], (*ml)[0] = (*ml)[0], (*ml)[ix]
				break
			}
		}
	}
	pos1, pos2 := noCaptIx, noCaptIx+1
	if (*ml)[pos1].cmp(transMove) {
		noCaptIx++
		pos1++
		pos2++
	}

	if len(*ml)-noCaptIx > 2 {
		// place killers first among non captures
		cnt := 0
		for ix := noCaptIx; ix < len(*ml); ix++ {
			mv := (*ml)[ix]
			switch {
			case killers[ply].k1.cmpFrTo(mv) && !mv.cmpFrTo(transMove) && b.sq[mv.to()] == empty:
				mv.packMove(mv.fr(), mv.to(), b.sq[mv.fr()], b.sq[mv.to()], mv.pr(), b.ep, b.castlings)
				(*ml)[ix] = mv
				(*ml)[ix], (*ml)[pos1] = (*ml)[pos1], (*ml)[ix]
				cnt++
			case killers[ply].k2.cmpFrTo(mv) && !mv.cmpFrTo(transMove) && b.sq[mv.to()] == empty:
				mv.packMove(mv.fr(), mv.to(), b.sq[mv.fr()], b.sq[mv.to()], mv.pr(), b.ep, b.castlings)
				(*ml)[ix] = mv
				(*ml)[ix], (*ml)[pos2] = (*ml)[pos2], (*ml)[ix]
				cnt++
			}
			if cnt >= 2 {
				break
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
func (h *historyStruct) get(fr, to int, stm color) uint {
	return h[stm][fr][to]
}

func (h *historyStruct) clear() {
	for fr := 0; fr < 64; fr++ {
		for to := 0; to < 64; to++ {
			h[0][fr][to] = 0
			h[1][fr][to] = 0
		}
	}
}

func (h historyStruct) Print(n int) {
	fmt.Println("history top", n)
	type top50 struct{ fr, to, sd, sc uint }
	var hTab = make([]top50, n, n)
	for ix := range hTab {
		hTab[ix].fr, hTab[ix].to, hTab[ix].sd, hTab[ix].sc = 0, 0, 0, 0
	}

	W, B := uint(WHITE), uint(BLACK)
	for fr := uint(0); fr < 64; fr++ {
		for to := uint(0); to < 64; to++ {
			sc := h.get(int(fr), int(to), WHITE)
			for ix := range hTab {
				if sc > hTab[ix].sc {
					for ix2 := n - 2; ix2 >= ix; ix2-- {
						hTab[ix2+1] = hTab[ix2]
					}
					hTab[ix].fr, hTab[ix].to, hTab[ix].sd, hTab[ix].sc = fr, to, W, sc
					break
				}
			}

			sc = h.get(int(fr), int(to), BLACK)
			for ix := range hTab {
				if sc > hTab[ix].sc {
					for ix2 := n - 2; ix2 >= ix; ix2-- {
						hTab[ix2+1] = hTab[ix2]
					}
					hTab[ix].fr, hTab[ix].to, hTab[ix].sd, hTab[ix].sc = fr, to, B, sc
					break
				}
			}
		}
	}
	for ix, ht := range hTab {
		if ht.fr == 0 && ht.to == 0 {
			continue
		}
		fmt.Printf("%2v: %v %v-%v   %v  \n", ix+1, color(ht.sd).String(), sq2Fen[int(ht.fr)], sq2Fen[int(ht.to)], ht.sc)
	}
}

var history historyStruct

/////////////////////////// Next move /////////////////////////////////////
var next func(*genInfoStruct, *boardStruct) (move, string) // or nextKEvasion or nextQS

const (
	initNext = iota
	nextTr
	nextFirstGoodCp
	nextGoodCp
	nextK1
	nextK2
	nextCounterMv
	nextFirstNonCp
	nextNonCp
	nextBadCp
	nextEnd
)

type genInfoStruct struct {
	// to be filled in, before first call to the next-function
	sv, ply   int
	transMove move

	// handle by the next-function
	captures, nonCapt moveList
	counterMv         move
}

func nextNormal(genInfo *genInfoStruct, b *boardStruct) (move, string) {
	switch genInfo.sv {
	case initNext:
		genInfo.sv = nextTr
		fallthrough
	case nextTr:
		genInfo.sv = nextFirstGoodCp
		if genInfo.transMove != noMove {
			if b.isLegal(genInfo.transMove) {
				return genInfo.transMove, "transMove"
			}
			genInfo.transMove = noMove
		}
		fallthrough
	case nextFirstGoodCp:
		genInfo.captures.new(20)
		b.genAllCaptures(&genInfo.captures)
		// pick a good capt - use see - not transMove
		bs := -1
		bIx := 0
		ml := &genInfo.captures
		for ix := 0; ix < len(*ml); ix++ {
			if (*ml)[ix].cmp(genInfo.transMove) {
				continue
			}
			sc := see((*ml)[ix].fr(), (*ml)[ix].to(), b)
			(*ml)[ix].packEval(sc)
			if sc > bs {
				bs = sc
				bIx = ix
			}
		}
		if bs >= 0 {
			mv := (*ml)[bIx]
			(*ml)[bIx], (*ml)[len(*ml)-1] = (*ml)[len(*ml)-1], (*ml)[bIx]
			*ml = (*ml)[:len(*ml)-1]
			genInfo.sv = nextGoodCp
			return mv, "first good capt"
		}

		genInfo.sv = nextK1
		fallthrough
	case nextGoodCp:
		// pick a good capt - use see - not transMove
		bs := -1
		bIx := 0
		ml := &genInfo.captures
		for ix := 0; ix < len(*ml); ix++ {
			if (*ml)[ix].cmp(genInfo.transMove) {
				continue
			}
			sc := (*ml)[ix].eval()
			if sc > bs {
				bs = sc
				bIx = ix
			}
		}
		if bs >= 0 {
			mv := (*ml)[bIx]
			(*ml)[bIx], (*ml)[len(*ml)-1] = (*ml)[len(*ml)-1], (*ml)[bIx]
			*ml = (*ml)[:len(*ml)-1]
			bs, bIx = minEval, -1
			return mv, "good capt"
		}
		genInfo.sv = nextK1
		fallthrough
	case nextK1: // not transMove
		genInfo.sv = nextK2
		if killers[genInfo.ply].k1 != noMove && !genInfo.transMove.cmpFrToP(killers[genInfo.ply].k1) {
			if b.isLegal(killers[genInfo.ply].k1) {
				var mv move
				mv.packMove(killers[genInfo.ply].k1.fr(), killers[genInfo.ply].k1.to(), b.sq[killers[genInfo.ply].k1.fr()], b.sq[killers[genInfo.ply].k1.to()], killers[genInfo.ply].k1.pr(), b.ep, b.castlings)
				return mv, "K1"
			}
		}

		fallthrough
	case nextK2: // not transMove
		genInfo.sv = nextCounterMv
		if killers[genInfo.ply].k2 != noMove && !genInfo.transMove.cmpFrToP(killers[genInfo.ply].k2) {
			if b.isLegal(killers[genInfo.ply].k2) {
				var mv move
				mv.packMove(killers[genInfo.ply].k2.fr(), killers[genInfo.ply].k2.to(), b.sq[killers[genInfo.ply].k2.fr()], b.sq[killers[genInfo.ply].k2.to()], killers[genInfo.ply].k2.pr(), b.ep, b.castlings)
				return mv, "K2"
			}
		}

		fallthrough
	case nextCounterMv: // not transMove, not killer1, not killer2
		genInfo.counterMv = noMove
		genInfo.sv = nextFirstNonCp
		//	if genInfo.counterMv != noMove && !genInfo.counterMv.cmpFrTo(genInfo.transMove) &&
		//     genInfo.counterMv.cmpFrTo(killers[genInfo.ply].k1) &&  genInfo.counterMv.cmpFrTo(killers[genInfo.ply].k2) {
		//		var mv move
		//		mv.packMove(counterMv.fr(), counterMv.to(),b.sq[counterMv.fr()],b.sq[counterMv.to()],counterMv.pr(), b.ep,b.castlings)
		//check if it is a valid move
		//		return sv, counterMovex[genInfo.ply][mv.to()]
		//	}

		fallthrough
	case nextFirstNonCp: // not transMove, not counterMove, not killer1, not killer2
		genInfo.nonCapt.new(50)
		ml := &genInfo.nonCapt
		b.genAllNonCaptures(ml)
		// pick by HistoryTab (see will probably not give anything) - I don't want to sort it. hist may change between moves
		bs := minEval
		bIx := -1
		for ix := 0; ix < len(*ml); ix++ {
			if (*ml)[ix].cmpFrToP(genInfo.transMove) || (*ml)[ix].cmpFrToP(genInfo.counterMv) ||
				(*ml)[ix].cmpFrToP(killers[genInfo.ply].k1) || (*ml)[ix].cmpFrToP(killers[genInfo.ply].k2) {
				continue
			}
			sc := int(history.get((*ml)[ix].fr(), (*ml)[ix].to(), b.stm))
			if sc > bs {
				bs = sc
				bIx = ix
			}
		}
		if bIx >= 0 {
			mv := (*ml)[bIx]

			(*ml)[bIx], (*ml)[len(*ml)-1] = (*ml)[len(*ml)-1], (*ml)[bIx]
			*ml = (*ml)[:len(*ml)-1]

			genInfo.sv = nextNonCp

			return mv, "first non capt"
		}

		genInfo.sv = nextBadCp
		fallthrough
	case nextNonCp: // not transMove, not counterMove, not killer1, not killer2
		// pick by HistoryTab (see will probably not give anything)
		bs := minEval
		bIx := -1
		ml := &genInfo.nonCapt
		for ix := 0; ix < len(*ml); ix++ {
			if (*ml)[ix].cmpFrToP(genInfo.transMove) || (*ml)[ix].cmpFrToP(genInfo.counterMv) ||
				(*ml)[ix].cmpFrToP(killers[genInfo.ply].k1) || (*ml)[ix].cmpFrToP(killers[genInfo.ply].k2) {
							continue
			}
			sc := int(history.get((*ml)[ix].fr(), (*ml)[ix].to(), b.stm))
			if sc > bs {
				bs = sc
				bIx = ix
			}
		}
		if bIx >= 0 {
			mv := (*ml)[bIx]
			(*ml)[bIx], (*ml)[len(*ml)-1] = (*ml)[len(*ml)-1], (*ml)[bIx]
			*ml = (*ml)[:len(*ml)-1]

			return mv, "non Capt"
		}

		genInfo.sv = nextBadCp
		fallthrough
	case nextBadCp: // not transMove
		// pick a bad capt  - use see?
		mv := noMove
		ml := &genInfo.captures
		for ix := len(*ml) - 1; ix >= 0; ix-- {
			if (*ml)[ix].cmp(genInfo.transMove) {
				*ml = (*ml)[:len(*ml)-1]
				continue
			}

			mv = (*ml)[ix]
			//		(*ml)[ix], (*ml)[len(*ml)-1] = (*ml)[len(*ml)-1], (*ml)[ix]
			*ml = (*ml)[:len(*ml)-1]
			break
		}

		return mv, "bad capt"
	default: // shouldn't happen
		panic("neve come here! nextNormal sv=" + strconv.Itoa(genInfo.sv))
	}
}

// StartPerft starts the Perft command that generates all moves until the given depth.
// It counts the leafs only taht is printed out for each possible move from current pos
func startPerft(depth int, bd *boardStruct) uint64 {
	if depth <= 0 {
		fmt.Printf("Total:\t%v\n", 1)
		return 0
	}

	transMove := noMove
	transMove, _, _, _ = trans.retrieve(bd.fullKey(), depth, 0)

	totCount := uint64(0)
	var genInfo = genInfoStruct{sv: 0, ply: 0, transMove: transMove}
	next = nextNormal
	ix := 0
	for mv, msg := next(&genInfo, bd); mv != noMove; mv, msg = next(&genInfo, bd) {
		if !bd.move(mv) {
			continue
		}
		dbg := false
/* 
		/////////////////////////////////////////////////////////////
		if mv.fr() == D4 && mv.to() == F4 {
			dbg = true
		}
		/////////////////////////////////////////////////////////////
 */
		count := perft(dbg, depth-1, 1, bd)
		totCount += count
		fmt.Printf("%2d: %v \t%v \t%v\n", ix+1, mv.String(), count, msg)

		bd.unmove(mv)
		ix++
	}
	fmt.Println("------------------")
	fmt.Println()
	fmt.Printf("Total:\t%v\n", totCount)
	return totCount
}

func perft(dbg bool, depth, ply int, bd *boardStruct) uint64 {
	if depth == 0 {
		return 1
	}

	transMove := noMove
	transMove, _, _, _ = trans.retrieve(bd.fullKey(), depth, ply)
	ix := 0
	count := uint64(0)
	var genInfo = genInfoStruct{sv: 0, ply: ply, transMove: transMove}
	next = nextNormal
	for mv, msg := next(&genInfo, bd); mv != noMove; mv, msg = next(&genInfo, bd) {
		if !bd.move(mv) {
			continue
		}
		_ = msg
		deb := false
/* 
		////////////////////////////////////////////////////////////////
		if dbg && mv.fr() == F5 && mv.to() == F4 {
			deb = true
		}
		if dbg && mv.fr() == E2 && mv.to() == E4 {
			deb = true
		}
		////////////////////////////////////////////////////////////////
 */
		cnt := perft(deb, depth-1, ply+1, bd)
		count += cnt
/*
		/////////////////////////////////////////////
		if dbg && !deb  {
			fmt.Println(ix+1, ":(e4)     ", mv.String(), msg, "\t", cnt)
			if ix==1{
				fmt.Println("K1",killers[ply].k1.StringFull())
				fmt.Println("K2",killers[ply].k2.StringFull())
			}
		}
		////////////////////////////////////////////
*/
		bd.unmove(mv)
		ix++
	}

	return count
}
