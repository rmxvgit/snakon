package client

import "net"

type Client struct {
	server_addr *net.UDPAddr
	client_addr *net.UDPAddr
	recv_conn   *net.UDPConn
}
