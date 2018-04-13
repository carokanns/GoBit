package main

import (
	"fmt"
	"strconv"
	"strings"
)

// various consts
const (
	nP12     = 12
	nP       = 6
	WHITE    = color(0)
	BLACK    = color(1)
	startpos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - "
	row1     = bitBoard(0x00000000000000FF)
	row2     = bitBoard(0x000000000000FF00)
	row3     = bitBoard(0x0000000000FF0000)
	row4     = bitBoard(0x00000000FF000000)
	row5     = bitBoard(0x000000FF00000000)
	row6     = bitBoard(0x0000FF0000000000)
	row7     = bitBoard(0x00FF000000000000)
	row8     = bitBoard(0xFF00000000000000)
	fileA    = bitBoard(0x0101010101010101)
	fileB    = bitBoard(0x0202020202020202)
	fileG    = bitBoard(0x4040404040404040)
	fileH    = bitBoard(0x8080808080808080)
)

var atksKnights [64]bitBoard
var atksKings [64]bitBoard

// initialize all possible knight attacks
func initAtksKnights() {
	for fr := A1; fr <= H8; fr++ {
		toBB := bitBoard(0)
		rk := fr / 8
		fl := fr % 8
		// NNE  2,1
		if rk+2 < 8 && fl+1 < 8 {
			to := uint((rk+2)*8 + fl + 1)
			toBB.set(to)
		}

		// ENE  1,2
		if rk+1 < 8 && fl+2 < 8 {
			to := uint((rk+1)*8 + fl + 2)
			toBB.set(to)
		}

		// ESE  -1,2
		if rk-1 >= 0 && fl+2 < 8 {
			to := uint((rk-1)*8 + fl + 2)
			toBB.set(to)
		}

		// SSE  -2,+1
		if rk-2 >= 0 && fl+1 < 8 {
			to := uint((rk-2)*8 + fl + 1)
			toBB.set(to)
		}

		// NNW  2,-1
		if rk+2 < 8 && fl-1 >= 0 {
			to := uint((rk+2)*8 + fl - 1)
			toBB.set(to)
		}

		// WNW  1,-2
		if rk+1 < 8 && fl-2 >= 0 {
			to := uint((rk+1)*8 + fl - 2)
			toBB.set(to)
		}

		// WSW  -1,-2
		if rk-1 >= 0 && fl-2 >= 0 {
			to := uint((rk-1)*8 + fl - 2)
			toBB.set(to)
		}

		// SSW  -2,-1
		if rk-2 >= 0 && fl-1 >= 0 {
			to := uint((rk-2)*8 + fl - 1)
			toBB.set(to)
		}
		atksKnights[fr] = toBB
	}
}

// initialize all possible King attacks
func initAtksKings() {
	fmt.Println("init atksKings")

	for fr := A1; fr <= H8; fr++ {
		toBB := bitBoard(0)
		rk := fr / 8
		fl := fr % 8
		//N 1,0
		if rk+1 < 8 {
			to := uint((rk+1)*8 + fl)
			toBB.set(to)
		}

		//NE 1,1
		if rk+1 < 8 && fl+1 < 8 {
			to := uint((rk+1)*8 + fl + 1)
			toBB.set(to)
		}

		//E   0,1
		if fl+1 < 8 {
			to := uint((rk)*8 + fl + 1)
			toBB.set(to)
		}

		//SE -1,1
		if rk-1 >= 0 && fl+1 < 8 {
			to := uint((rk-1)*8 + fl + 1)
			toBB.set(to)
		}

		//S  -1,0
		if rk-1 >= 0 {
			to := uint((rk-1)*8 + fl)
			toBB.set(to)
		}

		//SW -1,-1
		if rk-1 >= 0 && fl-1 >= 0 {
			to := uint((rk-1)*8 + fl - 1)
			toBB.set(to)
		}

		//W   0,-1
		if fl-1 >= 0 {
			to := uint((rk)*8 + fl - 1)
			toBB.set(to)
		}

		//NW  1,-1
		if rk+1 < 8 && fl-1 >= 0 {
			to := uint((rk+1)*8 + fl - 1)
			toBB.set(to)
		}
		atksKings[fr] = toBB
	}
}

type boardStruct struct {
	sq      [64]int
	wbBB    [2]bitBoard
	pieceBB [nP]bitBoard
	King    [2]int
	ep      int
	castlings
	stm    color
	count  [12]int
	rule50 int //set to 0 if a pawn or capt move otherwise increment
}
type color int

var board = boardStruct{}

func (b *boardStruct) allBB() bitBoard {
	return b.wbBB[0] | b.wbBB[1]
}

// clear the board, flags, bitboards etc
func (b *boardStruct) clear() {
	b.stm = WHITE
	b.rule50 = 0
	b.sq = [64]int{}
	b.King = [2]int{}
	b.ep = 0
	b.castlings = 0

	for ix := A1; ix <= H8; ix++ {
		b.sq[ix] = empty
	}

	for ix := 0; ix < nP12; ix++ {
		b.count[ix] = 0
	}

	// bitBoards
	b.wbBB[WHITE], b.wbBB[BLACK] = 0, 0
	for ix := 0; ix < nP; ix++ {
		b.pieceBB[ix] = 0
	}
}

