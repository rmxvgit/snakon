package gametypes

type Position struct {
	X, Y int32
}

type PlayerPos struct {
	ID  int32
	Pos Position
}

type PlayerState struct {
	ID  int32
	Pos Position
}
