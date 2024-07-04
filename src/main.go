package main

import (
	"jy.org/thumbgen/src/config"
	"jy.org/thumbgen/src/files"
	"jy.org/thumbgen/src/logging"
	"jy.org/thumbgen/src/process"
)

var logger = logging.Logger

func main() {
    defer logger.INFO.Println("End of the program")
    logger.INFO.Println("Start of the program")

    cfg := config.Config
    logger.INFO.Printf("Config: %+v\n", cfg)

    err := files.WalkAndDo(cfg.Dirs.Input, 
    func(file string) {
        logger.INFO.Printf("[Generation start] source: %v\n", file) // TODO no log when skipped

        // get target dir
        tdir, err := files.GetTargetDir(file, true)
        if err != nil {
            return
        }

        // generate gif
        vid := process.NewVideo(file, tdir)
        err = vid.GenPreviewGif()
        if err != nil {
            logger.ERROR.Printf("[Generation end][ERROR] %v\n", err)
            return
        }

        // generate img
        err = vid.GenPreviewImg()
        if err != nil {
            logger.ERROR.Printf("[Generation end][ERROR] %v\n", err)
            return
        }

        logger.INFO.Printf("[Generation end][ok]")
    }, 
    func(dir string) {
        logger.INFO.Printf("[Generation start] source: %v/\n", dir)

        // get target dir
        tdir, err := files.GetTargetDir(dir, true)
        if err != nil {
            return
        }

        dirp, err := process.NewDirProcessor(dir, tdir)
        if err != nil {
            logger.ERROR.Printf("[Generation end][ERROR] Error when creating directory processor: %v\n", err)
            return
        }
        err = dirp.GenPreviewGif()
        if err != nil {
            logger.ERROR.Printf("[Generation end][ERROR] Error when generating gif for directory: %v\n", err)
            return
        }

        err = dirp.GenPreviewImg()
        if err != nil {
            logger.ERROR.Printf("[Generation end][ERROR] Error when generating png for directory: %v\n", err)
            return
        }

        logger.INFO.Printf("[Generation end][ok]")
    }, 
    func(dir string) error {
        logger.INFO.Printf("Processing directory: %v/\n", dir)
        _, err := files.MkTargetDir(dir)
        if err != nil {
            logger.ERROR.Printf("Error when creating target directory: %v\n", err)
            return err
        }
        return nil
    })
    if err != nil {
        logger.ERROR.Println("Error when walking through directory: ", err)
    }
}

