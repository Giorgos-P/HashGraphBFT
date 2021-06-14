package app

import (
	"HashGraphBFT/logger"
	"HashGraphBFT/variables"
	"time"
)

func sentToAllEverything() bool {
	N := variables.N
	for i := 0; i < N; i++ {
		k := HashGraph[i].events
		if len(k) == 0 {
			continue
		}
		for _, v := range k {
			for syncWith := 0; syncWith < N; syncWith++ {
				if v.Know[syncWith] == false {
					return false
				}
			}
		}
	}
	return true
}

var sentChecked bool = false

var countChecked bool = false

var currentTotal int = -1

var prevTot int = 0

func countAll() {

	NumOfTransactions := 10
	EmptyTransaction := 20

	if countChecked || NumOfTransactions < 2 {
		return
	}

	count := len(sortedEvents)
	// for i := 0; i < variables.N; i++ {
	// 	k := HashGraph[i].events
	// 	count += len(k)
	// }

	tot := variables.Clients*(NumOfTransactions+EmptyTransaction)*variables.N + variables.N

	if currentTotal == -1 {
		logger.InfoLogger.Println(tot)
		logger.InfoLogger.Println("---------------")
		currentTotal = 0
	}

	//testSame()

	if currentTotal != count {
		logger.InfoLogger.Println(count)
		currentTotal = count
	}

	// currentSentToAll := sentToAllEverything()
	// if prevSendtoAll != currentSentToAll {
	// 	logger.InfoLogger.Println(currentSentToAll)
	// 	prevSendtoAll = currentSentToAll
	// }

	if count >= tot && count > prevTot {
		end = count
		// prevTot = count
		// logger.InfoLogger.Println("MUST START ")

		// showHashGraph() // write Hashgraph in Hashgraph log file
		// logger.InfoLogger.Println("Show HashGraph Done ")

		// DivideRounds()
		// DecideFame()
		// Ord()
		// showWitness()

		startTime := time.Now()
		DivideRounds()
		logger.InfoLogger.Println("Divide Rounds Done :", time.Since(startTime).Seconds())

		startTime = time.Now()
		DecideFame()
		logger.InfoLogger.Println("Decide Fame Done :", time.Since(startTime).Seconds())

		startTime = time.Now()
		Ord()
		logger.InfoLogger.Println("Find Order Done :", time.Since(startTime).Seconds())

		showWitness() // write witnesses in info log file
		logger.InfoLogger.Println("Show Witness Done ")

		// countChecked = false
	}

}

func HashGraphSize() int {
	count := 0
	for i := 0; i < variables.N; i++ {
		k := HashGraph[i].events
		count += len(k)
	}
	return count
}

var hashGraphShow bool = false

func showHashGraph() {
	// if hashGraphShow {
	// 	return
	// }
	// hashGraphShow = true
	logger.HashGraphLogger.Println("Display HashGraph in time order")
	logger.HashGraphLogger.Println("# of Elements : ", len(sortedEvents))

	for _, v := range sortedEvents {

		//logger.HashGraphLogger.Print(v.EventMessage.Timestamp, " ", v.EventMessage.Transaction, " ", v.EventMessage.PreviousHash, " ", v.EventMessage.ParentHash, " ", v.EventMessage.Owner, " ", v.OwnHash)
		//logger.HashGraphLogger.Print(v.EventMessage.Transaction, " ", v.EventMessage.Number, " ", v.EventMessage.ClientID, " ", v.Witness, " ", v.Famous)
		logger.HashGraphLogger.Print(v.OwnHash)

	}

	// for _, v := range sortedEvents {
	// 	logger.InfoLogger.Println(v.OwnHash, " ", v.Round)
	// }

}

var next int = 100
var end int = -1
var minLen int = -1

