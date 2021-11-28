package udplevis

import (
	"levis"
	"log"
	"net"
	"time"
)

// Handler is a type that handles Levis messages.
type Handler interface {
	// Handle the message and optionally return a response message.
	ServeLevis(l *net.UDPConn, a *net.UDPAddr, m *levis.Message) *levis.Message
}

type funcHandler func(l *net.UDPConn, a *net.UDPAddr, m *levis.Message) *levis.Message

func (f funcHandler) ServeLevis(l *net.UDPConn, a *net.UDPAddr, m *levis.Message) *levis.Message {
	return f(l, a, m)
}

// FuncHandler builds a handler from a function.
func FuncHandler(f func(l *net.UDPConn, a *net.UDPAddr, m *levis.Message) *levis.Message) Handler {
	return funcHandler(f)
}

func handlePacket(l *net.UDPConn, data []byte, u *net.UDPAddr,
	rh Handler) {

	msg, err := levis.ParseMessage(data)
	if err != nil {
		log.Printf("Error parsing %v", err)
		return
	}

	rv := rh.ServeLevis(l, u, &msg)
	if rv != nil {
		Transmit(l, u, *rv)
	}
}

// ListenAndServe binds to the given address and serve requests forever.
func ListenAndServe(n, addr string, rh Handler) error {
	uaddr, err := net.ResolveUDPAddr(n, addr)
	if err != nil {
		return err
	}

	l, err := net.ListenUDP(n, uaddr)
	if err != nil {
		return err
	}

	return Serve(l, rh)
}

// Serve processes incoming UDP packets on the given listener, and processes
// these requests forever (or until the listener is closed).
func Serve(listener *net.UDPConn, rh Handler) error {
	buf := make([]byte, levis.MaxPktLen)
	for {
		nr, addr, err := listener.ReadFromUDP(buf)
		if err != nil {
			if neterr, ok := err.(net.Error); ok && (neterr.Temporary() || neterr.Timeout()) {
				time.Sleep(5 * time.Millisecond)
				continue
			}
			return err
		}
		tmp := make([]byte, nr)
		copy(tmp, buf)
		go handlePacket(listener, tmp, addr, rh)
	}
}
