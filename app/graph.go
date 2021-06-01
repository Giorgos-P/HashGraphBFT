package app

import (
	"HashGraphBFT/config"
	"HashGraphBFT/logger"
	"HashGraphBFT/threshenc"
	"HashGraphBFT/types"
	"HashGraphBFT/variables"
	"math/rand"
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
	EventMessage  *types.EventMessage
	PreviousEvent *EventNode
	ParentEvent   *EventNode //gossipParent
	OwnHash       string
	Know          []bool
	//HaveInsertedFirst bool
	Round         int
	Witness       bool
	Famous        bool
	RoundReceived int
	ConsensusTime int64
	Orderedplace  int
}

type OrphanType struct {
	Orphan *types.EventMessage
	From   int
}

var (
	//HashGraph -
	HashGraph []History

	// EventNodes without parent which are in the hashgaph
	orphanParent []*EventNode

	// EventNodes without parent which are in the hashgaph
	orphanPrevious []*EventNode

	orphans []*OrphanType

	//all the hashgraph events in chronological order
	sortedEvents []*EventNode

	// MU - mutex to access hashgraph
	MU_HashGraph sync.RWMutex

	// MuChannel - Incoming Event channel mutex
	//MuChannel sync.RWMutex

	// MU_OrderedEvents - for the list of ordered events
	//MU_OrderedEvents sync.RWMutex

	// Witnesses -
	Witnesses map[int][]*EventNode
)

// // MuMessageChannel - Incoming Event channel mutex
// var MuMessageChannel sync.Mutex

// InitGraph -
func InitGraph() {

	N := variables.N
	HashGraph = make([]History, N)

	for i := 0; i < N; i++ {
		HashGraph[i].events = make([]*EventNode, 0)
	}

	orphanParent = make([]*EventNode, 0)
	orphanPrevious = make([]*EventNode, 0)
	orphans = make([]*OrphanType, 0)
	sortedEvents = make([]*EventNode, 0)
	Witnesses = make(map[int][]*EventNode, 0)

	insertFirstTransactions()
	owner = make([]bool, variables.N, variables.N) // count the owners we visited

}

var owner []bool

//StartHashGraph -
func StartHashGraph() {

	if config.Scenario == "IDLE" {

		if variables.ID < variables.T {
			go sendGossip()
		}
	} else {
		go sendGossip()
	}

	owner = make([]bool, variables.N, variables.N) // all are false
	//Simulates the client
	//go ClientTransactionCreation()

	go manageClientRequest()
	go ManageIncomingGossip()
}

//The first transactions, needed to init the graph and
//become the parent/previous of the next
//Inserts one empty transaction for each node
func insertFirstTransactions() {

	for i := 0; i < variables.N; i++ {
		transaction := "First" + strconv.Itoa(i)

		eventMessage := types.NewEventMessage([]byte(strconv.Itoa(i)), 0, transaction, "0", "0", i, 0, -1, i)

		eventNode := newEventNode(&eventMessage, nil, nil)
		for i := 0; i < variables.N; i++ { // we dont send this event - all nodes create the exact same
			eventNode.Know[i] = true
		}
		eventNode.Orderedplace = -1

		k := HashGraph[i].events
		k = append(k, eventNode)

		HashGraph[i].events = k
		sortedEvents = append(sortedEvents, eventNode)
		DivideRoundEvent(eventNode, sortedEvents)

	}

	SortSlice(sortedEvents)
}

func sleepScecario() {
	if config.Scenario == "Sleep" {
		if variables.ID >= variables.T {

			waitTillSleep := 15

			sleepTime := 1

			goToSleep(waitTillSleep, sleepTime)
			//goToSleep(waitTillSleep*2, sleepTime)

			//goToSleep(waitTillSleep, sleepTime)
		}
	}
}

func goToSleep(waitTillSleep int, sleepTime int) {
	time.Sleep(time.Duration(waitTillSleep) * time.Second)
	sleepNow = true
	time.Sleep(time.Duration(sleepTime) * time.Second)
	sleepNow = false
}

var sleepNow bool = false

var TotalGossip int = 0
var sendedEvents int = 0

