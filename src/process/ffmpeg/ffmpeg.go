package ffmpeg

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"jy.org/thumbgen/src/config"
	"jy.org/thumbgen/src/files"
)

type Ffmpeg struct {
    cfg *config.FfmpegCfg
}

func NewFfmpeg(cfg *config.FfmpegCfg) *Ffmpeg {
    return &Ffmpeg{cfg: cfg}
}

type FfVideo struct {
    Path string
    TargetDir string
    duration float64
    cfg *config.FfmpegCfg
}

func (ff Ffmpeg) NewFfVideo(fpath, tdir string) *FfVideo {
    return &FfVideo{
        Path: fpath,
        TargetDir: tdir,
        cfg: ff.cfg,
    }
}

func (vid *FfVideo) GenPreviewGif() error {
    // example: /mnt/f/aaa/bbb.mp4 -> /mnt/f/aaa/bbb.gif
    outFile := path.Join(vid.TargetDir, files.GetBaseName(vid.Path, false) + ".gif")
    if _, err := os.Stat(outFile); err == nil {
        logger.INFO.Printf("Preview gif already exists for %v: %v\n", vid.Path, outFile)
        return nil
    }

    logger.INFO.Println("Generating preview gif for", vid.Path, "to", outFile)

    // make tmp dir
    tmpDir, err := files.MkTmpDir(vid.Path)
    defer os.RemoveAll(tmpDir)
    if err != nil {
        return err
    }

    // get duration
    duration, err := vid.GetDuration()
    if err != nil {
        return err
    }

    // get cut start points
    starts := getCuts(duration, vid.cfg.CutDuration, vid.cfg.MaxCuts)
    if len(starts) == 0 {
        return errors.New(fmt.Sprintf("No cuts could be made for %v, duration: %v", vid.Path, duration))
    }

    // get gifs for each cut, save in tmp dir
    gifs := make([]string, len(starts))
    for i, start := range starts {
        gif := path.Join(tmpDir, "range" + strconv.Itoa(i) + ".gif")
        err := vid.GenGif(gif, start)
        if err != nil {
            return err
        }
        gifs[i] = gif
    }

    // combine gifs
    err = CombineGifs(gifs, outFile)
    if err != nil {
        return err
    }

    logger.INFO.Printf("Generated preview gif for %v\n", vid.Path)
    return nil
}

func (vid *FfVideo) GenPreviewImg() error {
    // example: /mnt/f/aaa/bbb.mp4 -> /mnt/f/aaa/bbb.png
    outFile := path.Join(vid.TargetDir, files.GetBaseName(vid.Path, false) + ".png")
    if _, err := os.Stat(outFile); err == nil {
        logger.INFO.Printf("Preview png already exists for %v: %v\n", vid.Path, outFile)
        return nil
    }

    logger.INFO.Println("Generating preview img for", vid.Path, "to", outFile)

    // generate img at half duration
    dur, err := vid.GetDuration()
    if err != nil {
        return err
    }
    snapTime := dur / 2
    vid.GenImg(outFile, snapTime)

    return nil
}

// TODO move to outside package
