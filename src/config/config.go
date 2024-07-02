package config

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"jy.org/videop/src/logging"
)

var logger = logging.Logger

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
}

type FfmpegCfg struct {
    PlaybackSpeed float64 `yaml:"playbackSpeed"`
    CutDuration float64 `yaml:"cutDuration"`
    MaxCuts int `yaml:"maxCuts"`
    ScaleWidth int `yaml:"scaleWidth"`
    ScaleHeight int `yaml:"scaleHeight"`
    Fps int `yaml:"fps"`
}

type config struct {
    Ffmpeg FfmpegCfg `yaml:"ffmpeg"`
    Dirs directories `yaml:"directories"`
    Files files `yaml:"files"`
}

var basePath = "/soft/video-prep/config/"
var configPath = path.Join(basePath, "config.yml")

func readYmlConfig(cfg *config) {
    f, err := os.Open(configPath)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    s, _ := f.Stat()
    if s.Size() == 0 {
        return
    }

    decoder := yaml.NewDecoder(f)
    err = decoder.Decode(&cfg)
    if err != nil {
        log.Fatal(err)
    }
}

func initConfig() config {
    err := godotenv.Load(path.Join(basePath, ".env"))
    if err != nil {
        logger.ERROR.Println("Error loading .env file")
    }

    var cfg config
    readYmlConfig(&cfg)
    cfg.Dirs.IgnoreMap = stringToMap(cfg.Dirs.IgnoreStr)
    cfg.Files.VideoExtMap = stringToMap(cfg.Files.VideoExtStr)

    logger.INFO.Printf("%+v\n", cfg)
    return cfg
}
var Config = initConfig()

func stringToMap(s string) map[string]bool {
    m := make(map[string]bool)
    if s == "" {
        return m
    }

    for _, v := range strings.Split(s, ":") {
        m[v] = true
    }
    return m
}

