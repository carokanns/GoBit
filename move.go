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
type moveList struct {
	mv []move
}
func init() {
	pieceRules[Rook] = append(pieceRules[Rook], E)
	pieceRules[Rook] = append(pieceRules[Rook], W)
	pieceRules[Rook] = append(pieceRules[Rook], N)
	pieceRules[Rook] = append(pieceRules[Rook], S)
}

func (ml *moveList) add(mv move) {
	ml.mv = append(ml.mv, mv)
}

type move uint64

var ml = moveList{}
