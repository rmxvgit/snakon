package messages

import gametypes "snakon/gameTypes"

type MessageFlag byte

const PACKET_SIZE int = 1024

const (
	NEW_PLAYER_MESSAGE MessageFlag = iota
	PLAYER_POS_MESSAGE
	MANY_PLAYER_POS_MSG
	GAME_STATE_MESSAGE
)

type NewPlayerDto struct {
	PlayerID int32
}

type PlayerPositionDto struct {
	PlayerID int32
	Pos      gametypes.Position
}

type ManyPlayerPositionDto struct {
	Positions []PlayerPositionDto
}
