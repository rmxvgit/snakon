package network

import (
	"net"
	"snakon/server"
	"sync"
)

type GameNetwork struct {
	client_addr      *net.UDPAddr
	server_addr      *net.UDPAddr
	conn             *net.UDPConn
	conn_mutex       sync.Mutex
	game_state_mutex sync.Mutex
	game_state       *server.GameState
}
