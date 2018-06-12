package main

import "fmt"

// directions
const (
	E  = +1
	W  = -1
	N  = 8
	S  = -8
	NW = +7
	NE = +9
	SW = -NE
	SE = -NW

	frMask     = 0x0000003f                 //0000 0000  0000 0000  0000 0000  0011 1111
	toMask     = 0x00000fd0                 //0000 0000  0000 0000  0000 1111  1100 0000
	pcMask     = 0x0000f000                 //0000 0000  0000 0000  1111 0000  0000 0000
	cpMask     = 0x000f0000                 //0000 0000  0000 1111  0000 0000  0000 0000
	prMask     = 0x00f00000                 //0000 0000  1111 0000  0000 0000  0000 0000
	epMask     = 0x0f000000                 //0000 1111  0000 0000  0000 0000  0000 0000
	castlMask  = 0xf0000000                 //1111 0000  0000 0000  0000 0000  0000 0000
	evalMask   = uint64(0xffff000000000000) // The 16 first bits in uint64
	toShift    = 6
	pcShift    = 12 //6+6
	cpShift    = 16 //6+6+4
	prShift    = 20 //6+6+4+4
	epShift    = 24 //6+6+4+4+4
	castlShift = 28 //6+6+4+4+4+4
	evalShift  = 64 - 16
	noMove     = move(0)
)

var pieceRules [nPt][]int // not pawns

type move uint64

func (m move) String() string {
	s := m.StringFull()
	s = s[1:3] + s[5:]
	return s
}
func (m move) StringFull() string {
	fr := sq2Fen[int(m.fr())]
	to := sq2Fen[int(m.to())]
	pc := pc2Fen(int(m.pc()))
	cp := pc2Fen(int(m.cp())) + " "
	pr := pc2Fen(int(m.pr()))
	return trim(fmt.Sprintf("%v%v-%v%v%v", pc, fr, cp[:1], to, pr))
}

func (m *move) packMove(fr, to, pc, cp, pr, epSq int, castl castlings) {
	// 6 bits fr, 6 bits to, 4 bits pc, 4 bits cp, 4 bits prom, 4 bits ep, 4 bits castl = 32 bits
	if epSq == empty{panic("wtf ep")}
	epFile := 0
	if epSq != 0 {
		epFile = epSq%8 + 1
	}
	*m = move(fr | (to << toShift) | (pc << pcShift) |
		(cp << cpShift) | (pr << prShift) | (epFile << epShift) | int(castl<<castlShift))
}
func (m *move) packEval(score int) {
	(*m) &= move(^evalMask) //clear eval
	(*m) |= move(score+30000) << evalShift
}

// compare two moves - only frSq and toSq
func (m move) cmpFrTo(m2 move) bool {
	// return (m & move(^evalMask)) == (m2 & move(^evalMask))
	return m.fr() == m2.fr() && m.to() == m2.to()
}

// compare two moves - only frSq, toSq and pc
func (m move) cmpFrToP(m2 move) bool {
	// return (m & move(^evalMask)) == (m2 & move(^evalMask))
	return m.fr() == m2.fr() && m.to() == m2.to() && m.pc() == m2.pc()
}

// compare two moves
func (m move) cmp(m2 move) bool {
	return (m & move(^evalMask)) == (m2 & move(^evalMask))
}

func (m move) eval() int {
	return int((uint(m)&uint(evalMask))>>evalShift) - 30000
}
func (m move) fr() int {
	return int(m & frMask)
}
func (m move) to() int {
	return int(m&toMask) >> toShift
}

func (m move) pc() int {
	return int(m&pcMask) >> pcShift
}
func (m move) cp() int {
	return int(m&cpMask) >> cpShift
}
func (m move) pr() int {
	return int(m&prMask) >> prShift
}

func (m move) ep(sd color) int {
	// sd is the side that can capture
	file := int(m&epMask) >> epShift
	if file == 0 {
		return 0 // no ep
	}

	// there is an ep sq
	rank := 5
	if sd == BLACK {
		rank = 2
	}

	return rank*8 + file - 1

}

func (m move) castl() castlings {
	return castlings(m&castlMask) >> castlShift
}

//move without eval
func (m move) onlyMv() move {
	return m & move(^evalMask)
}

type moveList []move
func (ml *moveList) new(size int) {
	*ml = make(moveList, 0, size)
}

func (ml *moveList) clear() {
	*ml = (*ml)[:0]
}

func (ml *moveList) add(mv move) {
	*ml = append(*ml, mv)
}

func (ml *moveList) remove(ix int) {
	if len(*ml) > ix && ix >= 0 {
		*ml = append((*ml)[:ix], (*ml)[ix+1:]...)
	}
}

// Sort is sorting the moves in the Score/Move list according to the score per move
func (ml *moveList) sort() {
	bSwap := true
	for bSwap {
		bSwap = false
		for i := 0; i < len(*ml)-1; i++ {
			if (*ml)[i+1].eval() > (*ml)[i].eval() {
				(*ml)[i], (*ml)[i+1] = (*ml)[i+1], (*ml)[i]
				bSwap = true
			}
		}
	}
}

func (ml moveList) String() string {
	theString := ""
	for _, mv := range ml {
		theString += mv.String() + " "
	}
	return theString
}
