package game

import (
	"snakon/game/input"
	"snakon/game/network"
	"snakon/game/render"
	gametypes "snakon/gameTypes"
)

type Game struct {
	input             *input.GameInput
	internet          *network.GameNetwork
	renderer          *render.Renderer
	player_id         int32
	player_pos        gametypes.Position
	player_move_dir   gametypes.Direction
	other_players_pos map[int32]gametypes.Position
}
