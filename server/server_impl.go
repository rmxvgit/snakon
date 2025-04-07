package server

import (
	"net"
	"os"
	"snakon/utils"
)

func Run() {
	addr := os.Args[2]
	server, err := NewServer(addr)
	utils.PanicOnError(err)
	server.Listen()
}

func NewServer(addr_str string) (serv *Server, err error) {
	serv = &Server{}

	serv.game_state = NewGameState()

	err = serv.SetAddr(addr_str)
	if err != nil {
		return nil, err
	}

	err = serv.SetConn()
	if err != nil {
		return nil, err
	}

	return serv, nil
}

func (server *Server) SetAddr(addr_str string) (err error) {
	server.addr, err = net.ResolveUDPAddr("udp4", addr_str)
	return err
}

func (server *Server) SetConn() (err error) {
	server.conn, err = net.ListenUDP("udp4", server.addr)
	return err
}

func (server *Server) Listen() {
	for {
		buf := make([]byte, 1024)

		_, remote_addr, err := server.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		go server.HandleMessage(remote_addr, buf)
	}
}

func (server *Server) HandleMessage(remote_addr *net.UDPAddr, data []byte) {
	flag := MessageFlag(data[0])

	switch flag {
	case NEW_PLAYER_MESSAGE:
		server.HandleNewPlayer(remote_addr, data)
	case PLAYER_POS_MESSAGE:
		server.HandlePositionMessage(remote_addr, data)
	}

	// send response
	server.SendGameState(remote_addr, data)
}

func (server *Server) HandlePositionMessage(remote_addr *net.UDPAddr, data []byte) {
	msg, err := DecodePlayerPositionMessage(data)
	utils.PanicOnError(err)

	server.state_mutex.Lock()
	if server.game_state.PlayersState[msg.PlayerID] == nil {
		server.game_state.PlayersState[msg.PlayerID] = NewPlayerState(msg.PlayerID)
	}
	server.game_state.PlayersState[msg.PlayerID].Pos.X = msg.Pos.X
	server.game_state.PlayersState[msg.PlayerID].Pos.Y = msg.Pos.Y
	server.state_mutex.Unlock()
}

func (server *Server) HandleNewPlayer(remote_addr *net.UDPAddr, data []byte) {
	msg, err := DecodeNewPlayerMessage(data)
	utils.PanicOnError(err)

	server.state_mutex.Lock()
	server.game_state.PlayersState[msg.PlayerID] = NewPlayerState(msg.PlayerID)
	server.state_mutex.Unlock()
}

func (server *Server) SendGameState(remote_addr *net.UDPAddr, data []byte) {
	data_packets := EncodeGameStateMessage(server.game_state)

	for _, packet := range data_packets {
		server.conn.WriteToUDP(packet, remote_addr)
	}
}

func EncodeGameStateMessage(game_state *GameState) [][]byte {

	players_position_data := game_state.ManyPlayerPositionMessage()
	encoded_position_data_packets := players_position_data.EncodeManyPlayerPositionMessage()

	//you can add more packets here if needed

	return encoded_position_data_packets
}
