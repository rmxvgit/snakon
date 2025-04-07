package server

type GameState struct {
	PlayersState map[int32]*PlayerState
}

type PlayerState struct {
	ID  int32
	Pos PlayerPos
}

type PlayerPos struct {
	X int32
	Y int32
}
