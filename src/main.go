package main

import (
	"jy.org/videop/src/config"
	"jy.org/videop/src/ffmpeg"
	"jy.org/videop/src/files"
	"jy.org/videop/src/logging"
)

var logger = logging.Logger

func main() {
    defer logger.INFO.Println("End of the program")
    logger.INFO.Println("Start of the program")

    cfg := config.Config
    ff := ffmpeg.NewFfmpeg()

    err := files.WalkAndDo(cfg.Dirs.Input, 
    func(path string) {
        logger.INFO.Printf("[Generation start] source: %v\n", path)
        tpath, err := files.GetTargetFile(path)
        if err != nil {
            return
        }
        err = ff.GenPreviewGif(path, tpath)
        if err != nil {
            logger.ERROR.Printf("[Generation end][ERROR] %v\n", err)
            return
        }
        logger.INFO.Printf("[Generation end][ok] output: %v\n", tpath)
    }, 
    func(path string) error {
        err := files.MkTargetDir(path)
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

