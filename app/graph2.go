package app

// import (
// 	"HashGraphBFT/logger"
// 	"HashGraphBFT/types"
// 	"HashGraphBFT/variables"
// 	"errors"
// 	"fmt"
// 	"math/rand"
// 	"sort"
// 	"strconv"
// 	"sync"
// 	"time"
// )

// //History -
// type History struct {
// 	events []*EventNode
// }

// // EventNode -
// type EventNode struct {
// 	EventMessage      *types.EventMessage
// 	PreviousEvent     *EventNode
// 	ParentEvent       *EventNode
// 	OwnHash           string
// 	Know              []bool
// 	HaveInsertedFirst bool
// }

// //HashGraph -
// var HashGraph []History

// // MU - mutex to access hashgraph
// var MU sync.Mutex

// // MuChannel - Incoming Event channel mutex
// var MuChannel sync.Mutex

// // MuMessageChannel - Incoming Event channel mutex
// var MuMessageChannel sync.Mutex

// // InitGraph -
// func InitGraph() {

// 	N := variables.N
// 	HashGraph = make([]History, N)

// 	for i := 0; i < N; i++ {
// 		HashGraph[i].events = make([]*EventNode, 0)
// 	}
// 	insertFirst()
// 	if variables.ID == 0 {
// 		makeTransaction("b")
// 	}
// 	//makeTransaction("c")
// 	go SendGossip()
// 	go ManageIncomingGossip()

// }

// func insertFirst() {
// 	//ID := variables.ID

// 	for i := 0; i < variables.N; i++ {
// 		eventMessage := new(types.EventMessage)
// 		eventMessage.Signature = strconv.Itoa(i)
// 		eventMessage.Timestamp = 0
// 		eventMessage.Transaction = strconv.Itoa(i) + "a"
// 		eventMessage.PreviousHash = "0"
// 		eventMessage.ParentHash = "0"
// 		eventMessage.Owner = i

// 		eventNode := newEventNode(eventMessage, nil, nil, true)
// 		for i := 0; i < variables.N; i++ {
// 			eventNode.Know[i] = true
// 		}
// 		k := HashGraph[i].events
// 		k = append(k, eventNode)

// 		HashGraph[i].events = k
// 	}

// }

// func makeTransaction(transactionLetter string) {
// 	ID := variables.ID
// 	// if ID >= 4 {
// 	// 	return
// 	// }

// 	k := HashGraph[ID].events
// 	myLastElement := k[len(k)-1]

// 	eventMessage := new(types.EventMessage)
// 	eventMessage.Signature = strconv.Itoa(ID)
// 	eventMessage.Timestamp = time.Now().UnixNano()
// 	eventMessage.Transaction = strconv.Itoa(ID) + transactionLetter
// 	eventMessage.PreviousHash = myLastElement.OwnHash
// 	eventMessage.ParentHash = "0"
// 	eventMessage.Owner = ID

// 	eventNode := newEventNode(eventMessage, myLastElement, nil, true)
// 	eventNode.Know[ID] = true

// 	k = append(k, eventNode)

// 	HashGraph[ID].events = k
// }

// func makeHash(ev *types.EventMessage) string {
// 	return MD5(fmt.Sprintf("%#v", ev))
// }

// func newEventNode(ev *types.EventMessage, PreviousEvent *EventNode, ParentEvent *EventNode,
// 	HaveInsertedFirst bool) *EventNode {

// 	eventNode := new(EventNode)
// 	eventNode.EventMessage = ev
// 	eventNode.PreviousEvent = PreviousEvent
// 	eventNode.ParentEvent = ParentEvent
// 	eventNode.Know = make([]bool, variables.N)
// 	for i := 0; i < variables.N; i++ {
// 		eventNode.Know[i] = false
// 	}

// 	eventNode.HaveInsertedFirst = HaveInsertedFirst
// 	eventNode.OwnHash = makeHash(ev)
// 	return eventNode
// }

// func newEventNodeInMyRow(parentEvent *EventNode) *EventNode {

// 	ID := variables.ID
// 	k := HashGraph[variables.ID].events
// 	myLastElement := k[len(k)-1]

// 	eventMessage := new(types.EventMessage)
// 	eventMessage.Signature = strconv.Itoa(ID)
// 	eventMessage.Timestamp = time.Now().UnixNano()
// 	eventMessage.Transaction = parentEvent.EventMessage.Transaction
// 	eventMessage.PreviousHash = myLastElement.OwnHash
// 	eventMessage.ParentHash = parentEvent.OwnHash
// 	eventMessage.Owner = ID

// 	eventNode := newEventNode(eventMessage, myLastElement, parentEvent, true)
// 	eventNode.Know[ID] = true

// 	k = append(k, eventNode)
// 	HashGraph[ID].events = k
// 	return eventNode
// }

