package main

import (
	"HashGraphBFT/app"
	"HashGraphBFT/config"
	"HashGraphBFT/logger"
	"HashGraphBFT/threshenc"

	//	"HashGraphBFT/types"

	"HashGraphBFT/variables"

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

	// if id == 0 {
	// 	//threshenc.GenerateKeys(n, "./keys/")
	// 	eventMessage := new(types.EventMessage)
	// 	eventMessage.Signature = []byte("abc")
	// 	eventMessage.Timestamp = time.Now().UnixNano()
	// 	eventMessage.Transaction = "0"
	// 	eventMessage.PreviousHash = "0"
	// 	eventMessage.ParentHash = "0"
	// 	eventMessage.Owner = 0

	// 	fmt.Println("------------------------------")
	// 	fmt.Println(*eventMessage)
	// 	fmt.Println("------------------------------")

	// 	app.CreateSignature(eventMessage)
	// 	fmt.Println(*eventMessage)
	// 	fmt.Println("------------------------------")

	// 	verif := app.VerifySignature(eventMessage)
	// 	fmt.Println(verif)
	// 	fmt.Println("------------------------------")

	// 	fmt.Println(eventMessage)
	// 	fmt.Println("------------------------------")

	// } else {
	// 	return
	// }

	// Initialize the message transmition and handling for the servers
	app.Subscribe()
	go app.TransmitMessages()

	go app.InitGraph()

	_ = <-done
}