func sendGossip() {
	rand.Seed(time.Now().UnixNano())
	var syncWith int
	N := variables.N
	ID := variables.ID

	go sleepScecario()

	for {
		for sleepNow {
			//wait till sleepNow is false
		}
		//Choose another node to sync
		//Not my self
		syncWith = rand.Intn(N)
		for syncWith == ID {
			syncWith = rand.Intn(N)
		}

		MU_HashGraph.Lock()
		sendedEvents = 0
		sentGossipTo(syncWith)

		if sendedEvents > 0 {
			TotalGossip += sendedEvents
			logger.InfoLogger.Printf("Gossip - %d (%v)", TotalGossip, sentToAllEverything())

		}
		// newSyncWith := rand.Intn(N)
		// for newSyncWith == ID || newSyncWith == syncWith {
		// 	newSyncWith = rand.Intn(N)
		// }
		// sentGossipTo(newSyncWith)

		MU_HashGraph.Unlock()

		showHashGraphSize()
	}

}

var prevHashGraphSize int = -1

func showHashGraphSize() {
	newHashGraphSize := len(sortedEvents)
	if newHashGraphSize > prevHashGraphSize {
		prevHashGraphSize = newHashGraphSize
		logger.InfoLogger.Println("Size:", newHashGraphSize)

		// if newHashGraphSize > 20 && newHashGraphSize < 50 {
		// 	logger.WitnessLogger.Println("Size:", len(sortedEvents[18].EventMessage.Signature))
		// 	logger.WitnessLogger.Println("Size:", len(sortedEvents[18].EventMessage.Signature))
		// 	logger.WitnessLogger.Println(len(sortedEvents[18].EventMessage.PreviousHash))

		// }
		// for _, k := range sortedEvents {

		// 	logger.InfoLogger.Println(k.EventMessage.Timestamp, " ", k.EventMessage.Transaction, " ", k.EventMessage.Owner, " ", k.EventMessage.Number, k.EventMessage.ClientID)
		// 	logger.InfoLogger.Println("\t", k.OwnHash, " ", k.EventMessage.PreviousHash, " ", k.EventMessage.ParentHash)
		// }
		// logger.InfoLogger.Println("--------------------------------------")

		//printHashGraph()
	}

	//logger.InfoLogger.Println(sentToAllEverything())

	// if newHashGraphSize == 16 {
	// 	for _, k := range sortedEvents {

	// 		logger.InfoLogger.Println(k.EventMessage.Timestamp, " ", k.EventMessage.Transaction, " ", k.EventMessage.Owner, " ", k.EventMessage.Number, k.EventMessage.ClientID)
	// 	}
	// 	logger.InfoLogger.Println("--------------------------------------")

	// }

	// for _, v := range sortedEvents {
	// 	for _, k := range sortedEvents {
	// 		if v.EventMessage.Transaction == k.EventMessage.Transaction && v.EventMessage.Owner == k.EventMessage.Owner && v.EventMessage.Number == k.EventMessage.Number && v.EventMessage.ClientID == k.EventMessage.ClientID {
	// 			if v.EventMessage.Timestamp != k.EventMessage.Timestamp {
	// 				logger.InfoLogger.Println(k.EventMessage.Timestamp, " ", v.EventMessage.Timestamp, " ", k.EventMessage.Transaction, " ", k.EventMessage.Owner)

	// 			}
	// 		}
	// 	}
	// }
	//logger.InfoLogger.Println("--------------------------------------")

}

func printHashGraph() {
	for _, k := range sortedEvents {

		logger.InfoLogger.Println(k.EventMessage.Transaction, " ", k.EventMessage.Timestamp, " ", k.EventMessage.Owner, " ")
		logger.InfoLogger.Println("\t ", k.OwnHash, " ", k.EventMessage.PreviousHash, " ", k.EventMessage.ParentHash)

	}
	logger.InfoLogger.Println("--------------------------------------")

}

var sendfaulty bool = false
var numOfFaults int = 0
var totalFauts int = 20

