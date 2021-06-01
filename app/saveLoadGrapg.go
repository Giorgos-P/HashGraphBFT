package app

import (
	"HashGraphBFT/types"
	"HashGraphBFT/variables"
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"

	"os"
)

// // EventMessage - the gossip message
// type EventMessage struct {
// 	Signature    []byte
// 	Timestamp    int64
// 	Transaction  string
// 	PreviousHash string
// 	ParentHash   string
// 	Owner        int
// }

// // EventNode -
// type EventNode struct {
// 	EventMessage      *types.EventMessage
// 	PreviousEvent     *EventNode
// 	ParentEvent       *EventNode //gossipParent
// 	OwnHash           string
// 	Know              []bool
// 	HaveInsertedFirst bool
// 	Round             int
// 	Witness           bool
// 	InRow             int //only for debugging
// 	Famous            bool
// }

func LoadGraph() {

	N := variables.N
	HashGraph = make([]History, N)
	allEvents := make([]*EventNode, 0)

	for i := 0; i < N; i++ {
		HashGraph[i].events = make([]*EventNode, 0)
	}

	file, err := os.Open("graph2.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		eventMessage := new(types.EventMessage)
		eventMessage.Signature = []byte{'a'}
		eventMessage.Timestamp, _ = strconv.ParseInt(s[0], 10, 64)
		eventMessage.Transaction = s[1]
		eventMessage.PreviousHash = s[2]
		eventMessage.ParentHash = s[3]
		eventMessage.Owner, _ = strconv.Atoi(s[4])

		eventNode := new(EventNode)

		parent := false
		previous := false

		eventNode.EventMessage = eventMessage
		if eventMessage.PreviousHash == "0" {
			eventNode.PreviousEvent = nil
			previous = true
		}
		if eventMessage.ParentHash == "0" {
			eventNode.ParentEvent = nil
			parent = true
		}
		eventNode.OwnHash = s[5]

		for _, v := range allEvents {
			if !parent && v.OwnHash == eventMessage.ParentHash {
				eventNode.ParentEvent = v
				parent = true
			}
			if !previous && v.OwnHash == eventMessage.PreviousHash {
				eventNode.PreviousEvent = v
				previous = true
			}
		}

		allEvents = append(allEvents, eventNode)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK")

	sortedEvents = allEvents
	DivideRounds()
	//showWit()
	os.Exit(0)
}
