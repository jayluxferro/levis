package main

import (
	"log"
	"os"
  "udplevis"
  "levis"
	"fmt"
	"time"
)

func main() {
  if len(os.Args) < 2 {
		os.Exit(1)
	}

	payload := ""

	if len(os.Args) == 2 {
		payload = ""
	}else {
		payload = os.Args[2]
	}
	req := levis.Message{
		Type:						levis.Confirmable,
		Code:						levis.GET,
		MessageID:			255,
		ContentFormat:	levis.TextPlain,
		Payload:				[]byte(payload),
  }

	start := time.Now()
	c, err := udplevis.Dial("udp", os.Args[1])
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
