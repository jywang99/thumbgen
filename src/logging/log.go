package logging

import (
	"io"
	"log"
	"os"

	"jy.org/thumbgen/src/config"
)

type logger struct {
    ERROR *log.Logger
    WARN *log.Logger
    INFO *log.Logger
}

func initLoggers() *logger {
    logFile, err := os.OpenFile(config.Config.Log.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }

    flags := log.LstdFlags | log.Lshortfile
    return &logger{
        ERROR: log.New(io.MultiWriter(os.Stderr, logFile), "ERROR:", flags),
        WARN: log.New(io.MultiWriter(os.Stdout, logFile), "WARN:", flags),
        INFO: log.New(io.MultiWriter(os.Stdout, logFile), "INFO:", flags),
    }
}

var Logger = initLoggers()

