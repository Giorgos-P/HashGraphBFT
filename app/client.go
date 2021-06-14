package app

var ClientChannel = make(chan struct {
	clientTransaction string
	transactionNumber int
	clientID          int
}, 10000)
