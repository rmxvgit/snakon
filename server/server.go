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
	conn    *net.UDPConn           // the connection used for sending and receiving messages
	clients map[string]*ClientInfo // a map containing client information
	state   *ServerState           // the game state stored in the server
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
	Mutex sync.Mutex
	Pos   gametypes.Position
}
