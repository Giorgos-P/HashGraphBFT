package app

import (
	"HashGraphBFT/variables"
	"fmt"

	"time"
)

var ClientChannel = make(chan struct {
	clientTransaction string
	transactionNumber int
	clientID          int
}, 10000)

var NumOfTransactions int
var EmptyTransaction int

func ClientTransactionCreation() {

	NumOfTransactions = 8
	//EmptyTransaction = 4

}

func ClientTransactionCreation1() {

	//num := 1

	frequency := 2
	NumOfTransactions = 2

	character := rune('a' + variables.ID)
	//transaction := fmt.Sprintf("%c", character)

	EmptyTransaction = 4

	for i := 1; i <= NumOfTransactions; i++ {
		transaction := fmt.Sprintf("%d%c", i, character)

		ClientChannel <- struct {
			clientTransaction string
			transactionNumber int
			clientID          int
		}{clientTransaction: transaction, transactionNumber: i, clientID: 0}

		fmt.Println(transaction, " ", i, " to: ", 0)

		time.Sleep(time.Duration(frequency) * time.Second)
	}

	for i := NumOfTransactions + 1; i <= NumOfTransactions+EmptyTransaction; i++ {
		transaction := fmt.Sprintf("empty%d", variables.ID)

		ClientChannel <- struct {
			clientTransaction string
			transactionNumber int
			clientID          int
		}{clientTransaction: transaction, transactionNumber: i, clientID: 0}

		fmt.Println(transaction, " ", i, " to: ", 0)

		time.Sleep(time.Duration(frequency) * time.Second)
	}

}

func ClientTransactionCreation2() {
	frequency := 2

	NumOfTransactions = 2
	Nodes := 4

	total1 := NumOfTransactions * Nodes

	EmptyTransaction = 4
	total2 := EmptyTransaction * Nodes

	character := rune('a' + variables.ID)

	for i := 1; i <= total1; i += 4 {

		for k := 0; k < Nodes; k++ {
			if k != variables.ID {
				continue
			}
			num := i + k
			transaction := fmt.Sprintf("%d%c", i, character)

			//transaction := fmt.Sprintf("%da", num)
			fmt.Println(transaction, " ", num, " to: ", k)
			ClientChannel <- struct {
				clientTransaction string
				transactionNumber int
				clientID          int
			}{clientTransaction: transaction, transactionNumber: num, clientID: 0}
		}
		time.Sleep(time.Duration(frequency) * time.Second)

	}

	for i := total1 + 1; i < total1+1+total2; i += 4 {

		for k := 0; k < Nodes; k++ {
			if k != variables.ID {
				continue
			}
			num := i + k
			transaction := fmt.Sprintf("empty%d", num)
			fmt.Println(transaction, " ", num, " to: ", k)
			ClientChannel <- struct {
				clientTransaction string
				transactionNumber int
				clientID          int
			}{clientTransaction: transaction, transactionNumber: num, clientID: 0}
		}
		time.Sleep(time.Duration(frequency) * time.Second)

	}

}
