package client

import gonet "net"

func (this *Client) handleConn(conn gonet.Conn) error {
	var b = make([]byte, 1024)
	for {
		conn.Read(b)
	}
}

func (this *Client) remoteServe() error {
	conn, err := this.connectSvr()
	if err != nil {
		return err
	}
	return this.handleConn(conn)
}
