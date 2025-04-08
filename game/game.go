package game

import (
	"snakon/game/input"
	"snakon/game/network"
	"snakon/game/render"
)

type Game struct {
	input     *input.GameInput
	internet  *network.GameNetwork
	renderer  *render.Renderer
	player_id int32
}
