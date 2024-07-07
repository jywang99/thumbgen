package config

import (
	"errors"
	"os"
	"path/filepath"
)

type FfmpegCfg struct {
    PlaybackSpeed float64 `yaml:"playbackSpeed"`
    CutDuration float64 `yaml:"cutDuration"`
    MaxCuts int `yaml:"maxCuts"`
    ScaleWidth int `yaml:"scaleWidth"`
    ScaleHeight int `yaml:"scaleHeight"`
    Fps int `yaml:"fps"`
}

type directories struct {
    Input string `yaml:"input"`
    Output string `yaml:"output"`
    Temp string `yaml:"temp"`
    IgnoreStr string `yaml:"ignore"`
    IgnoreMap map[string]bool
    MaxDepth int `yaml:"maxDepth"`
}

type files struct {
    VideoExtStr string `yaml:"videoExt"`
    VideoExtMap map[string]bool
    ImageExtStr string `yaml:"imageExt"`
    ImageExtMap map[string]bool
    DotFiles bool `yaml:"dotfiles"`
    Index string `yaml:"index"`
}

type logs struct {
    LogPath string `yaml:"file"`
}

type config struct {
    Ffmpeg FfmpegCfg `yaml:"ffmpeg"`
    Dirs directories `yaml:"directories"`
    Files files `yaml:"files"`
    Log logs `yaml:"logging"`
}

var Config = &config{
    Ffmpeg: FfmpegCfg{
        PlaybackSpeed: 1.0,
        CutDuration: 3,
        MaxCuts: 5,
        ScaleWidth: 320,
        ScaleHeight: 240,
        Fps: 20,
    },
    Dirs: directories{
        Temp: "/tmp",
        MaxDepth: 1,
    },
    Files: files{
        VideoExtMap: map[string]bool{ "mp4": true, "mkv": true, "avi": true, "mov": true, "wmv": true, "webm": true },
        ImageExtMap: map[string]bool{ "jpg": true, "jpeg": true, "png": true, "gif": true, "webp": true },
        DotFiles: false,
    },
    Log: logs{},
}

func Validate() error {
    dirs := Config.Dirs
    if !dirExists(dirs.Input) {
        return errors.New("Input directory does not exist")
    }
    if !dirExists(dirs.Output) {
        return errors.New("Output directory does not exist")
    }
    if !dirExists(dirs.Temp) {
        return errors.New("Temp directory does not exist")
    }

    files := Config.Files
    if len(files.VideoExtMap) == 0 && len(files.ImageExtMap) == 0 {
        return errors.New("No video or image extensions")
    }
    if !parentDirExists(files.Index) {
        return errors.New("Index file does not exist")
    }

    logs := Config.Log
    if logs.LogPath != "" && !parentDirExists(logs.LogPath) {
        return errors.New("Invalid log file path")
    }

    return nil
}

func parentDirExists(path string) bool {
    parent := filepath.Dir(path)
    return dirExists(parent)
}

func dirExists(dir string) bool {
    if dir == "" {
        return false
    }
    _, err := os.Stat(dir)
    return err == nil
}

