package main

import (
	"path"

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
    ff := ffmpeg.NewFfmpeg(cfg.Ffmpeg)

    err := files.WalkAndDo(cfg.Dirs.Input, 
    func(file string) {
        logger.INFO.Printf("[Generation start] source: %v\n", file)
        tdir, err := files.MkTargetDir(file)
        if err != nil {
            return
        }
        err = ff.GenPreviewGif(file, path.Join(tdir, "preview.gif"))
        if err != nil {
            logger.ERROR.Printf("[Generation end][ERROR] %v\n", err)
            return
        }
        // TODO generate img
        logger.INFO.Printf("[Generation end][ok] output: %v\n", tdir)
    }, 
    func(dir string) {
        logger.INFO.Printf("[Generation start] source: %v/\n", dir)
        tpath, err := files.MkTargetDir(dir)
        if err != nil {
            return
        }
        // TODO generate gif
        // TODO generate img
        logger.INFO.Printf("[Generation end][ok] output: %v\n", tpath)
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

