package messages

import "encoding/binary"

/*
Estrutura do pacote:

	1byte		  4byte

| message_type | player_id |
*/
func (msg *NewPlayerDto) Encode() (packet []byte) {
	packet = make([]byte, 1024)
	packet[0] = byte(NEW_PLAYER_MESSAGE)
	binary.BigEndian.PutUint32(packet[1:5], uint32(msg.PlayerID))

	return
}

func DecodeNewPlayerMessage(data []byte) (msg NewPlayerDto) {
	msg = NewPlayerDto{}

	msg.PlayerID = int32(binary.BigEndian.Uint32(data[1:5]))

	return msg
}

/*
Estrutura de cada pacote gerado:

	1byte		  8byte				 1byte					12bytes_each

| message_type | message_ordering | number_of_positions | [ player_id | player_x | player_y ]
*/
func (msg *ManyPlayerPositionDto) Encode(ordering uint64) (packets [][]byte) {
	const SIZE_OF_ONE_POS = 12
	const MAX_NUM_POS_PER_PACKET = 5
	const PACKET_HEADER_SIZE = 10
	num_of_positions := len(msg.Positions)
	packets = make([][]byte, (num_of_positions/MAX_NUM_POS_PER_PACKET)+1)

	for pack_index := range packets {

		pack_first_pos := MAX_NUM_POS_PER_PACKET * pack_index
		pack_last_pos := min(pack_first_pos+MAX_NUM_POS_PER_PACKET, num_of_positions)

		packets[pack_index] = make([]byte, 1024)

		packets[pack_index][0] = byte(MANY_PLAYER_POS_MSG)
		binary.BigEndian.PutUint64(packets[pack_index][1:9], uint64(ordering))
		packets[pack_index][9] = byte(pack_last_pos - pack_first_pos)

		for pos_index, pos_msg := range msg.Positions[pack_first_pos:pack_last_pos] {
			offset := (SIZE_OF_ONE_POS * pos_index) + PACKET_HEADER_SIZE
			binary.BigEndian.PutUint32(packets[pack_index][offset:offset+4], uint32(pos_msg.PlayerID))
			binary.BigEndian.PutUint32(packets[pack_index][offset+4:offset+8], uint32(pos_msg.Pos.X))
			binary.BigEndian.PutUint32(packets[pack_index][offset+8:offset+12], uint32(pos_msg.Pos.Y))
		}
	}

	return packets
}

func DecodeManyPlayerPositionMessage(data []byte) (msg ManyPlayerPositionDto, ordering uint64) {
	const SIZE_OF_ONE_POS = 12
	const PACKET_HEADER_SIZE = 10

	msg = ManyPlayerPositionDto{}

	num_of_positions := int(data[9])
	ordering = binary.BigEndian.Uint64(data[1:9])
	msg.Positions = make([]PlayerPositionDto, num_of_positions)

	for pos_index := range num_of_positions {
		offset := (SIZE_OF_ONE_POS * pos_index) + PACKET_HEADER_SIZE
		msg.Positions[pos_index] = PlayerPositionDto{}
		msg.Positions[pos_index].PlayerID = int32(binary.BigEndian.Uint32(data[offset : offset+4]))
		msg.Positions[pos_index].Pos.X = int32(binary.BigEndian.Uint32(data[offset+4 : offset+8]))
		msg.Positions[pos_index].Pos.Y = int32(binary.BigEndian.Uint32(data[offset+8 : offset+12]))
	}

	return
}

/*
	1byte		    8bytes			   4bytes	   4bytes	  4bytes

| message_type | message_ordering | player_id | player_x | player_y |
*/
func (msg *PlayerPositionDto) Encode(ordering uint64) (packet []byte) {
	packet = make([]byte, 1024)

	packet[0] = byte(PLAYER_POS_MESSAGE)
	binary.BigEndian.PutUint64(packet[1:9], ordering)
	binary.BigEndian.PutUint32(packet[9:13], uint32(msg.PlayerID))
	binary.BigEndian.PutUint32(packet[13:17], uint32(msg.Pos.X))
	binary.BigEndian.PutUint32(packet[17:21], uint32(msg.Pos.Y))

	return
}

func DecodePlayerPositionMessage(data []byte) (msg PlayerPositionDto, ordering uint64) {
	msg = PlayerPositionDto{}

	ordering = binary.BigEndian.Uint64(data[1:9])
	msg.PlayerID = int32(binary.BigEndian.Uint32(data[9:13]))
	msg.Pos.X = int32(binary.BigEndian.Uint32(data[13:17]))
	msg.Pos.Y = int32(binary.BigEndian.Uint32(data[17:21]))

	return
}
