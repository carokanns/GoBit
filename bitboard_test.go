package main

import (
	"testing"
)

func Test_bitBoard_some(t *testing.T) {

	tests := []struct {
		name string
		b    bitBoard
		pos  uint
	}{
		{"", 0xF, 63},
		{"", 0x0, 0},
		{"", 0x1, 0},
		{"", 0xFFFF, 63},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.set(tt.pos)
			if !tt.b.test(tt.pos) {
				t.Fatalf("set(%v) gives %v in %b. Want %v", tt.pos, false, tt.b, true)
			}

			tt.b.clr(tt.pos)
			if tt.b.test(tt.pos) {
				t.Errorf("clr(%v) gives %v in %b. Want %v", tt.pos, true, tt.b, false)
			}
		})
	}
}

func Test_bitBoard_firstOne(t *testing.T) {
	tests := []struct {
		name string
		b    bitBoard
		want int
	}{
		{"",0x0,64},
		{"",0x1,0},
		{"",0xFFFFFFFFFFFFFFFF,0},
		{"",0xFFFFFFFFFFFFFF00,8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x:=tt.b
			if got := tt.b.firstOne(); got != tt.want {
				t.Errorf("bitBoard.firstOne(%x) = %v, want %v", x,got, tt.want)
			}
		})
	}
}
