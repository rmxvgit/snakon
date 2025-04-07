package server

func NewGameState() *GameState {
	return &GameState{
		PlayersState: make(map[int32]*PlayerState),
	}
}

func (game_state *GameState) ManyPlayerPositionMessage() ManyPlayerPositionMessage {
	players_pos := make([]PlayerPositionMessage, len(game_state.PlayersState))
	pos_index_to_fill := 0

	for _, player_state := range game_state.PlayersState {
		players_pos[pos_index_to_fill] = player_state.PlayerPositionMessage()
		pos_index_to_fill++
	}

	return ManyPlayerPositionMessage{Positions: players_pos}
}

func NewPlayerState(player_id int32) *PlayerState {
	player_state := PlayerState{
		ID: player_id,
	}
	return &player_state
}

func (player_state *PlayerState) PlayerPositionMessage() PlayerPositionMessage {
	return PlayerPositionMessage{
		PlayerID: player_state.ID,
		Pos: PlayerPos{
			X: player_state.Pos.X,
			Y: player_state.Pos.Y,
		},
	}
}
