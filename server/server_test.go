package server_test

import (
	"fmt"
	"snakon/server"
	"testing"
)

func TestEncodeManyPlayerPositionMessage(t *testing.T) {

	msg := server.ManyPlayerPositionMessage{
		Positions: []server.PlayerPositionMessage{
			{PlayerID: 1, Pos: server.PlayerPos{X: 1, Y: 1}},
			{PlayerID: 2, Pos: server.PlayerPos{X: 2, Y: 2}},
			{PlayerID: 3, Pos: server.PlayerPos{X: 3, Y: 3}},
			{PlayerID: 4, Pos: server.PlayerPos{X: 4, Y: 4}},
			{PlayerID: 5, Pos: server.PlayerPos{X: 5, Y: 5}},
			{PlayerID: 6, Pos: server.PlayerPos{X: 6, Y: 6}},
			{PlayerID: 7, Pos: server.PlayerPos{X: 7, Y: 7}},
			{PlayerID: 8, Pos: server.PlayerPos{X: 8, Y: 8}},
			{PlayerID: 9, Pos: server.PlayerPos{X: 9, Y: 9}},
		},
	}

	packets := msg.EncodeManyPlayerPositionMessage()
	for _, pack := range packets {
		fmt.Print(pack[0:80])
		fmt.Println()
	}
}
