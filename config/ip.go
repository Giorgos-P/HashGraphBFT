package config

import (
	"HashGraphBFT/variables"
	"strconv"
)

var addresses = []string{
	"192.36.94.2",
	"141.22.213.35",
	"139.30.241.191",
	"132.227.123.14",
	"129.242.19.196",
	"141.24.249.131",
	"130.192.157.138",
	"141.22.213.34",
	"192.33.193.18",
	"192.33.193.16",
	"131.246.19.201",
	"155.185.54.249",
	"128.232.103.202",
	"195.251.248.180",
	"194.42.17.164",
	"128.232.103.201",
	"193.1.201.27",
	"193.226.19.30",
	"132.65.240.103",
	"193.1.201.26",
	"129.16.20.70",
	"129.16.20.71",
	"195.113.161.13",
}

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

	for i := 0; i < variables.N; i++ {
		RepAddressesIP[i] = "tcp://*:" + strconv.Itoa(4000+i)
		ReqAddressesIP[i] = "tcp://" + addresses[i] + ":" + strconv.Itoa(4000+i)
	}
	for i := 0; i < variables.Clients; i++ {
		ServerAddressesIP[i] = "tcp://*:" + strconv.Itoa(7000+variables.ID*100+i)
		ResponseAddressesIP[i] = "tcp://*:" + strconv.Itoa(10000+variables.ID*100+i)
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
