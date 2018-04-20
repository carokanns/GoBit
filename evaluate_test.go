package main

import "testing"

func Test_evaluate(t *testing.T) {
	tests := []struct {
		name string
		pos  string
		want int
	}{
		{"", "position startpos", 0},
		{"e4", "position startpos moves e2e4", 24},              
		{"Nf3", "position startpos moves g1f3", 22},             
		{"d4", "position startpos moves d2d4", 20},               
		{"c4", "position startpos moves c2c4", 11},               
		{"Nc3", "position startpos moves b1c3", 21},              
		{"e4+d4", "position startpos moves e2e4 e7e5 d2d4",20},   
		{"e4+Nf3", "position startpos moves e2e4 e7e5 g1f3", 22}, 
		{"e4+Be2", "position startpos moves e2e4 e7e5 f1e2", 15},
		{"e4+Bd3", "position startpos moves e2e4 e7e5 f1d3", 17},
		{"e4+Bc4", "position startpos moves e2e4 e7e5 f1c4", 18},
		{"e4+Bb5", "position startpos moves e2e4 e7e5 f1b5", 17},
		{"e4+Ba6", "position startpos moves e2e4 e7e5 f1a6", 8},
	}

	for _, tt := range tests {
		handlePosition(tt.pos)
		t.Run(tt.name, func(t *testing.T) {
			if got := evaluate(&board); got != tt.want {
				t.Errorf("%v: evaluate() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
