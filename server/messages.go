package server

type MessageFlag byte

const (
	NEW_PLAYER_MESSAGE MessageFlag = iota
	PLAYER_POS_MESSAGE
	MANY_PLAYER_POS_MSG
	GAME_STATE_MESSAGE
)

type NewPlayerMessage struct {
	PlayerID int32
}

type PlayerPositionMessage struct {
	PlayerID int32
	Pos      PlayerPos
}

type ManyPlayerPositionMessage struct {
	Positions []PlayerPositionMessage
}
