package main

import (
	"fmt"
	"os"
	"snakon/game"
	"snakon/sender"
	"snakon/server"
)

// to run the program in serever mode:
// ./<name_of_the_executable> serv listening_port
func main() {
	run_mode := os.Args[1] // serv / cli

	if run_mode == "serv" {
		server.Run()
	} else if run_mode == "cli" {
		game.Run()
	} else if run_mode == "send" {
		sender.Run()
	} else {
		err := fmt.Errorf("invalid run mode: %s\n", run_mode)
		panic(err)
	}
}
