package ffmpeg

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// example: ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 /mnt/f/aaa/bbb.mp4
func (vid *FfVideo) getDuration() (float64, error) {
    // already got duration
    if vid.duration != 0 {
        return vid.duration, nil
    }

    logger.INFO.Println("Getting duration for", vid.Path)
    cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", vid.Path)
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

    // cache
    vid.duration = res
    return res, nil
}

// example: fmpeg -i /mnt/f/aaa/bbb.mp4 -ss 00:00:00 -t 5 -vf "fps=20,scale=320:-1:flags=lanczos" -c:v pam -f image2pipe - | convert -delay 2 -loop 0 - range1.gif
func (vid *FfVideo) genGif(output string, start float64) error {
    logger.INFO.Printf("Generating gif to %v starting at %v\n", output, start)

    // commands
    ffCmd := exec.Command(
        "ffmpeg", 
        "-i", vid.Path,
        "-ss", strconv.FormatFloat(start, 'f', 0, 64), 
        "-t", strconv.FormatFloat(vid.cfg.CutDuration, 'f', 0, 64), 
        "-vf", "fps=" + strconv.Itoa(vid.cfg.Fps) + ",scale=" + strconv.Itoa(vid.cfg.ScaleWidth) + ":" + strconv.Itoa(vid.cfg.ScaleHeight) + ":flags=lanczos", 
        "-c:v", "pam", "-f", "image2pipe", "-",
    )
    mgkCmd := exec.Command("magick", "-delay", "2", "-loop", "0", "-", output)

    // pipe together and execute
    err := execPipeCmd(ffCmd, mgkCmd)
    if err != nil {
        logger.ERROR.Printf("genGif command error. Executed command:\n\t%v | %v\n", ffCmd.String(), mgkCmd.String())
        return err
    }

    return nil
}

// example: convert -delay 5 -loop 0 range1.gif range2.gif combined.gif
func (vid *FfVideo) combineGifs(gifs []string, outGif string) error {
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

// example: ffmpeg -i input_video.mp4 -ss 00:00:10 -vframes 1 -q:v 2 -vf "scale=640:-1" output_image.png
func (vid *FfVideo) genImg(outPng string, time float64) error {
    logger.INFO.Println("Generating img for", vid.Path, "to", outPng)
    cmd := exec.Command("ffmpeg", 
        "-i", vid.Path, 
        "-ss", strconv.FormatFloat(time, 'f', 0, 64), 
        "-vframes", "1", 
        "-q:v", "2", 
        "-vf", fmt.Sprintf("scale=%d:-1", vid.cfg.ScaleWidth), 
        outPng)
    _, err := execCmd(cmd)
    if err != nil {
        logger.ERROR.Println("genImg command error")
        return err
    }
    return nil
}

