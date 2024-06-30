package ffmpeg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type Ffmpeg struct {
    playbackSpeed float64
    cutDuration float64
    maxCuts int
    scaleWidth int
    scaleHeight int
    fps int
}

func NewFfmpeg() *Ffmpeg {
    return &Ffmpeg{
        playbackSpeed: 2.0,
        cutDuration: 5,
        maxCuts: 5,
        scaleWidth: 320,
        scaleHeight: -1,
        fps: 20,
    }
}

func (f *Ffmpeg) GenPreviewGif(inFile string, outFile string) error {
    if _, err := os.Stat(outFile); err == nil {
        return nil
    }

    logger.INFO.Println("Generating preview gif for", inFile, "to", outFile)

    // make tmp dir
    tmpDir, err := mkTmpDir(inFile)
    defer os.RemoveAll(tmpDir)
    if err != nil {
        return err
    }

    // get duration
    duration, err := f.getDuration(inFile)
    if err != nil {
        return err
    }

    // get gifs for each cut, save in tmp dir
    starts := getCuts(duration, f.cutDuration, f.maxCuts)
    gifs := make([]string, len(starts))
    for i, start := range starts {
        gif := path.Join(tmpDir, "range" + strconv.Itoa(i) + ".gif")
        err := f.genGif(inFile, gif, start)
        if err != nil {
            return err
        }
        gifs[i] = gif
    }

    // combine gifs
    err = f.combineGifs(gifs, outFile)
    if err != nil {
        return err
    }

    return nil
}

// example: ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 /mnt/f/MetalBondage/mb218.mp4
func (f *Ffmpeg) getDuration(file string) (float64, error) {
    logger.INFO.Println("Getting duration for", file)
    cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", file)
    out, err := execCmd(cmd)
    if err != nil {
        return 0, err
    }

    // convert result to float
    res, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
    if err != nil {
        log.Printf("getDuration conversion error: %v\n", err)
        return 0, err
    }

    return res, nil
}

// example: fmpeg -i /mnt/f/MetalBondage/mb218.mp4 -ss 00:00:00 -t 5 -vf "fps=20,scale=320:-1:flags=lanczos" -c:v pam -f image2pipe - | convert -delay 2 -loop 0 - range1.gif
func (f *Ffmpeg) genGif(input string, output string, start float64) error {
    logger.INFO.Println("Generating gif for", input, "to", output, "starting at", start)
    cmd := exec.Command("sh", "-c",
        strings.Join([]string{
            "ffmpeg", 
            "-i", fmt.Sprintf("'%s'", input),
            "-ss", strconv.FormatFloat(start, 'f', 0, 64), 
            "-t", strconv.FormatFloat(f.cutDuration, 'f', 0, 64), 
            "-vf", "fps=" + strconv.Itoa(f.fps) + ",scale=" + strconv.Itoa(f.scaleWidth) + ":" + strconv.Itoa(f.scaleHeight) + ":flags=lanczos", 
            "-c:v", "pam", "-f", "image2pipe", "-",
            "|", "magick", "-delay", "2", "-loop", "0", "-", fmt.Sprintf("'%s'", output)}, " "))
    _, err := execCmd(cmd)
    if err != nil {
        log.Println("genGif command error")
        return err
    }
    return nil
}

// example: convert -delay 5 -loop 0 range1.gif range2.gif combined.gif
func (f *Ffmpeg) combineGifs(gifs []string, outGif string) error {
    logger.INFO.Println("Combining gifs to", outGif)
    args := []string{"convert", "-delay", "5", "-loop", "0"}
    for _, gif := range gifs {
        args = append(args, gif)
    }
    args = append(args, outGif)

    cmd := exec.Command("magick", args...)
    _, err := execCmd(cmd)
    if err != nil {
        log.Println("combineGifs command error")
        return err
    }

    return nil
}

