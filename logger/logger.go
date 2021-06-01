package logger

import (
	"HashGraphBFT/variables"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	// OutLogger - Log the outputs
	OutLogger *log.Logger

	// ErrLogger - Log the errors
	ErrLogger *log.Logger

	HashGraphLogger *log.Logger

	OrderLogger *log.Logger

	InfoLogger *log.Logger

	WitnessLogger *log.Logger
)

// InitializeLogger - Initializes the Out and Err loggers
func InitializeLogger() {
	// outFolder := "/home/giorgos/logs/" // changes depending on the PC
	// errFolder := "/home/giorgos/logs/"
	// logger.InitializeLogger("./logs/out/", "./logs/error/")

	outFolder := "./logs/out/"
	errFolder := "./logs/error/"

	// switch config.TestCase {
	// case config.NORMAL:
	// 	outFolder += "normal/out/"
	// 	errFolder += "normal/err/"
	// 	break
	// case config.STALE_VIEWS:
	// 	outFolder += "stale_views/out/"
	// 	errFolder += "stale_views/err/"
	// 	break
	// case config.STALE_STATES:
	// 	outFolder += "stale_states/out/"
	// 	errFolder += "stale_states/err/"
	// 	break
	// case config.STALE_REQUESTS:
	// 	outFolder += "stale_requests/out/"
	// 	errFolder += "stale_requests/err/"
	// 	break
	// case config.BYZANTINE_PRIM:
	// 	outFolder += "byzantine_prim/out/"
	// 	errFolder += "byzantine_prim/err/"
	// 	break
	// case config.NON_SS:
	// 	outFolder += "non_ss/out/"
	// 	errFolder += "non_ss/err/"
	// }
	outputFilePath := outFolder + "output_" + strconv.Itoa(variables.ID) + "_" +
		time.Now().Format("01-02-2006 15:04:05") + ".txt"
	errorFilePath := errFolder + "error_" + strconv.Itoa(variables.ID) + "_" +
		time.Now().Format("01-02-2006 15:04:05") + ".txt"
	graphFilePath := outFolder + "graph_" + strconv.Itoa(variables.ID) + "_" + ".txt"

	orderFilePath := outFolder + "order_" + strconv.Itoa(variables.ID) + "_" + ".txt"

	infoFilePath := outFolder + "info_" + strconv.Itoa(variables.ID) + "_" + ".txt"

	witnessFilePath := outFolder + "witness_" + strconv.Itoa(variables.ID) + "_" + ".txt"

	outFile, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	errFile, err := os.OpenFile(errorFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	graphFile, err := os.OpenFile(graphFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	orderFile, err := os.OpenFile(orderFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	infoFile, err := os.OpenFile(infoFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	witnessFile, err := os.OpenFile(witnessFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	OutLogger = log.New(
		outFile,
		"INFO:\t",
		log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	ErrLogger = log.New(
		errFile,
		"ERROR:\t",
		log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	HashGraphLogger = log.New(
		graphFile,
		"INFO:\t", 0)
	//logger.HashGraphLogger.Println("")

	OrderLogger = log.New(
		orderFile,
		"INFO:\t", 0)
	//logger.OrderLogger.Println("")

	InfoLogger = log.New(
		infoFile,
		"INFO:\t", 0)
	//logger.OrderLogger.Println("")

	WitnessLogger = log.New(
		witnessFile,
		"INFO:\t", 0)
	//logger.WitnessLogger.Println("")
}
