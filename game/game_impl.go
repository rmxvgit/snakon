package game

import (
	"os"
	"snakon/game/input"
	"snakon/game/network"
	"snakon/game/render"
	gametypes "snakon/gameTypes"
	"snakon/utils"
	"strconv"
	"time"
)

func Run() {
	player_id, err := strconv.Atoi(os.Args[4])
	utils.PanicOnError(err)

	client_addr := os.Args[2]
	server_addr := os.Args[3]

	game := NewGame(int32(player_id), client_addr, server_addr)

	game.start()
}

func NewGame(player_id int32, client_addr string, server_addr string) (game *Game) {
	game = &Game{}

	game.internet = network.NewGameNetwork(client_addr, server_addr)
	go game.internet.Listen()

	game.input = input.SetupGameInput()

	game.renderer = render.SetupRender(10, 10)

	game.player_id = player_id
	err := game.internet.SendNewPlayerNotification(player_id)
	utils.PanicOnError(err)

	return
}

func (game *Game) start() {
	for {
		game.gather_input()
		game.process()
		time.Sleep(10000000)
	}
}

func (game *Game) process() {

	err := game.internet.SendPlayerPosition(game.player_id, game.player_pos)
	utils.PanicOnError(err)
}

func (game *Game) gather_input() {
	key := game.input.Keyboard.ConsumeLastKey()

	switch key {
	case 'w':
		game.player_move_dir = gametypes.DirectionUp
	case 's':
		game.player_move_dir = gametypes.DirectionDown
	case 'a':
		game.player_move_dir = gametypes.DirectionLeft
	case 'd':
		game.player_move_dir = gametypes.DirectionRight
	default:
		game.player_move_dir = gametypes.DirectionNone
	}

	game.other_players_pos = game.internet.GetOtherPlayersPositions()
}
