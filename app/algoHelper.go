package app

import (
	"HashGraphBFT/variables"
	"sort"
	"strings"
)

const ROUND = 0
const FAME = 1

var ownerSee []bool
var visited []*EventNode

//StronglySee -
func StronglySee(x *EventNode, y *EventNode) bool {
	ownerSee = make([]bool, variables.N, variables.N) // count the owners we visited - init to false
	visited = make([]*EventNode, 0)                   // keep visited nodes so we dont visit them twice

	// for i := 0; i < variables.N; i++ {
	// 	ownerSee[i] = false
	// }

	canSee := see(x, y, ROUND)

	if !canSee {
		return false
	}

	count := 0
	for i := 0; i < variables.N; i++ {
		if ownerSee[i] {
			count++
		}
	}

	limit := variables.T

	return count >= limit
}

func see(x *EventNode, y *EventNode, choice int) bool {
	if x.OwnHash == y.OwnHash {
		if choice == ROUND {
			ownerSee[x.EventMessage.Owner] = true
		}
		return true
	}
	if x.Round < y.Round {
		return false
	}

	if choice == ROUND {
		ownerSee[x.EventMessage.Owner] = true
	}

	var a, b bool = false, false
	if x.ParentEvent != nil && !isVisited(x.ParentEvent) {
		visited = append(visited, x.ParentEvent)
		a = see(x.ParentEvent, y, choice)
	}
	if a {
		return a
	}

	if x.PreviousEvent != nil && !isVisited(x.PreviousEvent) {
		visited = append(visited, x.PreviousEvent)
		b = see(x.PreviousEvent, y, choice)
	}

	return a || b
}

func isVisited(x *EventNode) bool {
	for _, v := range visited {
		if x.OwnHash == v.OwnHash {
			return true
		}
	}
	return false
}

func sameTransaction(eventNode1 *EventNode, eventNode2 *EventNode) bool {

	tra1 := eventNode1.EventMessage.Transaction
	tra2 := eventNode2.EventMessage.Transaction
	cl1 := eventNode1.EventMessage.ClientID
	cl2 := eventNode2.EventMessage.ClientID
	num1 := eventNode1.EventMessage.Number
	num2 := eventNode2.EventMessage.Number

	if tra1 == tra2 && cl1 == cl2 && num1 == num2 {
		return true
	}

	return false

}

func emptyTransaction(eventNode *EventNode) bool {
	transaction := eventNode.EventMessage.Transaction

	if strings.HasPrefix(transaction, "empty") || strings.HasPrefix(transaction, "sync") {
		return true
	}

	return false
}

func SortOrderedEvents(events []*EventNode) {
	sort.Slice(events[:], func(i, j int) bool {
		if events[i].RoundReceived < events[j].RoundReceived {
			return true
		} else if events[i].RoundReceived == events[j].RoundReceived {
			return events[i].ConsensusTime < events[j].ConsensusTime
		}

		return false
	})
}

func seeParent(base *EventNode) bool {
	owner := make([]bool, variables.N, variables.N)

	visited = make([]*EventNode, 0) // keep visited nodes so we dont visit them twice

	for _, v := range sortedEvents {
		visited = make([]*EventNode, 0)
		if sameTransactionParent(v, base) {
			visited = make([]*EventNode, 0)
			if seePar(v, base) {
				owner[v.EventMessage.Owner] = true
			}
		}
	}

	count := 0
	for _, v := range owner {

		if v {
			count += 1
		}
	}

	return count >= variables.T

}

func seePar(x *EventNode, y *EventNode) bool {
	if x.OwnHash == y.OwnHash {
		return true
	}
	if x.Round < y.Round {
		return false
	}

	var a bool = false
	if x.ParentEvent != nil && !isVisited(x.ParentEvent) {
		visited = append(visited, x.ParentEvent)
		a = seePar(x.ParentEvent, y)
	}

	return a
}

func sameTransactionParent(x *EventNode, y *EventNode) bool {

	return sameTransaction(x, y) && x.EventMessage.FirstOwner == y.EventMessage.FirstOwner

}
