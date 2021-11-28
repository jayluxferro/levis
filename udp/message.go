package udplevis

import (
  "net"
  "time"
  "levis"
)

// Transmit a message
func Transmit(l *net.UDPConn, a *net.UDPAddr, m levis.Message) error {
  d, err := m.MarshalBinary()
  if err != nil {
    return err
  }

  if a == nil {
    _, err = l.Write(d)
  } else {
    _, err = l.WriteTo(d, a)
  }

  return err
}

// Receive a message
func Receive(l *net.UDPConn, buf []byte) (levis.Message, error) {
  l.SetReadDeadline(time.Now().Add(levis.ResponseTimeout)) 

  nr, _, err := l.ReadFromUDP(buf)
  if err != nil {
    return levis.Message{}, err
  }

  return levis.ParseMessage(buf[:nr])
}
