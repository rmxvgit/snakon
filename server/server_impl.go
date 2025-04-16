package server

import (
	"fmt"
	"net"
	"os"
	gametypes "snakon/gameTypes"
	"snakon/internet/messages"
	"snakon/utils"
	"sync"
)

func Run() {
	addr := os.Args[2] // listening addr

	server, err := NewServer(addr)
	utils.PanicOnError(err)

	server.Listen() // the main loop of the server program
}

func NewServer(addr_str string) (serv *Server, err error) {
	serv = &Server{}

	serv.state = &ServerState{}

	serv.addr, err = net.ResolveUDPAddr("udp4", addr_str)
	if err != nil {
		return nil, err
	}

	serv.conn, err = net.ListenUDP("udp4", serv.addr)
	if err != nil {
		return nil, err
	}

	return serv, nil
}

func (server *Server) Listen() {
	for {
		buf := make([]byte, 1024)

		_, remote_addr, err := server.conn.ReadFromUDP(buf)
		if err != nil {
			server.LogError(err)
			continue
		}

		go server.HandleMessage(remote_addr, buf)
	}
}

func (server *Server) LogError(err error) {}

func (server *Server) HandleMessage(remote_addr *net.UDPAddr, data []byte) {
	flag := messages.MessageFlag(data[0])

	var err_code MsgErrCode = MsgOk
	var msg_err error = nil

	switch flag {

	case messages.NEW_PLAYER_MESSAGE:
		msg_err, err_code = server.HandleNewPlayerMessage(remote_addr, data)

	case messages.PLAYER_POS_MESSAGE:
		msg_err, err_code = server.HandlePlayerPositionMessage(remote_addr, data)
	}

	if msg_err != nil {
		server.HandleMessageError(msg_err, err_code, remote_addr, data)
	}

	// send response
	server.SendResponse(remote_addr, flag)
}

func (server *Server) HandleMessageError(err error, errCode MsgErrCode, remote_addr *net.UDPAddr, data []byte) {
}

func (server *Server) HandlePlayerPositionMessage(remote_addr *net.UDPAddr, msg_data []byte) (error, MsgErrCode) {
	msg_info, ordering := messages.DecodePlayerPositionMessage(msg_data)

	player_state, exists := server.state.players[msg_info.PlayerID]
	if !exists {
		return fmt.Errorf("non existent player: %d", msg_info.PlayerID), NullPlayerReference
	}

	player_state.Mutex.Lock()
	defer player_state.Mutex.Unlock()

	if player_state.LastMessage > ordering {
		return fmt.Errorf("packet dropped"), MessageDrop
	}

	player_state.Pos.X = msg_info.Pos.X
	player_state.Pos.Y = msg_info.Pos.Y

	return nil, MsgOk
}

func (server *Server) HandleNewPlayerMessage(remote_addr *net.UDPAddr, data []byte) (error, MsgErrCode) {
	msg_info := messages.DecodeNewPlayerMessage(data)

	_, exists := server.state.players[msg_info.PlayerID]
	if exists {
		return fmt.Errorf("redeclaration of player: %d", msg_info.PlayerID), PlayerRedeclaration
	}

	server.state.players[msg_info.PlayerID] = &PlayerServerState{
		LastMessage: 0,
		Mutex:       sync.Mutex{},
		Pos:         gametypes.Position{},
	}

	return nil, MsgOk
}

func (server *Server) SendResponse(remote_addr *net.UDPAddr, flag messages.MessageFlag) {
	position_messages := server.generateAllPlayerPositionMessages()
	position_packets := position_messages.Encode(0)

	server.SendPackets(remote_addr, position_packets)
}

func (server *Server) generateAllPlayerPositionMessages() messages.ManyPlayerPositionDto {
	return messages.ManyPlayerPositionDto{}
}

func (server *Server) SendPackets(dest_addr *net.UDPAddr, packets [][]byte) {
	for i := range packets {
		n, err := server.conn.WriteToUDP(packets[i], dest_addr)

		if err != nil {
			server.LogError(err)
			continue
		}

		if n != PACKET_SIZE {
			err = fmt.Errorf("não foi possível escrever toda a mensagem")
			server.LogError(err)
			continue
		}
	}
}

func (server *Server) SendPacket(dest_addr *net.UDPAddr, packets []byte) {}
