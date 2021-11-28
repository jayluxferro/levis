package tcplevis

import (
	"levis"
	"net"
)

// Conn is a Levis TCP client connection
type Conn struct {
	conn *net.TCPConn
	buf  []byte
}

// connects a Levis client over TCP
func Dial(n, addr string) (*Conn, error) {
	uaddr, err := net.ResolveTCPAddr(n, addr)
	if err != nil {
		return nil, err
	}

	s, err := net.DialTCP("tcp", nil, uaddr)

	if err != nil {
		return nil, err
	}

	return &Conn{s, make([]byte, levis.MaxPktLen)}, nil
}

// send a message and receive a response
func (c *Conn) Send(req levis.Message) (*levis.Message, error) {
	err := Transmit(c.conn, nil, req)
	if err != nil {
		return nil, err
	}

	if !req.IsConfirmable() {
		return nil, nil
	}

	rv, err := Receive(c.conn, c.buf)
	if err != nil {
		return nil, err
	}

	return &rv, nil
}
