package main

import (
	"fmt"
	"math/bits"
	"strings"
)

type bitBoard uint64

func (b bitBoard) count() int {
	return bits.OnesCount64(uint64(b))
}

func (b *bitBoard) set(pos uint) {
	*b |= bitBoard(uint64(1) << pos)
}

func (b bitBoard) test(pos uint) bool {
	return (b & bitBoard(uint64(1)<<pos)) != 0
}

func (b *bitBoard) clr(pos uint) {
	*b &= bitBoard(^(uint64(1) << pos))
}

func (b *bitBoard) firstOne() int {
	bit := bits.TrailingZeros64(uint64(*b))
	if bit == 64 {
		return 64
	}
	*b = (*b >> uint(bit+1)) << uint(bit+1)
	return bit
}

//TODO: lastOne() here

// returns the full bitstring (with leading zeroes) of the bitBoard
func (b bitBoard) String() string {
	zeroes := ""
	for ix := 0; ix < 64; ix++ {
		zeroes = zeroes + "0"
	}

	bits := zeroes + fmt.Sprintf("%b", b)
	return bits[len(bits)-64:]
}

// returns the bitboard string 8x8
func (b bitBoard) Stringln() string {
	s := b.String()
	row := [8]string{}
	row[0] = s[0:8]
	row[1] = s[8:16]
	row[2] = s[16:24]
	row[3] = s[24:32]
	row[4] = s[32:40]
	row[5] = s[40:48]
	row[6] = s[48:56]
	row[7] = s[56:]
	for ix, r := range row {
		row[ix] = fmt.Sprintf("%v%v%v%v%v%v%v%v\n", r[7:8], r[6:7], r[5:6], r[4:5], r[3:4], r[2:3], r[1:2], r[0:1])
	}

	s = strings.Join(row[:], "")
	s = strings.Replace(s, "1", "1 ", -1)
	s = strings.Replace(s, "0", "0 ", -1)
	return s
}
