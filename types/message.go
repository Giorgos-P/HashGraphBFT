package types

import (
	"HashGraphBFT/logger"
	"HashGraphBFT/variables"
	"bytes"
	"encoding/gob"
)

// Message - The general message struct
type Message struct {
	Payload []byte
	Type    string
	From    int
}

// EventMessage - the gossip message
type EventMessage struct {
	Signature    []byte
	Timestamp    int64
	Transaction  string
	PreviousHash string
	ParentHash   string
	Owner        int // the node which created this / also shows in which row it is in the hashGraph
	Number       int // increased number
	ClientID     int // which client sent me this transaction
	FirstOwner   int
}

// NewMessage - Creates a new payload message
func NewMessage(payload []byte, Type string) Message {
	return Message{Payload: payload, Type: Type, From: variables.ID}
}

// NewEventMessage - Creates a new event message
// func NewEventMessage(Signature []byte, Timestamp int64, Transaction string, PreviousHash string, ParentHash string, Owner int, Number int, ClientID int) EventMessage {
// 	return EventMessage{Signature: Signature, Timestamp: Timestamp, Transaction: Transaction, PreviousHash: PreviousHash, ParentHash: ParentHash, Owner: Owner, Number: Number, ClientID: ClientID}
// }
func NewEventMessage(Signature []byte, Timestamp int64, Transaction string, PreviousHash string, ParentHash string, Owner int, Number int, ClientID int, first int) EventMessage {
	return EventMessage{Signature: Signature, Timestamp: Timestamp, Transaction: Transaction, PreviousHash: PreviousHash, ParentHash: ParentHash, Owner: Owner, Number: Number, ClientID: ClientID, FirstOwner: first}
}

// GobEncode - Message encoder
func (m Message) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(m.Payload)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(m.Type)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(m.From)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	return w.Bytes(), nil
}

// GobDecode - Message decoder
func (m *Message) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&m.Payload)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&m.Type)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&m.From)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	return nil
}

// GobEncodeEvent - EventMessage encoder
func (m EventMessage) GobEncodeEvent() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(m.Signature)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(m.Timestamp)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(m.Transaction)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(m.PreviousHash)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(m.ParentHash)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(m.Owner)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(m.Number)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(m.ClientID)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(m.FirstOwner)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecodeEvent - Message decoder
func (m *EventMessage) GobDecodeEvent(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&m.Signature)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&m.Timestamp)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&m.Transaction)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&m.PreviousHash)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&m.ParentHash)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&m.Owner)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&m.Number)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&m.ClientID)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&m.FirstOwner)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}