func sentGossipTo(syncWith int) {

	//syncWith = (variables.ID + 1) % variables.N

	myEvents := make([]*types.EventMessage, 0)
	othersEvents := make([]*types.EventMessage, 0)

	for _, v := range sortedEvents {
		if v.Know[syncWith] == false {

			isNormal := config.Scenario == "NORMAL" || (config.Scenario == "IDLE" && variables.ID < variables.T) || config.Scenario == "Sleep"
			if isNormal {
				//SendEvent(v.EventMessage, syncWith)

				if v.EventMessage.Owner == variables.ID {
					myEvents = append(myEvents, v.EventMessage)
				} else {
					othersEvents = append(othersEvents, v.EventMessage)
				}

				sendedEvents += 1
			} else if config.Scenario == "TimestampChange" {
				timestampScenario1(v, syncWith)
			} else if config.Scenario == "SmallTimestamp" {
				timestampScenario2(v, syncWith)
			} else if config.Scenario == "SmallTimestampToSome" {
				timestampScenario3(v, syncWith)
			} else if config.Scenario == "Fork" {
				timestampScenarioFork(v, syncWith)
			}

			v.Know[syncWith] = true

		}
	}

	for _, v := range myEvents {
		SendEvent(v, syncWith)
	}
	for _, v := range othersEvents {
		SendEvent(v, syncWith)
	}

	tempEventMessage := types.NewEventMessage([]byte("0"), 0, "0", "0", "0", -1, -1, -1, -1)
	eventMessage := &tempEventMessage
	threshenc.CreateSignature(eventMessage)
	SendEvent(eventMessage, syncWith)
	//sendedEvents += 1
}

func getMaxTimestampFromPrevEvent(v *EventNode) int64 {
	var timestamp int64 = int64(0)

	if v.PreviousEvent != nil && v.PreviousEvent.EventMessage != nil {
		timestamp = v.PreviousEvent.EventMessage.Timestamp
	}
	if v.ParentEvent != nil && v.ParentEvent.EventMessage != nil {
		timestamp = max64(timestamp, v.ParentEvent.EventMessage.Timestamp)
	}
	return timestamp
}

// func timestampScenario111(v *EventNode, syncWith int) {
// 	isFaulty := variables.ID >= variables.T && v.EventMessage.Owner == variables.ID && syncWith == 0

// 	//if isFaulty && !sendfaulty {
// 	if isFaulty && !sendfaulty {

// 		var tempEventMessage types.EventMessage
// 		var now int64 = time.Now().UnixNano()

// 		maxPrev := getMaxTimestampFromPrevEvent(v)
// 		//maxPrev := int64()
// 		if maxPrev > 0 {

// 			newTimestamp := rand.Int63n(now-maxPrev) + 2
// 			tempEventMessage = types.NewEventMessage([]byte("0"), newTimestamp, v.EventMessage.Transaction, v.EventMessage.PreviousHash, v.EventMessage.ParentHash, v.EventMessage.Owner, v.EventMessage.Number, v.EventMessage.ClientID)
// 		} else {
// 			SendEvent(v.EventMessage, syncWith)
// 			v.Know[syncWith] = true
// 			return
// 		}

// 		eventMessage := &tempEventMessage
// 		threshenc.CreateSignature(eventMessage)
// 		SendEvent(eventMessage, syncWith)
// 		numOfFaults++
// 		if numOfFaults == totalFauts {
// 			sendfaulty = true
// 		}
// 	} else {
// 		SendEvent(v.EventMessage, syncWith)
// 	}

// }

func timestampScenario1(v *EventNode, syncWith int) {
	isFaulty := variables.ID >= variables.T && v.EventMessage.Owner == variables.ID && v.EventMessage.ParentHash == "0" && syncWith < variables.N/2
	//totalFauts = 50
	//if isFaulty && !sendfaulty {
	if isFaulty {

		var tempEventMessage types.EventMessage
		newTimestamp := int64(100 + numOfFaults)
		tempEventMessage = types.NewEventMessage([]byte("0"), newTimestamp, v.EventMessage.Transaction, "0", "0", v.EventMessage.Owner, v.EventMessage.Number, v.EventMessage.ClientID, v.EventMessage.FirstOwner)
		eventMessage := &tempEventMessage

		threshenc.CreateSignature(eventMessage)
		SendEvent(eventMessage, syncWith)
		sendedEvents += 1
		numOfFaults++
		if numOfFaults == totalFauts {
			sendfaulty = true

		}
	} else {
		SendEvent(v.EventMessage, syncWith)
		sendedEvents += 1
	}

}

