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
)

// InitializeLogger - Initializes the Out and Err loggers
func InitializeLogger() {
	outFolder := "/home/giorgos/logs/" // vasilis changes depending on the PC
	errFolder := "/home/giorgos/logs/"
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

	outFile, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	errFile, err := os.OpenFile(errorFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
}
