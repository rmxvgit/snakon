package network

import (
	"encoding/binary"
	"fmt"
	"net"
	"snakon/server"
)

func NewGameNetwork(my_addr, serv_addr string) (network *GameNetwork, err error) {
	network = &GameNetwork{}

	network.game_state = server.NewGameState()

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
	flag := server.MessageFlag(data[0])

	switch flag {
	case server.MANY_PLAYER_POS_MSG:
		network.HandleManyPlayerPosMessage(data)
	}
}

func (network *GameNetwork) HandleManyPlayerPosMessage(data []byte) {
	msg := server.DecodeManyPlayerPositionMessage(data)

	players_state := network.game_state.PlayersState

	network.game_state_mutex.Lock()
	for _, other_player := range msg.Positions {
		if players_state[other_player.PlayerID] == nil {
			players_state[other_player.PlayerID] = server.NewPlayerState(other_player.PlayerID)
		}
		players_state[other_player.PlayerID].Pos = other_player.Pos
	}
	network.game_state_mutex.Unlock()
}

func (network *GameNetwork) NewPlayer(id uint32) {
	msg := make([]byte, 1024)
	msg[0] = byte(server.NEW_PLAYER_MESSAGE)

	binary.BigEndian.PutUint32(msg[1:5], uint32(id))

	network.conn_mutex.Lock()
	network.conn.WriteToUDP(msg, network.server_addr)
	network.conn_mutex.Unlock()

}

func (network *GameNetwork) SendEntityPos(id int32, pos server.PlayerPos) {
	msg := make([]byte, 1024)
	msg[0] = byte(server.PLAYER_POS_MESSAGE)

	binary.BigEndian.PutUint32(msg[1:5], uint32(id))
	binary.BigEndian.PutUint32(msg[5:9], uint32(pos.X))
	binary.BigEndian.PutUint32(msg[9:13], uint32(pos.Y))

	network.conn_mutex.Lock()
	network.conn.WriteToUDP(msg, network.server_addr)
	network.conn_mutex.Unlock()
}

func (network *GameNetwork) GetServerGameState() server.GameState {
	network.game_state_mutex.Lock()
	defer network.game_state_mutex.Unlock()
	return *network.game_state
}

func (network *GameNetwork) GetPlayerState(id int32) (server.PlayerState, error) {
	network.game_state_mutex.Lock()
	state, exists := network.game_state.PlayersState[id]
	defer network.game_state_mutex.Unlock()
	if !exists {
		return server.PlayerState{}, fmt.Errorf("non existent player")
	}
	return *state, nil
}
