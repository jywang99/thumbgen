package main

import (
	"path"
	"path/filepath"

	"jy.org/thumbgen/src/files"
	"jy.org/thumbgen/src/process"
)

func doForParentDir(dir string) error {
    logger.INFO.Printf("Processing directory: %v/\n", dir)
    _, err := files.MkTargetDir(dir)
    if err != nil {
        logger.ERROR.Printf("Error when creating target directory: %v\n", err)
        return err
    }
    return nil
}

func doForLeaf(file string, isDir bool) {
    logger.INFO.Printf("[Generation start] source: %v\n", file)

    // get relative path
    rel, err := filepath.Rel(cfg.Dirs.Input, file)
    if err != nil {
        logger.ERROR.Printf("[Generation end][ERROR] Error when getting relative path: %v\n", err)
        return
    }

    // write to index file
    err = idx.WriteLine(rel)
    if err != nil {
        logger.ERROR.Printf("[Generation end][ERROR] Error when writing to index file: %v\n", err)
        return
    }

    // get target dir
    tdir, err := files.GetTargetDir(file, true)
    if err != nil {
        logger.ERROR.Printf("Error when getting target dir: %v\n", err)
        return
    }

    // check if targets already exists
    outBase := path.Join(tdir, files.GetBaseName(file, false))
    outGif := outBase + ".gif" // example: /mnt/f/aaa/bbb.mp4 -> /mnt/g/aaa/bbb.gif
    vExist := files.CheckFileExists(outGif)
    outImg := outBase + ".png" // example: /mnt/f/aaa/bbb.mp4 -> /mnt/g/aaa/bbb.png
    iExist := files.CheckFileExists(outImg)
    if vExist && iExist {
        logger.INFO.Printf("[Generation end] Already exist: %v, %v\n", outGif, outImg)
        return
    }

    // process the found file or directory
    if isDir {
        err = processDir(file, outGif, outImg)
    } else {
        err = processFile(file, outGif, outImg)
    }
    if err != nil {
        logger.ERROR.Printf("[Generation end][ERROR] Error processing: %v\n", err)
        return
    }

    logger.INFO.Printf("[Generation end][ok]")
}

func processFile(file, outGif, outImg string) error {
    // generate gif
    vid := process.NewVideo(file)
    err := vid.GenPreviewGif(outGif)
    if err != nil {
        logger.ERROR.Printf("Error when generating gif for %v: %v\n", file, err)
        return err
    }

    // generate img
    err = vid.GenPreviewImg(outImg)
    if err != nil {
        logger.ERROR.Printf("Error when generating img for %v: %v\n", file, err)
        return err
    }

    return nil
}

func processDir(dir, outGif, outImg string) error {
    dirp, err := process.NewDirProcessor(dir)
    if err != nil {
        logger.ERROR.Printf("Error when creating directory processor: %v\n", err)
        return err
    }

    err = dirp.GenPreviewGif(outGif)
    if err != nil {
        logger.ERROR.Printf("Error when generating gif for directory: %v\n", err)
        return err
    }

    err = dirp.GenPreviewImg(outImg)
    if err != nil {
        logger.ERROR.Printf("Error when generating png for directory: %v\n", err)
        return err
    }

    return nil
}
