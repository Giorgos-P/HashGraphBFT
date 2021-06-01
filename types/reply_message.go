package types

import (
	"HashGraphBFT/logger"
	"HashGraphBFT/variables"
	"bytes"
	"encoding/gob"
)

// Reply struct
type Reply struct {
	From  int
	Value int
}

// NewReplyMessage - Creates a new Reply
func NewReplyMessage(value int) Reply {
	return Reply{From: variables.ID, Value: value}
}

// GobEncode - Reply message encoder
func (r Reply) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(r.From)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(r.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Reply message decoder
func (r *Reply) GobDecode(buf []byte) error {
	d := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(d)
	err := decoder.Decode(&r.From)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&r.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}
