package client

import gonet "net"

func (this *Client) localServe() error {
	//create tun
	select {
	case conn := <-this.onConn:
		this.localHandle(conn)
	}
	return nil
}

func (this *Client) localHandle(remoteConn gonet.Conn) {
	for {

	}
}
