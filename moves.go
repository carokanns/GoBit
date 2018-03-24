package main

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
)

var pieceRules [nP][]int // not pawns

func init() {
	pieceRules[Rook] = append(pieceRules[Rook], E)
	pieceRules[Rook] = append(pieceRules[Rook], W)
	pieceRules[Rook] = append(pieceRules[Rook], N)
	pieceRules[Rook] = append(pieceRules[Rook], S)
}

type move uint64
type moveList []move
func (mvs *moveList) add(mv move) {
	*mvs = append(*mvs, mv)
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