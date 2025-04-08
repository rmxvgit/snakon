package game

import (
	"fmt"
	"os"
	"snakon/game/input"
	"snakon/game/network"
	"snakon/game/render"
	"snakon/utils"
	"strconv"
	"time"
)

func Run() {
	game, err := SetupGame()
	utils.PanicOnError(err)
	game.start()
}

func SetupGame() (game *Game, err error) {
	game = &Game{}

	id, err := strconv.Atoi(os.Args[4])
	utils.PanicOnError(err)
	game.player_id = int32(id)

	game.input = input.SetupGameInput()

	client_addr := os.Args[2]
	server_addr := os.Args[3]
	game.internet, err = network.NewGameNetwork(client_addr, server_addr)
	if err != nil {
		return nil, err
	}
	go game.internet.Listen()

	game.renderer = render.SetupRender(10, 10)
	return
}

func (game *Game) start() {
	game.internet.NewPlayer(uint32(game.player_id))
	for {
		game.process()
		game.render()
		time.Sleep(10000000)
	}
}

func (game *Game) render() {
	game.renderer.CleanBuffer()
	for _, player_state := range game.internet.GetServerGameState().PlayersState {
		x := player_state.Pos.X
		y := player_state.Pos.Y
		game.renderer.WriteChar(int(x), int(y), '&')
	}
	game.renderer.Flush()
}

func (game *Game) diagnosticRender() {
	state, err := game.internet.GetPlayerState(4)
	if err != nil {
		return
	}
	fmt.Printf("id: %d x: %d, y: %d\n", state.ID, state.Pos.X, state.Pos.Y)
}

func (game *Game) process() {
	myState, err := game.internet.GetPlayerState(game.player_id)
	if err != nil {
		return
	}

	key := game.input.Keyboard.ConsumeLastKey()
	switch key {
	case 'w':
		myState.Pos.Y--
	case 'a':
		myState.Pos.X--
	case 's':
		myState.Pos.Y++
	case 'd':
		myState.Pos.X++
	}
	game.internet.SendEntityPos(myState.ID, myState.Pos)
}
