package tcplevis

import (
	"levis"
	"log"
	"net"
	"time"
)

// Handler is a type that handles Levis TCP Messages.
type Handler interface {
	// Handle the message and optionally return a response message.
	ServeLevis(l *net.TCPConn, a *net.TCPAddr, m *levis.Message) *levis.Message
}

type funcHandler func(l *net.TCPConn, a *net.TCPAddr, m *levis.Message) *levis.Message

func (f funcHandler) ServeLevis(l *net.TCPConn, a *net.TCPAddr, m *levis.Message) *levis.Message {
	return f(l, a, m)
}

// FuncHandler builds a handler from a function.
func FuncHandler(f func(l *net.TCPConn, a *net.TCPAddr, m *levis.Message) *levis.Message) Handler {
	return funcHandler(f)
}

func handlePacket(l *net.TCPConn, data []byte, u *net.TCPAddr, rh Handler) {
	msg, err := levis.ParseMessage(data)
	if err != nil {
		log.Printf("Error parsing: %v", err)
		return
	}

	rv := rh.ServeLevis(l, u, &msg)
	if rv != nil {
		Transmit(l, u, *rv)
	}
}

// ListenAndServe binds to the given address and serve requests forever.
func ListenAndServe(n, addr string, rh Handler) error {
	taddr, err := net.ResolveTCPAddr(n, addr)

	if err != nil {
		return err
	}

	l, err := net.ListenTCP(n, taddr)
	if err != nil {
		return err
	}

	return Serve(l, rh)
}

// Serve processes incoming TCP packets on the given listener, and proccesses these requests forever (or until the listener is closed).
func Serve(listener *net.TCPListener, rh Handler) error {
	buf := make([]byte, levis.MaxPktLen)

	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			return err
		}

		nr, err := conn.Read(buf)

		if err != nil {
			if neterr, ok := err.(net.Error); ok && (neterr.Temporary() || neterr.Timeout()) {
				time.Sleep(5 * time.Millisecond)
				continue
			}
			return err
		}

		tmp := make([]byte, nr)
		copy(tmp, buf)

		go handlePacket(conn, tmp, nil, rh)
	}
}
