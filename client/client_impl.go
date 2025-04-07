package client

import "net"

func NewClient(recv_addr, send_addr, server_addr string) (client *Client, err error) {
	client = &Client{}

	err = client.SetRecvAddr(recv_addr)
	if err != nil {
		return
	}

	err = client.SetSendAddr(send_addr)
	if err != nil {
		return
	}

	err = client.SetServerAddr(server_addr)
	if err != nil {
		return
	}

	err = client.MakeConnection()
	if err != nil {
		return
	}

	return
}

func (client *Client) Recv() error {
	return nil
}

func Send(msg []byte) error {
	return nil
}

func (client *Client) SetServerAddr(addr string) (err error) {
	client.server_addr, err = net.ResolveUDPAddr("udp4", addr)
	return
}

func (client *Client) SetRecvAddr(addr string) (err error) {
	client.server_addr, err = net.ResolveUDPAddr("udp4", addr)
	return
}

func (client *Client) SetSendAddr(addr string) (err error) {
	client.client_addr, err = net.ResolveUDPAddr("udp4", addr)
	return
}

func (client *Client) MakeConnection() (err error) {
	client.recv_conn, err = net.ListenUDP("udp4", client.client_addr)
	return err
}

func (client *Client) GetRecvConn() *net.UDPConn {
	return client.recv_conn
}