// var show bool = false
// var displayCount int = 0
// var transactionsCreated int = 8

// func showHashGraph(writeToLog bool) {

// 	var events []*EventNode
// 	events = make([]*EventNode, 0, 10)

// 	for i := 0; i < variables.N; i++ {
// 		k := HashGraph[i].events
// 		if len(k) == 0 {
// 			continue
// 		}
// 		for _, v := range k {
// 			events = append(events, v)

// 		}
// 	}
// 	sortSlice(events)

// 	// if len(events) < 24 {
// 	// 	return
// 	// }

// 	// if displayCount == len(events) {
// 	// 	return
// 	// } else {
// 	// 	displayCount = len(events)
// 	// 	fmt.Println(len(events))
// 	// }

// 	// total := 4 * (transactionsCreated + 1)
// 	// if len(events) != total {
// 	// 	return
// 	// }
// 	// transactionsCreated = -1

// 	if writeToLog {
// 		logger.HashGraphLogger.Println("Display HashGraph in time order")
// 		logger.HashGraphLogger.Println("# of Elements : ", len(events))

// 		for _, v := range events {
// 			logger.HashGraphLogger.Println("---", v.OwnHash, v.Know, v.HaveInsertedFirst)
// 			logger.HashGraphLogger.Println(v.EventMessage)
// 		}
// 	} else {
// 		fmt.Println("Display HashGraph in time order")
// 		fmt.Println("# of Elements : ", len(events))

// 		for _, v := range events {
// 			fmt.Println("---", v.OwnHash, v.Know, v.HaveInsertedFirst)
// 			fmt.Println(v.EventMessage)
// 		}

// 	}
// }

// func verifySignature(ev *types.EventMessage) bool {

// 	//sign := ev.Signature
// 	k := *ev
// 	k.Signature = "0"

// 	return true
// }

// // // ByTimestamp
// // // type ByTimestamp []*EventNode

// // // func (a ByTimestamp) Len() int           { return len(a) }
// // // func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// // // func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }

// func sorting() {

// 	for i := 0; i < variables.N; i++ {
// 		k := HashGraph[i].events
// 		//sort.Sort(ByTimestamp(k))

// 		sort.Slice(k[:], func(i, j int) bool {
// 			return k[i].EventMessage.Timestamp < k[j].EventMessage.Timestamp
// 		})

// 	}

// }

// func sortSlice(events []*EventNode) {
// 	sort.Slice(events[:], func(i, j int) bool {
// 		return events[i].EventMessage.Timestamp < events[j].EventMessage.Timestamp
// 	})
// }

// var changes bool = true

// // SendGossip -
// func SendGossip() {
// 	// if variables.ID >= 2 {
// 	// 	return
// 	// }

// 	rand.Seed(time.Now().UnixNano())
// 	var syncWith int
// 	for {
// 		syncWith = rand.Intn(variables.N)

// 		for syncWith == variables.ID {
// 			syncWith = rand.Intn(variables.N)
// 		}

// 		var events []*EventNode
// 		events = make([]*EventNode, 0, 10)

// 		for i := 0; i < variables.N; i++ {
// 			k := HashGraph[i].events
// 			if len(k) == 0 {
// 				continue
// 			}
// 			for _, v := range k {
// 				if v.Know[syncWith] == false {

// 					events = append(events, v)
// 				}
// 			}
// 		}

// 		if len(events) == 0 {
// 			if changes {
// 				fmt.Println("------------------------")
// 				fmt.Println("No HashGrpah Changes")
// 				changes = false
// 			}
// 			continue
// 		}

// 		sortSlice(events)

// 		for _, v := range events {
// 			SendEvent(v.EventMessage, syncWith)

// 			v.Know[syncWith] = true
// 		}
// 		showHashGraph(false)
// 		changes = true
// 	}

// }

// //ManageIncomingGossip -
// func ManageIncomingGossip() {

// 	for {
// 		MuChannel.Lock()
// 		count := len(EventChannel)

// 		if count > 0 {
// 			gossip := <-EventChannel

// 			MuChannel.Unlock()
// 			from := gossip.From
// 			msg := gossip.Ev
// 			checkGossip(msg, from)
// 			showHashGraph(true)
// 		} else {
// 			MuChannel.Unlock()

// 		}

// 	}

// }

// func checkGossip(ev *types.EventMessage, from int) {
// 	verify := verifySignature(ev) // verify message signature
// 	if !verify {
// 		return
// 	}
// 	logger.HashGraphLogger.Println("\n\n")

// 	logger.HashGraphLogger.Println("--------------------------")
// 	logger.HashGraphLogger.Println("Receice msg from ", from)
// 	logger.HashGraphLogger.Println(*ev)

