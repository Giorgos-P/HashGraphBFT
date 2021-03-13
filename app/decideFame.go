package app

import (
	"fmt"
)

// import (
// 	"HashGraphBFT/variables"
// 	"fmt"
// )

// func showWitnesses() {
// 	for i := 0; i < len(Witnesses); i++ {
// 		currentWitnesses := Witnesses[i].events
// 		for _, v := range currentWitnesses {
// 			fmt.Println(v.EventMessage.Transaction, " ", v.Round, " ", v.Witness, " ", v.Famous, " ")
// 		}
// 	}

// }

func showWit() {
	// keys := make([]int, 0)
	// for k, _ := range Witnesses {
	// 	keys = append(keys, k)
	// }
	// sort.Ints(keys)
	// for _, k := range keys {
	// 	fmt.Println(k, Witnesses[k])
	// }
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	for k, _ := range Witnesses {
		fmt.Println(k)
		for _, v := range Witnesses[k] {
			fmt.Println(v.EventMessage.Transaction, " ", v.Round, " ", v.Witness, " ", v.Famous, " ")
		}
	}
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

}

// func DecideFame() {

// 	for i := 0; i < len(Witnesses)-1; i++ {
// 		currentWitnesses := Witnesses[i].events
// 		if len(currentWitnesses) == 0 {
// 			continue
// 		}
// 		for _, cur := range currentWitnesses {
// 			canSee := 0

// 			for _, next := range Witnesses[i+1].events {
// 				visitedFamous = make([]*EventNode, 0)
// 				if seeFamous(next, cur) {
// 					canSee++
// 				}
// 				if canSee > 2*variables.N/3 {
// 					cur.Famous = true
// 				}
// 			}

// 		}
// 	}
// 	showWitnesses()
// }

// var visitedFamous []*EventNode

// func seeFamous(x *EventNode, y *EventNode) bool {
// 	if x.OwnHash == y.OwnHash {
// 		return true
// 	}
// 	if x.Round < y.Round {
// 		return false
// 	}

// 	var a, b bool = false, false
// 	if x.ParentEvent != nil && !isVisitedFamous(x.ParentEvent) {
// 		visitedFamous = append(visitedFamous, x.ParentEvent)
// 		a = seeFamous(x.ParentEvent, y)
// 	}
// 	if x.PreviousEvent != nil && !isVisitedFamous(x.PreviousEvent) {
// 		visitedFamous = append(visitedFamous, x.PreviousEvent)
// 		b = seeFamous(x.PreviousEvent, y)
// 	}
// 	return a || b
// }

// func isVisitedFamous(x *EventNode) bool {
// 	for _, v := range visitedFamous {
// 		if x.OwnHash == v.OwnHash {
// 			return true
// 		}
// 	}
// 	return false
// }
