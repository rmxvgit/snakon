package server

import (
	"net"
	"sync"
)

type Server struct {
	addr        *net.UDPAddr
	conn        *net.UDPConn
	state_mutex sync.Mutex
	game_state  *GameState
}