func timestampScenario2(v *EventNode, syncWith int) {
	isFaulty := variables.ID >= variables.T && v.EventMessage.Owner == variables.ID && v.EventMessage.ParentHash == "0"

	//if isFaulty && !sendfaulty {
	if isFaulty {

		var tempEventMessage types.EventMessage
		newTimestamp := int64(100 + numOfFaults)
		tempEventMessage = types.NewEventMessage([]byte("0"), newTimestamp, v.EventMessage.Transaction, v.EventMessage.PreviousHash, v.EventMessage.ParentHash, v.EventMessage.Owner, v.EventMessage.Number, v.EventMessage.ClientID, v.EventMessage.FirstOwner)
		eventMessage := &tempEventMessage

		threshenc.CreateSignature(eventMessage)
		SendEvent(eventMessage, syncWith)
		numOfFaults++
		if numOfFaults == totalFauts {
			sendfaulty = true
		}
	} else {
		SendEvent(v.EventMessage, syncWith)
	}

}

func timestampScenario3(v *EventNode, syncWith int) {
	isFaulty := variables.ID >= variables.T && v.EventMessage.Owner == variables.ID && v.EventMessage.ParentHash == "0" && syncWith == 0

	//if isFaulty && !sendfaulty {
	if isFaulty {
		var tempEventMessage types.EventMessage
		newTimestamp := int64(100 + numOfFaults)
		tempEventMessage = types.NewEventMessage([]byte("0"), newTimestamp, v.EventMessage.Transaction, v.EventMessage.PreviousHash, v.EventMessage.ParentHash, v.EventMessage.Owner, v.EventMessage.Number, v.EventMessage.ClientID, v.EventMessage.FirstOwner)
		eventMessage := &tempEventMessage

		threshenc.CreateSignature(eventMessage)
		SendEvent(eventMessage, syncWith)
		numOfFaults++
		if numOfFaults == totalFauts {
			sendfaulty = true
		}
	} else {
		SendEvent(v.EventMessage, syncWith)
	}

}

//needs changes
func timestampScenarioFork(v *EventNode, syncWith int) {
	isFaulty := variables.ID >= variables.T && v.EventMessage.Owner == variables.ID && v.EventMessage.ParentHash == "0" && syncWith == 0
	totalFauts = 50
	if isFaulty && !sendfaulty {
		//if isFaulty {

		var tempEventMessage types.EventMessage
		var now int64 = time.Now().UnixNano()

		maxPrev := getMaxTimestampFromPrevEvent(v)
		var newTimestamp int64

		//maxPrev := int64()
		if maxPrev > 0 {
			newTimestamp = rand.Int63n(now-maxPrev-100) + 2
			tempEventMessage = types.NewEventMessage([]byte("0"), newTimestamp, v.EventMessage.Transaction, v.EventMessage.PreviousHash, v.EventMessage.ParentHash, v.EventMessage.Owner, v.EventMessage.Number, v.EventMessage.ClientID, v.EventMessage.FirstOwner)
		} else {
			SendEvent(v.EventMessage, syncWith)
			v.Know[syncWith] = true
			return
		}

		eventMessage := &tempEventMessage
		threshenc.CreateSignature(eventMessage)
		SendEvent(eventMessage, syncWith)
		sendedEvents += 1
		numOfFaults++
		if numOfFaults == totalFauts {
			sendfaulty = true
		}
	} else {
		SendEvent(v.EventMessage, syncWith)
		sendedEvents += 1
	}

}

func manageClientRequest() {

	for msg := range ClientChannel {
		MU_HashGraph.Lock()

		if InsertClientTransactionInGraph(msg.clientTransaction, msg.transactionNumber, msg.clientID) {
			//executeAlgorithms()
		}

		MU_HashGraph.Unlock()
	}
}

var countTrans int = 0

func executeAlgorithms() {

	//countTrans++

	//if countTrans%(variables.N*10) == 0 {
	DivideRounds()
	DecideFame()
	Ord()
	//}
	countAll()

}

func ManageIncomingGossip() {

	//lastSync := -1

	for gossip := range EventChannel {

		from := gossip.From
		msg := gossip.Ev

		MU_HashGraph.Lock()

		if msg.Transaction == "0" {
			owner[from] = true
		} else {
			res := checkGossip(msg, from)
			if res == 0 {
				countTrans++

			}
			// else if res == 2 {
			// 	orphan := new(OrphanType)
			// 	orphan.Orphan = msg
			// 	orphan.From = from

			// 	orphans = append(orphans, orphan)
			// }
		}

		for testSame2() {
		}

		// if len(EventChannel) == 0 {

		// 	countSyncOwners := 0
		// 	for _, v := range owner {
		// 		if v {
		// 			countSyncOwners++
		// 		}
		// 	}

		// 	//if countSyncOwners >= variables.T-1 {

		// 	// for testSame() {

		// 	// }

		// 	testSame()

		// 	owner = make([]bool, variables.N, variables.N) // count the owners we visited
		// 	countTrans = 0
		// 	//countAll()
		// 	//}
		// 	//insertNotOrphansInGraph()
		// }

		MU_HashGraph.Unlock()

	}

}

