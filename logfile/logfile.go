package logfile

import (
	"log"
	"os"
)

var errLog error

// FileOpen - logfile
func init() {
	// File- file os that

	// If the file doesn't exist, create it or append to the file
	File, errLog := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if errLog != nil {
		log.Fatal("LOG FILE ERROR: ", errLog)
	}

	defer File.Close()

}
