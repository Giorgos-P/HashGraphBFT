package app

import (
	"HashGraphBFT/logger"
	"HashGraphBFT/variables"
)

var witnessChangeRound int = 0
var witnessChange bool = true

// DivideRounds -
func DivideRounds() {

	countRec := 0

	hasChange := false
	for i, v := range sortedEvents {
		if i >= end {
			break
		}
		if hasChange {
			DivideRoundEvent(v, sortedEvents)
			v.Orderedplace = i
			countRec += 1
		} else {
			if i == v.Orderedplace {
				continue
			} else {
				hasChange = true
				v.Orderedplace = i
				DivideRoundEvent(v, sortedEvents)
				countRec += 1
			}
		} //else hasChange
	} // for range

	logger.InfoLogger.Println("Recursions:", countRec)

}

func insertFirstWitness(eventNode *EventNode) {
	var events []*EventNode
	events = make([]*EventNode, 0, variables.N)
	events = append(events, eventNode)
	Witnesses[eventNode.Round] = events
}

// DivideRoundEvent -
func DivideRoundEvent(eventNode *EventNode, events []*EventNode) {

	if eventNode.PreviousEvent == nil && eventNode.EventMessage.Timestamp == 0 {
		eventNode.Round = 1
		eventNode.Witness = true
		round := 1
		_, ok := Witnesses[round]
		if !ok {
			insertFirstWitness(eventNode)
			witnessChangeRound = 1
			witnessChange = true
		} else {
			checkInsertWitness(eventNode)
			witnessChangeRound = 1
			witnessChange = true
		} // else !ok
		return
	}

	if isOrphan(eventNode) {
		eventNode.Round = -1
		eventNode.Orderedplace = -1
		return
	}

	round := eventNode.PreviousEvent.Round

	if eventNode.ParentEvent != nil && eventNode.ParentEvent.Round > round {
		round = eventNode.ParentEvent.Round
	}
	eventNode.Round = round

	numStrongSee := 0
	limit := variables.T // threshold 2N/3

	for _, v := range events {
		if numStrongSee >= limit {
			break
		}
		if v.Round == round && StronglySee(eventNode, v) {
			numStrongSee++
		}
	} //for

	if numStrongSee >= limit {
		round++
		if eventNode.Witness { // if it is witness in the previous round we remove it - it is in the next Round now
			ev, ok := Witnesses[round-1]
			if ok {
				place := -1
				for i, v := range ev {
					if v.OwnHash == eventNode.OwnHash {
						place = i
						break
					}

				}
				if place > -1 {
					ev = removeElementAtIndex(ev, place)
					eventNode.Witness = false
					Witnesses[round-1] = ev
				}
				witnessChange = true
				witnessChangeRound = min(witnessChangeRound, round-1)

			} //if ok
		}
	}
	eventNode.Round = round

	var witness bool

	witness = eventNode.PreviousEvent == nil || round > eventNode.PreviousEvent.Round

	eventNode.Witness = witness

	if witness {
		_, ok := Witnesses[round]
		if !ok {
			insertFirstWitness(eventNode)
			if witnessChange {
				witnessChangeRound = min(witnessChangeRound, round)
			} else {
				witnessChangeRound = round
			}
			witnessChange = true
		} else {
			checkInsertWitness(eventNode)
			if witnessChange {
				witnessChangeRound = min(witnessChangeRound, round)
			} else {
				witnessChangeRound = round
			}
			witnessChange = true
		} //else !ok

	}
}

func checkInsertWitness(eventNode *EventNode) {
	ev := Witnesses[eventNode.Round]

	ev = append(ev, eventNode)
	min := eventNode.EventMessage.Timestamp
	place := len(ev) - 1
	node := eventNode

	for i, v := range ev {
		if v.EventMessage.Owner == eventNode.EventMessage.Owner && v.EventMessage.Timestamp < min {
			v.Witness = true
			min = v.EventMessage.Timestamp
			place = i
			node = v
		}
	}

	for i := 0; i < len(ev); i++ {
		//for i, v := range ev {
		v := ev[i]
		if v.EventMessage.Owner == eventNode.EventMessage.Owner && i != place {
			v.Witness = false
			v.Famous = false
			ev = append(ev[:i], ev[i+1:]...)
			i--
			place--
		}
	}
	node.Witness = true
	Witnesses[eventNode.Round] = ev

}
