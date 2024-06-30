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

    err := files.WalkAndDo(cfg.Dirs.Input, func(path string) error {
        tpath, err := files.GetTargetFile(path)
        if err != nil {
            return err
        }
        return ff.GenPreviewGif(path, tpath)
    }, files.MkTargetDir)
    if err != nil {
        logger.ERROR.Println("Error when walking through directory: ", err)
    }
}

