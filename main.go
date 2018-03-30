package main

func main() {
	tell("info string Starting GoBit")

	uci(input())

	tell("info string quits GOBIT")
}

func init(){
	initFenSq2Int()
	initMagic()
	initAtksKings()
	initAtksKnights()
}