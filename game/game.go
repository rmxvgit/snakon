package game

import (
	"snakon/game/input"
	"snakon/game/network"
	"snakon/game/render"
)

type Game struct {
	input      *input.GameInput
	net_handle *network.GameNetManager
	renderer   *render.Renderer
}
