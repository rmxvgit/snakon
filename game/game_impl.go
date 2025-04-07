package game

import (
	"fmt"
	"os"
	"snakon/game/input"
	"snakon/game/network"
	"snakon/game/render"
	"snakon/utils"
	"time"
)

func Run() {
	SetupGame().start()
}

func SetupGame() (game *Game) {
	var err error
	addr := os.Args[2]
	game = &Game{}
	game.input = input.SetupGameInput()
	game.net_handle, err = network.SetupNetManager(addr)
	utils.PanicOnError(err)
	game.renderer = render.SetupRender()
	return
}

func (game *Game) start() {
	for {
		game.net_handle.SendMyPos(2, 3, 3)
		game.render()
		time.Sleep(1000)
	}
}

func (game *Game) render() {
	game.renderer.Scr_buffer.Write([]byte("\033[H\033[2J"))
	for _, value := range game.net_handle.Game_state.PlayersState {
		str := fmt.Sprint("id:", value.ID, "posx:", value.Pos.X, "posy:", value.Pos.Y)
		game.renderer.Scr_buffer.Write([]byte(str))
	}
	game.renderer.Scr_buffer.Flush()
}

func (game *Game) process() {

}
