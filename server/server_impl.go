package server

import (
	"fmt"
	"net"
	"os"
	gametypes "snakon/gameTypes"
	"snakon/internet/messages"
	"snakon/utils"
	"strings"
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

	serv.state = EmptyServerState()
	serv.clients = make(map[string]*ClientInfo)

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

		fmt.Println(server.state.String())

		go server.HandleMessage(remote_addr, buf)
	}
}

func (server *Server) LogError(err error) {
	println(err)
}

func (server *Server) HandleMessage(remote_addr *net.UDPAddr, data []byte) {
	server.AcountMessage(remote_addr)

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
	client := server.clients[remote_addr.String()]

	msg_info, ordering := messages.DecodePlayerPositionMessage(msg_data)

	last_msg_ordering := client.getLastMessage()

	player_state, exists := server.state.players[msg_info.PlayerID]
	if !exists {
		return fmt.Errorf("non existent player: %d", msg_info.PlayerID), NullPlayerReference
	}

	// dropa o pacote caso a mensagem recebida seja defasada
	if last_msg_ordering > ordering {
		return fmt.Errorf("packet dropped"), MessageDrop
	}

	client.setLastMessage(ordering)

	player_state.Mutex.Lock()
	defer player_state.Mutex.Unlock()

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
		Mutex: sync.Mutex{},
		Pos:   gametypes.Position{},
	}

	return nil, MsgOk
}

func (server *Server) SendResponse(remote_addr *net.UDPAddr, flag messages.MessageFlag) {
	client := server.clients[remote_addr.String()]

	position_messages := server.generateAllPlayerPositionMessages()
	position_packets := position_messages.Encode(client.getNumOfReceivedMessages())
	client.increaseNumOfReceivedMessages()

	server.SendPackets(remote_addr, position_packets)
}

func (server *Server) generateAllPlayerPositionMessages() messages.ManyPlayerPositionDto {
	players_pos_msg := messages.ManyPlayerPositionDto{
		Positions: make([]messages.PlayerPositionDto, len(server.state.players)),
	}

	list_pos := 0
	for id, player_state := range server.state.players {
		player_state.Mutex.Lock()
		players_pos_msg.Positions[list_pos].PlayerID = id
		players_pos_msg.Positions[list_pos].Pos.X = player_state.Pos.X
		players_pos_msg.Positions[list_pos].Pos.Y = player_state.Pos.Y
		player_state.Mutex.Unlock()
	}

	return players_pos_msg
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

// Increases the counter of messages received from this client by 1. If the client was not registered yet, it registers it.
// WARNING: this function does not set the number of the last message received. (variable used to keep track of the correct order of the messages
// received by a certain client)
func (server *Server) AcountMessage(remote_addr *net.UDPAddr) {
	client, exists := server.clients[remote_addr.String()]
	if !exists {
		server.RegisterNewClient(remote_addr)
		return
	}

	client.mutex.Lock()
	client.n_msg_recv++
	client.mutex.Unlock()
}

func (server *Server) RegisterNewClient(remote_addr *net.UDPAddr) {
	server.clients[remote_addr.String()] = &ClientInfo{
		n_msg_sent:    0,
		n_msg_recv:    1,
		last_msg_recv: 0,
	}
}

func (client *ClientInfo) setLastMessage(ordering uint64) {
	client.mutex.Lock()
	client.last_msg_recv = ordering
	client.mutex.Unlock()
}

func (client *ClientInfo) getLastMessage() uint64 {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	return client.last_msg_recv
}

func (client *ClientInfo) getNumOfReceivedMessages() uint64 {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	return client.n_msg_recv
}

func (client *ClientInfo) increaseNumOfReceivedMessages() {
	client.mutex.Lock()
	client.n_msg_recv++
	client.mutex.Unlock()
}

func (state *ServerState) String() string {
	buff := strings.Builder{}
	for key, value := range state.players {
		value.Mutex.Lock()
		value_representation := fmt.Sprintf("id:%d x:%d y:%d\n", key, value.Pos.X, value.Pos.Y)
		buff.WriteString(value_representation)
		value.Mutex.Unlock()
	}
	return buff.String()
}

func EmptyServerState() *ServerState {
	return &ServerState{
		players: make(map[int32]*PlayerServerState),
	}
}
