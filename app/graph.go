package app

import (
	"HashGraphBFT/logger"
	"HashGraphBFT/threshenc"
	"HashGraphBFT/types"
	"HashGraphBFT/variables"

	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

//History -
type History struct {
	events []*EventNode
}

// EventNode -
type EventNode struct {
	EventMessage      *types.EventMessage
	PreviousEvent     *EventNode
	ParentEvent       *EventNode //gossipParent
	OwnHash           string
	Know              []bool
	HaveInsertedFirst bool
	Round             int
	Witness           bool
	InRow             int //only for debugging
	Famous            bool
}

//HashGraph -
var HashGraph []History

// EventNodes without parent which are in the hashgaph
var orphanParent []*EventNode

// EventNodes without parent which are in the hashgaph
var orphanPrevious []*EventNode

//all the hashgraph events in chronological order
var sortedEvents []*EventNode

// MU - mutex to access hashgraph
var MU sync.Mutex

// MuChannel - Incoming Event channel mutex
var MuChannel sync.Mutex

// MuMessageChannel - Incoming Event channel mutex
var MuMessageChannel sync.Mutex

// InitGraph -
func InitGraph() {

	N := variables.N
	HashGraph = make([]History, N)

	for i := 0; i < N; i++ {
		HashGraph[i].events = make([]*EventNode, 0)
	}

	orphanParent = make([]*EventNode, 0)
	orphanPrevious = make([]*EventNode, 0)

	insertFirst()

	transactionCount := 2
	makeTra(transactionCount)

	transactionsCreated = ((transactionCount * variables.N) * variables.N) + variables.N
	go SendGossip()
	go ManageIncomingGossip()

}

func makeTra(num int) {
	for i := 0; i < num; i++ {
		character := rune('b' + i)
		tra := fmt.Sprintf("%c", character)
		makeTransaction(tra)
	}
}

func insertOrderedEvent(eventNode *EventNode) {
	//sortedEvents

	atTheEnd := true
	insertAtIndex := -1
	for i, v := range sortedEvents {
		if isBefore(eventNode, v) {
			atTheEnd = false
			insertAtIndex = i
			break
		}
	}

	if atTheEnd {
		sortedEvents = append(sortedEvents, eventNode)
	} else {
		temp := append([]*EventNode{}, sortedEvents[insertAtIndex:]...) // sub array from index till end
		sortedEvents = append(sortedEvents[0:insertAtIndex], eventNode) // sub array from zero to index-1 (inclusive) and append EventNode
		sortedEvents = append(sortedEvents, temp...)                    //merge the 2 subarrays
	}
}

func isBefore(newEventNode *EventNode, old *EventNode) bool {
	if newEventNode.EventMessage.Timestamp < old.EventMessage.Timestamp {
		return true
	} else if newEventNode.EventMessage.Timestamp == old.EventMessage.Timestamp {
		return newEventNode.EventMessage.Owner < old.EventMessage.Owner
	}
	return false
}

func insertFirst() {

	for i := 0; i < variables.N; i++ {
		eventMessage := new(types.EventMessage)
		eventMessage.Timestamp = 0
		eventMessage.Transaction = strconv.Itoa(i) + "a"
		eventMessage.PreviousHash = "0"
		eventMessage.ParentHash = "0"
		eventMessage.Owner = i
		eventMessage.Signature = []byte(strconv.Itoa(i))

		eventNode := newEventNode(eventMessage, nil, nil, true)
		for i := 0; i < variables.N; i++ {
			eventNode.Know[i] = true
		}
		k := HashGraph[i].events
		k = append(k, eventNode)

		HashGraph[i].events = k
		sortedEvents = append(sortedEvents, eventNode)
	}

	SortSlice(sortedEvents)
}

func makeTransaction(transactionLetter string) {
	ID := variables.ID
	// if ID >= 4 {
	// 	return
	// }

	k := HashGraph[ID].events
	myLastElement := k[len(k)-1]

	eventMessage := new(types.EventMessage)
	//eventMessage.Signature = strconv.Itoa(ID)
	eventMessage.Timestamp = time.Now().UnixNano()
	eventMessage.Transaction = strconv.Itoa(ID) + transactionLetter
	eventMessage.PreviousHash = myLastElement.OwnHash
	eventMessage.ParentHash = "0"
	eventMessage.Owner = ID
	CreateSignature(eventMessage)

	eventNode := newEventNode(eventMessage, myLastElement, nil, true)
	eventNode.Know[ID] = true

	k = append(k, eventNode)

	HashGraph[ID].events = k
	insertOrderedEvent(eventNode)
}

func makeHash(ev *types.EventMessage) string {
	return MD5(fmt.Sprintf("%#v", ev))
}

func newEventNode(ev *types.EventMessage, PreviousEvent *EventNode, ParentEvent *EventNode,
	HaveInsertedFirst bool) *EventNode {

	eventNode := new(EventNode)
	eventNode.EventMessage = ev
	eventNode.PreviousEvent = PreviousEvent
	eventNode.ParentEvent = ParentEvent
	eventNode.Know = make([]bool, variables.N)
	for i := 0; i < variables.N; i++ {
		eventNode.Know[i] = false
	}

	//eventNode.HaveInsertedFirst = HaveInsertedFirst
	eventNode.OwnHash = makeHash(ev)
	return eventNode
}

func newEventNodeInMyRow(parentEvent *EventNode) *EventNode {

	ID := variables.ID
	k := HashGraph[variables.ID].events
	myLastElement := k[len(k)-1]

	eventMessage := new(types.EventMessage)
	//eventMessage.Signature = strconv.Itoa(ID)
	eventMessage.Timestamp = time.Now().UnixNano()
	eventMessage.Transaction = parentEvent.EventMessage.Transaction
	eventMessage.PreviousHash = myLastElement.OwnHash
	eventMessage.ParentHash = parentEvent.OwnHash
	eventMessage.Owner = ID
	CreateSignature(eventMessage)

	eventNode := newEventNode(eventMessage, myLastElement, parentEvent, true)
	eventNode.Know[ID] = true

	k = append(k, eventNode)
	HashGraph[ID].events = k
	insertOrderedEvent(eventNode)
	return eventNode
}

var show bool = false
var displayCount int = 0
var transactionsCreated int = 4

func showHashGraph(writeToLog bool) {

	var events []*EventNode
	events = make([]*EventNode, 0, 10)

	for i := 0; i < variables.N; i++ {
		k := HashGraph[i].events
		if len(k) == 0 {
			continue
		}
		for _, v := range k {
			v.InRow = i
			events = append(events, v)

		}
	}
	SortSlice(events)

	// if len(events) < 24 {
	// 	return
	// }

	// if displayCount == len(events) {
	// 	return
	// } else {
	// 	displayCount = len(events)
	// 	//fmt.Println(len(events))
	// }

	total := transactionsCreated

	if (len(events) == total) && (sentToAllEverything()) {
	} else {
		return
	}

	//DivideRounds()
	//DecideFame()

	//showWit()

	transactionsCreated = -1
	writeToLog = true
	if writeToLog {
		logger.HashGraphLogger.Println("Display HashGraph in time order")
		logger.HashGraphLogger.Println("# of Elements : ", len(events))

		for _, v := range events {
			//logger.HashGraphLogger.Println("---", v.OwnHash, v.Know, v.HaveInsertedFirst)
			logger.HashGraphLogger.Println(v.Round, v.Witness, v.Famous)

			//logger.HashGraphLogger.Println(v.EventMessage.Signature, v.EventMessage.Timestamp, v.EventMessage.Transaction, v.EventMessage.PreviousHash, v.EventMessage.ParentHash, v.EventMessage.Owner)
			logger.HashGraphLogger.Println(v.EventMessage.Timestamp, v.EventMessage.Transaction, v.EventMessage.PreviousHash, v.EventMessage.ParentHash, v.EventMessage.Owner)

			// if v.PreviousEvent == nil {
			// 	logger.HashGraphLogger.Println(v.EventMessage.Transaction)
			// } else {

			// 	logger.HashGraphLogger.Print(v.InRow, " - ", v.EventMessage.Transaction, " - ", v.PreviousEvent.EventMessage.Transaction, " - ")
			// 	if v.ParentEvent != nil {
			// 		logger.HashGraphLogger.Println(v.ParentEvent.InRow, " - ", v.ParentEvent.EventMessage.Transaction)
			// 	} else {
			// 		logger.HashGraphLogger.Println(" - No parent")
			// 	}
			// }
		}
	} else {
		fmt.Println("Display HashGraph in time order")
		fmt.Println("# of Elements : ", len(events))

		for _, v := range events {
			fmt.Println("---", v.OwnHash, v.Know, v.HaveInsertedFirst)
			fmt.Println(v.EventMessage)
		}

	}

	go checkExit()
}

func checkExit() {
	// logger.HashGraphLogger.Println("len(MessageChannel)")

	// logger.HashGraphLogger.Println(len(MessageChannel))
	time.Sleep(time.Second * 10)
	if len(MessageChannel) == 0 && len(EventChannel) == 0 {
		os.Exit(0)
	}
}

// VerifySignature -
func VerifySignature(ev *types.EventMessage) bool {

	//sign := ev.Signature
	k := *ev
	k.Signature = []byte("0")
	msg := fmt.Sprintf("%#v", k)
	ans := threshenc.VerifyMessage([]byte(msg), ev.Signature, ev.Owner)
	return ans
}

// CreateSignature -
func CreateSignature(ev *types.EventMessage) {

	//sign := ev.Signature
	k := *ev
	k.Signature = []byte("0")
	msg := fmt.Sprintf("%#v", k)
	ev.Signature = threshenc.SignMessage([]byte(msg))

}

// // ByTimestamp
// // type ByTimestamp []*EventNode

// // func (a ByTimestamp) Len() int           { return len(a) }
// // func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// // func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }

// func sorting() {

// 	for i := 0; i < variables.N; i++ {
// 		k := HashGraph[i].events
// 		//sort.Sort(ByTimestamp(k))

// 		sort.Slice(k[:], func(i, j int) bool {
// 			return k[i].EventMessage.Timestamp < k[j].EventMessage.Timestamp
// 		})

// 	}

// }

// SortSlice -
func SortSlice(events []*EventNode) {
	sort.Slice(events[:], func(i, j int) bool {
		if events[i].EventMessage.Timestamp < events[j].EventMessage.Timestamp {
			return true
		} else if events[i].EventMessage.Timestamp == events[j].EventMessage.Timestamp {
			return events[i].EventMessage.Owner < events[j].EventMessage.Owner
		}
		return false
	})
}

var changes bool = true

// SendGossip -
func SendGossip() {
	// if variables.ID >= 2 {
	// 	return
	// }

	rand.Seed(time.Now().UnixNano())
	var syncWith int
	for {
		syncWith = rand.Intn(variables.N)

		for syncWith == variables.ID {
			syncWith = rand.Intn(variables.N)
		}

		count := 0

		var events []*EventNode
		events = make([]*EventNode, 0, 10)

		for _, v := range sortedEvents {
			if v.Know[syncWith] == false {
				events = append(events, v)
			}
		}

		// for i := 0; i < variables.N; i++ {
		// 	k := HashGraph[i].events
		// 	count += len(k)
		// 	if len(k) == 0 {
		// 		continue
		// 	}
		// 	for _, v := range k {
		// 		if v.Know[syncWith] == false {

		// 			events = append(events, v)
		// 		}
		// 	}
		// }

		if len(events) == 0 {
			if changes {
				fmt.Println("------------------------")
				fmt.Println("No HashGrpah Changes, Event# = ", count)
				fmt.Println("Ordered Events# = ", len(sortedEvents))

				changes = false
			}
			continue
		}

		//SortSlice(events)

		for _, v := range events {
			SendEvent(v.EventMessage, syncWith)
			v.Know[syncWith] = true

			fmt.Println("------------------------\n")
			fmt.Println("Sending")
			fmt.Println("Send ", v.EventMessage)
			fmt.Println("To ", syncWith)

		}

		showHashGraph(false)
		changes = true
	}

}

func sentToAllEverything() bool {
	for i := 0; i < variables.N; i++ {
		k := HashGraph[i].events
		if len(k) == 0 {
			continue
		}
		for _, v := range k {
			for syncWith := 0; syncWith < 4; syncWith++ {
				if v.Know[syncWith] == false {
					return false
				}
			}
		}
	}
	return true
}

//ManageIncomingGossip -
func ManageIncomingGossip() {

	for {
		MuChannel.Lock()
		count := len(EventChannel)

		if count > 0 {
			gossip := <-EventChannel

			MuChannel.Unlock()
			from := gossip.From
			msg := gossip.Ev
			checkGossip(msg, from)
			showHashGraph(true)
		} else {
			MuChannel.Unlock()

		}

	}

}

func checkGossip(ev *types.EventMessage, from int) {
	verify := VerifySignature(ev) // verify message signature
	if !verify {
		//logger.HashGraphLogger.Println("Not Verify")
		return
	}
	// logger.HashGraphLogger.Println("\n\n")

	// logger.HashGraphLogger.Println("--------------------------")
	// logger.HashGraphLogger.Println("Receice msg from ", from)
	// logger.HashGraphLogger.Println(*ev)

	//exists, existEventNode := checkIfExistsInHashGraph(ev, from)
	exists, _ := checkIfExistsInHashGraph(ev, from)
	//logger.HashGraphLogger.Println("Exists ", exists)

	if exists == 1 {
		return
	}
	//else it does not exists in hashGraph

	//it is not nill - always has a previous
	previousEvent := getPreviousEvent(ev)
	// if previousEvent == nil {
	// 	return
	// }

	orphan := false
	parentEvent, err := getParentEvent(ev)
	if err != nil { // if we didnt have the parent it will be inserted in orphaned list
		orphan = true
	}

	// logger.HashGraphLogger.Println("Orphan ", orphan)
	// logger.HashGraphLogger.Println("New Event ")

	eventNode := insertIncomingGossip(ev, previousEvent, parentEvent)
	//logger.HashGraphLogger.Println(eventNode)

	if previousEvent == nil {
		orphanPrevious = append(orphanPrevious, eventNode)
	}

	if orphan {
		orphanParent = append(orphanParent, eventNode)
	}

	findOrphanParent(eventNode)
	findOrphanPrevious(eventNode)
	//checkIfItIsANewTransaction(eventNode, from)
	checkIfItIsANewTransaction2(eventNode, from)

}

func findOrphanPrevious(eventNode *EventNode) {
	var toRemove []*EventNode

	toRemove = make([]*EventNode, 0)
	for _, v := range orphanPrevious {
		if v.EventMessage.PreviousHash == eventNode.OwnHash {
			v.PreviousEvent = eventNode
			toRemove = append(toRemove, v)
		}
	}
	removeEventNodePrevious(toRemove)
}

// RemoveEventNode -
func removeEventNodePrevious(toRemove []*EventNode) {
	for _, v := range toRemove {
		for i, k := range orphanPrevious {
			if v == k {
				orphanPrevious = RemoveIndex(orphanPrevious, i)
				break
			}
		}
	}
}

func findOrphanParent(eventNode *EventNode) {
	var toRemove []*EventNode

	toRemove = make([]*EventNode, 0)
	for _, v := range orphanParent {
		if v.EventMessage.ParentHash == eventNode.OwnHash {
			v.ParentEvent = eventNode
			toRemove = append(toRemove, v)
		}
	}
	RemoveEventNode(toRemove)
}

// RemoveEventNode -
func RemoveEventNode(toRemove []*EventNode) {
	for _, v := range toRemove {
		for i, k := range orphanParent {
			if v == k {
				orphanParent = RemoveIndex(orphanParent, i)
				break
			}
		}
	}
}

// RemoveIndex -
func RemoveIndex(s []*EventNode, index int) []*EventNode {
	return append(s[:index], s[index+1:]...)
}

// Check if i have the message
//returns true if i have it in the hashmap and sets that sender knows that message
// returns 0 = not exists
// returns 1 = exists in hashgraph
func checkIfExistsInHashGraph(ev *types.EventMessage, from int) (int, *EventNode) {
	hash := makeHash(ev)
	var exists = false
	var eventNode *EventNode = nil

Loop:
	for i := 0; i < variables.N; i++ {
		k := HashGraph[i].events
		if len(k) == 0 {
			continue
		}
		for _, v := range k {
			if v.OwnHash == hash {
				exists = true
				eventNode = v
				break Loop
			}
		}
	}

	if exists {
		eventNode.Know[from] = true
		return 1, eventNode
	} else {
		return 0, nil

	}

}

// find previous event
// the event in owner map
func getPreviousEvent(ev *types.EventMessage) *EventNode {

	k := HashGraph[ev.Owner].events
	if len(k) == 0 {
		return nil
	}

	for _, v := range k {
		if v.OwnHash == ev.PreviousHash {
			return v
		}
	}
	return nil
}

// find parent event
func getParentEvent(ev *types.EventMessage) (*EventNode, error) {

	if ev.ParentHash == "0" {
		return nil, nil
	}

	for i := 0; i < variables.N; i++ {
		if i == ev.Owner {
			continue
		}
		k := HashGraph[i].events
		if len(k) == 0 {
			continue
		}

		for _, v := range k {
			if v.OwnHash == ev.ParentHash {
				return v, nil
			}
		}
	}

	return nil, errors.New("Parent not exists")
}

func insertIncomingGossip(ev *types.EventMessage, previousEvent *EventNode, parentEvent *EventNode) *EventNode {
	eventNode := newEventNode(ev, previousEvent, parentEvent, false)
	eventNode.Know[variables.ID] = true
	insertEventNodeInHashGraph(eventNode)
	return eventNode
}

func insertEventNodeInHashGraph(eventNode *EventNode) {
	OwnerID := eventNode.EventMessage.Owner
	k := HashGraph[OwnerID].events
	k = append(k, eventNode)
	HashGraph[OwnerID].events = k
	insertOrderedEvent(eventNode)
}

func checkIfItIsANewTransaction(eventNode *EventNode, from int) {

	if eventNode.EventMessage.Owner != from {
		return
	}

	haveInserted := CheckHaveInsertedFirst(eventNode)
	if haveInserted {
		return
	}
	//	logger.HashGraphLogger.Println("Inserting also new Transaction")

	newEvent := newEventNodeInMyRow(eventNode)
	abc(newEvent)
}

func abc(eventNode *EventNode) {
	for eventNode != nil {
		eventNode.HaveInsertedFirst = true
		eventNode = eventNode.ParentEvent
	}
}

func checkIfItIsANewTransaction2(eventNode *EventNode, from int) {
	//logger.HashGraphLogger.Println("Check if its a new Transaction")

	if haveTransaction(eventNode) {
		return
	}
	//logger.HashGraphLogger.Println("Insert New Event in my row")
	newEventNodeInMyRow(eventNode)
	//ev := newEventNodeInMyRow(eventNode)
	//logger.HashGraphLogger.Println(ev, "\n\n")

	//abc(newEvent)
}

func haveTransaction(eventNode *EventNode) bool {
	k := HashGraph[variables.ID].events
	for _, v := range k {
		if v.EventMessage.Transaction == eventNode.EventMessage.Transaction {
			return true
		}
	}
	return false
}

//CheckHaveInsertedFirst -
// follow parent events and check i the first event of this transaction(=Event without parent)
// check haveInsertetFirst
// that shows if i have insert this transaction in my row on graph
// returns if i have insert it
// makes HaveInsertedFirst variable in all eventNodes in the path true
// because i will insert it or it is already true
func CheckHaveInsertedFirst(eventNode *EventNode) bool {
	var haveInserted bool
	if eventNode.ParentEvent != nil {
		haveInserted = CheckHaveInsertedFirst(eventNode.ParentEvent)
	} else { //eventNode.ParentEvent == nil
		haveInserted = eventNode.HaveInsertedFirst
	}

	eventNode.HaveInsertedFirst = true
	return haveInserted
}
