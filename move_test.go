package main

import "testing"

func Test_moveList_add(t *testing.T) {
	tests := []struct {
		name string
		ml   moveList
		mv move
	}{
		{"", moveList{},0},
		{"", moveList{},1},
		{"", moveList{},2},
		{"", moveList{},3},
	}

	ml = moveList{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ml = tt.ml
			tt.ml.add(tt.mv)
			ix:=len(tt.ml.mv)-1
			if ix<0{
				t.Fatalf("wrong len %v: %v %v %v",ix,tt.ml,tt.ml.mv,tt.ml.mv[0])
			}
			testMv := tt.ml.mv[ix]
			if  testMv != tt.mv{
				t.Errorf("tt.ml.add() = %v. want %v",testMv,tt.mv)
			}
		})
	}
}
