package util

import (
	"encoding/json"
	"log"
	"os"

	"github.com/DuGlaser/atc/internal"
)

var infoLogger = log.New(os.Stderr, "INFO: ", log.LstdFlags)
var errorLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags)

func JsonLog(v ...interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		errorLogger.Println(err)
	}

	return string(data)
}

func InfoLog(v ...interface{}) {
	if internal.Verbose {
		infoLogger.Println(v...)
	}
}
