package udplevis

import (
  "net"
  "levis"
)

// Conn is a Levis client connection.
type Conn struct {
  conn *net.UDPConn
  buf []byte
}

// Connects a Levis client over UDP.
func Dial(n, addr string) (*Conn, error) {
  uaddr, err := net.ResolveUDPAddr(n, addr)
  if err != nil {
    return nil, err
  }

  s, err := net.DialUDP("udp", nil, uaddr)
  if err != nil {
    return nil, err
  }

  return &Conn{s, make([]byte, levis.MaxPktLen)}, nil
}

// send a message and receive a response if any.
func (c *Conn) Send(req levis.Message) (*levis.Message, error){
  err := Transmit(c.conn, nil, req)
  if err != nil {
    return nil, err
  }

  if !req.IsConfirmable(){
    return nil, nil
  }

  rv, err := Receive(c.conn, c.buf)
  if err != nil {
    return nil, err
  }

  return &rv, nil
}
