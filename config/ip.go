package config

import (
	"HashGraphBFT/variables"
	"strconv"
)

var address0 = []string{}

var address1 = []string{}

var (
	// RepAddressesIP - Initialize the address of IP REP sockets
	RepAddressesIP map[int]string

	// ReqAddressesIP - Initialize the address of IP REQ sockets
	ReqAddressesIP map[int]string

	// ServerAddressesIP - Initialize the address of IP Server sockets
	ServerAddressesIP map[int]string

	// ResponseAddressesIP - Initialize the address of IP Response sockets
	ResponseAddressesIP map[int]string
)

// InitializeIP - Initializes system with ips.
func InitializeIP() {
	RepAddressesIP = make(map[int]string, variables.N)
	ReqAddressesIP = make(map[int]string, variables.N)
	ServerAddressesIP = make(map[int]string, variables.Clients)
	ResponseAddressesIP = make(map[int]string, variables.Clients)

	address := address1
	if variables.N == 4 {
		address = address0
	}

	for i := 0; i < variables.N; i++ {
		ad := i
		if i >= len(address) {
			ad = (i % 2) + 2
		}

		RepAddressesIP[i] = "tcp://*:" + strconv.Itoa(27000+i+(variables.ID*variables.N))
		ReqAddressesIP[i] = "tcp://" + address[ad] + ":" + strconv.Itoa(27000+variables.ID+(i*variables.N))
	}
	for i := 0; i < variables.Clients; i++ {
		ServerAddressesIP[i] = "tcp://*:" + strconv.Itoa(27250+i+(variables.ID*variables.Clients))
		ResponseAddressesIP[i] = "tcp://*:" + strconv.Itoa(27625+i+(variables.ID*variables.Clients))
	}
}

// GetRepAddress - Returns the IP REP address of the server with that id
func GetRepAddress(id int) string {
	return RepAddressesIP[id]
}

// GetReqAddress - Returns the IP REQ address of the server with that id
func GetReqAddress(id int) string {
	return ReqAddressesIP[id]
}

// GetServerAddress - Returns the IP Server address of the server with that id
func GetServerAddress(id int) string {
	return ServerAddressesIP[id]
}

// GetResponseAddress - Returns the IP Response address of the server with that id
func GetResponseAddress(id int) string {
	return ResponseAddressesIP[id]
}

// var addresses = []string{
// }

// var (
// 	// RepAddressesIP - Initialize the address of IP REP sockets
// 	RepAddressesIP map[int]string

// 	// ReqAddressesIP - Initialize the address of IP REQ sockets
// 	ReqAddressesIP map[int]string

// 	// ServerAddressesIP - Initialize the address of IP Server sockets
// 	ServerAddressesIP map[int]string

// 	// ResponseAddressesIP - Initialize the address of IP Response sockets
// 	ResponseAddressesIP map[int]string
// )

// // InitializeIP - Initializes system with ips.
// func InitializeIP() {
// 	RepAddressesIP = make(map[int]string, variables.N)
// 	ReqAddressesIP = make(map[int]string, variables.N)
// 	ServerAddressesIP = make(map[int]string, variables.Clients)
// 	ResponseAddressesIP = make(map[int]string, variables.Clients)

// 	for i := 0; i < variables.N; i++ {
// 		RepAddressesIP[i] = "tcp://*:" + strconv.Itoa(4000+i)
// 		ReqAddressesIP[i] = "tcp://" + addresses[i] + ":" + strconv.Itoa(4000+i)
// 	}
// 	for i := 0; i < variables.Clients; i++ {
// 		ServerAddressesIP[i] = "tcp://*:" + strconv.Itoa(7000+variables.ID*100+i)
// 		ResponseAddressesIP[i] = "tcp://*:" + strconv.Itoa(10000+variables.ID*100+i)
// 	}
// }

// // GetRepAddress - Returns the IP REP address of the server with that id
// func GetRepAddress(id int) string {
// 	return RepAddressesIP[id]
// }

// // GetReqAddress - Returns the IP REQ address of the server with that id
// func GetReqAddress(id int) string {
// 	return ReqAddressesIP[id]
// }

// // GetServerAddress - Returns the IP Server address of the server with that id
// func GetServerAddress(id int) string {
// 	return ServerAddressesIP[id]
// }

// // GetResponseAddress - Returns the IP Response address of the server with that id
// func GetResponseAddress(id int) string {
// 	return ResponseAddressesIP[id]
// }
