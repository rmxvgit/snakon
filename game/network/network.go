package network

import (
	"net"
	"snakon/server"
	"sync"
)

type GameNetManager struct {
	Addr             *net.UDPAddr
	Conn             *net.UDPConn
	Conn_lock_mutex  sync.Mutex
	Game_state_mutex sync.Mutex
	Game_state       *server.GameState
}