// make a move
func (b *boardStruct) move(mv move) bool {
	newEp := 0
	// we assume that the move is legally correct (except inChekc())
	fr := int(mv.fr())
	to := int(mv.to())
	pr := int(mv.pr())
	p12 := b.sq[fr]
	switch {
	case p12 == wK:
		b.castlings.off(shortW | longW)
		if abs(int(to)-int(fr)) == 2 {
			if to == G1 {
				b.setSq(wR, F1)
				b.setSq(empty, H1)
			} else {
				b.setSq(wR, D1)
				b.setSq(empty, A1)
			}
		}
	case p12 == bK:
		b.castlings.off(shortB | longB)
		if abs(int(to)-int(fr)) == 2 {
			if to == G8 {
				b.setSq(bR, F8)
				b.setSq(empty, H8)
			} else {
				b.setSq(bR, D8)
				b.setSq(empty, A8)
			}
		}
	case p12 == wR:
		if fr == A1 {
			b.off(longW)
		} else if fr == H1 {
			b.off(shortW)
		}
	case p12 == bR:
		if fr == A8 {
			b.off(longB)
		} else if fr == H8 {
			b.off(shortB)
		}

	case p12 == wP && b.sq[to] == empty: // ep move or set ep
		if to-fr == 16 {
			newEp = fr + 8
		} else if to-fr == 7 { // must be ep
			b.setSq(empty, to-8)
		} else if to-fr == 9 { // must be ep
			b.setSq(empty, to-8)
		}
	case p12 == bP && b.sq[to] == empty: //  ep move or set ep
		if fr-to == 16 {
			newEp = to + 8
		} else if fr-to == 7 { // must be ep
			b.setSq(empty, to+8)
		} else if fr-to == 9 { // must be ep
			b.setSq(empty, to+8)
		}
	}
	b.ep = newEp
	b.setSq(empty, fr)

	if pr != empty {
		b.setSq(pr, to)
	} else {
		b.setSq(p12, to)
	}

	b.stm = b.stm ^ 0x1
	if b.isAttacked(b.King[b.stm ^0x1], b.stm) {
		b.unmove(mv)
		return false
	}

	return true
}

func (b *boardStruct) unmove(mv move) {
	b.ep = int(mv.ep())
	b.castlings = mv.castl()
	p12 := int(mv.p12())
	fr := int(mv.fr())
	to := int(mv.to())
	b.setSq(int(mv.cp()), to)
	b.setSq(p12, fr)

	if piece(p12) == Pawn {
		if to == b.ep { // ep move
			b.setSq(empty, to)
			switch fr - to {
			case NW, NE:
				b.setSq(bP, to-N)
			case SW, SE:
				b.setSq(wP, to-S)
			}
		}
	} else if piece(p12) == King {
		sd := p12Color(p12)
		if fr-to == 2 { // long castling
			b.setSq(castl[sd].rook, int(castl[sd].rookL))
			b.setSq(empty, fr-1)
		} else if fr-to == -2 { // short castling
			b.setSq(castl[sd].rook, int(castl[sd].rookSh))
			b.setSq(empty, fr+1)
		}
	}
	b.stm = b.stm ^ 0x1
}

func (b *boardStruct) setSq(p12, sq int) {
	p := piece(p12)
	sd := p12Color(p12)

	if b.sq[sq] != empty { // capture
		cp := b.sq[sq]
		b.count[cp]--
		b.wbBB[sd^0x1].clr(uint(sq))
		b.pieceBB[piece(cp)].clr(uint(sq))
	}
	b.sq[sq] = p12

	if p12 == empty {
		b.wbBB[WHITE].clr(uint(sq))
		b.wbBB[BLACK].clr(uint(sq))
		for p := 0; p < nP; p++ {
			b.pieceBB[p].clr(uint(sq))
		}
		return
	}

	b.count[p12]++

	if p == King {
		b.King[sd] = sq
	}

	b.wbBB[sd].set(uint(sq))
	b.pieceBB[p].set(uint(sq))
}

func (b *boardStruct) newGame() {
	b.stm = WHITE
	b.clear()
	parseFEN(startpos)
}

func (b *boardStruct) genRookMoves(ml *moveList) {
	sd := b.stm
	allRBB := b.pieceBB[Rook] & b.wbBB[sd]
	p12 := uint(pc2P12(Rook, color(sd)))
	var mv move
	for fr := allRBB.firstOne(); fr != 64; fr = allRBB.firstOne() {
		toBB := mRookTab[fr].atks(b) & (^b.wbBB[sd])
		for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {
			mv.packMove(uint(fr), uint(to), p12, uint(b.sq[to]), empty, uint(b.ep), uint(b.castlings))
			ml.add(mv)
		}
	}
}

func (b *boardStruct) genBishopMoves(ml *moveList) {
	sd := b.stm
	allBBB := b.pieceBB[Bishop] & b.wbBB[sd]
	p12 := uint(pc2P12(Bishop, color(sd)))
	ep := uint(b.ep)
	castlings := uint(b.castlings)
	var mv move

	for fr := allBBB.firstOne(); fr != 64; fr = allBBB.firstOne() {
		toBB := mBishopTab[fr].atks(b) & (^b.wbBB[sd])
		for to := toBB.lastOne(); to != 64; to = toBB.lastOne() {
			mv.packMove(uint(fr), uint(to), p12, uint(b.sq[to]), empty, ep, castlings)
			ml.add(mv)
		}
	}
}

