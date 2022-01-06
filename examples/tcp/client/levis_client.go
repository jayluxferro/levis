package main

import (
	"levis"
	"log"
	"tcplevis"
  "os"
	"fmt"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	
	payload := ""

	if len(os.Args) == 2 {
		payload = "{'data': ''}"
	}else {
		payload = "{'data': '" + os.Args[2] + "'}"
	}


	req := levis.Message{
		Type:						levis.Confirmable,
		Code:						levis.GET,
		MessageID:			255,
		ContentFormat:	levis.AppJSON,
		Payload:				[]byte(payload),
	}

	start := time.Now()
	c, err := tcplevis.Dial("tcp", os.Args[1])
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
