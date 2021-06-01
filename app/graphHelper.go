package app

import (
	"HashGraphBFT/threshenc"
	"HashGraphBFT/types"
	"HashGraphBFT/variables"

	"errors"
	"sort"
	"time"
)

func newEventNode(ev *types.EventMessage, PreviousEvent *EventNode, ParentEvent *EventNode) *EventNode {

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

func SortSlice(events []*EventNode) {
	sort.Slice(events[:], func(i, j int) bool {
		return isBefore(events[i], events[j])

		// if events[i].EventMessage.Timestamp < events[j].EventMessage.Timestamp {
		// 	return true
		// } else if events[i].EventMessage.Timestamp == events[j].EventMessage.Timestamp {
		// 	return events[i].EventMessage.Owner < events[j].EventMessage.Owner
		// }
		// return false

	})
}

//ExistInMyGraph -
//Checks the client message if it exists in my graph
// I assume that client sent its request to one node
// otherwise i have to check all the graph if i have that node
// same transaction & same ID/Number & same Client Number
// func ExistInMyGraph(transactionLetter string, transactionNumber int, clientID int) bool {
// 	k := HashGraph[variables.ID].events

// 	for _, v := range k {
// 		if v.EventMessage.Transaction == transactionLetter && v.EventMessage.Number == transactionNumber && v.EventMessage.ClientID == clientID {
// 			return true
// 		}
// 	}
// 	return false
// }

func haveClientTransaction(transactionLetter string, transactionNumber int, clientID int) bool {

	for _, v := range sortedEvents {
		if v.EventMessage.Transaction == transactionLetter && v.EventMessage.Number == transactionNumber && v.EventMessage.ClientID == clientID {
			return true
		}
	}
	return false
}

func checkIfItIsANewTransaction(eventNode *EventNode) {

	if haveTransaction(eventNode) {
		// transaction := eventNode.EventMessage.Transaction
		// if !strings.HasPrefix(transaction, "sync") {

		// 	newEventNodeInMyRowSync(eventNode)
		// }

		return
	}
	newEventNodeInMyRow(eventNode)

}

func haveTransaction(eventNode *EventNode) bool {
	k := HashGraph[variables.ID].events
	for _, v := range k {
		if v.EventMessage.Transaction == eventNode.EventMessage.Transaction && v.EventMessage.Number == eventNode.EventMessage.Number && v.EventMessage.ClientID == eventNode.EventMessage.ClientID && v.EventMessage.FirstOwner == eventNode.EventMessage.FirstOwner {
			return true
		}
	}
	return false
}

func haveOwnerTransaction(eventMessage *types.EventMessage) bool {
	k := HashGraph[eventMessage.Owner].events
	for _, v := range k {
		if v.EventMessage.Transaction == eventMessage.Transaction && v.EventMessage.Number == eventMessage.Number && v.EventMessage.ClientID == eventMessage.ClientID {

			if v.EventMessage.Timestamp < eventMessage.Timestamp {
				changedHash := v.OwnHash

				v.EventMessage.Timestamp = eventMessage.Timestamp
				v.EventMessage.Signature = []byte("0")
				threshenc.CreateSignature(v.EventMessage)
				v.OwnHash = makeHash(v.EventMessage)
				for i := 0; i < variables.N; i++ {
					v.Know[i] = false
				}
				v.Know[variables.ID] = true

				for _, ev := range sortedEvents {
					if ev.EventMessage.PreviousHash == changedHash {
						ev.EventMessage.PreviousHash = v.OwnHash
					}
					if ev.EventMessage.ParentHash == changedHash {
						ev.EventMessage.ParentHash = v.OwnHash
					}
				}

			}

			return true
		}
	}
	return false
}

// Check if i have the message
//returns true if i have it in the hashgraph and sets that sender knows that message
// returns nil = not exists
// returns the eventNode if exists in hashgraph
func checkIfExistsInHashGraph(ev *types.EventMessage, from int) bool { //*EventNode {
	//	hash := makeHash(ev)
	var exists = false
	var eventNode *EventNode = nil

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

	//***************
	for _, v := range sortedEvents {
		if v.EventMessage.Timestamp == ev.Timestamp && v.EventMessage.Owner == ev.Owner && v.EventMessage.Transaction == ev.Transaction && v.EventMessage.Number == ev.Number && v.EventMessage.ClientID == ev.ClientID && v.EventMessage.FirstOwner == ev.FirstOwner {
			exists = true
			eventNode = v
		}
	}

	// event := checkSameTransaction(ev)
	// if event != nil {
	// 	if event.EventMessage.Timestamp == ev.Timestamp {
	// 		exists = true
	// 		eventNode = event
	// 	}
	// }
	//***************

	if exists {
		eventNode.Know[from] = true
		return true
		//return eventNode
	} else {
		return false
		//return nil
	}

}

//InsertSorderedEvent -
func insertSorderedEvent(eventNode *EventNode) {
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

// get the previous Event in the same row
// return nil if there isn't one
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
// if err == nil
// 		- it does not have a parent or
//		- it have a parent and i found it in the graph
// if err != nil - it have a parent and i did not found it in the graph
func getParentEvent(ev *types.EventMessage) (*EventNode, error) {

	//it does not have a parent
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

	//if i did not find its parent
	return nil, errors.New("Parent not exists")
}

var syncNum int = 0

// func newEventNodeInMyRowSync(parentEvent *EventNode) *EventNode {

// 	transaction := fmt.Sprintf("sync%d-%d", variables.ID, syncNum)
// 	syncNum += 1
// 	//transaction := fmt.Sprintf("sync")

// 	ID := variables.ID
// 	k := HashGraph[variables.ID].events
// 	myLastElement := k[len(k)-1]

// 	tempEventMessage := types.NewEventMessage([]byte("0"), time.Now().UnixNano(), transaction, myLastElement.OwnHash, parentEvent.OwnHash, ID, parentEvent.EventMessage.Number, parentEvent.EventMessage.ClientID)
// 	eventMessage := &tempEventMessage
// 	threshenc.CreateSignature(eventMessage)

// 	eventNode := insertIncomingGossip(eventMessage, myLastElement, parentEvent)

// 	return eventNode
// }

// create event in my row when i receive a new transaction i did not know
func newEventNodeInMyRow(parentEvent *EventNode) *EventNode {

	ID := variables.ID
	k := HashGraph[variables.ID].events
	myLastElement := k[len(k)-1]

	tempEventMessage := types.NewEventMessage([]byte("0"), time.Now().UnixNano(), parentEvent.EventMessage.Transaction, myLastElement.OwnHash, parentEvent.OwnHash, ID, parentEvent.EventMessage.Number, parentEvent.EventMessage.ClientID, parentEvent.EventMessage.FirstOwner)
	eventMessage := &tempEventMessage
	threshenc.CreateSignature(eventMessage)

	eventNode := insertIncomingGossip(eventMessage, myLastElement, parentEvent)

	return eventNode
}

func insertIncomingGossip(ev *types.EventMessage, previousEvent *EventNode, parentEvent *EventNode) *EventNode {
	eventNode := newEventNode(ev, previousEvent, parentEvent)
	eventNode.Know[variables.ID] = true // i know that i insert it in my graph
	insertEventNodeInHashGraph(eventNode)
	return eventNode
}

func insertEventNodeInHashGraph(eventNode *EventNode) {
	OwnerID := eventNode.EventMessage.Owner
	k := HashGraph[OwnerID].events
	k = append(k, eventNode)
	HashGraph[OwnerID].events = k
	insertSorderedEvent(eventNode)
}

func insertNotOrphansInGraph() {

	var toRemove []*OrphanType
	toRemove = make([]*OrphanType, 0)

	for _, msg := range orphans {
		orphan := msg.Orphan
		hasParent := false
		hasPrevious := false
		for _, v := range sortedEvents {

			if v.OwnHash == orphan.ParentHash {
				hasParent = true
			}
			if v.OwnHash == orphan.PreviousHash {
				hasPrevious = true
			}

		}

		if hasPrevious && hasParent {
			toRemove = append(toRemove, msg)
		}

	}
	RemoveEventMsg(toRemove, orphans)

	for _, notOrphan := range toRemove {
		EventChannel <- struct {
			Ev   *types.EventMessage
			From int
		}{Ev: notOrphan.Orphan, From: notOrphan.From}
	}
}

func canOrder(eventNode *EventNode) bool {
	return eventNode.ParentEvent == nil && !isOrphan(eventNode)
}

func findOrphanParent(eventNode *EventNode) {
	findOrphan(eventNode, orphanParent, 0)
}

func findOrphanPrevious(eventNode *EventNode) {
	findOrphan(eventNode, orphanPrevious, 1)
}

func isOrphan(eventNode *EventNode) bool {
	if eventNode.PreviousEvent == nil {
		return true
	}

	if eventNode.ParentEvent == nil {
		for _, v := range orphanParent {
			if v.OwnHash == eventNode.OwnHash {
				return true
			}
		}
	}

	return false
}

func findOrphan(eventNode *EventNode, list []*EventNode, choice int) {
	var toRemove []*EventNode
	toRemove = make([]*EventNode, 0)

	for _, v := range list {

		if choice == 0 {
			if v.EventMessage.ParentHash == eventNode.OwnHash {
				v.ParentEvent = eventNode
				toRemove = append(toRemove, v)
			}
		} else {
			if v.EventMessage.PreviousHash == eventNode.OwnHash {
				v.PreviousEvent = eventNode
				toRemove = append(toRemove, v)
			}
		}

	}
	RemoveEventNode(toRemove, list)
}

// RemoveEventNode -
func RemoveEventNode(toRemove []*EventNode, list []*EventNode) {
	for i := 0; i < len(list); i++ {
		url := list[i]
		for _, rem := range toRemove {
			if url == rem {
				list = append(list[:i], list[i+1:]...)
				i-- // Important: decrease index
				break
			}
		}
	}
}

func removeElementAtIndex(list []*EventNode, index int) []*EventNode {
	list = append(list[:index], list[index+1:]...)
	return list
}

func RemoveEventMsg(toRemove []*OrphanType, list []*OrphanType) {
	for i := 0; i < len(list); i++ {
		url := list[i]
		for _, rem := range toRemove {
			if url == rem {
				list = append(list[:i], list[i+1:]...)
				i-- // Important: decrease index
				break
			}
		}
	}
}

func removeElementAtIndexEventMsg(list []*OrphanType, index int) []*OrphanType {
	list = append(list[:index], list[index+1:]...)
	return list
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func max64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// // insert Transaction in my graph
func InsertClientTransactionInGraph(transactionLetter string, transactionNumber int, clientID int) bool {

	if haveClientTransaction(transactionLetter, transactionNumber, clientID) {
		return false
	}

	ID := variables.ID

	k := HashGraph[ID].events
	myLastElement := k[len(k)-1]

	tempEventMessage := types.NewEventMessage([]byte("0"), time.Now().UnixNano(), transactionLetter, myLastElement.OwnHash, "0", ID, transactionNumber, clientID, variables.ID)

	eventMessage := &tempEventMessage

	threshenc.CreateSignature(eventMessage)

	eventNode := newEventNode(eventMessage, myLastElement, nil)
	eventNode.Know[ID] = true

	k = append(k, eventNode)
	HashGraph[ID].events = k

	insertSorderedEvent(eventNode)

	return true
}
