package app

import (
	"HashGraphBFT/variables"
)

// Witnesses -
var Witnesses map[int][]*EventNode

// DivideRounds -
func DivideRounds() {

	Witnesses = make(map[int][]*EventNode, 0)

	graph := HashGraph

	var events []*EventNode
	events = make([]*EventNode, 0, 10)

	for i := 0; i < variables.N; i++ {
		k := graph[i].events
		if len(k) == 0 {
			continue
		}
		for _, v := range k {
			events = append(events, v)
		}
	}
	SortSlice(events)
	for _, v := range events {
		DivideRoundEvent(v, events)
	}
}

// DivideRoundEvent -
func DivideRoundEvent(eventNode *EventNode, events []*EventNode) {

	if eventNode.PreviousEvent == nil {
		eventNode.Round = 1
		eventNode.Witness = true

		round := 1
		ev, ok := Witnesses[round]
		if !ok {
			var events []*EventNode
			events = make([]*EventNode, 0, variables.N)
			events = append(events, eventNode)
			Witnesses[round] = events
		} else {
			ev = append(ev, eventNode)
			Witnesses[round] = ev
		}

		// if len(Witnesses) < round {
		// 	newRoundWitnesses := new(History)
		// 	Witnesses = append(Witnesses, *newRoundWitnesses)
		// 	Witnesses[round-1].events = make([]*EventNode, 0)
		// 	Witnesses[round-1].events = append(Witnesses[round-1].events, eventNode)
		//  } else {
		//  	Witnesses[round-1].events = append(Witnesses[round-1].events, eventNode)
		// }

		return
	}

	round := eventNode.PreviousEvent.Round

	if eventNode.ParentEvent != nil && eventNode.ParentEvent.Round > round {
		round = eventNode.ParentEvent.Round
	}
	eventNode.Round = round

	numStrongSee := 0

	for _, v := range events {
		if eventNode.EventMessage.Timestamp < v.EventMessage.Timestamp {
			continue
		}
		if numStrongSee > 2*variables.N/3 {
			break
		}
		if v.Round == round && StronglySee(eventNode, v) {
			numStrongSee++
		}
	} //for

	if numStrongSee > 2*variables.N/3 {
		round++
	}
	eventNode.Round = round

	var witness bool
	witness = eventNode.PreviousEvent == nil || round > eventNode.PreviousEvent.Round
	eventNode.Witness = witness

	if witness {
		ev, ok := Witnesses[round]
		if !ok {
			var events []*EventNode
			events = make([]*EventNode, 0, variables.N)
			events = append(events, eventNode)
			Witnesses[round] = events
		} else {

			ev = append(ev, eventNode)
			Witnesses[round] = ev
		}

		// 	if len(Witnesses) < round {
		// 		fmt.Println("------------")
		// 		fmt.Println("len(Witnesses) = ", len(Witnesses))
		// 		fmt.Println("round = ", round)
		// 		newRoundWitnesses := new(History)
		// 		Witnesses = append(Witnesses, *newRoundWitnesses)
		// 		fmt.Println("len(Witnesses) = ", len(Witnesses))

		// 		Witnesses[round-1].events = make([]*EventNode, 0)
		// 		Witnesses[round-1].events = append(Witnesses[round-1].events, eventNode)
		// 	} else {
		// 		Witnesses[round-1].events = append(Witnesses[round-1].events, eventNode)
		// 	}
	}
}

var ownerSee []bool
var visited []*EventNode

//StronglySee -
func StronglySee(x *EventNode, y *EventNode) bool {
	ownerSee = make([]bool, variables.N, variables.N)
	visited = make([]*EventNode, 0)

	for i := 0; i < variables.N; i++ {
		ownerSee[i] = false
	}

	canSee := see(x, y)
	if !canSee {
		return false
	}

	count := 0
	for i := 0; i < variables.N; i++ {
		if ownerSee[i] {
			count++
		}
	}

	if count > 2*variables.N/3 {
		return true
	} else {
		return false
	}

}

func see(x *EventNode, y *EventNode) bool {
	if x.OwnHash == y.OwnHash {
		ownerSee[x.EventMessage.Owner] = true
		return true
	}
	if x.Round < y.Round {
		return false
	}

	ownerSee[x.EventMessage.Owner] = true
	var a, b bool = false, false
	if x.ParentEvent != nil && !isVisited(x.ParentEvent) {
		visited = append(visited, x.ParentEvent)
		a = see(x.ParentEvent, y)
	}
	if x.PreviousEvent != nil && !isVisited(x.PreviousEvent) {
		visited = append(visited, x.PreviousEvent)
		b = see(x.PreviousEvent, y)
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
