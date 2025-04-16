package network

import (
	"net"
	"snakon/internet/messages"
)

func NewGameNetwork(my_addr, serv_addr string) (network *GameNetwork, err error) {
	network = &GameNetwork{}

	//network.game_state = server.NewGameState()

	network.client_addr, err = net.ResolveUDPAddr("udp4", my_addr)
	if err != nil {
		return nil, err
	}

	network.server_addr, err = net.ResolveUDPAddr("udp4", serv_addr)
	if err != nil {
		return nil, err
	}

	network.conn, err = net.ListenUDP("udp4", network.client_addr)
	if err != nil {
		return nil, err
	}

	return network, nil
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

	switch flag {
	case messages.MANY_PLAYER_POS_MSG:
		network.HandleManyPlayerPosMessage(data)
	}
}

func (network *GameNetwork) HandleManyPlayerPosMessage(data []byte) {
}
