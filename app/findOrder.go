package app

import (
	"HashGraphBFT/logger"
	"HashGraphBFT/types"
	"HashGraphBFT/variables"

	"fmt"
	"math"
)

var OrderedTransactions []string = make([]string, 0)
var OrderedEvents []*EventNode = make([]*EventNode, 0)
var tempOrder []*EventNode = make([]*EventNode, 0)
var maxOrderedRoundReceive int = 2

// find the max Round number which has witness
func maxRound() int {
	max := -1
	for i := 1; ; i++ {
		_, exist := Witnesses[i]
		if exist {
			max = i
		} else {
			break
		}
	}
	return max
}

func Ord() {
	maxRoundWitness := maxRound()

	maxDesiredRound := maxRoundWitness - 3 //inclusive
	if maxDesiredRound < 2 {
		return
	}

	for _, v := range OrderedEvents {
		for _, event := range sortedEvents {
			if sameTransaction(v, event) {
				event.RoundReceived = v.RoundReceived
			}
		}
		for _, event := range tempOrder {
			if sameTransaction(v, event) {
				event.RoundReceived = v.RoundReceived
			}
		}
	}

	for _, v := range tempOrder {
		for _, event := range sortedEvents {
			if sameTransaction(v, event) {
				event.RoundReceived = v.RoundReceived
			}
		}
	}

	for key, currentEvent := range sortedEvents {
		if key >= end {
			break
		}

		whoKnows := 0
		for _, know := range currentEvent.Know {
			if know {
				whoKnows++
			}
		}
		if whoKnows < variables.T {
			continue
		}

		//if its an empty transaction or one of the firsts
		if currentEvent.EventMessage.Timestamp == 0 || currentEvent.Round == -1 { // || emptyTransaction(currentEvent) {
			continue
		}
		//if its a recent transaction or one that already been ordered
		if currentEvent.Round > maxDesiredRound || currentEvent.RoundReceived != 0 || emptyTransaction(currentEvent) {
			continue
		}

		if !seeParent(currentEvent) {
			continue
		}
		if !canOrder(currentEvent) {
			continue
		}
		eventRound := currentEvent.Round + 1

		for {

			currentWitnesses, ok := Witnesses[eventRound]
			if eventRound > maxDesiredRound {
				break
			}

			if !ok {
				//fmt.Println("Empty wintess") // Normally this cant happen
				break
			}
			countSee := 0
			countFamous := 0
			for _, witness := range currentWitnesses {
				visited = make([]*EventNode, 0)
				if witness.Famous {
					countFamous++
				} else {
					continue
				}
				if see(witness, currentEvent, FAME) {
					countSee++
				}
			}

			if countSee == countFamous && countFamous > 0 {
				setRoundRec(currentEvent, eventRound, key)

				break
			} else if countFamous == 0 {
				break
			}
			eventRound++

		}

	}
	SortOrderedEvents(tempOrder)

	currentMaxReceivedRound := -1
	maxPlace := -1

	for _, v := range tempOrder {
		if v.RoundReceived > currentMaxReceivedRound {
			currentMaxReceivedRound = v.RoundReceived
		}
	}
	if currentMaxReceivedRound > maxOrderedRoundReceive {
		maxOrderedRoundReceive = currentMaxReceivedRound
	}

	for i, v := range tempOrder {
		exist := false
		for _, orderEvent := range OrderedEvents {
			if sameTransaction(v, orderEvent) {
				fmt.Println("Same transaction ", v.EventMessage.Transaction, v.EventMessage.Number)
				exist = true
			}
		}

		if v.RoundReceived < maxOrderedRoundReceive {
			if i > maxPlace {
				maxPlace = i
			}
			if !exist {
				OrderedEvents = append(OrderedEvents, v)
			}
			if !emptyTransaction(v) && !exist {
				logger.InfoLogger.Println(v.RoundReceived, " ", v.Round, " T:", v.EventMessage.Transaction, " Client:", v.EventMessage.ClientID, " Num:", v.EventMessage.Number)
				//logger.OrderLogger.Println(v.RoundReceived, " ", v.EventMessage.ClientID, " ", v.EventMessage.Transaction, " ", v.EventMessage.Number, " ", v.OwnHash, " ", v.ConsensusTime, " ", v.Round)
				logger.OrderLogger.Println(v.EventMessage.ClientID, " ", v.EventMessage.Transaction, " ", v.EventMessage.Number)

				reply := types.NewReplyMessage(v.EventMessage.Number)

				go ReplyClient(reply, v.EventMessage.ClientID)
				//ReplyClient(reply, v.EventMessage.ClientID)

			}
		}
	}

	if maxPlace >= 0 {
		tempOrder = tempOrder[maxPlace+1:]
	}

}

func setRoundRec(event *EventNode, roundReceived int, key int) {
	event.RoundReceived = roundReceived

	times := make([]int64, 0)

	for _, v := range sortedEvents {
		if sameTransaction(event, v) {
			v.RoundReceived = roundReceived

			if v.Round >= event.Round && v.Round < roundReceived {

				visited = make([]*EventNode, 0)
				if see(v, event, FAME) {
					times = append(times, v.EventMessage.Timestamp)

				}
			}
		}

	}

	eventsNum := len(times)
	median := int(math.Ceil(float64(eventsNum)/2.0) - 1) // take the midean timestamp

	event.ConsensusTime = times[median]
	tempOrder = append(tempOrder, event)
}

var secondTime bool = false

func clearAll() {
	secondTime = true
	tempOrder = make([]*EventNode, 0)
	for _, v := range sortedEvents {
		v.RoundReceived = 0
		v.ConsensusTime = 0
	}
	OrderedEvents = make([]*EventNode, 0)

}
