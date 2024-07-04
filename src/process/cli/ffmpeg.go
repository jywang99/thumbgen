package cli

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// example: ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 /mnt/f/aaa/bbb.mp4
func GetVidDuration(path string) (float64, error) {
    logger.INFO.Println("Getting duration for", path)
    cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", path)
    out, err := execCmd(cmd)
    if err != nil {
        return 0, err
    }

    // convert result to float
    res, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
    if err != nil {
        logger.ERROR.Printf("Error when parsing duration: %v\n", err)
        return 0, err
    }

    return res, nil
}

// example: fmpeg -i /mnt/f/aaa/bbb.mp4 -ss 00:00:00 -t 5 -vf "fps=20,scale=320:-1:flags=lanczos" -c:v pam -f image2pipe - | convert -delay 2 -loop 0 - range1.gif
func GenGif(input string, output string, start float64) error {
    logger.INFO.Printf("Generating gif to %v starting at %v\n", output, start)

    // commands
    ffCmd := exec.Command(
        "ffmpeg", 
        "-i", input,
        "-ss", strconv.FormatFloat(start, 'f', 0, 64), 
        "-t", strconv.FormatFloat(cfg.CutDuration, 'f', 0, 64), 
        "-vf", "fps=" + strconv.Itoa(cfg.Fps) + ",scale=-1:" + strconv.Itoa(cfg.ScaleHeight) + ":flags=lanczos", 
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

// example: ffmpeg -i input_video.mp4 -ss 00:00:10 -vframes 1 -q:v 2 -vf "scale=640:-1" output_image.png
func GetVidFrame(input, output string, time float64) error {
    logger.INFO.Println("Generating img for", input, "to", output)
    cmd := exec.Command("ffmpeg", 
        "-i", input, 
        "-ss", strconv.FormatFloat(time, 'f', 0, 64), 
        "-vframes", "1", 
        "-q:v", "2", 
        "-vf", fmt.Sprintf("scale=%d:-1", cfg.ScaleWidth), 
        output)
    _, err := execCmd(cmd)
    if err != nil {
        logger.ERROR.Println("genImg command error")
        return err
    }
    return nil
}

