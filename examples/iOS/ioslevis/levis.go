package ioslevis

import (
	"levis"
	"log"
	"net"
	"tcplevis"
	"udplevis"
	"time"
	"fmt"
)

/**** TCP ****/
func StartTCPServer(endpoint string){
	mux := tcplevis.NewServeMux()
	mux.Handle(tcplevis.FuncHandler(handleTCP))

	log.Println("TCP Server Started: " + endpoint)
	log.Fatal(tcplevis.ListenAndServe("tcp", endpoint, mux))
}

func handleTCP(l *net.TCPConn, a *net.TCPAddr, m *levis.Message) *levis.Message {
	log.Printf("Got message in handleTCP: %#v from %v", m, l.RemoteAddr().String())
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

func StartTCPClient(endpoint string, payload string){
	req := levis.Message{
		Type:						levis.Confirmable,
		Code:						levis.GET,
		MessageID:			255,
		ContentFormat:	levis.AppJSON,
		Payload:				[]byte(payload),
	}

	start := time.Now()
	c, err := tcplevis.Dial("tcp", endpoint)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		//log.Printf("Response payload: %s", rv.Payload)
		fmt.Printf("%v\n", time.Since(start).Seconds())
	}
}


/*** UDP ***/
func StartUDPServer(endpoint string){
	mux := udplevis.NewServeMux()
	mux.Handle(udplevis.FuncHandler(handleUDP))

	log.Println("UDP Server Started: " + endpoint)
	log.Fatal(udplevis.ListenAndServe("udp", endpoint, mux))
}

func StartUDPClient(endpoint string, payload string){
	req := levis.Message{
		Type:						levis.Confirmable,
		Code:						levis.GET,
		MessageID:			255,
		ContentFormat:	levis.AppJSON,
		Payload:				[]byte(payload),
  }

	start := time.Now()
	c, err := udplevis.Dial("udp", endpoint)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		//log.Printf("Response payload: %s", rv.Payload)
		fmt.Printf("%v\n", time.Since(start).Seconds())
	}
}

func handleUDP(l *net.UDPConn, a *net.UDPAddr, m *levis.Message) *levis.Message {
	log.Printf("Got message in handleUDP: %#v from %v", m, a.String())
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

