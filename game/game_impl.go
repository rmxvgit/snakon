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
		game.render()
		time.Sleep(10000000)
	}
}

func (game *Game) process() {
	game.movePlayer()
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

func (game *Game) movePlayer() {
	switch game.player_move_dir {
	case gametypes.DirectionUp:
		game.player_pos.Y--
	case gametypes.DirectionDown:
		game.player_pos.Y++
	case gametypes.DirectionLeft:
		game.player_pos.X--
	case gametypes.DirectionRight:
		game.player_pos.X++
	}
}

func (game *Game) render() {
	game.renderer.CleanBuffer()

	// Render player
	game.renderer.WriteChar(int(game.player_pos.X), int(game.player_pos.Y), 'P')

	// Render other players
	for _, pos := range game.other_players_pos {
		game.renderer.WriteChar(int(pos.X), int(pos.Y), 'O')
	}

	game.renderer.Flush()
}
