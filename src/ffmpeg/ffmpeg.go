package ffmpeg

import (
	"os"
	"path"
	"strconv"

	"jy.org/videop/src/config"
	"jy.org/videop/src/files"
)

type Ffmpeg struct {
    cfg config.FfmpegCfg
}

func NewFfmpeg(cfg config.FfmpegCfg) *Ffmpeg {
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
        cfg: &ff.cfg,
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
    tmpDir, err := mkTmpDir(vid.Path)
    defer os.RemoveAll(tmpDir)
    if err != nil {
        return err
    }

    // get duration
    duration, err := vid.getDuration()
    if err != nil {
        return err
    }

    // get gifs for each cut, save in tmp dir
    starts := getCuts(duration, vid.cfg.CutDuration, vid.cfg.MaxCuts)
    gifs := make([]string, len(starts))
    for i, start := range starts {
        gif := path.Join(tmpDir, "range" + strconv.Itoa(i) + ".gif")
        err := vid.genGif(gif, start)
        if err != nil {
            return err
        }
        gifs[i] = gif
    }

    // combine gifs
    err = vid.combineGifs(gifs, outFile)
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
    dur, err := vid.getDuration()
    if err != nil {
        return err
    }
    snapTime := dur / 2
    vid.genImg(outFile, snapTime)

    return nil
}

