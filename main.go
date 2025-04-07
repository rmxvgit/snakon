package main

import (
	"fmt"
	"os"
	"snakon/game"
	"snakon/server"
)

// to run the program in serever mode:
// ./<name_of_the_executable> serv listening_port
func main() {
	run_mode := os.Args[1]

	if run_mode == "serv" {
		server.Run()
	} else if run_mode == "cli" {
		game.Run()
	} else {
		err := fmt.Errorf("invalid run mode: %s\n", run_mode)
		panic(err)
	}
}
