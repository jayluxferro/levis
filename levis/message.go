package levis

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// LevisType represents the message type
type LevisType uint8

const (
	// Confirmable message; requires acknowledgement.
	Confirmable LevisType = 0

	// Non-confirmable message; doesn't require acknowledgement.
	NonConfirmable LevisType = 1

	// Acknowledgement; indicates a response to a confirmable message.
	Acknowledgement LevisType = 2

	// Reset; indicates a permanent negative acknowledgement.
	// Receiver is unable to process a non-confirmable message.
	Reset LevisType = 3

	// ResponseTimeout is the amount of time to wait for a response.
	ResponseTimeout = time.Second * 2 // 2 seconds

	// ResponseRandomFactor is a multiplier for response backoff.
	ResponseRandomFactor = 1.5

	// MaxRetransmit is the maximum number of times a message will be retransmitted.
	MaxRetransmit = 4

	// Maximum Packet Len
	MaxPktLen = 1500

	// Default Header Size
	DefaultHeaderSize = 2 // bytes

	// ProtocolName
	ProtocolName = "levis"
)

var typeNames = [4]string{
	Confirmable:     "Confirmable",
	NonConfirmable:  "NonConfirmable",
	Acknowledgement: "Acknowledgement",
	Reset:           "Reset",
}

func init() {
	for i := range typeNames {
		if typeNames[i] == "" {
			typeNames[i] = fmt.Sprintf("Unknown (0x%x)", i)
		}
	}
}

func (t LevisType) String() string {
	return typeNames[t]
}

// LevisCode represents the REST request/response codes
type LevisCode uint8

const (
	// Request/Response Codes
	GET										LevisCode = 0
	POST									LevisCode = 1
	BadRequest            LevisCode = 2
	ServiceUnavailable    LevisCode = 3
)

var codeNames = [16]string{
	GET:                   "GET",
	POST:                  "POST",
	BadRequest:            "BadRequest",
	ServiceUnavailable:    "ServiceUnavailable",
}

func init() {
	for i := range codeNames {
		if codeNames[i] == "" {
			codeNames[i] = fmt.Sprintf("Unknown (0x%x)", i)
		}
	}
}

func (c LevisCode) String() string {
	return codeNames[c]
}

// MediaType specifies the content type of a message.
type MediaType uint8

// Content types.
const (
	AppXML        MediaType = 0  // application/xml
	AppJSON       MediaType = 1  // application/json
)


// Message is a Levis message
type Message struct {
	Type           LevisType
	Code           LevisCode
	ContentFormat	 MediaType
	MessageID      uint8
	Payload				 []byte
}


// MarshalBinary produces the binary form of the Message
func (m *Message) MarshalBinary() ([]byte, error) {
	tmpbuf := []byte{0, 0}
	binary.BigEndian.PutUint16(tmpbuf, uint16(m.MessageID))

	/*** Message Format ***
	 0                   1                   2                   3
	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|Ver| T |  TKL  |      Code     |          Message ID           |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|   Options (if any) ...
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|1 1 1 1 1 1 1 1|    Payload (if any) ...
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	***/

	buf := bytes.Buffer{}
	buf.Write([]byte{
		(1 << 6) | (uint8(m.Type) << 4) | (uint8(m.Code) << 2) | uint8(m.ContentFormat),
		tmpbuf[1], // last 8 bits for Message ID
	})

	buf.Write(m.Payload)

	return buf.Bytes(), nil
}

// Check if message is confirmable; boolean
func (m Message) IsConfirmable() bool {
	return m.Type == Confirmable
}

// ParseMessage extracts the Message from the buffer
func ParseMessage(data []byte) (Message, error) {
	rv := Message{}
	return rv, rv.UnmarshalBinary(data)
}

// UnmarshalBinary parses the given binary slice as a Message
func (m *Message) UnmarshalBinary(data []byte) error {
	// fmt.Println(string(data))
	if len(data) < DefaultHeaderSize {
		return errors.New("short packet")
	}

	// check version match; 1
	if data[0]>>6 != 1 {
		return errors.New("invalid version")
	}

	m.Type = LevisType((data[0] >> 4) & 0x3)

	m.Code = LevisCode((data[0] >> 2 ) & 0x3)
	m.ContentFormat = MediaType(data[0] & 0x3)
	m.MessageID = uint8(data[1])

	m.Payload = data[DefaultHeaderSize:]
	return nil
}

// TCPMessage is a Levis Message that can encode itself for TCP
// transport.
type TCPMessage struct {
	Message
}
