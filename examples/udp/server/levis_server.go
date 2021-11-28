package main

import (
	"log"
	"net"
	"levis"
  "udplevis"
	"os"
)

func handleA(l *net.UDPConn, a *net.UDPAddr, m *levis.Message) *levis.Message {
	log.Printf("Got message in handleA: %#v from %v", m, a.String())
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

	mux := udplevis.NewServeMux()
	mux.Handle(udplevis.FuncHandler(handleA))

	log.Fatal(udplevis.ListenAndServe("udp", os.Args[1], mux))
}
