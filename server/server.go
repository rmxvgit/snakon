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
	addr    *net.UDPAddr
	conn    *net.UDPConn
	clients map[string]*ClientInfo
	state   *ServerState
}

type ClientInfo struct {
	mutex         sync.Mutex
	n_msg_sent    uint64
	n_msg_recv    uint64
	last_msg_recv uint64
}

type ServerState struct {
	players map[int32]*PlayerServerState
}

type PlayerServerState struct {
	LastMessage uint64
	Mutex       sync.Mutex
	Pos         gametypes.Position
}
