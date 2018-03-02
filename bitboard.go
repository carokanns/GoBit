package main

import "math/bits"

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
