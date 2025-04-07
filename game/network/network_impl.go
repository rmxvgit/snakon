package network

import (
	"encoding/binary"
	"net"
	"snakon/server"
	"snakon/utils"
)

func SetupNetManager(addr_str string) (net_manager *GameNetManager, err error) {
	net_manager = &GameNetManager{}

	net_manager.Game_state = server.NewGameState()

	err = net_manager.SetAddr(addr_str)
	if err != nil {
		return nil, err
	}
	err = net_manager.SetConn()
	if err != nil {
		return nil, err
	}

	go net_manager.Listen()

	return net_manager, nil
}

func (net_magager *GameNetManager) SetAddr(addr_str string) (err error) {
	net_magager.Addr, err = net.ResolveUDPAddr("udp4", addr_str)
	return err
}

func (net_manager *GameNetManager) SetConn() (err error) {
	net_manager.Conn, err = net.ListenUDP("udp4", net_manager.Addr)
	return err
}

func (net_manager *GameNetManager) Listen() {
	for {
		buf := make([]byte, 1024)

		_, remote_addr, err := net_manager.Conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		go net_manager.HandleMessage(remote_addr, buf)
	}
}

func (net_manager *GameNetManager) HandleMessage(remote_addr *net.UDPAddr, data []byte) {
	flag := server.MessageFlag(data[0])

	switch flag {
	case server.MANY_PLAYER_POS_MSG:
		net_manager.HandleManyPlayerPosMessage(data)
	}
}

func (net_manager *GameNetManager) HandleManyPlayerPosMessage(data []byte) {
	msg := server.DecodeManyPlayerPositionMessage(data)

	for _, pl_pos := range msg.Positions {
		net_manager.Game_state_mutex.Lock()

		if net_manager.Game_state.PlayersState[pl_pos.PlayerID] == nil {
			net_manager.Game_state.PlayersState[pl_pos.PlayerID] = &server.PlayerState{
				ID: pl_pos.PlayerID,
				Pos: server.PlayerPos{
					X: pl_pos.Pos.X,
					Y: pl_pos.Pos.Y,
				},
			}
		} else {
			net_manager.Game_state.PlayersState[pl_pos.PlayerID].Pos = pl_pos.Pos
		}
		net_manager.Game_state_mutex.Lock()
	}
}

func (net_manager *GameNetManager) SendId(id uint32) {
	serv_addr, err := net.ResolveUDPAddr("udp4", ":3001")
	utils.PanicOnError(err)
	msg := make([]byte, 1024)
	msg[0] = 0
	binary.BigEndian.PutUint32(msg[1:5], id)
	net_manager.Conn.WriteToUDP(msg, serv_addr)
}

func (net_manager *GameNetManager) SendMyPos(id, x, y uint32) {
	serv_addr, err := net.ResolveUDPAddr("udp4", ":3001")
	utils.PanicOnError(err)
	msg := make([]byte, 1024)
	msg[0] = 1
	binary.BigEndian.PutUint32(msg[1:5], id)
	binary.BigEndian.PutUint32(msg[5:9], x)
	binary.BigEndian.PutUint32(msg[9:13], y)

	net_manager.Conn.WriteToUDP(msg, serv_addr)
}
