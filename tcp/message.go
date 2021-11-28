package tcplevis

import (
	"levis"
	"net"
	"time"
)

// Transmit a message
func Transmit(l *net.TCPConn, a *net.TCPAddr, m levis.Message) error {
	d, err := m.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = l.Write(d)
	return err
}

// Receive a message
func Receive(l *net.TCPConn, buf []byte) (levis.Message, error) {
	l.SetReadDeadline(time.Now().Add(levis.ResponseTimeout))

	_, err := l.Read(buf)
	if err != nil {
		return levis.Message{}, err
	}

	return levis.ParseMessage(buf)
}
