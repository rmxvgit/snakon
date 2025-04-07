package internet

import (
	"errors"
	"net"
)

func NewSimpleSender() *SimpleSender {
	return &SimpleSender{}
}

func (sender *SimpleSender) CompleteSend(src_addr, dest_addr string, message []byte) (err error) {
	sender.src_addr, err = net.ResolveUDPAddr("udp4", src_addr)
	if err != nil {
		return err
	}

	sender.destination_addr, err = net.ResolveUDPAddr("udp4", dest_addr)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp4", sender.src_addr, sender.destination_addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	bytesWritten, err := conn.Write(message)
	if err != nil {
		return err
	}

	if bytesWritten != len(message) {
		return errors.New("not all bytes were written")
	}

	return
}
