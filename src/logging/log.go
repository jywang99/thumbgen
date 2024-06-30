package logging

import (
	"log"
	"os"
)

type logger struct {
    ERROR *log.Logger
    WARN *log.Logger
    INFO *log.Logger
}

func initLoggers() *logger {
    return &logger{
        ERROR: log.New(os.Stderr, "ERROR:", log.LstdFlags),
        WARN: log.New(os.Stdout, "WARN:", log.LstdFlags),
        INFO: log.New(os.Stdout, "INFO:", log.LstdFlags),
    }
}

var Logger = initLoggers()

