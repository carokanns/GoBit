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
)

type boardStruct struct {
	sq      [64]int
	wbBB    [2]bitBoard
	pieceBB [nP]bitBoard
	King    [2]int
	ep      int
	castlings
	stm    color
	count  [12]int
	rule50 int //set to 0 if pawn or capt nmove otherwise increment
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

// make a pseudomove
func (b *boardStruct) move(fr, to, pr int) bool {

	newEp := 0
	// we assume that the move is legally correct
	p12 := b.sq[fr]
	switch {
	case p12 == wK:
		b.castlings.off(shortW | longW)
		if abs(to-fr) == 2 {
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
		if abs(to-fr) == 2 {
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

	// TODO isInCheck() needs to be made (when move generation is finished)
	/*
		if b.isInCheck(b.stm) {
			b.stm = b.stm ^ 0x1
			return false
		}
	*/
	b.stm = b.stm ^ 0x1
	return true
}

func (b *boardStruct) setSq(p12, s int) {
	if b.sq[s] != empty {
		b.count[b.sq[s]]--
	}
	b.sq[s] = p12

	if p12 == empty {
		b.wbBB[WHITE].clr(uint(s))
		b.wbBB[BLACK].clr(uint(s))
		for p := 0; p < nP; p++ {
			b.pieceBB[p].clr(uint(s))
		}
		return
	}

	b.count[p12]++

	p := piece(p12)
	sd := p12Color(p12)

	if p == King {
		b.King[sd] = s
	}

	b.wbBB[sd].set(uint(s))
	b.pieceBB[p].set(uint(s))
}

func (b *boardStruct) newGame() {
	b.stm = WHITE
	b.clear()
	parseFEN(startpos)
}
func (b *boardStruct) genRookMoves(ml *moveList, sd color) {
	// TODO: explain magic numbers/bitboards
	// TODO: create tables and methods to use magic numbers
	// TODO: generate rook moves with magic bitBoards - Benchmark

	allRBB := b.pieceBB[Rook] & b.wbBB[sd]
	p12 := uint(pc2P12(Rook, color(sd)))
	ep := uint(b.ep)
	castl := uint(b.castlings)
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
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castl)
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
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castl)
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
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castl)
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
			mv.packMove(uint(fr), to, p12, cp, empty, ep, castl)
			ml.add(mv)
			if cp != empty {
				break
			}
		}
	}
}

func (b *boardStruct) genFrMoves(p12 int, frBB bitBoard, ml *moveList) {

}

//////////////////////////////////// my own commands - NOT UCI /////////////////////////////////////

func (b *boardStruct) printAllMvs() {

	fmt.Println("magic")
	// var ml moveList
	// board.genRookMovesM(&ml)
	// fmt.Println(ml.String())

	var ml1 moveList
	fmt.Println("simple")
	board.genRookMoves(&ml1, b.stm)
	fmt.Println(ml1.String())

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
			if !strings.ContainsAny(mv[4:5], "qrbn") {
				tell("info string promotion piece in ", mv, " in the position command is not correct")
				return
			}
			pr = fen2Int(mv[4:5])
			pr = pc2P12(pr, board.stm)
		}
		board.move(fr, to, pr)
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