func checkSameTransaction(ev *types.EventMessage) *EventNode {
	for _, v := range sortedEvents {
		if v.EventMessage.Owner == ev.Owner && v.EventMessage.Transaction == ev.Transaction && v.EventMessage.Number == ev.Number && v.EventMessage.ClientID == ev.ClientID {
			return v
		}
	}
	return nil
}

// func check(ev *types.EventMessage) bool {

// 	event := checkSameTransaction(ev)
// 	if event != nil {
// 		if event.EventMessage.Timestamp != ev.Timestamp {
// 			logger.WitnessLogger.Printf("Owner:%d From:%d", ev.Owner, event.EventMessage.Owner)
// 			return true
// 		}
// 	}
// 	return false

// }

//false ==1
//true ==0
//orphan ==2

func checkGossip(ev *types.EventMessage, from int) int {
	//isFaulty := check(ev)

	// logger.InfoLogger.Println("Receive from ", from)
	// logger.InfoLogger.Println(ev.Timestamp, " ", ev.Owner, " ", ev.PreviousHash, " ", ev.ParentHash)

	if ev.Timestamp > time.Now().UnixNano() {
		//logger.WitnessLogger.Printf("After Now")

		return 1
	}

	verify := threshenc.VerifySignature(ev) // verify message signature
	if !verify {
		//logger.WitnessLogger.Printf("Not verified")

		return 1
	}

	exist := checkIfExistsInHashGraph(ev, from)
	if exist {
		return 1
	}

	// exist = haveOwnerTransaction(ev)
	// if exist {
	// 	return true
	// }

	previousEvent := getPreviousEvent(ev)

	orphan := false
	parentEvent, err := getParentEvent(ev)
	if err != nil { // if we didnt have the parent it will be inserted in orphaned list
		orphan = true
	} else {
		// its not orphan from parent
		// the timestamp has to be after parents
		if parentEvent != nil {
			if !isAfter(ev, parentEvent) {
				return 1
			}
		}
	}

	// its not orphan from previous
	if previousEvent != nil {
		if !isAfter(ev, previousEvent) {
			return 1
		}
	}

	// if previousEvent == nil || orphan {
	// 	//orphans = append(orphans, ev)
	// 	return 2
	// }

	// found := false
	// for _, v := range sortedEvents {
	// 	if v.EventMessage.Transaction == ev.Transaction && v.EventMessage.Owner == ev.Owner && v.EventMessage.Number == ev.Number && v.EventMessage.ClientID == ev.ClientID {
	// 		logger.InfoLogger.Println("Found")
	// 		found = true
	// 		if v.EventMessage.Timestamp != ev.Timestamp {
	// 			logger.InfoLogger.Println(ev.Timestamp, " ", v.EventMessage.Timestamp, " ", v.EventMessage.Transaction, " ", v.EventMessage.Owner, " ", v.EventMessage.Number, v.EventMessage.ClientID)
	// 		} else {
	// 			logger.InfoLogger.Println(v.EventMessage.Timestamp, " ", v.EventMessage.Transaction, " ", v.EventMessage.Owner, " ", v.EventMessage.Number, v.EventMessage.ClientID)
	// 			logger.InfoLogger.Println("-----------------------")

	// 			return false
	// 		}
	// 	}
	// }

	// if found {
	// 	logger.InfoLogger.Println("-----------------------")

	// }

	eventNode := insertIncomingGossip(ev, previousEvent, parentEvent)
	if previousEvent == nil {
		orphanPrevious = append(orphanPrevious, eventNode)
	}
	if orphan {
		orphanParent = append(orphanParent, eventNode)
	}

	findOrphanParent(eventNode)
	findOrphanPrevious(eventNode)

	checkIfItIsANewTransaction(eventNode)
	return 0
}

func isAfter(ev *types.EventMessage, previous *EventNode) bool {

	return ev.Timestamp > previous.EventMessage.Timestamp
}
