package types

import (
	"HashGraphBFT/logger"
	"bytes"
	"container/list"
	"encoding/gob"
	"reflect"
)

const (
	PRE_PREP int = 0
	PREP     int = 1
	COMMIT   int = 2
)

// PendReqs - List of pending requests
type PendReqs *list.List

// ReqQ - List of request status
type ReqQ *list.List

// ------------------------------------------------------------------------------------ //

// LogTuple - Data structure used in rLog
type LogTuple struct {
	// Request TODO [] : Check between AcceptedRequest and Request
	Req *AcceptedRequest
	// TODO: Type to be defined
	XSet []int
}

// GobEncode - LogTuple encoder
func (lt *LogTuple) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(lt.Req)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(lt.XSet)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - LogTuple decoder
func (lt *LogTuple) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&lt.Req)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&lt.XSet)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}

// Equals - Checks if LogTuples are equal
func (lt *LogTuple) Equals(tp *LogTuple) bool {
	if len(lt.XSet) != len(tp.XSet) {
		return false
	}
	for i := range lt.XSet {
		if lt.XSet[i] != tp.XSet[i] {
			return false
		}
	}
	return lt.Req.Equals(tp.Req)
}

// ------------------------------------------------------------------------------------ //

// RequestStatus - Struct request status
type RequestStatus struct {
	Req    *AcceptedRequest
	Status int
}

// GobEncode - RequestStatus encoder
func (rs *RequestStatus) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(rs.Req)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rs.Status)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - RequestStatus decoder
func (rs *RequestStatus) GobDecode(buf []byte) error {
	read := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(read)
	err := decoder.Decode(&rs.Req)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&rs.Status)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}

// Equals - RequestStatus equals check
func (rs *RequestStatus) Equals(r *RequestStatus) bool {
	return rs.Req.Equals(r.Req) && rs.Status == r.Status
}

// ------------------------------------------------------------------------------------ //

// ReplicaStructure - TODO Check for capacities
type ReplicaStructure struct {
	RepState    RepState
	RLog        []*LogTuple
	PendReqs    *list.List //Queue Of Request
	ReqQ        *list.List //Queue Of RequestStatus
	LastReq     []*RequestReply
	ConFlag     bool
	ViewChanged bool
	Prim        int
}

// GobEncode - ReplicaStructure Encoder
func (rs *ReplicaStructure) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(rs.RepState)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rs.RLog)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	var pendReqs []*Request
	for e := rs.PendReqs.Front(); e != nil; e = e.Next() {
		pendReqs = append(pendReqs, e.Value.(*Request))
	}
	err = encoder.Encode(pendReqs)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	var reqQ []*RequestStatus
	for e := rs.ReqQ.Front(); e != nil; e = e.Next() {
		reqQ = append(reqQ, e.Value.(*RequestStatus))
	}
	err = encoder.Encode(reqQ)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rs.LastReq)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rs.ConFlag)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rs.ViewChanged)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rs.Prim)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - ReplicaStructure Decoder
func (rs *ReplicaStructure) GobDecode(buf []byte) error {
	read := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(read)
	err := decoder.Decode(&rs.RepState)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&rs.RLog)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	var pendReqs []*Request
	err = decoder.Decode(&pendReqs)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	logger.OutLogger.Println("Received |pendReqs|=", len(pendReqs))
	rs.PendReqs = list.New()
	for _, req := range pendReqs {
		rs.PendReqs.PushBack(req)
	}
	var reqQ []*RequestStatus
	err = decoder.Decode(&reqQ)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	logger.OutLogger.Println("Received |reqQ|=", len(reqQ))

	rs.ReqQ = list.New()
	for _, req := range reqQ {
		rs.ReqQ.PushBack(req)
	}
	err = decoder.Decode(&rs.LastReq)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&rs.ConFlag)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&rs.ViewChanged)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&rs.Prim)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}

// Equals - ReplicaStructure equals check
func (rs *ReplicaStructure) Equals(repl *ReplicaStructure) bool {
	return rs.RepState.Equals(repl.RepState) &&
		RLogEquals(rs.RLog, repl.RLog) &&
		ListEquals(rs.PendReqs, repl.PendReqs) &&
		ListEquals(rs.ReqQ, repl.ReqQ) &&
		LastReqEquals(rs.LastReq, repl.LastReq) &&
		rs.ConFlag == repl.ConFlag &&
		rs.ViewChanged == repl.ViewChanged &&
		rs.Prim == repl.Prim
}

// ------------------------------------------------------------------------------------ //

// LastExec - Returns the last executed replica
func (rs *ReplicaStructure) LastExec() *LogTuple {
	if len(rs.RLog) < 1 {
		return nil
	}
	last := rs.RLog[0]
	for _, element := range rs.RLog {
		if element.Req.Sq > last.Req.Sq {
			last = element
		}
	}
	return last
}

// TODO: Apply this to all sections necessary.
func ContainsRequest(x *Request, arr []*Request) bool {
	for _, req := range arr {
		if req.Equals(x) {
			return true
		}
	}
	return false
}

