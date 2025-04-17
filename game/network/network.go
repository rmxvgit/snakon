package network

import (
	"net"
	gametypes "snakon/gameTypes"
	"sync"
)

type GameNetwork struct {
	client_addr *net.UDPAddr
	server      *ServerInfo
	conn        *net.UDPConn
	conn_mutex  sync.Mutex
	state       *NetworkGameState
}

type ServerInfo struct {
	Mutex             sync.Mutex
	addr              *net.UDPAddr
	last_msg_ordering uint64
	n_msgs_received   uint64
}

type NetworkGameState struct {
	Mutex   sync.Mutex
	players map[int32]*PlayerClientState
}

type PlayerClientState struct {
	Mutex sync.Mutex
	Pos   gametypes.Position
}
