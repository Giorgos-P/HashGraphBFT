package main

import (
	"HashGraphBFT/app"
	"HashGraphBFT/config"
	"HashGraphBFT/logger"
	"HashGraphBFT/types"

	//	"HashGraphBFT/types"

	"HashGraphBFT/variables"

	"fmt"

	"log"
	"os"
	"os/signal"
	"strconv"
)

// Initializer - Method that initializes all required processes
func initializer(id int, n int, t int, clients int, scenario config.Scenario) {
	variables.Initialize(id, n, t, clients)
	config.InitializeLocal()
	config.InitializeIP()
	config.InitializeScenario(scenario)

	logger.InitializeLogger()
	logger.OutLogger.Println(
		"ID:", variables.ID,
		"| N:", variables.N,
		"| F:", variables.F,
		"| T:", variables.T,
		"| Clients:", variables.Clients,
	)

	app.InitializeMessenger()
	// app.InitializeReplication()

	// Initialize the message transmition and handling for the servers
	app.Subscribe()
	go app.TransmitMessages()
	// go app.ByzantineReplication()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	go func() {

		for range terminate {
			for i := 0; i < n; i++ {
				if i == id {
					continue
				}
				app.ReceiveSockets[i].Close()
				app.SendSockets[i].Close()
			}
			os.Exit(0)
		}
	}()
}

func main() {
	done := make(chan interface{})

	args := os.Args[1:]
	if len(args) < 5 {
		log.Fatal("Arguments should be '<id> <n> <f> <k> <scenario>")
	}

	id, _ := strconv.Atoi(args[0])
	n, _ := strconv.Atoi(args[1])
	t, _ := strconv.Atoi(args[2])
	clients, _ := strconv.Atoi(args[3])
	tmp, _ := strconv.Atoi(args[4])
	scenario := config.Scenario(tmp)

	initializer(id, n, t, clients, scenario)
	app.Subscribe()
	go app.TransmitMessages()

	//	fmt.Println(config.GetServerAddressLocal(id))
	//	fmt.Println(config.GetServerAddressLocal(id))

	s := "abcd"

	//message := new(types.Message)

	//	message := types.NewMessage([]byte(s), "ReplicaStructure")
	message := types.NewMessage([]byte(s), "a")

	fmt.Println(message)

	if id == 2 {
		message.Type = "b"
		message.Hash = 200
		app.SendMessage(message, 1)
	}

	_ = <-done
}