func (b *boardStruct) genQueenMoves(mlq *moveList) {
	sd := b.stm
	allQBB := b.pieceBB[Queen] & b.wbBB[sd]
	p12 := uint(pc2P12(Queen, color(sd)))
	ep := uint(b.ep)
	castlings := uint(b.castlings)
	var mv move

	for fr := allQBB.firstOne(); fr != 64; fr = allQBB.firstOne() {
		toBB := mBishopTab[fr].atks(b) & (^b.wbBB[sd])
		toBB |= mRookTab[fr].atks(b) & (^b.wbBB[sd])
		for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {
			mv.packMove(uint(fr), uint(to), p12, uint(b.sq[to]), empty, ep, castlings)
			mlq.add(mv)
		}
	}
}

func (b *boardStruct) genKnightMoves(ml *moveList) {
	sd:=b.stm
	allNBB := b.pieceBB[Knight] & b.wbBB[sd]
	p12 := uint(pc2P12(Knight, color(sd)))
	ep := uint(b.ep)
	castlings := uint(b.castlings)
	var mv move
	for fr := allNBB.firstOne(); fr != 64; fr = allNBB.firstOne() {
		toBB := atksKnights[fr] & (^b.wbBB[sd])
		for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {
			mv.packMove(uint(fr), uint(to), p12, uint(b.sq[to]), empty, ep, castlings)
			ml.add(mv)
		}
	}
}

func (b *boardStruct) genKingMoves(ml *moveList) {
	sd:=b.stm
	// 'normal' moves
	p12 := uint(pc2P12(King, color(sd)))
	ep := uint(b.ep)
	castlings := uint(b.castlings)
	var mv move

	toBB := atksKings[b.King[sd]] & (^b.wbBB[sd])
	for to := toBB.firstOne(); to != 64; to = toBB.firstOne() {
		mv.packMove(uint(b.King[sd]), uint(to), p12, uint(b.sq[to]), empty, ep, castlings)
		ml.add(mv)
	}

	if b.King[sd] == castl[sd].kingPos { // NOTE: Maybe not needed. We should know that the king is there if the flags are ok
		// short castling
		if b.sq[castl[sd].rookSh] == castl[sd].rook && // NOTE: Maybe not needed. We should know that the rook is there if the flags are ok
			(castl[sd].betweenSh&b.allBB()) == 0 {
			if b.castlingShOk(sd) {
				// TODO: not in check and between not attacked
				mv.packMove(uint(b.King[sd]), uint(b.King[sd]+2), uint(b.sq[b.King[sd]]), empty, empty, uint(b.ep), uint(b.castlings))
				ml.add(mv)
			}
		}

		// long castling
		if b.sq[castl[sd].rookL] == castl[sd].rook && // NOTE: Maybe not needed. We should know that the rook is there if the flags are ok
			(castl[sd].betweenL&b.allBB()) == 0 {
			// TODO: not in check and between not attacked
			if b.castlingLOk(sd) {
				mv.packMove(uint(b.King[sd]), uint(b.King[sd]-2), uint(b.sq[b.King[sd]]), empty, empty, uint(b.ep), uint(b.castlings))
				ml.add(mv)
			}
		}
	}
}

// check if short castlings is legal
func (b *boardStruct) castlingShOk(sd color) bool {
	opp := sd ^ 0x1
	if castl[sd].pawnsSh&b.pieceBB[Pawn]&b.wbBB[opp] != 0 {
		return false
	}
	if castl[sd].pawnsSh&b.pieceBB[King]&b.wbBB[opp] != 0 {
		return false
	}
	if castl[sd].knightsSh&b.pieceBB[Knight]&b.wbBB[opp] != 0 {
		return false
	}

	// sliding e1/e8	//NOTE: Maybe not needed during search because we know if we are in check
	sq := b.King[sd]
	if (mBishopTab[sq].atks(b) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq].atks(b) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// sliding f1/f8
	if (mBishopTab[sq+1].atks(b) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq+1].atks(b) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// sliding g1/g8		//NOTE: Maybe not needed because we always make inCheck() before a move
	if (mBishopTab[sq+2].atks(b) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq+2].atks(b) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	return true
}

// check if long castlings is legal
func (b *boardStruct) castlingLOk(sd color) bool {
	opp := sd ^ 0x1
	if castl[sd].knightsL&b.pieceBB[Knight]&b.wbBB[opp] != 0 {
		return false
	}

	// sliding e1/e8
	sq := b.King[sd]
	if (mBishopTab[sq].atks(b) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq].atks(b) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// sliding d1/d8
	if (mBishopTab[sq-1].atks(b) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq-1].atks(b) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}

	// sliding c1/c8	//NOTE: Maybe not needed because we always make inCheck() before a move
	if (mBishopTab[sq-2].atks(b) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	if (mRookTab[sq-2].atks(b) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[opp]) != 0 {
		return false
	}
	return true
}

