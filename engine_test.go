package main

import "testing"

func Test_evaluate(t *testing.T) {
	tests := []struct {
		name string
		pos string
		want int
	}{
		{"","position startpos",0},
		{"","position startpos moves e2e4 d7d5 e4d5",100},
	}
	pSqInit()
	for _, tt := range tests {
		handlePosition(tt.pos)
		t.Run(tt.name, func(t *testing.T) {
			if got := board.evaluate(); got != tt.want {
				t.Errorf("%v: evaluate() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
