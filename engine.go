package main

func engine() (toEngine chan string, frEngine chan string) {
	tell("info string Hello from engine")
	frEngine = make(chan string)
	toEngine = make(chan string)
	go func() {
		for cmd := range toEngine {
			tell("info string engine got ", cmd)
			switch cmd {
			case "stop":
			case "quit":

			}
		}
	}()

	return
}
