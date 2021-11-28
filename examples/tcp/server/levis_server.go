package main

import (
	"levis"
	"log"
	"net"
	"tcplevis"
	"os"
)

func handleA(l *net.TCPConn, a *net.TCPAddr, m *levis.Message) *levis.Message {
	log.Printf("Got message in handleA: %#v from %v", m, l.RemoteAddr().String())
	if m.IsConfirmable() {
		res := &levis.Message{
			Type:						levis.Acknowledgement,
			Code:						levis.POST,
			ContentFormat:	m.ContentFormat,
			MessageID:			m.MessageID,
			Payload:				m.Payload,
		}
		log.Printf("Transmitting from A %#v", res)
		return res
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	mux := tcplevis.NewServeMux()
	mux.Handle(tcplevis.FuncHandler(handleA))

	log.Fatal(tcplevis.ListenAndServe("tcp", os.Args[1], mux))
}
