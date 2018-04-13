package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var tell = mainTell
var trim = strings.TrimSpace
var low = strings.ToLower

var saveBm = ""

func uci(input chan string) {
	fmt.Println("info string Hello from uci")

	toEng, frEng := engine()
	bInfinite := false
	var cmd string
	var bm string
	quit := false
	for !quit {
		select {
		case cmd = <-input:
			//			tell("info string uci got ", cmd, "\n")
		case bm = <-frEng:
			handleBm(bm, bInfinite)
			continue
		}
		words := strings.Split(cmd, " ")
		words[0] = trim(low(words[0]))

		switch words[0] {
		case "uci":
			handleUci()
		case "setoption":
			handleSetOption(words)
		case "isready":
			handleIsReady()
		case "ucinewgame":
			handleNewgame()
		case "position":
			handlePosition(cmd)
		case "debug":
			handleDebug(words)
		case "register":
			handleRegister(words)
		case "go":
			handleGo(words)
		case "ponderhit":
			handlePonderhit()
		case "stop":
			handleStop(toEng, &bInfinite)
		case "quit", "q":
			handleQuit(toEng)
			quit = true
			continue
		case "pb":
			board.Print()
		case "pbb":
			board.printAllBB()
		case "pm":
			board.printAllLegals()
		default:
			tell("info string unknown cmd ", cmd)
		}
	}

	tell("info string leaving uci()")
}

func handleUci() {
	tell("id name GoBit")
	tell("id author Carokanns")

	tell("option name Hash type spin default 128 min 16 max 1024")
	tell("option name Threads type spin default 1 min 1 max 16")
	tell("uciok")
}

func handleIsReady() {
	tell("readyok")
}

func handleNewgame() {
	board.newGame()
}

func handlePosition(cmd string) {
	// position [fen <fenstring> | startpos ]  moves <move1> .... <movei>
	board.newGame()
	cmd = trim(strings.TrimPrefix(cmd, "position"))
	parts := strings.Split(cmd, "moves")
	if len(cmd) == 0 || len(parts) > 2 {
		err := fmt.Errorf("%v wrong length=%v", parts, len(parts))
		tell("info string Error", fmt.Sprint(err))
		return
	}

	alt := strings.Split(parts[0], " ")
	alt[0] = trim(alt[0])
	if alt[0] == "startpos" {
		parts[0] = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	} else if alt[0] == "fen" {
		parts[0] = trim(strings.TrimPrefix(parts[0], "fen"))
	} else {
		err := fmt.Errorf("%#v must be %#v or %#v", alt[0], "fen", "startpos")
		tell("info string Error", err.Error())
		return
	}
	// Now parts[0] is the fen-string only

	// start the parsing
	//fmt.Printf("info string parse %#v\n", parts[0])
	parseFEN(parts[0])

	if len(parts) == 2 {
		parts[1] = low(trim(parts[1]))
		//fmt.Printf("info string parse %#v\n", parts[1])
		parseMvs(parts[1])
	}
}

func handleStop(toEng chan string, bInfinite *bool) {
	if *bInfinite {
		if saveBm != "" {
			tell(saveBm)
			saveBm = ""
		}

		toEng <- "stop"
		*bInfinite = false
	}
	tell("info string stop not implemented")
}

// handleQuit not really necessary
func handleQuit(toEng chan string) {
	toEng <- "stop"
}

func handleBm(bm string, bInfinite bool) {
	if bInfinite {
		saveBm = bm
		return
	}
	tell(bm)
}

// not implemented uci commands
func handleSetOption(words []string) {
	// setoption name Hash value 256
	tell("info string setoption not implemented")
}

// go  searchmoves <move1-moveii>/ponder/wtime <ms>/ btime <ms>/winc <ms>/binc <ms>/movestogo <x>/
//     depth <x>/nodes <x>/movetime <ms>/mate <x>/infinite
func handleGo(words []string) {
	// TODO: Start with moveTime and infinite
	if len(words) > 1 {
		words[1] = trim(low(words[1]))
		switch words[1] {
		case "searchmoves":
			tell("info string go searchmoves not implemented")
		case "ponder":
			tell("info string go ponder not implemented")
		case "wtime":
			tell("info string go wtime not implemented")
		case "btime":
			tell("info string go btime not implemented")
		case "winc":
			tell("info string go winc not implemented")
		case "binc":
			tell("info string go binc not implemented")
		case "movestogo":
			tell("info string go movestogo not implemented")
		case "depth":
			tell("info string go depth not implemented")
		case "nodes":
			tell("info string go nodes not implemented")
		case "movetime":
			tell("info string go movetime not implemented")
		case "mate": // mate <x>  mate in x moves
			tell("info string go mate not implemented")
		case "infinite":
			tell("info string go infinite not implemented")
		case "register":
			tell("info string go register not implemented")
		default:
			tell("info string go ", words[1], " not implemented")
		}
	} else {
		tell("info string go not implemented")
	}
}

func handlePonderhit() {
	tell("info string ponderhit not implemented")
}

func handleDebug(words []string) {
	// debug [ on | off ]
	tell("info string debug not implemented")
}

func handleRegister(words []string) {
	// register later/name <x>/code <y>
	tell("info string register not implemented")
}

//------------------------------------------------------
func mainTell(text ...string) {
	toGUI := ""
	for _, t := range text {
		toGUI += t
	}
	fmt.Println(toGUI)
}

func input() chan string {
	line := make(chan string)
	var reader *bufio.Reader
	reader = bufio.NewReader(os.Stdin)
	go func() {
		for {
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if err != io.EOF && len(text) > 0 {
				line <- text
			}
		}
	}()
	return line
}
