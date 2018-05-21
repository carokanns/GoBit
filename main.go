package main

import (
	"strings"
)

func main() {
	tell("info string Starting GoBit")

	uci(input())

	tell("info string quits GOBIT")
}

func init() {
	initFen2Sq()
	initMagic()
	initKeys()
	initAtksKings()
	initAtksKnights()
	initCastlings()
	pcSqInit()
	board.newGame()
	handleSetOption(strings.Split("setoption name hash value 32", " "))

}