// 	exists, existEventNode := checkIfExistsInHashGraph(ev, from)
// 	if exists == 2 { //already exists in HashGraph and i have inserted event in my row
// 		return
// 	} else if exists == 1 { //already exists in HashGraph and i DONT have inserted event in my row
// 		checkIfItIsANewTransaction(existEventNode, from)
// 		return
// 	}
// 	//else it does not exists in hashGraph

// 	//it is not nill - always has a previous
// 	previousEvent := getPreviousEvent(ev)
// 	if previousEvent == nil {
// 		return
// 	}

// 	parentEvent, err := getParentEvent(ev)
// 	if err != nil {
// 		return
// 	}

// 	eventNode := insertIncomingGossip(ev, previousEvent, parentEvent)
// 	checkIfItIsANewTransaction(eventNode, from)
// }

// // Check if i have the message
// //returns true if i have it in the hashmap and sets that sender knows that message
// // returns 0 = not exists
// // returns 1 = exists in hashgraph & I DONT have Inserted eventNode in my row
// // returns 2 = exists in hashgraph & I HAVE Inserted eventNode in my row
// func checkIfExistsInHashGraph(ev *types.EventMessage, from int) (int, *EventNode) {
// 	hash := makeHash(ev)
// 	var exists = false
// 	var eventNode *EventNode = nil

// Loop:
// 	for i := 0; i < variables.N; i++ {
// 		k := HashGraph[i].events
// 		if len(k) == 0 {
// 			continue
// 		}
// 		for _, v := range k {
// 			if v.OwnHash == hash {
// 				exists = true
// 				eventNode = v
// 				break Loop
// 			}
// 		}
// 	}

// 	if exists {
// 		eventNode.Know[from] = true
// 		if eventNode.HaveInsertedFirst {
// 			return 2, nil
// 		} else {
// 			return 1, eventNode
// 		}
// 	} else {
// 		return 0, nil

// 	}

// }

// // find previous event
// // the event in owner map
// func getPreviousEvent(ev *types.EventMessage) *EventNode {

// 	k := HashGraph[ev.Owner].events
// 	if len(k) == 0 {
// 		return nil
// 	}

// 	for _, v := range k {
// 		if v.OwnHash == ev.PreviousHash {
// 			return v
// 		}
// 	}
// 	return nil
// }

// // find parent event
// func getParentEvent(ev *types.EventMessage) (*EventNode, error) {

// 	if ev.ParentHash == "0" {
// 		return nil, nil
// 	}

// 	for i := 0; i < variables.N; i++ {
// 		if i == ev.Owner {
// 			continue
// 		}
// 		k := HashGraph[i].events
// 		if len(k) == 0 {
// 			continue
// 		}

// 		for _, v := range k {
// 			if v.OwnHash == ev.ParentHash {
// 				return v, nil
// 			}
// 		}
// 	}

// 	return nil, errors.New("Parent not exists")
// }

// func insertIncomingGossip(ev *types.EventMessage, previousEvent *EventNode, parentEvent *EventNode) *EventNode {
// 	eventNode := newEventNode(ev, previousEvent, parentEvent, false)
// 	eventNode.Know[variables.ID] = true
// 	insertEventNodeInHashGraph(eventNode)
// 	return eventNode
// }

// func insertEventNodeInHashGraph(eventNode *EventNode) {
// 	OwnerID := eventNode.EventMessage.Owner
// 	k := HashGraph[OwnerID].events
// 	k = append(k, eventNode)
// 	HashGraph[OwnerID].events = k
// }

// func checkIfItIsANewTransaction(eventNode *EventNode, from int) {

// 	if eventNode.EventMessage.Owner != from {
// 		return
// 	}

// 	haveInserted := CheckHaveInsertedFirst(eventNode)
// 	if haveInserted {
// 		return
// 	}
// 	logger.HashGraphLogger.Println("Inserting also new Transaction")

// 	newEvent := newEventNodeInMyRow(eventNode)
// 	abc(newEvent)
// }

// func abc(eventNode *EventNode) {
// 	for eventNode != nil {
// 		eventNode.HaveInsertedFirst = true
// 		eventNode = eventNode.ParentEvent
// 	}
// }

// //CheckHaveInsertedFirst -
// // follow parent events and check i the first event of this transaction(=Event without parent)
// // check haveInsertetFirst
// // that shows if i have insert this transaction in my row on graph
// // returns if i have insert it
// // makes HaveInsertedFirst variable in all eventNodes in the path true
// // because i will insert it or it is already true
// func CheckHaveInsertedFirst(eventNode *EventNode) bool {
// 	var haveInserted bool
// 	if eventNode.ParentEvent != nil {
// 		haveInserted = CheckHaveInsertedFirst(eventNode.ParentEvent)
// 	} else { //eventNode.ParentEvent == nil
// 		haveInserted = eventNode.HaveInsertedFirst
// 	}

// 	eventNode.HaveInsertedFirst = true
// 	return haveInserted
// }
