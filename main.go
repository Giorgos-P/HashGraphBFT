package main

import (
	"HashGraphBFT/app"
	"HashGraphBFT/config"
	"HashGraphBFT/logger"
	"HashGraphBFT/threshenc"

	//"HashGraphBFT/types"

	"HashGraphBFT/variables"

	"os"
	"os/signal"
	"strconv"
)

// Initializer - Method that initializes all required processes
func initializer(id int, n int, clients int, scenario int, rem int) {
	variables.Initialize(id, n, clients, rem)
	if rem == 0 {
		config.InitializeLocal()
	} else {
		config.InitializeIP()
	}
	config.InitializeScenario(scenario)

	logger.InitializeLogger()
	logger.OutLogger.Println(
		"ID:", variables.ID,
		"| N:", variables.N,
		"| F:", variables.F,
		"| T:", variables.T,
		"| Clients:", variables.Clients,
	)
	threshenc.ReadKeys("./keys/")

	app.InitializeMessenger()

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
	args := os.Args[1:]

	if len(args) == 2 && string(args[0]) == "generate_keys" {
		N, _ := strconv.Atoi(args[1])
		threshenc.GenerateKeys(N, "./keys/")
		return
	}

	// if len(args) < 4 {
	// 	log.Fatal("Arguments should be '<ID> <N> <Clients> <Scenario>")
	// }

	done := make(chan interface{})

	id, _ := strconv.Atoi(args[0])
	n, _ := strconv.Atoi(args[1])
	clients, _ := strconv.Atoi(args[2])
	scenario, _ := strconv.Atoi(args[3])
	rem, _ := strconv.Atoi(args[4])

	initializer(id, n, clients, scenario, rem)

	//os.Exit(0)

	// Initialize the message transmition and handling for the servers
	app.Subscribe()
	go app.TransmitMessages()

	//app.LoadGraph()

	app.InitGraph()
	app.StartHashGraph()

	_ = <-done
}
