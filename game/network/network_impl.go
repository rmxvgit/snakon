package network

import (
	"fmt"
	"net"
	"snakon/internet/messages"
	"snakon/utils"
)

func NewGameNetwork(my_addr, serv_addr string) (network *GameNetwork) {
	network = &GameNetwork{}
	var err error

	//network.game_state = server.NewGameState()

	network.client_addr, err = net.ResolveUDPAddr("udp4", my_addr)
	utils.PanicOnError(err)

	network.server.addr, err = net.ResolveUDPAddr("udp4", serv_addr)
	utils.PanicOnError(err)

	network.server.n_msgs_received = 0
	network.server.last_msg_ordering = 0

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
