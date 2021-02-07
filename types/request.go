package types

import (
	"HashGraphBFT/logger"
	"bytes"
	"encoding/gob"
	"time"
)

// ------------------------------------------------------------------------------------ //

// Operation struct
type Operation struct {
	Op    int
	Value rune
}

// GobEncode - Operation encoder
func (op Operation) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(op.Op)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(op.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Operation decoder
func (op *Operation) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&op.Op)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&op.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}

// Equals - Checks if operations are equal
func (op *Operation) Equals(operation Operation) bool {
	return op.Value == operation.Value && op.Op == operation.Op
}

// ------------------------------------------------------------------------------------ //

// Request struct
type Request struct {
	Client    int
	TimeStamp time.Time
	Operation Operation
}

// GobEncode - Request encoder
func (r *Request) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(r.Client)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(r.TimeStamp)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(r.Operation)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Request decoder
func (r *Request) GobDecode(buf []byte) error {
	read := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(read)
	err := decoder.Decode(&r.Client)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&r.TimeStamp)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&r.Operation)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}

// Equals - Checks if requests are equal
func (r *Request) Equals(request *Request) bool {
	return r.Client == request.Client && r.Operation.Equals(request.Operation) &&
		r.TimeStamp.Equal(request.TimeStamp)
}

// MergeRequestSets - Merges two requests into one set
func MergeRequestSets(r []*Request, q []*Request) []*Request {
	var set []*Request
	set = append([]*Request(nil), r...)
	for _, el := range r {
		flag := false
		for _, e := range q {
			if el.Equals(e) {
				flag = true
				break
			}
		}
		if !flag {
			set = append(set, el)
		}
	}
	return set
}

// ------------------------------------------------------------------------------------ //

// Reply struct
type Reply struct {
	View      int
	TimeStamp time.Time
	Client    int
	ID        int
	Result    RepState
}

// Equals - Checks if replies are equal
func (rep *Reply) Equals(reply *Reply) bool {
	return rep.Client == reply.Client &&
		rep.ID == reply.ID && rep.Result.Equals(reply.Result)
}

// ------------------------------------------------------------------------------------ //

// RequestReply struct
type RequestReply struct {
	Req    *Request // Request TODO []: Check between AcceptedRequest and Request
	Client int
	Rep    *Reply
}

// GobEncode - RequestReply encoder
func (rr *RequestReply) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(rr.Req)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rr.Client)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rr.Rep)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - RequestReply decoder
func (rr *RequestReply) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&rr.Req)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&rr.Client)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&rr.Rep)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}

// Equals - Checks if requestreplies are equal
func (rr *RequestReply) Equals(req *RequestReply) bool {
	return rr.Req.Equals(req.Req) && (rr.Rep == nil || rr.Rep.Equals(rr.Rep))
}

// ------------------------------------------------------------------------------------ //

// AcceptedRequest struct
type AcceptedRequest struct {
	Request *Request
	View    int
	Sq      int // Sequence Number
}

// GobEncode - AcceptedRequest encoder
func (r *AcceptedRequest) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(r.Request)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(r.View)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(r.Sq)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), err
}

// GobDecode - AcceptedRequest decoder
func (r *AcceptedRequest) GobDecode(buf []byte) error {
	read := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(read)
	err := decoder.Decode(&r.Request)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&r.View)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&r.Sq)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}

// Equals - Checks if acceptedrequests are equal
func (r *AcceptedRequest) Equals(req *AcceptedRequest) bool {
	//("r.Sq=",r.Sq)
	//("req.Sq=",req.Sq)
	return r.Request.Equals(req.Request) && r.View == req.View &&
		r.Sq == req.Sq
}
