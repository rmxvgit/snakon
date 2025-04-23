package network

import (
	"fmt"
	"net"
	gametypes "snakon/gameTypes"
	"snakon/internet/messages"
	"snakon/utils"
	"sync"
)

func NewGameNetwork(my_addr, serv_addr string) (network *GameNetwork) {
	network = &GameNetwork{}
	var err error

	network.state = NewGameState()
	network.server = NewServerInfo()

	network.server.n_msgs_received = 0
	network.server.last_msg_ordering = 0

	network.client_addr, err = net.ResolveUDPAddr("udp4", my_addr)
	utils.PanicOnError(err)

	network.server.addr, err = net.ResolveUDPAddr("udp4", serv_addr)
	utils.PanicOnError(err)

	println(my_addr)

	network.conn, err = net.ListenUDP("udp4", network.client_addr)
	utils.PanicOnError(err)

	return network
}

func (network *GameNetwork) Listen() {
	for {
		buf := make([]byte, 1024)

		_, remote_addr, err := network.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		go network.HandleMessage(remote_addr, buf)
	}
}

func (network *GameNetwork) HandleMessage(remote_addr *net.UDPAddr, data []byte) {
	flag := messages.MessageFlag(data[0])
	var err error

	switch flag {
	case messages.MANY_PLAYER_POS_MSG:
		err = network.HandleManyPlayerPosMessage(remote_addr, data)
	}

	if err != nil {
		network.LogError(err)
	}
}

func (network *GameNetwork) LogError(err error) {
	fmt.Println(err)
}

func (network *GameNetwork) HandleManyPlayerPosMessage(remote_addr *net.UDPAddr, data []byte) (err error) {
	msg, ordering := messages.DecodeManyPlayerPositionMessage(data)

	network.server.Mutex.Lock()

	network.server.n_msgs_received++
	if network.server.last_msg_ordering > ordering {
		return fmt.Errorf("dropped packet")
	}

	network.server.Mutex.Unlock()

	network.state.Mutex.Lock()

	for _, pos_msg := range msg.Positions {
		player_state, exists := network.state.players[pos_msg.PlayerID]

		if !exists {
			player_state = &PlayerClientState{}
		}

		player_state.Pos.X = pos_msg.Pos.X
		player_state.Pos.Y = pos_msg.Pos.Y
	}

	network.state.Mutex.Unlock()

	return nil
}

func (network *GameNetwork) SendNewPlayerNotification(player_id int32) error {
	network.conn_mutex.Lock()
	defer network.conn_mutex.Unlock()

	msg := messages.NewPlayerDto{
		PlayerID: player_id,
	}

	data := msg.Encode()
	_, err := network.conn.WriteToUDP(data, network.server.addr)
	utils.PanicOnError(err)

	return nil
}

func (network *GameNetwork) SendPlayerPosition(player_id int32, pos gametypes.Position) error {
	network.conn_mutex.Lock()
	defer network.conn_mutex.Unlock()

	msg := messages.PlayerPositionDto{
		PlayerID: player_id,
		Pos:      pos,
	}

	data := msg.Encode(network.server.n_msgs_received)
	network.server.n_msgs_received++

	_, err := network.conn.WriteToUDP(data, network.server.addr)
	utils.PanicOnError(err)

	return nil
}

func (network *GameNetwork) GetOtherPlayersPositions() map[int32]gametypes.Position {
	positions := make(map[int32]gametypes.Position)

	network.state.Mutex.Lock()
	for id, player_state := range network.state.players {
		positions[id] = player_state.Pos
	}
	network.state.Mutex.Unlock()

	return positions
}

func NewGameState() *NetworkGameState {
	return &NetworkGameState{
		Mutex:   sync.Mutex{},
		players: make(map[int32]*PlayerClientState),
	}
}

func NewServerInfo() *ServerInfo {
	return &ServerInfo{
		Mutex:             sync.Mutex{},
		addr:              nil,
		last_msg_ordering: 0,
		n_msgs_received:   0,
	}
}
