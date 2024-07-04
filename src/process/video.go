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
    TargetDir string
    duration float64
}

func NewVideo(path string, targetDir string) *Video {
    return &Video{
        Path: path,
        TargetDir: targetDir,
    }
}

func (vid *Video) GenPreviewGif() error {
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
    duration, err := vid.getDuration()
    if err != nil {
        return err
    }

    // get cut start points
    starts := vid.getCuts()
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

func (vid *Video) GenPreviewImg() error {
    // example: /mnt/f/aaa/bbb.mp4 -> /mnt/g/aaa/bbb.png
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

// list of cut start points
// evenly distribute cuts
func (vid *Video) getCuts() []float64 {
    dur, _ := vid.getDuration()
    cuts := make([]float64, 0)
    for i := 0; i < ffcfg.MaxCuts; i++ {
        start := float64(i) * dur / float64(ffcfg.MaxCuts)
        if len(cuts) > 0 && cuts[len(cuts) - 1] + ffcfg.CutDuration > start {
            continue
        }
        cuts = append(cuts, start)
    }
    return cuts
}

