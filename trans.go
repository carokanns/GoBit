package main

import (
	"fmt"
	"math/rand"
)

////////////////////////////////////////////////////////
//////////////////////// HASH //////////////////////////
var randPcSq [12 * 64]uint64 // keyvalues for 'pc on sq'
var randEp [8]uint64         // keyvalues for 8 ep files
var randCastl [16]uint64     // keyvalues for castling states

// setup random generator with seed
var rnd = (*rand.Rand)(rand.New(rand.NewSource(1013))) //usage: rnd.Intn(n) NOTE: n > 0

// Rand64 creates one 64 bit random number
func rand64() uint64 {
	rand := uint64(0)

	for i := 0; i < 4; i++ {
		rand = uint64(int(rand<<16) | rnd.Intn(1<<16))
	}

	return rand
}

// initKeys computes random hash keyvalues for pc/sq, ep and castlings
func initKeys() {
	for i := 0; i < 12*64; i++ {
		randPcSq[i] = rand64()
	}
	for i := 0; i < 8; i++ {
		randEp[i] = rand64()
	}
	for i := 0; i < 16; i++ {
		randCastl[i] = rand64()
	}
}

// for color we just flip with XOR ffffffffffffffff
// hash key after change color
func flipSide(key uint64) uint64 {
	return ^key
}

// pcSqKey returns the keyvalue för piece on square
func pcSqKey(pc, sq int) uint64 {
	return randPcSq[pc*sq]
}

// epKey returns the keyvalue for the current ep state
func epKey(epSq int) uint64 {
	if epSq == 0 {
		return 0
	}
	return randEp[epSq%8]
}

// castlKey returns the keyvalue for the current castling state
func castlKey(castling uint) uint64 {
	return randCastl[castling]
}

////////////////////////////////////////////////////////
//////////////////////// TRANS /////////////////////////
const entrySize = 128 / 8

type ttEntry struct {
	lock      uint32 // the lock, extra safety
	move      uint32 // the best move from the search
	_         uint16 // alignement, not used
	score     int16  // the score from the search
	age       uint8  // the age of this entry
	depth     int8   // the depth that the score is based on
	scoreType uint8  // the score has this score type
	_         uint8  // alignement, not used
}

// clear one entry
func (e *ttEntry) clear() {
	//Obs entry skall vara 16 bytes
	e.lock = 0
	e.move = uint32(noMove)
	// entry.utfyllnad = 0  behövs inte
	e.score = 0
	e.age = 0
	e.depth = -1
	e.scoreType = 0
	//entry.tomt = 0   behövs inte
}

type transpStruct struct {
	entries uint // number of entries
	mask    uint // mask for the index
	cntUsed int  // The transposition table usage
	age     int  // current age
	tab     []ttEntry
	// for health tests
	cStores int
	cTried  int
	cFound  int
	cPrune  int
	cBest   int
}

var trans transpStruct

// allocate a new transposition table with the size from GUI
func (t *transpStruct) new(mB int) error {
	if mB > 4000 {
		return fmt.Errorf("max transtable size is 4GB (~4000 MB)")
	}

	byteSize := mB << 20
	bits := sizeToBits(byteSize)

	t.entries = 1 << uint(bits)
	t.mask = t.entries - 1

	t.age = 0
	t.cntUsed = 0

	t.tab = make([]ttEntry, t.entries, t.entries)
	t.clear()
	tell(fmt.Sprintf("info string allocated %v MB to %v entries", len(t.tab)*entrySize/(1024*1024), t.entries))
	return nil
}

// returns how many bits the mask will need to cover the table size
func sizeToBits(size int) uint {
	bits := uint(0)
	for cntEntries := size / entrySize; cntEntries > 1; cntEntries /= 2 {
		bits++
	}

	return bits
}

// clear all entries
func (t *transpStruct) clear() {
	var e ttEntry
	e.clear()

	for i := uint(0); i < t.entries; i++ {
		t.tab[i] = e
	}

	t.age = 0
	t.cntUsed = 0

	//counts
	t.cFound, t.cStores, t.cTried, t.cPrune, t.cBest = 0, 0, 0, 0, 0
}

// index uses the Key to compute an index into the table
func (t *transpStruct) index(fullKey uint64) int64 {
	return int64(fullKey)
}

// Lock extracts the lock value from the hash key
func (t *transpStruct) lock(fullKey uint64) uint32 {
	return uint32(fullKey >> 32)
}

//
func (t *transpStruct) initSearch() {
	t.incAge()
	t.cntUsed = 0

	// Health check counts
	t.cFound, t.cStores, t.cTried, t.cPrune, t.cBest = 0, 0, 0, 0, 0
}

// incAge increments the date for the hahs table.
// We are reborned after the age 255
func (t *transpStruct) incAge() {
	t.age = (t.age + 1) % 256
}

func (b *boardStruct) fullKey() uint64 {
	key := b.key ^ epKey(b.ep)
	key ^= castlKey(uint(b.castlings))
	return key
}

