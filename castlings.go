package main

import (
	
	"strings"
)

type castlings uint

const (
	shortW = uint(0x1) // white can castle short
	longW  = uint(0x2) // white can castle long
	shortB = uint(0x4) // black can castle short
	longB  = uint(0x8) // black can castle short
)

type castlOptions struct {
	short                                uint // flag
	long                                 uint // flag
	rook                                 int  // rook pc (wR/bR)
	kingPos                              int  // king pos
	rookSh                               uint // rook pos short
	rookL                                uint // rook pos long
	betweenSh                            bitBoard
	betweenL                             bitBoard
	pawnsSh, pawnsL, knightsSh, knightsL bitBoard
}

var castl = [2]castlOptions{
	{shortW, longW, wR, E1, H1, A1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{shortB, longB, bR, E8, H8, A8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
}

// castling privileges
func (c castlings) flags(sd color) bool {
	return c.shortFlag(sd) || c.longFlag(sd)
}
func (c castlings) shortFlag(sd color) bool {
	return (castl[sd].short & uint(c)) != 0
}
func (c castlings) longFlag(sd color) bool {
	return (castl[sd].long & uint(c)) != 0
}

func (c *castlings) on(val uint) {
	(*c) |= castlings(val)
}

func (c *castlings) off(val uint) {
	(*c) &= castlings(^val)
}

func (c castlings) String() string {
	flags := ""
	if uint(c)&shortW != 0 {
		flags = "K"
	}
	if uint(c)&longW != 0 {
		flags += "Q"
	}
	if uint(c)&shortB != 0 {
		flags += "k"
	}
	if uint(c)&longB != 0 {
		flags += "q"
	}
	if flags == "" {
		flags = "-"
	}
	return flags
}

// parse castling rights in fenstring
func parseCastlings(fenCastl string) castlings {
	c := uint(0)
	if fenCastl == "-" {
		return castlings(0)
	}

	if strings.Index(fenCastl, "K") >= 0 {
		c |= shortW
	}
	if strings.Index(fenCastl, "Q") >= 0 {
		c |= longW
	}
	if strings.Index(fenCastl, "k") >= 0 {
		c |= shortB
	}
	if strings.Index(fenCastl, "q") >= 0 {
		c |= longB
	}

	return castlings(c)
}

func initCastlings() {
	// squares between K and R short castling
	castl[WHITE].betweenSh.set(F1)
	castl[WHITE].betweenSh.set(G1)
	castl[BLACK].betweenSh.set(F8)
	castl[BLACK].betweenSh.set(G8)

	// squares between K and R long castling
	castl[WHITE].betweenL.set(B1)
	castl[WHITE].betweenL.set(C1)
	castl[WHITE].betweenL.set(D1)
	castl[BLACK].betweenL.set(B8)
	castl[BLACK].betweenL.set(C8)
	castl[BLACK].betweenL.set(D8)

	// pawns stop short castling W
	castl[WHITE].pawnsSh.set(D2)
	castl[WHITE].pawnsSh.set(E2)
	castl[WHITE].pawnsSh.set(F2)
	castl[WHITE].pawnsSh.set(G2)
	castl[WHITE].pawnsSh.set(H2)
	// pawns stop long castling W
	castl[WHITE].pawnsL.set(B2)
	castl[WHITE].pawnsL.set(C2)
	castl[WHITE].pawnsL.set(D2)
	castl[WHITE].pawnsL.set(E2)
	castl[WHITE].pawnsL.set(F2)

	// pawns stop short castling B
	castl[BLACK].pawnsSh.set(D7)
	castl[BLACK].pawnsSh.set(E7)
	castl[BLACK].pawnsSh.set(F7)
	castl[BLACK].pawnsSh.set(G7)
	castl[BLACK].pawnsSh.set(H7)
	// pawns stop long castling B
	castl[BLACK].pawnsL.set(B7)
	castl[BLACK].pawnsL.set(C7)
	castl[BLACK].pawnsL.set(D7)
	castl[BLACK].pawnsL.set(E7)
	castl[BLACK].pawnsL.set(F7)

	// knights stop short castling W
	castl[WHITE].knightsSh.set(C2)
	castl[WHITE].knightsSh.set(D2)
	castl[WHITE].knightsSh.set(E2)
	castl[WHITE].knightsSh.set(G2)
	castl[WHITE].knightsSh.set(H2)
	castl[WHITE].knightsSh.set(D3)
	castl[WHITE].knightsSh.set(E3)
	castl[WHITE].knightsSh.set(F3)
	castl[WHITE].knightsSh.set(G3)
	castl[WHITE].knightsSh.set(H3)
	// knights stop long castling W
	castl[WHITE].knightsL.set(A2)
	castl[WHITE].knightsL.set(B2)
	castl[WHITE].knightsL.set(C2)
	castl[WHITE].knightsL.set(E2)
	castl[WHITE].knightsL.set(F2)
	castl[WHITE].knightsL.set(G2)
	castl[WHITE].knightsL.set(B3)
	castl[WHITE].knightsL.set(C3)
	castl[WHITE].knightsL.set(D3)
	castl[WHITE].knightsL.set(E3)
	castl[WHITE].knightsL.set(F3)

	// knights stop short castling B
	castl[BLACK].knightsSh.set(C7)
	castl[BLACK].knightsSh.set(D7)
	castl[BLACK].knightsSh.set(E7)
	castl[BLACK].knightsSh.set(G7)
	castl[BLACK].knightsSh.set(H7)
	castl[BLACK].knightsSh.set(D6)
	castl[BLACK].knightsSh.set(E6)
	castl[BLACK].knightsSh.set(F6)
	castl[BLACK].knightsSh.set(G6)
	castl[BLACK].knightsSh.set(H6)
	// knights stop long castling B
	castl[BLACK].knightsL.set(A7)
	castl[BLACK].knightsL.set(B7)
	castl[BLACK].knightsL.set(C7)
	castl[BLACK].knightsL.set(E7)
	castl[BLACK].knightsL.set(F7)
	castl[BLACK].knightsL.set(G7)
	castl[BLACK].knightsL.set(B6)
	castl[BLACK].knightsL.set(C6)
	castl[BLACK].knightsL.set(D6)
	castl[BLACK].knightsL.set(E6)
	castl[BLACK].knightsL.set(F6)

}
