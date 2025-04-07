package internet

import "net"

type SimpleSender struct {
	destination_addr *net.UDPAddr
	src_addr         *net.UDPAddr
	send_conn        *net.UDPConn
}
