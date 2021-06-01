package app

import (
	"HashGraphBFT/config"
	"HashGraphBFT/logger"
	"HashGraphBFT/types"
	"HashGraphBFT/variables"
	"bytes"
	"encoding/gob"

	"github.com/pebbe/zmq4"
)

var (
	// Context to initialize sockets
	Context *zmq4.Context

	// SendSockets - Send messages to other servers
	SendSockets map[int]*zmq4.Socket

	// ReceiveSockets - Receive messages from other servers
	ReceiveSockets map[int]*zmq4.Socket

	// ServerSockets - Get the client requests
	ServerSockets map[int]*zmq4.Socket

	// ResponseSockets - Send responses to clients
	ResponseSockets map[int]*zmq4.Socket

	// MessageChannel - Channel to put the messages that need to be transmitted in
	MessageChannel = make(chan struct {
		Message types.Message
		To      int
	}, 10000)

	// EventChannel - Channel to put the received event messages(gossip) in
	EventChannel = make(chan struct {
		Ev   *types.EventMessage
		From int
	}, 10000)

	count = 0
)

// InitializeMessenger - Initializes the 0MQ sockets ( between Servers and Clients)
func InitializeMessenger() {
	Context, err := zmq4.NewContext()
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	// Sockets REP & PUB to communicate with each one of the clients
	ServerSockets = make(map[int]*zmq4.Socket, variables.Clients)
	ResponseSockets = make(map[int]*zmq4.Socket, variables.Clients)
	for i := 0; i < variables.Clients; i++ {

		// ServerSockets initialization to get clients requests
		ServerSockets[i], err = Context.NewSocket(zmq4.REP)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var serverAddr string
		if !variables.Remote {
			serverAddr = config.GetServerAddressLocal(i)
		} else {
			serverAddr = config.GetServerAddress(i)
		}
		err = ServerSockets[i].Bind(serverAddr)
		for err != nil {
			//logger.ErrLogger.Fatal(err)
			err = ServerSockets[i].Bind(serverAddr)
		}
		logger.OutLogger.Println("Requests from Client", i, "on", serverAddr)

		// ResponseSockets initialization to publish the response back to the clients.
		ResponseSockets[i], err = Context.NewSocket(zmq4.PUB)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var responseAddr string
		if !variables.Remote {
			responseAddr = config.GetResponseAddressLocal(i)
		} else {
			responseAddr = config.GetResponseAddress(i)
		}
		err = ResponseSockets[i].Bind(responseAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("Response to Client", i, "on", responseAddr)
	}

	// A socket pair (REP/REQ) to communicate with each one of the other servers
	ReceiveSockets = make(map[int]*zmq4.Socket)
	SendSockets = make(map[int]*zmq4.Socket)
	for i := 0; i < variables.N; i++ {
		// Not myself
		if i == variables.ID {
			continue
		}

		// ReceiveSockets initialization to get information from other servers
		ReceiveSockets[i], err = Context.NewSocket(zmq4.REP)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var receiveAddr string
		if !variables.Remote {
			receiveAddr = config.GetRepAddressLocal(i)
		} else {
			receiveAddr = config.GetRepAddress(i)
		}
		err = ReceiveSockets[i].Bind(receiveAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("Receive from Server", i, "on", receiveAddr)

		// SendSockets initialization to send information to other servers
		SendSockets[i], err = Context.NewSocket(zmq4.REQ)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var sendAddr string
		if !variables.Remote {
			sendAddr = config.GetReqAddressLocal(i)
		} else {
			sendAddr = config.GetReqAddress(i)
		}
		err = SendSockets[i].Connect(sendAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("Send to Server", i, "on", sendAddr)
	}
}

// Broadcast - Broadcasts a message to other servers
func Broadcast(message types.Message) {
	for i := 0; i < variables.N; i++ {
		// Not myself
		if i == variables.ID {
			continue
		}
		SendMessage(message, i)
	}
}

// func SendMessageClient(event *types.ClientMsg, to int) {
// 	w := new(bytes.Buffer)
// 	encoder := gob.NewEncoder(w)
// 	err := encoder.Encode(event)
// 	if err != nil {
// 		logger.ErrLogger.Fatal(err)
// 	}
// 	message := types.NewMessage(w.Bytes(), "ClientMessage")

// 	SendMessage(message, to)
// }

// SendEvent -
func SendEvent(event *types.EventMessage, to int) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(event)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	message := types.NewMessage(w.Bytes(), "Event")

	SendMessage(message, to)
}

// SendMessage - Puts the messages in the message channel to be transmitted with TransmitMessages
func SendMessage(message types.Message, to int) {

	MessageChannel <- struct {
		Message types.Message
		To      int
	}{Message: message, To: to}

}

// TransmitMessages - Transmites the messages to the other servers [go started from main]
func TransmitMessages() {

	for messageTo := range MessageChannel {

		to := messageTo.To
		message := messageTo.Message
		w := new(bytes.Buffer)
		encoder := gob.NewEncoder(w)
		err := encoder.Encode(message)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		_, err = SendSockets[to].SendBytes(w.Bytes(), 0)

		for err != nil {
			_, err = SendSockets[to].SendBytes(w.Bytes(), 0)

			//logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("SENT Message to", to)

		_, err = SendSockets[to].Recv(0)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("OKAY from", to)
		//time.Sleep(time.Second * 1)
	}

}

// Subscribe - Handles the inputs from both clients and other servers
func Subscribe() {
	//fmt.Println("Clients:", variables.Clients)
	// Gets requests from clients and handles them
	for i := 0; i < variables.Clients; i++ {
		go func(i int) { // Initialize them with a goroutine and waits forever
			for {
				message, err := ServerSockets[i].RecvBytes(0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
				logger.OutLogger.Println("Request Received")

				go handleClientRequest(message, i)

				_, err = ServerSockets[i].Send("", 0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
			}
		}(i)
	}

	// Gets messages from other servers and handles them
	for i := 0; i < variables.N; i++ {
		// Not myself
		if i == variables.ID {
			continue
		}
		go func(i int) { // Initializes them with a goroutine and waits forever
			for {
				message, err := ReceiveSockets[i].RecvBytes(0)
				//for err != nil {
				//	message, err = ReceiveSockets[i].RecvBytes(0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}

				// i = who sent as the message
				go handleMessage(message)

				_, err = ReceiveSockets[i].Send("OK.", 0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
			}
		}(i)
	}
}

// Put client's message in RequestChannel
func handleClientRequest(msg []byte, from int) {
	logger.OutLogger.Println("RECEIVED REQ from", from)

	message := new(types.ClientMsg)
	buffer := bytes.NewBuffer([]byte(msg))
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&message)

	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	ClientChannel <- struct {
		clientTransaction string
		transactionNumber int
		clientID          int
	}{clientTransaction: message.ClientTransaction, transactionNumber: message.TransactionNumber, clientID: from}

}

// Handles the messages from the other servers (i think only ReplicaStructure concern us)
func handleMessage(msg []byte) {
	count++
	logger.OutLogger.Println("Message Count:", count)

	message := new(types.Message)
	buffer := bytes.NewBuffer([]byte(msg))
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&message)

	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	switch message.Type {
	// case "ReplicaStructure":
	// 	replica := new(types.ReplicaStructure)
	// 	buf := bytes.NewBuffer(message.Payload)
	// 	dec := gob.NewDecoder(buf)
	// 	err = dec.Decode(&replica)
	// 	if err != nil {
	// 		logger.ErrLogger.Fatal(err)
	// 	}
	// 	ReplicaChannel <- struct {
	// 		Rep  *types.ReplicaStructure
	// 		From int
	// 	}{Rep: replica, From: message.From}
	case "Event":
		replica := new(types.EventMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&replica)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		EventChannel <- struct {
			Ev   *types.EventMessage
			From int
		}{Ev: replica, From: message.From}
	}
}

// ReplyClient - Sends back a response to the client
func ReplyClient(reply types.Reply, to int) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(reply)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	_, err = ResponseSockets[to].SendBytes(w.Bytes(), 0)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	logger.OutLogger.Println("REPLIED Client", to, "-", reply.Value)
}
