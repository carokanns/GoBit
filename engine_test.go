package main

import (
	"testing"
)

func Test_genAndSort(t *testing.T) {
	tests := []struct {
		name string
		pos string
		want1Mv string
		want1Ev int	
	}{
		{"","position startpos ","e2e4",24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ml moveList 
			genAndSort(&board, &ml)
			if tt.want1Mv != trim(ml[0].String()) {
			t.Errorf("%v: %#v should be best move. Got %#v", tt.name, tt.want1Mv, trim(ml[0].String()))
		}
		if tt.want1Ev != ml[0].eval() {
			t.Errorf("%v: %v should be best score. Got %v", tt.name, tt.want1Ev, ml[0].eval())
		}

		})
	}
}
