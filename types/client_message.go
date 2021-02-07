package types

import (
	"HashGraphBFT/logger"
	"bytes"
	"encoding/gob"
)

// ClientMessage - Client message struct
type ClientMessage struct {
	Req *Request
	Ack bool
}

// GobEncode - Client message encoder
func (cm *ClientMessage) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(cm.Req)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(cm.Ack)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Client message decoder
func (cm *ClientMessage) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&cm.Req)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&cm.Ack)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}

// Equals - Checks if client messages are equal
func (cm *ClientMessage) Equals(cmsg *ClientMessage) bool {
	return cm.Req.Equals(cmsg.Req) && cm.Ack == cmsg.Ack
}