var genPawns = [2]func(*boardStruct, *moveList){(*boardStruct).genWPawnMoves, (*boardStruct).genBPawnMoves}

func (b *boardStruct) genPawnMoves(ml *moveList) {
	genPawns[b.stm](b, ml)
}
func (b *boardStruct) genWPawnMoves(ml *moveList) {
	var mv move
	wPawns := b.pieceBB[Pawn] & b.wbBB[WHITE]

	// one step
	to1Step := (wPawns << N) & ^b.allBB()
	// two steps,
	to2Step := ((to1Step & row3) << N) & ^b.allBB()
	// captures
	toCapL := ((wPawns & ^fileA) << NW) & b.wbBB[BLACK]
	toCapR := ((wPawns & ^fileH) << NE) & b.wbBB[BLACK]
	// prom
	prom := (to1Step | toCapL | toCapR) & row8

	if prom != 0 {
		for to := prom.firstOne(); to != 64; to = prom.firstOne() {
			cp := empty
			if b.sq[to] != empty {
				cp = b.sq[to]
				if toCapL.test(uint(to)) {
					fr := to - NW
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wQ, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wR, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wN, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wB, uint(b.ep), uint(b.castlings))
					ml.add(mv)
				}
				if toCapR.test(uint(to)) {
					fr := to - NE
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wQ, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wR, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wN, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), wP, uint(cp), wB, uint(b.ep), uint(b.castlings))
					ml.add(mv)
				}
			} else {
				fr := to - N
				mv.packMove(uint(fr), uint(to), wP, uint(cp), wQ, uint(b.ep), uint(b.castlings))
				ml.add(mv)
				mv.packMove(uint(fr), uint(to), wP, uint(cp), wR, uint(b.ep), uint(b.castlings))
				ml.add(mv)
				mv.packMove(uint(fr), uint(to), wP, uint(cp), wN, uint(b.ep), uint(b.castlings))
				ml.add(mv)
				mv.packMove(uint(fr), uint(to), wP, uint(cp), wB, uint(b.ep), uint(b.castlings))
				ml.add(mv)
			}
		}
		to1Step &= ^row8
		toCapL &= ^row8
		toCapR &= ^row8

	}
	// ep move
	if b.ep != 0 {
		epBB := bitBoard(1) << uint(b.ep)
		// ep left
		epToL := ((wPawns & ^fileA) << NW) & epBB
		if epToL != 0 {
			mv.packMove(uint(b.ep-NW), uint(b.ep), wP, bP, empty, uint(b.ep), uint(b.castlings))
			ml.add(mv)
		}
		epToR := ((wPawns & ^fileH) << NE) & epBB
		if epToR != 0 {
			mv.packMove(uint(b.ep-NE), uint(b.ep), wP, bP, empty, uint(b.ep), uint(b.castlings))
			ml.add(mv)
		}
	}
	// Add one step forward
	for to := to1Step.firstOne(); to != 64; to = to1Step.firstOne() {
		mv.packMove(uint(to-N), uint(to), wP, empty, empty, uint(b.ep), uint(b.castlings))
		ml.add(mv)
	}
	// Add two steps forward
	for to := to2Step.firstOne(); to != 64; to = to2Step.firstOne() {
		mv.packMove(uint(to-2*N), uint(to), wP, empty, empty, uint(b.ep), uint(b.castlings))
		ml.add(mv)
	}

	// add Captures left
	for to := toCapL.firstOne(); to != 64; to = toCapL.firstOne() {
		mv.packMove(uint(to-NW), uint(to), wP, uint(b.sq[to]), empty, uint(b.ep), uint(b.castlings))
		ml.add(mv)
	}

	// add Captures right
	for to := toCapR.firstOne(); to != 64; to = toCapR.firstOne() {
		mv.packMove(uint(to-NE), uint(to), wP, uint(b.sq[to]), empty, uint(b.ep), uint(b.castlings))
		ml.add(mv)
	}
}

