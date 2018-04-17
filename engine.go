package main

import (
	"fmt"
	"math"
)

type searchLimits struct {
	depth    int
	nodes    uint64
	moveTime int // in milliseconds
	infinite bool
	//////////////// current //////////
	stop bool
}

var limits searchLimits

func (s *searchLimits) init() {
	s.depth = 9999
	s.nodes = math.MaxUint64
	s.moveTime = 99999999999
	s.infinite = false
}

func (s *searchLimits) setStop(st bool) {
	s.stop = st
}
func (s *searchLimits) setDepth(d int) {
	s.depth = d
}
func (s *searchLimits) setMoveTime(m int) {
	s.moveTime = m
}
func (s *searchLimits) setInfinite(b bool) {
	s.infinite = b
}

func engine() (toEngine chan bool, frEngine chan string) {
	fmt.Println("info string Hello from engine")
	frEngine = make(chan string)
	toEngine = make(chan bool)
	go rootx(toEngine, frEngine)

	return
}

func root(toEngine chan bool, frEngine chan string) {
	for _ = range toEngine {
		tell("info string engine got go!")
		// genAllMoves
		// evaluate and sort
		// for each move{
		// 		score := search()
		// 		store score in move
		//}
		// reply to GUI with the best move
	}
}

func search(b *boardStruct) int {

	return b.evaluate()
}