func ContainsAcceptedRequest(x *AcceptedRequest, arr []*AcceptedRequest) bool {
	for _, req := range arr {
		if req.Equals(x) {
			return true
		}
	}
	return false
}

func GetRequestsListFromLog(rLog []*LogTuple) *list.List {
	reqs := list.New()
	for _, log := range rLog {
		reqs.PushBack(log.Req)
	}
	return reqs
}

func GetRequestsFromLog(rLog []*LogTuple) []*AcceptedRequest {
	var reqs []*AcceptedRequest
	for _, log := range rLog {
		reqs = append(reqs, log.Req)
	}
	return reqs
}

/*************** OPERATORS **************/

/*
Enqueue - adds an element (or set of elements) x to a queue.
If any element enqueued already exists, then only the most
recent copy of it is kept and it is carried back to the
queue.
*/
func (rs *ReplicaStructure) Enqueue(el ...interface{}) {
	if len(el) == 0 {
		return
	}
	elType := reflect.TypeOf(el[0]).String()
	logger.OutLogger.Println("Enqueue Called with ", elType, ".")
	switch elType {
	case "*types.Request":
		for _, element := range el {
			element := element.(*Request)
			for e := rs.PendReqs.Front(); e != nil; e = e.Next() {
				if e.Value.(*Request).Equals(element) {
					rs.PendReqs.Remove(e)
				}
			}
			rs.PendReqs.PushBack(element)
		}
		logger.OutLogger.Println("Enqueued Request in PendReqs")
		break
	case "*types.RequestStatus":
		for _, element := range el {
			element := element.(*RequestStatus)
			for e := rs.ReqQ.Front(); e != nil; e = e.Next() {
				if e.Value.(*RequestStatus).Equals(element) {
					rs.ReqQ.Remove(e)
				}
			}
			rs.ReqQ.PushBack(element)
		}
		logger.OutLogger.Println("Enqueued RequestStatus in ReqQ")
		break
	case "*types.RequestReply":
		for _, element := range el {
			element := element.(*RequestReply)
			for i, reqRep := range rs.LastReq {
				if element.Equals(reqRep) {
					rs.LastReq = append(rs.LastReq[:i], rs.LastReq[i+1:]...)
				}
			}
			rs.LastReq = append(rs.LastReq, element)
		}
		logger.OutLogger.Println("Enqueued RequestReply in LastReq")

		break
	}
}

// Remove - removes element x from a structure
func (rs *ReplicaStructure) Remove(el ...interface{}) {
	if len(el) == 0 {
		return
	}
	elType := reflect.TypeOf(el[0]).String()
	logger.OutLogger.Println("Remove Called with ", elType)
	switch elType {
	case "*types.Request":
		for _, element := range el {
			element := element.(*Request)
			for e := rs.PendReqs.Front(); e != nil; e = e.Next() {
				if e.Value.(*Request).Equals(element) {
					rs.PendReqs.Remove(e)
				}
			}
		}

	case "*types.AcceptedRequest":
		for _, element := range el {
			element := element.(*AcceptedRequest)
			count := 0
			for e := rs.ReqQ.Front(); e != nil; e = e.Next() {
				if e.Value.(*RequestStatus).Req.Request.Equals(element.Request) {
					count++
					rs.ReqQ.Remove(e)
				}
			}
		}
	case "*types.RequestReply":
		for _, element := range el {
			element := element.(*RequestReply)
			for i, reqRep := range rs.LastReq {
				if element.Equals(reqRep) {
					rs.LastReq = append(rs.LastReq[:i], rs.LastReq[i+1:]...)
				}
			}
		}
	}
}

// Add - adds element x to a structure
func (rs *ReplicaStructure) Add(el ...interface{}) {
	if len(el) == 0 {
		return
	}
	elType := reflect.TypeOf(el[0]).String()
	logger.OutLogger.Println("Add Called with ", elType)
	switch elType {
	case "*types.Request":
		for _, element := range el {
			element := element.(*Request)
			rs.PendReqs.PushBack(element)
		}
		break
	case "*types.RequestStatus":
		for _, element := range el {
			element := element.(*RequestStatus)
			val := new(RequestStatus)
			val.Status = element.Status
			//*val = *element
			val.Req = &AcceptedRequest{
				Request: &Request{
					Client:    element.Req.Request.Client,
					TimeStamp: element.Req.Request.TimeStamp,
					Operation: element.Req.Request.Operation},
				View: element.Req.View,
				Sq:   element.Req.Sq,
			}
			for e := rs.ReqQ.Front(); e != nil; e = e.Next() {
				if e.Value.(*RequestStatus).Equals(val) {
					return
				}
			}
			rs.ReqQ.PushBack(val)

		}
		break
	case "*types.AcceptedRequest":
		for _, element := range el {
			element := element.(*AcceptedRequest)
			rs.ReqQ.PushBack(&RequestStatus{Req: element, Status: PRE_PREP})
		}
		break
	case "*types.RequestReply":
		for _, element := range el {
			element := element.(*RequestReply)
			rs.LastReq = append(rs.LastReq, element)
		}
		break
	case "LogTuple":
		for _, element := range el {
			element := element.(*LogTuple)
			rs.RLog = append(rs.RLog, element)
		}

	}
}