func (b *boardStruct) genBPawnMoves(ml *moveList) {
	var mv move
	bPawns := b.pieceBB[Pawn] & b.wbBB[BLACK]

	// one step
	to1Step := (bPawns >> (-S)) & ^b.allBB()
	// two steps,
	to2Step := ((to1Step & row6) >> (-S)) & ^b.allBB()
	// captures
	toCapL := ((bPawns & ^fileA) >> (-SW)) & b.wbBB[WHITE]
	toCapR := ((bPawns & ^fileH) >> (-SE)) & b.wbBB[WHITE]
	// prom
	prom := (to1Step | toCapL | toCapR) & row1
	if prom != 0 {
		for to := prom.firstOne(); to != 64; to = prom.firstOne() {
			cp := empty
			if b.sq[to] != empty {
				cp = b.sq[to]
				if toCapL.test(uint(to)) {
					fr := to - SW
					mv.packMove(uint(fr), uint(to), bP, uint(cp), bQ, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), bP, uint(cp), bR, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), bP, uint(cp), bN, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), bP, uint(cp), bB, uint(b.ep), uint(b.castlings))
					ml.add(mv)
				}
				if toCapR.test(uint(to)) {
					fr := to - SE
					mv.packMove(uint(fr), uint(to), bP, uint(cp), bQ, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), bP, uint(cp), bR, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), bP, uint(cp), bN, uint(b.ep), uint(b.castlings))
					ml.add(mv)
					mv.packMove(uint(fr), uint(to), bP, uint(cp), bB, uint(b.ep), uint(b.castlings))
					ml.add(mv)
				}
			} else {
				fr := to - S
				mv.packMove(uint(fr), uint(to), bP, uint(cp), bQ, uint(b.ep), uint(b.castlings))
				ml.add(mv)
				mv.packMove(uint(fr), uint(to), bP, uint(cp), bR, uint(b.ep), uint(b.castlings))
				ml.add(mv)
				mv.packMove(uint(fr), uint(to), bP, uint(cp), bN, uint(b.ep), uint(b.castlings))
				ml.add(mv)
				mv.packMove(uint(fr), uint(to), bP, uint(cp), bB, uint(b.ep), uint(b.castlings))
				ml.add(mv)
			}
		}
		to1Step &= ^row1
		toCapL &= ^row1
		toCapR &= ^row1
	}
	// ep move
	if b.ep != 0 {
		epBB := bitBoard(1) << uint(b.ep)
		// ep left
		epToL := ((bPawns & ^fileA) >> (-SW)) & epBB
		if epToL != 0 {
			mv.packMove(uint(b.ep-SW), uint(b.ep), bP, wP, empty, uint(b.ep), uint(b.castlings))
			ml.add(mv)
		}
		epToR := ((bPawns & ^fileH) >> (-SE)) & epBB
		if epToR != 0 {
			mv.packMove(uint(b.ep-SE), uint(b.ep), bP, wP, empty, uint(b.ep), uint(b.castlings))
			ml.add(mv)
		}
	}
	// Add one step forward
	for to := to1Step.firstOne(); to != 64; to = to1Step.firstOne() {
		mv.packMove(uint(to-S), uint(to), bP, empty, empty, uint(b.ep), uint(b.castlings))
		ml.add(mv)
	}
	// Add two steps forward
	for to := to2Step.firstOne(); to != 64; to = to2Step.firstOne() {
		mv.packMove(uint(to-2*S), uint(to), bP, empty, empty, uint(b.ep), uint(b.castlings))
		ml.add(mv)
	}

	// add Captures left
	for to := toCapL.firstOne(); to != 64; to = toCapL.firstOne() {
		mv.packMove(uint(to-SW), uint(to), bP, uint(b.sq[to]), empty, uint(b.ep), uint(b.castlings))
		ml.add(mv)
	}

	// add Captures right
	for to := toCapR.firstOne(); to != 64; to = toCapR.firstOne() {
		mv.packMove(uint(to-SE), uint(to), bP, uint(b.sq[to]), empty, uint(b.ep), uint(b.castlings))
		ml.add(mv)
	}
}

// generates pseudolegal moves (without inCheck() )
func (b *boardStruct) genAllMoves(ml *moveList) {
	b.genPawnMoves(ml)
	b.genKnightMoves(ml)
	b.genBishopMoves(ml)
	b.genRookMoves(ml)
	b.genQueenMoves(ml)
	b.genKingMoves(ml)

	//
	b.filterLegals(ml)
}

// generate all legal moves
func (b *boardStruct) filterLegals(ml *moveList) {
	for ix:=len(*ml)-1;ix >=0;ix--{
		mov := (*ml)[ix]
		if b.move(mov) {
			b.unmove(mov)
		} else {
			ml.remove(ix)
		}
	}
}

func (b *boardStruct) genFrMoves(p12 int, toBB bitBoard, ml *moveList) {

}

// is sq attacked by the sd color side
func (b *boardStruct) isAttacked(sq int, sd color) bool {
	// TODO: Fixar inte alla attacker både true och false
	if pawnAtks[sd](b, sq) {
		return true
	}

	if atksKnights[sq]&b.pieceBB[Knight]&b.wbBB[sd] != 0 {
		return true
	}
	if atksKings[sq]&b.pieceBB[King]&b.wbBB[sd] != 0 {
		return true
	}
	if (mBishopTab[sq].atks(b) & (b.pieceBB[Bishop] | b.pieceBB[Queen]) & b.wbBB[sd]) != 0 {
		return true
	}
	if (mRookTab[sq].atks(b) & (b.pieceBB[Rook] | b.pieceBB[Queen]) & b.wbBB[sd]) != 0 {
		return true
	}

	return false
}

var pawnAtks = [2]func(*boardStruct, int) bool{(*boardStruct).wPawnAtks, (*boardStruct).bPawnAtks}

func (b *boardStruct) wPawnAtks(sq int) bool {
	sqBB := bitBoard(1) << uint(sq)

	wPawns := b.pieceBB[Pawn] & b.wbBB[WHITE]

	// Attacks left and right
	toCap := ((wPawns & ^fileA) << NW) & b.wbBB[BLACK]
	toCap |= ((wPawns & ^fileH) << NE) & b.wbBB[BLACK]
	if toCap&sqBB == 0 {
		return false
	}
	return true
}
func (b *boardStruct) bPawnAtks(sq int) bool {
	sqBB := bitBoard(1) << uint(sq)

	bPawns := b.pieceBB[Pawn] & b.wbBB[BLACK]

	// Attacks left and right
	toCap := ((bPawns & ^fileA) >> (-SW)) & b.wbBB[WHITE]
	toCap |= ((bPawns & ^fileH) >> (-SE)) & b.wbBB[WHITE]

	if toCap&sqBB == 0 {
		return false
	}
	return true
}

