package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"jy.org/thumbgen/src/config"
)

type logger struct {
    ERROR *log.Logger
    WARN *log.Logger
    INFO *log.Logger
}

func isValidLogfile(path string) bool {
    // check if parent exists
    parent := filepath.Dir(path)
    if _, err := os.Stat(parent); err != nil {
        return false
    }
    return true
}

var Logger = &logger{
    ERROR: log.New(os.Stderr, "ERROR:", log.LstdFlags|log.Lshortfile),
    WARN: log.New(os.Stdout, "WARN:", log.LstdFlags|log.Lshortfile),
    INFO: log.New(os.Stdout, "INFO:", log.LstdFlags|log.Lshortfile),
}

func InitLogFiles() {
    logPath := config.Config.Log.LogPath
    var logFile *os.File
    var err error
    if isValidLogfile(logPath) {
        logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
            logFile = nil
        }
    }

    if logFile == nil {
        Logger.ERROR.Println("Invalid log file path")
        return
    }

    flags := log.LstdFlags | log.Lshortfile
    Logger.ERROR = log.New(io.MultiWriter(os.Stderr, logFile), "ERROR:", flags)
    Logger.WARN = log.New(io.MultiWriter(os.Stdout, logFile), "WARN:", flags)
    Logger.INFO = log.New(io.MultiWriter(os.Stdout, logFile), "INFO:", flags)
}

