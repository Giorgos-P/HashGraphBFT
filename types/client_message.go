package types

import (
	"HashGraphBFT/logger"
	"bytes"
	"encoding/gob"
)

type ClientMsg struct {
	ClientTransaction string
	TransactionNumber int
}

// GobEncode - Client message encoder
func (cm *ClientMsg) GobEncodeMsg() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(cm.ClientTransaction)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(cm.TransactionNumber)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Client message decoder
func (cm *ClientMsg) GobDecodeMsg(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&cm.ClientTransaction)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&cm.TransactionNumber)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}
