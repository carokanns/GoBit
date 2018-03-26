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

	frMask     = 0x00000003f                //00 0000 0000  0000 0000  0000 0000  0011 1111
	toMask     = 0x000000fd0                //00 0000 0000  0000 0000  0000 1111  1100 0000
	p12Mask    = 0x00000f000                //00 0000 0000  0000 0000  1111 0000  0000 0000
	cpMask     = 0x0000f0000                //00 0000 0000  0000 1111  0000 0000  0000 0000
	prMask     = 0x000f00000                //00 0000 0000  1111 0000  0000 0000  0000 0000
	epMask     = 0x03f000000                //00 0011 1111  0000 0000  0000 0000  0000 0000
	castlMask  = 0x3c0000000                //11 1100 0000  0000 0000  0000 0000  0000 0000
	evalMask   = uint64(0xffff000000000000) // The 16 first bits in uint64
	toShift    = 6
	p12Shift   = 12 //6+6
	cpShift    = 16 //6+6+4
	prShift    = 20 //6+6+4+4
	epShift    = 24 //6+6+4+4+4
	castlShift = 30 //6+6+4+4+4+6
	evalShift  = 64 - 16
)

var pieceRules [nP][]int // not pawns

func init() {
	pieceRules[Rook] = append(pieceRules[Rook], E)
	pieceRules[Rook] = append(pieceRules[Rook], W)
	pieceRules[Rook] = append(pieceRules[Rook], N)
	pieceRules[Rook] = append(pieceRules[Rook], S)
}

type move uint64

func (m move) String() string {
	s := m.StringFull()
	s = s[1:3] + s[5:]
	return s
}
func (m move) StringFull() string {
	fr := sq2Fen[int(m.fr())]
	to := sq2Fen[int(m.to())]
	p := int2Fen(int(m.p12()))
	cp := int2Fen(int(m.cp())) + " "
	pr := int2Fen(int(m.pr()))
	return fmt.Sprintf("%v%v-%v%v%v", p, fr, cp[:1], to, pr)
}


func (m *move) packMove(fr, to, p12, cp, pr, epSq, castl uint) {
	// 6 bits fr, 6 bits to, 4 bits p12, 4 bits cp, 4 bits prom, 4 bits ep, 4 bits castl = 32 bits
	*m = move(fr | (to << toShift) | (p12 << p12Shift) |
		(cp << cpShift) | (pr << prShift) | (epSq << epShift) | (castl << castlShift))
}

func (m move) fr() uint {
	return uint(m & frMask)
}
func (m move) to() uint {
	return uint(m&toMask) >> toShift
}

func (m move) p12() uint {
	return uint(m&p12Mask) >> p12Shift
}
func (m move) cp() uint {
	return uint(m&cpMask) >> cpShift
}
func (m move) pr() uint {
	return uint(m&prMask) >> prShift
}
func (m move) ep() uint {
	return uint(m&epMask) >> epShift
}
func (m move) castl() castlings {
	return castlings(m&castlMask) >> castlShift
}

type moveList []move

func (mvs *moveList) add(mv move) {
	*mvs = append(*mvs, mv)
}

func (mvs moveList) String() string {
	theString := ""
	for _, mv := range mvs {
		theString += mv.String() + " "
	}
	return theString
}

var ml = moveList{}

///////////////////////////////
/*
func (m move) pack(fr, to, cp, pr, ep, castl int) move {
	packed := uint32(0)
	packed = uint32(fr)
	packed |= uint32(to)<<6  // fr needs 6 bits first
	packed |= uint32(cp)<<(6+6) // fr+to needs 12 bits first
	packed |= uint32(pr)<<(6+6+4) // cp needs 4 bits (with color = p12)
	packed |= uint32(ep)<<(6+6+4+4) // pr needs 4 bits
	packed |= uint32(castl)<<(6+6+4+4+6) // ep needs 6 bits ( and castl needs 4)
	// we use 6+6+4+4+6+4 = 30  That means we have 32 bits left for move value
	return move(packed)
}
*/