// store current position in the transp table.
// The key is computed from the position. The lock value is the 32 first bits in the key
// From the key we get an index to the table.
// We will try 4 entries in a sequence if a lock is found
// We always try to replace another age and/or a lower searched depth

func (t *transpStruct) store(fullKey uint64, mv move, depth, ply, sc, scoreType int) {
	//	fmt.Println("store:", mv, depth, ply, sc, scoreType)
	t.cStores++
	sc = removeMatePly(sc, ply)

	index := fullKey & uint64(t.mask)
	lock := t.lock(fullKey)

	var newEntry *ttEntry
	bestDep := -2000

	for i := uint64(0); i < 4; i++ {
		idx := (uint64(index) + i) & uint64(t.mask)

		entry := &t.tab[idx]

		if entry.lock == lock {
			if int(entry.age) != t.age {
				entry.age = uint8(t.age)
				t.cntUsed++
			}

			if depth >= int(entry.depth) {
				if mv != noMove {
					entry.move = uint32(mv.onlyMv())
				}
				entry.depth = int8(depth)
				entry.score = int16(sc)
				entry.scoreType = uint8(scoreType)
				return
			}

			if entry.move == uint32(noMove) {
				entry.move = uint32(mv.onlyMv())
			}

			return
		}

		selDepth := -int(entry.depth)
		if entry.age != uint8(t.age) {
			selDepth += 1000
		}

		if selDepth > bestDep {
			newEntry = entry
			bestDep = selDepth
		}
	}

	if newEntry.age != uint8(t.age) {
		t.cntUsed++
	}

	newEntry.lock = lock
	newEntry.age = uint8(t.age)
	newEntry.depth = int8(depth)
	newEntry.move = uint32(mv)
	newEntry.score = int16(sc)
	newEntry.scoreType = uint8(scoreType)
}

// retrieve get move and score to the current position from the transp Table if the key and lock is correct
// if no entry is matching return false else return true, depth not ok return false but with move filled in
// We will try the 4 entries in sequence until lock match otherwise return false
func (t *transpStruct) retrieve(fullKey uint64, depth, ply int) (mv move, sc, scoreType int, ok bool) {
	t.cTried++
	mv = noMove
	ok = false
	sc = noScore
	scoreType = 0

	index := fullKey & uint64(t.mask)
	lock := t.lock(fullKey)

	for i := uint64(0); i < 4; i++ {

		idx := uint(index+i) & t.mask
		entry := &t.tab[idx]

		if entry.lock == lock { // there is a matching position already here
			t.cFound++

			if int(entry.age) != t.age { // from another generation?
				entry.age = uint8(t.age) // touch entry
				t.cntUsed++
			}
			mv = move(entry.move)
			sc = addMatePly(int(entry.score), ply)
			scoreType = int(entry.scoreType)
			ok = true
			if int(entry.depth) >= depth {
				//				fmt.Println("retrieve (key,depth,ply)",fullKey,depth,ply,"(mv,sc,scTyp):", mv, sc, scoreType)
				return
			}

			if isMateScore(sc) {
				scoreType &= ^scoreTypeUpper
				if sc < 0 {
					scoreType &= ^scoreTypeLower
				}
				//				fmt.Println("retrieve (key,depth,ply)",fullKey,depth,ply,"(mv,sc,scTyp):", mv, sc, scoreType)
				return
			}
			//			fmt.Println("Nothing to retrieve-depth? (key,depth,ply):", fullKey, depth,ply)
			ok = false
			return
		}
	}
	//	fmt.Println("Nothing to retrieve - no hit (key,depth,ply):", fullKey, depth,ply)
	ok = false
	return
}

// isMateScore returns true if the score is a mate score
func isMateScore(sc int) bool {
	return sc < minEval+maxPly || sc > maxEval-maxPly
}

// removeMatePly removes ply from the score value (score - ply) if mate
// in order to mix up different depths
func removeMatePly(sc, ply int) int {
	if sc < minEval+maxPly {
		return -mateEval
	} else if sc > maxEval-maxPly {
		return mateEval
	} else {
		return sc
	}
}

// addMatePly adjusts mate value with ply if mate score
func addMatePly(sc, ply int) int {
	if sc < minEval+maxPly {
		return mateEval - ply
	} else if sc > maxEval-maxPly {
		return -mateEval+ply
	}
	return sc
}

const (
	//no scoretype = 0
	scoreTypeLower   = 0x1                             // sc > alpha
	scoreTypeUpper   = 0x2                             // sc < beta
	scoreTypeBetween = scoreTypeLower | scoreTypeUpper // alpha < sc < beta
)

// scoreType sets if it is an upper or lower score
func scoreType(sc, alpha, beta int) int {

	scoreType := 0
	if sc > alpha {
		scoreType |= scoreTypeLower
	}
	if sc < beta {
		scoreType |= scoreTypeUpper
	}

	return scoreType
}
