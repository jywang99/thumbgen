package main

import (
	"log"

	"jy.org/thumbgen/src/config"
	"jy.org/thumbgen/src/files"
	"jy.org/thumbgen/src/logging"
)

var cfg = config.Config
var logger = logging.Logger
var idx *files.Indexer

func main() {
    // read config
    args := parseArgs()
    config.Override(args.config)
    err := config.Validate()
    if err != nil {
        log.Fatal(err)
        return
    }

    // init loggers
    logging.InitLogFiles()
    logger.INFO.Println("Starting thumbgen")
    defer logger.INFO.Println("Exiting thumbgen")
    logger.INFO.Printf("Config: %+v\n", cfg)

    // create index file
    idx, err = files.NewIndexer(cfg.Files.Index)
    if err != nil {
        logger.ERROR.Printf("Error when creating indexer: %v\n", err)
        return
    }
    defer idx.Close() // TODO handle SIGTERM
    logger.INFO.Printf("Index file created at %v\n", cfg.Files.Index)

    err = files.WalkAndDo(cfg.Dirs.Input, doForLeaf, doForParentDir)
    if err != nil {
        logger.ERROR.Println("Error when walking through directory: ", err)
    }
}