//////////////////////////////////// my own commands - NOT UCI /////////////////////////////////////

// print all legal moves
func (b *boardStruct) printAllLegals() {
	// TODO: det blir felaktigheter från inCheck() move eller unmove
	var ml moveList
	b.genAllMoves(&ml)
	fmt.Println(ml.String())
}

func (b *boardStruct) Print() {
	txtStm := "BLACK"
	if b.stm == WHITE {
		txtStm = "WHITE"
	}
	txtEp := "-"
	if b.ep != 0 {
		txtEp = sq2Fen[b.ep]
	}

	fmt.Printf("%v to move; ep: %v  castling:%v\n", txtStm, txtEp, b.castlings.String())

	fmt.Println("  +------+------+------+------+------+------+------+------+")
	for lines := 8; lines > 0; lines-- {
		fmt.Println("  |      |      |      |      |      |      |      |      |")
		fmt.Printf("%v |", lines)
		for ix := (lines - 1) * 8; ix < lines*8; ix++ {
			if b.sq[ix] == bP {
				fmt.Printf("   o  |")
			} else {
				fmt.Printf("   %v  |", int2Fen(b.sq[ix]))
			}
		}
		fmt.Println()
		fmt.Println("  |      |      |      |      |      |      |      |      |")
		fmt.Println("  +------+------+------+------+------+------+------+------+")
	}

	fmt.Printf("       A      B      C      D      E      F      G      H\n")
}

