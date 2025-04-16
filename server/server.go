package server

import (
	"net"
	gametypes "snakon/gameTypes"
	"sync"
)

type MsgErrCode int

const PACKET_SIZE int = 1024

const (
	MsgOk MsgErrCode = iota
	MessageDrop
	NullPlayerReference
	PlayerRedeclaration
)

type Server struct {
	addr        *net.UDPAddr
	conn        *net.UDPConn
	state_mutex sync.Mutex
	state       *ServerState
}

type ServerState struct {
	players map[int32]*PlayerServerState
}

type PlayerServerState struct {
	LastMessage uint64
	Mutex       sync.Mutex
	Pos         gametypes.Position
}
