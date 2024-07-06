package process

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"jy.org/thumbgen/src/files"
	"jy.org/thumbgen/src/process/cli"
)

var ffcfg = cfg.Ffmpeg

type Video struct {
    Path string
    duration float64
}

func NewVideo(path string) *Video {
    return &Video{
        Path: path,
    }
}

func (vid *Video) GenPreviewGif(outFile string) error {
    logger.INFO.Println("Generating preview gif for", vid.Path, "to", outFile)

    // make tmp dir
    tmpDir, err := files.MkTmpDir(vid.Path)
    defer os.RemoveAll(tmpDir)
    if err != nil {
        return err
    }

    // get duration
    duration, err := vid.getDuration()
    if err != nil {
        return err
    }

    // get cut start points
    starts := getCuts(duration)
    if len(starts) == 0 {
        return errors.New(fmt.Sprintf("No cuts could be made for %v, duration: %v", vid.Path, duration))
    }

    // get gifs for each cut, save in tmp dir
    for i, start := range starts {
        gif := path.Join(tmpDir, "range" + strconv.Itoa(i) + ".gif")
        err := cli.GenGif(vid.Path, gif, start)
        if err != nil {
            return err
        }
    }

    // combine gifs
    err = cli.CombineGifs(tmpDir, outFile)
    if err != nil {
        return err
    }

    logger.INFO.Printf("Generated preview gif for %v\n", vid.Path)
    return nil
}

func (vid *Video) GenPreviewImg(outFile string) error {
    logger.INFO.Println("Generating preview img for", vid.Path, "to", outFile)

    // generate img at half duration
    dur, err := vid.getDuration()
    if err != nil {
        return err
    }
    snapTime := dur / 2
    return cli.GetVidFrame(vid.Path, outFile, snapTime)
}

func (vid *Video) getDuration() (float64, error) {
    // already got duration
    if vid.duration != 0 {
        return vid.duration, nil
    }

    // get duration
    dur, err := cli.GetVidDuration(vid.Path)
    if err != nil {
        return 0, err
    }

    vid.duration = dur
    return dur, nil
}

