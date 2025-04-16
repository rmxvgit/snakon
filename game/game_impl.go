package game

import (
	"os"
	"snakon/game/input"
	"snakon/game/network"
	"snakon/game/render"
	"snakon/utils"
	"strconv"
	"time"
)

func Run() {
	game := SetupGame()
	game.start()
}

func SetupGame() (game *Game) {
	game = &Game{}

	player_id, err := strconv.Atoi(os.Args[4])
	utils.PanicOnError(err)
	game.player_id = int32(player_id)

	client_addr := os.Args[2]
	server_addr := os.Args[3]

	game.internet, err = network.NewGameNetwork(client_addr, server_addr)
	utils.PanicOnError(err)

	go game.internet.Listen()

	game.input = input.SetupGameInput()

	game.renderer = render.SetupRender(10, 10)
	return
}

func (game *Game) start() {
	for {
		game.process()
		time.Sleep(10000000)
	}
}

func (game *Game) process() {
}
