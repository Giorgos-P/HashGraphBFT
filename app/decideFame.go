package app

import (
	"HashGraphBFT/variables"
)

func DecideFame() {

	maxRoundWitness := maxRound()

	round := 1
	round = max(witnessChangeRound-1, 1)

	if !witnessChange {
		return
	}
	witnessChange = false

	for {

		if round >= maxRoundWitness-2 {
			witnessChangeRound = maxRoundWitness - 2
			witnessChange = true
			return
		}

		currentWitnesses, okCur := Witnesses[round]
		nextWitnesses, okNext := Witnesses[round+1]

		if !okCur || !okNext {
			return
		}

		for _, current := range currentWitnesses {

			count := 0
			for _, next := range nextWitnesses {
				visited = make([]*EventNode, 0)
				if see(next, current, FAME) {
					count++
				}
			}

			limit := variables.T
			if count >= limit {
				current.Famous = true
			}
		}

		round++
	}

}
