package server

import (
	"encoding/binary"
)

func DecodeNewPlayerMessage(data []byte) (*NewPlayerMessage, error) {
	msg := &NewPlayerMessage{}

	_, err := binary.Decode(data[1:5], binary.BigEndian, &msg.PlayerID)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func DecodePlayerPositionMessage(data []byte) (msg *PlayerPositionMessage, err error) {
	msg = &PlayerPositionMessage{}

	_, err = binary.Decode(data[1:5], binary.BigEndian, &msg.PlayerID)
	if err != nil {
		return nil, err
	}

	_, err = binary.Decode(data[5:9], binary.BigEndian, &msg.Pos.X)
	if err != nil {
		return nil, err
	}

	_, err = binary.Decode(data[9:13], binary.BigEndian, &msg.Pos.Y)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// encoded message structure:
// |message_type, number_of_positions, []position|
// position -> [player_id, x, y]
func (msg *ManyPlayerPositionMessage) EncodeManyPlayerPositionMessage() (packets [][]byte) {
	const SIZE_OF_ONE_POS = 12
	const MAX_NUM_POS_PER_PACKET = 5
	num_of_positions := len(msg.Positions)
	packets = make([][]byte, (num_of_positions/MAX_NUM_POS_PER_PACKET)+1)

	for pack_index := range packets {

		pack_first_pos := MAX_NUM_POS_PER_PACKET * pack_index
		pack_last_pos := min(pack_first_pos+MAX_NUM_POS_PER_PACKET, num_of_positions)

		packets[pack_index] = make([]byte, 1024)
		packets[pack_index][0] = byte(MANY_PLAYER_POS_MSG)
		packets[pack_index][1] = byte(pack_last_pos - pack_first_pos)

		for pos_index, pos_msg := range msg.Positions[pack_first_pos:pack_last_pos] {
			offset := (SIZE_OF_ONE_POS * pos_index) + 2
			binary.BigEndian.PutUint32(packets[pack_index][offset:offset+4], uint32(pos_msg.PlayerID))
			binary.BigEndian.PutUint32(packets[pack_index][offset+4:offset+8], uint32(pos_msg.Pos.X))
			binary.BigEndian.PutUint32(packets[pack_index][offset+8:offset+12], uint32(pos_msg.Pos.Y))
		}
	}

	return packets
}

func DecodeManyPlayerPositionMessage(data []byte) (msg *ManyPlayerPositionMessage) {
	const SIZE_OF_ONE_POS = 12
	msg = &ManyPlayerPositionMessage{}

	num_of_positions := int(data[1])
	msg.Positions = make([]PlayerPositionMessage, num_of_positions)

	for pos_index := range num_of_positions {
		offset := (SIZE_OF_ONE_POS * pos_index) + 2
		msg.Positions[pos_index] = PlayerPositionMessage{}
		msg.Positions[pos_index].PlayerID = int32(binary.BigEndian.Uint32(data[offset : offset+4]))
		msg.Positions[pos_index].Pos.X = int32(binary.BigEndian.Uint32(data[offset+4 : offset+8]))
		msg.Positions[pos_index].Pos.Y = int32(binary.BigEndian.Uint32(data[offset+8 : offset+12]))
	}

	return msg
}
