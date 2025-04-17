package sender

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	gametypes "snakon/gameTypes"
	"snakon/internet/messages"
	"snakon/utils"
)

func Run() {
	source_addr, err := net.ResolveUDPAddr("udp4", os.Args[2])
	utils.PanicOnError(err)

	dest_addr, err := net.ResolveUDPAddr("udp4", os.Args[3])
	utils.PanicOnError(err)

	conn, err := net.ListenUDP("udp4", source_addr)
	utils.PanicOnError(err)
	defer conn.Close()

	msg := messages.PlayerPositionDto{PlayerID: 5, Pos: gametypes.Position{X: 30, Y: 40}}

	var ordering_counter uint64 = 0

	go Listen(conn, dest_addr)

	new_player_message := messages.NewPlayerDto{PlayerID: 5}
	data := new_player_message.Encode()
	_, err = conn.WriteToUDP(data, dest_addr)
	utils.PanicOnError(err)

	for {
		data := msg.Encode(uint64(ordering_counter))
		msg.Pos.X = rand.Int31n(20)
		msg.Pos.Y = rand.Int31n(20)
		ordering_counter++

		fmt.Scanln()
		println("batata")

		_, err := conn.WriteToUDP(data, dest_addr)
		utils.PanicOnError(err)
	}
}

func Listen(conn *net.UDPConn, send_addr *net.UDPAddr) {
	for {
		data := make([]byte, 1024)
		_, addr, err := conn.ReadFromUDP(data)
		println(err)
		utils.PanicOnError(err)

		fmt.Printf("Received message from %s: %v\n", addr.String(), data)
	}
}
