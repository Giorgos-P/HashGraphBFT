package config

import (
	"HashGraphBFT/variables"
	"strconv"
)

var (
	// RepAddresses - Initialize the address of local Rep sockets
	RepAddresses map[int]string

	// ReqAddresses - Initialize the address of local Req sockets
	ReqAddresses map[int]string

	// ServerAddresses - Initialize the address of local Server sockets
	ServerAddresses map[int]string

	// ResponseAddresses - Initialize the address of local Response sockets
	ResponseAddresses map[int]string
)

// InitializeLocal - Initializes system locally.
func InitializeLocal() {
	RepAddresses = make(map[int]string, variables.N)
	ReqAddresses = make(map[int]string, variables.N)
	ServerAddresses = make(map[int]string, variables.Clients)
	ResponseAddresses = make(map[int]string, variables.Clients)

	for i := 0; i < variables.N; i++ {
		RepAddresses[i] = "tcp://*:" + strconv.Itoa(4000+variables.ID*100+i)
		ReqAddresses[i] = "tcp://localhost:" + strconv.Itoa(4000+i*100+variables.ID)
	}
	for i := 0; i < variables.Clients; i++ {
		ServerAddresses[i] = "tcp://*:" + strconv.Itoa(7000+variables.ID*100+i)
		ResponseAddresses[i] = "tcp://*:" + strconv.Itoa(10000+variables.ID*100+i)
	}

}

// GetRepAddressLocal - Returns the local REP address of the server with that id
func GetRepAddressLocal(id int) string {
	return RepAddresses[id]
}

// GetReqAddressLocal - Returns the local REQ address of the server with that id
func GetReqAddressLocal(id int) string {
	return ReqAddresses[id]
}

// GetServerAddressLocal - Returns the local Server address of the server with that id
func GetServerAddressLocal(id int) string {
	return ServerAddresses[id]
}

// GetResponseAddressLocal - Returns the local Response address of the server with that id
func GetResponseAddressLocal(id int) string {
	return ResponseAddresses[id]
}
