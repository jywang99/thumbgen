package main

import (
	"jy.org/thumbgen/src/config"
	"jy.org/thumbgen/src/files"
	"jy.org/thumbgen/src/logging"
)

var cfg = config.Config
var logger = logging.Logger
var idx *files.Indexer

func main() {
    logger.INFO.Println("Start of the program")
    defer logger.INFO.Println("End of the program")

    // read config
    logger.INFO.Printf("Config: %+v\n", cfg)

    // create index file
    var err error
    idx, err = files.NewIndexer(cfg.Files.Index)
    if err != nil {
        logger.ERROR.Printf("Error when creating indexer: %v\n", err)
        return
    }
    defer idx.Close()
    logger.INFO.Printf("Index file created at %v\n", cfg.Files.Index)

    err = files.WalkAndDo(cfg.Dirs.Input, doForLeaf, doForParentDir)
    if err != nil {
        logger.ERROR.Println("Error when walking through directory: ", err)
    }
}

