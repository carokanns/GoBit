package main

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
}
