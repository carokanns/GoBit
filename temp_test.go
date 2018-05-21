package main

import (
	"testing"
)

////////////////////////////// BENCHMARKS ///////////////////////////////
func Benchmark_genAllCapturesx(b *testing.B) {
	var ml = make(moveList, 0, 60)

	handleNewgame()
	handlePosition("position startpos moves d2d4 d7d5 c1f4 g8f6 e2e3 e7e6 b1d2 c7c5 c2c3 b8c6 g1f3 f8e7 f1d3 c8d7")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		board.genAllCaptures(&ml)
	}
}

func Benchmark_genAllCapturesy(b *testing.B) {
	var ml = make(moveList, 0, 60)

	handleNewgame()
	handlePosition("position startpos moves d2d4 d7d5 c1f4 g8f6 e2e3 e7e6 b1d2 c7c5 c2c3 b8c6 g1f3 f8e7 f1d3 c8d7")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		board.genAllCapturesy(&ml)
	}

}