func (b *boardStruct) printAllBB() {
	txtStm := "BLACK"
	if b.stm == WHITE {
		txtStm = "WHITE"
	}
	txtEp := "-"
	if b.ep != 0 {
		txtEp = sq2Fen[b.ep]
	}
	fmt.Printf("%v to move; ep: %v   castling:%v\n", txtStm, txtEp, b.castlings.String())

	fmt.Println("white pieces")
	fmt.Println(b.wbBB[WHITE].Stringln())
	fmt.Println("black pieces")
	fmt.Println(b.wbBB[BLACK].Stringln())

	fmt.Println("wP")
	fmt.Println((b.pieceBB[Pawn] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wN")
	fmt.Println((b.pieceBB[Knight] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wB")
	fmt.Println((b.pieceBB[Bishop] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wR")
	fmt.Println((b.pieceBB[Rook] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wQ")
	fmt.Println((b.pieceBB[Queen] & b.wbBB[WHITE]).Stringln())
	fmt.Println("wK")
	fmt.Println((b.pieceBB[King] & b.wbBB[WHITE]).Stringln())

	fmt.Println("bP")
	fmt.Println((b.pieceBB[Pawn] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bN")
	fmt.Println((b.pieceBB[Knight] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bB")
	fmt.Println((b.pieceBB[Bishop] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bR")
	fmt.Println((b.pieceBB[Rook] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bQ")
	fmt.Println((b.pieceBB[Queen] & b.wbBB[BLACK]).Stringln())
	fmt.Println("bK")
	fmt.Println((b.pieceBB[King] & b.wbBB[BLACK]).Stringln())
}

// parse a FEN string and setup that position
func parseFEN(FEN string) {
	board.clear()
	fenIx := 0
	sq := 0
	for row := 7; row >= 0; row-- {
		for sq = row * 8; sq < row*8+8; {

			char := string(FEN[fenIx])
			fenIx++
			if char == "/" {
				continue
			}

			if i, err := strconv.Atoi(char); err == nil { //numeriskt
				for j := 0; j < i; j++ {
					board.setSq(empty, sq)
					sq++
				}
				continue
			}

			if strings.IndexAny(p12ToFen, char) == -1 {
				tell("info string invalid piece ", char, " try next one")
				continue
			}

			board.setSq(fen2Int(char), sq)

			sq++
		}
	}

	remaining := strings.Split(trim(FEN[fenIx:]), " ")

	// stm
	if len(remaining) > 0 {
		if remaining[0] == "w" {
			board.stm = WHITE
		} else if remaining[0] == "b" {
			board.stm = BLACK
		} else {
			r := fmt.Sprintf("%v; sq=%v;  fenIx=%v", strings.Join(remaining, " "), sq, fenIx)

			tell("info string remaining=", r, ";")
			tell("info string ", remaining[0], " invalid stm color")
			board.stm = WHITE
		}
	}

	// castling
	board.castlings = 0
	if len(remaining) > 1 {
		board.castlings = parseCastlings(remaining[1])
	}

	// ep square
	board.ep = 0
	if len(remaining) > 2 {
		if remaining[2] != "-" {
			board.ep = fenSq2Int[remaining[2]]
		}
	}

	// 50-move
	board.rule50 = 0
	if len(remaining) > 3 {
		board.rule50 = parse50(remaining[3])
	}
}

// parse 50 move rue in fenstring
func parse50(fen50 string) int {
	r50, err := strconv.Atoi(fen50)
	if err != nil || r50 < 0 {
		tell("info string 50 move rule in fenstring ", fen50, " is not a valid number >= 0 ")
		return 0
	}
	return r50
}

// parse and make the moves in position command from GUI
func parseMvs(mvstr string) {

	mvs := strings.Fields(low(mvstr))

	for _, mv := range mvs {
		mv = trim(mv)
		if len(mv) < 4 || len(mv) > 5 {
			tell("info string ", mv, " in the position command is not a correct move")
			return
		}
		// is fr square ok?
		fr, ok := fenSq2Int[mv[:2]]
		if !ok {
			tell("info string ", mv, " in the position command is not a correct fr square")
			return
		}

		p12 := board.sq[fr]
		if p12 == empty {
			tell("info string ", mv, " in the position command. fr_sq is an empty square")
			return
		}
		pCol := p12Color(p12)
		if pCol != board.stm {
			tell("info string ", mv, " in the position command. fr piece has the wrong color")
			return
		}

		// is to square ok?
		to, ok := fenSq2Int[mv[2:4]]
		if !ok {
			tell("info string ", mv, " in the position has an incorrect to square")
			return
		}

		// is the prom piece ok?
		pr := empty
		if len(mv) == 5 { //prom
			if !strings.ContainsAny(mv[4:5], "QRNBqrnb") {
				tell("info string promotion piece in ", mv, " in the position command is not correct")
				return
			}
			pr = fen2Int(mv[4:5])
			pr = pc2P12(pr, board.stm)
		}
		cp := board.sq[to]

		var intMv move // internal move format
		intMv.packMove(uint(fr), uint(to), uint(p12), uint(cp), uint(pr), uint(board.ep), uint(board.castlings))

		if !board.move(intMv){
			tell(fmt.Sprintf("tell info string %v-%v is an illegal move", sq2Fen[fr], sq2Fen[to]) )
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// fen2Int convert pieceString to p12 int
func fen2Int(c string) int {
	for p, x := range p12ToFen {
		if string(x) == c {
			return p
		}
	}
	return empty
}

// int2fen convert p12 to fenString
func int2Fen(p12 int) string {
	if p12 == empty {
		return " "
	}
	return string(p12ToFen[p12])
}

// piece returns the pc from p12
func piece(p12 int) int {
	return p12 >> 1
}

// p12Color returns the color of a p12 form
func p12Color(p12 int) color {
	return color(p12 & 0x1)
}

// pc2P12 returns p12 from pc and sd
func pc2P12(pc int, sd color) int {
	return (pc << 1) | int(sd)
}

// map fen-sq to int
var fenSq2Int = make(map[string]int)

// map int-sq to fen
var sq2Fen = make(map[int]string)

// init the square map from string to int and int to string
func initFenSq2Int() {
	fenSq2Int["a1"] = A1
	fenSq2Int["a2"] = A2
	fenSq2Int["a3"] = A3
	fenSq2Int["a4"] = A4
	fenSq2Int["a5"] = A5
	fenSq2Int["a6"] = A6
	fenSq2Int["a7"] = A7
	fenSq2Int["a8"] = A8

	fenSq2Int["b1"] = B1
	fenSq2Int["b2"] = B2
	fenSq2Int["b3"] = B3
	fenSq2Int["b4"] = B4
	fenSq2Int["b5"] = B5
	fenSq2Int["b6"] = B6
	fenSq2Int["b7"] = B7
	fenSq2Int["b8"] = B8

	fenSq2Int["c1"] = C1
	fenSq2Int["c2"] = C2
	fenSq2Int["c3"] = C3
	fenSq2Int["c4"] = C4
	fenSq2Int["c5"] = C5
	fenSq2Int["c6"] = C6
	fenSq2Int["c7"] = C7
	fenSq2Int["c8"] = C8

	fenSq2Int["d1"] = D1
	fenSq2Int["d2"] = D2
	fenSq2Int["d3"] = D3
	fenSq2Int["d4"] = D4
	fenSq2Int["d5"] = D5
	fenSq2Int["d6"] = D6
	fenSq2Int["d7"] = D7
	fenSq2Int["d8"] = D8

	fenSq2Int["e1"] = E1
	fenSq2Int["e2"] = E2
	fenSq2Int["e3"] = E3
	fenSq2Int["e4"] = E4
	fenSq2Int["e5"] = E5
	fenSq2Int["e6"] = E6
	fenSq2Int["e7"] = E7
	fenSq2Int["e8"] = E8

	fenSq2Int["f1"] = F1
	fenSq2Int["f2"] = F2
	fenSq2Int["f3"] = F3
	fenSq2Int["f4"] = F4
	fenSq2Int["f5"] = F5
	fenSq2Int["f6"] = F6
	fenSq2Int["f7"] = F7
	fenSq2Int["f8"] = F8

	fenSq2Int["g1"] = G1
	fenSq2Int["g2"] = G2
	fenSq2Int["g3"] = G3
	fenSq2Int["g4"] = G4
	fenSq2Int["g5"] = G5
	fenSq2Int["g6"] = G6
	fenSq2Int["g7"] = G7
	fenSq2Int["g8"] = G8

	fenSq2Int["h1"] = H1
	fenSq2Int["h2"] = H2
	fenSq2Int["h3"] = H3
	fenSq2Int["h4"] = H4
	fenSq2Int["h5"] = H5
	fenSq2Int["h6"] = H6
	fenSq2Int["h7"] = H7
	fenSq2Int["h8"] = H8

	// -------------- sq2Fen
	sq2Fen[A1] = "a1"
	sq2Fen[A2] = "a2"
	sq2Fen[A3] = "a3"
	sq2Fen[A4] = "a4"
	sq2Fen[A5] = "a5"
	sq2Fen[A6] = "a6"
	sq2Fen[A7] = "a7"
	sq2Fen[A8] = "a8"

	sq2Fen[B1] = "b1"
	sq2Fen[B2] = "b2"
	sq2Fen[B3] = "b3"
	sq2Fen[B4] = "b4"
	sq2Fen[B5] = "b5"
	sq2Fen[B6] = "b6"
	sq2Fen[B7] = "b7"
	sq2Fen[B8] = "b8"

	sq2Fen[C1] = "c1"
	sq2Fen[C2] = "c2"
	sq2Fen[C3] = "c3"
	sq2Fen[C4] = "c4"
	sq2Fen[C5] = "c5"
	sq2Fen[C6] = "c6"
	sq2Fen[C7] = "c7"
	sq2Fen[C8] = "c8"

	sq2Fen[D1] = "d1"
	sq2Fen[D2] = "d2"
	sq2Fen[D3] = "d3"
	sq2Fen[D4] = "d4"
	sq2Fen[D5] = "d5"
	sq2Fen[D6] = "d6"
	sq2Fen[D7] = "d7"
	sq2Fen[D8] = "d8"

	sq2Fen[E1] = "e1"
	sq2Fen[E2] = "e2"
	sq2Fen[E3] = "e3"
	sq2Fen[E4] = "e4"
	sq2Fen[E5] = "e5"
	sq2Fen[E6] = "e6"
	sq2Fen[E7] = "e7"
	sq2Fen[E8] = "e8"

	sq2Fen[F1] = "f1"
	sq2Fen[F2] = "f2"
	sq2Fen[F3] = "f3"
	sq2Fen[F4] = "f4"
	sq2Fen[F5] = "f5"
	sq2Fen[F6] = "f6"
	sq2Fen[F7] = "f7"
	sq2Fen[F8] = "f8"

	sq2Fen[G1] = "g1"
	sq2Fen[G2] = "g2"
	sq2Fen[G3] = "g3"
	sq2Fen[G4] = "g4"
	sq2Fen[G5] = "g5"
	sq2Fen[G6] = "g6"
	sq2Fen[G7] = "g7"
	sq2Fen[G8] = "g8"

	sq2Fen[H1] = "h1"
	sq2Fen[H2] = "h2"
	sq2Fen[H3] = "h3"
	sq2Fen[H4] = "h4"
	sq2Fen[H5] = "h5"
	sq2Fen[H6] = "h6"
	sq2Fen[H7] = "h7"
	sq2Fen[H8] = "h8"
}

// 6 piece types - no color (P)
const (
	Pawn int = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

// 12 pieces with color (P12)
const (
	wP = iota
	bP
	wN
	bN
	wB
	bB
	wR
	bR
	wQ
	bQ
	wK
	bK
	empty = 15
)

// piece char definitions
const (
	pc2Char  = "PNBRQK?"
	p12ToFen = "PpNnBbRrQqKk"
)

// square names
const (
	A1 = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1

	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2

	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3

	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4

	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5

	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6

	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7

	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

//////////////////////////////// TODO: remove this after benchmarking ////////////////////////////////////////
func (b *boardStruct) genSimpleRookMoves(ml *moveList, sd color) {
	allRBB := b.pieceBB[Rook] & b.wbBB[sd]
	p12 := uint(pc2P12(Rook, color(sd)))
	ep := uint(b.ep)
	castlings := uint(b.castlings)
	var mv move
	for fr := allRBB.firstOne(); fr != 64; fr = allRBB.firstOne() {
		rk := fr / 8
		fl := fr % 8
		//N
		for r := rk + 1; r < 8; r++ {
			to := uint(r*8 + fl)
			cp := uint(b.sq[to])
			if cp != empty && p12Color(int(cp)) == sd {
				break
			}
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castlings)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
		//S
		for r := rk - 1; r >= 0; r-- {
			to := uint(r*8 + fl)
			cp := uint(b.sq[to])
			if cp != empty && p12Color(int(cp)) == sd {
				break
			}
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castlings)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
		//E
		for f := fl + 1; f < 8; f++ {
			to := uint(rk*8 + f)
			cp := uint(b.sq[to])
			if cp != empty && p12Color(int(cp)) == sd {
				break
			}
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castlings)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
		//W
		for f := fl - 1; f >= 0; f-- {
			to := uint(rk*8 + f)
			cp := uint(b.sq[to])
			if cp != empty && p12Color(int(cp)) == sd {
				break
			}
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castlings)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
	}
}