func testSame() bool {

	size := HashGraphSize()

	if size <= next {
		return false
	}

	//start := next - 1100
	//end = next - 500
	//end = size - 500

	end = size - 500

	logger.InfoLogger.Println("------Starts------------------")
	logger.InfoLogger.Println("Size -", size)

	logger.InfoLogger.Println("HashGraph Bounds till: ", end)

	DivideRounds()
	DecideFame()
	Ord()

	// for i := 0; i < end; i++ {
	// 	v := sortedEvents[i]
	// 	//logger.InfoLogger.Print(i, ": ", v.Round, " ", v.Witness, " ", v.Famous, " ", v.OwnHash)
	// 	logger.InfoLogger.Print(i, ": ", v.EventMessage.ClientID, " ", v.EventMessage.Transaction, " ", v.EventMessage.Number, " ", v.EventMessage.Owner)
	// 	//	logger.InfoLogger.Print("\t", v.OwnHash, " ", v.EventMessage.PreviousHash, " ", v.EventMessage.ParentHash)
	// 	if v.Witness {
	// 		logger.InfoLogger.Print("\t", v.Round, " ", v.Witness, " ", v.Famous)
	// 	}
	// }

	//next += 100

	// min1 := min(len(HashGraph[0].events), len(HashGraph[1].events))

	// min2 := max(len(HashGraph[0].events), len(HashGraph[1].events))

	// for i := 2; i < variables.N; i++ {
	// 	mikos := len(HashGraph[i].events)
	// 	if mikos >= min2 {
	// 		continue
	// 	} else if mikos <= min1 {
	// 		min2 = min1
	// 		min1 = mikos
	// 	} else if mikos > min1 && mikos < min2 {
	// 		min2 = mikos
	// 	}

	// }
	// minLen = min2

	// for i := 0; i < variables.N; i++ {
	// 	mikos := len(HashGraph[i].events)
	// 	logger.InfoLogger.Println("--", mikos)

	// }

	logger.InfoLogger.Println("------END------------------")

	//logger.InfoLogger.Println("    --------------------         ")
	return true
}

func executeAlgorithms() bool {

	size := HashGraphSize()
	//next = size + 1
	if size <= next {
		return false
	}

	//start := next - 1100
	//end = next - 500
	//end = size - 500

	end = size - 100
	//end = 2404
	logger.InfoLogger.Println("------Starts------------------")
	logger.InfoLogger.Println("Size -", size)

	logger.InfoLogger.Println("HashGraph Bounds till: ", end)

	// DivideRounds()
	// DecideFame()
	// Ord()

	callAlgos()

	logger.InfoLogger.Println("------END------------------")
	next += 100
	//logger.InfoLogger.Println("    --------------------         ")
	return true
}

func showWitness() {

	logger.WitnessLogger.Println("Witnesses")

	for i := 1; ; i++ {
		list, exist := Witnesses[i]
		if !exist {
			return
		}
		SortSlice(list)
		var owners []bool
		owners = make([]bool, variables.N)

		for _, v := range list {
			if i != v.Round {
				logger.WitnessLogger.Println(" Different witness Round with ROUND")
			}
			if !v.Witness {
				logger.WitnessLogger.Println("Witness must be true")
			}
			if owners[v.EventMessage.Owner] {
				logger.WitnessLogger.Println("Owner - ", v.EventMessage.Owner, " already has a witness")
			}

			logger.WitnessLogger.Println(i, " ", v.Round, v.EventMessage.Transaction, " ", v.EventMessage.Owner, " ", v.Witness, " ", v.Famous, " ", v.OwnHash, v.EventMessage.Timestamp)
			owners[v.EventMessage.Owner] = true
		}

	}

}

func callAlgos() {
	startTime := time.Now()
	DivideRounds()
	logger.InfoLogger.Println("Divide Rounds Done :", time.Since(startTime).Seconds())

	startTime = time.Now()
	DecideFame()
	logger.InfoLogger.Println("Decide Fame Done :", time.Since(startTime).Seconds())

	startTime = time.Now()
	Ord()
	logger.InfoLogger.Println("Find Order Done :", time.Since(startTime).Seconds())

	showWitness() // write witnesses in info log file
	logger.InfoLogger.Println("Show Witness Done ")
}

func showOrdering() {
	logger.OrderLogger.Println("Ordered Events: ", len(OrderedEvents))
	for _, v := range OrderedEvents {
		//logger.OrderLogger.Println(v.EventMessage.Transaction, " ", v.Round, " ", v.Witness, " ", v.Famous, " ", v.RoundReceived, " ", v.ConsensusTime)
		logger.OrderLogger.Println(v.RoundReceived, " ", v.ConsensusTime, " ", v.EventMessage.Transaction, " ", v.EventMessage.Number, " ", v.EventMessage.ClientID)
	}

}
